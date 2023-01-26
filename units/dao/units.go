package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/postgresql"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

var (
	ErrNotFound = errors.New("not found")

	selectUnitBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("id", "data", "created_at").From("units")
)

type Units struct {
	db postgresql.DB
}

func NewUnits(db postgresql.DB) *Units {
	return &Units{db: db}
}

func (u *Units) Create(ctx context.Context, unit *models.Unit) error {
	const op = "units.Units.Create"

	_, err := u.db.ExecCtx(
		ctx,
		`insert into units(
			id, data, created_at
		) 
		values(
			$1, $2, $3
		) 
		on conflict(id) do nothing`,
		unit.ID, unit.Data, unit.CreatedAt)
	if err != nil {
		return wrap(op, err)
	}

	return nil
}

func (u *Units) Update(ctx context.Context, id string, data []byte) (*models.Unit, error) {
	const op = "units.Units.Update"

	unit := &models.Unit{
		ID:   id,
		Data: data,
	}

	err := u.db.QueryRowCtx(ctx, `update units set data = $2 where id = $1 returning created_at`, unit.ID, unit.Data).Scan(&unit.CreatedAt)
	switch err {
	case pgx.ErrNoRows:
		return nil, ErrNotFound
	case nil:
		return unit, nil
	default:
		return nil, wrap(op, err)
	}
}

func (u *Units) Delete(ctx context.Context, id string) error {
	const op = "units.Units.Delete"

	tag, err := u.db.ExecCtx(ctx, `delete from units where id = $1`, id)
	if err != nil {
		return wrap(op, err)
	}

	if tag.RowsAffected() == 0 {
		return wrap(op, ErrNotFound)
	}

	return nil
}

func (u *Units) FindByID(ctx context.Context, id string) (*models.Unit, error) {
	const op = "units.Units.FindByID"
	units, err := u.queryUnits(ctx, selectUnitBuilder.Where(sq.Eq{"id": id}))
	if err != nil {
		return nil, wrap(op, err)
	}
	if len(units) == 0 {
		return nil, ErrNotFound
	}
	return units[0], nil
}

func (u *Units) FindByIDs(ctx context.Context, ids []string) (models.Units, error) {
	const op = "units.Units.FindByIDs"
	units, err := u.queryUnits(ctx, selectUnitBuilder.Where(sq.Eq{"id": ids}))
	if err != nil {
		return nil, wrap(op, err)
	}
	return units, nil
}

func (u *Units) FetchAll(ctx context.Context) (models.Units, error) {
	const op = "units.Units.FetchAll"
	units, err := u.queryUnits(ctx, selectUnitBuilder)
	if err != nil {
		return nil, wrap(op, err)
	}
	return units, nil
}

func (u *Units) queryUnits(ctx context.Context, builder sq.SelectBuilder) (models.Units, error) {
	units := make(models.Units, 0)
	rows, err := u.db.QueryxCtx(ctx, builder)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := &models.Unit{}
		err := scanUnit(rows, u)
		if err != nil {
			return nil, err
		}
		units = append(units, u)
	}

	return units, nil
}

func wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s, %w", op, err)
}

func scanUnit(row pgx.Row, unit *models.Unit) error {
	return row.Scan(
		&unit.ID,
		&unit.Data,
		&unit.CreatedAt,
	)
}
