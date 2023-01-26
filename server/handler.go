package server

import (
	"time"

	"github.com/AltMax/art-test/units"
)

type UnitService struct {
	units             units.Units
	fetchUnitsTimeout time.Duration
}

func NewUnitService(units units.Units, d time.Duration) *UnitService {
	return &UnitService{
		units:             units,
		fetchUnitsTimeout: d,
	}
}
