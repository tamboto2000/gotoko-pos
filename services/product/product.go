package product

import (
	"context"

	"github.com/tamboto2000/gotoko-pos/domain/products"
	"go.uber.org/zap"
)

type ProductServiceConfig func(*ProductService) error

type ProductService struct {
	prodRepo *products.ProductsRepository
	log      *zap.Logger
}

func NewProductService(cfgs ...ProductServiceConfig) (*ProductService, error) {
	prodSvc := new(ProductService)
	for _, cfg := range cfgs {
		if err := cfg(prodSvc); err != nil {
			return nil, err
		}
	}

	return prodSvc, nil
}

func WithProductsRepository(repo *products.ProductsRepository) ProductServiceConfig {
	return func(ps *ProductService) error {
		ps.prodRepo = repo
		return nil
	}
}

func WithLogger(log *zap.Logger) ProductServiceConfig {
	return func(ps *ProductService) error {
		ps.log = log
		return nil
	}
}

func (prodSvc *ProductService) CreateProduct(ctx context.Context, prod *products.Product) error {
	if err := prod.Validate(); err != nil {
		return err
	}

	return prodSvc.prodRepo.CreateProduct(ctx, prod)
}

func (prodSvc *ProductService) GetProductDetail(ctx context.Context, id int) (*products.Product, error) {
	return prodSvc.prodRepo.GetProductDetail(ctx, id)
}

func (prodSvc *ProductService) GetProductList(ctx context.Context, limit, skip, categoryId int, qs string) (*products.ProductList, error) {
	return prodSvc.prodRepo.GetProductList(ctx, limit, skip, categoryId, qs)
}

func (prodSvc *ProductService) UpdateProduct(ctx context.Context, prod *products.Product) error {
	if err := prod.ValidateForUpdate(); err != nil {
		return err
	}

	return prodSvc.prodRepo.UpdateProduct(ctx, prod)
}

func (prodSvc *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return prodSvc.prodRepo.DeleteProduct(ctx, id)
}

func (prodSvc *ProductService) CreateCategory(ctx context.Context, cat *products.Category) error {
	if err := cat.Validate(); err != nil {
		return err
	}

	return prodSvc.prodRepo.CreateCategory(ctx, cat)
}

func (prodSvc *ProductService) GetCategoryDetail(ctx context.Context, id int) (*products.Category, error) {
	return prodSvc.prodRepo.GetCategoryDetail(ctx, id)
}

func (prodSvc *ProductService) GetCategoryList(ctx context.Context, limit, skip int) (*products.CategoryList, error) {
	return prodSvc.prodRepo.GetCategoryList(ctx, limit, skip)
}

func (prodSvc *ProductService) UpdateCategory(ctx context.Context, cat *products.Category) error {
	if err := cat.ValidateForUpdate(); err != nil {
		return err
	}

	return prodSvc.prodRepo.UpdateCategory(ctx, cat)
}

func (prodSvc *ProductService) DeleteCategory(ctx context.Context, id int) error {
	return prodSvc.prodRepo.DeleteCategory(ctx, id)
}
