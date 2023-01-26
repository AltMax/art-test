package server

import (
	"context"
	"testing"
	"time"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Create_Negative_MissingParameter(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	resp, err := handler.unitServiceClient.Create(ctx, &services.CreateUnitRequest{Data: []byte{}})
	require.NotNil(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, status.Code())
	require.Nil(t, resp)
}

func Test_Create_Positive(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	unitData := []byte("some data")
	unit := randomUnit()
	unit.Data = unitData

	handler.unitsMock.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(unit *models.Unit) bool {
			return unit.ID != "" &&
				string(unit.Data) == string(unitData) &&
				time.Now().After(unit.CreatedAt) &&
				time.Now().Before(unit.CreatedAt.Add(1*time.Second))
		}),
	).Return(nil)

	resp, err := handler.unitServiceClient.Create(ctx, &services.CreateUnitRequest{Data: unitData})
	require.NoError(t, err)
	require.Equal(t, unit.Data, resp.Data)
}
