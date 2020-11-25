package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server is a struct for HTTP router.
type Server struct {
	// config is the configuration mapping endpoints and methods to service.
	config Config
	// errorConfig is a map that matches status codes to Handler.
	errorConfig ErrorConfig
	// service is a map of all Handler.
	service Service
	// preprocessors are the collection of handlers executed before the main handler.
	// The order of the execution follows the order of every handler in the collection.
	preprocessors []Handler
	// postprocessors are the collection of handlers executed after the main handler.
	// The order of the execution follows the order of every handler in the collection.
	postprocessors []Handler
	// endpointConfig is a map with endpoint as key and routerConfig as value.
	endpointConfig EndpointConfig
}

// Default creates a default Server without any configurations.
func Default() *Server {
	return &Server{
		config:         Config{},
		errorConfig:    ErrorConfig{},
		endpointConfig: EndpointConfig{},
	}
}

// AddEndpoint links an endpoint to a service by name.
// If the path does not end with a slash and an asterisk ("/*"),
// only requests that match the path EXACTLY will be handled by the service.
// Paths ending with "/*" is considered as a prefix.
//
// For example:
//
// "/api/echo" can be handled by "/api/echo" or "/api/*", but not "/api" or "/".
//
// If multiple prefixes exist, the prefix that matches the most will be the handler.
//
// For example:
//
// "/api/foo/bar" will be handled by "/api/foo/*" but not "/api/*".
//
// The service name of an endpoint should be as specific as possible and should not contain asterisk (*).
func (s *Server) AddEndpoint(path string, method string, service string) {
	if s.config == nil {
		s.config = Config{}
	}
	if path == "" || !isValidPath(trimPrefix(path)) {
		panic("invalid path")
	}
	if service == baseServiceHandler || !isValidService(service) {
		panic("invalid service")
	}
	s.config[Endpoint{Path: path, Method: method}] = service
}

// AddService registers a service.
func (s *Server) AddService(name string, handler Handler) {
	if s.service == nil {
		s.service = Service{}
	}
	if handler == nil {
		panic("nil handler")
	}
	if !isValidService(name) {
		panic("invalid service")
	}
	s.service[name] = handler
}

// AddErrorConfig registers an ErrorConfig with a specific status code.
func (s *Server) AddErrorConfig(status int, handler Handler) {
	if s.errorConfig == nil {
		s.errorConfig = ErrorConfig{}
	}
	s.errorConfig[status] = handler
}

// AddPreprocessor registers a preprocessor to the Server.
func (s *Server) AddPreprocessor(handler Handler) {
	s.preprocessors = append(s.preprocessors, handler)
}

// AddPreprocessors registers preprocessors to the Server.
func (s *Server) AddPreprocessors(handlers ...Handler) {
	for _, h := range handlers {
		s.AddPreprocessor(h)
	}
}

// AddPostprocessor registers a postprocessor to the Server.
func (s *Server) AddPostprocessor(handler Handler) {
	s.postprocessors = append(s.postprocessors, handler)
}

// AddPostprocessors registers postprocessors to the Server.
func (s *Server) AddPostprocessors(handlers ...Handler) {
	for _, h := range handlers {
		s.AddPostprocessor(h)
	}
}

// prepare sets all configurations before running.
func (s *Server) prepare(addr string) *http.Server {
	if s.config == nil {
		s.config = Config{}
	}

	if s.errorConfig == nil {
		s.errorConfig = ErrorConfig{}
	}

	// parse service
	s.endpointConfig = EndpointConfig{}
	for endpoint, name := range s.config {
		matchedName, handler := s.matchService(name)
		if handler == nil {
			panic(fmt.Sprintf("handler not found: %s", name))
		}
		if s.endpointConfig[endpoint.Path] == nil {
			s.endpointConfig[endpoint.Path] = &routerConfig{}
		}
		(*s.endpointConfig[endpoint.Path])[endpoint.Method] = serviceInfo{name: matchedName, handler: handler}
	}

	svr := &http.Server{
		Addr:    addr,
		Handler: s,
	}
	return svr
}

// Run starts the server with the current Config.
func (s *Server) Run(addr string) error {
	svr := s.prepare(addr)
	return svr.ListenAndServe()
}

// RunWithShutdown starts the server with the current Config.
// It catches a SIGINT or SIGTERM as shutdown signal.
func (s *Server) RunWithShutdown(addr string, shutdownTimeout time.Duration) error {
	svr := s.prepare(addr)
	errChan := make(chan error)
	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	quit := make(chan os.Signal)
	// kill: SIGTERM
	// kill -2: SIGINT
	// kill -9: SIGKILL (cannot be caught)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := svr.Shutdown(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		}
	case err := <-errChan:
		return err
	}
}

// matchService finds service that is the closest to the given one.
// For example, if the given service is "foo.bar" and there are service like:
//
// foo, foo.bar.baz, foo.baz
//
// The service "foo" will be matched.
// The service "foo.bar.baz" is more specific than the given one.
// The service "foo.baz" has different sub-service "baz".
func (s *Server) matchService(name string) (string, Handler) {
	for thisName := name; thisName != ""; thisName = removeLastSubService(thisName) {
		if srv, ok := s.service[thisName]; ok {
			return thisName, srv
		}
	}

	if srv, ok := s.service[baseServiceHandler]; ok {
		return baseServiceHandler, srv
	} else {
		return "", nil
	}
}

// ServeHTTP serves HTTP requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:        req,
		StatusCode:     http.StatusOK,
		Response:       nil,
		Header:         map[string][]string{},
		responseWriter: w,
	}

	config := s.endpointConfig[req.URL.EscapedPath()]
	if config == nil {
		s.generalResponse(ctx, http.StatusNotFound)
		return
	}

	service, ok := (*config)[req.Method]
	if !ok {
		s.generalResponse(ctx, http.StatusNotFound)
		return
	}

	ctx.serviceName = service.name
	s.response(ctx, service.handler)
	return
}

// preprocess executes preprocessors on the context.
func (s *Server) preprocess(context *Context) {
	for _, h := range s.preprocessors {
		h(context)
		if context.isInterrupted {
			return
		}
	}
}

// postprocess executes postprocessors on the context.
func (s *Server) postprocess(context *Context) {
	if context.isInterrupted {
		return
	}
	for _, h := range s.postprocessors {
		h(context)
		if context.isInterrupted {
			return
		}
	}
}

// response generates HTTP response using the handler.
// ServeHTTP must return after calling this method.
func (s *Server) response(context *Context, handler Handler) {
	defer context.write()

	s.preprocess(context)
	if context.isInterrupted {
		return
	}
	handler(context)
	if context.isInterrupted {
		return
	}
	s.postprocess(context)
}

// generalResponse generates error messages depending on the status code.
// ServeHTTP must return after calling this method.
func (s *Server) generalResponse(context *Context, statusCode int) {
	defer context.write()

	context.StatusCode = statusCode
	s.preprocess(context)
	if context.isInterrupted {
		return
	}
	if handler, ok := s.errorConfig[statusCode]; ok {
		handler(context)
		if context.isInterrupted {
			return
		}
	}
	s.postprocess(context)
}
