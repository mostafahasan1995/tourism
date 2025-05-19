package member

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/gateway"
	gwmodels "larsa-tourism-microservices/pkg/gateway/models"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/member/filters"
	"larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/member/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (*models.Customer, error)
	GetOne(ctx context.Context, customerId string) (*models.Customer, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.CustomerWithPagination, error)
	Add(ctx context.Context, data *models.CustomerDto) (*models.Customer, error)
	Update(ctx context.Context, customerId string, data *models.CustomerDto) (*models.Customer, error)
	Delete(ctx context.Context, customerId string) error
}

type customerSvcs struct {
	repo        repo.CustomerRepo
	sortingsvcs dbsvcs.SortingSvcs
	usersgw     *gateway.UsersGw
	withtxn     *db.WithTxn
}

func NewCustomerSvcs(i *do.Injector) (CustomerSvcs, error) {
	return &customerSvcs{
		repo:        do.MustInvoke[repo.CustomerRepo](i),
		sortingsvcs: do.MustInvoke[dbsvcs.SortingSvcs](i),
		usersgw:     do.MustInvoke[*gateway.UsersGw](i),
		withtxn:     do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (c *customerSvcs) GetOne(ctx context.Context, customerId string) (*models.Customer, error) {
	_id, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, err
	}

	return c.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (c *customerSvcs) GetByFilter(ctx context.Context, filter bson.M) (*models.Customer, error) {
	return c.repo.GetByFilter(ctx, filter)
}

func (c *customerSvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.CustomerWithPagination, error) {
	match := bson.M{"trash": false}

	filters, err := filters.NewCustomerFilter(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := c.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Customer
	errAg := c.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.CustomerWithPagination{
		Customers:  result,
		Pagination: pagination,
	}, nil
}

func (c *customerSvcs) Add(ctx context.Context, data *models.CustomerDto) (*models.Customer, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	result, err := c.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		customer := &models.Customer{
			CustomerDto: *data,
			CreatedAt:   time.Now(),
			CreatedBy:   cfg.User.Id,
		}

		userId, err := c.AddCustomerCredentials(ctx, customer)
		if err != nil {
			return nil, err
		}

		customer.Id = userId

		seq, err := c.sortingsvcs.GetAndUpdateSourceSeq(ctx, "customer")
		if err != nil {
			return nil, err
		}

		customer.CustomerId = fmt.Sprintf("CUSTOMER-%d", seq)

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

	_id, err := primitive.ObjectIDFromHex(customerId) //same as user id
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

		_, err := c.UpdateCustomerCredentials(ctx, customer)
		if err != nil {
			return nil, err
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

func (c *customerSvcs) Delete(ctx context.Context, customerId string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = c.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *customerSvcs) AddCustomerCredentials(ctx context.Context, data *models.Customer) (userId primitive.ObjectID, err error) {
	password := data.Security.NewPassword
	if password == "" {
		password = util.GeneratePassword(8, 2, 2, 2)
	}

	user := &gwmodels.PostUserData{
		FirstName: data.Name,
		LastName:  "-",
		Email:     data.Security.Email,
		Password:  password,
		// Roles:        []primitive.ObjectID{}, //empty for default role
		// Capabilities: []primitive.ObjectID{},
	}

	return c.usersgw.AddUser(ctx, user)
}

func (c *customerSvcs) UpdateCustomerCredentials(ctx context.Context, data *models.Customer) (userId primitive.ObjectID, err error) {
	user := &gwmodels.PostUserData{
		FirstName: data.Name,
		LastName:  "-",
		Email:     data.Security.Email,
	}

	if data.Security.NewPassword != "" {
		user.Password = data.Security.NewPassword
	}

	return c.usersgw.UpdateUser(ctx, data.Id.Hex(), user)
}
