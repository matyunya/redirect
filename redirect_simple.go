package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

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

func handler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "redirect") {
		return
	}

	redirect := Redirect{}
	q := DB.Model(&redirect)

	for k, p := range r.URL.Query() {
		if k == "p1" {
			q.Where(`p1 = ?`, p)
		} else {
			q.Where(k+` @> ?`, pg.Array(p))
		}
	}

	if err := q.Column(`url`).Order(`id desc`).Limit(1).Select(); err != nil {
		fmt.Println(err)
	}

	if len(redirect.URL) > 0 {
		http.Redirect(w, r, redirect.URL, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/get", handler)
	http.ListenAndServe(":8069", nil)
}
