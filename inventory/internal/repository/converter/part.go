package converter

import (
	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func PartToModel(part repoModel.Part) model.Part {
	var dimensions *model.Dimensions
	if part.Dimensions != nil {
		dimensions = &model.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}
	var manufacturer *model.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = &model.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}
	metadata := make(map[string]*model.Value)
	for k, v := range part.Metadata {
		metadata[k] = &model.Value{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
		}
	}
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          append([]string{}, part.Tags...),
		Metadata:      metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func PartFromProto(part *inventoryV1.Part) repoModel.Part {
	var dimensions *repoModel.Dimensions
	if part.Dimensions != nil {
		dimensions = &repoModel.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}
	var manufacturer *repoModel.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = &repoModel.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}
	metadata := make(map[string]*repoModel.Value)
	for k, v := range part.Metadata {
		value := &repoModel.Value{}
		switch x := v.Kind.(type) {
		case *inventoryV1.Value_StringValue:
			s := x.StringValue
			value.StringValue = &s
		case *inventoryV1.Value_Int64Value:
			i := x.Int64Value
			value.Int64Value = &i
		case *inventoryV1.Value_DoubleValue:
			d := x.DoubleValue
			value.DoubleValue = &d
		case *inventoryV1.Value_BoolValue:
			b := x.BoolValue
			value.BoolValue = &b
		}
		metadata[k] = value
	}
	return repoModel.Part{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoModel.Category(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          append([]string{}, part.Tags...),
		Metadata:      metadata,
		CreatedAt:     part.CreatedAt.Seconds,
		UpdatedAt:     part.UpdatedAt.Seconds,
	}
}

func ToRepoPartsFilter(filter model.PartsFilter) repoModel.PartsFilter {
	categories := make([]repoModel.Category, len(filter.Categories))
	for i, c := range filter.Categories {
		categories[i] = repoModel.Category(c)
	}
	return repoModel.PartsFilter{
		UUIDs:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
