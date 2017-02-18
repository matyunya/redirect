package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pg "gopkg.in/pg.v5"
)

type Redirect struct {
	ID                 int
	URL, P1            string
	P2, P3, P4, P5, P6 []string
}

type List []Redirect

var DB *pg.DB

func init() {
	DB = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "mac",
		Database: "redirect",
		Password: os.Getenv("pass"),
	})

	pg.SetQueryLogger(log.New(os.Stdout, "", log.LstdFlags))
}

func main() {
	file, err := os.Open("redirects_fixtures.json")
	if err != nil {
		log.Println("Error opening file:", err)
		os.Exit(1)
	}

	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		os.Exit(1)
	}

	redirects := List{}
	json.Unmarshal(b, &redirects)

	for _, r := range redirects {
		_, err := DB.Query(nil,
			`insert into redirects (id, url, p1, p2, p3, p4, p5, p6) values (?, ?, ?, ?, ?, ?, ?, ?)`,
			r.ID, r.URL, r.P1, pg.Array(r.P2), pg.Array(r.P3), pg.Array(r.P4), pg.Array(r.P5), pg.Array(r.P6))

		if err != nil {
			log.Println("error inserting redirect", err)
		} else {
			log.Println("successfully inserted redirect", r)
		}
	}
}
