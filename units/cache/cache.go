package cache

import (
	"context"

	"github.com/AltMax/art-test/models"
	"github.com/AltMax/art-test/units"
	lru "github.com/hashicorp/golang-lru/v2"
)

type Cache struct {
	units.Units
	cache *lru.Cache[string, *models.Unit]
}

func NewCache(store units.Units, size int) (*Cache, error) {
	l, err := lru.New[string, *models.Unit](size)
	if err != nil {
		return nil, err
	}
	return &Cache{Units: store, cache: l}, nil
}

func (c *Cache) Create(ctx context.Context, unit *models.Unit) error {
	err := c.Units.Create(ctx, unit)
	if err != nil {
		return err
	}

	c.add(unit)

	return nil
}

func (c *Cache) Update(ctx context.Context, id string, data []byte) (*models.Unit, error) {
	updatedUnit, err := c.Units.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}

	c.add(updatedUnit)

	return updatedUnit, nil
}

func (c *Cache) Delete(ctx context.Context, id string) error {
	err := c.Units.Delete(ctx, id)
	if err != nil {
		return err
	}

	c.remove(id)

	return nil
}

func (c *Cache) FindByID(ctx context.Context, id string) (*models.Unit, error) {
	unit := c.getByID(id)
	if unit != nil {
		return unit, nil
	}

	unit, err := c.Units.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.add(unit)

	return unit, nil
}

func (c *Cache) FindByIDs(ctx context.Context, ids []string) (models.Units, error) {
	return units.FindByIDs(ctx, ids, c.getByIDs, c.add, c.Units)
}

func (c *Cache) add(units ...*models.Unit) {
	for _, unit := range units {
		c.cache.Add(unit.ID, unit)
	}
}

func (c *Cache) remove(id string) {
	c.cache.Remove(id)
}

func (c *Cache) getByID(id string) *models.Unit {
	unit, ok := c.cache.Get(id)
	if !ok {
		return nil
	}

	return unit
}

func (c *Cache) getByIDs(ids []string) []*models.Unit {
	units := make([]*models.Unit, 0, len(ids))
	for _, id := range ids {
		if unit, ok := c.cache.Get(id); ok {
			units = append(units, unit)
		}
	}
	return units
}

func (c *Cache) FetchAll(ctx context.Context) (models.Units, error) {
	units, err := c.Units.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, unit := range units {
		if c.cache.Contains(unit.ID) {
			c.cache.Add(unit.ID, unit)
		}
	}
	return units, nil
}
