package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aodin/config"
)

var conf = config.Config{
	Port: 8081,
}

func main() {
	http.HandleFunc("/", headers)
	log.Printf("blackstar: starting on %s", conf.Address())
	log.Fatal(http.ListenAndServe(conf.Address(), nil))
}

func headers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", r.URL.Path)
	r.Header.Write(w)
	return
}
