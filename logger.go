package gateway

import "github.com/LYZhelloworld/go-logger"

// AttachLogger attaches logger to the Server.
func (s *Server) AttachLogger(logger logger.Logger) {
	s.logger = logger
}
