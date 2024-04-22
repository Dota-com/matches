package main

import (
	"log/slog"
	"matches/internal/app"
	"matches/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.MustLoad()
	log := config.SetupLoger(conf.Env)
	log.Debug("Debug running")

	matches := app.New(log, conf.Grpc.Port, conf.ConfigDB)

	go matches.GRPC.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping app matches", slog.String("signal", sign.String()))

	matches.GRPC.Stop()

	log.Info("App matches stop")
}
