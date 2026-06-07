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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	order_api "github.com/Steadypim/rocket-factory/order/internal/api/order/v1"
	inventory_client "github.com/Steadypim/rocket-factory/order/internal/client/grpc/inventory/v1"
	payment_client "github.com/Steadypim/rocket-factory/order/internal/client/grpc/payment/v1"
	order_repository "github.com/Steadypim/rocket-factory/order/internal/repository/order"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort         = "8080"
	inventoryAddress = "localhost:50051"
	paymentAddress   = "localhost:50052"

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("failed to create inventory client", "error", err)
		os.Exit(1)
	}

	paymentConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			slog.Error("failed to close inventory client", "error", closeErr)
		}
		slog.Error("failed to create payment client", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			slog.Error("failed to close inventory client", "error", closeErr)
		}
	}()
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			slog.Error("failed to close payment client", "error", closeErr)
		}
	}()

	inventoryClient := inventory_client.NewClient(
		inventory_v1.NewInventoryServiceClient(inventoryConn),
	)
	paymentClient := payment_client.NewClient(
		payment_v1.NewPaymentServiceClient(paymentConn),
	)

	orderRepository := order_repository.NewOrderRepository()
	orderService := order_service.NewOrderService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)
	orderAPI := order_api.NewOrderAPI(orderService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	order_v1.HandlerFromMux(order_v1.NewStrictHandler(orderAPI, nil), router)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		slog.Info("order HTTP server started", "port", httpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("order HTTP server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shut down order HTTP server", "error", err)
		return
	}

	slog.Info("order HTTP server stopped")
}
