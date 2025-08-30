package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
)

func (r *postgresRepository) UpdateOrder(ctx context.Context, orderModel model.Order) error {
	repoOrder := converter.ToRepoOrder(orderModel)

	// Конвертируем []uuid.UUID в JSON для хранения в JSONB
	partUUIDsJSON, err := json.Marshal(repoOrder.PartUuids)
	if err != nil {
		return fmt.Errorf("failed to marshal part UUIDs: %w", err)
	}

	query := `
		UPDATE orders 
		SET user_uuid = $2, part_uuids = $3, total_price = $4, 
		    transaction_uuid = $5, payment_method = $6, status = $7
		WHERE order_uuid = $1
	`

	// Обрабатываем transaction_uuid - если это uuid.Nil, передаем NULL
	var transactionUUIDParam interface{}
	if repoOrder.TransactionUUID == uuid.Nil {
		transactionUUIDParam = nil
	} else {
		transactionUUIDParam = repoOrder.TransactionUUID
	}

	result, err := r.pool.Exec(ctx, query,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		partUUIDsJSON,
		repoOrder.TotalPrice,
		transactionUUIDParam,
		string(repoOrder.PaymentMethod),
		string(repoOrder.Status),
	)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
