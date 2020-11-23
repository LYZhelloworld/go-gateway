package gateway

import (
	"fmt"
	"net/http"
)

// Server is a struct for HTTP router.
type Server struct {
	Config   *Config
	Services []*Service
	router   router
}

// RegisterService registers a service.
func (s *Server) RegisterService(service *Service) *Server {
	s.Services = append(s.Services, service)
	return s
}

// Run starts the server with the current Config.
func (s *Server) Run(addr string) error {
	if s.Config == nil || s.Services == nil {
		panic("empty configuration")
	}

	// parse Services
	s.router = make(router)
	for endpoint, name := range *(s.Config) {
		service := s.matchService(name)
		if service == nil {
			panic(fmt.Sprintf("service not found: %s", name))
		}
		if s.router[endpoint.Path] == nil {
			s.router[endpoint.Path] = &routerConfig{}
		}
		switch endpoint.Method {
		case Get:
			s.router[endpoint.Path].getHandler = service.Handler
		case Post:
			s.router[endpoint.Path].postHandler = service.Handler
		case Put:
			s.router[endpoint.Path].putHandler = service.Handler
		case Delete:
			s.router[endpoint.Path].deleteHandler = service.Handler
		default:
			panic(fmt.Sprintf("invalid method: %s", endpoint.Method))
		}
	}

	svr := &http.Server{
		Addr:    addr,
		Handler: &s.router,
	}
	return svr.ListenAndServe()
}

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
