package product

type FilterType string

const (
	FilterPriceASC    FilterType = "price_asc"
	FilterPriceDESC   FilterType = "price_desc"
	FilterOldestFirst FilterType = "oldest_first"
	FilterNewestFirst FilterType = "newest_first"
	FilterDefault     FilterType = "id"
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

type Filter struct {
	Order    FilterType
	Color    []string
	MinPrice *int64
	MaxPrice *int64
}
