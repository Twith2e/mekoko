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

func (r *Repository) GetProducts(ctx context.Context, limit int, offset int, filter Filter) ([]domain.Product, int64, error) {
	var pOrderBy string
	var productsOrderBy string
	switch strings.TrimSpace(strings.ToLower(string(filter.Order))) {
	case strings.ToLower(string(FilterPriceASC)):
		pOrderBy = "p.base_price ASC"
		productsOrderBy = "products.base_price ASC"
	case strings.ToLower(string(FilterPriceDESC)):
		pOrderBy = "p.base_price DESC"
		productsOrderBy = "products.base_price DESC"
	case strings.ToLower(string(FilterOldestFirst)):
		pOrderBy = "p.created_at ASC"
		productsOrderBy = "products.created_at ASC"
	case strings.ToLower(string(FilterNewestFirst)):
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

	args := []any{limit, offset}
	conditions := []string{}
	filterConditions := []string{}

	if filter.Color != nil {
		args = append(args, filter.Color)
		filterConditions = append(filterConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM product_variants WHERE product_variants.product_id = products.id AND product_variants.color = ANY($%d))", len(args)-2))
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM product_variants WHERE product_variants.product_id = products.id AND product_variants.color = ANY($%d))", len(args)))
	}

	if filter.MaxPrice != nil {
		args = append(args, *filter.MaxPrice)
		filterConditions = append(filterConditions, fmt.Sprintf("base_price <= $%d", len(args)-2))
		conditions = append(conditions, fmt.Sprintf("products.base_price <= $%d", len(args)))
	}

	if filter.MinPrice != nil {
		args = append(args, *filter.MinPrice)
		filterConditions = append(filterConditions, fmt.Sprintf("base_price >= $%d", len(args)-2))
		conditions = append(conditions, fmt.Sprintf("products.base_price >= $%d", len(args)))
	}

	whereClause := ""
	filterWhereClause := ""

	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	if len(filterConditions) > 0 {
		filterWhereClause = " WHERE " + strings.Join(filterConditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.public_id, p.name, p.discount_percentage, p.base_price, p.description, product_variants.id, product_variants.public_id, product_variants.product_id, product_variants.color, product_variants.size, product_variants.image_url, product_variants.stock_quantity
		FROM (SELECT * FROM products %s ORDER BY %s LIMIT $1 OFFSET $2) AS p
		JOIN product_variants ON product_variants.product_id = p.id
		ORDER BY %s
	`, whereClause, productsOrderBy, pOrderBy)

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM products %s
	`, filterWhereClause)

	var count int64

	rows, err := r.db.QueryContext(ctx, query, args...)
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

	countRow := r.db.QueryRowContext(ctx, countQuery, args[2:]...)
	if err := countRow.Scan(&count); err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (r *Repository) GetProductByPublicID(ctx context.Context, publicID string) (*domain.Product, error) {
	query := `
		SELECT p.id, p.public_id, p.name, p.discount_percentage, p.base_price, p.description, product_variants.id, product_variants.public_id, product_variants.product_id, product_variants.color, product_variants.size, product_variants.image_url, product_variants.stock_quantity
		FROM products AS p
		JOIN product_variants ON product_variants.product_id = p.id
		WHERE p.public_id = $1`

	rows, err := r.db.QueryContext(ctx, query, publicID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var product *domain.Product
	for rows.Next() {
		var (
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

		if err := rows.Scan(&productPublicID, &productName, &productDiscountPercentage, &productBasePrice, &productDescription, &variantID, &variantPublicID, &variantProductID, &variantColor, &variantSize, &variantImageURL, &variantStockQuantity); err != nil {
			return nil, err
		}

		if product == nil {
			product = &domain.Product{
				Name:               productName,
				Description:        productDescription,
				PublicID:           productPublicID,
				BasePrice:          productBasePrice,
				DiscountPercentage: productDiscountPercentage,
			}
		}

		product.Variants = append(product.Variants, domain.ProductVariant{
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
		return nil, err
	}

	return product, nil
}
