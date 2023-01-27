package products

import (
	"context"

	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type CacheProductRepository struct {
	cc  *cache.Cache
	log *zap.Logger
}

func NewCacheProductsRepository(log *zap.Logger) *CacheProductRepository {
	return &CacheProductRepository{
		cc:  cache.New(cache.NoExpiration, cache.NoExpiration),
		log: log,
	}
}

func (cProdRepo *CacheProductRepository) GetProductList(ctx context.Context, limit, skip, categoryId int, qs string) (*ProductList, error) {
	

	return nil, nil
}
