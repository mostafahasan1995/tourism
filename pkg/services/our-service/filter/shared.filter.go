package filter



type PriceRangeFilter struct {
	From int `bson:"from" json:"from"`
	To   int `bson:"to" json:"to"`
}