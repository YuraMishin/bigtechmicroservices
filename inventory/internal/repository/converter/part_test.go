package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	serviceModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func TestToModel_Full(t *testing.T) {
	repo := repoModel.Part{
		UUID:          "id",
		Name:          "n",
		Description:   "d",
		Price:         10,
		StockQuantity: 2,
		Category:      repoModel.Category(3),
		Dimensions:    &repoModel.Dimensions{Length: 1, Width: 2, Height: 3, Weight: 4},
		Manufacturer:  &repoModel.Manufacturer{Name: "m", Country: "c", Website: "w"},
		Tags:          []string{"t1", "t2"},
		Metadata: map[string]*repoModel.Value{
			"a": {StringValue: ptr("s")},
			"b": {Int64Value: ptr[int64](5)},
			"c": {DoubleValue: ptr(1.5)},
			"d": {BoolValue: ptr(true)},
		},
		CreatedAt: 1,
		UpdatedAt: 2,
	}

	m := ToModel(repo)
	assert.Equal(t, serviceModel.Category(3), m.Category)
	assert.Equal(t, repo.UUID, m.UUID)
	assert.Equal(t, repo.Dimensions.Length, m.Dimensions.Length)
	assert.Equal(t, repo.Manufacturer.Country, m.Manufacturer.Country)
	assert.Equal(t, *repo.Metadata["a"].StringValue, *m.Metadata["a"].StringValue)
	assert.Equal(t, *repo.Metadata["b"].Int64Value, *m.Metadata["b"].Int64Value)
	assert.Equal(t, *repo.Metadata["c"].DoubleValue, *m.Metadata["c"].DoubleValue)
	assert.Equal(t, *repo.Metadata["d"].BoolValue, *m.Metadata["d"].BoolValue)
}

func TestFromProto_Full(t *testing.T) {
	p := &inventoryV1.Part{
		Uuid:          "id",
		Name:          "n",
		Description:   "d",
		Price:         10,
		StockQuantity: 2,
		Category:      inventoryV1.Category(4),
		Dimensions:    &inventoryV1.Dimensions{Length: 1, Width: 2, Height: 3, Weight: 4},
		Manufacturer:  &inventoryV1.Manufacturer{Name: "m", Country: "c", Website: "w"},
		Tags:          []string{"t"},
		Metadata: map[string]*inventoryV1.Value{
			"a": {Kind: &inventoryV1.Value_StringValue{StringValue: "s"}},
			"b": {Kind: &inventoryV1.Value_Int64Value{Int64Value: 5}},
			"c": {Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 1.5}},
			"d": {Kind: &inventoryV1.Value_BoolValue{BoolValue: true}},
		},
		CreatedAt: timestamppb.New(time.Unix(1, 0)),
		UpdatedAt: timestamppb.New(time.Unix(2, 0)),
	}

	r := FromProto(p)
	assert.Equal(t, repoModel.Category(4), r.Category)
	assert.Equal(t, p.Uuid, r.UUID)
	assert.Equal(t, p.Dimensions.Length, r.Dimensions.Length)
	assert.Equal(t, p.Manufacturer.Country, r.Manufacturer.Country)
	assert.Equal(t, "s", *r.Metadata["a"].StringValue)
	assert.Equal(t, int64(5), *r.Metadata["b"].Int64Value)
	assert.Equal(t, 1.5, *r.Metadata["c"].DoubleValue)
	assert.Equal(t, true, *r.Metadata["d"].BoolValue)
}

func TestToRepoPartsFilter(t *testing.T) {
	f := serviceModel.PartsFilter{
		UUIDs:                 []string{"id"},
		Names:                 []string{"n"},
		Categories:            []serviceModel.Category{serviceModel.CategoryEngine, serviceModel.CategoryFuel},
		ManufacturerCountries: []string{"c"},
		Tags:                  []string{"t"},
	}
	r := ToRepoPartsFilter(f)
	assert.Equal(t, []repoModel.Category{repoModel.Category(serviceModel.CategoryEngine), repoModel.Category(serviceModel.CategoryFuel)}, r.Categories)
	assert.Equal(t, f.UUIDs, r.UUIDs)
	assert.Equal(t, f.Names, r.Names)
	assert.Equal(t, f.ManufacturerCountries, r.ManufacturerCountries)
	assert.Equal(t, f.Tags, r.Tags)
}

func ptr[T any](v T) *T { return &v }
