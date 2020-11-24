package gateway

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
