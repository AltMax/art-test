package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(UpInit, DownInit)
}

var upInit = `
create table if not exists units (
	id text not null,
	data bytea not null,
	created_at timestamp not null,
	primary key(id)
)
`

func UpInit(tx *sql.Tx) error {
	// This code is executed when the migrations is applied.
	_, err := tx.Exec(upInit)
	return err
}

func DownInit(tx *sql.Tx) error {
	// This code is executed when the migrations is rolled back.
	return nil
}
