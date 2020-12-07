package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYZhelloworld/gateway/logger"
)

// Server is a struct for HTTP router.
type Server struct {
	// Config is the configuration mapping endpoints and methods to Service.
	Config Config
	// ErrorConfig is a map that matches status codes to Handler.
	ErrorConfig ErrorConfig
	// Service is a map of all Handler.
	Service Service
	// Preprocessors are a collection of middlewares executed before the main handler.
	// The order of the execution follows the order of every middleware in the collection.
	Preprocessors Middleware
	// Postprocessors are a collection of middlewares executed after the main handler.
	// The order of the execution follows the order of every middleware in the collection.
	Postprocessors Middleware
	// Logger is the logger assigned to the Server.
	Logger logger.Logger

	// endpointConfig is a map with endpoint as key and routerConfig as value.
	endpointConfig EndpointConfig
}

// Default creates a default Server without any configurations.
func Default() *Server {
	return &Server{
		Config:         Config{},
		ErrorConfig:    ErrorConfig{},
		Service:        Service{},
		Preprocessors:  Middleware{},
		Postprocessors: Middleware{},
		Logger:         logger.GetDefaultLogger(),
	}
}

// prepare sets all configurations before running.
func (s *Server) prepare(addr string) *http.Server {
	if s.Config == nil {
		s.Config = Config{}
		s.Logger.Warn("Config is nil. Use empty Config instead.")
	}

	if s.ErrorConfig == nil {
		s.ErrorConfig = ErrorConfig{}
		s.Logger.Warn("ErrorConfig is nil. Use empty ErrorConfig instead.")
	}

	// parse Service
	s.endpointConfig = EndpointConfig{}
	for endpoint, name := range s.Config {
		matchedName, handler := s.matchService(name)
		if handler == nil {
			s.Logger.WithField("endpoint", endpoint.Path).
				WithField("method", endpoint.Method).
				WithField("service", name).Fatal("handler not found")
			panic(fmt.Sprintf("handler not found: %s", name))
		}
		if s.endpointConfig[endpoint.Path] == nil {
			s.endpointConfig[endpoint.Path] = &routerConfig{}
		}
		(*s.endpointConfig[endpoint.Path])[endpoint.Method] = serviceInfo{name: matchedName, handler: handler}
		s.Logger.WithField("endpoint", endpoint.Path).
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
	s.Logger.Info("start server")
	return svr.ListenAndServe()
}

// RunWithShutdown starts the server with the current Config.
// It catches a SIGINT or SIGTERM as shutdown signal.
func (s *Server) RunWithShutdown(addr string, shutdownTimeout time.Duration) error {
	svr := s.prepare(addr)
	errChan := make(chan error)
	go func() {
		s.Logger.Info("start server")
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
			s.Logger.Info("shutdown server")
			return ctx.Err()
		case err := <-errChan:
			s.Logger.Info("shutdown server")
			return err
		}
	case err := <-errChan:
		s.Logger.Info("shutdown server")
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
// The Service "foo.baz" has different sub-Service "baz".
func (s *Server) matchService(name string) (string, Handler) {
	for thisName := name; thisName != ""; thisName = removeLastSubService(thisName) {
		if srv, ok := s.Service[thisName]; ok {
			return thisName, srv
		}
	}

	if srv, ok := s.Service[baseServiceHandler]; ok {
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
				log = s.Logger.WithError(err)
			} else {
				log = s.Logger.WithField("err", r)
			}
			log.Error("server panic")
		}
	}()

	ctx := createContext(w, req)
	path := req.URL.EscapedPath()
	method := req.Method
	log := s.Logger.WithField("path", path).WithField("method", method)

	config := s.endpointConfig[path]
	if config == nil {
		s.generalResponse(ctx, http.StatusNotFound)
		log.Warn("not found")
		return
	}

	service, ok := (*config)[method]
	if !ok {
		s.generalResponse(ctx, http.StatusNotFound)
		log.Warn("not found")
		return
	}

	ctx.serviceName = service.name
	log.WithField("service", service.name).Debug("handle http service")
	s.response(ctx, service.handler)
	return
}

// preprocess executes Preprocessors on the context.
func (s *Server) preprocess(context *Context) {
	for _, h := range s.Preprocessors.handlers {
		h(context)
		if context.isInterrupted {
			return
		}
	}
}

// postprocess executes Postprocessors on the context.
func (s *Server) postprocess(context *Context) {
	if context.isInterrupted {
		return
	}
	for _, h := range s.Postprocessors.handlers {
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
	if handler, ok := s.ErrorConfig[statusCode]; ok {
		handler(context)
		if context.isInterrupted {
			return
		}
	}
	s.postprocess(context)
}
