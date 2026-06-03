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

func (h *OrderHandler) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	if req.UserUUID == "" {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "user_uuid is required",
		}, nil
	}

	if len(req.PartUuids) == 0 {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "part_uuids is required",
		}, nil
	}

	partsResp, err := h.inventoryClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Uuids: req.PartUuids,
	})
	if err != nil {
		return &order_v1.InternalServerError{
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
			return &order_v1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("part with uuid %s not found", partUUID),
			}, nil
		}
		totalPrice += part.GetPrice()
	}

	orderUUID := uuid.NewString()

	order := &order_v1.Order{
		OrderUUID:  order_v1.NewOptString(orderUUID),
		UserUUID:   order_v1.NewOptString(req.UserUUID),
		PartUuids:  req.PartUuids,
		TotalPrice: order_v1.NewOptFloat64(totalPrice),
		Status:     order_v1.NewOptOrderStatus(order_v1.OrderStatusPENDINGPAYMENT),
	}

	h.storage.mu.Lock()
	h.storage.orders[orderUUID] = order
	h.storage.mu.Unlock()

	return &order_v1.CreatedOrder{
		UUID:       order_v1.NewOptString(orderUUID),
		TotalPrice: order_v1.NewOptFloat64(totalPrice),
	}, nil
}

func (h *OrderHandler) CancelOrder(_ context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	orderUUID := params.OrderUUID.String()

	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	storedOrder, ok := h.storage.orders[orderUUID]
	if !ok {
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	if storedOrder.Status.Value == order_v1.OrderStatusPAID {
		return &order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusConflict),
			Message: order_v1.NewOptString("paid order cannot be cancelled"),
		}, nil
	}

	storedOrder.Status = order_v1.NewOptOrderStatus(order_v1.OrderStatusCANCELLED)

	return &order_v1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) GetOrder(_ context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	orderUUID := params.OrderUUID.String()

	h.storage.mu.RLock()
	storedOrder, ok := h.storage.orders[orderUUID]
	h.storage.mu.RUnlock()

	if !ok {
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	return storedOrder, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	orderUUID := params.OrderUUID.String()

	h.storage.mu.RLock()
	storedOrder, ok := h.storage.orders[orderUUID]
	if !ok {
		h.storage.mu.RUnlock()
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	if storedOrder.Status.Value == order_v1.OrderStatusCANCELLED {
		h.storage.mu.RUnlock()
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "cancelled order cannot be paid",
		}, nil
	}

	userUUID := storedOrder.UserUUID.Value
	h.storage.mu.RUnlock()

	paymentMethod, orderPaymentMethod, err := mapPaymentMethod(req.PaymentMethod)
	if err != nil {
		return &order_v1.BadRequestError{
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
		return &order_v1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}

	h.storage.mu.Lock()
	defer h.storage.mu.Unlock()

	storedOrder, ok = h.storage.orders[orderUUID]
	if !ok {
		return &order_v1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	if storedOrder.Status.Value == order_v1.OrderStatusCANCELLED {
		return &order_v1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "cancelled order cannot be paid",
		}, nil
	}

	transactionUUID := paymentResp.GetTransactionUuid()

	storedOrder.Status = order_v1.NewOptOrderStatus(order_v1.OrderStatusPAID)
	storedOrder.TransactionUUID = order_v1.NewOptString(transactionUUID)
	storedOrder.PaymentMethod = order_v1.NewOptOrderPaymentMethod(orderPaymentMethod)

	return &order_v1.OrderPayment{
		TransactionUUID: transactionUUID,
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

func (*OrderHandler) NewError(_ context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: order_v1.GenericError{
			Code:    order_v1.NewOptInt(http.StatusInternalServerError),
			Message: order_v1.NewOptString(err.Error()),
		},
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

	storageServer, err := order_v1.NewServer(storageHandler)
	if err != nil {
		slog.Error("Ошибка создания сервера OpenAPI", "error", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", storageServer)

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
