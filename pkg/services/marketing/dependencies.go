package marketing

import (
	"larsa-tourism-microservices/pkg/services/marketing/repo"

	"github.com/samber/do"
)

func RegisterMarketingDependencies(i *do.Injector) {
	// Register repositories
	do.Provide(i, repo.NewExhibitorProfileRepo)
	do.Provide(i, repo.NewInquiryRepo)
	do.Provide(i, repo.NewVisitorRepo)
	do.Provide(i, repo.NewVisitorActivityRepo)

	// Register services
	do.Provide(i, NewExhibitorProfileSvcs)
	do.Provide(i, NewInquirySvcs)
	do.Provide(i, NewVisitorSvcs)
}
