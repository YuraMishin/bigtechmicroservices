package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	serviceModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
)

func TestNewRepository_InitData(t *testing.T) {
	r := NewRepository()
	assert.NotNil(t, r)
	assert.Greater(t, len(r.data), 0)
}

func TestRepository_GetPart_Found(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()

	var anyUUID string
	for k := range r.data {
		anyUUID = k
		break
	}

	p, err := r.GetPart(ctx, anyUUID)
	assert.NoError(t, err)
	assert.Equal(t, anyUUID, p.UUID)
}

func TestRepository_GetPart_NotFound(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()

	p, err := r.GetPart(ctx, "does-not-exist")
	assert.Error(t, err)
	assert.ErrorIs(t, err, serviceModel.ErrPartNotFound)
	assert.Equal(t, serviceModel.Part{}, p)
}

func TestRepository_ListParts_NoFilter(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()

	parts, err := r.ListParts(ctx, serviceModel.PartsFilter{})
	assert.NoError(t, err)
	assert.Equal(t, len(r.data), len(parts))
}

func TestRepository_ListParts_FilterByUUID(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()

	var uuids []string
	for k := range r.data {
		uuids = append(uuids, k)
		break
	}

	parts, err := r.ListParts(ctx, serviceModel.PartsFilter{UUIDs: uuids})
	assert.NoError(t, err)
	assert.Len(t, parts, 1)
	assert.Equal(t, uuids[0], parts[0].UUID)
}

func TestRepository_ListParts_FilterByTagsAndCountry(t *testing.T) {
	r := NewRepository()
	ctx := context.Background()

	var country string
	var tag string
	var found bool
	for _, v := range r.data {
		if v.Manufacturer != nil && len(v.Tags) > 0 {
			country = v.Manufacturer.Country
			tag = v.Tags[0]
			found = true
			break
		}
	}
	if !found {
		t.Skip("no sample data with manufacturer and tags")
	}

	parts, err := r.ListParts(ctx, serviceModel.PartsFilter{
		ManufacturerCountries: []string{country},
		Tags:                  []string{tag},
	})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(parts), 1)
	for _, p := range parts {
		assert.Equal(t, country, p.Manufacturer.Country)
		assert.Contains(t, p.Tags, tag)
	}
}
