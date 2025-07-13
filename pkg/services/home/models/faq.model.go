package models

import (
	"time"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FAQ Page represents a static FAQ section (About, Our Agents, etc.)
type FaqPage struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `bson:"name" json:"name"`               // "About", "Our Agents", etc.
	Slug        string             `bson:"slug" json:"slug"`               // URL-friendly version
	Description string             `bson:"description" json:"description"` // Optional description
	IsActive    bool               `bson:"isActive" json:"isActive"`       // Enable/disable page
	SortOrder   int                `bson:"sortOrder" json:"sortOrder"`     // Display order
	Trash       bool               `bson:"trash" json:"trash"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
}

// FAQ Group represents a group of questions within a FAQ page
type FaqGroup struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	FaqPageId   primitive.ObjectID `bson:"faqPageId" json:"faqPageId"`     // Reference to FAQ page
	Name        string             `bson:"name" json:"name"`               // "General Questions", etc.
	Description string             `bson:"description" json:"description"` // Optional description
	IsActive    bool               `bson:"isActive" json:"isActive"`       // Enable/disable group
	SortOrder   int                `bson:"sortOrder" json:"sortOrder"`     // Display order within page
	Trash       bool               `bson:"trash" json:"trash"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
}

// FAQ Question represents a question-answer pair
type FaqQuestion struct {
	Id         primitive.ObjectID  `bson:"_id" json:"_id"`
	FaqPageId  primitive.ObjectID  `bson:"faqPageId" json:"faqPageId"`   // Reference to FAQ page
	FaqGroupId *primitive.ObjectID `bson:"faqGroupId" json:"faqGroupId"` // Optional: null for general questions
	Question   string              `bson:"question" json:"question"`     // The question text
	Answer     string              `bson:"answer" json:"answer"`         // The answer text
	IsActive   bool                `bson:"isActive" json:"isActive"`     // Enable/disable question
	SortOrder  int                 `bson:"sortOrder" json:"sortOrder"`   // Display order within group/page
	Trash      bool                `bson:"trash" json:"trash"`
	CreatedAt  time.Time           `bson:"createdAt" json:"createdAt"`
	CreatedBy  primitive.ObjectID  `bson:"createdBy" json:"createdBy"`
	UpdatedAt  time.Time           `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy  primitive.ObjectID  `bson:"updatedBy" json:"updatedBy"`
}

// DTOs for API requests
type FaqPageDto struct {
	Name        string `json:"name" validate:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
	SortOrder   int    `json:"sortOrder"`
}

type FaqGroupDto struct {
	FaqPageId   primitive.ObjectID `json:"faqPageId,omitempty"` // Optional - can be provided via URL
	Name        string             `json:"name" validate:"required"`
	Description string             `json:"description"`
	IsActive    bool               `json:"isActive"`
	SortOrder   int                `json:"sortOrder"`
}

type FaqQuestionDto struct {
	FaqPageId  primitive.ObjectID  `json:"faqPageId,omitempty"`  // Optional - can be provided via URL
	FaqGroupId *primitive.ObjectID `json:"faqGroupId,omitempty"` // Optional - can be provided via URL or null for general questions
	Question   string              `json:"question" validate:"required"`
	Answer     string              `json:"answer" validate:"required"`
	IsActive   bool                `json:"isActive"`
	SortOrder  int                 `json:"sortOrder"`
}

// Response models with populated data
type FaqPageWithStats struct {
	FaqPage
	GroupCount    int `json:"groupCount"`
	QuestionCount int `json:"questionCount"`
}

type FaqGroupWithQuestions struct {
	FaqGroup
	Questions []FaqQuestion `json:"questions"`
}

type FaqPageComplete struct {
	FaqPage
	GeneralQuestions []FaqQuestion           `json:"generalQuestions"` // Questions not in any group
	Groups           []FaqGroupWithQuestions `json:"groups"`           // Groups with their questions
}

// Search result model for returning matching questions/answers
type FaqSearchResult struct {
	Id        string    `json:"id"`
	PageId    string    `json:"pageId"`
	GroupId   string    `json:"groupId,omitempty"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	IsActive  bool      `json:"isActive"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Pagination models
type FaqPagePagination struct {
	FaqPages   []FaqPageWithStats `json:"faqPages"`
	Pagination common.Pagination  `json:"pagination"`
}

type FaqGroupPagination struct {
	FaqGroups  []FaqGroup        `json:"faqGroups"`
	Pagination common.Pagination `json:"pagination"`
}

type FaqQuestionPagination struct {
	FaqQuestions []FaqQuestion     `json:"faqQuestions"`
	Pagination   common.Pagination `json:"pagination"`
}

// Search result pagination
type FaqSearchResultPagination struct {
	FaqSearchResults []FaqSearchResult `json:"faqSearchResults"`
	Pagination       common.Pagination `json:"pagination"`
}
