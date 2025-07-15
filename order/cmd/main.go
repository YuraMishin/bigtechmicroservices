package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	inventoryGrpcPort = 50051
	paymentGrpcPort   = 50052
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
	inventoryConn   *grpc.ClientConn
	paymentConn     *grpc.ClientConn
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	inventoryConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", inventoryGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}

	paymentConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", paymentGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}

	inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

	handler := &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}

	// Store connections for cleanup
	handler.inventoryConn = inventoryConn
	handler.paymentConn = paymentConn

	return handler
}

func (o OrderHandler) CreateNewOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateNewOrderRes, error) {
	orderUUID := uuid.New()

	var partUUIDs []string
	for _, partUUID := range req.PartUuids {
		partUUIDs = append(partUUIDs, partUUID.String())
	}

	filter := &inventoryV1.PartsFilter{
		Uuids: partUUIDs,
	}

	listRequest := &inventoryV1.ListPartsRequest{
		Filter: filter,
	}

	inventoryResponse, err := o.inventoryClient.ListParts(ctx, listRequest)
	if err != nil {
		log.Printf("Error calling InventoryService: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}

	if len(inventoryResponse.Parts) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Some parts not found in inventory",
		}, nil
	}

	var totalPrice float32
	for _, part := range inventoryResponse.Parts {
		totalPrice += float32(part.Price)
	}

	order := &orderV1.OrderDto{
		OrderUUID:       orderUUID,
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      totalPrice,
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
		TransactionUUID: uuid.Nil,                                              // Будет заполнено при оплате
		PaymentMethod:   orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED, // Будет заполнено при оплате
	}

	o.storage.mu.Lock()
	o.storage.orders[orderUUID.String()] = order
	o.storage.mu.Unlock()

	log.Printf("Created new order: %s with total price: %.2f", orderUUID.String(), totalPrice)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (o OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	o.storage.mu.RLock()
	order, exists := o.storage.orders[params.OrderUUID.String()]
	o.storage.mu.RUnlock()

	if !exists {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	var paymentMethod paymentV1.PaymentMethod
	switch req.PaymentMethod {
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED:
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Payment method must be specified",
		}, nil
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return &orderV1.BadRequestError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}

	paymentRequest := &paymentV1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	}

	paymentResponse, err := o.paymentClient.PayOrder(ctx, paymentRequest)
	if err != nil {
		log.Printf("Error calling PaymentService: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}

	transactionUUID, err := uuid.Parse(paymentResponse.TransactionUuid)
	if err != nil {
		log.Printf("Error parsing transaction UUID: %v", err)
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal error",
		}, nil
	}

	var orderPaymentMethod orderV1.OrderDtoPaymentMethod
	switch req.PaymentMethod {
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODUNSPECIFIED:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODCARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODSBP
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODCREDITCARD
	case orderV1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODINVESTORMONEY
	default:
		orderPaymentMethod = orderV1.OrderDtoPaymentMethodPAYMENTMETHODUNSPECIFIED
	}

	o.storage.mu.Lock()
	order.Status = orderV1.OrderDtoStatusPAID
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = orderPaymentMethod
	o.storage.mu.Unlock()

	log.Printf("Order %s paid successfully with transaction %s", params.OrderUUID.String(), transactionUUID.String())

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

func (o OrderHandler) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	if params.OrderUUID == uuid.Nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid order UUID",
		}, nil
	}

	o.storage.mu.RLock()
	order, exists := o.storage.orders[params.OrderUUID.String()]
	o.storage.mu.RUnlock()

	if !exists {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	return order, nil
}

func (o OrderHandler) CancelOrderByUUID(ctx context.Context, params orderV1.CancelOrderByUUIDParams) (orderV1.CancelOrderByUUIDRes, error) {
	o.storage.mu.RLock()
	order, exists := o.storage.orders[params.OrderUUID.String()]
	o.storage.mu.RUnlock()

	if !exists {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order not found",
		}, nil
	}

	switch order.Status {
	case orderV1.OrderDtoStatusPENDINGPAYMENT:
		o.storage.mu.Lock()
		order.Status = orderV1.OrderDtoStatusCANCELLED
		o.storage.mu.Unlock()
		log.Printf("Order %s cancelled successfully", params.OrderUUID.String())
		return &orderV1.CancelOrderByUUIDNoContent{}, nil
	case orderV1.OrderDtoStatusPAID:
		return &orderV1.Conflict{
			Code:    409,
			Message: "Order is already paid and cannot be cancelled",
		}, nil
	case orderV1.OrderDtoStatusCANCELLED:
		return &orderV1.CancelOrderByUUIDNoContent{}, nil
	default:
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "Invalid order status",
		}, nil
	}
}

func (o OrderHandler) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	log.Printf("Internal error: %v", err)
	return &orderV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(500),
			Message: orderV1.NewOptString("Internal error"),
		},
	}
}

func (o OrderHandler) Close() {
	if o.inventoryConn != nil {
		if err := o.inventoryConn.Close(); err != nil {
			log.Printf("Error closing inventory connection: %v", err)
		}
	}
	if o.paymentConn != nil {
		if err := o.paymentConn.Close(); err != nil {
			log.Printf("Error closing payment connection: %v", err)
		}
	}
}

func main() {
	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage)
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		orderHandler.Close()
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	defer orderHandler.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v", err)
	}

	log.Println("✅  Сервер остановлен")
}
