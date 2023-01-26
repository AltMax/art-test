package main

import (
	"database/sql"

	"github.com/AltMax/art-test/config"
	_ "github.com/AltMax/art-test/migrations/common/migrations"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("read config")
	}
	connStr := conf.Postgresql.ConnString()
	migrationsConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("open migrations db connection")
	}
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("set dialect")
	}
	if err := goose.Up(migrationsConn, "."); err != nil {
		log.Fatal().Err(err).Msgf("migrate up. %s", connStr)
	}
	_ = migrationsConn.Close()
}
