package gateway

import "strings"

// endpointConfig is a map that matches string endpoint to routerConfig.
type endpointConfig map[string]*routerConfig

// getEndpointConfig gets corresponding endpoint Config from the Server.
func (e *endpointConfig) get(path string) *routerConfig {
	// remove trailing slash
	path = strings.TrimSuffix(path, "/")
	if config, ok := (*e)[path]; ok {
		return config
	} else {
		// check prefixes
		for p := removeLastDir(path); p != ""; p = removeLastDir(p) {
			if config, ok := (*e)[path+"/*"]; ok {
				return config
			}
		}
		return nil
	}
}

// routerConfig holds Service for different methods, with method string as the key
type routerConfig map[string]serviceInfo
