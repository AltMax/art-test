package server

import (
	"context"
	"testing"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_GetUnits_Negative_MissingParameter(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	resp, err := handler.unitServiceClient.GetUnits(ctx, &services.GetUnitsRequest{Ids: []string{}})
	require.NotNil(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, status.Code())
	require.Nil(t, resp)
}

func Test_GetUnits_Positive_NotFound(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	handler.unitsMock.On("FindByIDs", mock.Anything, []string{"notExistID"}).Return(models.Units{}, nil)
	resp, err := handler.unitServiceClient.GetUnits(ctx, &services.GetUnitsRequest{Ids: []string{"notExistID"}})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Units, 0)
}

func Test_GetUnits_Positive(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	units := models.Units{
		randomUnit(),
		randomUnit(),
		randomUnit(),
	}

	unitIDs := []string{units[0].ID, units[1].ID, units[2].ID}

	handler.unitsMock.On("FindByIDs", mock.Anything, unitIDs).Return(units, nil)

	resp, err := handler.unitServiceClient.GetUnits(ctx, &services.GetUnitsRequest{Ids: unitIDs})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, units.Proto(), resp.Units)
}
