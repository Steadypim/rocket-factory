package part

import (
	"context"
	"errors"
	"testing"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
)

type repositoryStub struct {
	part  *model.Part
	parts []model.Part
	err   error
}

func (r repositoryStub) Get(_ context.Context, _ string) (*model.Part, error) {
	return r.part, r.err
}

func (r repositoryStub) List(_ context.Context, _ model.PartsFilter) ([]model.Part, error) {
	return r.parts, r.err
}

func TestGetReturnsRepositoryPart(t *testing.T) {
	t.Parallel()

	service := NewService(repositoryStub{part: &model.Part{UUID: "part-1"}})

	part, err := service.Get(context.Background(), "part-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if part.UUID != "part-1" {
		t.Fatalf("UUID = %q, want part-1", part.UUID)
	}
}

func TestGetReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	service := NewService(repositoryStub{err: model.ErrPartNotFound})

	_, err := service.Get(context.Background(), "missing")
	if !errors.Is(err, model.ErrPartNotFound) {
		t.Fatalf("Get error = %v, want ErrPartNotFound", err)
	}
}

func TestListReturnsRepositoryParts(t *testing.T) {
	t.Parallel()

	service := NewService(repositoryStub{parts: []model.Part{{UUID: "part-1"}, {UUID: "part-2"}}})

	parts, err := service.List(context.Background(), model.PartsFilter{})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(parts) != 2 {
		t.Fatalf("parts len = %d, want 2", len(parts))
	}
}
