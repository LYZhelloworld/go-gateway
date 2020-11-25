package gateway

// serviceInfo contains the name and handler of a Service.
type serviceInfo struct {
	// name is the name of a Service.
	name string
	// handler is the Handler of a Service.
	handler Handler
}
