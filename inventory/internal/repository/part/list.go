package inventory

import (
	"context"
	"slices"

	serviceModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/model"
	repoConverter "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/converter"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filter serviceModel.PartsFilter) ([]serviceModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoFilter := repoConverter.ToRepoPartsFilter(filter)

	var filteredParts []repoModel.Part
	for _, part := range r.data {
		if r.matchesFilter(part, repoFilter) {
			filteredParts = append(filteredParts, part)
		}
	}
	var serviceParts []serviceModel.Part
	for _, part := range filteredParts {
		serviceParts = append(serviceParts, repoConverter.PartToModel(part))
	}
	return serviceParts, nil
}

func containsCategory(slice []repoModel.Category, item repoModel.Category) bool {
	for _, c := range slice {
		if c == item {
			return true
		}
	}
	return false
}

func matchesUUIDFilter(part repoModel.Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}
	return slices.Contains(uuids, part.UUID)
}

func matchesNameFilter(part repoModel.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	return slices.Contains(names, part.Name)
}

func matchesCategoryFilter(part repoModel.Part, categories []repoModel.Category) bool {
	if len(categories) == 0 {
		return true
	}
	return containsCategory(categories, part.Category)
}

func matchesManufacturerCountryFilter(part repoModel.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	if part.Manufacturer == nil {
		return false
	}
	return slices.Contains(countries, part.Manufacturer.Country)
}

func matchesTagsFilter(part repoModel.Part, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, filterTag := range filterTags {
		if slices.Contains(part.Tags, filterTag) {
			return true
		}
	}
	return false
}

func (r *repository) matchesFilter(part repoModel.Part, filter repoModel.PartsFilter) bool {
	return matchesUUIDFilter(part, filter.UUIDs) &&
		matchesNameFilter(part, filter.Names) &&
		matchesCategoryFilter(part, filter.Categories) &&
		matchesManufacturerCountryFilter(part, filter.ManufacturerCountries) &&
		matchesTagsFilter(part, filter.Tags)
}
