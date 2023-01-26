package server

import (
	"context"

	"github.com/AltMax/art-test/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UnitService) GetUnits(ctx context.Context, req *services.GetUnitsRequest) (*services.GetUnitsResponse, error) {
	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one id is required")
	}

	units, err := h.units.FindByIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	return &services.GetUnitsResponse{
		Units: units.Proto(),
	}, nil
}
