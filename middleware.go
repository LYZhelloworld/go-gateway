package gateway

// UsePreprocessor registers a handler as the preprocessor.
func (s *Server) UsePreprocessor(handler Handler) {
	checkNonNilHandler(handler)
	s.preprocessors = append(s.preprocessors, handler)
}

// UsePreprocessors registers handlers as the preprocessors.
func (s *Server) UsePreprocessors(handlers ...Handler) {
	for _, h := range handlers {
		s.UsePreprocessor(h)
	}
}

// UsePostprocessor registers a handler as the postprocessor.
func (s *Server) UsePostprocessor(handler Handler) {
	checkNonNilHandler(handler)
	s.postprocessors = append(s.postprocessors, handler)
}

// UsePostprocessors registers handlers as the postprocessors.
func (s *Server) UsePostprocessors(handlers ...Handler) {
	for _, h := range handlers {
		s.UsePostprocessor(h)
	}
}
