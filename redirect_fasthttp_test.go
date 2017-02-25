package redirect_test

import (
	"log"
	"math/rand"
	"net/http"
	"testing"

	"github.com/matyunya/redirect/redirect"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

func runServer() {
	ln, err := reuseport.Listen("tcp4", "localhost:8069")
	if err != nil {
		log.Fatalf("error in reuseport listener: %s", err)
	}

	go func() {
		if err = fasthttp.Serve(ln, redirect.Handler); err != nil {
			log.Fatalf("error in fasthttp Server: %s", err)
		}
	}()
}

var urls = []string{
	"p1=63c89c35-1ced-43af-85d5-f45fd2dad28d",
	"p6=f9f7e2db-c207-404a-8bf6-5f7192660ef5",
	"p1=a0fde34c-a00b-4035-a0b9-a6e961ffc82e",
	"p1=a95410e6-2f98-4bb4-9adc-77ad9bc0ffb4",
	"p2=CPM",
	"p4=GS",
	"p3=Paragonia",
	"p3=Endicil&p2=CPC",
	"p5=news&p3=Zillar",
	"p3=Magmina&p5=news",
	"p4=SM",
	"p3=Quilch",
	"p3=Venoflex",
}

var urlsResponses = map[string]string{
	"p1=63c89c35-1ced-43af-85d5-f45fd2dad28d": "/redirect/16",
	"p6=f9f7e2db-c207-404a-8bf6-5f7192660ef5": "/redirect/38",
	"p1=a0fde34c-a00b-4035-a0b9-a6e961ffc82e": "/redirect/40",
	"p1=a95410e6-2f98-4bb4-9adc-77ad9bc0ffb4": "/redirect/55",
	"p2=CPM":             "/redirect/100",
	"p4=GS":              "/redirect/19",
	"p3=Paragonia":       "/redirect/39",
	"p3=Endicil&p2=CPC":  "/redirect/42",
	"p5=news&p3=Zillar":  "/redirect/45",
	"p3=Magmina&p5=news": "/redirect/48",
	"p4=SM":              "/redirect/52",
	"p3=Quilch":          "/redirect/59",
	"p3=Venoflex":        "/redirect/66",
}

var origin = `http://` + redirect.ServerHost

func r() string { return urls[rand.Intn(len(urls))] }

func BenchmarkHandler(b *testing.B) {
	runServer()
	for n := 0; n < b.N; n++ {
		fasthttp.Get(nil, origin+"/get?"+r())
	}
}

func TestHandler(t *testing.T) {
	runServer()
	for origURL, expectedURL := range urlsResponses {
		testRequestCtxRedirect(t, origin+"/get?"+origURL, redirect.RedirectHost+expectedURL)
	}
}

func testRequestCtxRedirect(t *testing.T, origURL, expectedURL string) {
	if resp, _ := http.Get(origURL); resp.Request.URL.String() != expectedURL {
		t.Log("Test failed for origin", origURL, ", Expected: ", expectedURL, ", Recieved:", resp.Request.URL.String())
		t.Fail()
	}
}
