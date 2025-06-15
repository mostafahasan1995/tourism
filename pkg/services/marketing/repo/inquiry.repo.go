package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/marketing/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type InquiryRepo interface {
	dbrepo.MainRepo[models.Inquiry]
}

type inquiryRepo struct {
	dbrepo.MainRepoImpl[models.Inquiry]
}

func NewInquiryRepo(i *do.Injector) (InquiryRepo, error) {
	return &inquiryRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Inquiry]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "marketingInquiries",
		},
	}, nil
}
