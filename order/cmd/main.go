package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/melkomukovki/go-or-die/order/pkg/app"
	inventoryv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/melkomukovki/go-or-die/shared/pkg/proto/payment/v1"
)

const (
	// Адреса сервисов.
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"

	// HTTP-сервер параметры.
	httpAddress           = ":8080"
	httpReadHeaderTimeout = 5 * time.Second
	httpReadTimeout       = 15 * time.Second
	httpWriteTimeout      = 15 * time.Second
	httpIdleTimeout       = 60 * time.Second
	httpShutdownTimeout   = 5 * time.Second

	// gRPC-клиент keepalive параметры.
	grpcKeepaliveTime    = 5 * time.Minute
	grpcKeepaliveTimeout = 1 * time.Second

	envFileLocation = "order.env"
)

func main() {
	// Создать gRPC соединение с InventoryService
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}))
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		return
	}
	defer func() {
		closeErr := inventoryConn.Close()
		if closeErr != nil {
			slog.Error("ошибка при закрытии grpc соединения с InventoryService", "error", closeErr)
		} else {
			slog.Info("grpc соединение с InventoryService успешно закрыто")
		}
	}()

	// Создать gRPC клиент PaymentService
	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}))
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		return
	}
	defer func() {
		closeErr := paymentConn.Close()
		if closeErr != nil {
			slog.Error("ошибка при закрытии grpc соединения с PaymentService", "error", closeErr)
		} else {
			slog.Info("grpc соединение с PaymentService успешно закрыто")
		}
	}()

	inventoryClient := inventoryv1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentv1.NewPaymentServiceClient(paymentConn)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = godotenv.Load(envFileLocation)
	if err != nil {
		slog.Error("ошибка загрузки переменных окружения", "error", err, "file", envFileLocation)
	}

	pool, txManager, err := newDB(ctx)
	if err != nil {
		slog.Error("ошибка инициализации БД", "error", err)
	}

	handler, err := app.NewHTTPHandler(pool, txManager, inventoryClient, paymentClient)
	if err != nil {
		slog.Error("не удалось создать HTTP обработчик", "error", err)
		return
	}

	httpServer := &http.Server{
		Addr:              httpAddress,
		Handler:           handler,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	go func() {
		slog.Info("HTTP сервер запущен", "address", httpAddress)
		if serverErr := httpServer.ListenAndServe(); serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска HTTP сервера", "error", serverErr)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("останавливаем HTTP сервер")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer shutdownCancel()

	if shutdownErr := httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
		slog.Error("ошибка при остановке HTTP сервера", "error", shutdownErr)
	}
	slog.Info("HTTP сервер остановлен")
}

func newDB(ctx context.Context) (*pgxpool.Pool, *manager.Manager, error) {
	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		return nil, nil, errors.New("переменная окружения DB_URI не установлена")
	}

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	defer pool.Close()

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("ошибка при создании транзакционного менеджера", "error", err)
		return nil, nil, fmt.Errorf("ошибка при создании транзакционного менеджера: %w", err)
	}
	return pool, txManager, nil
}
