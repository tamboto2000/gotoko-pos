package cashiers

import (
	"context"
	"sync"

	"github.com/tamboto2000/gotoko-pos/models"
)

type CashiersMemRepository struct {
	smap *sync.Map
}

func NewMemory() *CashiersMemRepository {
	return &CashiersMemRepository{
		smap: new(sync.Map),
	}
}

func (cashierMRepo *CashiersMemRepository) CreateSession(ctx context.Context, c *models.CashierSession) error {
	val, ok := cashierMRepo.smap.Load(c.Id)
	if !ok {
		cashierMRepo.smap.Store(c.CashierId, []*models.CashierSession{c})
		return nil
	}

	sessions := val.([]*models.CashierSession)
	sessions = append(sessions, c)
	cashierMRepo.smap.Store(c.CashierId, sessions)
	return nil
}

func (cashierMRepo *CashiersMemRepository) DeleteSession(ctx context.Context, id int) error {
	cashierMRepo.smap.Delete(id)
	return nil
}

func (cashierMRepo *CashiersMemRepository) GetSession(ctx context.Context, cashierId int, issuedAt int64) (*models.CashierSession, error) {
	val, ok := cashierMRepo.smap.Load(cashierId)
	if !ok {
		return nil, ErrCashierSessionNotFound
	}

	sessions := val.([]*models.CashierSession)
	for _, s := range sessions {
		if s.CashierId == cashierId && s.IssuedAt == issuedAt {
			return s, nil
		}
	}

	return nil, ErrCashierSessionNotFound
}
