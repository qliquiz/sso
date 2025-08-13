package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/psql"
	"time"
)

type App struct {
	GrpcApp *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, dsn string, tokenTTL time.Duration) *App {
	storage, err := psql.New(dsn)
	if err != nil {
		panic(err)
	}

	authService := auth.NewAuth(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}
}
