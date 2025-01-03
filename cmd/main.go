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
	security  = flag.String("security", "", "the file used to validate authentication of users, it should be a CSV structured file, with a user and a hashed password (SHA512) per line. If this variable is nil it won't use authentication for any of the internal endpoints")
)

func main() {
	flag.Parse()

	level := logger.From(*level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level.ToSlogLevel()})
	logger := slog.New(handler)

	server, _ := internal.NewServer(&internal.ServerOpts{
		Logger:             logger,
		Directory:          *directory,
		AuthenticationFile: *security,
	})

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
