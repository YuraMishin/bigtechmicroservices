package converter

import (
	"github.com/YuraMishin/bigtechmicroservices/order/internal/model"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func ToModelParts(parts []*inventoryV1.Part) []model.Part {
	partsModel := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		partsModel = append(partsModel, ToModelPart(part))
	}
	return partsModel
}

func ToModelPart(part *inventoryV1.Part) model.Part {
	var dimensions *model.Dimensions
	if part.Dimensions != nil {
		dims := ToModelDimensions(part.Dimensions)
		dimensions = &dims
	}
	var manufacturer *model.Manufacturer
	if part.Manufacturer != nil {
		man := ToModelManufacturer(part.Manufacturer)
		manufacturer = &man
	}
	metadata := make(map[string]*model.Value)
	for k, v := range part.Metadata {
		if v == nil || v.Kind == nil {
			continue
		}
		value := &model.Value{}
		switch x := v.Kind.(type) {
		case *inventoryV1.Value_StringValue:
			value.StringValue = &x.StringValue
		case *inventoryV1.Value_Int64Value:
			value.Int64Value = &x.Int64Value
		case *inventoryV1.Value_DoubleValue:
			value.DoubleValue = &x.DoubleValue
		case *inventoryV1.Value_BoolValue:
			value.BoolValue = &x.BoolValue
		}
		metadata[k] = value
	}
	return model.Part{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      metadata,
		CreatedAt:     part.CreatedAt.Seconds,
		UpdatedAt:     part.UpdatedAt.Seconds,
	}
}

func ToModelDimensions(dimensions *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ToModelManufacturer(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToProtoPartsFilter(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}
	return &inventoryV1.PartsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
