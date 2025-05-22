package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/invoice/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceRepo interface {
	dbrepo.MainRepo[models.Invoice]
}

type invoicerepo struct {
	dbrepo.MainRepoImpl[models.Invoice]
}

func NewInvoiceRepo(i *do.Injector) (InvoiceRepo, error) {
	return &invoicerepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Invoice]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismInvoices",
		},
	}, nil
}
