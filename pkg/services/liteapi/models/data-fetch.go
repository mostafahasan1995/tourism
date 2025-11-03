package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DataFetch struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Language     string             `json:"language" bson:"language"`
	PlaceId      string             `json:"placeId" bson:"placeId"`
	FullyFetched bool               `json:"fullyFetched" bson:"fullyFetched"`
	TotalCount   int                `json:"totalCount" bson:"totalCount"`
	FetchedCount int                `json:"fetchedCount" bson:"fetchedCount"`
	Count        int                `json:"count" bson:"count"`
}

func (d *DataFetch) IsFullyFetched() bool {
	return d.FetchedCount >= d.TotalCount
}

func (d *DataFetch) IsPartiallyFetched() bool {
	return d.FetchedCount > 0 && d.FetchedCount < d.TotalCount
}

func (d *DataFetch) IsNotFetched() bool {
	return d.FetchedCount == 0
}
