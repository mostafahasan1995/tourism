package ourservice

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/member"
	membermodels "larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"reflect"
	"time"

	"larsa-tourism-microservices/pkg/types"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var programTotalCostStages = []bson.M{
	{
		"$addFields": bson.M{
			"resolvedDestinations": bson.M{
				"$concatArrays": bson.A{
					bson.M{"$ifNull": bson.A{"$customType.destinations", bson.A{}}},
					bson.M{"$ifNull": bson.A{"$customType.vipCar.destinations", bson.A{}}},
					bson.M{"$ifNull": bson.A{"$customType.flightTicketRequest.destinations", bson.A{}}},
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"destinationCosts": bson.M{
				"$map": bson.M{
					"input": "$resolvedDestinations",
					"as":    "d",
					"in": bson.M{
						"$let": bson.M{
							"vars": bson.M{
								"accomSum": bson.M{
									"$cond": bson.M{
										"if": bson.M{"$isArray": "$$d.accommodation"},
										"then": bson.M{
											"$reduce": bson.M{
												"input":        "$$d.accommodation",
												"initialValue": 0,
												"in": bson.M{
													"$add": bson.A{
														"$$value",
														bson.M{"$ifNull": bson.A{"$$this.totalStayCost", 0}},
													},
												},
											},
										},
										"else": 0,
									},
								},
							},
							"in": bson.M{
								"$add": bson.A{
									bson.M{"$ifNull": bson.A{"$$d.totalCost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.flightTickets.totalCost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.transportation.totalCost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.activities.totalCost", 0}},
									"$$accomSum",
									bson.M{"$ifNull": bson.A{"$$d.services.onGroundAssistance.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.travelInsurance.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.visaAssistance.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.welcomeKit.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.freeSimCardWifi.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.complimentaryGifts.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.vipAirportServices.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.personalTravelConsultant.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.childcareServices.cost", 0}},
									bson.M{"$ifNull": bson.A{"$$d.services.accessibilitySupport.cost", 0}},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"totalDestinationsCost": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$isArray": "$destinationCosts"},
					"then": bson.M{"$sum": "$destinationCosts"},
					"else": 0,
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"topLevelTotalsSum": bson.M{
				"$add": bson.A{
					bson.M{"$ifNull": bson.A{"$customType.vipCar.totalCost", 0}},
					bson.M{"$ifNull": bson.A{"$customType.flightTicketRequest.totalCost", 0}},
					bson.M{"$ifNull": bson.A{"$customType.hotelBooking.totalCost", 0}},
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"customType.totalCost": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$gt": bson.A{bson.M{"$size": "$resolvedDestinations"}, 0}},
					"then": "$totalDestinationsCost",
					"else": "$topLevelTotalsSum",
				},
			},
		},
	},
	{
		"$addFields": bson.M{
			"totalCost": bson.M{"$ifNull": bson.A{"$customType.totalCost", 0}},
		},
	},
	{
		"$project": bson.M{
			"resolvedDestinations":  0,
			"destinationCosts":      0,
			"totalDestinationsCost": 0,
			"topLevelTotalsSum":     0,
		},
	},
}

var programLookup = func() []bson.M {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"$expr": bson.M{
					"$eq": bson.A{"$_id", "$$programId"},
				},
			},
		},
	}

	for _, stage := range programTotalCostStages {
		pipeline = append(pipeline, stage)
	}

	return []bson.M{
		{
			"$lookup": bson.M{
				"from":     "tourismPrograms",
				"let":      bson.M{"programId": "$program"},
				"pipeline": pipeline,
				"as":       "programData",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$programData",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
}()

var travelReqCustomerLookup = []bson.M{
	{"$lookup": bson.M{
		"from":         "tourismCustomers",
		"localField":   "customerId",
		"foreignField": "_id",
		"as":           "customerData",
	}},
	{"$unwind": bson.M{
		"path":                       "$customerData",
		"preserveNullAndEmptyArrays": true,
	}},
}

var travelReqPackageLookup = []bson.M{
	{"$lookup": bson.M{
		"from":         "tourismPackages",
		"localField":   "package",
		"foreignField": "_id",
		"as":           "packageData",
	}},
	{"$unwind": bson.M{
		"path":                       "$packageData",
		"preserveNullAndEmptyArrays": true,
	}},
}

var hotelLookup = []bson.M{
	{
		"$lookup": bson.M{
			"from": "tourismHotels",
			"let":  bson.M{"hotelId": "$hotelBooking.hotelId"},
			"pipeline": []bson.M{
				{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$ne": []interface{}{"$$hotelId", nil}},
								{"$eq": []interface{}{"$_id", "$$hotelId"}},
							},
						},
					},
				},
			},
			"as": "hotelData",
		},
	},
	{
		"$unwind": bson.M{
			"path":                       "$hotelData",
			"preserveNullAndEmptyArrays": true,
		},
	},
	{
		"$set": bson.M{"hotelOwner": "$hotelData.owner"},
	},
	{
		"$project": bson.M{
			"hotelData": 0,
		},
	},
}

type TravelRequestSvcs interface {
	Get(ctx context.Context, skip, limit int64, query any) (*models.TravelRequestPagination, error)
	GetCustomerRequests(ctx context.Context, customerId string, skip, limit int64, query any) (*models.CustomerTravelRequestPagination, error)
	GetAll(ctx context.Context, query any) ([]models.TravelRequestRes, error)
	GetOne(ctx context.Context, id string) (*models.TravelRequest, error)
	GetCompanyTransactions(ctx context.Context, skip, limit int64, query any) (*models.CompanyTransactionPagination, error)
	Add(ctx context.Context, data *models.TravelRequestDto) (*models.TravelRequest, error)
	Update(ctx context.Context, id string, data *models.TravelRequestDto) (*models.TravelRequest, error)
	MyRequests(ctx context.Context, status string) ([]models.TravelRequestRes, error)
	Patch(ctx context.Context, filter, update bson.M) (*models.TravelRequest, error)
	BulkWrite(ctx context.Context, writes []mongo.WriteModel) (*mongo.BulkWriteResult, error)
	Approve(ctx context.Context, id string) (*models.TravelRequest, error)
	Reject(ctx context.Context, id string, data *models.RejectMyReq) (*models.TravelRequest, error)
	SetAsCompleted(ctx context.Context, id string) (*models.TravelRequest, error)
	GetAgentTransactions(ctx context.Context, agentId string, skip, limit int64, query any) (*models.AgentTransactionPagination, error)

	Count(ctx context.Context, filter any) (int64, error)
	//v2
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.TravelRequestPagination, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.TravelRequestRes, error)
}

type travelrequestsvcs struct {
	repo                  repo.TravelRequestRepo
	invoicesvcs           InvoiceSvcs
	sortingsvcs           dbsvcs.SortingSvcs
	agentsvcs             member.AgentSvcs
	customersvcs          member.CustomerSvcs
	financialsettingssvcs FinancialSettingsSvcs
	agentfinancialsvcs    AgentFinancialSvcs
	programrepo           repo.ProgramRepo
	withtxn               *db.WithTxn
}

func NewTravelRequestSvcs(i *do.Injector) (TravelRequestSvcs, error) {
	return &travelrequestsvcs{
		repo:                  do.MustInvoke[repo.TravelRequestRepo](i),
		invoicesvcs:           do.MustInvoke[InvoiceSvcs](i),
		sortingsvcs:           do.MustInvoke[dbsvcs.SortingSvcs](i),
		agentsvcs:             do.MustInvoke[member.AgentSvcs](i),
		customersvcs:          do.MustInvoke[member.CustomerSvcs](i),
		financialsettingssvcs: do.MustInvoke[FinancialSettingsSvcs](i),
		agentfinancialsvcs:    do.MustInvoke[AgentFinancialSvcs](i),
		programrepo:           do.MustInvoke[repo.ProgramRepo](i),
		withtxn:               do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (t *travelrequestsvcs) GetOne(ctx context.Context, id string) (*models.TravelRequest, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return t.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (t *travelrequestsvcs) buildUserPipeline(ctx context.Context, query any) ([]bson.M, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	userId := cfg.User.Id
	//check if user can get other travel requests
	check, ok := ctx.Value(util.ReqCapabilityCheck).(*types.CapabilityCheck)
	if !ok {
		return nil, errors.New("error check user capability")
	}

	var pipeline []bson.M

	if check.Capability == "tourismGetOtherTravelRequests" && check.IsAllowed {
		match := bson.M{"trash": false}

		f, err := helpers.ParseFilters[filter.TravelReqFilters](query)

		if err != nil {
			return nil, errors.New("invalid query")
		}

		pipeline = f.BuildPipeline(match)

	} else {

		pipeline = []bson.M{
			{"$match": bson.M{
				"trash": false,
			}},
		}

		pipeline = append(pipeline, hotelLookup...)
		pipeline = append(pipeline, bson.M{"$match": bson.M{
			"$or": bson.A{
				bson.M{"departureAgent": userId},
				bson.M{"tripCoordinator": userId},
				bson.M{"hotelOwner": userId},
			},
		}})

		f, err := helpers.ParseFilters[filter.TravelReqFilters](query)
		if err != nil {
			return nil, errors.New("invalid query")
		}
		filterPipeline := f.BuildPipeline(bson.M{})
		pipeline = append(pipeline, filterPipeline...)

	}

	return pipeline, nil

}

func (t *travelrequestsvcs) Get(ctx context.Context, skip, limit int64, query any) (*models.TravelRequestPagination, error) {
	pipeline, err := t.buildUserPipeline(ctx, query)
	if err != nil {
		return nil, err
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := t.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	pipeline = append(pipeline, programLookup...)
	pipeline = append(pipeline, travelReqCustomerLookup...)
	pipeline = append(pipeline, travelReqPackageLookup...)

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
	pipeline, err := t.buildUserPipeline(ctx, query)
	if err != nil {
		return nil, err
	}
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})

	pipeline = append(pipeline, programLookup...)
	pipeline = append(pipeline, travelReqCustomerLookup...)
	pipeline = append(pipeline, travelReqPackageLookup...)

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
		userId := cfg.User.Id // same as customerId
		// check if the user making the request is customer
		customer, err := t.customersvcs.GetOne(ctx, userId.Hex())
		if err != nil {
			return nil, errors.New("error get customer, check if user is customer")
		}

		if data.ClientEmail != customer.Security.Email {
			return nil, errors.New("client email is not the same as customer email")
		}

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
			CustomerId:       customer.Id,
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

func (t *travelrequestsvcs) MyRequests(ctx context.Context, status string) ([]models.TravelRequestRes, error) {
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

	pipeline = append(pipeline, programLookup...)
	pipeline = append(pipeline, travelReqCustomerLookup...)
	pipeline = append(pipeline, travelReqPackageLookup...)

	var requests []models.TravelRequestRes
	err = t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &requests)
	})

	if err != nil {
		return nil, err
	}

	// Calculate fees for each request
	financialSettings, err := t.financialsettingssvcs.Get(ctx)
	if err != nil {
		return nil, errors.New("error get settings")
	}
	profitRatio := financialSettings.ProfitRatio

	for i := range requests {
		if requests[i].ProgramData.CustomType != nil && requests[i].ProgramData.CustomType.TotalCost > 0 {
			subTotal := requests[i].ProgramData.CustomType.TotalCost
			fees := subTotal * (profitRatio / 100)
			total := subTotal + fees

			requests[i].SubTotal = subTotal
			requests[i].Fees = fees
			requests[i].Total = total
		}
	}

	return requests, nil
}

func (t *travelrequestsvcs) Approve(ctx context.Context, id string) (*models.TravelRequest, error) {
	// Get profit ratio before transaction to avoid context issues
	financialSettings, err := t.financialsettingssvcs.Get(ctx)
	if err != nil {
		return nil, errors.New("error get financial settings")
	}
	profitRatio := financialSettings.ProfitRatio

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

		//Shouldn't be able to approve multiple times
		if request.Status == enums.TravelReqStatusApproved {
			return nil, errors.New("travel request already approved")
		}

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
		// Update program status to "approved" if it's "waiting"
		if program.ProgramDto.Status == "waiting" {
			filter := bson.M{"_id": program.Id}
			update := bson.M{"$set": bson.M{
				"status":    "approved",
				"updatedAt": time.Now(),
				"updatedBy": cfg.User.Id,
			}}
			_, err := t.programrepo.Patch(ctx, filter, update)
			if err != nil {
				return nil, fmt.Errorf("error updating program status: %w", err)
			}
		}
		svcss, err := program.CustomType.GetAllServicePricing()

		if err != nil {
			return nil, errors.New("error get program service list pricing")
		}

		invoiceDto.Services = svcss

		invoice, err := t.invoicesvcs.AddInvoiceForTravelRequest(ctx, request.Id, invoiceDto, profitRatio)
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

	// Get travel request to find associated program
	travelReq, err := t.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update travel request status
	filter := bson.M{"_id": _id}
	rejectReason := ""
	if data != nil {
		rejectReason = data.Reason
	}
	update := bson.M{"$set": bson.M{
		"status":       enums.TravelReqStatusRejected,
		"rejectReason": rejectReason,
		"updatedAt":    time.Now(),
		"updatedBy":    cfg.User.Id,
	}}

	updatedRequest, err := t.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Update program status to "rejected" if program exists
	if travelReq.Program != primitive.NilObjectID {
		// Get program to check current status
		program, err := t.programrepo.GetByFilter(ctx, bson.M{"_id": travelReq.Program, "trash": false})
		if err == nil && program != nil {
			// Check if status is "waiting" or "pending" and update to "rejected"
			currentStatus := program.ProgramDto.Status
			if currentStatus == "waiting" || currentStatus == "pending" {
				filter := bson.M{"_id": travelReq.Program}
				update := bson.M{"$set": bson.M{
					"status":    "rejected",
					"updatedAt": time.Now(),
					"updatedBy": cfg.User.Id,
				}}
				_, err := t.programrepo.Patch(ctx, filter, update)
				if err != nil {
					// Log error but don't fail the reject operation
					fmt.Printf("error updating program status: %v\n", err)
				}
			}
		}
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
	_id, err := primitive.ObjectIDFromHex(id) // travel request id
	if err != nil {
		return nil, err
	}

	travelReq, err := t.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	if travelReq.Status == enums.TravelReqStatusCompleted {
		return travelReq, nil
	}

	if travelReq.Program == primitive.NilObjectID {
		return nil, errors.New("travel request has no program")
	}

	filter := bson.M{"_id": _id, "program": bson.M{"$ne": primitive.NilObjectID}}
	update := bson.M{"$set": bson.M{"status": enums.TravelReqStatusCompleted}}

	updatedTravelReq, err := t.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, errors.New("error updated status , check if travel request has program")
	}

	if err := t.distributeAgentProfits(ctx, updatedTravelReq); err != nil {
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

	financialSettings, err := t.financialsettingssvcs.Get(ctx)
	if err != nil {
		return nil, errors.New("error get settings")
	}

	profitRatio := financialSettings.ProfitRatio // platform profit ratio
	var transactions []models.AgentTransaction

	loadAgentFinancial := func(agentData membermodels.Agent, agentID primitive.ObjectID) (membermodels.AgentFinancial, error) {
		financial := agentData.Financial
		if !financial.ProfitOfTourismProgram || financial.Ratio == 0 {
			agentDoc, err := t.agentsvcs.GetByFilter(ctx, bson.M{"_id": agentID, "trash": false})
			if err != nil {
				return financial, err
			}
			financial = agentDoc.Financial
		}
		return financial, nil
	}

	for _, r := range result {
		if r.Invoice.Id == primitive.NilObjectID {
			return nil, errors.New("invoice not found")
		}

		// Use SubTotal for profit calculation (before fees)
		// If SubTotal is 0, fallback to Total (for older invoices)
		invoiceSubTotal := r.Invoice.SubTotal
		if invoiceSubTotal == 0 {
			invoiceSubTotal = r.Invoice.Total
		}

		// Calculate client profit: invoiceSubTotal * (profitRatio / 100)

		clientProfit := invoiceSubTotal * (profitRatio / 100)

		// Initialize commission
		var commission float64

		// Determine which agent to calculate commission for
		var targetFinancial membermodels.AgentFinancial
		var hasValidFinancial bool

		if r.DepartureAgent == r.TripCoordinator {
			// Same agent for both roles
			if r.DepartureAgent == _id {
				financial, err := loadAgentFinancial(r.DepartureAgentData, r.DepartureAgent)
				if err == nil && financial.ProfitOfTourismProgram && financial.Ratio > 0 {
					targetFinancial = financial
					hasValidFinancial = true
				}
			}
		} else {
			// Different agents
			if r.DepartureAgent == _id {
				financial, err := loadAgentFinancial(r.DepartureAgentData, r.DepartureAgent)
				if err == nil && financial.ProfitOfTourismProgram && financial.Ratio > 0 {
					targetFinancial = financial
					hasValidFinancial = true
				}
			} else if r.TripCoordinator == _id {
				financial, err := loadAgentFinancial(r.TripCoordinatorData, r.TripCoordinator)
				if err == nil && financial.ProfitOfTourismProgram && financial.Ratio > 0 {
					targetFinancial = financial
					hasValidFinancial = true
				}
			}
		}

		// Calculate commission: clientProfit * (agentRatio / 100)
		// Example: 82.5 * (5 / 100) = 4.125
		if hasValidFinancial {
			commission = clientProfit * (targetFinancial.Ratio / 100)
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

func (t *travelrequestsvcs) GetCompanyTransactions(ctx context.Context, skip, limit int64, query any) (*models.CompanyTransactionPagination, error) {
	match := bson.M{
		"status": enums.TravelReqStatusCompleted,
		"trash":  false,
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
		{
			"$lookup": bson.M{
				"from":         "tourismInvoices",
				"localField":   "invoiceId",
				"foreignField": "_id",
				"as":           "invoice",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$invoice",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	customerLookup := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "tourismCustomers",
				"localField":   "customerId",
				"foreignField": "_id",
				"as":           "customer",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$customer",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	pipeline = append(pipeline,
		bson.M{"$sort": bson.M{"_id": -1}},
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
	)

	pipeline = append(pipeline, invoiceLookup...)
	pipeline = append(pipeline, customerLookup...)

	type companyAux struct {
		models.TravelRequest `bson:",inline"`
		Invoice              models.Invoice        `bson:"invoice" json:"invoice"`
		Customer             membermodels.Customer `bson:"customer" json:"customer"`
	}

	var result []companyAux
	err = t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	financialSettings, err := t.financialsettingssvcs.Get(ctx)
	if err != nil {
		return nil, errors.New("error get settings")
	}

	profitRatio := financialSettings.ProfitRatio

	var transactions []models.CompanyTransaction
	for _, r := range result {
		if r.Invoice.Id == primitive.NilObjectID {
			continue
		}

		invoiceTotal := r.Invoice.Total
		if invoiceTotal == 0 {
			invoiceTotal = r.Invoice.SubTotal
		}

		profit := invoiceTotal * (profitRatio / 100)

		transaction := models.CompanyTransaction{
			TravelRequestId: r.Id,
			InvoiceId:       r.InvoiceId,
			Date:            r.Date,
			OrderId:         r.Invoice.InvoiceId,
			CustomerName:    r.Customer.Name,
			Profit:          profit,
		}

		transactions = append(transactions, transaction)
	}

	var totalPages float64 = 0
	if limit > 0 {
		totalPages = math.Ceil(float64(count) / float64(limit))
	}

	return &models.CompanyTransactionPagination{
		Transactions: transactions,
		Pagination: types.Pagination{
			TotalPages: totalPages,
			PerPage:    limit,
			TotalCount: count,
		},
	}, nil
}

func (t *travelrequestsvcs) distributeAgentProfits(ctx context.Context, travelReq *models.TravelRequest) error {
	if travelReq == nil {
		return nil
	}

	if travelReq.InvoiceId == primitive.NilObjectID {
		return nil
	}

	invoice, err := t.invoicesvcs.GetOne(ctx, travelReq.InvoiceId.Hex())
	if err != nil {
		return err
	}

	invoiceTotal := invoice.Total
	if invoiceTotal == 0 {
		invoiceTotal = invoice.SubTotal
	}

	if invoiceTotal <= 0 {
		return nil
	}

	financialSettings, err := t.financialsettingssvcs.Get(ctx)
	if err != nil {
		return errors.New("error get settings")
	}

	clientProfit := invoiceTotal * (financialSettings.ProfitRatio / 100)
	if clientProfit <= 0 {
		return nil
	}

	type agentTarget struct {
		agent *membermodels.Agent
		ratio float64
	}

	targets := map[primitive.ObjectID]agentTarget{}

	if travelReq.DepartureAgent != primitive.NilObjectID {
		agent, err := t.agentsvcs.GetByFilter(ctx, bson.M{"_id": travelReq.DepartureAgent, "trash": false})
		if err == nil && agent.Financial.ProfitOfTourismProgram && agent.Financial.Ratio > 0 {
			targets[agent.Id] = agentTarget{agent: agent, ratio: agent.Financial.Ratio}
		}
	}

	if travelReq.TripCoordinator != primitive.NilObjectID {
		if _, exists := targets[travelReq.TripCoordinator]; !exists {
			agent, err := t.agentsvcs.GetByFilter(ctx, bson.M{"_id": travelReq.TripCoordinator, "trash": false})
			if err == nil && agent.Financial.ProfitOfTourismProgram && agent.Financial.Ratio > 0 {
				targets[agent.Id] = agentTarget{agent: agent, ratio: agent.Financial.Ratio}
			}
		}
	}

	for _, target := range targets {
		commission := clientProfit * (target.ratio / 100)
		if commission <= 0 {
			continue
		}

		if _, err := t.agentfinancialsvcs.AddProfit(ctx, target.agent, commission); err != nil {
			return err
		}
	}

	return nil
}

func (t *travelrequestsvcs) Count(ctx context.Context, filter any) (int64, error) {
	return t.repo.Count(ctx, filter)
}

// v2

func (t *travelrequestsvcs) buildUserPipelineV2(ctx context.Context, query *query.Conditions) ([]bson.M, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	userId := cfg.User.Id

	// Check if user has agent role by checking roleNames
	// Use reflection to access RoleNames field which might not be in the struct definition
	isAgent := false
	userDataValue := reflect.ValueOf(cfg.User.UserData)
	if userDataValue.Kind() == reflect.Struct {
		roleNamesField := userDataValue.FieldByName("RoleNames")
		if roleNamesField.IsValid() && roleNamesField.Kind() == reflect.Slice {
			roleNames := roleNamesField.Interface().([]string)
			for _, roleName := range roleNames {
				if roleName == "agent" {
					isAgent = true
					break
				}
			}
		}
	}

	var pipeline []bson.M

	if isAgent {
		// Agent: filter by program.agentId = userId
		pipeline = []bson.M{
			{"$match": bson.M{"trash": false}},
		}

		// Lookup program to check agentId
		programLookupForFilter := []bson.M{
			{
				"$lookup": bson.M{
					"from":         "tourismPrograms",
					"localField":   "program",
					"foreignField": "_id",
					"as":           "programForFilter",
				},
			},
			{
				"$unwind": bson.M{
					"path":                       "$programForFilter",
					"preserveNullAndEmptyArrays": false, // Only include travel requests with programs
				},
			},
			{
				"$match": bson.M{
					"programForFilter.agentId": userId,
				},
			},
			{
				"$project": bson.M{
					"programForFilter": 0,
				},
			},
		}

		pipeline = append(pipeline, programLookupForFilter...)

		if err := query.CheckValid(); err != nil {
			return nil, err
		}
		filter, err := query.ConvertToMongo()
		if err != nil {
			return nil, err
		}

		pipeline = append(pipeline, bson.M{"$match": filter})

	} else {
		// Administrator: return all travel requests (no agent filter)
		if err := query.CheckValid(); err != nil {
			return nil, err
		}
		filter, err := query.ConvertToMongo()
		if err != nil {
			return nil, err
		}

		pipeline = []bson.M{
			{"$match": bson.M{"trash": false}},
			{"$match": filter},
		}
	}

	return pipeline, nil

}

func (t *travelrequestsvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.TravelRequestPagination, error) {
	pipeline, err := t.buildUserPipelineV2(ctx, query)
	if err != nil {
		return nil, err
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := t.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	pipeline = append(pipeline, programLookup...)
	pipeline = append(pipeline, travelReqCustomerLookup...)
	pipeline = append(pipeline, travelReqPackageLookup...)

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

func (t *travelrequestsvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.TravelRequestRes, error) {
	pipeline, err := t.buildUserPipelineV2(ctx, query)
	if err != nil {
		return nil, err
	}
	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})

	pipeline = append(pipeline, programLookup...)
	pipeline = append(pipeline, travelReqCustomerLookup...)
	pipeline = append(pipeline, travelReqPackageLookup...)

	var result []models.TravelRequestRes
	errAg := t.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}
