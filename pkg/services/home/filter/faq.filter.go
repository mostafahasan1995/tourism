package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FaqPageFilter for filtering FAQ pages
type FaqPageFilter struct {
	Page     int    `bson:"page" json:"page"`
	Size     int    `bson:"size" json:"size"`
	Name     string `bson:"name" json:"name"`         // Search by name
	IsActive *bool  `bson:"isActive" json:"isActive"` // Filter by active status
}

func (f *FaqPageFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add name search if provided
	if f.Name != "" {
		filterConditions = append(filterConditions, bson.M{
			"name": bson.M{"$regex": f.Name, "$options": "i"},
		})
	}

	// Add active status filter if provided
	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	return bson.M{"$and": filterConditions}
}

// FaqGroupFilter for filtering FAQ groups
type FaqGroupFilter struct {
	Page      int                `bson:"page" json:"page"`
	Size      int                `bson:"size" json:"size"`
	FaqPageId primitive.ObjectID `bson:"faqPageId" json:"faqPageId"` // Filter by FAQ page
	Name      string             `bson:"name" json:"name"`           // Search by name
	Search    string             `bson:"search" json:"search"`       // Search in questions and answers
	IsActive  *bool              `bson:"isActive" json:"isActive"`   // Filter by active status
}

func (f *FaqGroupFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add FAQ page filter if provided
	if !f.FaqPageId.IsZero() {
		filterConditions = append(filterConditions, bson.M{"faqPageId": f.FaqPageId})
	}

	// Add name search if provided
	if f.Name != "" {
		filterConditions = append(filterConditions, bson.M{
			"name": bson.M{"$regex": f.Name, "$options": "i"},
		})
	}

	// Add active status filter if provided
	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	return bson.M{"$and": filterConditions}
}

// HasQuestionSearch returns true if the filter includes search in questions/answers
func (f *FaqGroupFilter) HasQuestionSearch() bool {
	return f.Search != ""
}

// GetQuestionSearchTerm returns the search term for questions/answers
func (f *FaqGroupFilter) GetQuestionSearchTerm() string {
	return f.Search
}

// FaqQuestionFilter for filtering FAQ questions
type FaqQuestionFilter struct {
	Page       int                 `bson:"page" json:"page"`
	Size       int                 `bson:"size" json:"size"`
	FaqPageId  primitive.ObjectID  `bson:"faqPageId" json:"faqPageId"`   // Filter by FAQ page
	FaqGroupId *primitive.ObjectID `bson:"faqGroupId" json:"faqGroupId"` // Filter by group (null for general)
	Question   string              `bson:"question" json:"question"`     // Search in question text
	Answer     string              `bson:"answer" json:"answer"`         // Search in answer text
	Search     string              `bson:"search" json:"search"`         // Search in both question and answer
	IsActive   *bool               `bson:"isActive" json:"isActive"`     // Filter by active status
}

func (f *FaqQuestionFilter) ToBsonFilter() bson.M {
	filterConditions := []bson.M{
		{"trash": bson.M{"$ne": true}},
	}

	// Add FAQ page filter if provided
	if !f.FaqPageId.IsZero() {
		filterConditions = append(filterConditions, bson.M{"faqPageId": f.FaqPageId})
	}

	// Add FAQ group filter if provided
	if f.FaqGroupId != nil {
		if f.FaqGroupId.IsZero() {
			// Filter for general questions (no group)
			filterConditions = append(filterConditions, bson.M{
				"$or": []bson.M{
					{"faqGroupId": bson.M{"$exists": false}},
					{"faqGroupId": nil},
				},
			})
		} else {
			// Filter for specific group
			filterConditions = append(filterConditions, bson.M{"faqGroupId": *f.FaqGroupId})
		}
	}

	// Add question search if provided
	if f.Question != "" {
		filterConditions = append(filterConditions, bson.M{
			"question": bson.M{"$regex": f.Question, "$options": "i"},
		})
	}

	// Add answer search if provided
	if f.Answer != "" {
		filterConditions = append(filterConditions, bson.M{
			"answer": bson.M{"$regex": f.Answer, "$options": "i"},
		})
	}

	// Add general search if provided (searches both question and answer)
	if f.Search != "" {
		filterConditions = append(filterConditions, bson.M{
			"$or": []bson.M{
				{"question": bson.M{"$regex": f.Search, "$options": "i"}},
				{"answer": bson.M{"$regex": f.Search, "$options": "i"}},
			},
		})
	}

	// Add active status filter if provided
	if f.IsActive != nil {
		filterConditions = append(filterConditions, bson.M{"isActive": *f.IsActive})
	}

	return bson.M{"$and": filterConditions}
}
