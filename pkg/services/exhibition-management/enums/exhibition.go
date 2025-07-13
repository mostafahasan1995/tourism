package enums

// Exhibition status constants
const (
	ExhibitionStatusActive   = "active"
	ExhibitionStatusClosed   = "closed"
	ExhibitionStatusUpcoming = "upcoming"
)

// GetValidExhibitionStatuses returns all valid exhibition status values
func GetValidExhibitionStatuses() []string {
	return []string{
		ExhibitionStatusActive,
		ExhibitionStatusClosed,
		ExhibitionStatusUpcoming,
	}
}

// IsValidExhibitionStatus checks if the provided status is valid
func IsValidExhibitionStatus(status string) bool {
	validStatuses := GetValidExhibitionStatuses()
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}
