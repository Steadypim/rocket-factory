package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiv1 "github.com/Steadypim/rocket-factory/payment/internal/api/payment/v1"
	paymentservice "github.com/Steadypim/rocket-factory/payment/internal/service/payment"
	paymentv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return
	}

	grpcServer := grpc.NewServer()

	paymentService := paymentservice.NewService()
	paymentAPI := apiv1.NewAPI(paymentService)

	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentAPI)
	reflection.Register(grpcServer)

	go func() {
		slog.Info("🚀 gRPC payment server listening", "port", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("failed to serve", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🛑 Shutting down gRPC payment server...")
	grpcServer.GracefulStop()
	slog.Info("✅ Payment server stopped")
}
