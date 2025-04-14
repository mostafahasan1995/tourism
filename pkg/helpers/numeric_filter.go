package helpers

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type NumericFilter struct {
	Eq  *int `json:"eq,omitempty"`
	Gt  *int `json:"gt,omitempty"`
	Gte *int `json:"gte,omitempty"`
	Lt  *int `json:"lt,omitempty"`
	Lte *int `json:"lte,omitempty"`
}

func (nf *NumericFilter) BuildFilter() bson.M {
	if nf == nil {
		return nil
	}

	filter := bson.M{}

	if nf.Eq != nil {
		filter["$eq"] = *nf.Eq
	}
	if nf.Gt != nil {
		filter["$gt"] = *nf.Gt
	}
	if nf.Gte != nil {
		filter["$gte"] = *nf.Gte
	}
	if nf.Lt != nil {
		filter["$lt"] = *nf.Lt
	}
	if nf.Lte != nil {
		filter["$lte"] = *nf.Lte
	}

	if len(filter) == 0 {
		return nil
	}
	return filter
}

func (nf *NumericFilter) Validate() error {
	if nf == nil {
		return nil
	}

	// Check for logical consistency in ranges
	if nf.Gt != nil && nf.Gte != nil {
		return errors.New("cannot specify both gt and gte")
	}
	if nf.Lt != nil && nf.Lte != nil {
		return errors.New("cannot specify both lt and lte")
	}

	// Check that ranges make sense
	if nf.Gt != nil && nf.Lt != nil && *nf.Gt >= *nf.Lt {
		return errors.New("gt must be less than lt")
	}
	if nf.Gte != nil && nf.Lte != nil && *nf.Gte > *nf.Lte {
		return errors.New("gte must be less than or equal to lte")
	}

	return nil
}
