package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	now := timestamppb.Now()

	engineUUID := uuid.New().String()
	fuelUUID := uuid.New().String()
	wingUUID := uuid.New().String()
	portholeUUID := uuid.New().String()

	s.parts[engineUUID] = &inventoryV1.Part{
		Uuid:          engineUUID,
		Name:          "Main Engine",
		Description:   "Primary propulsion engine for spacecraft",
		Price:         50000.0,
		StockQuantity: 5,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 2.5,
			Width:  1.2,
			Height: 1.8,
			Weight: 1500.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "SpaceTech Industries",
			Country: "Germany",
			Website: "https://spacetech.de",
		},
		Tags: []string{"propulsion", "main", "high-thrust"},
		Metadata: map[string]*inventoryV1.Value{
			"thrust": {
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 1000000.0},
			},
			"fuel_type": {
				Kind: &inventoryV1.Value_StringValue{StringValue: "liquid_hydrogen"},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.parts[fuelUUID] = &inventoryV1.Part{
		Uuid:          fuelUUID,
		Name:          "Fuel Tank",
		Description:   "Cryogenic fuel storage tank",
		Price:         25000.0,
		StockQuantity: 10,
		Category:      inventoryV1.Category_CATEGORY_FUEL,
		Dimensions: &inventoryV1.Dimensions{
			Length: 3.0,
			Width:  2.0,
			Height: 2.0,
			Weight: 800.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "CryoSystems Ltd",
			Country: "USA",
			Website: "https://cryosystems.com",
		},
		Tags: []string{"storage", "cryogenic", "fuel"},
		Metadata: map[string]*inventoryV1.Value{
			"capacity": {
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 5000.0},
			},
			"temperature": {
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: -253.0},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.parts[wingUUID] = &inventoryV1.Part{
		Uuid:          wingUUID,
		Name:          "Main Wing",
		Description:   "Primary wing assembly for atmospheric flight",
		Price:         35000.0,
		StockQuantity: 3,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Dimensions: &inventoryV1.Dimensions{
			Length: 8.0,
			Width:  2.5,
			Height: 0.3,
			Weight: 1200.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "AeroDynamics Corp",
			Country: "France",
			Website: "https://aerodynamics.fr",
		},
		Tags: []string{"aerodynamics", "main", "wing"},
		Metadata: map[string]*inventoryV1.Value{
			"span": {
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 16.0},
			},
			"material": {
				Kind: &inventoryV1.Value_StringValue{StringValue: "carbon_fiber"},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.parts[portholeUUID] = &inventoryV1.Part{
		Uuid:          portholeUUID,
		Name:          "Observation Porthole",
		Description:   "Reinforced glass porthole for space observation",
		Price:         15000.0,
		StockQuantity: 8,
		Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 0.8,
			Width:  0.8,
			Height: 0.1,
			Weight: 50.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "OptiGlass GmbH",
			Country: "Germany",
			Website: "https://optiglass.de",
		},
		Tags: []string{"observation", "glass", "window"},
		Metadata: map[string]*inventoryV1.Value{
			"thickness": {
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 0.05},
			},
			"transparency": {
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: 0.95},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

//
//nolint:unused
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

//
//nolint:unused
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
	return containsString(uuids, part.Uuid)
}

func matchesNameFilter(part *inventoryV1.Part, names []string) bool {
	if len(names) == 0 {
		return true
	}
	return containsString(names, part.Name)
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
	return containsString(countries, part.Manufacturer.Country)
}

func matchesTagsFilter(part *inventoryV1.Part, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	for _, filterTag := range filterTags {
		if containsString(part.Tags, filterTag) {
			return true
		}
	}
	return false
}

func (s *inventoryService) matchesFilter(part *inventoryV1.Part, filter *inventoryV1.PartsFilter) bool {
	// If no filter is provided, return true
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

	if in.Uuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "UUID cannot be empty")
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

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	service.initializeSampleData()

	inventoryV1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

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
