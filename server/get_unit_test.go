package server

import (
	"context"
	"testing"

	"github.com/AltMax/art-test/services"
	"github.com/AltMax/art-test/units/dao"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_GetUnit_Negative_MissingParameter(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	resp, err := handler.unitServiceClient.GetUnit(ctx, &services.GetUnitRequest{Id: ""})
	require.NotNil(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, status.Code())
	require.Nil(t, resp)
}

func Test_GetUnit_Negative_NotFound(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	handler.unitsMock.On("FindByID", mock.Anything, "notExistID").Return(nil, dao.ErrNotFound)
	resp, err := handler.unitServiceClient.GetUnit(ctx, &services.GetUnitRequest{Id: "notExistID"})
	require.NotNil(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, status.Code())
	require.Nil(t, resp)
}

func Test_GetUnit_Positive(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	unit := randomUnit()

	handler.unitsMock.On("FindByID", mock.Anything, unit.ID).Return(unit, nil)

	resp, err := handler.unitServiceClient.GetUnit(ctx, &services.GetUnitRequest{Id: unit.ID})
	require.NoError(t, err)
	require.Equal(t, unit.Proto(), resp)
}
