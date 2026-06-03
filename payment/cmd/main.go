package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

// paymentService реализует gRPC сервис оплаты заказов.
type paymentService struct {
	payment_v1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(_ context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transactionUUID := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return
	}

	s := grpc.NewServer()
	payment_v1.RegisterPaymentServiceServer(s, &paymentService{})
	reflection.Register(s)

	go func() {
		slog.Info("🚀 gRPC payment server listening", "port", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			slog.Error("failed to serve", "error", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("🛑 Shutting down gRPC payment server...")
	s.GracefulStop()
	slog.Info("✅ Payment server stopped")
}
