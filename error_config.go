package gateway

// ErrorConfig is a map that matches status codes to Handler.
type ErrorConfig map[int]Handler

// Add registers a Handler with a specific status code.
func (e *ErrorConfig) Add(status int, handler Handler) {
	(*e)[status] = handler
}
