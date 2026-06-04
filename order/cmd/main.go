package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
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

// OrderStorage потокобезопасное хранилище данных о заказах
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*order_v1.Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*order_v1.Order),
	}
}

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventory_v1.InventoryServiceClient
	paymentClient   payment_v1.PaymentServiceClient
}

func NewOrderHandler(
	storage *OrderStorage,
	inventoryClient inventory_v1.InventoryServiceClient,
	paymentClient payment_v1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (h *OrderHandler) CreateOrder(
	ctx context.Context,
	request order_v1.CreateOrderRequestObject,
) (order_v1.CreateOrderResponseObject, error) {
	if request.Body == nil {
		return order_v1.CreateOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "request body is required",
		}, nil
	}

	req := request.Body
	if req.UserUuid == "" {
		return order_v1.CreateOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "user_uuid is required",
		}, nil
	}

	if len(req.PartUuids) == 0 {
		return order_v1.CreateOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "part_uuids is required",
		}, nil
	}

	partsResp, err := h.inventoryClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Uuids: req.PartUuids,
	})
	if err != nil {
		return order_v1.CreateOrder500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	foundParts := make(map[string]*inventory_v1.Part, len(partsResp.GetParts()))
	for _, part := range partsResp.GetParts() {
		foundParts[part.GetUuid()] = part
	}

	var totalPrice float64

	for _, partUUID := range req.PartUuids {
		part, ok := foundParts[partUUID]
		if !ok {
			return order_v1.CreateOrder400JSONResponse{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("part with uuid %s not found", partUUID),
			}, nil
		}
		totalPrice += part.GetPrice()
	}

	orderUUID := uuid.NewString()

	order := &order_v1.Order{
		OrderUuid:  orderUUID,
		UserUuid:   req.UserUuid,
		PartUuids:  req.PartUuids,
		TotalPrice: float32(totalPrice),
		Status:     order_v1.PENDINGPAYMENT,
	}

	h.storage.mu.Lock()
	h.storage.orders[orderUUID] = order
	h.storage.mu.Unlock()

	return order_v1.CreateOrder200JSONResponse{
		Uuid:       orderUUID,
		TotalPrice: float32(totalPrice),
	}, nil
}

func (h *OrderHandler) CancelOrder(
	_ context.Context,
	request order_v1.CancelOrderRequestObject,
) (order_v1.CancelOrderResponseObject, error) {
	orderUUID := request.OrderUuid.String()

	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	storedOrder, ok := h.storage.orders[orderUUID]
	if !ok {
		return order_v1.CancelOrder404JSONResponse{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	if storedOrder.Status == order_v1.PAID {
		return order_v1.CancelOrder409JSONResponse{
			Code:    http.StatusConflict,
			Message: "paid order cannot be cancelled",
		}, nil
	}

	storedOrder.Status = order_v1.CANCELLED

	return order_v1.CancelOrder204Response{}, nil
}

func (h *OrderHandler) GetOrder(
	_ context.Context,
	request order_v1.GetOrderRequestObject,
) (order_v1.GetOrderResponseObject, error) {
	orderUUID := request.OrderUuid.String()

	h.storage.mu.RLock()
	storedOrder, ok := h.storage.orders[orderUUID]
	h.storage.mu.RUnlock()

	if !ok {
		return order_v1.GetOrder404JSONResponse{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	return order_v1.GetOrder200JSONResponse(*storedOrder), nil
}

func (h *OrderHandler) PayOrder(
	ctx context.Context,
	request order_v1.PayOrderRequestObject,
) (order_v1.PayOrderResponseObject, error) {
	orderUUID := request.OrderUuid.String()

	if request.Body == nil {
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "request body is required",
		}, nil
	}

	req := request.Body

	h.storage.mu.RLock()
	storedOrder, ok := h.storage.orders[orderUUID]
	if !ok {
		h.storage.mu.RUnlock()
		return order_v1.PayOrder404JSONResponse{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	if storedOrder.Status == order_v1.CANCELLED {
		h.storage.mu.RUnlock()
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "cancelled order cannot be paid",
		}, nil
	}

	userUUID := storedOrder.UserUuid
	h.storage.mu.RUnlock()

	paymentMethod, orderPaymentMethod, err := mapPaymentMethod(req.PaymentMethod)
	if err != nil {
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	paymentResp, err := h.paymentClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return order_v1.PayOrder500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	storedOrder, ok = h.storage.orders[orderUUID]
	if !ok {
		return order_v1.PayOrder404JSONResponse{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	if storedOrder.Status == order_v1.CANCELLED {
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "cancelled order cannot be paid",
		}, nil
	}

	transactionUUID := paymentResp.GetTransactionUuid()

	storedOrder.Status = order_v1.PAID
	storedOrder.TransactionUuid = transactionUUID
	storedOrder.PaymentMethod = orderPaymentMethod

	return order_v1.PayOrder200JSONResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func mapPaymentMethod(method order_v1.PayOrderRequestPaymentMethod) (payment_v1.PaymentMethod, order_v1.OrderPaymentMethod, error) {
	switch method {
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD, order_v1.OrderPaymentMethodPAYMENTMETHODCARD, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP, order_v1.OrderPaymentMethodPAYMENTMETHODSBP, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, order_v1.OrderPaymentMethodPAYMENTMETHODCREDITCARD, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, order_v1.OrderPaymentMethodPAYMENTMETHODINVESTORMONEY, nil
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN, order_v1.OrderPaymentMethodPAYMENTMETHODUNKNOWN, fmt.Errorf("unsupported payment method %q", method)
	}
}

func main() {
	storage := NewOrderStorage()

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

	storageHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	order_v1.HandlerFromMux(order_v1.NewStrictHandler(storageHandler, nil), r)

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
