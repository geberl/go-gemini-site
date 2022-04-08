package main

import (
	"context"
	"os"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/mattn/go-isatty"
)

func makeLogger(logLevel string) log.Logger {
	var logger log.Logger

	w := log.NewSyncWriter(os.Stdout)
	if isatty.IsTerminal(os.Stdout.Fd()) {
		logger = log.NewLogfmtLogger(w)
	} else {
		logger = log.NewJSONLogger(w)
	}

	switch logLevel {
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	return logger
}

func loggingMiddleware(h gemini.Handler, logger log.Logger) gemini.Handler {
	return gemini.HandlerFunc(func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		startTime := time.Now()
		lw := &logResponseWriter{CorrelationId: uuid.New().String(), rw: w}

		level.Info(logger).Log("msg", "request received", "correlation_id", lw.CorrelationId, "url", r.URL)
		h.ServeGemini(ctx, lw, r)

		duration := time.Now().Sub(startTime)
		level.Info(logger).Log(
			"msg", "request handled",
			"correlation_id", lw.CorrelationId,
			"url", r.URL,
			"status", lw.Status,
			"bytes", lw.Wrote,
			"duration", duration.Milliseconds(),
		)
	})
}

type logResponseWriter struct {
	CorrelationId string
	Status        gemini.Status
	Wrote         int
	rw            gemini.ResponseWriter
	mediatype     string
	wroteHeader   bool
}

func (w *logResponseWriter) SetMediaType(mediatype string) {
	w.mediatype = mediatype
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		meta := w.mediatype
		if meta == "" {
			meta = "text/gemini" // default media type
		}
		w.WriteHeader(gemini.StatusSuccess, meta)
	}
	n, err := w.rw.Write(b)
	w.Wrote += n
	return n, err
}

func (w *logResponseWriter) WriteHeader(status gemini.Status, meta string) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.Status = status
	w.Wrote += len(meta) + 5
	w.rw.WriteHeader(status, meta)
}

func (w *logResponseWriter) Flush() error {
	return nil
}
