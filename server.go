package gateway

import (
	"fmt"
	"net/http"
)

// Server is a struct for HTTP router.
type Server struct {
	// Config is the configuration mapping endpoints and methods to services.
	Config Config
	// ErrorConfig is a map that matches status codes to ServiceHandler.
	ErrorConfig ErrorConfig
	// Services are the collection of all services.
	Services []*Service
	// Preprocessors are the collection of handlers executed before the main handler.
	// The order of the execution follows the order of every handler in the collection.
	Preprocessors []ServiceHandler
	// Postprocessors are the collection of handlers executed after the main handler.
	// The order of the execution follows the order of every handler in the collection.
	Postprocessors []ServiceHandler

	// endpointConfig is a map with endpoint as key and routerConfig as value.
	endpointConfig EndpointConfig
}

// Default creates a default Server without any configurations.
func Default() *Server {
	return &Server{
		Config:      Config{},
		ErrorConfig: ErrorConfig{},
	}
}

// AddEndpoint links an endpoint to a service by name.
func (s *Server) AddEndpoint(path string, method string, service ServiceName) {
	if s.Config == nil {
		s.Config = Config{}
	}
	s.Config[Endpoint{Path: path, Method: method}] = service
}

// AddService registers a service.
func (s *Server) AddService(name ServiceName, handler ServiceHandler) {
	if handler == nil {
		panic("nil handler")
	}
	s.Services = append(s.Services, &Service{Name: name, Handler: handler})
}

// AddErrorConfig registers an ErrorConfig with a specific status code.
func (s *Server) AddErrorConfig(status int, handler ServiceHandler) {
	if s.ErrorConfig == nil {
		s.ErrorConfig = ErrorConfig{}
	}
	s.ErrorConfig[status] = handler
}

// AddPreprocessor registers a preprocessor to the Server.
func (s *Server) AddPreprocessor(handler ServiceHandler) {
	s.Preprocessors = append(s.Preprocessors, handler)
}

// AddPreprocessors registers preprocessors to the Server.
func (s *Server) AddPreprocessors(handlers ...ServiceHandler) {
	for _, h := range handlers {
		s.AddPreprocessor(h)
	}
}

// AddPostprocessor registers a postprocessor to the Server.
func (s *Server) AddPostprocessor(handler ServiceHandler) {
	s.Postprocessors = append(s.Postprocessors, handler)
}

// AddPostprocessors registers postprocessors to the Server.
func (s *Server) AddPostprocessors(handlers ...ServiceHandler) {
	for _, h := range handlers {
		s.AddPostprocessor(h)
	}
}

// Run starts the server with the current Config.
func (s *Server) Run(addr string) error {
	if s.Config == nil {
		s.Config = Config{}
	}

	if s.ErrorConfig == nil {
		s.ErrorConfig = ErrorConfig{}
	}

	// parse Services
	s.endpointConfig = EndpointConfig{}
	for endpoint, name := range s.Config {
		service := s.matchService(name)
		if service == nil {
			panic(fmt.Sprintf("service not found: %s", name))
		}
		if s.endpointConfig[endpoint.Path] == nil {
			s.endpointConfig[endpoint.Path] = &routerConfig{}
		}
		if ok := s.endpointConfig[endpoint.Path].setService(endpoint.Method, service); !ok {
			panic(fmt.Sprintf("invalid method: %s", endpoint.Method))
		}
	}

	svr := &http.Server{
		Addr:    addr,
		Handler: s,
	}
	return svr.ListenAndServe()
}

// matchService finds service that is the closest to the given one.
// For example, if the given service is "foo.bar" and there are services like:
//
// foo, foo.bar.baz, foo.baz
//
// The service "foo" will be matched.
// The service "foo.bar.baz" is more specific than the given one.
// The service "foo.baz" has different sub-service "baz".
func (s *Server) matchService(name ServiceName) *Service {
	var found *Service = nil
	distance := int(^uint(0) >> 1) // largest int
	for _, srv := range s.Services {
		ok, d := srv.match(name)
		if ok {
			if d < distance {
				found = srv
				distance = d
			}
		}
	}
	return found
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

	config := s.endpointConfig[req.URL.RawPath]
	if config == nil {
		s.generalResponse(ctx, http.StatusNotFound)
		return
	}

	service, ok := config.getService(req.Method)
	if !ok {
		s.generalResponse(ctx, http.StatusMethodNotAllowed)
		return
	}
	if service == nil {
		s.generalResponse(ctx, http.StatusNotFound)
		return
	}

	ctx.serviceName = service.Name
	s.response(ctx, service.Handler)
	return
}

// preprocess executes preprocessors on the context.
func (s *Server) preprocess(context *Context) {
	for _, h := range s.Preprocessors {
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
	for _, h := range s.Postprocessors {
		h(context)
		if context.isInterrupted {
			return
		}
	}
}

// response generates HTTP response using the handler.
// ServeHTTP must return after calling this method.
func (s *Server) response(context *Context, handler ServiceHandler) {
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
	if handler, ok := s.ErrorConfig[statusCode]; ok {
		handler(context)
		if context.isInterrupted {
			return
		}
	}
	s.postprocess(context)
}
