package gateway

type Middleware struct {
	// handlers store all the middleware handlers.
	handlers []Handler
}

// Add registers a handler to the Middleware.
func (m *Middleware) Add(handler Handler) {
	m.handlers = append(m.handlers, handler)
}

// AddAll registers handlers to the Middleware.
func (m *Middleware) AddAll(handlers ...Handler) {
	for _, h := range handlers {
		m.Add(h)
	}
}
