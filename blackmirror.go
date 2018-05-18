package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"

	"github.com/aodin/config"
	cli "gopkg.in/urfave/cli.v2"
)

var version string // Set by build, e.g. -ldflags "-X main.version=0.0.1"
var conf config.Config

func main() {
	conf.Version = version
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
	// The server will gracefully exit on any interrupt or SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	conf.Port = ctx.Int("port")
	conf.Domain = ctx.String("host")

	// Serve the dump handler on all paths
	srv := &http.Server{
		Addr:    conf.Address(),
		Handler: http.HandlerFunc(dump),
	}

	go func() {
		log.Printf(
			"blackmirror (%s): starting on %s", conf.Version, conf.Address(),
		)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Printf("blackmirror (%s): shutting down", conf.Version)
	srv.Shutdown(context.Background())
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
}
