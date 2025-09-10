package models

type HotelDetails map[string]any

type HotelDetailsData struct {
	Data HotelDetails `json:"data"`
}
