package member

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/gateway"
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

var ErrUserNotFound = errors.New("user not found")
var ErrDupliateEmail = errors.New("duplicated email")
var ErrUnKnowen = errors.New("unknowen error")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")

type CustomerSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (*models.Customer, error)
	GetOne(ctx context.Context, customerId string) (*models.Customer, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.CustomerWithPagination, error)
	GetAll(ctx context.Context) ([]models.Customer, error)
	Add(ctx context.Context, data *models.CustomerDto) (*models.Customer, error)
	Update(ctx context.Context, customerId string, data *models.CustomerDto) (*models.Customer, error)
	Delete(ctx context.Context, customerId string) error
	//
	RegisterAsCustomer(ctx context.Context, data *models.CustomerRegisterData) (*models.Customer, error)
}

type customerSvcs struct {
	repo        repo.CustomerRepo
	sortingsvcs dbsvcs.SortingSvcs
	usersgw     *gateway.UsersGw
	gateway     gateway.Gateway
	withtxn     *db.WithTxn
}

func NewCustomerSvcs(i *do.Injector) (CustomerSvcs, error) {
	return &customerSvcs{
		repo:        do.MustInvoke[repo.CustomerRepo](i),
		sortingsvcs: do.MustInvoke[dbsvcs.SortingSvcs](i),
		usersgw:     do.MustInvoke[*gateway.UsersGw](i),
		gateway:     do.MustInvoke[gateway.Gateway](i),
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

func (c *customerSvcs) GetAll(ctx context.Context) ([]models.Customer, error) {
	match := bson.M{"trash": false}
	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"_id": -1}},
	}

	var result []models.Customer
	err := c.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
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

		userId, err := c.AddUpdateCustomerCredentials(ctx, customer)
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

		_, err := c.AddUpdateCustomerCredentials(ctx, customer)
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

func (c *customerSvcs) AddUpdateCustomerCredentials(ctx context.Context, data *models.Customer) (userId primitive.ObjectID, err error) {
	zeroId := primitive.NilObjectID

	var path, method string
	if data.Id == primitive.NilObjectID {
		path = "users/"
		method = "POST"
	} else {
		path = "users/" + data.Id.Hex()
		method = "PATCH"
	}

	password := data.Security.NewPassword
	if method == "POST" && password == "" {
		password = util.GeneratePassword(8, 2, 2, 2)
	}

	user := map[string]any{
		"firstName": data.Name,
		"lastName":  "-",
		"email":     data.Security.Email,
	}

	if password != "" {
		user["password"] = password
	}

	resp, err := c.gateway.Request(ctx, "users", path, method, "", user)

	if err != nil {
		return zeroId, errors.New("error adding user")
	} else if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 409:
			return zeroId, ErrDupliateEmail
		case 401:
			return zeroId, ErrUnauthorized
		case 403:
			return zeroId, ErrForbidden
		default:
			return zeroId, errors.New("error adding user")
		}

	}

	type TempUser struct {
		Id primitive.ObjectID `json:"_id"`
	}

	type AddedUser struct {
		User TempUser `json:"user"`
	}

	var _data AddedUser
	if errDec := json.NewDecoder(resp.Body).Decode(&_data); errDec != nil {
		return zeroId, errDec
	}

	return _data.User.Id, nil
}

func (c *customerSvcs) RegisterCustomerUser(ctx context.Context, data *models.Customer) (userId primitive.ObjectID, err error) {
	zeroId := primitive.NilObjectID

	user := map[string]any{
		"firstName":            data.Name,
		"lastName":             "-",
		"email":                data.Security.Email,
		"password":             data.Security.NewPassword,
		"passwordConfirmation": data.Security.NewPassword,
	}

	resp, err := c.gateway.Request(ctx, "users", "users/register", "POST", "", user)

	if err != nil {
		return zeroId, errors.New("error adding user")
	} else if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 409:
			return zeroId, ErrDupliateEmail
		case 401:
			return zeroId, ErrUnauthorized
		case 403:
			return zeroId, ErrForbidden
		default:
			return zeroId, errors.New("error adding user")
		}

	}

	type TempUser struct {
		Id primitive.ObjectID `json:"_id"` //user id
	}

	var _data TempUser
	if errDec := json.NewDecoder(resp.Body).Decode(&_data); errDec != nil {
		return zeroId, errDec
	}

	return _data.Id, nil
}

func (c *customerSvcs) RegisterAsCustomer(ctx context.Context, data *models.CustomerRegisterData) (*models.Customer, error) {
	result, err := c.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		customer := &models.Customer{
			CustomerDto: models.CustomerDto{
				Name:        data.ClientName,
				Nationality: data.Nationality,
				ClientContact: models.MemberContact{
					Mobile: data.ClientPhone,
					Email:  data.ClientEmail,
				},
				Security: models.MemberSecurity{
					Email:       data.ClientEmail,
					NewPassword: data.Password,
				},
			},
			CreatedAt: time.Now(),
		}

		if data.Password == "" {
			return nil, errors.New("password is required")
		}

		userId, err := c.RegisterCustomerUser(ctx, customer)
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
