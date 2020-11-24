package gateway

import (
	"net/http"
)

// Config is a map that matches endpoints to Services.
type Config map[Endpoint]ServiceName

// Endpoint is a struct of path and method.
type Endpoint struct {
	Path   string
	Method string
}

// EndpointConfig is a map that matches string endpoint to routerConfig.
type EndpointConfig map[string]*routerConfig

// ErrorConfig is a map that matches status codes to ServiceHandler.
type ErrorConfig map[int]*ServiceHandler

// routerConfig holds services for different methods.
type routerConfig struct {
	get    *Service
	post   *Service
	put    *Service
	delete *Service
}

// setService assigns Service to the specific method.
func (r *routerConfig) setService(method string, service *Service) (ok bool) {
	switch method {
	case http.MethodGet:
		r.get = service
	case http.MethodPost:
		r.post = service
	case http.MethodPut:
		r.put = service
	case http.MethodDelete:
		r.delete = service
	default:
		return false
	}
	return true
}

// getService gets Service of the specific method.
func (r *routerConfig) getService(method string) (service *Service, ok bool) {
	switch method {
	case http.MethodGet:
		return r.get, true
	case http.MethodPost:
		return r.post, true
	case http.MethodPut:
		return r.put, true
	case http.MethodDelete:
		return r.delete, true
	default:
		return nil, false
	}
}
