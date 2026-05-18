package product

type Variant struct {
	Color         string `json:"color" binding:"required"`
	ImageURL      string `json:"image_url" binding:"required"`
	Size          string `json:"size"`
	StockQuantity int64  `json:"stock_quantity" binding:"omitempty"`
}

type VariantResponse struct {
	ID            string `json:"id"`
	Color         string `json:"color"`
	Size          string `json:"size"`
	ImageURL      string `json:"image_url"`
	StockQuantity int64  `json:"stock_quantity"`
}

type AddProductsRequest struct {
	Name               string    `json:"name" binding:"required"`
	Description        string    `json:"description" binding:"required"`
	BasePrice          int64     `json:"base_price" binding:"required,gt=0"`
	DiscountPercentage int       `json:"discount_percentage" binding:"omitempty"`
	Variants           []Variant `json:"variants" binding:"required"`
}

type GetProductsResponse struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name" binding:"required"`
	Description        string            `json:"description" binding:"required"`
	BasePrice          int64             `json:"base_price" binding:"required,gt=0"`
	DiscountPercentage int               `json:"discount_percentage" binding:"omitempty"`
	Variants           []VariantResponse `json:"variants" binding:"required"`
}

type GetProductsQuery struct {
	Filter string `form:"filter" binding:"omitempty,oneof=price_asc price_desc oldest_first newest_first"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
}
