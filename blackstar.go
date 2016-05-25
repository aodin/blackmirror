package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/aodin/config"
)

var conf = config.Config{}

var (
	DefaultPort = 8080
	DefaultHost = ""
)

func init() {
	if port, _ := strconv.Atoi(os.Getenv("PORT")); port > 0 {
		DefaultPort = port
	}
	if host := os.Getenv("HOST"); host != "" {
		DefaultHost = host
	}

	flag.IntVar(&conf.Port, "port", DefaultPort, "port for service")
	flag.StringVar(&conf.Domain, "host", DefaultHost, "host for service")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", dump)
	log.Printf("blackstar: starting on %s", conf.Address())
	log.Fatal(http.ListenAndServe(conf.Address(), nil))
}

func dump(w http.ResponseWriter, r *http.Request) {
	log.Printf("blackstar: %s %s", r.Method, r.URL.Path)
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", dump)
	return
}
