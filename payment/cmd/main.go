package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	payment_api "github.com/Steadypim/rocket-factory/payment/internal/api/payment/v1"
	payment_repository "github.com/Steadypim/rocket-factory/payment/internal/repository/payment"
	payment_service "github.com/Steadypim/rocket-factory/payment/internal/service/payment"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	paymentRepository := payment_repository.NewPaymentRepository()
	paymentService := payment_service.NewPaymentService(paymentRepository)
	paymentAPI := payment_api.NewPaymentAPI(paymentService)

	server := grpc.NewServer()
	payment_v1.RegisterPaymentServiceServer(server, paymentAPI)
	reflection.Register(server)

	go func() {
		slog.Info("payment gRPC server started", "port", grpcPort)
		if err := server.Serve(listener); err != nil {
			slog.Error("payment gRPC server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.GracefulStop()
	slog.Info("payment gRPC server stopped")
}
