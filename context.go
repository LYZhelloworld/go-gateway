package gateway

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Context is the context of a request.
type Context struct {
	// Request is the pointer to the http.Request.
	Request *http.Request
	// StatusCode holds the status code of the response.
	StatusCode int
	// Response holds the response body.
	Response io.ReadCloser
	// Header hold HTTP headers in the response.
	Header http.Header

	// serviceName is the name of the service of the request.
	serviceName ServiceName
	// responseWriter is the http.ResponseWriter from the handler.
	responseWriter http.ResponseWriter
	// isWritten is a flag shows whether the response has been written to the http.ResponseWriter.
	isWritten bool
	// isInterrupted is a flag shows whether the execution of handler chain is interrupted.
	isInterrupted bool
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
		res, err := ioutil.ReadAll(c.Response)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(res)
		if err != nil {
			panic(err)
		}
	}
}

// GetServiceName gets service name of the request.
func (c *Context) GetServiceName() ServiceName {
	return c.serviceName
}

// Interrupt stops the following handlers from executing, but does not stop the current handler.
// This method can be used in either pre-/post-processors or the main handler.
// Calling this method multiple times does not have side effects.
func (c *Context) Interrupt() {
	c.isInterrupted = true
}
