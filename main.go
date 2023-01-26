package main

import (
	"context"
	"net"
	"time"

	"github.com/AltMax/art-test/config"
	"github.com/AltMax/art-test/postgresql"
	"github.com/AltMax/art-test/server"
	"github.com/AltMax/art-test/services"
	"github.com/AltMax/art-test/units/cache"
	"github.com/AltMax/art-test/units/dao"
	"github.com/AltMax/art-test/units/store"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config creation")
	}

	postgresDB, err := postgresql.NewConnectionPool(conf.Postgresql)
	if err != nil {
		log.Fatal().Err(err).Msg("create postgres session")
	}
	unitsDao := dao.NewUnits(postgresDB)
	store := store.NewStore(unitsDao)
	cache, err := cache.NewCache(store, conf.LRUCacheSize)
	if err != nil {
		log.Fatal().Err(err).Int("lru-cache-size", conf.LRUCacheSize).Msg("create lru cache with size")
	}
	ctx := context.Background()

	fetchUnitsTimeout := time.Duration(conf.FetchUnitsTimeout) * time.Second
	handler := server.NewUnitService(cache, fetchUnitsTimeout)

	//первая синхронизация при запуске
	_, err = cache.FetchAll(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("first units fetch")
	}

	//синхронизация каждые [conf.FetchUnitsTimeout] секунд
	go handler.FetchUnitsSometimes(ctx)

	unitServer := server.New(&conf)
	services.RegisterUnitServiceServer(unitServer, handler)
	lis, err := net.Listen("tcp", conf.ServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen unit service")
	}
	log.Info().Msg("unit server started")
	if err := unitServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("listen unit server")
	}
}
