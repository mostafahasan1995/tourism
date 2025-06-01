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

type AgentSvcs interface {
	GetByFilter(ctx context.Context, filter bson.M) (*models.Agent, error)
	GetOne(ctx context.Context, agentId string) (*models.Agent, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.AgentWithPagination, error)
	Add(ctx context.Context, data *models.AgentDto) (*models.Agent, error)
	Update(ctx context.Context, agentId string, data *models.AgentDto) (*models.Agent, error)
	Delete(ctx context.Context, agentId string) error
	//agent join
	GetOneAgentJoin(ctx context.Context, agentJoinId string) (*models.AgentJoin, error)
	GetJoinRequests(ctx context.Context, skip, limit int64, query string) (*models.AgentJoinPagination, error)
	Join(ctx context.Context, data *models.AgentJoinDto) (*models.AgentJoin, error)
	ConvertToAgent(ctx context.Context, agentId string, data *models.AgentJoinDto) (*models.Agent, error)
	RejectJoin(ctx context.Context, agentId string) error
}

type agentsvcs struct {
	repo          repo.AgentRepo
	agentjoinrepo repo.AgentJoinRepo
	sortingsvcs   dbsvcs.SortingSvcs
	usersgw       *gateway.UsersGw
	gateway       gateway.Gateway
	withtxn       *db.WithTxn
}

func NewAgentSvcs(i *do.Injector) (AgentSvcs, error) {
	return &agentsvcs{
		repo:          do.MustInvoke[repo.AgentRepo](i),
		agentjoinrepo: do.MustInvoke[repo.AgentJoinRepo](i),
		sortingsvcs:   do.MustInvoke[dbsvcs.SortingSvcs](i),
		usersgw:       do.MustInvoke[*gateway.UsersGw](i),
		gateway:       do.MustInvoke[gateway.Gateway](i),
		withtxn:       do.MustInvoke[*db.WithTxn](i),
	}, nil
}

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

// add from dashboard
func (a *agentsvcs) Add(ctx context.Context, data *models.AgentDto) (*models.Agent, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}
	result, err := a.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		agent := &models.Agent{
			AgentDto:  *data,
			CreatedAt: time.Now(),
			CreatedBy: cfg.User.Id,
			Status:    "inactive", //active, inactive
		}

		userId, err := a.AddUpdateAgentCredentials(ctx, agent)
		if err != nil {
			return nil, err
		}

		agent.Id = userId

		seq, err := a.sortingsvcs.GetAndUpdateSourceSeq(ctx, "agent")
		if err != nil {
			return nil, err
		}

		agent.AgentId = fmt.Sprintf("AG-%d", seq)

		if err := a.repo.Add(ctx, agent); err != nil {
			return nil, err
		}

		return agent, nil
	})

	if err != nil {
		return nil, err
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

	result, err := a.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		agent := &models.Agent{
			Id:        _id,
			AgentDto:  *data,
			UpdatedAt: time.Now(),
			UpdatedBy: cfg.User.Id,
		}

		_, err := a.AddUpdateAgentCredentials(ctx, agent)
		if err != nil {
			return nil, err
		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": agent}

		updatedAgent, err := a.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedAgent, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Agent), nil
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

func (a *agentsvcs) AddUpdateAgentCredentials(ctx context.Context, data *models.Agent) (userId primitive.ObjectID, err error) {
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

	resp, err := a.gateway.Request(ctx, "users", path, method, "", user)

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
		Status:       "pending",
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
				Phone: models.AgentPhone{
					Pre:     "",
					Content: data.Phone.Content,
				},
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

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": bson.M{
			"status": "converted",
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

	filter := bson.M{"_id": _id}
	update := bson.M{"$set": bson.M{
		"status": "rejected",
	}}

	_, errUp := a.agentjoinrepo.Patch(ctx, filter, update)
	if errUp != nil {
		return errUp
	}

	return nil

}

func (a *agentsvcs) GetJoinRequests(ctx context.Context, skip, limit int64, query string) (*models.AgentJoinPagination, error) {
	match := bson.M{}

	filter, err := filters.NewAgentJoinFilter(query)
	if err != nil {
		return nil, err
	}

	pipeline := filter.BuildPipeline(match)

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
