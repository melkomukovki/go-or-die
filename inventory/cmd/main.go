package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/melkomukovki/go-or-die/inventory/pkg/app"
)

const (
	// gRPC параметры.
	grpcAddress               = ":50051"
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute

	envFileLocation = "inventory.env"
)

func main() {
	//nolint:noctx,gosec // Контекст здесь не нужен: GracefulStop() сам закроет listener и прервёт Accept()
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		return
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime: grpcMinPingInterval,
		}),
	)

	// Создаем подключение к БД
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = godotenv.Load(envFileLocation)
	if err != nil {
		slog.Error("ошибка загрузки переменных окружения", "error", err, "file", envFileLocation)
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		slog.Error("переменная окружения DB_URI не установлена")
		return
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		slog.Error("ошибка подключения к БД", "error", err)
		return
	}
	defer pool.Close()

	app.RegisterServices(grpcServer, pool)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	go func() {
		slog.Info("запуск InventoryService", "адрес", grpcAddress)
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("остановка сервера")

	grpcServer.GracefulStop()
	slog.Info("gRPC сервер остановлен")
}
