package appgrpc

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	matches_server "matches/internal/grpc/matches"
	"net"
)

type App struct {
	log  *slog.Logger
	grpc *grpc.Server
	port int
}

func New(log *slog.Logger, port int) *App {
	grpcService := grpc.NewServer()
	matches_server.RegisterMatchesServer(grpcService)
	reflection.Register(grpcService)
	return &App{
		log:  log,
		port: port,
		grpc: grpcService,
	}
}

func (a *App) run() error {
	log := a.log.With("Match service running")

	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", a.port))

	if err != nil {
		log.Error("Error running Match service")
		return fmt.Errorf("error running Match service")
	}
	if err := a.grpc.Serve(l); err != nil {
		log.Error("Error Grpc service Matches")
		return fmt.Errorf("Error Grpc service Matches")
	}
	log.Info("Server started")

	return nil
}

func (a *App) Stop() {
	a.log.With("Stop Matches service")
	a.grpc.GracefulStop()
	a.log.Info("Matches service stop")
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}
