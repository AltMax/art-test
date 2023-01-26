package server

import (
	"context"
	"errors"

	"github.com/AltMax/art-test/services"
	"github.com/AltMax/art-test/units/dao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UnitService) Update(ctx context.Context, req *services.UpdateUnitRequest) (*services.Unit, error) {
	if len(req.Data) == 0 {
		return nil, status.Error(codes.InvalidArgument, "data is required")
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	unit, err := h.units.Update(ctx, req.Id, req.Data)
	if errors.Is(err, dao.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "unit not found")
	}
	if err != nil {
		return nil, err
	}

	return unit.Proto(), nil
}
