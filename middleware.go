package gateway

// UseMiddleware registers a middleware.
func (s *Server) UseMiddleware(handler Handler) {
	checkNonNilHandler(handler)
	s.middleware = append(s.middleware, handler)
}

// UseMiddlewares registers middlewares.
func (s *Server) UseMiddlewares(handlers ...Handler) {
	for _, h := range handlers {
		s.UseMiddleware(h)
	}
}
