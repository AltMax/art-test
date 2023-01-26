package store

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
	testStore := NewStore(unitsMock)

	ctx := context.Background()

	unit := randomUnit()

	unitsMock.On("Create", mock.Anything, unit).Return(nil)
	err := testStore.Create(ctx, unit)
	require.NoError(t, err)

	storedUnit := testStore.getByID(unit.ID)
	require.Equal(t, unit, storedUnit)
}

func Test_Update(t *testing.T) {
	unitsMock := &mocks.Units{}
	testStore := NewStore(unitsMock)

	ctx := context.Background()

	unit := randomUnit()

	testStore.saveUnits(unit)
	storedUnit := testStore.getByID(unit.ID)
	require.Equal(t, unit, storedUnit)

	unit.Data = []byte("updated data")
	unitsMock.On("Update", mock.Anything, unit.ID, unit.Data).Return(unit, nil)
	updatedUnit, err := testStore.Update(ctx, unit.ID, unit.Data)
	require.NoError(t, err)
	require.Equal(t, unit, updatedUnit)

	storedUnit = testStore.getByID(unit.ID)
	require.Equal(t, unit, storedUnit)
}

func Test_Delete(t *testing.T) {
	unitsMock := &mocks.Units{}
	testStore := NewStore(unitsMock)

	ctx := context.Background()

	unit := randomUnit()

	unitsMock.On("Create", mock.Anything, unit).Return(nil)
	err := testStore.Create(ctx, unit)
	require.NoError(t, err)

	unitsMock.On("Delete", mock.Anything, unit.ID).Return(nil)
	err = testStore.Delete(ctx, unit.ID)
	require.NoError(t, err)

	storedUnit := testStore.getByID(unit.ID)
	require.Nil(t, storedUnit)
}

func Test_FindByID(t *testing.T) {
	unitsMock := &mocks.Units{}
	testStore := NewStore(unitsMock)

	ctx := context.Background()

	unit := randomUnit()

	//first call gets from db
	unitsMock.On("FindByID", mock.Anything, unit.ID).Return(unit, nil).Once()
	actualUnit, err := testStore.FindByID(ctx, unit.ID)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)

	//second call finds in store
	actualUnit, err = testStore.FindByID(ctx, unit.ID)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)
}

func Test_FindByIDs(t *testing.T) {
	unitsMock := &mocks.Units{}
	testStore := NewStore(unitsMock)

	ctx := context.Background()

	units := models.Units{
		randomUnit(),
		randomUnit(),
		randomUnit(),
	}
	testStore.saveUnits(units...)

	dbUnit := randomUnit()

	unitsMock.On("FindByIDs", mock.Anything, []string{dbUnit.ID}).Return(models.Units{dbUnit}, nil)
	actualUnits, err := testStore.FindByIDs(ctx, []string{units[0].ID, units[1].ID, units[2].ID, dbUnit.ID})
	require.NoError(t, err)
	require.Equal(t, append(units, dbUnit), actualUnits)
}

func Test_FetchAll(t *testing.T) {
	unitsMock := &mocks.Units{}
	testStore := NewStore(unitsMock)

	ctx := context.Background()

	units := models.Units{
		randomUnit(),
		randomUnit(),
		randomUnit(),
	}

	unitsMock.On("FetchAll", mock.Anything).Return(units, nil)
	actualUnits, err := testStore.FetchAll(ctx)
	require.NoError(t, err)
	require.Equal(t, units, actualUnits)

	storedUnits := testStore.getByIDs([]string{units[0].ID, units[1].ID, units[2].ID})
	require.Equal(t, units, models.Units(storedUnits))
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
