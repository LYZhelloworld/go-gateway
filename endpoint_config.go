package gateway

import "strings"

// EndpointConfig is a map that matches string endpoint to routerConfig.
type EndpointConfig map[string]*routerConfig

// getEndpointConfig gets corresponding endpoint Config from the Server.
func (e *EndpointConfig) get(path string) *routerConfig {
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
