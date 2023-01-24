package products

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	commonerr "github.com/tamboto2000/gotoko-pos/common/errors"
	"github.com/tamboto2000/gotoko-pos/models"
	"go.uber.org/zap"
)

const (
	productsTable    = "products"
	productSkusTable = "product_skus"
	discountsTable   = "discounts"
	categoriesTable  = "categories"
)

type ProductsRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewProductsRepository(db *sql.DB, log *zap.Logger) *ProductsRepository {
	return &ProductsRepository{db: db, log: log}
}

func (prodRepo *ProductsRepository) CreateProduct(ctx context.Context, prod *Product) error {
	tx, err := prodRepo.db.BeginTx(ctx, nil)
	if err != nil {
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	q := fmt.Sprintf(`INSERT INTO %s (
		name, 
		stock,
		price,
		image_url,
		category_id		
	) VALUES (?,?,?,?,?)`, productsTable)

	res, err := tx.ExecContext(
		ctx,
		q,
		prod.Name,
		prod.Stock,
		prod.Price,
		prod.ImageUrl,
		prod.CategoryId,
	)

	if err != nil {
		tx.Rollback()
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	if prod.Discount != nil {
		disc := prod.Discount
		q := fmt.Sprintf(`INSERT INTO %s (
			product_id,
			min_qty,
			type,
			result,
			expired_at
		) VALUES (?,?,?,?,?)`, discountsTable)

		_, err := tx.ExecContext(
			ctx,
			q,
			lastId,
			disc.MinQty,
			disc.Type,
			disc.Result,
			disc.ExpiredAt,
		)

		if err != nil {
			tx.Rollback()
			prodRepo.log.Error(err.Error())
			return commonerr.ErrInternal
		}
	}

	q = fmt.Sprintf(`INSERT INTO %s (product_id, sku) VALUES(?,?)`, productSkusTable)
	sku := fmt.Sprintf("ID%03d", lastId)
	if _, err = tx.ExecContext(ctx, q, lastId, sku); err != nil {
		tx.Rollback()
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	prod.Id = int(lastId)
	prod.Sku = sku

	return nil
}

func (prodRepo *ProductsRepository) GetProductDetail(ctx context.Context, id int) (*Product, error) {
	q := fmt.Sprintf(
		`SELECT 
			p.id,
			p.name, 
			pk.sku,
			p.stock, 
			p.price,
			p.image_url,
			p.category_id,
			c.name AS category_name,
			d.min_qty AS discount_min_qty,
			d.type AS discount_type,
			d.result AS discount_result,
			d.expired_at AS discount_expired_at
		FROM %s p
		INNER JOIN %s pk ON pk.product_id = p.id
		LEFT JOIN %s c ON c.id = p.category_id
		LEFT JOIN %s d ON d.product_id = p.id
		WHERE p.id = ?`,
		productsTable,
		productSkusTable,
		categoriesTable,
		discountsTable,
	)

	row := prodRepo.db.QueryRowContext(ctx, q, id)
	prod := new(Product)
	prod.Category = new(Category)
	prod.Discount = new(models.Discount)
	if err := row.Scan(
		&prod.Id,
		&prod.Name,
		&prod.Sku,
		&prod.Stock,
		&prod.Price,
		&prod.ImageUrl,
		&prod.Category.Id,
		&prod.Category.Name,
		&prod.Discount.MinQty,
		&prod.Discount.Type,
		&prod.Discount.Result,
		&prod.Discount.ExpiredAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductNotFound
		}

		prodRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	return prod, nil
}

func (prodRepo *ProductsRepository) GetProductList(ctx context.Context, limit, skip, categoryId int, qs string) (*ProductList, error) {
	// getProductCount(in_category_id INT, in_qs TEXT)
	q := `CALL getProductCount(?,?)`
	row := prodRepo.db.QueryRowContext(ctx, q, categoryId, qs)
	var count int
	if err := row.Scan(&count); err != nil {
		prodRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	// getProductList(in_category_id INT, in_qs TEXT, in_limit INT, in_skip INT)
	q = `CALL getProductList(?, ?, ?, ?);`

	if limit <= 0 {
		limit = count
	}

	rows, err := prodRepo.db.QueryContext(ctx, q, categoryId, qs, limit, skip)
	if err != nil {
		prodRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	prodList := new(ProductList)
	for rows.Next() {
		var prod Product
		prod.Category = new(Category)
		prod.Discount = new(models.Discount)
		if err := rows.Scan(
			&prod.Id,
			&prod.Name,
			&prod.Sku,
			&prod.Stock,
			&prod.Price,
			&prod.ImageUrl,
			&prod.Category.Id,
			&prod.Category.Name,
			&prod.Discount.MinQty,
			&prod.Discount.Type,
			&prod.Discount.Result,
			&prod.Discount.ExpiredAt,
		); err != nil {
			if err == sql.ErrNoRows {
				break
			}

			prodRepo.log.Error(err.Error())
			return nil, commonerr.ErrInternal
		}

		prodList.Products = append(prodList.Products, prod)
	}

	prodList.Meta = Meta{
		Total: count,
		Limit: limit,
		Skip:  skip,
	}

	return prodList, nil
}

func (prodRepo *ProductsRepository) UpdateProduct(ctx context.Context, prod *Product) error {
	q := fmt.Sprintf(`SELECT id FROM %s WHERE id = ?`, productsTable)
	row := prodRepo.db.QueryRowContext(ctx, q, prod.Id)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrProductNotFound
		}

		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	/*
		updateProduct(
			in_id INT,
			in_category_id INT,
			in_name VARCHAR(100),
			in_image_url TEXT,
			in_price INT UNSIGNED,
			in_stock INT UNSIGNED
		)
	*/
	q = `CALL updateProduct(?,?,?,?,?,?)`

	if _, err := prodRepo.db.ExecContext(
		ctx,
		q,
		id,
		prod.CategoryId,
		prod.Name,
		prod.ImageUrl,
		prod.Price,
		prod.Stock,
	); err != nil {
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	return nil
}

func (prodRepo *ProductsRepository) DeleteProduct(ctx context.Context, id int) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, productsTable)
	res, err := prodRepo.db.ExecContext(ctx, q, id)
	if err != nil {
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (prodRepo *ProductsRepository) CreateCategory(ctx context.Context, cat *Category) error {
	q := fmt.Sprintf(`INSERT INTO %s (name) VALUES (?)`, categoriesTable)
	res, err := prodRepo.db.ExecContext(ctx, q, cat.Name)
	if err != nil {
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	q = fmt.Sprintf(`SELECT created_at, updated_at FROM %s WHERE id = ?`, categoriesTable)
	row := prodRepo.db.QueryRowContext(ctx, q, lastId)
	cat.CreatedAt = new(time.Time)
	cat.UpdatedAt = new(time.Time)
	if err := row.Scan(
		&cat.CreatedAt,
		&cat.UpdatedAt,
	); err != nil {
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	cat.Id = int(lastId)

	return nil
}

func (prodRepo *ProductsRepository) GetCategoryDetail(ctx context.Context, id int) (*Category, error) {
	q := fmt.Sprintf(`SELECT id, name FROM %s WHERE id = ?`, categoriesTable)
	row := prodRepo.db.QueryRowContext(ctx, q, id)
	cat := new(Category)
	if err := row.Scan(
		&cat.Id,
		&cat.Name,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCategoryNotFound
		}

		prodRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	return cat, nil
}

func (prodRepo *ProductsRepository) GetCategoryList(ctx context.Context, limit, skip int) (*CategoryList, error) {
	q := fmt.Sprintf(`SELECT COUNT(id) FROM %s`, categoriesTable)
	var count int
	row := prodRepo.db.QueryRowContext(ctx, q)
	if err := row.Scan(&count); err != nil {
		prodRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	if limit <= 0 {
		limit = count
	}

	q = fmt.Sprintf(`SELECT id, name FROM %s ORDER BY id ASC LIMIT ? OFFSET ?`, categoriesTable)
	rows, err := prodRepo.db.QueryContext(ctx, q, limit, skip)
	if err != nil {
		prodRepo.log.Error(err.Error())
		return nil, commonerr.ErrInternal
	}

	catList := new(CategoryList)
	for rows.Next() {
		var cat Category
		if err := rows.Scan(
			&cat.Id,
			&cat.Name,
		); err != nil {
			if err == sql.ErrNoRows {
				break
			}

			prodRepo.log.Error(err.Error())
			return nil, commonerr.ErrInternal
		}

		catList.Categories = append(catList.Categories, cat)
	}

	catList.Meta = Meta{
		Total: count,
		Limit: limit,
		Skip:  skip,
	}

	return catList, nil
}

func (prodRepo *ProductsRepository) UpdateCategory(ctx context.Context, cat *Category) error {
	q := fmt.Sprintf(`SELECT id FROM %s WHERE id = ?`, categoriesTable)
	row := prodRepo.db.QueryRowContext(ctx, q, cat.Id)
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return ErrCategoryNotFound
		}

		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	q = `CALL updateCategory(?,?)`
	_, err := prodRepo.db.ExecContext(ctx, q, id, cat.Name)
	if err != nil {
		prodRepo.log.Error(err.Error())
		return commonerr.ErrInternal
	}

	return nil
}

func (prodRepo *ProductsRepository) DeleteCategory(ctx context.Context, id int) error {
	q := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, categoriesTable)
	res, err := prodRepo.db.ExecContext(ctx, q, id)
	if err != nil {
		prodRepo.log.Error(err.Error())
	}

	if count, _ := res.RowsAffected(); count == 0 {
		return ErrCategoryNotFound
	}

	return nil
}
