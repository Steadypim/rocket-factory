package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiv1 "github.com/Steadypim/rocket-factory/inventory/internal/api/inventory/v1"
	partrepo "github.com/Steadypim/rocket-factory/inventory/internal/repository/part"
	partservice "github.com/Steadypim/rocket-factory/inventory/internal/service/part"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return
	}

	grpcServer := grpc.NewServer()

	partRepository := partrepo.NewRepository()
	partService := partservice.NewService(partRepository)
	inventoryAPI := apiv1.NewAPI(partService)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, inventoryAPI)
	reflection.Register(grpcServer)

	go func() {
		slog.Info("🚀 gRPC inventory server listening", "port", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("failed to serve", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🛑 Shutting down gRPC inventory server...")
	grpcServer.GracefulStop()
	slog.Info("✅ Inventory server stopped")
}
