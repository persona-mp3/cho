package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

const (
	LogDir          = "./logs"
	LogDest         = "./logs/structured_logs.txt"
	DefaultDuration = 530 * time.Millisecond
)

type Log struct {
	Diagnostics string
}

func (l *Log) String() string {
	return fmt.Sprintf("%s", l.Diagnostics)
}

/*
{
	"time":"2026-06-21T00:35:22.088749179+01:00",
  "level":"INFO",
	"source":{
		"function":"main.main",
		"file":"/home/daniel/dev/cho/tools/log.go",
		"line":59
	},
	"diagnostics":"Generating random logs"
}

*/

func main() {

	if err := os.MkdirAll(LogDir, 0777); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatalf("could not create %s. Reason: %s\n", LogDir, err)
		}
	}

	file, err := os.Create(LogDest)
	if err != nil {
		log.Fatalf("could not create logFile %s. Reason: %s\n", LogDest, err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "diagnostics"
			}

			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("02-01-2006 15:04:05"))
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(file, opts)
	logger := slog.New(handler)

	ticker := time.NewTicker(DefaultDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l := Log{Diagnostics: "Generating random logs"}
			logger.Error(l.String())
		}

	}
}
