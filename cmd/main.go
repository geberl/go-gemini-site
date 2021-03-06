package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"

	"go-gemini-site/internal/gemlog"
	"go-gemini-site/internal/hn"
)

func main() {
	baseUrl := os.Getenv("GEMINI_SITE_BASE_URL")
	if baseUrl == "" {
		baseUrl = "localhost"
	}
	logLevel := os.Getenv("GEMINI_SITE_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logger := makeLogger(logLevel)
	level.Info(logger).Log("msg", "app starting", "baseurl", baseUrl, "loglevel", logLevel)

	certificates := &certificate.Store{}
	certificates.Register(baseUrl)
	if err := certificates.Load("certs"); err != nil {
		level.Error(logger).Log("msg", "unable to load certs", "err", err)
		log.Fatal(err)
	}

	var (
		g   run.Group
		ctx = context.Background()
	)
	ctx, cancel := context.WithCancel(ctx)
	{
		mux := &gemini.Mux{}

		mux.HandleFunc("/", HandlerIndex(baseUrl, logger))

		mux.Handle("/gemlog/", gemlog.FileServer(baseUrl, logger, "gemlog", os.DirFS("gemlog")))

		mux.HandleFunc("/hn/", hn.HandlerIndex(baseUrl, logger))
		mux.HandleFunc("/hn/item/", hn.HandlerItem(baseUrl, logger))
		mux.HandleFunc("/hn/user/", hn.HandlerUser(baseUrl, logger))

		server := &gemini.Server{
			Handler:        loggingMiddleware(mux, logger),
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   1 * time.Minute,
			GetCertificate: certificates.Get,
		}

		g.Add(func() error {
			ctx := context.Background()
			level.Info(logger).Log("msg", "gemini server starting")
			return server.ListenAndServe(ctx)
		}, func(err error) {
			level.Info(logger).Log("msg", "gemini server closing", "err", err)
			server.Shutdown(ctx)
			cancel()
		})
	}
	{
		execute, interrupt := run.SignalHandler(ctx, syscall.SIGTERM, syscall.SIGINT)
		g.Add(func() error {
			level.Debug(logger).Log("msg", "signal listener starting")
			err := execute()
			if se, ok := err.(run.SignalError); ok {
				level.Info(logger).Log("msg", "signal received", "signal", se.Signal)
				return nil
			}
			return err
		}, func(err error) {
			level.Debug(logger).Log("msg", "signal listener closing")
			interrupt(err)
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "error running groups", "err", err)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "app exiting")
}
