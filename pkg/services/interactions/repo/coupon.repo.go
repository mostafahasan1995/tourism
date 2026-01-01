package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type CouponRepo interface {
	dbrepo.MainRepo[models.Coupon]
}

type couponrepo struct {
	dbrepo.MainRepoImpl[models.Coupon]
}

func NewCouponRepo(i *do.Injector) (CouponRepo, error) {
	return &couponrepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.Coupon]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "tourismCoupons",
		},
	}, nil
}
