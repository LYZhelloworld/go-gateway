package gateway

import (
	"net/http"
)

// Context is the context of a request.
type Context struct {
	// Request is the pointer to the http.Request.
	Request *http.Request
	// statusCode holds the status code of the response.
	statusCode int
	// response holds the response body.
	response []byte
}

// SetStatusCode sets status code.
func (c *Context) SetStatusCode(code int) {
	c.statusCode = code
}

// SetResponse sets response.
func (c *Context) SetResponse(resp []byte) {
	c.response = resp
}
