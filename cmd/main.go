package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"gitlab.com/kokishin/server/internal"
)

var (
	directory = flag.String("directory", ".", "the directory to be served")
)

func main() {
	flag.Parse()

	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	server := internal.NewServer(&internal.ServerOpts{Logger: logger, Directory: *directory})
	go server.Start()

	signals := make(chan os.Signal)

	signal.Notify(signals, os.Interrupt, os.Kill)

	select {
	case <-server.Done():
		slog.Info("bye...")
	case <-signals:
		slog.Info("being killed is not nice...")
		os.Exit(0)
	}
}
