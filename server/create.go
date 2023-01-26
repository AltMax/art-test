package server

import (
	"context"
	"time"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/services"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UnitService) Create(ctx context.Context, req *services.CreateUnitRequest) (*services.Unit, error) {
	if len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "data is required")
	}

	id := uuid.New().String()
	unit := &models.Unit{
		ID:        id,
		Data:      req.Data,
		CreatedAt: time.Now().UTC(),
	}

	err := h.units.Create(ctx, unit)
	if err != nil {
		return nil, err
	}

	return unit.Proto(), nil
}
