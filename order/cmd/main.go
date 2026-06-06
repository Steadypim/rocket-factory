package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderhttp "github.com/Steadypim/rocket-factory/order/internal/api/http/order"
	inventoryclient "github.com/Steadypim/rocket-factory/order/internal/client/inventory"
	paymentclient "github.com/Steadypim/rocket-factory/order/internal/client/payment"
	orderrepository "github.com/Steadypim/rocket-factory/order/internal/repository/order"
	orderservice "github.com/Steadypim/rocket-factory/order/internal/service/order"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort       = "8080"
	paymentAddress = "localhost:50052"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("Ошибка подключения к InventoryService", "error", err)
		os.Exit(1)
	}
	defer inventoryConn.Close()

	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("Ошибка подключения к PaymentService", "error", err)
		os.Exit(1)
	}
	defer paymentConn.Close()

	paymentClient := payment_v1.NewPaymentServiceClient(paymentConn)

	orderRepository := orderrepository.NewOrderRepository()
	inventoryGateway := inventoryclient.NewClient(inventoryClient)
	paymentGateway := paymentclient.NewClient(paymentClient)
	orderService := orderservice.NewOrderService(orderRepository, inventoryGateway, paymentGateway)
	orderHandler := orderhttp.NewHandler(orderService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	order_v1.HandlerFromMux(order_v1.NewStrictHandler(orderHandler, nil), r)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		slog.Info("🚀 HTTP-сервер запущен на порту ", "port", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("❌ Ошибка запуска сервера:", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("❌ Ошибка при остановке сервера", "error", err)
	}

	slog.Info("✅ Сервер остановлен")

}
