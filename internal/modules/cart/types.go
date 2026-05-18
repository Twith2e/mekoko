package cart

type CartForUI struct {
	ID                   string `json:"id"`
	VariantID            string `json:"variant_id"`
	Quantity             int64  `json:"quantity"`
	UnitPriceAtSelection int64  `json:"unit_price_at_selection"`
	ImageURL             string `json:"image_url"`
	Name                 string `json:"name"`
	Color                string `json:"color"`
}
