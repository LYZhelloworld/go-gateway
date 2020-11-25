package gateway

import (
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

// ErrorConfig is a map that matches status codes to Handler.
type ErrorConfig map[int]Handler

// routerConfig holds service for different methods, with method string as the key
type routerConfig map[string]serviceInfo

// serviceInfo contains the name and handler of a service.
type serviceInfo struct {
	// name is the name of a service.
	name string
	// handler is the Handler of a service.
	handler Handler
}
