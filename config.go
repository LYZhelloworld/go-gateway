package gateway

// Config is a map that matches endpoints to Services.
type Config map[Endpoint]ServiceName

// Endpoint is a struct of path and method.
type Endpoint struct {
	Path   string
	Method string
}
