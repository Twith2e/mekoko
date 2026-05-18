package product

const (
	FilterPriceASC  = "price_asc"
	FilterPriceDESC = "price_desc"
	OldestFirst     = "oldest_first"
	NewestFirst     = "newest_first"
	IdFilter        = "id"
)

type NewProduct struct {
	PublicID           string
	Name               string
	Description        string
	DiscountPercentage int
	BasePrice          int64
}

type NewVariant struct {
	PublicID      string
	ProductID     int64
	Color         string
	StockQuantity int64
	Size          string
	ImageURL      string
}
