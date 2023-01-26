package units

import (
	"context"

	"github.com/AltMax/art-test/models"
)

type Units interface {
	Create(ctx context.Context, unit *models.Unit) error
	Update(ctx context.Context, id string, data []byte) (*models.Unit, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*models.Unit, error)
	FindByIDs(ctx context.Context, ids []string) (models.Units, error)
	FetchAll(ctx context.Context) (models.Units, error)
}

func deduplicateIDs(ids []string) []string {
	idSet := make(map[string]struct{}, len(ids))
	uniqueIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		if _, isDuplicate := idSet[id]; !isDuplicate {
			idSet[id] = struct{}{}
			uniqueIDs = append(uniqueIDs, id)
		}
	}
	return uniqueIDs
}

func FindByIDs(
	ctx context.Context,
	ids []string,
	getByIDs func(ids []string) []*models.Unit,
	saveUnits func(units ...*models.Unit),
	nextLayer Units,
) (models.Units, error) {
	uniqueIDs := deduplicateIDs(ids)

	units := getByIDs(uniqueIDs)
	if len(units) == len(ids) {
		return units, nil
	}

	fetchedUnits := make(map[string]struct{}, len(units))
	for _, unit := range units {
		fetchedUnits[unit.ID] = struct{}{}
	}

	missedIDs := uniqueIDs[:0]
	for _, id := range uniqueIDs {
		if _, isFetched := fetchedUnits[id]; !isFetched {
			missedIDs = append(missedIDs, id)
		}
	}

	if len(missedIDs) > 0 {
		dbUnits, err := nextLayer.FindByIDs(ctx, missedIDs)
		if err != nil {
			return nil, err
		}
		saveUnits(dbUnits...)
		units = append(units, dbUnits...)
	}

	return units, nil
}
