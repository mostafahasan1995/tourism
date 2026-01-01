package enums

type PackageStatus string

const (
	PackageStatusActive   PackageStatus = "active"
	PackageStatusInactive PackageStatus = "inactive"
	PackageStatusDraft    PackageStatus = "draft"
	PackageStatusArchived PackageStatus = "archived"
)
