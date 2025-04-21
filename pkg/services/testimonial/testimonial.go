package testimonial

import (
	"context"
	"larsa-tourism-microservices/pkg/services/testimonial/models"
	"larsa-tourism-microservices/pkg/services/testimonial/repo"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestimonialService interface {
	All(ctx context.Context) ([]models.Testimonial, error)
	Add(ctx context.Context, testimonial *models.TestimonialDto) (*models.Testimonial, error)
}

type testimonialsvcs struct {
	repo repo.TestimonialRepo
}

func NewTestimonialService(i *do.Injector) (TestimonialService, error) {
	return &testimonialsvcs{
		repo: do.MustInvoke[repo.TestimonialRepo](i),
	}, nil
}

func (t *testimonialsvcs) All(ctx context.Context) ([]models.Testimonial, error) {
	var result []models.Testimonial
	err := t.repo.Aggregate(ctx, []bson.M{}, func(cur *mongo.Cursor) error {
		if err := cur.All(ctx, &result); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *testimonialsvcs) Add(ctx context.Context, testimonial *models.TestimonialDto) (*models.Testimonial, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	testiminial := &models.Testimonial{
		Id:             primitive.NewObjectID(),
		UserId:         cfg.User.Id,
		TestimonialDto: *testimonial,
		CreatedAt:      time.Now(),
	}
	if err := t.repo.Add(ctx, testiminial); err != nil {
		return nil, err
	}

	return testiminial, nil
}
