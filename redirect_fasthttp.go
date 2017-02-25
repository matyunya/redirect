package main

import (
	"log"

	"github.com/matyunya/redirect/redirect"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

func main() {
	ln, err := reuseport.Listen("tcp4", redirect.ServerHost)
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	if err = fasthttp.Serve(ln, redirect.Handler); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)
	}
}
