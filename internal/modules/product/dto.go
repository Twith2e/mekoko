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
	Order    string   `form:"order" binding:"omitempty,oneof=price_asc price_desc oldest_first newest_first"`
	Color    []string `form:"color" binding:"omitempty"`
	MinPrice *int64   `form:"min_price" binding:"omitempty,gt=0"`
	MaxPrice *int64   `form:"max_price" binding:"omitempty,gt=0"`
	Page     int      `form:"page" binding:"omitempty,min=1"`
	Limit    int      `form:"limit" binding:"omitempty,min=1,max=100"`
}
