package server

import (
	"context"
	"errors"

	"github.com/AltMax/art-test/services"
	"github.com/AltMax/art-test/units/dao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UnitService) Delete(ctx context.Context, req *services.DeleteUnitRequest) (*services.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := h.units.Delete(ctx, req.Id)
	if errors.Is(err, dao.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "unit not found")
	}
	if err != nil {
		return nil, err
	}

	return &services.Empty{}, nil
}
