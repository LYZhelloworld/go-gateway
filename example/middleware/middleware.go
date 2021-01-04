package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LYZhelloworld/go-gateway"
	"github.com/LYZhelloworld/go-gateway/middleware"
)

// This example shows how the middleware works.
//
// The middleware 1 will be executed, and then middleware 2 (before Next()), middleware 3, main handler,
// then the middleware 2 (after Next()).
//
// If you change `context.Next()` to `context.Interrupt()`, the execution will stop at "middleware 2 end" and
// the following handlers will not be run.
func main() {
	s := gateway.Default()
	cfg := gateway.Config{}
	cfg.Add("/hello", http.MethodGet, "api.gateway.hello")
	s.UseConfig(cfg)
	s.Register("api.gateway.hello", func(context *gateway.Context) {
		context.Logger.Info("body start")
		context.Response = []byte("hello, world")
		context.Logger.Info("body end")
	})
	s.UseMiddlewares(middleware.Logger(), func(context *gateway.Context) {
		context.Logger.Info("middleware 1")
	}, func(context *gateway.Context) {
		context.Logger.Info("middleware 2 start")
		context.Next() // Change this to `context.Interrupt()` to get different result.
		context.Logger.Info("middleware 2 end")
	}, func(context *gateway.Context) {
		context.Logger.Info("middleware 3")
	})
	if err := s.RunWithShutdown(":8080", 5*time.Second); err != nil {
		fmt.Println(err)
	}
}
