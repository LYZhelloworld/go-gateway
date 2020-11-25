# Gateway
Gateway is an HTTP server written in Golang.

## Quick Start
```
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LYZhelloworld/gateway"
)

func main() {
	s := gateway.Default()
	s.AddEndpoint("/hello", http.MethodGet, "api.gateway.hello")
	s.AddService("api.gateway.hello", func(context *gateway.Context) {
		context.Response = []byte("hello, world")
	})
	if err := s.RunWithShutdown(":8080", 5*time.Second); err != nil {
		fmt.Println(err)
	}
}
```

## Endpoint and Service
An endpoint is the URL path of the HTTP request. For example: `/hello`.

A service is a handler with a group of identifiers separated by dots (`.`) as its name. For example: `api.gateway.hello`.
Every endpoint points to a service name. A service with exact the same service name can handle the request.

However, a service that is "more generic" than the service name indicated by the endpoint is still capable of handling
the service, only if there is no other services that is "more specific". For example: if there is `api.gateway` but
no `api.gateway.hello`, it can still handle `api.gateway.hello` request.

A service with name `*` will handle all requests if no other service handler exists and matches the service name given.

## Server
`Server.Config` maps endpoints to service names. The endpoint here is a struct of both the path and the method.

`Server.Services` is a collection of all services with their name.

`Server.Run()` starts a server without shutting down procedure.

`Server.RunWithShutdown()` starts a server with shutdown timeout and will shutdown the server gracefully.

## Context
`Context` is the thing that the handler requires when the server is running.

`Context.Request` contains all the information of the request.

`Context.StatusCode` is the status code of the response. It can be changed in the handler.

`Context.Response` is a byte array which contains the response body.

`Context.Header` contains headers of the response.

## Preprocessors and Postprocessors
Preprocessors will be executed before a request, similarly, postprocessors after a request.
They share the same context during the request flow.

You can call `Context.Interrupt()` at any time inside these handlers.
After the handler returns, the following handlers will not be executed, but the response will still be written.
