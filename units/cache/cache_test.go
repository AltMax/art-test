package cache

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/units/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Create(t *testing.T) {
	unitsMock := &mocks.Units{}
	testCache, err := NewCache(unitsMock, 10)
	require.NoError(t, err)

	ctx := context.Background()

	unit := randomUnit()

	unitsMock.On("Create", mock.Anything, unit).Return(nil)
	err = testCache.Create(ctx, unit)
	require.NoError(t, err)

	chachedUnit := testCache.getByID(unit.ID)
	require.Equal(t, unit, chachedUnit)
}

func Test_Update(t *testing.T) {
	unitsMock := &mocks.Units{}
	testCache, err := NewCache(unitsMock, 10)
	require.NoError(t, err)

	ctx := context.Background()

	unit := randomUnit()

	testCache.add(unit)
	chachedUnit := testCache.getByID(unit.ID)
	require.Equal(t, unit, chachedUnit)

	unit.Data = []byte("updated data")
	unitsMock.On("Update", mock.Anything, unit.ID, unit.Data).Return(unit, nil)
	updatedUnit, err := testCache.Update(ctx, unit.ID, unit.Data)
	require.NoError(t, err)
	require.Equal(t, unit, updatedUnit)

	chachedUnit = testCache.getByID(unit.ID)
	require.Equal(t, unit, chachedUnit)
}

func Test_Delete(t *testing.T) {
	unitsMock := &mocks.Units{}
	testCache, err := NewCache(unitsMock, 10)
	require.NoError(t, err)

	ctx := context.Background()

	unit := randomUnit()

	unitsMock.On("Create", mock.Anything, unit).Return(nil)
	err = testCache.Create(ctx, unit)
	require.NoError(t, err)

	unitsMock.On("Delete", mock.Anything, unit.ID).Return(nil)
	err = testCache.Delete(ctx, unit.ID)
	require.NoError(t, err)

	chachedUnit := testCache.getByID(unit.ID)
	require.Nil(t, chachedUnit)
}

func Test_FindByID(t *testing.T) {
	unitsMock := &mocks.Units{}
	testCache, err := NewCache(unitsMock, 10)
	require.NoError(t, err)

	ctx := context.Background()

	unit := randomUnit()

	//first call gets from db
	unitsMock.On("FindByID", mock.Anything, unit.ID).Return(unit, nil).Once()
	actualUnit, err := testCache.FindByID(ctx, unit.ID)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)

	//second call finds in store
	actualUnit, err = testCache.FindByID(ctx, unit.ID)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)
}

func Test_FindByIDs(t *testing.T) {
	unitsMock := &mocks.Units{}
	testCache, err := NewCache(unitsMock, 10)
	require.NoError(t, err)

	ctx := context.Background()

	units := models.Units{
		randomUnit(),
		randomUnit(),
		randomUnit(),
	}
	testCache.add(units...)

	dbUnit := randomUnit()

	unitsMock.On("FindByIDs", mock.Anything, []string{dbUnit.ID}).Return(models.Units{dbUnit}, nil)
	actualUnits, err := testCache.FindByIDs(ctx, []string{units[0].ID, units[1].ID, units[2].ID, dbUnit.ID})
	require.NoError(t, err)
	require.Equal(t, append(units, dbUnit), actualUnits)
}

func Test_FetchAll(t *testing.T) {
	unitsMock := &mocks.Units{}
	testCache, err := NewCache(unitsMock, 10)
	require.NoError(t, err)

	ctx := context.Background()

	units := models.Units{
		randomUnit(),
		randomUnit(),
		randomUnit(),
	}
	testCache.add(units...)

	unitsMock.On("FetchAll", mock.Anything).Return(units, nil)
	actualUnits, err := testCache.FetchAll(ctx)
	require.NoError(t, err)
	require.Equal(t, units, actualUnits)

	chachedUnits := testCache.getByIDs([]string{units[0].ID, units[1].ID, units[2].ID})
	require.Equal(t, units, models.Units(chachedUnits))
}

func randomUnit() *models.Unit {
	buf := make([]byte, 50)
	rand.Read(buf)
	return &models.Unit{
		ID:        uuid.New().String(),
		Data:      buf,
		CreatedAt: time.Now().UTC(),
	}
}
