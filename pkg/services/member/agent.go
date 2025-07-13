package member

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/helpers"
	"larsa-tourism-microservices/pkg/query"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	interactionsModels "larsa-tourism-microservices/pkg/services/interactions/models"
	"larsa-tourism-microservices/pkg/services/member/enums"
	"larsa-tourism-microservices/pkg/services/member/filters"
	"larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/member/repo"
	"larsa-tourism-microservices/pkg/services/messaging"
	messagingmodels "larsa-tourism-microservices/pkg/services/messaging/models"
	"larsa-tourism-microservices/pkg/services/picklist"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	"github.com/goccy/go-json"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AgentSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (*models.Agent, error)
	GetOne(ctx context.Context, agentId string) (*models.Agent, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.AgentWithPagination, error)
	GetAll(ctx context.Context, query string) ([]models.Agent, error)
	Add(ctx context.Context, data *models.AgentDto) (*models.Agent, error)
	Update(ctx context.Context, agentId string, data *models.AgentDto) (*models.Agent, error)
	UpdateStatus(ctx context.Context, agentId string, status string) (*models.Agent, error)
	Delete(ctx context.Context, agentId string) error
	//agent join
	GetOneAgentJoin(ctx context.Context, agentJoinId string) (*models.AgentJoin, error)
	GetJoinRequests(ctx context.Context, skip, limit int64, query any) (*models.AgentJoinPagination, error)
	Join(ctx context.Context, data *models.AgentJoinDto) (*models.AgentJoin, error)
	ConvertToAgent(ctx context.Context, agentId string, data *models.AgentJoinDto) (*models.Agent, error)
	RejectJoin(ctx context.Context, agentId string) error
	SetAsPending(ctx context.Context, agentId string) error
	//destination agent
	GetAgentByDestination(ctx context.Context, destinationId string) (*models.Agent, error)
	GetDestinationAgents(ctx context.Context, query string) ([]models.Agent, error)

	//v2
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.AgentV2Pagination, error)
	GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Agent, error)
	GetJoinRequestsV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.AgentJoinPagination, error)
}

type agentsvcs struct {
	repo            repo.AgentRepo
	agentjoinrepo   repo.AgentJoinRepo
	sortingsvcs     dbsvcs.SortingSvcs
	memberAuthSvcs  MemberAuthSvcs
	messagesvcs     messaging.MessageSvcs
	destinationsvcs picklist.DestinationSvcs
	withtxn         *db.WithTxn
}

func NewAgentSvcs(i *do.Injector) (AgentSvcs, error) {
	return &agentsvcs{
		repo:            do.MustInvoke[repo.AgentRepo](i),
		agentjoinrepo:   do.MustInvoke[repo.AgentJoinRepo](i),
		sortingsvcs:     do.MustInvoke[dbsvcs.SortingSvcs](i),
		memberAuthSvcs:  do.MustInvoke[MemberAuthSvcs](i),
		messagesvcs:     do.MustInvoke[messaging.MessageSvcs](i),
		destinationsvcs: do.MustInvoke[picklist.DestinationSvcs](i),
		withtxn:         do.MustInvoke[*db.WithTxn](i),
	}, nil
}

// agent
func (a *agentsvcs) GetOne(ctx context.Context, agentId string) (*models.Agent, error) {
	_id, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	return a.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
}

func (a *agentsvcs) GetByFilter(ctx context.Context, filter bson.M) (*models.Agent, error) {
	return a.repo.GetByFilter(ctx, filter)
}

func (a *agentsvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.AgentWithPagination, error) {
	match := bson.M{"trash": false}

	filters, err := filters.NewAgentFilter(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := a.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Agent
	errAg := a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.AgentWithPagination{
		Agents:     result,
		Pagination: pagination,
	}, nil
}

func (a *agentsvcs) GetAll(ctx context.Context, query string) ([]models.Agent, error) {
	match := bson.M{"trash": false}

	filters, err := filters.NewAgentFilter(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	var result []models.Agent
	errAg := a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

// add from dashboard
func (a *agentsvcs) Add(ctx context.Context, data *models.AgentDto) (*models.Agent, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	var userId primitive.ObjectID
	var msg *messagingmodels.Message

	result, err := a.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		agent := &models.Agent{
			AgentDto:  *data,
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
			Status:    enums.AgentStatusInactive,
		}

		var password string
		var err error

		userId, password, err = a.memberAuthSvcs.AddCredentials(ctx, agent)
		if err != nil {
			return nil, err
		}

		agent.Id = userId

		seq, err := a.sortingsvcs.GetAndUpdateSourceSeq(ctx, "agent")
		if err != nil {
			return nil, err
		}

		agent.AgentId = fmt.Sprintf("AG-%d", seq)
		agent.Security.NewPassword = ""

		if err := a.repo.Add(ctx, agent); err != nil {
			return nil, err
		}

		msg, err = a.memberAuthSvcs.GetInvitationEmail(ctx, password, agent)
		if err != nil {
			return nil, err
		}

		return agent, nil
	})

	if err != nil {
		if userId != primitive.NilObjectID {
			if err := a.memberAuthSvcs.DeleteCredentials(ctx, userId.Hex()); err != nil {
				fmt.Println("error deleting user", err)
			}
			fmt.Println("user deleted")
		}
		return nil, err
	}

	if err := a.messagesvcs.SendEmail(ctx, msg); err != nil {
		fmt.Println(err)
	}

	return result.(*models.Agent), nil
}

func (a *agentsvcs) Update(ctx context.Context, agentId string, data *models.AgentDto) (*models.Agent, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	var msg *messagingmodels.Message

	result, err := a.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		agent := &models.Agent{
			Id:        _id,
			AgentDto:  *data,
			UpdatedAt: time.Now(),
			UpdatedBy: cfg.User.Id,
		}

		pass := data.Security.NewPassword
		agent.Security.NewPassword = ""

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": agent}

		updatedAgent, err := a.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		msg, err = a.memberAuthSvcs.GetAccountUpdatedEmail(ctx, pass, updatedAgent)
		if err != nil {
			return nil, err
		}

		if err := a.memberAuthSvcs.UpdateCredentials(ctx, pass, agent); err != nil {
			return nil, err
		}

		return updatedAgent, nil
	})

	if err != nil {
		return nil, err
	}

	if err := a.messagesvcs.SendEmail(ctx, msg); err != nil {
		fmt.Println(err)
	}

	return result.(*models.Agent), nil
}

func (a *agentsvcs) UpdateStatus(ctx context.Context, agentId string, status string) (*models.Agent, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": _id, "trash": false}
	update := bson.M{"$set": bson.M{
		"status":    status,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	updatedAgent, err := a.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updatedAgent, nil
}

func (a *agentsvcs) Delete(ctx context.Context, agentId string) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	_id, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"trash":     true,
		"updatedAt": time.Now(),
		"updatedBy": cfg.User.Id,
	}}

	_, err = a.repo.Patch(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// agent join
func (a *agentsvcs) GetOneAgentJoin(ctx context.Context, agentJoinId string) (*models.AgentJoin, error) {
	_id, err := primitive.ObjectIDFromHex(agentJoinId)
	if err != nil {
		return nil, err
	}
	return a.agentjoinrepo.GetByFilter(ctx, bson.M{"_id": _id})
}

func (a *agentsvcs) Join(ctx context.Context, data *models.AgentJoinDto) (*models.AgentJoin, error) {
	agent := &models.AgentJoin{
		Id:           primitive.NewObjectID(),
		AgentJoinDto: *data,
		Status:       enums.AgentJoinStatusPending,
	}

	if err := a.agentjoinrepo.Add(ctx, agent); err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *agentsvcs) ConvertToAgent(ctx context.Context, agentId string, data *models.AgentJoinDto) (*models.Agent, error) {
	result, err := a.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		_id, err := primitive.ObjectIDFromHex(agentId) //agent join id
		if err != nil {
			return nil, err
		}

		agentDto := &models.AgentDto{
			Name:        data.FullName,
			Nationality: data.Nationality,
			SpokenLangs: data.SpokenLangs,
			Company:     data.CompanyName,
			CompanyLogo: data.CompanyLogo,
			Bio:         data.Bio,
			Countries:   data.Countries,
			Contact: models.AgentContact{
				Phone: data.Phone,
				Email: data.Email,
				Web:   "",
			},
			Security: models.MemberSecurity{
				Email:       data.Email,
				NewPassword: "",
			},
		}

		agent, err := a.Add(ctx, agentDto)
		if err != nil {
			return nil, err
		}

		filter := bson.M{"_id": _id, "status": enums.AgentJoinStatusPending}
		update := bson.M{"$set": bson.M{
			"status": enums.AgentJoinStatusConverted,
		}}

		_, errUp := a.agentjoinrepo.Patch(ctx, filter, update)
		if errUp != nil {
			return nil, errUp
		}

		return agent, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Agent), nil
}

func (a *agentsvcs) RejectJoin(ctx context.Context, agentId string) error {
	_id, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id, "status": enums.AgentJoinStatusPending}
	update := bson.M{"$set": bson.M{
		"status": enums.AgentJoinStatusRejected,
	}}

	_, errUp := a.agentjoinrepo.Patch(ctx, filter, update)
	if errUp != nil {
		return errors.New("error update status, only pending joins can be rejected")
	}

	return nil

}

func (a *agentsvcs) SetAsPending(ctx context.Context, agentId string) error {
	_id, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": _id, "status": enums.AgentJoinStatusRejected}
	update := bson.M{"$set": bson.M{
		"status": enums.AgentJoinStatusPending,
	}}

	_, errUp := a.agentjoinrepo.Patch(ctx, filter, update)
	if errUp != nil {
		return errors.New("error update status, only rejected joins can set to pending")
	}

	return nil

}

func (a *agentsvcs) GetJoinRequests(ctx context.Context, skip, limit int64, query any) (*models.AgentJoinPagination, error) {
	match := bson.M{"status": bson.M{"$ne": "converted"}}

	f, err := helpers.ParseFilters[filters.AgentJoinFilter](query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := f.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := a.agentjoinrepo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.AgentJoin
	errAg := a.agentjoinrepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.AgentJoinPagination{
		Agents:     result,
		Pagination: pagination,
	}, nil
}

func (a *agentsvcs) GetAgentByDestination(ctx context.Context, destinationId string) (*models.Agent, error) {
	destination, err := a.destinationsvcs.GetOne(ctx, destinationId)
	if err != nil {
		return nil, errors.New("error get departure destination")
	}

	country := destination.Name

	pipeline := []bson.M{
		{"$match": bson.M{
			"countries": bson.M{"$in": []string{country}},
			"status":    enums.AgentStatusActive,
			"trash":     false,
		}},
	}

	var result []models.Agent
	errAg := a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errors.New("error get agent by destination")
	}

	if len(result) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &result[0], nil
}

//test

func (a *agentsvcs) GetDestinationAgents(ctx context.Context, query string) ([]models.Agent, error) {
	if query == "" {
		return nil, errors.New("query is required")
	}

	type aux struct {
		Ids []primitive.ObjectID `json:"ids"`
	}

	var data aux
	if err := json.Unmarshal([]byte(query), &data); err != nil {
		return nil, err
	}

	ids := data.Ids

	match := bson.M{"trash": false}

	pipeline := []bson.M{
		{"$match": match},
	}

	var destinationLookup = []bson.M{
		{
			"$lookup": bson.M{
				"from": "tourismDestinations",
				"let":  bson.M{"countries": "$countries"},
				"pipeline": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": bson.A{
									bson.M{"$in": bson.A{"$name", "$$countries"}},
									bson.M{"$in": bson.A{"$_id", ids}},
								},
							},
						},
					},
				},
				"as": "destinations",
			},
		},
		{
			"$match": bson.M{
				"destinations": bson.M{"$ne": bson.A{}},
			},
		},
		{
			"$project": bson.M{
				"destinations": 0,
			},
		},
	}

	pipeline = append(pipeline, destinationLookup...)

	var result []models.Agent
	errAg := a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

// v2
func (a *agentsvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.AgentV2Pagination, error) {
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

	count, err := a.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	cfg, err := util.GetReqAppCfg(ctx)
	if err == nil && cfg.User != nil {
		pipeline = append(pipeline, interactionsModels.BuildFavoritePipeline(cfg.User.Id, interactionsModels.FaveTypeAgent)...)
	} else {
		pipeline = append(pipeline, interactionsModels.BuildDefaultFavorite())
	}

	var result []models.AgentRes
	errAg := a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.AgentV2Pagination{
		Agents:     result,
		Pagination: pagination,
	}, nil
}

func (a *agentsvcs) GetAllV2(ctx context.Context, query *query.Conditions) ([]models.Agent, error) {
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

	var result []models.Agent
	errAg := a.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})
	if errAg != nil {
		return nil, errAg
	}

	return result, nil
}

func (a *agentsvcs) GetJoinRequestsV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.AgentJoinPagination, error) {
	if err := query.CheckValid(); err != nil {
		return nil, err
	}

	filter, err := query.ConvertToMongo()
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"status": bson.M{"$ne": "converted"}}},
		{"$match": filter},
	}

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := a.agentjoinrepo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.AgentJoin
	errAg := a.agentjoinrepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.AgentJoinPagination{
		Agents:     result,
		Pagination: pagination,
	}, nil
}
