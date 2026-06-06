package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type partServiceStub struct {
	part  *model.Part
	parts []model.Part
	err   error
}

func (s partServiceStub) Get(_ context.Context, _ string) (*model.Part, error) {
	return s.part, s.err
}

func (s partServiceStub) List(_ context.Context, _ model.PartsFilter) ([]model.Part, error) {
	return s.parts, s.err
}

func TestGetPartReturnsProtoPart(t *testing.T) {
	t.Parallel()

	api := NewAPI(partServiceStub{part: &model.Part{UUID: "part-1", Name: "Engine"}})

	resp, err := api.GetPart(context.Background(), &inventoryv1.GetPartRequest{Uuid: "part-1"})
	if err != nil {
		t.Fatalf("GetPart returned error: %v", err)
	}
	if resp.GetPart().GetUuid() != "part-1" {
		t.Fatalf("uuid = %q, want part-1", resp.GetPart().GetUuid())
	}
}

func TestGetPartMapsNotFoundToGrpcCode(t *testing.T) {
	t.Parallel()

	api := NewAPI(partServiceStub{err: model.ErrPartNotFound})

	_, err := api.GetPart(context.Background(), &inventoryv1.GetPartRequest{Uuid: "missing"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("code = %v, want NotFound", status.Code(err))
	}
}

func TestListPartsReturnsProtoParts(t *testing.T) {
	t.Parallel()

	api := NewAPI(partServiceStub{parts: []model.Part{{UUID: "part-1"}, {UUID: "part-2"}}})

	resp, err := api.ListParts(context.Background(), &inventoryv1.ListPartsRequest{})
	if err != nil {
		t.Fatalf("ListParts returned error: %v", err)
	}
	if len(resp.GetParts()) != 2 {
		t.Fatalf("parts len = %d, want 2", len(resp.GetParts()))
	}
}

func TestListPartsMapsUnexpectedErrorToInternal(t *testing.T) {
	t.Parallel()

	api := NewAPI(partServiceStub{err: errors.New("database is unavailable")})

	_, err := api.ListParts(context.Background(), &inventoryv1.ListPartsRequest{})
	if status.Code(err) != codes.Internal {
		t.Fatalf("code = %v, want Internal", status.Code(err))
	}
}
