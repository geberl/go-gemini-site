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
)

func main() {
	baseUrl := os.Getenv("HN_BASE_URL")
	if baseUrl == "" {
		baseUrl = "localhost"
	}
	logLevel := os.Getenv("HN_LOG_LEVEL")
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
		mux.HandleFunc("/", frontHandler(baseUrl, logger))
		mux.HandleFunc("/about", aboutHandler(baseUrl, logger))
		mux.HandleFunc("/item/", itemHandler(baseUrl, logger))

		server := &gemini.Server{
			Handler:        gemini.LoggingMiddleware(mux),
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
			level.Debug(logger).Log("msg", "signal func g add")
			err := execute()
			if se, ok := err.(run.SignalError); ok {
				level.Info(logger).Log("signal", se.Signal)
				return nil
			}
			return err
		}, func(err error) {
			level.Debug(logger).Log("msg", "signal func g interrupt")
			interrupt(err)
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("msg", "error running groups", "err", err)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "app exiting")

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)

	// errch := make(chan error)
	// go func() {
	// 	ctx := context.Background()
	// 	errch <- server.ListenAndServe(ctx)
	// }()

	// select {
	// case err := <-errch:
	// 	log.Fatal(err)
	// case <-c:
	// 	log.Println("Shutting down...")
	// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// 	defer cancel()
	// 	err := server.Shutdown(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}
