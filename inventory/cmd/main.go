package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventory_api "github.com/Steadypim/rocket-factory/inventory/internal/api/inventory/v1"
	inventory_repository "github.com/Steadypim/rocket-factory/inventory/internal/repository/inventory"
	inventory_service "github.com/Steadypim/rocket-factory/inventory/internal/service/inventory"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	inventoryRepository := inventory_repository.NewInventoryRepository()
	inventoryService := inventory_service.NewInventoryService(inventoryRepository)
	inventoryAPI := inventory_api.NewInventoryAPI(inventoryService)

	server := grpc.NewServer()
	inventory_v1.RegisterInventoryServiceServer(server, inventoryAPI)
	reflection.Register(server)

	go func() {
		slog.Info("inventory gRPC server started", "port", grpcPort)
		if err := server.Serve(listener); err != nil {
			slog.Error("inventory gRPC server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.GracefulStop()
	slog.Info("inventory gRPC server stopped")
}
