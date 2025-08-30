package order

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/YuraMishin/bigtechmicroservices/order/internal/repository"
)

var _ def.OrderRepository = (*postgresRepository)(nil)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) (*postgresRepository, error) {
	if pool == nil {
		return nil, errors.New("orderRepository is nil")
	}

	return &postgresRepository{
		pool: pool,
	}, nil
}
