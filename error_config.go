package gateway

// SetErrorHandler sets error handler of the given HTTP status code.
func (s *Server) SetErrorHandler(status int, handler Handler) {
	checkNonNilHandler(handler)
	s.errorConfig[status] = handler
}

// RemoveErrorHandler removes error handler of the given HTTP status code.
func (s *Server) RemoveErrorHandler(status int) {
	delete(s.errorConfig, status)
}
