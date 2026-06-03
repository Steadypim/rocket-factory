package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const grpcPort = 50051

// inventoryService реализует gRPC сервис, отвечающий за хранение и предоставление информации о деталях для сборки космических кораблей
type inventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*inventory_v1.Part
}

// GetPart Возвращает информацию о детали по её UUID.
func (s *inventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}

	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*inventory_v1.Part, 0)
	uuids := mergeFilters(req.GetUuids(), req.GetFilter().GetUuids())
	names := mergeFilters(req.GetNames(), req.GetFilter().GetNames())
	categories := mergeFilters(req.GetCategories(), req.GetFilter().GetCategories())
	manufacturerCountries := mergeFilters(req.GetManufacturerCountries(), req.GetFilter().GetManufacturerCountries())
	tags := mergeFilters(req.GetTags(), req.GetFilter().GetTags())

	for _, part := range s.parts {
		if !matchesAny(uuids, part.GetUuid()) {
			continue
		}

		if !matchesAny(names, part.GetName()) {
			continue
		}

		if !matchesAny(categories, part.GetCategory()) {
			continue
		}

		if !matchesAny(manufacturerCountries, part.GetManufacturer().GetCountry()) {
			continue
		}

		if !matchesTags(part, tags) {
			continue
		}

		result = append(result, part)
	}

	return &inventory_v1.ListPartsResponse{
		Parts: result,
	}, nil
}

func mergeFilters[T comparable](topLevel []T, nested []T) []T {
	if len(nested) == 0 {
		return topLevel
	}

	if len(topLevel) == 0 {
		return nested
	}

	merged := make([]T, 0, len(topLevel)+len(nested))
	merged = append(merged, topLevel...)
	merged = append(merged, nested...)
	return merged
}

func matchesAny[T comparable](filter []T, value T) bool {
	if len(filter) == 0 {
		return true
	}

	return slices.Contains(filter, value)
}

func matchesTags(part *inventory_v1.Part, tags []string) bool {
	if len(tags) == 0 {
		return true
	}

	partTags := make(map[string]struct{}, len(part.GetTags()))
	for _, tag := range part.GetTags() {
		partTags[tag] = struct{}{}
	}

	for _, tag := range tags {
		if _, ok := partTags[tag]; ok {
			return true
		}
	}

	return false
}

func seedParts() map[string]*inventory_v1.Part {
	parts := []*inventory_v1.Part{
		{
			Uuid:          "11111111-1111-1111-1111-111111111111",
			Name:          "Main engine",
			Description:   "Основной двигатель ракеты",
			Price:         1500000,
			StockQuantity: 4,
			Category:      inventory_v1.Category_ENGINE,
			Dimensions: &inventory_v1.Dimensions{
				Length: 250,
				Width:  120,
				Height: 120,
				Weight: 850,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "Rocket Dynamics",
				Country: "USA",
				Website: "https://rocket-dynamics.example.com",
			},
			Tags: []string{"engine", "rocket", "heavy"},
			Metadata: map[string]*inventory_v1.Value{
				"fuel_type": {
					Value: &inventory_v1.Value_StringValue{
						StringValue: "kerosene",
					},
				},
				"reusable": {
					Value: &inventory_v1.Value_BoolValue{
						BoolValue: true,
					},
				},
				"thrust_kn": {
					Value: &inventory_v1.Value_Int64Value{
						Int64Value: 7600,
					},
				},
			},
		},
		{
			Uuid:          "22222222-2222-2222-2222-222222222222",
			Name:          "Fuel tank",
			Description:   "Топливный бак первой ступени",
			Price:         420000,
			StockQuantity: 8,
			Category:      inventory_v1.Category_FUEL,
			Dimensions: &inventory_v1.Dimensions{
				Length: 600,
				Width:  200,
				Height: 200,
				Weight: 1200,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "Orbital Parts",
				Country: "Germany",
				Website: "https://orbital-parts.example.com",
			},
			Tags: []string{"fuel", "tank", "stage-1"},
			Metadata: map[string]*inventory_v1.Value{
				"material": {
					Value: &inventory_v1.Value_StringValue{
						StringValue: "aluminum-lithium",
					},
				},
				"capacity_liters": {
					Value: &inventory_v1.Value_Int64Value{
						Int64Value: 50000,
					},
				},
			},
		},
		{
			Uuid:          "33333333-3333-3333-3333-333333333333",
			Name:          "Porthole glass",
			Description:   "Иллюминатор из многослойного стекла",
			Price:         75000,
			StockQuantity: 20,
			Category:      inventory_v1.Category_PORTHOLE,
			Dimensions: &inventory_v1.Dimensions{
				Length: 80,
				Width:  80,
				Height: 12,
				Weight: 35,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "Luna Glass",
				Country: "Japan",
				Website: "https://luna-glass.example.com",
			},
			Tags: []string{"porthole", "glass", "crew"},
			Metadata: map[string]*inventory_v1.Value{
				"radiation_protection": {
					Value: &inventory_v1.Value_BoolValue{
						BoolValue: true,
					},
				},
				"layers": {
					Value: &inventory_v1.Value_Int64Value{
						Int64Value: 5,
					},
				},
			},
		},
		{
			Uuid:          "44444444-4444-4444-4444-444444444444",
			Name:          "Stabilizer wing",
			Description:   "Стабилизирующее крыло ракеты",
			Price:         210000,
			StockQuantity: 12,
			Category:      inventory_v1.Category_WING,
			Dimensions: &inventory_v1.Dimensions{
				Length: 320,
				Width:  90,
				Height: 25,
				Weight: 180,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "Cosmo Engineering",
				Country: "Russia",
				Website: "https://cosmo-engineering.example.com",
			},
			Tags: []string{"wing", "stabilizer", "aero"},
			Metadata: map[string]*inventory_v1.Value{
				"heat_resistant": {
					Value: &inventory_v1.Value_BoolValue{
						BoolValue: true,
					},
				},
				"max_temperature_c": {
					Value: &inventory_v1.Value_Int64Value{
						Int64Value: 1800,
					},
				},
			},
		},
	}

	result := make(map[string]*inventory_v1.Part, len(parts))
	for _, part := range parts {
		result[part.GetUuid()] = part
	}

	return result
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return
	}

	// Создаем grpc сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &inventoryService{
		parts: seedParts(),
	}

	inventory_v1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		slog.Info("🚀 gRPC server listening", "port", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			slog.Error("failed to serve", "error", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	slog.Info("✅ Server stopped")
}
