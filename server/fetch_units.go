package server

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (h *UnitService) FetchUnitsSometimes(ctx context.Context) {
	fetchTiker := time.NewTicker(h.fetchUnitsTimeout)
	defer fetchTiker.Stop()

	//чтобы сихронизация стора с базой не отъедала время тикера
	resetTiker := func() {
		fetchTiker.Stop()
		fetchTiker = time.NewTicker(h.fetchUnitsTimeout)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-fetchTiker.C:
			_, err := h.units.FetchAll(ctx)
			if err != nil {
				log.Error().Err(err).Msg("fetch units")
			}
			resetTiker()
		}
	}
}
