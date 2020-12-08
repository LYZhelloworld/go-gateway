package gateway

// Register registers a service.
//
// Service is the configuration of one Service.
// The key is the identifier of the Service, separated by dots, with parent Service before sub Service.
// For example: "foo.bar.baz".
//
// The value is the handler function of a Service.
//
// A Service can be handled by a more generic Service name (the request of which can be forwarded to other Service).
// For example: "foo.bar" can handle "foo.bar.baz" requests.
// But "foo.bar.baz" cannot handle "foo.bar".
//
// An asterisk (*) means a Service handler for all Service, if there is no other Service that are more specific.
func (s *Server) Register(name string, handler Handler) {
	checkNonNilHandler(handler)
	if !isValidService(name) {
		panic("invalid service")
	}
	s.service[name] = handler
}