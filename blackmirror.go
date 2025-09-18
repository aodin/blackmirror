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
	"github.com/urfave/cli/v3"
)

var version string // Set by build, e.g. -ldflags "-X main.version=0.0.1"
var conf config.Config

func main() {
	conf.Version = version

	app := &cli.Command{
		Name:  "blackmirror",
		Usage: "reflect HTTP requests back as a response",
	}

	// Set version
	if app.Version = conf.Version; app.Version == "" {
		conf.Version = "dev"
		app.Version = conf.Version
	}

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "port, p",
			Value:   8080,
			Usage:   "server port",
			Sources: cli.EnvVars("PORT"),
		},
		&cli.StringFlag{
			Name:    "host, h",
			Value:   "",
			Usage:   "server host",
			Sources: cli.EnvVars("HOST"),
		},
	}

	app.Action = server

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func server(ctx context.Context, cmd *cli.Command) error {
	// The server will gracefully exit on any interrupt or SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	conf.Port = cmd.Int("port")
	conf.Domain = cmd.String("host")

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
