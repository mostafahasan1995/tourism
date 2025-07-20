package models

import (
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/services/exhibition-management/enums"
	marketingModels "larsa-tourism-microservices/pkg/services/marketing/models"
	"larsa-tourism-microservices/pkg/types"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdDto for creating/updating ads
type AdDto struct {
	Title       transl.Localizable[string]   `json:"title" validate:"required"`
	Description transl.Localizable[string]   `json:"description"`
	Tags        []transl.Localizable[string] `json:"tags"`
	Images      []types.FileField            `json:"images" bson:"images"` // Store as ObjectIDs in DB
	WebsiteUrl  string                       `json:"websiteUrl"`
	// Payment & Activation
	PackageType    string    `json:"packageType" validate:"required,oneof=basic premium gold"` // basic, premium, gold
	Duration       int       `json:"duration" validate:"required,min=1"`                       // in days
	PaymentMethod  string    `json:"paymentMethod" validate:"required,oneof=credit_card paypal bank_transfer"`
	PaymentStatus  string    `json:"paymentStatus" validate:"required,oneof=pending paid failed"`
	ActivationDate time.Time `json:"activationDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	Price          float64   `json:"price" validate:"required,min=0"`
	IsActive       bool      `json:"isActive"`
}

// Ad represents an advertisement within an exhibition
type Ad struct {
	Id          primitive.ObjectID           `bson:"_id" json:"_id"`
	Title       transl.Localizable[string]   `bson:"title" json:"title"`
	Description transl.Localizable[string]   `bson:"description" json:"description"`
	Tags        []transl.Localizable[string] `bson:"tags" json:"tags"`
	Images      []types.FileField            `bson:"images" json:"images"` // Store as ObjectIDs in DB
	WebsiteUrl  string                       `bson:"websiteUrl" json:"websiteUrl"`
	// Payment & Activation
	PackageType    string    `bson:"packageType" json:"packageType"`
	Duration       int       `bson:"duration" json:"duration"`
	PaymentMethod  string    `bson:"paymentMethod" json:"paymentMethod"`
	PaymentStatus  string    `bson:"paymentStatus" json:"paymentStatus"`
	ActivationDate time.Time `bson:"activationDate" json:"activationDate"`
	ExpirationDate time.Time `bson:"expirationDate" json:"expirationDate"`
	Price          float64   `bson:"price" json:"price"`
	IsActive       bool      `bson:"isActive" json:"isActive"`
	// Audit fields
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
}

// ExhibitionDto for API requests and updates
type ExhibitionDto struct {
	Title              transl.Localizable[string]   `json:"title" bson:"title" validate:"required"`
	Description        transl.Localizable[string]   `json:"description" bson:"description" validate:"required"`
	LocationType       transl.Localizable[string]   `json:"locationType" bson:"locationType" validate:"required"`
	LocationAddress    transl.Localizable[string]   `json:"locationAddress" bson:"locationAddress"`
	StartDate          time.Time                    `json:"startDate" bson:"startDate" validate:"required"`
	EndDate            time.Time                    `json:"endDate" bson:"endDate" validate:"required"`
	Images             []types.FileField            `json:"images" bson:"images"`
	RelatedExhibitions []primitive.ObjectID         `json:"relatedExhibitions" bson:"relatedExhibitions"`
	Tags               []transl.Localizable[string] `json:"tags" bson:"tags"`
	IsActive           bool                         `json:"isActive" bson:"isActive"`
	Ads                []AdDto                      `json:"ads" bson:"ads"` // Ads can be included in exhibition creation
	Status             string                       `json:"status" bson:"status" validate:"required,oneof=active closed upcoming"`
}

// Exhibition main model with nested ads
type Exhibition struct {
	Id primitive.ObjectID `bson:"_id" json:"_id"`
	ExhibitionDto
	Ads       []Ad               `bson:"ads" json:"ads"` // Nested ads array
	IsFav     bool               `bson:"isFav" json:"isFav"`
	Trash     bool               `bson:"trash" json:"trash"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
}

// ExhibitionWithPagination for paginated responses
type ExhibitionWithPagination struct {
	Exhibitions []Exhibition     `json:"exhibitions"`
	Pagination  types.Pagination `json:"pagination"`
}

// ExhibitionStats for dashboard statistics
type ExhibitionStats struct {
	ActiveCount        int64 `json:"activeCount"`
	UpcomingCount      int64 `json:"upcomingCount"`
	ClosedCount        int64 `json:"closedCount"`
	TotalCount         int64 `json:"totalCount"`
	RegisteredVisitors int64 `json:"registeredVisitors"`
}

// ExhibitionWithVisitorStats for detailed exhibition statistics including visitor data
type ExhibitionWithVisitorStats struct {
	Exhibition     Exhibition                    `json:"exhibition"`
	VisitorCount   int64                         `json:"visitorCount"`
	ActiveVisitors int64                         `json:"activeVisitors"`
	VIPVisitors    int64                         `json:"vipVisitors"`
	VisitorStats   *marketingModels.VisitorStats `json:"visitorStats,omitempty"`
}

// Validation methods
func (e *ExhibitionDto) Validate(v *validator.Validate) error {
	if err := v.Struct(e); err != nil {
		return err
	}

	// Custom validation: end date should be after start date
	if e.EndDate.Before(e.StartDate) {
		return helpers.BadRequest("End date must be after start date")
	}

	// Custom validation: status must be one of the valid enum values
	if !enums.IsValidExhibitionStatus(e.Status) {
		return helpers.BadRequest("Status must be one of: active, closed, upcoming")
	}

	// Validate each ad if provided
	for _, ad := range e.Ads {
		if err := ad.Validate(v); err != nil {
			return err
		}
	}

	return nil
}

func (a *AdDto) Validate(v *validator.Validate) error {
	if err := v.Struct(a); err != nil {
		return err
	}

	// Custom validation: expiration date should be after activation date
	if !a.ExpirationDate.IsZero() && !a.ActivationDate.IsZero() {
		if a.ExpirationDate.Before(a.ActivationDate) {
			return helpers.BadRequest("Expiration date must be after activation date")
		}
	}

	return nil
}

// Helper method to calculate ad pricing based on package and duration
func CalculateAdPrice(packageType string, duration int) float64 {
	basePrice := map[string]float64{
		"basic":   50.0,
		"premium": 100.0,
		"gold":    200.0,
	}

	if price, exists := basePrice[packageType]; exists {
		return price * float64(duration)
	}
	return 0
}

// GetExhibitionStatusByDates automatically determines exhibition status based on current date and exhibition dates
func GetExhibitionStatusByDates(startDate, endDate time.Time) string {
	now := time.Now()

	if now.Before(startDate) {
		return enums.ExhibitionStatusUpcoming
	} else if now.After(endDate) {
		return enums.ExhibitionStatusClosed
	} else {
		return enums.ExhibitionStatusActive
	}
}

// SetAutoStatus automatically sets the status based on start and end dates if no status is provided
func (e *ExhibitionDto) SetAutoStatus() {
	if e.Status == "" {
		e.Status = GetExhibitionStatusByDates(e.StartDate, e.EndDate)
	}
}
