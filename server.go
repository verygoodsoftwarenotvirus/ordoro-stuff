package main

import (
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Query())
}

func main() {
	http.HandleFunc("/test", handleRequest)
	http.ListenAndServe(":3000", nil)
}
