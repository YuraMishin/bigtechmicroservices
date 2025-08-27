package inventory

import (
	"sync"

	def "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: initializeSampleData(),
	}
}
