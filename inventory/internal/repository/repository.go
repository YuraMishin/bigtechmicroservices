package repository

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter repoModel.PartsFilter) ([]model.Part, error)
}
