package cart

type CartForUI struct {
	ID                   string
	VariantID            string
	Quantity             int64
	UnitPriceAtSelection int64
	ImageURL             string
	Name                 string
	Color                string
}
