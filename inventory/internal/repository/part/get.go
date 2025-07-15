package inventory

import (
	"context"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoConverter "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	repoPart, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}
	return repoConverter.PartToModel(repoPart), nil
}
