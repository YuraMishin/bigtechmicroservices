package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
	repoModel "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/model"
)

func (r *postgresRepository) CreateOrder(ctx context.Context, orderModel model.Order) (model.Order, error) {
	repoOrder := converter.ToRepoOrder(orderModel)

	// Конвертируем []uuid.UUID в JSON для хранения в JSONB
	partUUIDsJSON, err := json.Marshal(repoOrder.PartUuids)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to marshal part UUIDs: %w", err)
	}

	query := `
		INSERT INTO orders (
			order_uuid, user_uuid, part_uuids, total_price, 
			transaction_uuid, payment_method, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING order_uuid, user_uuid, part_uuids, total_price, 
		          transaction_uuid, payment_method, status, 
		          created_at, updated_at, cancelled_at
	`

	var (
		orderUUID       uuid.UUID
		userUUID        uuid.UUID
		partUUIDsJSONB  json.RawMessage
		totalPrice      float32
		transactionUUID *uuid.UUID
		paymentMethod   string
		status          string
		createdAt       time.Time
		updatedAt       time.Time
		cancelledAt     *time.Time
	)

	// Обрабатываем transaction_uuid - если это uuid.Nil, передаем NULL
	var transactionUUIDParam interface{}
	if repoOrder.TransactionUUID == uuid.Nil {
		transactionUUIDParam = nil
	} else {
		transactionUUIDParam = repoOrder.TransactionUUID
	}

	err = r.pool.QueryRow(ctx, query,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		partUUIDsJSON,
		repoOrder.TotalPrice,
		transactionUUIDParam,
		string(repoOrder.PaymentMethod),
		string(repoOrder.Status),
	).Scan(
		&orderUUID, &userUUID, &partUUIDsJSONB, &totalPrice,
		&transactionUUID, &paymentMethod, &status,
		&createdAt, &updatedAt, &cancelledAt,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	// Конвертируем JSON обратно в []uuid.UUID
	var partUUIDs []uuid.UUID
	if err := json.Unmarshal(partUUIDsJSONB, &partUUIDs); err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal part UUIDs: %w", err)
	}

	// Создаем результирующий объект
	resultOrder := repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      totalPrice,
		TransactionUUID: uuid.Nil, // Будет заполнено при оплате
		PaymentMethod:   repoOrder.PaymentMethod,
		Status:          repoOrder.Status,
	}

	if transactionUUID != nil {
		resultOrder.TransactionUUID = *transactionUUID
	}

	return converter.ToModelOrder(resultOrder), nil
}
