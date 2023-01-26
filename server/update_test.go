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

func Test_Update_Negative_MissingParameter(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	resp, err := handler.unitServiceClient.Update(ctx, &services.UpdateUnitRequest{Id: ""})
	require.NotNil(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Nil(t, resp)

	resp, err = handler.unitServiceClient.Update(ctx, &services.UpdateUnitRequest{Id: "randomID", Data: []byte{}})
	require.NotNil(t, err)
	st, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
	require.Nil(t, resp)
}

func Test_Update_Negative_NotFound(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	unitID := "notExistID"
	unitData := []byte("new data")
	handler.unitsMock.On("Update", mock.Anything, unitID, unitData).Return(nil, dao.ErrNotFound)
	resp, err := handler.unitServiceClient.Update(ctx, &services.UpdateUnitRequest{Id: unitID, Data: unitData})
	require.NotNil(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, status.Code())
	require.Nil(t, resp)
}

func Test_Update_Positive(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()
	unit := randomUnit()

	handler.unitsMock.On("Update", mock.Anything, unit.ID, unit.Data).Return(unit, nil)
	resp, err := handler.unitServiceClient.Update(ctx, &services.UpdateUnitRequest{Id: unit.ID, Data: unit.Data})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, unit.Proto(), resp)
}
