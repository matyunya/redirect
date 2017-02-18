package main

import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/get", handler)
	http.ListenAndServe(":8069", nil)
}
