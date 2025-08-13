package suite

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	ssov1 "github.com/qliquiz/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log/slog"
	"net"
	"os"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/services/auth"
	"sso/internal/storage/psql"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("error loading .env file: %v", err)
	}

	key := "TEST_CONFIG_PATH"
	var configPath string
	if configPath = os.Getenv(key); configPath == "" {
		panic("config path is not set: " + key)
	}
	cfg := config.MustLoadByPath(configPath)

	storageDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.DBName,
		cfg.Storage.SSLMode,
	)

	storage, err := psql.New(storageDSN)
	if err != nil {
		t.Fatalf("failed to init storage: %v", err)
	}

	silentLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	authService := auth.NewAuth(silentLogger, storage, storage, storage, cfg.TokenTTL)
	application := grpcapp.New(silentLogger, authService, cfg.GRPC.Port)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err = application.Serve(listener); err != nil {
			t.Errorf("failed to serve gRPC server: %v", err)
		}
	}()
	t.Cleanup(func() {
		application.Stop()
	})

	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create gRPC client: %v", err)
	}
	t.Cleanup(func() {
		err = conn.Close()
		if err != nil {
			t.Errorf("failed to close grpc connection: %v", err)
		}
	})

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(cancelCtx)

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(conn),
	}
}
