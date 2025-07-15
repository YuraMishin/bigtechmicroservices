package inventory

import (
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	repoConverter "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/converter"
	repoModel "github.com/YuraMishin/bigtechmicroservices/inventory/internal/repository/model"
	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

func initializeSampleData() map[string]repoModel.Part {
	data := make(map[string]repoModel.Part)
	if err := gofakeit.Seed(42); err != nil {
		log.Printf("gofakeit.Seed error: %v", err)
	}
	for i := 0; i < 10; i++ {
		partUUID := uuid.New().String()
		now := timestamppb.Now()
		categories := []inventoryV1.Category{
			inventoryV1.Category_CATEGORY_ENGINE,
			inventoryV1.Category_CATEGORY_FUEL,
			inventoryV1.Category_CATEGORY_WING,
			inventoryV1.Category_CATEGORY_PORTHOLE,
		}
		category := categories[gofakeit.IntRange(0, len(categories)-1)]
		countries := []string{"USA", "Germany", "France", "Japan", "UK", "Canada", "Italy", "Sweden"}
		country := countries[gofakeit.IntRange(0, len(countries)-1)]
		var tags []string
		switch category {
		case inventoryV1.Category_CATEGORY_ENGINE:
			tags = []string{"propulsion", "thrust", "combustion", "high-performance", "main-engine"}
		case inventoryV1.Category_CATEGORY_FUEL:
			tags = []string{"storage", "cryogenic", "fuel", "tank", "pressurized"}
		case inventoryV1.Category_CATEGORY_WING:
			tags = []string{"aerodynamics", "lift", "control", "carbon-fiber", "wing"}
		case inventoryV1.Category_CATEGORY_PORTHOLE:
			tags = []string{"observation", "glass", "window", "reinforced", "transparent"}
		}
		gofakeit.ShuffleStrings(tags)
		selectedTags := tags[:gofakeit.IntRange(2, min(4, len(tags)))]
		metadata := make(map[string]*inventoryV1.Value)
		switch category {
		case inventoryV1.Category_CATEGORY_ENGINE:
			metadata["thrust"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: gofakeit.Float64Range(500000, 2000000)},
			}
			metadata["fuel_type"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_StringValue{StringValue: gofakeit.RandomString([]string{"liquid_hydrogen", "kerosene", "methane", "solid_fuel"})},
			}
		case inventoryV1.Category_CATEGORY_FUEL:
			metadata["capacity"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: gofakeit.Float64Range(1000, 10000)},
			}
			metadata["temperature"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: gofakeit.Float64Range(-260, -200)},
			}
		case inventoryV1.Category_CATEGORY_WING:
			metadata["span"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: gofakeit.Float64Range(10, 25)},
			}
			metadata["material"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_StringValue{StringValue: gofakeit.RandomString([]string{"carbon_fiber", "titanium", "aluminum", "composite"})},
			}
		case inventoryV1.Category_CATEGORY_PORTHOLE:
			metadata["thickness"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: gofakeit.Float64Range(0.02, 0.1)},
			}
			metadata["transparency"] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: gofakeit.Float64Range(0.85, 0.98)},
			}
		}
		var companyName string
		switch country {
		case "USA":
			companyName = gofakeit.RandomString([]string{"SpaceTech USA", "AeroDynamics Inc", "RocketCorp", "Orbital Systems"})
		case "Germany":
			companyName = gofakeit.RandomString([]string{"Deutsche SpaceTech", "Bavarian Aerospace", "Berlin Dynamics", "Hamburg Systems"})
		case "France":
			companyName = gofakeit.RandomString([]string{"Aérospatiale France", "Paris Dynamics", "Lyon Aerospace", "Marseille Tech"})
		case "Japan":
			companyName = gofakeit.RandomString([]string{"Tokyo SpaceTech", "Osaka Aerospace", "Kyoto Dynamics", "Yokohama Systems"})
		default:
			companyName = gofakeit.Company() + " Aerospace"
		}
		protoPart := &inventoryV1.Part{
			Uuid: partUUID,
			Name: gofakeit.RandomString([]string{
				"Main Engine", "Auxiliary Engine", "Thruster", "Propulsion Unit",
				"Fuel Tank", "Cryogenic Tank", "Storage Container", "Fuel Cell",
				"Main Wing", "Control Surface", "Aileron", "Flap Assembly",
				"Observation Window", "Porthole", "Viewport", "Transparent Panel",
			}) + " " + gofakeit.RandomString([]string{"Alpha", "Beta", "Gamma", "Delta", "Echo"}),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Float64Range(5000, 100000),
			StockQuantity: int64(gofakeit.IntRange(1, 20)),
			Category:      category,
			Dimensions: &inventoryV1.Dimensions{
				Length: gofakeit.Float64Range(0.5, 10.0),
				Width:  gofakeit.Float64Range(0.3, 5.0),
				Height: gofakeit.Float64Range(0.1, 3.0),
				Weight: gofakeit.Float64Range(50, 2000),
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    companyName,
				Country: country,
				Website: "https://" + gofakeit.DomainName(),
			},
			Tags:      selectedTags,
			Metadata:  metadata,
			CreatedAt: now,
			UpdatedAt: now,
		}
		data[partUUID] = repoConverter.PartFromProto(protoPart)
	}
	return data
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
