package redirect

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
	pg "gopkg.in/pg.v5"
)

type Redirect struct {
	ID                 int
	URL, P1            string
	P2, P3, P4, P5, P6 []string `pg:",array"`
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

	q.Column(`url`).Order(`id desc`).Limit(1).Select()
	if len(redirect.URL) > 0 {
		ctx.Redirect(redirect.URL, fasthttp.StatusFound)
	} else {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}
