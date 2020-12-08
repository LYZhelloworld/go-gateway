package gateway

import (
	"net/http"

	"github.com/LYZhelloworld/gateway/logger"
)

// Context is the context of a request.
type Context struct {
	// Request is the pointer to the http.Request.
	Request *http.Request
	// StatusCode holds the status code of the response.
	StatusCode int
	// Response holds the response body.
	Response []byte
	// Header holds HTTP headers in the response.
	Header http.Header
	// Data is a map that holds data of any type for value exchange between middlewares and main handler.
	Data map[string]interface{}
	// Logger is the current logger used in this context.
	Logger logger.Logger

	// serviceName is the name of the Service of the request.
	serviceName string
	// responseWriter is the http.ResponseWriter from the handler.
	responseWriter http.ResponseWriter
	// isWritten is a flag shows whether the response has been written to the http.ResponseWriter.
	isWritten bool

	// handlerSeq is a pointer to the handlers going to be run.
	handlerSeq []Handler
	// handlerCounter is a counter of the current handler.
	handlerCounter int
}

// createContext creates an empty Context.
func createContext(w http.ResponseWriter, req *http.Request, server *Server) *Context {
	ctx := &Context{
		Request:        req,
		StatusCode:     http.StatusOK,
		Header:         map[string][]string{},
		Data:           map[string]interface{}{},
		Logger:         server.logger,
		responseWriter: w,
	}
	ctx.handlerSeq = make([]Handler, 0, len(server.middleware)+1)
	for _, m := range server.middleware {
		ctx.handlerSeq = append(ctx.handlerSeq, m)
	}
	return ctx
}

// write writes response to the http.ResponseWriter.
func (c *Context) write() {
	if !c.isWritten {
		c.isWritten = true
		w := c.responseWriter

		for key, values := range c.Header {
			w.Header().Del(key)
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(c.StatusCode)
		_, err := w.Write(c.Response)
		if err != nil {
			panic(err)
		}
	}
}

// run runs all handlers.
func (c *Context) run() {
	for c.handlerCounter = 0; !c.isDone(); {
		c.runCurrentHandler()
	}
}

// runCurrentHandler runs handler based on the handlerCounter.
func (c *Context) runCurrentHandler() {
	if !c.isDone() {
		oldCounter := c.handlerCounter
		c.handlerSeq[c.handlerCounter](c)
		// after running handler, the counter should increase at least once
		// if not, increase it manually
		if c.handlerCounter == oldCounter {
			c.handlerCounter++
		}
	}
}

// isDone checks if the current Context has run all the middlewares or has been interrupted.
func (c *Context) isDone() bool {
	return c.handlerCounter >= len(c.handlerSeq)
}

// Next continues with the next handler, and will return if the following handlers have been run.
func (c *Context) Next() {
	c.handlerCounter++
	for ; !c.isDone(); {
		c.runCurrentHandler()
	}

}

// GetServiceName gets Service name of the request.
func (c *Context) GetServiceName() string {
	return c.serviceName
}

// Interrupt stops the following handlers from executing, but does not stop the current handler.
// This method can be used in either pre-/post-processors or the main handler.
// Calling this method multiple times does not have side effects.
func (c *Context) Interrupt() {
	c.handlerCounter = len(c.handlerSeq)
}
