package ourservice

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/helpers"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/member"
	membermodels "larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"larsa-tourism-microservices/pkg/types"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TravelRequestSvcs interface {
	Get(ctx context.Context, skip, limit int64, query any) (*models.TravelRequestPagination, error)
	GetCustomerRequests(ctx context.Context, customerId string, skip, limit int64, query any) (*models.CustomerTravelRequestPagination, error)
	GetAll(ctx context.Context, query any) ([]models.TravelRequestRes, error)
	GetOne(ctx context.Context, id string) (*models.TravelRequest, error)
	Add(ctx context.Context, data *models.TravelRequestDto) (*models.TravelRequest, error)
	Update(ctx context.Context, id string, data *models.TravelRequestDto) (*models.TravelRequest, error)
	MyRequests(ctx context.Context, status string) ([]models.TravelRequest, error)
	//UpdateStatus(ctx context.Context, id string, data *models.ChangeStatusDto) (*models.TravelRequest, error)
	Patch(ctx context.Context, filter, update bson.M) (*models.TravelRequest, error)
	BulkWrite(ctx context.Context, writes []mongo.WriteModel) (*mongo.BulkWriteResult, error)
	Approve(ctx context.Context, id string) (*models.TravelRequest, error)
	Reject(ctx context.Context, id string, data *models.RejectMyReq) (*models.TravelRequest, error)
	SetAsCompleted(ctx context.Context, id string) (*models.TravelRequest, error)
	GetAgentTransactions(ctx context.Context, agentId string, skip, limit int64, query any) (*models.AgentTransactionPagination, error)
}

type travelrequestsvcs struct {
	repo         repo.TravelRequestRepo
	invoicesvcs  InvoiceSvcs
	sortingsvcs  dbsvcs.SortingSvcs
	agentsvcs    member.AgentSvcs
	settingssvcs SettingsSvcs
	withtxn      *db.WithTxn
}

func NewTravelRequestSvcs(i *do.Injector) (TravelRequestSvcs, error) {
	return &travelrequestsvcs{
		repo:         do.MustInvoke[repo.TravelRequestRepo](i),
		invoicesvcs:  do.MustInvoke[InvoiceSvcs](i),
		sortingsvcs:  do.MustInvoke[dbsvcs.SortingSvcs](i),
		agentsvcs:    do.MustInvoke[member.AgentSvcs](i),
		settingssvcs: do.MustInvoke[SettingsSvcs](i),
		withtxn:      do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (t *travelrequestsvcs) GetOne(ctx context.Context, id string) (*models.TravelRequest, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return t.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (t *travelrequestsvcs) Get(ctx context.Context, skip, limit int64, query any) (*models.TravelRequestPagination, error) {
	match := bson.M{}

	f, err := helpers.ParseFilters[filter.TravelReqFilters](query)

	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := f.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := t.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "tourismPrograms",
		"localField":   "program",
		"foreignField": "_id",
		"as":           "programData",
	}})
	pipeline = append(pipeline, bson.M{"$unwind": bson.M{
		"path":                       "$programData",
		"preserveNullAndEmptyArrays": true,
	}})

	var result []models.TravelRequestRes
	errAg := t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))

	return &models.TravelRequestPagination{
		Requests: result,
		Pagination: types.Pagination{
			TotalPages: totalPages,
			PerPage:    limit,
			TotalCount: count,
		},
	}, nil
}

func (t *travelrequestsvcs) GetAll(ctx context.Context, query any) ([]models.TravelRequestRes, error) {
	match := bson.M{}

	f, err := helpers.ParseFilters[filter.TravelReqFilters](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := f.BuildPipeline(match)
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "tourismPrograms",
		"localField":   "program",
		"foreignField": "_id",
		"as":           "programData",
	}})
	pipeline = append(pipeline, bson.M{"$unwind": bson.M{
		"path":                       "$programData",
		"preserveNullAndEmptyArrays": true,
	}})

	var result []models.TravelRequestRes
	errAg := t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

func (t *travelrequestsvcs) GetCustomerRequests(ctx context.Context, customerId string, skip, limit int64, query any) (*models.CustomerTravelRequestPagination, error) {
	_id, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, err
	}

	f, err := helpers.ParseFilters[filter.TravelReqFilters](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	f.CustomerId = &_id

	result, err := t.Get(ctx, skip, limit, f)
	if err != nil {
		return nil, err
	}

	var r []models.CustomerTravelRequest
	for _, req := range result.Requests {
		var price float64
		var err error
		if req.ProgramData.Id != primitive.NilObjectID {
			price, err = req.ProgramData.CustomType.GetTotalPrice()
			if err != nil {
				return nil, err
			}

		}

		r = append(r, models.CustomerTravelRequest{TravelRequestRes: req, Price: price})
	}

	return &models.CustomerTravelRequestPagination{
		Requests:   r,
		Pagination: result.Pagination,
	}, nil
}

func (t *travelrequestsvcs) Add(ctx context.Context, data *models.TravelRequestDto) (*models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	result, err := t.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		seq, err := t.sortingsvcs.GetAndUpdateSourceSeq(ctx, "travelRequest")
		if err != nil {
			return nil, err
		}

		reqId := fmt.Sprintf("RQ-%d-%d", time.Now().Year(), seq)

		request := &models.TravelRequest{
			Id:               primitive.NewObjectID(),
			ReqId:            reqId,
			TravelRequestDto: *data,
			Date:             time.Now(),
			Status:           enums.TravelReqStatusPending,
			CustomerId:       cfg.User.Id,
			CreatedAt:        time.Now(),
			CreatedBy:        cfg.User.Id,
		}

		departureAgentId, err := t.getTravelRequestDepartureAgent(ctx, request)
		if err != nil {
			return nil, err
		}

		request.DepartureAgent = *departureAgentId

		if err := t.repo.Add(ctx, request); err != nil {
			return nil, err
		}

		return request, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.TravelRequest), nil
}

func (t *travelrequestsvcs) getTravelRequestDepartureAgent(ctx context.Context, travelReq *models.TravelRequest) (*primitive.ObjectID, error) {
	departureDestinationId, err := travelReq.GetDepartureDestinationId()
	if err != nil {
		return nil, err
	}

	agent, err := t.agentsvcs.GetAgentByDestination(ctx, departureDestinationId.Hex())
	if err != nil {
		return nil, errors.New("error get departure destination agent, check if agent has destination and is active")
	}

	return &agent.Id, nil

}

func (t *travelrequestsvcs) Update(ctx context.Context, id string, data *models.TravelRequestDto) (*models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	result, err := t.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		travelReq, err := t.GetOne(ctx, id)
		if err != nil {
			return nil, errors.New("travel request not found")
		}

		travelReq.TravelRequestDto = *data
		travelReq.UpdatedAt = time.Now()
		travelReq.UpdatedBy = cfg.User.Id

		departureAgentId, err := t.getTravelRequestDepartureAgent(ctx, travelReq)
		if err != nil {
			return nil, err
		}

		travelReq.DepartureAgent = *departureAgentId

		filter := bson.M{"_id": travelReq.Id, "status": enums.TravelReqStatusPending}
		update := bson.M{"$set": travelReq}

		updatedRequest, err := t.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, errors.New("error update travel request, check if travel request is pending")
		}

		return updatedRequest, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.TravelRequest), nil

}

func (t *travelrequestsvcs) MyRequests(ctx context.Context, status string) ([]models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	match := bson.M{"customerId": cfg.User.Id}

	if status != "all" {
		match["status"] = status
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$sort": bson.M{"_id": -1}},
	}

	var requests []models.TravelRequest
	err = t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &requests)
	})

	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (t *travelrequestsvcs) Approve(ctx context.Context, id string) (*models.TravelRequest, error) {
	result, err := t.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, err
		}
		_id, err := primitive.ObjectIDFromHex(id) // travel request id
		if err != nil {
			return nil, err
		}

		pipeline := []bson.M{
			{"$match": bson.M{"_id": _id}},
			{"$lookup": bson.M{
				"from":         "tourismPrograms",
				"localField":   "program",
				"foreignField": "_id",
				"as":           "program",
			}},
			{"$unwind": bson.M{
				"path":                       "$program",
				"preserveNullAndEmptyArrays": true,
			}},
			{"$lookup": bson.M{
				"from":         "tourismCustomers",
				"localField":   "customerId",
				"foreignField": "_id",
				"as":           "customer",
			}},
			{"$unwind": bson.M{
				"path":                       "$customer",
				"preserveNullAndEmptyArrays": true,
			}},
			{"$limit": 1},
		}

		var result []models.TravelRequestWithProgram
		errAg := t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
			return cur.All(ctx, &result)
		})

		if errAg != nil || len(result) == 0 {
			return nil, errors.New("error fetch travel request")
		}

		request := result[0]

		if request.CustomerId != cfg.User.Id {
			return nil, errors.New("only owner of this travel request can approve it")
		}

		program := request.Program
		customer := request.Customer

		if customer.Id == primitive.NilObjectID {
			return nil, errors.New("customer not found")
		}
		if program.Id == primitive.NilObjectID {
			return nil, errors.New("program not found")
		}

		if program.ProgramType != "custom" {
			return nil, errors.New("error program type")
		}

		invoiceDto := &models.InvoiceDto{
			DateOfIssue: time.Now(),
			Customer: models.InvoiceContact{
				Name:    customer.Name,
				Address: "",
				Phone:   customer.ClientContact.Mobile,
				Email:   customer.Security.Email,
				Website: customer.ClientContact.Website,
			},
			Company:     models.InvoiceContact{},
			ProgramName: program.Title,
			TravelStart: program.StartDate,
			TravelEnd:   program.EndDate,
			Services:    []models.InvoiceService{},
			Adjustments: []models.InvoiceAdjustment{},
			Note:        "",
		}

		svcss, err := program.CustomType.GetAllServicePricing()
		if err != nil {
			return nil, errors.New("error get program service list pricing")
		}

		invoiceDto.Services = svcss

		invoice, err := t.invoicesvcs.AddInvoiceForTravelRequest(ctx, request.Id, invoiceDto)
		if err != nil {
			return nil, errors.New("error add invoice")
		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": bson.M{
			"status":    enums.TravelReqStatusApproved,
			"invoiceId": invoice.Id,
			"updatedAt": time.Now(),
			"updatedBy": cfg.User.Id,
		}}

		updatedRequest, err := t.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedRequest, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.TravelRequest), nil

}

func (t *travelrequestsvcs) Reject(ctx context.Context, id string, data *models.RejectMyReq) (*models.TravelRequest, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	_id, err := primitive.ObjectIDFromHex(id) // travel request id
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"status":       enums.TravelReqStatusRejected,
		"rejectReason": data.Reason,
		"updatedAt":    time.Now(),
		"updatedBy":    cfg.User.Id,
	}}

	updatedRequest, err := t.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedRequest, nil

}

func (t *travelrequestsvcs) Patch(ctx context.Context, filter, update bson.M) (*models.TravelRequest, error) {
	return t.repo.Patch(ctx, filter, update)
}

func (t *travelrequestsvcs) BulkWrite(ctx context.Context, writes []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	return t.repo.BulkWrite(ctx, writes)
}

func (t *travelrequestsvcs) SetAsCompleted(ctx context.Context, id string) (*models.TravelRequest, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"status": enums.TravelReqStatusCompleted}}

	updatedTravelReq, err := t.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedTravelReq, nil

}

func (t *travelrequestsvcs) GetAgentTransactions(ctx context.Context, agentId string, skip, limit int64, query any) (*models.AgentTransactionPagination, error) {
	_id, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	match := bson.M{
		"status": enums.TravelReqStatusCompleted,
		"trash":  false,
		"$or": []bson.M{
			{
				"departureAgent": _id,
			},
			{
				"tripCoordinator": _id,
			},
		},
	}

	pipeline := []bson.M{
		{"$match": match},
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := t.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	invoiceLookup := []bson.M{
		{"$lookup": bson.M{
			"from":         "tourismInvoices",
			"localField":   "invoiceId",
			"foreignField": "_id",
			"as":           "invoice",
		}},
		{"$unwind": bson.M{
			"path":                       "$invoice",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	customerLookup := []bson.M{
		{"$lookup": bson.M{
			"from":         "tourismCustomers",
			"localField":   "customerId",
			"foreignField": "_id",
			"as":           "customer",
		}},
		{"$unwind": bson.M{
			"path":                       "$customer",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	departureAgentLookup := []bson.M{
		{"$lookup": bson.M{
			"from":         "tourismAgents",
			"localField":   "departureAgent",
			"foreignField": "_id",
			"as":           "departureAgentData",
		}},
		{"$unwind": bson.M{
			"path":                       "$departureAgentData",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	tripCoordinatorAgentLookup := []bson.M{
		{"$lookup": bson.M{
			"from":         "tourismAgents",
			"localField":   "tripCoordinator",
			"foreignField": "_id",
			"as":           "tripCoordinatorData",
		}},
		{"$unwind": bson.M{
			"path":                       "$tripCoordinatorData",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	pipeline = append(pipeline, invoiceLookup...)
	pipeline = append(pipeline, customerLookup...)
	pipeline = append(pipeline, departureAgentLookup...)
	pipeline = append(pipeline, tripCoordinatorAgentLookup...)

	type aux struct {
		models.TravelRequest `bson:",inline"`
		Invoice              models.Invoice        `bson:"invoice" json:"invoice"`
		Customer             membermodels.Customer `bson:"customer" json:"customer"`
		DepartureAgentData   membermodels.Agent    `bson:"departureAgentData" json:"departureAgentData"`
		TripCoordinatorData  membermodels.Agent    `bson:"tripCoordinatorData" json:"tripCoordinatorData"`
	}

	var result []aux
	err = t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	if err != nil {
		return nil, err
	}

	settings, err := t.settingssvcs.Get(ctx)
	if err != nil {
		return nil, errors.New("error get settings")
	}

	profitRatio := settings.ProfitRatio //platform profit ratio

	var transactions []models.AgentTransaction
	for _, r := range result {
		if r.Invoice.Id == primitive.NilObjectID {
			return nil, errors.New("invoice not found")
		}

		clientProfit := r.Invoice.Total * (profitRatio / 100)

		var commission float64

		if r.DepartureAgent == r.TripCoordinator {
			profitOfTourismProgram := r.DepartureAgentData.Financial.ProfitOfTourismProgram
			if profitOfTourismProgram {
				commission = clientProfit * (r.DepartureAgentData.Financial.Ratio / 100)

			}

		} else {
			if r.DepartureAgent == _id {
				profitOfTourismProgram := r.DepartureAgentData.Financial.ProfitOfTourismProgram
				if profitOfTourismProgram {
					commission = clientProfit * (r.DepartureAgentData.Financial.Ratio / 100)
				}

			} else if r.TripCoordinator == _id {
				profitOfTourismProgram := r.TripCoordinatorData.Financial.ProfitOfTourismProgram
				if profitOfTourismProgram {
					commission = clientProfit * (r.TripCoordinatorData.Financial.Ratio / 100)
				}
			}
		}

		transaction := models.AgentTransaction{
			TravelRequestId: r.Id,
			InvoiceId:       r.InvoiceId,
			Date:            r.Date,
			OrderId:         r.Invoice.InvoiceId,
			CustomerName:    r.Customer.Name,
			Commission:      commission,
		}

		transactions = append(transactions, transaction)
	}

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))

	return &models.AgentTransactionPagination{
		Transactions: transactions,
		Pagination: types.Pagination{
			TotalPages: totalPages,
			PerPage:    limit,
			TotalCount: count,
		},
	}, nil

}

func (t *travelrequestsvcs) Count(ctx context.Context, filter any) (int64, error) {
	return t.repo.Count(ctx, filter)
}
