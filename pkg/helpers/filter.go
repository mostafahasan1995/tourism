package helpers

import (
	"encoding/json"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
)

type Filtering interface {
	BuildPipeline(m bson.M) []bson.M
}

// ParseFilters converts a query of any type to the specified filter type T.
// It handles:
// - *T: returns it directly
// - T: converts to *T
// - string: unmarshals JSON to *T
// - map: marshals to JSON then unmarshals to *T
// - other types: returns a new empty *T
func ParseFilters[T Filtering](query any) (*T, error) {
	// Case 1: query is already *T
	if t, ok := query.(*T); ok {
		return t, nil
	}

	// Case 2: query is T (non-pointer)
	t := new(T)
	targetType := reflect.TypeOf(*t)
	if reflect.TypeOf(query) == targetType {
		val := reflect.ValueOf(query)
		converted := reflect.New(targetType).Elem()
		converted.Set(val)
		result := converted.Interface().(T)
		return &result, nil
	}

	// Case 3: query is a string (JSON)
	if str, ok := query.(string); ok {
		if str == "" {
			return new(T), nil
		}

		t := new(T)
		if err := json.Unmarshal([]byte(str), t); err != nil {
			return nil, err
		}
		return t, nil
	}

	// Case 4: query is a map or other convertible type
	if query != nil {
		jsonBytes, err := json.Marshal(query)
		if err == nil {
			t := new(T)
			if err := json.Unmarshal(jsonBytes, t); err == nil {
				return t, nil
			}
		}
	}

	// Default: return empty filter
	return new(T), nil
}
