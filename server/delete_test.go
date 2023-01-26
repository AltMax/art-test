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

func Test_Delete_Negative_MissingParameter(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	resp, err := handler.unitServiceClient.Delete(ctx, &services.DeleteUnitRequest{Id: ""})
	require.NotNil(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, status.Code())
	require.Nil(t, resp)
}

func Test_Delete_Negative_NotFound(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	handler.unitsMock.On("Delete", mock.Anything, "notExistID").Return(dao.ErrNotFound)
	resp, err := handler.unitServiceClient.Delete(ctx, &services.DeleteUnitRequest{Id: "notExistID"})
	require.NotNil(t, err)
	status, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, status.Code())
	require.Nil(t, resp)
}

func Test_Delete_Positive(t *testing.T) {
	handler := newTestHandler()
	ctx := context.Background()

	handler.unitsMock.On("Delete", mock.Anything, "existsID").Return(nil)
	resp, err := handler.unitServiceClient.Delete(ctx, &services.DeleteUnitRequest{Id: "existsID"})
	require.NoError(t, err)
	require.NotNil(t, resp)
}
