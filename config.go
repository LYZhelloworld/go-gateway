package gateway

// Config is a map that matches endpoints to Service.
type Config map[Endpoint]string

// Add links an endpoint to a Service by name.
// If the path does not end with a slash and an asterisk ("/*"),
// only requests that match the path EXACTLY will be handled by the Service.
// Paths ending with "/*" is considered as a prefix.
//
// For example:
//
// "/api/echo" can be handled by "/api/echo" or "/api/*", but not "/api" or "/".
//
// If multiple prefixes exist, the prefix that matches the most will be the handler.
//
// For example:
//
// "/api/foo/bar" will be handled by "/api/foo/*" but not "/api/*".
//
// The Service name of an endpoint should be as specific as possible and should not contain asterisk (*).
func (c *Config) Add(path string, method string, service string) {
	if path == "" || !isValidPath(trimPrefix(path)) {
		panic("invalid path")
	}
	if service == baseServiceHandler || !isValidService(service) {
		panic("invalid Service")
	}
	(*c)[Endpoint{Path: path, Method: method}] = service
}

// Get gets service name of the specific path and method.
func (c *Config) Get(path string, method string) string {
	return (*c)[Endpoint{Path: path, Method: method}]
}
