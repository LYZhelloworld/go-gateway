package gateway

import (
	"fmt"
	"net/http"
)

// Server is a struct for HTTP router.
type Server struct {
	Config      *Config
	ErrorConfig *ErrorConfig
	Services    []*Service
	router      *router
}

// RegisterService registers a service.
func (s *Server) RegisterService(service *Service) *Server {
	s.Services = append(s.Services, service)
	return s
}

// Run starts the server with the current Config.
func (s *Server) Run(addr string) error {
	if s.Config == nil {
		s.Config = &Config{}
	}

	if s.ErrorConfig == nil {
		s.ErrorConfig = &ErrorConfig{}
	}

	if s.Services == nil {
		s.Services = []*Service{}
	}

	// parse Services
	s.router = &router{endpointConfig: make(EndpointConfig), errorConfig: *s.ErrorConfig}
	for endpoint, name := range *(s.Config) {
		service := s.matchService(name)
		if service == nil {
			panic(fmt.Sprintf("service not found: %s", name))
		}
		if s.router.endpointConfig[endpoint.Path] == nil {
			s.router.endpointConfig[endpoint.Path] = &routerConfig{}
		}
		switch endpoint.Method {
		case http.MethodGet:
			s.router.endpointConfig[endpoint.Path].getHandler = service.Handler
		case http.MethodPost:
			s.router.endpointConfig[endpoint.Path].postHandler = service.Handler
		case http.MethodPut:
			s.router.endpointConfig[endpoint.Path].putHandler = service.Handler
		case http.MethodDelete:
			s.router.endpointConfig[endpoint.Path].deleteHandler = service.Handler
		default:
			panic(fmt.Sprintf("invalid method: %s", endpoint.Method))
		}
	}

	svr := &http.Server{
		Addr:    addr,
		Handler: s.router,
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
