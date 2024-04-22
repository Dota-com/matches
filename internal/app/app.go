package app

import (
	"log/slog"
	appgrpc "matches/internal/app/app"
	"matches/internal/config"
	service_matches "matches/internal/service/matches"
	storage_matches "matches/internal/storage"
)

type App struct {
	GRPC *appgrpc.App
}

func New(log *slog.Logger, port int, storage config.DB) *App {
	storageMatches, err := storage_matches.New(storage)
	if err != nil {
		panic(err)
	}
	service := appgrpc.New(log, port)

	service_matches.New(log, storageMatches)
	return &App{GRPC: service}
}
