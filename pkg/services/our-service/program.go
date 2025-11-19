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
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/services/picklist"
	pModels "larsa-tourism-microservices/pkg/services/picklist/models"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	interactionsModels "larsa-tourism-microservices/pkg/services/interactions/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var updatedByUserLookup = []bson.M{
	{"$lookup": bson.M{
		"from":         "users",
		"localField":   "updatedBy",
		"foreignField": "_id",
		"as":           "updatedByUser",
	}},
	{"$unwind": bson.M{
		"path":                       "$updatedByUser",
		"preserveNullAndEmptyArrays": true,
	}},
	{
		"$set": bson.M{
			"updatedByName": bson.M{
				"$concat": []interface{}{
					"$updatedByUser.firstName",
					" ",
					"$updatedByUser.lastName",
				},
			},
		},
	},
	{
		"$project": bson.M{
			"updatedByUser": 0,
		},
	},
}

var packageLookup = []bson.M{
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
	{
		"$set": bson.M{
			"packageName": "$packageData.name",
		},
	},
	{
		"$project": bson.M{
			"packageData": 0,
		},
	},
}

var customerLookup = []bson.M{
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
	{
		"$set": bson.M{
			"customerName": "$customer.name",
		},
	},
	{
		"$project": bson.M{
			"customer": 0,
		},
	},
}

var durationLookup = bson.M{
	"$addFields": bson.M{
		"duration": bson.M{
			"$cond": bson.M{
				"if": bson.M{
					"$and": []interface{}{
						bson.M{"$ne": []interface{}{"$startDate", nil}},
						bson.M{"$ne": []interface{}{"$endDate", nil}},
					},
				},
				"then": bson.M{
					"$dateDiff": bson.M{
						"startDate": "$startDate",
						"endDate":   "$endDate",
						"unit":      "day",
					},
				},
				"else": 0,
			},
		},
	},
}

type ProgramSvcs interface {
	GetOne(ctx context.Context, id string) (*models.ProgramRes, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error)
	GetAll(ctx context.Context, query string) ([]models.Program, error)
	GetAuth(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error)
	GetAllAuth(ctx context.Context, query string) ([]models.Program, error)
	Add(ctx context.Context, data *models.ProgramDto) (*models.Program, error)
	Update(ctx context.Context, id string, data *models.ProgramDto) (*models.Program, error)
	UpdateIsFav(ctx context.Context, id string, isFav bool) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, filter any) (int64, error)
	//v2
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.ProgramPagination, error)
}

type programsvcs struct {
	repo                  repo.ProgramRepo
	travelreqsvcs         TravelRequestSvcs
	travelrequestrepo     repo.TravelRequestRepo
	withtxn               *db.WithTxn
	activitiesSvcs        picklist.ActivitiesSvcs
	financialsettingssvcs FinancialSettingsSvcs
	customersvcs          member.CustomerSvcs
	sortingsvcs           dbsvcs.SortingSvcs
	agentsvcs             member.AgentSvcs
}

func NewProgramSvcs(i *do.Injector) (ProgramSvcs, error) {
	return &programsvcs{
		repo:                  do.MustInvoke[repo.ProgramRepo](i),
		travelreqsvcs:         do.MustInvoke[TravelRequestSvcs](i),
		travelrequestrepo:     do.MustInvoke[repo.TravelRequestRepo](i),
		withtxn:               do.MustInvoke[*db.WithTxn](i),
		activitiesSvcs:        do.MustInvoke[picklist.ActivitiesSvcs](i),
		financialsettingssvcs: do.MustInvoke[FinancialSettingsSvcs](i),
		customersvcs:          do.MustInvoke[member.CustomerSvcs](i),
		sortingsvcs:           do.MustInvoke[dbsvcs.SortingSvcs](i),
		agentsvcs:             do.MustInvoke[member.AgentSvcs](i),
	}, nil
}

//

func (p *programsvcs) GetOne(ctx context.Context, id string) (*models.ProgramRes, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Build aggregation pipeline for detailed program information
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":   _id,
				"trash": false,
			},
		},
	}
	// build totalCost aggregation stages
	totalCostStages := []bson.M{
		// 1) build a single array with all possible destination sources
		{
			"$addFields": bson.M{
				"resolvedDestinations": bson.M{
					"$concatArrays": []interface{}{
						bson.M{"$ifNull": []interface{}{"$customType.destinations", bson.A{}}},
						bson.M{"$ifNull": []interface{}{"$customType.vipCar.destinations", bson.A{}}},
						bson.M{"$ifNull": []interface{}{"$customType.flightTicketRequest.destinations", bson.A{}}},
					},
				},
			},
		},

		// 2) compute per-destination cost array (including accommodation reduce)
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
														"$add": []interface{}{
															"$$value",
															bson.M{"$ifNull": []interface{}{"$$this.totalStayCost", 0}},
														},
													},
												},
											},
											"else": 0,
										},
									},
								},
								"in": bson.M{
									"$add": []interface{}{
										// use precomputed destination totalCost if exists (vipCar often has this)
										bson.M{"$ifNull": []interface{}{"$$d.totalCost", 0}},
										// common per-destination parts (some may be absent)
										bson.M{"$ifNull": []interface{}{"$$d.flightTickets.totalCost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.transportation.totalCost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.activities.totalCost", 0}},
										"$$accomSum",
										// per-destination service costs (guard with ifNull)
										bson.M{"$ifNull": []interface{}{"$$d.services.onGroundAssistance.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.travelInsurance.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.visaAssistance.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.welcomeKit.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.freeSimCardWifi.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.complimentaryGifts.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.vipAirportServices.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.personalTravelConsultant.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.childcareServices.cost", 0}},
										bson.M{"$ifNull": []interface{}{"$$d.services.accessibilitySupport.cost", 0}},
									},
								},
							},
						},
					},
				},
			},
		},

		// 3) sum destinationCosts -> totalDestinationsCost
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

		// 4) compute fallback sum of top-level totals (used when resolvedDestinations is empty)
		{
			"$addFields": bson.M{
				"topLevelTotalsSum": bson.M{
					"$add": []interface{}{
						bson.M{"$ifNull": []interface{}{"$customType.vipCar.totalCost", 0}},
						bson.M{"$ifNull": []interface{}{"$customType.flightTicketRequest.totalCost", 0}},
						bson.M{"$ifNull": []interface{}{"$customType.hotelBooking.totalCost", 0}},
						// Add other top-level cost fields here if you have them
					},
				},
			},
		},

		// 5) set customType.totalCost: prefer destinations sum if any, else fallback to top-level totals
		{
			"$addFields": bson.M{
				"customType.totalCost": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$gt": []interface{}{bson.M{"$size": "$resolvedDestinations"}, 0}},
						"then": "$totalDestinationsCost",
						"else": "$topLevelTotalsSum",
					},
				},
			},
		},

		// 6) expose root-level totalCost and cleanups
		{
			"$addFields": bson.M{
				"totalCost": bson.M{"$ifNull": []interface{}{"$customType.totalCost", 0}},
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

	pipeline = append(pipeline, totalCostStages...)

	// Add lookups for comprehensive data
	pipeline = append(pipeline, customerLookup...)
	pipeline = append(pipeline, packageLookup...)
	pipeline = append(pipeline, updatedByUserLookup...)
	pipeline = append(pipeline, durationLookup)
	// add favorite pipeline
	pipeline = append(pipeline, interactionsModels.BuildFavoritePipelineWithAuth(ctx, interactionsModels.FaveTypeProgram)...)
	var result []models.ProgramRes
	errAg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	if len(result) == 0 {
		return nil, errors.New("program not found")
	}

	// Calculate fees using financial settings
	financialSettings, err := p.financialsettingssvcs.Get(ctx)
	if err != nil {
		return nil, errors.New("error get settings")
	}
	profitRatio := financialSettings.ProfitRatio

	// Calculate subTotal, fees, and total for the program
	if result[0].Program.CustomType != nil && result[0].Program.CustomType.TotalCost > 0 {
		subTotal := result[0].Program.CustomType.TotalCost
		fees := subTotal * (profitRatio / 100)
		total := subTotal + fees

		result[0].SubTotal = subTotal
		result[0].Fees = fees
		result[0].Total = total
	}

	return &result[0], nil
}

func (p *programsvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error) {
	match := bson.M{"trash": false}

	f, err := helpers.ParseFilters[filter.ProgramFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := f.BuildPipeline(match)

	// Add customer, package, and user lookups first
	pipeline = append(pipeline, customerLookup...)
	pipeline = append(pipeline, packageLookup...)
	pipeline = append(pipeline, updatedByUserLookup...)

	// Calculate duration in days between startDate and endDate
	pipeline = append(pipeline, durationLookup)

	// add favorite pipeline
	pipeline = append(pipeline, interactionsModels.BuildFavoritePipelineWithAuth(ctx, interactionsModels.FaveTypeProgram)...)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := p.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.ProgramRes
	errAg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	// The isFav field should already be populated correctly from the pipeline
	// No manual assignment needed as it comes directly from the aggregation

	var totalPages float64 = math.Ceil(float64(count) / float64(limit))
	pagination := types.Pagination{
		TotalPages: totalPages,
		PerPage:    limit,
		TotalCount: count,
	}

	return &models.ProgramPagination{
		Programs:   result,
		Pagination: pagination,
	}, nil
}

func (p *programsvcs) GetAll(ctx context.Context, query string) ([]models.Program, error) {
	match := bson.M{"trash": false}

	filters, err := helpers.ParseFilters[filter.ProgramFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)
	// add favorite pipeline
	pipeline = append(pipeline, interactionsModels.BuildFavoritePipelineWithAuth(ctx, interactionsModels.FaveTypeProgram)...)
	var result []models.Program
	err = p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAuth retrieves programs with pagination, requiring authentication
func (p *programsvcs) GetAuth(ctx context.Context, skip, limit int64, query string) (*models.ProgramPagination, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as Get but ensure user is authenticated
	return p.Get(ctx, skip, limit, query)
}

// GetAllAuth retrieves all programs without pagination, requiring authentication
func (p *programsvcs) GetAllAuth(ctx context.Context, query string) ([]models.Program, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, helpers.Unauthorized("Authentication required")
	}
	if cfg.User == nil {
		return nil, helpers.Unauthorized("Authentication required")
	}

	// Use the same logic as GetAll but ensure user is authenticated
	return p.GetAll(ctx, query)
}

// add general or custom program - update related travel request
func (p *programsvcs) Add(ctx context.Context, data *models.ProgramDto) (*models.Program, error) {
	result, err := p.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, err
		}
		// Check if this is a general program and has daily itinerary
		if data.ProgramType == "general" && data.GeneralType != nil {
			// Sum up all durations from destinations before saving/returning
			totalDuration := 0
			for _, dest := range data.GeneralType.Destinations {
				totalDuration += dest.Duration
			}
			data.GeneralType.GeneralDuration = totalDuration
		}

		if data.ProgramType == "general" && data.GeneralType != nil {
			// Loop through daily itinerary
			for i, day := range data.GeneralType.DailyItinerary {
				// Check if NewActions exists and has elements
				if len(day.NewActions) > 0 {
					for _, actionName := range day.NewActions {
						newActivity := &pModels.ActivitiesDto{
							Name:        actionName,
							Description: nil, // You can customize this
							// Add other fields as needed
						}
						createdActivity, err := p.activitiesSvcs.Add(ctx, newActivity)
						if err != nil {
							return nil, errors.New("failed to create new activity: " + err.Error())
						}
						data.GeneralType.DailyItinerary[i].Actions = append(data.GeneralType.DailyItinerary[i].Actions, createdActivity.Id)
					}
				}
			}
		}

		program := &models.Program{
			Id:         primitive.NewObjectID(),
			ProgramDto: *data,
			CreatedAt:  time.Now(),
			CreatedBy:  cfg.User.Id,
			UpdatedAt:  time.Now(),
			UpdatedBy:  cfg.User.Id,
		}
		if program.ProgramType == "custom" {
			program.ProgramDto.Status = "waiting"
		}
		if program.TravelReqId != primitive.NilObjectID {
			if err := p.AssignProgramToTravelRequest(ctx, program); err != nil {
				return nil, errors.New("error updating travel request, check if it is already assigned to a program")
			}
		} else if program.ProgramType == "custom" && program.CustomerId != nil {
			// Reverse flow: Create travel request from program when no TravelReqId is provided
			travelReq, err := p.createTravelRequestFromProgram(ctx, program)
			if err != nil {
				return nil, fmt.Errorf("error creating travel request from program: %w", err)
			}
			// Link program to the created travel request
			program.TravelReqId = travelReq.Id
		}

		if err := p.repo.Add(ctx, program); err != nil {
			return nil, err
		}

		return program, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Program), nil
}

func (p *programsvcs) AssignProgramToTravelRequest(ctx context.Context, program *models.Program) error {
	filter := bson.M{
		"_id": program.TravelReqId,
		"$or": []bson.M{
			{"program": primitive.NilObjectID},
			{"program": bson.M{"$exists": false}},
		},
		"trash": false,
	}
	update := bson.M{"$set": bson.M{
		"program":     program.Id,
		"package":     program.Package,
		"status":      enums.TravelReqStatusWaiting,
		"revisionNum": 1,
	}}

	if _, err := p.travelreqsvcs.Patch(ctx, filter, update); err != nil {
		return errors.New("error updating travel request, check if it is already assigned to a program")
	}

	return nil
}

// createTravelRequestFromProgram creates a travel request from a custom program (reverse flow)
// This is called when admin creates a custom program without a travel request
func (p *programsvcs) createTravelRequestFromProgram(ctx context.Context, program *models.Program) (*models.TravelRequest, error) {
	if program.CustomerId == nil {
		return nil, errors.New("customerId is required to create travel request from program")
	}

	if program.CustomType == nil {
		return nil, errors.New("customType is required for custom program")
	}

	// Get customer information
	customer, err := p.customersvcs.GetByFilter(ctx, bson.M{"_id": *program.CustomerId, "trash": false})
	if err != nil {
		return nil, fmt.Errorf("error getting customer: %w", err)
	}

	// Map ProgramServiceType to ServiceType
	serviceType, err := p.mapProgramServiceTypeToServiceType(program.ServiceType)
	if err != nil {
		return nil, err
	}

	// Calculate trip duration from program dates
	tripDuration := 0
	if !program.StartDate.IsZero() && !program.EndDate.IsZero() {
		duration := program.EndDate.Sub(program.StartDate)
		tripDuration = int(duration.Hours() / 24)
		if tripDuration < 1 {
			tripDuration = 1
		}
	}

	// Build TravelRequestDto from Program
	travelReqDto := &models.TravelRequestDto{
		Package:         program.Package,
		ClientName:      customer.Name.GetContentByLang("en"), // Default to English
		ClientPhone:     customer.ClientContact.Mobile,
		ClientEmail:     customer.Security.Email,
		Nationality:     customer.Nationality,
		TripDuration:    tripDuration,
		ServiceType:     serviceType,
		TripCoordinator: program.AgentId,
		ContactMethod:   []string{"email"}, // Default contact method
		SpecialReq:      program.Purpose,
	}

	// Map CustomType fields to TravelRequestDto
	if program.CustomType.Delegation.OrganizationName != nil {
		travelReqDto.Delegation = &program.CustomType.Delegation
	}
	if program.CustomType.BusinessMan.Purpose != "" {
		travelReqDto.BusinessMan = &program.CustomType.BusinessMan
	}
	if program.CustomType.CustomPlan.TripType != "" {
		travelReqDto.CustomPlan = &program.CustomType.CustomPlan
	}
	if program.CustomType.HotelBooking.HotelId != primitive.NilObjectID {
		travelReqDto.HotelBooking = &program.CustomType.HotelBooking
	}

	// Convert ProgramDestination to Destination
	if len(program.CustomType.Destinations) > 0 {
		destinations := make([]models.Destination, 0, len(program.CustomType.Destinations))
		for _, pd := range program.CustomType.Destinations {
			dest := models.Destination{
				DestinationFrom: pd.DestinationFrom,
				DestinationTo:   pd.DestinationTo,
				TripDetails:     pd.TripDetails,
				Accommodation:   make([]models.Accommodation, 0, len(pd.Accommodation)),
				FlightTickets:   pd.FlightTickets.FlightTicket,
				Transportation:  pd.Transportation.Transportation,
				Activities:      []string{}, // ProgramActivities uses transl.Localizable, convert if needed
				Agenda:          pd.Agenda,
				Services:        p.convertProgramServicesToServices(pd.Services),
			}

			// Convert ProgramAccommodation to Accommodation
			for _, pa := range pd.Accommodation {
				dest.Accommodation = append(dest.Accommodation, pa.Accommodation)
			}

			destinations = append(destinations, dest)
		}
		travelReqDto.Destinations = destinations
	}

	// Convert ProgramVipCar to VipCar
	if len(program.CustomType.VipCar.Destinations) > 0 {
		vipCarDests := make([]models.VipCarDest, 0, len(program.CustomType.VipCar.Destinations))
		for _, pvd := range program.CustomType.VipCar.Destinations {
			// ProgramTransportation embeds Transportation, access embedded fields directly
			vipCarDest := models.VipCarDest{
				DestinationFrom: pvd.DestinationFrom,
				DestinationTo:   pvd.DestinationTo,
				Transportation: models.Transportation{
					TransType:             pvd.Transportation.TransType,
					Capacity:              pvd.Transportation.Capacity,
					DriverLanguagesSpoken: pvd.Transportation.DriverLanguagesSpoken,
					LuxuryFeatures:        pvd.Transportation.LuxuryFeatures,
					StartDate:             pvd.Transportation.StartDate,
					EndDate:               pvd.Transportation.EndDate,
					StartTime:             pvd.Transportation.StartTime,
					EndTime:               pvd.Transportation.EndTime,
					CarTypeId:             pvd.Transportation.CarTypeId,
				},
				Services: p.convertProgramServicesToServices(pvd.Services),
			}
			vipCarDests = append(vipCarDests, vipCarDest)
		}
		travelReqDto.VipCar = &models.VipCar{
			Destinations: vipCarDests,
		}
	}

	// Convert ProgramFlightTicketRequest to FlightTicketRequest
	if len(program.CustomType.FlightTicketRequest.Destinations) > 0 {
		flightDests := make([]models.FlightTicketDestination, 0, len(program.CustomType.FlightTicketRequest.Destinations))
		for _, pft := range program.CustomType.FlightTicketRequest.Destinations {
			flightDests = append(flightDests, models.FlightTicketDestination{
				FlightTicket: pft.FlightTicket,
				Services:     p.convertProgramServicesToServices(pft.Services),
			})
		}
		travelReqDto.FlightTicketRequest = &models.FlightTicketRequest{
			Destinations: flightDests,
		}
	}

	// Create travel request
	seq, err := p.sortingsvcs.GetAndUpdateSourceSeq(ctx, "travelRequest")
	if err != nil {
		return nil, fmt.Errorf("error getting sequence: %w", err)
	}

	reqId := fmt.Sprintf("RQ-%d-%d", time.Now().Year(), seq)

	travelReq := &models.TravelRequest{
		Id:               primitive.NewObjectID(),
		ReqId:            reqId,
		TravelRequestDto: *travelReqDto,
		Date:             time.Now(),
		Status:           enums.TravelReqStatusWaiting,
		CustomerId:       *program.CustomerId,
		Program:          program.Id, // Link to program
		CreatedAt:        time.Now(),
		CreatedBy:        program.CreatedBy,
	}

	// Get departure agent using the same logic as normal flow
	departureDestinationId, err := travelReq.GetDepartureDestinationId()
	if err != nil {
		return nil, fmt.Errorf("error getting departure destination: %w", err)
	}

	agent, err := p.agentsvcs.GetAgentByDestination(ctx, departureDestinationId.Hex())
	if err != nil {
		return nil, fmt.Errorf("error get departure destination agent, check if agent has destination and is active: %w", err)
	}

	travelReq.DepartureAgent = agent.Id

	// Save travel request
	if err := p.travelrequestrepo.Add(ctx, travelReq); err != nil {
		return nil, fmt.Errorf("error saving travel request: %w", err)
	}

	return travelReq, nil
}

// mapProgramServiceTypeToServiceType maps ProgramServiceType to ServiceType
func (p *programsvcs) mapProgramServiceTypeToServiceType(programServiceType enums.ProgramServiceType) (enums.ServiceType, error) {
	switch programServiceType {
	case enums.ProgramServiceTypeDelegation:
		return enums.ServiceTypeDelegation, nil
	case enums.ProgramServiceTypeCustomProgram:
		return enums.ServiceTypeCustomPlan, nil
	case enums.ProgramServiceTypeBusinessManTravel:
		return enums.ServiceTypeBusinessMan, nil
	case enums.ProgramServiceTypeVipCar:
		return enums.ServiceTypeVipCar, nil
	case enums.ProgramServiceTypeFlightTicket:
		return enums.ServiceTypeFlightRequest, nil
	case enums.ProgramServiceTypeHotelBooking:
		return enums.ServiceTypeHotelBooking, nil
	default:
		return "", fmt.Errorf("unsupported program service type: %s", programServiceType)
	}
}

// convertProgramServicesToServices converts ProgramServices to Services
func (p *programsvcs) convertProgramServicesToServices(ps models.ProgramServices) models.Services {
	return models.Services{
		OnGroundAssistance:       ps.OnGroundAssistance.Active,
		TravelInsurance:          ps.TravelInsurance.Active,
		VisaAssistance:           ps.VisaAssistance.Active,
		WelcomeKit:               ps.WelcomeKit.Active,
		FreeSimCardWifi:          ps.FreeSimCardWifi.Active,
		ComplimentaryGifts:       ps.ComplimentaryGifts.Active,
		VipAirportServices:       ps.VipAirportServices.Active,
		PersonalTravelConsultant: ps.PersonalTravelConsultant.Active,
		ChildcareServices:        ps.ChildcareServices.Active,
		AccessibilitySupport:     ps.AccessibilitySupport.Active,
	}
}

func (p *programsvcs) UnassignProgramFromTravelRequest(ctx context.Context, program *models.Program) error {
	filter := bson.M{
		"program": program.Id,
	}
	update := bson.M{"$set": bson.M{
		"program":     primitive.NilObjectID,
		"package":     primitive.NilObjectID,
		"revisionNum": 0,
	}}

	_, err := p.travelreqsvcs.Patch(ctx, filter, update)

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	return nil
}

func (p *programsvcs) UpdateTravelRequestRevisionNum(ctx context.Context, program *models.Program) error {
	filter := bson.M{
		"_id":     program.TravelReqId,
		"program": program.Id,
	}
	update := bson.M{
		"$inc": bson.M{"revisionNum": 1},
		"$set": bson.M{"status": enums.TravelReqStatusWaiting},
	}

	if _, err := p.travelreqsvcs.Patch(ctx, filter, update); err != nil {
		return errors.New("error updating travel request")
	}

	return nil
}

func (p *programsvcs) Update(ctx context.Context, id string, data *models.ProgramDto) (*models.Program, error) {
	result, err := p.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		cfg, err := util.GetReqAppCfg(ctx)
		if err != nil {
			return nil, err
		}

		_id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		oldProgram, err := p.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
		if err != nil {
			return nil, errors.New("error fetching program")
		}
		if oldProgram.Status == "rejected" {
			data.Status = "waiting"
		}
		program := &models.Program{
			Id:         _id,
			ProgramDto: *data,
			UpdatedAt:  time.Now(),
			UpdatedBy:  cfg.User.Id,
		}

		if program.TravelReqId == primitive.NilObjectID {
			if err := p.UnassignProgramFromTravelRequest(ctx, program); err != nil {
				return nil, err
			}
		} else {
			currentProgram, err := p.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
			if err != nil {
				return nil, errors.New("error fetching program")
			}

			if currentProgram.TravelReqId == primitive.NilObjectID {
				if err := p.AssignProgramToTravelRequest(ctx, program); err != nil {
					return nil, err
				}
			} else {
				if currentProgram.TravelReqId == program.TravelReqId {
					if err := p.UpdateTravelRequestRevisionNum(ctx, program); err != nil {
						return nil, err
					}
				} else {

					if err := p.UnassignProgramFromTravelRequest(ctx, program); err != nil {
						return nil, err
					}

					if err := p.AssignProgramToTravelRequest(ctx, program); err != nil {
						return nil, err
					}
				}
			}

		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": program}

		updatedProgram, err := p.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedProgram, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Program), nil
}

func (p *programsvcs) UpdateIsFav(ctx context.Context, id string, isFav bool) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{"isFav": isFav}}

	if _, err := p.repo.Patch(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (p *programsvcs) Delete(ctx context.Context, id string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	if _, err := p.repo.Patch(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (p *programsvcs) Count(ctx context.Context, filter any) (int64, error) {
	return p.repo.Count(ctx, filter)
}

// v2
func (p *programsvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.ProgramPagination, error) {

	if err := query.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := query.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"trash": false}},
		{"$match": filter},
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := p.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	pipeline = append(pipeline, interactionsModels.BuildFavoritePipelineWithAuth(ctx, interactionsModels.FaveTypeProgram)...)

	var result []models.ProgramRes
	errAg := p.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.ProgramPagination{
		Programs:   result,
		Pagination: pagination,
	}, nil

}
