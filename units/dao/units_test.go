package dao

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/AltMax/art-test/config"
	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

func Test_Create_Positive(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	unit := randomUnit()

	actualUnit := &models.Unit{}
	row := postgresDB.QueryRowCtx(ctx, `select id, data, created_at from units where id = $1`, unit.ID)
	err = scanUnit(row, actualUnit)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	err = testUnits.Create(ctx, unit)
	require.NoError(t, err)

	actualUnit = &models.Unit{}
	row = postgresDB.QueryRowCtx(ctx, `select id, data, created_at from units where id = $1`, unit.ID)
	err = scanUnit(row, actualUnit)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)
}

func Test_Create_Duplicate(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	unit := randomUnit()

	err = testUnits.Create(ctx, unit)
	require.NoError(t, err)

	actualUnit := &models.Unit{}
	row := postgresDB.QueryRowCtx(ctx, `select id, data, created_at from units where id = $1`, unit.ID)
	err = scanUnit(row, actualUnit)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)

	newUnit := randomUnit()
	newUnit.ID = unit.ID

	err = testUnits.Create(ctx, unit)
	require.NoError(t, err)

	actualUnit = &models.Unit{}
	row = postgresDB.QueryRowCtx(ctx, `select id, data, created_at from units where id = $1`, unit.ID)
	err = scanUnit(row, actualUnit)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)
}

func Test_Update_Positive(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	unit := randomUnit()

	err = testUnits.Create(ctx, unit)
	require.NoError(t, err)

	unit.Data = []byte("updated data")

	updatedUnit, err := testUnits.Update(ctx, unit.ID, unit.Data)
	require.NoError(t, err)
	require.Equal(t, unit, updatedUnit)

	actualUnit, err := testUnits.FindByID(ctx, unit.ID)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)
}

func Test_Update_NotFound(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	unit := randomUnit()

	err = testUnits.Create(ctx, unit)
	require.NoError(t, err)

	newUnit := randomUnit()

	updatedUnit, err := testUnits.Update(ctx, newUnit.ID, newUnit.Data)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, updatedUnit)

	actualUnit, err := testUnits.FindByID(ctx, unit.ID)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)
}

func Test_Delete_Positive(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	unit := randomUnit()

	err = testUnits.Create(ctx, unit)
	require.NoError(t, err)

	actualUnit, err := testUnits.FindByID(ctx, unit.ID)
	require.NoError(t, err)
	require.Equal(t, unit, actualUnit)

	err = testUnits.Delete(ctx, unit.ID)
	require.NoError(t, err)

	deletedUnit, err := testUnits.FindByID(ctx, unit.ID)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, deletedUnit)
}

func Test_Delete_NotFound(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	err = testUnits.Delete(ctx, uuid.New().String())
	require.ErrorIs(t, err, ErrNotFound)
}

func Test_FindByIDs_Positive(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	unit := randomUnit()
	unit2 := randomUnit()
	unit3 := randomUnit()

	err = testUnits.Create(ctx, unit)
	require.NoError(t, err)
	err = testUnits.Create(ctx, unit2)
	require.NoError(t, err)
	err = testUnits.Create(ctx, unit3)
	require.NoError(t, err)

	expectedUnits := models.Units{unit, unit2, unit3}

	actualUnits, err := testUnits.FindByIDs(ctx, []string{unit.ID, unit2.ID, unit3.ID})
	require.NoError(t, err)
	require.Equal(t, expectedUnits, actualUnits)
}

func Test_FetchAll(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	require.NoError(t, err)
	defer postgresDB.Close()

	testUnits := NewUnits(postgresDB)
	ctx := context.Background()

	_, err = postgresDB.ExecCtx(ctx, `delete from units`)
	require.NoError(t, err)

	units := models.Units{}
	for i := 0; i < 10; i++ {
		unit := randomUnit()
		err := testUnits.Create(ctx, unit)
		require.NoError(t, err)
		units = append(units, unit)
	}

	actualUnits, err := testUnits.FetchAll(ctx)
	require.NoError(t, err)
	require.Equal(t, units, actualUnits)
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
