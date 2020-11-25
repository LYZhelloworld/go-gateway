package gateway

import (
	"net/http"
	"strings"
)

// Config is a map that matches endpoints to service.
type Config map[Endpoint]string

// Endpoint is a struct of path and method.
type Endpoint struct {
	Path   string
	Method string
}

// EndpointConfig is a map that matches string endpoint to routerConfig.
type EndpointConfig map[string]*routerConfig

// getEndpointConfig gets corresponding endpoint config from the Server.
func (e *EndpointConfig) get(path string) *routerConfig {
	// remove trailing slash
	path = strings.TrimSuffix(path, "/")
	if config, ok := (*e)[path]; ok {
		return config
	} else {
		// check prefixes
		for p := removeLastDir(path); p != ""; p = removeLastDir(p) {
			if config, ok := (*e)[path + "/*"]; ok {
				return config
			}
		}
		return nil
	}
}

// ErrorConfig is a map that matches status codes to ServiceHandler.
type ErrorConfig map[int]ServiceHandler

// routerConfig holds service for different methods.
type routerConfig struct {
	get    *serviceInfo
	post   *serviceInfo
	put    *serviceInfo
	delete *serviceInfo
}

// serviceInfo contains the name and handler of a service.
type serviceInfo struct {
	// name is the name of a service.
	name string
	// handler is the ServiceHandler of a service.
	handler ServiceHandler
}

// setService assigns Service to the specific method.
func (r *routerConfig) setService(method string, service *serviceInfo) (ok bool) {
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
func (r *routerConfig) getService(method string) (service *serviceInfo, ok bool) {
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
