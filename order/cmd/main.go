package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderHandler "github.com/melkomukovki/go-or-die/order/pkg/handler"
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

	// Создаём хранилище и обработчик
	store := orderHandler.NewOrderStore()
	h := orderHandler.NewOrderHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		store,
	)

	// Создать OpenAPI сервер
	orderServer, err := orderHandler.SetupServer(h)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		return
	}

	httpServer := &http.Server{
		Addr:              httpAddress,
		Handler:           orderServer,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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
