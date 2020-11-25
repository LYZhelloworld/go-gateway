package gateway

// serviceInfo contains the name and handler of a service.
type serviceInfo struct {
	// name is the name of a service.
	name string
	// handler is the Handler of a service.
	handler Handler
}
