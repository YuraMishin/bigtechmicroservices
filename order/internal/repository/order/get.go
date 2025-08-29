package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	repoConverter "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
	repoModel "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/model"
)

func (r *postgresRepository) GetOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	query := `
		SELECT order_uuid, user_uuid, part_uuids, total_price, 
		       transaction_uuid, payment_method, status, 
		       created_at, updated_at, cancelled_at
		FROM orders 
		WHERE order_uuid = $1
	`

	var (
		orderUUIDResult  uuid.UUID
		userUUID         uuid.UUID
		partUUIDsJSONB   json.RawMessage
		totalPrice       float32
		transactionUUID  *uuid.UUID
		paymentMethodStr string
		statusStr        string
		createdAt        time.Time
		updatedAt        time.Time
		cancelledAt      *time.Time
	)

	err := r.pool.QueryRow(ctx, query, orderUUID).Scan(
		&orderUUIDResult, &userUUID, &partUUIDsJSONB, &totalPrice,
		&transactionUUID, &paymentMethodStr, &statusStr,
		&createdAt, &updatedAt, &cancelledAt,
	)
	if err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	// Конвертируем JSON обратно в []uuid.UUID
	var partUUIDs []uuid.UUID
	if err := json.Unmarshal(partUUIDsJSONB, &partUUIDs); err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal part UUIDs: %w", err)
	}

	// Конвертируем строки обратно в enum типы
	paymentMethod, err := parsePaymentMethod(paymentMethodStr)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to parse payment method: %w", err)
	}

	status, err := parseOrderStatus(statusStr)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to parse order status: %w", err)
	}

	// Создаем результирующий объект
	repoOrder := repoModel.Order{
		OrderUUID:       orderUUIDResult,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      totalPrice,
		TransactionUUID: uuid.Nil, // По умолчанию
		PaymentMethod:   paymentMethod,
		Status:          status,
	}

	if transactionUUID != nil {
		repoOrder.TransactionUUID = *transactionUUID
	}

	return repoConverter.ToModelOrder(repoOrder), nil
}
