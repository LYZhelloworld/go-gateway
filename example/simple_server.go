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
