package main

import (
	"ayinke-llc/gophercrunch/testing-go/cmd/config"
	"ayinke-llc/gophercrunch/testing-go/server"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	var httpPort int = 3200

	flag.IntVar(&httpPort, "http.port", 3200, "http port to run web server on")

	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("could not load configuration details")
	}

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = logrus.DebugLevel
	}

	logrus.SetLevel(lvl)

	h, _ := os.Hostname()

	logger := logrus.WithField("host", h).
		WithField("app", "testing-go")

	logger.Debug("starting app")

	srv := server.New(cfg, httpPort)

	go func() {
		logger.Debug("starting HTTP server")
		if err := srv.ListenAndServe(); err != nil {
			logger.WithError(err).Error("an error occured while shutting down http server")
		}
	}()

	<-sig
	logger.Debug("shutting down app")
}
