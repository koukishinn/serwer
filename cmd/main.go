package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"gitlab.com/kokishin/serwer/internal"
	"gitlab.com/kokishin/serwer/internal/logger"
)

var (
	directory = flag.String("directory", ".", "the directory to be served")
	level     = flag.String("level", "warn", "the level used in the logger (info|warn|error|debug)")
)

func main() {
	flag.Parse()

	level := logger.From(*level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level.ToSlogLevel()})
	logger := slog.New(handler)

	server := internal.NewServer(&internal.ServerOpts{Logger: logger, Directory: *directory})
	go server.Start()

	signals := make(chan os.Signal)

	signal.Notify(signals, os.Interrupt, os.Kill)

	select {
	case <-server.Done():
		logger.Info("bye...")
	case <-signals:
		logger.Info("being killed is not nice...")
		os.Exit(0)
	}
}
