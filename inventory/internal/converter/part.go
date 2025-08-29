package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	partV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func ToProtoPart(part model.Part) *partV1.Part {
	var dimensions *partV1.Dimensions
	if part.Dimensions != nil {
		dimensions = &partV1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}

	var manufacturer *partV1.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = &partV1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*partV1.Value)
	for k, v := range part.Metadata {
		value := &partV1.Value{}
		switch {
		case v.StringValue != nil:
			value.Kind = &partV1.Value_StringValue{StringValue: *v.StringValue}
		case v.Int64Value != nil:
			value.Kind = &partV1.Value_Int64Value{Int64Value: *v.Int64Value}
		case v.DoubleValue != nil:
			value.Kind = &partV1.Value_DoubleValue{DoubleValue: *v.DoubleValue}
		case v.BoolValue != nil:
			value.Kind = &partV1.Value_BoolValue{BoolValue: *v.BoolValue}
		}
		metadata[k] = value
	}

	return &partV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      partV1.Category(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          append([]string{}, part.Tags...),
		Metadata:      metadata,
		CreatedAt:     &timestamppb.Timestamp{Seconds: part.CreatedAt},
		UpdatedAt:     &timestamppb.Timestamp{Seconds: part.UpdatedAt},
	}
}

func ToProtoPartList(parts []model.Part) []*partV1.Part {
	protoParts := make([]*partV1.Part, len(parts))
	for idx, part := range parts {
		protoParts[idx] = ToProtoPart(part)
	}
	return protoParts
}

func ToModelPartsFilter(partsFilter *partV1.PartsFilter) model.PartsFilter {
	if partsFilter == nil {
		return model.PartsFilter{}
	}

	categories := make([]model.Category, len(partsFilter.Categories))
	for i, category := range partsFilter.Categories {
		categories[i] = model.Category(category)
	}

	return model.PartsFilter{
		UUIDs:                 partsFilter.Uuids,
		Names:                 partsFilter.Names,
		Categories:            categories,
		ManufacturerCountries: partsFilter.ManufacturerCountries,
		Tags:                  partsFilter.Tags,
	}
}
