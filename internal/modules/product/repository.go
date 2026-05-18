package product

import (
	"context"
	"database/sql"
	"fmt"
	"mekoko/internal/domain"
	"strings"
)

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

type Repository struct {
	db DBTX
}

func NewRepository(db DBTX) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) WithTx(db *sql.Tx) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddProduct(ctx context.Context, newProduct NewProduct) (*domain.Product, error) {
	query := `
		INSERT INTO products (public_id, name, description, discount_percentage, base_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, public_id, name, description, discount_percentage, base_price
	`

	var product domain.Product

	if err := r.db.QueryRowContext(ctx, query, newProduct.PublicID, newProduct.Name, newProduct.Description, newProduct.DiscountPercentage, newProduct.BasePrice).Scan(&product.ID, &product.PublicID, &product.Name, &product.Description, &product.DiscountPercentage, &product.BasePrice); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *Repository) AddProductVariant(ctx context.Context, newVariant NewVariant) error {
	query := `
		INSERT INTO product_variants (public_id, product_id, color, size, stock_quantity, image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, newVariant.PublicID, newVariant.ProductID, newVariant.Color, newVariant.Size, newVariant.StockQuantity, newVariant.ImageURL)

	return err
}

func (r *Repository) GetProducts(ctx context.Context, limit int, offset int, filter string) ([]domain.Product, int64, error) {
	var pOrderBy string
	var productsOrderBy string
	switch strings.TrimSpace(strings.ToLower(filter)) {
	case FilterPriceASC:
		pOrderBy = "p.base_price ASC"
		productsOrderBy = "products.base_price ASC"
	case FilterPriceDESC:
		pOrderBy = "p.base_price DESC"
		productsOrderBy = "products.base_price DESC"
	case OldestFirst:
		pOrderBy = "p.created_at ASC"
		productsOrderBy = "products.created_at ASC"
	case NewestFirst:
		pOrderBy = "p.created_at DESC"
		productsOrderBy = "products.created_at DESC"
	default:
		pOrderBy = "p.id ASC"
		productsOrderBy = "products.id ASC"
	}

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.public_id, p.name, p.discount_percentage, p.base_price, p.description, product_variants.id, product_variants.public_id, product_variants.product_id, product_variants.color, product_variants.size, product_variants.image_url, product_variants.stock_quantity
		FROM (SELECT * FROM products ORDER BY %s LIMIT $1 OFFSET $2) AS p
		JOIN product_variants ON product_variants.product_id = p.id
		ORDER BY %s
	`, productsOrderBy, pOrderBy)

	countQuery := `
		SELECT COUNT(*) FROM products
	`

	var count int64

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	productMap := make(map[int64]*domain.Product)

	for rows.Next() {
		var (
			productID                 int64
			productPublicID           string
			productName               string
			productBasePrice          int64
			productDiscountPercentage int
			productDescription        string
			variantID                 int64
			variantPublicID           string
			variantProductID          int64
			variantColor              string
			variantSize               string
			variantImageURL           string
			variantStockQuantity      int64
		)

		if err := rows.Scan(&productID, &productPublicID, &productName, &productDiscountPercentage, &productBasePrice, &productDescription, &variantID, &variantPublicID, &variantProductID, &variantColor, &variantSize, &variantImageURL, &variantStockQuantity); err != nil {
			return nil, 0, err
		}

		if productMap[productID] == nil {
			productMap[productID] = &domain.Product{
				ID:                 productID,
				Name:               productName,
				Description:        productDescription,
				PublicID:           productPublicID,
				BasePrice:          productBasePrice,
				DiscountPercentage: productDiscountPercentage,
			}
		}

		productMap[productID].Variants = append(productMap[productID].Variants, domain.ProductVariant{
			ID:            variantID,
			Color:         variantColor,
			Size:          variantSize,
			StockQuantity: variantStockQuantity,
			PublicID:      variantPublicID,
			ProductID:     variantProductID,
			ImageURL:      variantImageURL,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	products := make([]domain.Product, 0, len(productMap))
	for _, p := range productMap {
		products = append(products, *p)
	}

	countRow := r.db.QueryRowContext(ctx, countQuery)
	if err := countRow.Scan(&count); err != nil {
		return nil, 0, err
	}

	return products, count, nil
}
