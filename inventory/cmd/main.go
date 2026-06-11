package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventory_api "github.com/Steadypim/rocket-factory/inventory/internal/api/inventory/v1"
	inventory_repository "github.com/Steadypim/rocket-factory/inventory/internal/repository/inventory"
	inventory_service "github.com/Steadypim/rocket-factory/inventory/internal/service/inventory"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort       = 50051
	startupTimeout = 10 * time.Second
)

func main() {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("failed to load .env file", "error", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		slog.Error("MONGO_URI is required")
		return
	}

	databaseName := envOrDefault("MONGO_DATABASE", "inventory-service")
	collectionName := envOrDefault("MONGO_COLLECTION", "parts")

	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		slog.Error("failed to create MongoDB client", "error", err)
		return
	}
	defer func() {
		if disconnectErr := mongoClient.Disconnect(context.Background()); disconnectErr != nil {
			slog.Error("failed to disconnect MongoDB client", "error", disconnectErr)
		}
	}()

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		slog.Error("MongoDB is unavailable", "error", err)
		return
	}

	inventoryRepository := inventory_repository.NewInventoryRepository(
		mongoClient.Database(databaseName).Collection(collectionName),
	)
	if err := inventoryRepository.Seed(ctx); err != nil {
		slog.Error("failed to seed inventory", "error", err)
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return
	}

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

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
