package server

import (
	"testing"
	"time"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/units/mocks"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func Test_FetchUnitsSometimes(t *testing.T) {
	unitsMock := &mocks.Units{}

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	handler := NewUnitService(unitsMock, 1*time.Second)

	unitsMock.On("FetchAll", mock.Anything).Return(models.Units{randomUnit()}, nil)

	go handler.FetchUnitsSometimes(ctx)

	<-ctx.Done()

	unitsMock.AssertNumberOfCalls(t, "FetchAll", 2)
}
