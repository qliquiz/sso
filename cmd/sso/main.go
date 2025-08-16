package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger/handlers/slogpretty"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// TODO: Вынести JWT и ротацию в internal/jwk + internal/authn/jwt.go, добавить JWKS сервер (через http/grpc-gateway)
// TODO: Разделить services/auth (аутентика) и services/identity (профиль/пароль/сброс)
// TODO: Добавить permissions сервис и gRPC-интерцептор для authz
// TODO: В storage/psql разнести репозитории (users/sessions/tokens/permissions/jwk)
// TODO: Тесты: сделать E2E suite, который крутит миграции и чистит БД между тестами
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	if cfg.Env == envLocal {
		log.Info("starting application", slog.Any("config", cfg))
	}

	storageDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.DBName,
		cfg.Storage.SSLMode,
	)

	application := app.New(log, cfg.GRPC.Port, storageDSN, cfg.TokenTTL)
	go application.GrpcApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	application.GrpcApp.Stop()
	log.Info("application stopped with signal:", sign.String())
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
