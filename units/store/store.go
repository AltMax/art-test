package store

import (
	"context"
	"sync"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/units"
)

type Store struct {
	units.Units
	sync.RWMutex
	store map[string]*models.Unit
}

func NewStore(dao units.Units) *Store {
	return &Store{
		Units: dao,
		store: make(map[string]*models.Unit),
	}
}

func (s *Store) Create(ctx context.Context, unit *models.Unit) error {
	err := s.Units.Create(ctx, unit)
	if err != nil {
		return err
	}

	s.saveUnits(unit)

	return nil
}

func (s *Store) Update(ctx context.Context, id string, data []byte) (*models.Unit, error) {
	updatedUnit, err := s.Units.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}

	s.saveUnits(updatedUnit)

	return updatedUnit, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	err := s.Units.Delete(ctx, id)
	if err != nil {
		return nil
	}

	s.removeUnit(id)

	return nil
}

func (s *Store) FindByID(ctx context.Context, id string) (u *models.Unit, err error) {
	u = s.getByID(id)
	if u != nil {
		return u, nil
	}

	u, err = s.Units.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.saveUnits(u)

	return u, nil
}

func (s *Store) FindByIDs(ctx context.Context, ids []string) (models.Units, error) {
	return units.FindByIDs(ctx, ids, s.getByIDs, s.saveUnits, s.Units)
}

func (s *Store) saveUnits(units ...*models.Unit) {
	s.Lock()
	defer s.Unlock()
	for _, unit := range units {
		s.store[unit.ID] = unit
	}
}

func (s *Store) removeUnit(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.store, id)
}

func (s *Store) getByID(id string) *models.Unit {
	s.RLock()
	defer s.RUnlock()
	if t, ok := s.store[id]; ok {
		return t
	}
	return nil
}

func (s *Store) getByIDs(ids []string) []*models.Unit {
	s.RLock()
	defer s.RUnlock()

	units := make([]*models.Unit, 0, len(ids))

	for _, id := range ids {
		if u, ok := s.store[id]; ok {
			units = append(units, u)
		}
	}

	return units
}

func (s *Store) FetchAll(ctx context.Context) (models.Units, error) {
	units, err := s.Units.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	s.saveUnits(units...)
	return units, nil
}
