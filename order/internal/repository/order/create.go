package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	"github.com/YuraMishin/bigtechmicroservices/order/internal/repository/converter"
	repoModel "github.com/YuraMishin/bigtechmicroservices/order/internal/repository/model"
)

func (r *postgresRepository) CreateOrder(ctx context.Context, orderModel model.Order) (model.Order, error) {
	repoOrder := converter.ToRepoOrder(orderModel)

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
		          transaction_uuid, payment_method, status
	`

	var transactionUUIDParam interface{}
	if repoOrder.TransactionUUID == uuid.Nil {
		transactionUUIDParam = nil
	} else {
		transactionUUIDParam = repoOrder.TransactionUUID
	}

	rows, err := r.pool.Query(ctx, query,
		repoOrder.OrderUUID,
		repoOrder.UserUUID,
		partUUIDsJSON,
		repoOrder.TotalPrice,
		transactionUUIDParam,
		string(repoOrder.PaymentMethod),
		string(repoOrder.Status),
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to create order: %w", err)
	}
	defer rows.Close()

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.Order])
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to collect order row: %w", err)
	}

	return converter.ToModelOrder(row), nil
}
