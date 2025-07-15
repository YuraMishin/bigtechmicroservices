package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/YuraMishin/bigtechmicroservices/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
)

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) initializeSampleData() {
	s.mu.Lock()
	defer s.mu.Unlock()
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
		s.parts[partUUID] = &inventoryV1.Part{
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
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func containsCategory(slice []inventoryV1.Category, item inventoryV1.Category) bool {
	for _, c := range slice {
		if c == item {
			return true
		}
	}
	return false
}

func matchesUUIDFilter(part *inventoryV1.Part, uuids []string) bool {
	if len(uuids) == 0 {
		return true
	}
	return slices.Contains(uuids, part.Uuid)
}

func matchesNameFilter(part *inventoryV1.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	return slices.Contains(names, part.Name)
}

func matchesCategoryFilter(part *inventoryV1.Part, categories []inventoryV1.Category) bool {
	if len(categories) == 0 {
		return true
	}
	return containsCategory(categories, part.Category)
}

func matchesManufacturerCountryFilter(part *inventoryV1.Part, countries []string) bool {
	if len(countries) == 0 {
		return true
	}
	if part.Manufacturer == nil {
		return false
	}
	return slices.Contains(countries, part.Manufacturer.Country)
}

func matchesTagsFilter(part *inventoryV1.Part, filterTags []string) bool {
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

func (s *inventoryService) matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	if filter == nil {
		return true
	}

	return matchesUUIDFilter(part, filter.Uuids) &&
		matchesNameFilter(part, filter.Names) &&
		matchesCategoryFilter(part, filter.Categories) &&
		matchesManufacturerCountryFilter(part, filter.ManufacturerCountries) &&
		matchesTagsFilter(part, filter.Tags)
}

func (s *inventoryService) GetPart(ctx context.Context, in *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err := uuid.Parse(in.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid UUID format: %v", err)
	}

	part, exists := s.parts[in.Uuid]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Part with UUID %s not found", in.Uuid)
	}

	return &inventoryV1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(ctx context.Context, in *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var filteredParts []*inventoryV1.Part
	for _, part := range s.parts {
		if s.matchesFilter(part, in.GetFilter()) {
			filteredParts = append(filteredParts, part)
		}
	}
	return &inventoryV1.ListPartsResponse{
		Parts: filteredParts,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	reflection.Register(s)
	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}
	service.initializeSampleData()
	inventoryV1.RegisterInventoryServiceServer(s, service)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
