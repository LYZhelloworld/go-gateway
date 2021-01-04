package middleware

import "github.com/LYZhelloworld/go-gateway"

// Logger provides a default middleware for logging requests and responses.
func Logger() gateway.Handler {
	return func(context *gateway.Context) {
		log := context.Logger
		log.WithField("path", context.Request.URL.EscapedPath()).
			WithField("method", context.Request.Method).
			WithField("service", context.GetServiceName()).Info("request")
		context.Next()
		log.WithField("status", context.StatusCode).
			WithField("response_length", len(context.Response)).Info("response")
	}
}
