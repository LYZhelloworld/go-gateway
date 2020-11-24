# Gateway
Gateway is an HTTP server written in Golang.

## Quick Start
```
package main

import "github.com/LYZhelloworld/gateway"

func main() {
	s := gateway.Server{
		Config: &gateway.Config{
			gateway.Endpoint{
				Path:   "/echo",
				Method: "GET",
			}: "api.gateway.echo",
		},
		Services: []*gateway.Service{
			{
				Name: "api.gateway.echo",
				Handler: func(context *gateway.Context) {
					context.Response = []byte("hello, world")
				},
			},
		},
	}
	_ = s.Run(":8080")
}
```

## Endpoint and Service
An endpoint is the URL path of the HTTP request. For example: `/api/echo`.

A service is a handler with a group of identifiers separated by dots (`.`) as its name. For example: `api.gateway.echo`.
Every endpoint points to a service name. A service with exact the same service name can handle the request.

However, a service that is "more generic" than the service name indicated by the endpoint is still capable of handling
the service, only if there is no other services that is "more specific". For example: if there is `api.gateway` but
no `api.gateway.echo`, it can still handle `api.gateway.echo` request.

A service with name `*` will handle all requests if no other service handler exists and matches the service name given.

## Server
`Server.Config` maps endpoints to service names. The endpoint here is a struct of both the path and the method.

`Server.Services` is a collection of all services with their name.

By executing `Server.Run()`, these two arguments will be parsed, and the server starts.

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
