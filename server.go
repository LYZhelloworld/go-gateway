package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYZhelloworld/go-logger"
)

// Server is a struct for HTTP router.
type Server struct {
	// config is the configuration mapping endpoints and methods to service.
	config Config
	// errorConfig is a map that matches status codes to Handler.
	errorConfig map[int]Handler
	// service is a map of all Handler.
	service map[string]Handler
	// middleware is a collection of middlewares executed before/after the main handler.
	// The order of the execution follows the order of every middleware in the collection.
	// Use Context.Next() to continue with the next middleware
	// and it will return after the following middlewares are executed.
	middleware []Handler
	// logger is the logger assigned to the Server.
	logger logger.Logger
	// endpointConfig is a map with endpoint as key and routerConfig as value.
	endpointConfig endpointConfig
}

// Default creates a Server with default configurations.
func Default() *Server {
	return &Server{
		config:      Config{},
		errorConfig: map[int]Handler{},
		service:     map[string]Handler{},
		logger:      logger.GetDefaultLogger(),
	}
}

// prepare sets all configurations before running.
func (s *Server) prepare(addr string) *http.Server {
	if s.config == nil {
		s.config = Config{}
		s.logger.Warn("config is nil. Use empty config instead.")
	}

	if s.errorConfig == nil {
		s.errorConfig = map[int]Handler{}
		s.logger.Warn("errorConfig is nil. Use empty errorConfig instead.")
	}

	// parse service
	s.endpointConfig = endpointConfig{}
	for endpoint, name := range s.config {
		matchedName, handler := s.matchService(name)
		if handler == nil {
			s.logger.WithField("endpoint", endpoint.Path).
				WithField("method", endpoint.Method).
				WithField("service", name).Fatal("handler not found")
			panic(fmt.Sprintf("handler not found: %s", name))
		}
		if s.endpointConfig[endpoint.Path] == nil {
			s.endpointConfig[endpoint.Path] = &routerConfig{}
		}
		(*s.endpointConfig[endpoint.Path])[endpoint.Method] = serviceInfo{name: matchedName, handler: handler}
		s.logger.WithField("endpoint", endpoint.Path).
			WithField("method", endpoint.Method).
			WithField("service", matchedName).
			Info("service matched")
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
	s.logger.Info("start server")
	return svr.ListenAndServe()
}

// RunWithShutdown starts the server with the current Config.
// It catches a SIGINT or SIGTERM as shutdown signal.
func (s *Server) RunWithShutdown(addr string, shutdownTimeout time.Duration) error {
	svr := s.prepare(addr)
	errChan := make(chan error)
	go func() {
		s.logger.Info("start server")
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
			s.logger.Info("shutdown server")
			return ctx.Err()
		case err := <-errChan:
			s.logger.Info("shutdown server")
			return err
		}
	case err := <-errChan:
		s.logger.Info("shutdown server")
		return err
	}
}

// matchService finds Service that is the closest to the given one.
// For example, if the given Service is "foo.bar" and there are Service like:
//
// foo, foo.bar.baz, foo.baz
//
// The Service "foo" will be matched.
// The Service "foo.bar.baz" is more specific than the given one.
// The Service "foo.baz" has different sub-service "baz".
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
	// catch all panics here so that the panics from handlers will not make the server crash
	defer func() {
		if r := recover(); r != nil {
			var log logger.Logger
			if err, ok := r.(error); ok {
				log = s.logger.WithError(err)
			} else {
				log = s.logger.WithField("err", r)
			}
			log.Error("server panic")
		}
	}()

	ctx := createContext(w, req, s)
	path := req.URL.EscapedPath()
	method := req.Method

	config := s.endpointConfig[path]
	if config == nil {
		s.generalResponse(ctx, http.StatusNotFound)
		return
	}

	service, ok := (*config)[method]
	if !ok {
		s.generalResponse(ctx, http.StatusNotFound)
		return
	}

	ctx.serviceName = service.name
	s.response(ctx, service.handler)
	return
}

// response generates HTTP response using the handler.
// ServeHTTP must return after calling this method.
func (s *Server) response(context *Context, handler Handler) {
	defer context.write()

	context.handlerSeq = append(context.handlerSeq, handler)
	context.run()
}

// generalResponse generates error messages depending on the status code.
// ServeHTTP must return after calling this method.
func (s *Server) generalResponse(context *Context, statusCode int) {
	defer context.write()

	context.StatusCode = statusCode
}

// Handler is a function that handles the Service.
type Handler func(context *Context)
