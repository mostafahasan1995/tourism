package filter

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

// PriceRangeFilter represents a range of prices for filtering
type PriceRangeFilter struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Validate checks if the price range is valid
func (p *PriceRangeFilter) Validate() error {
	if p.From < 0 || p.To < 0 {
		return errors.New("price range values cannot be negative")
	}
	if p.From > 0 && p.To > 0 && p.From > p.To {
		return errors.New("from price cannot be greater than to price")
	}
	return nil
}

// BuildPriceFilter builds a MongoDB filter for price range
func (p *PriceRangeFilter) BuildPriceFilter() bson.M {
	if p.From <= 0 && p.To <= 0 {
		return nil
	}

	priceFilter := bson.M{}
	if p.From > 0 {
		priceFilter["$gte"] = p.From
	}
	if p.To > 0 {
		priceFilter["$lte"] = p.To
	}
	return priceFilter
}
