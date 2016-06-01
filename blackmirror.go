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

var version string // Set by build, e.g. -ldflags "-X main.version=0.0.1"
var conf config.Config

func main() {
	conf.Version = version // TODO Any way to set nested ldflags?
	app := cli.NewApp()
	app.Name = "blackmirror"
	app.Usage = "reflect HTTP requests back as a response"
	if app.Version = conf.Version; app.Version == "" {
		conf.Version = "unversioned"
		app.Version = conf.Version
	}
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
	log.Printf("blackmirror (%s): starting on %s", conf.Version, conf.Address())
	if err := http.ListenAndServe(conf.Address(), nil); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func dump(w http.ResponseWriter, r *http.Request) {
	// Add the version to the log and response
	log.Printf("blackmirror (%s): %s %s", conf.Version, r.Method, r.URL.Path)

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Blackmirror-Version", conf.Version)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", dump)
	return
}
