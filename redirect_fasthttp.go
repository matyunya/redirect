package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	pg "gopkg.in/pg.v5"
)

type Redirect struct {
	ID                 int
	URL, P1            string
	P2, P3, P4, P5, P6 []string
}

var DB *pg.DB

func init() {
	DB = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "mac",
		Database: "redirect",
		Password: os.Getenv("pass"),
	})
}

func main() {
	ln, err := reuseport.Listen("tcp4", "localhost:8069")
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	if err = fasthttp.Serve(ln, Handler); err != nil {
		log.Fatalf("error in fasthttp Server: %s", err)
	}
}

func Handler(ctx *fasthttp.RequestCtx) {
	if strings.Contains(string(ctx.Path()), "redirect") {
		ctx.SetStatusCode(fasthttp.StatusFound)
		return
	}

	args := ctx.QueryArgs()
	pargs, err := url.ParseQuery(args.String())
	if err != nil {
		fmt.Print(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	redirect := Redirect{}
	q := DB.Model(&redirect)
	for param, value := range pargs {
		if param == "p1" {
			q.Where(`p1 = ?`, value[0])
		} else {
			q.Where(param+` @> ?`, pg.Array(value))
		}
	}

	if err := q.Column(`url`).Order(`id desc`).Limit(1).Select(); err != nil {
		fmt.Println(err)
	}

	if len(redirect.URL) > 0 {
		ctx.Redirect(redirect.URL, fasthttp.StatusFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}
