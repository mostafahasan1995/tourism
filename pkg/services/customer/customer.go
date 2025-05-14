package customer

import (
	"context"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/services/customer/models"
	"larsa-tourism-microservices/pkg/services/customer/repo"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/util"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerSvcs interface {
	Add(ctx context.Context, data *models.CustomerDto) (*models.Customer, error)
	Update(ctx context.Context, customerId string, data *models.CustomerDto) (*models.Customer, error)
}

type customerSvcs struct {
	repo        repo.CustomerRepo
	sortingsvcs dbsvcs.SortingSvcs
	withtxn     *db.WithTxn
}

func NewCustomerSvcs(i *do.Injector) (CustomerSvcs, error) {
	return &customerSvcs{
		repo:        do.MustInvoke[repo.CustomerRepo](i),
		sortingsvcs: do.MustInvoke[dbsvcs.SortingSvcs](i),
		withtxn:     do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (c *customerSvcs) Add(ctx context.Context, data *models.CustomerDto) (*models.Customer, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	result, err := c.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		seq, err := c.sortingsvcs.GetAndUpdateSourceSeq(ctx, "customer")
		if err != nil {
			return nil, err
		}

		customerId := fmt.Sprintf("CUSTOMER-%d", seq)

		customer := &models.Customer{
			Id:          primitive.NewObjectID(),
			CustomerId:  customerId,
			CustomerDto: *data,
			CreatedAt:   time.Now(),
			CreatedBy:   cfg.User.Id,
		}

		if err := c.repo.Add(ctx, customer); err != nil {
			return nil, err
		}

		return customer, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Customer), nil
}

func (c *customerSvcs) Update(ctx context.Context, customerId string, data *models.CustomerDto) (*models.Customer, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, err
	}

	result, err := c.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {

		customer := &models.Customer{
			Id:          _id,
			CustomerDto: *data,
			UpdatedAt:   time.Now(),
			UpdatedBy:   cfg.User.Id,
		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": customer}

		updatedCustomer, err := c.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedCustomer, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Customer), nil

}
