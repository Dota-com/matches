package appgrpc

import (
	"google.golang.org/grpc"
	"log/slog"
	matches_server "matches/internal/grpc/matches"
)

type App struct {
	log  *slog.Logger
	grpc grpc.Server
	port int
}

func (a *App) New(log *slog.Logger, port int) *App {
	grpcService := grpc.NewServer()
	matches_server.RegisterMatchesServer(grpcService)

	return &App{
		log:  log,
		port: port,
		grpc: grpcService,
	}
}
