package repo

import (
	dbrepo "larsa-tourism-microservices/pkg/services/db/repo"
	"larsa-tourism-microservices/pkg/services/home/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpiritualGiftRepo interface {
	dbrepo.MainRepo[models.SpiritualGift]
}
type spiritualGiftRepo struct {
	dbrepo.MainRepoImpl[models.SpiritualGift]
}

func NewSpiritualGiftRepo(i *do.Injector) (SpiritualGiftRepo, error) {
	return &spiritualGiftRepo{
		MainRepoImpl: dbrepo.MainRepoImpl[models.SpiritualGift]{
			Db:       do.MustInvoke[*mongo.Client](i),
			CollName: "spiritualGifts",
		},
	}, nil
}
