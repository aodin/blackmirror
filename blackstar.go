package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/aodin/config"
	"gopkg.in/urfave/cli.v2" // imports as "cli"
)

var conf config.Config

func main() {
	app := cli.NewApp()
	app.Name = "blackstar"
	app.Usage = "reflect HTTP requests back as a response"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:    "port, p",
			Value:   8080,
			Usage:   "server port",
			EnvVars: []string{"PORT"},
		},
		cli.StringFlag{
			Name:    "host, h",
			Value:   "",
			Usage:   "server host",
			EnvVars: []string{"HOST"},
		},
	}
	app.Action = server
	app.Run(os.Args)
}

func server(ctx *cli.Context) error {
	conf.Port = ctx.Int("port")
	conf.Domain = ctx.String("host")
	http.HandleFunc("/", dump)
	log.Printf("blackstar: starting on %s", conf.Address())
	if err := http.ListenAndServe(conf.Address(), nil); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
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
