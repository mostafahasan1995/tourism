package filter

import (
	"larsa-tourism-microservices/pkg/services/marketing/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InquiryFilter struct {
	Page        int                   `bson:"page" json:"page"`
	Size        int                   `bson:"size" json:"size"`
	HotelId     primitive.ObjectID    `bson:"hotelId,omitempty" json:"hotelId,omitempty"`
	Status      *models.InquiryStatus `bson:"status,omitempty" json:"status,omitempty"`
	Priority    string                `bson:"priority" json:"priority"`
	Source      string                `bson:"source" json:"source"`
	IsRead      *bool                 `bson:"isRead,omitempty" json:"isRead,omitempty"`
	VisitorName string                `bson:"visitorName" json:"visitorName"`
	Email       string                `bson:"email" json:"email"`
	DateFrom    *time.Time            `bson:"dateFrom,omitempty" json:"dateFrom,omitempty"`
	DateTo      *time.Time            `bson:"dateTo,omitempty" json:"dateTo,omitempty"`
	HasReplies  *bool                 `bson:"hasReplies,omitempty" json:"hasReplies,omitempty"`
	SearchText  string                `bson:"searchText" json:"searchText"`
}

func (f *InquiryFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add hotel ID filter if provided
	if !f.HotelId.IsZero() {
		filterConditions = append(filterConditions, bson.M{"hotelId": f.HotelId})
	}

	// Add status filter if provided
	if f.Status != nil {
		filterConditions = append(filterConditions, bson.M{"status": *f.Status})
	}

	// Add priority filter if provided
	if f.Priority != "" {
		filterConditions = append(filterConditions, bson.M{"priority": f.Priority})
	}

	// Add source filter if provided
	if f.Source != "" {
		filterConditions = append(filterConditions, bson.M{"source": f.Source})
	}

	// Add read status filter if provided
	if f.IsRead != nil {
		filterConditions = append(filterConditions, bson.M{"isRead": *f.IsRead})
	}

	// Add visitor name search if provided
	if f.VisitorName != "" {
		filterConditions = append(filterConditions, bson.M{
			"visitorName": bson.M{"$regex": f.VisitorName, "$options": "i"},
		})
	}

	// Add email search if provided
	if f.Email != "" {
		filterConditions = append(filterConditions, bson.M{
			"email": bson.M{"$regex": f.Email, "$options": "i"},
		})
	}

	// Add date range filter if provided
	if f.DateFrom != nil && f.DateTo != nil {
		filterConditions = append(filterConditions, bson.M{
			"createdAt": bson.M{
				"$gte": *f.DateFrom,
				"$lte": *f.DateTo,
			},
		})
	} else if f.DateFrom != nil {
		filterConditions = append(filterConditions, bson.M{
			"createdAt": bson.M{"$gte": *f.DateFrom},
		})
	} else if f.DateTo != nil {
		filterConditions = append(filterConditions, bson.M{
			"createdAt": bson.M{"$lte": *f.DateTo},
		})
	}

	// Add replies filter if provided
	if f.HasReplies != nil {
		if *f.HasReplies {
			filterConditions = append(filterConditions, bson.M{
				"replies": bson.M{"$exists": true, "$not": bson.M{"$size": 0}},
			})
		} else {
			filterConditions = append(filterConditions, bson.M{
				"$or": []bson.M{
					{"replies": bson.M{"$exists": false}},
					{"replies": bson.M{"$size": 0}},
				},
			})
		}
	}

	// Add text search if provided
	if f.SearchText != "" {
		filterConditions = append(filterConditions, bson.M{
			"$or": []bson.M{
				{"message": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"subject": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"visitorName": bson.M{"$regex": f.SearchText, "$options": "i"}},
				{"email": bson.M{"$regex": f.SearchText, "$options": "i"}},
			},
		})
	}

	return bson.M{"$and": filterConditions}
}
