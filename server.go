/*
	Url proxy, uh, proxies a url, maybe so you don't have to expose an API key or something.
	To get up & running, make sure you have [go](http://golang.org) installed, `cd` to the
	directory, run `go build`, then `ENDPOINT=[url to proxy] ./url-proxy`.
	It reads from the environemnt variable `ENDPOINT`, so you can set that however you like.
*/
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var (
	// cfg is the global configuration for the server. It's read in at startup from
	// the config.json file and enviornment variables, see config.go for more info.
	cfg *config

	// When was the last alert sent out?
	// Use this value to avoid bombing alerts
	lastAlertSent *time.Time

	// log output
	log = logrus.New()
)

func init() {
	log.Out = os.Stderr
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}
}

func main() {
	var err error
	// setup config
	mode := os.Getenv("GOLANG_ENV")
	if mode == "" {
		mode = DEVELOP_MODE
	}

	cfg, err = initConfig(mode)
	if err != nil {
		// panic if the server is missing a vital configuration detail
		log.Fatal(fmt.Errorf("server configuration error: %s", err.Error()))
	}

	s := &http.Server{}
	m := http.NewServeMux()
	m.Handle("/", middleware(ProxyHandler(cfg.Endpoint)))
	m.HandleFunc("/.well-known/acme-challenge/", CertbotHandler)

	// connect mux to server
	s.Handler = m

	// fire it up!
	log.Printf("starting server on port %s in %s mode\n", cfg.Port, mode)
	// start server wrapped in a log.Fatal b/c http.ListenAndServe will not
	// return unless there's an error, which would be a program crash
	log.Fatal(StartServer(cfg, s))
}
