package models

import (
	"context"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"larsa-tourism-microservices/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FaveType string

const (
	FaveTypeProgram       FaveType = "program"       // DONE
	FaveTypeHotel         FaveType = "hotel"         // DONE
	FaveTypeTravelerStory FaveType = "travelerStory" // DONE
	FaveTypeClientStory   FaveType = "clientStory"   // DONE
	FaveTypeExhibition    FaveType = "exhibition"    // DONE
	FaveTypeAgent         FaveType = "agent"         // DONE
	FaveTypeDestination   FaveType = "destination"   // INPROGRESS
)

type FaveDto struct {
	Type  FaveType           `bson:"type" json:"type" validate:"required,oneof=program hotel travelerStory clientStory exhibition agent destination"` //program - hotel - diary - exhibition - agent - destination
	RefId primitive.ObjectID `bson:"refId" json:"refId" validate:"required"`
	IsFav bool               `bson:"isFav" json:"isFav"`
}

type Fave struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId" json:"userId"`
	FaveDto   `bson:",inline"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type FaveItem struct {
	Fave `bson:",inline"`
	Item any `bson:"item" json:"item"`
}

type FavePagination struct {
	Faves      []Fave           `bson:"faves" json:"faves"`
	Pagination types.Pagination `bson:"pagination" json:"pagination"`
}

func BuildFavoritePipelineWithAuth(ctx context.Context, faveType FaveType) []bson.M {
	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		return buildFavoritePipeline(cfg.User.Id, faveType)
	} else {
		return []bson.M{buildDefaultFavorite()}
	}
}

// buildFavoritePipeline create a pipeline to add favorite to items
// Used with any favtype (program, hotel, etc.)
func buildFavoritePipeline(userId primitive.ObjectID, faveType FaveType) []bson.M {
	return []bson.M{
		{
			"$lookup": bson.M{
				"from": "tourismFavorites",
				"let":  bson.M{"itemId": "$_id"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []bson.M{
									{"$eq": []interface{}{"$refId", "$$itemId"}},
									{"$eq": []interface{}{"$type", faveType}},
									{"$eq": []interface{}{"$userId", userId}},
									{"$eq": []interface{}{"$isFav", true}},
								},
							},
							"trash": bson.M{"$ne": true},
						},
					},
				},
				"as": "faveRecord",
			},
		},
		{
			"$addFields": bson.M{
				"isFav": bson.M{
					"$gt": []interface{}{
						bson.M{"$size": "$faveRecord"},
						0,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"faveRecord": 0, // Remove temporary field
			},
		},
	}
}

// buildDefaultFavorite adds a default isFav field set to false
// Use this when user is not authenticated
func buildDefaultFavorite() bson.M {
	return bson.M{
		"$addFields": bson.M{
			"isFav": false,
		},
	}
}
