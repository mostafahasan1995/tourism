package gateway

import (
	"context"
	"encoding/json"
	"errors"

	// "errors"

	"fmt"
	"larsa-tourism-microservices/pkg/gateway/models"
	"larsa-tourism-microservices/pkg/util"

	// "net/url"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserNotFound = errors.New("user not found")
var ErrDupliateEmail = errors.New("duplicated email")
var ErrUnKnowen = errors.New("unknowen error")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")

type UsersGw struct{}

func NewUsersGw(i *do.Injector) (*UsersGw, error) {
	return &UsersGw{}, nil
}

func (ugw *UsersGw) AddUser(ctx context.Context, user *models.PostUserData) (addedUserId primitive.ObjectID, errAdd error) {

	zeroId := primitive.NilObjectID

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return zeroId, err
	}

	newUserReqOp := &common.RequestParams{
		Service: "users",
		Path:    "users/",
		Method:  "POST",
		Data: map[string]interface{}{
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"password":     user.Password,
			"roles":        user.Roles,
			"capabilities": user.Capabilities,
		},
		Header: cfg.Hp,
	}

	res, err := common.CallService(newUserReqOp)

	if err != nil {
		return zeroId, errors.New("error adding user")
	} else if res.StatusCode != 200 {

		switch res.StatusCode {
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

	var data AddedUser
	if errDec := json.NewDecoder(res.Body).Decode(&data); errDec != nil {
		return zeroId, errDec
	}

	return data.User.Id, nil
}

func (ugw *UsersGw) UpdateUser(ctx context.Context, userId string, user *models.PostUserData) (UpdatedUserId primitive.ObjectID, errUpdate error) {

	zeroId := primitive.NilObjectID

	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return zeroId, err
	}

	updateUserReqOp := &common.RequestParams{
		Service: "users",
		Path:    fmt.Sprintf("users/%s", userId),
		Method:  "PATCH",
		Data: map[string]interface{}{
			"firstName":    user.FirstName,
			"lastName":     user.LastName,
			"email":        user.Email,
			"password":     user.Password,
			"roles":        user.Roles,
			"capabilities": user.Capabilities,
		},
		Header: cfg.Hp,
	}

	res, err := common.CallService(updateUserReqOp)

	if err != nil || res.StatusCode != 200 {
		return zeroId, errors.New("error updating user")
	}

	type TempUser struct {
		Id primitive.ObjectID `json:"_id"`
	}

	type UpdatedUser struct {
		User TempUser `json:"user"`
	}

	var data UpdatedUser
	if errDec := json.NewDecoder(res.Body).Decode(&data); errDec != nil {
		return zeroId, errDec
	}

	return data.User.Id, nil
}

// func (ugw *UsersGw) UpdateUserRoles(ctx context.Context, userId primitive.ObjectID, roles []primitive.ObjectID) error {
// 	if userId == primitive.NilObjectID || len(roles) == 0 {
// 		return errors.New("invalid input: userId or roles are empty")
// 	}

// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	updateRolesReqOp := &common.RequestParams{
// 		Service: "users",
// 		Path:    fmt.Sprintf("users/%s/roles", userId.Hex()),
// 		Method:  "PATCH",
// 		Data: map[string]interface{}{
// 			"roles": roles,
// 		},
// 		Header: cfg.Hp,
// 	}

// 	res, err := common.CallService(updateRolesReqOp)
// 	if err != nil {
// 		return errors.New("something went wrong while calling the service")
// 	} else if res.StatusCode != 200 {
// 		switch res.StatusCode {
// 		case 400:
// 			return errors.New("bad request: roles array may be empty")
// 		case 404:
// 			return errors.New("user not found")
// 		case 401:
// 			return types.ErrUnauthorized
// 		case 403:
// 			return types.ErrForbidden
// 		default:
// 			return errors.New("unknown error occurred")
// 		}
// 	}

// 	return nil
// }

// func (ugw *UsersGw) GetUserByEmail(ctx context.Context, email string) (*common.User, error) {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//check if user with same email exist
// 	query := fmt.Sprintf(`{"email":"%s"}`, email)
// 	params := url.Values{}
// 	params.Add("query", query)

// 	reqOp := &common.RequestParams{
// 		Service: "users",
// 		Path:    "users/all?" + params.Encode(),
// 		Method:  "GET",
// 		Data:    map[string]interface{}{},
// 		Header:  cfg.Hp,
// 	}

// 	res, err := common.CallService(reqOp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var userlist types.UserList
// 	if errDec := json.NewDecoder(res.Body).Decode(&userlist); errDec != nil {
// 		return nil, errDec
// 	}

// 	if len(userlist.Users) == 0 {
// 		return nil, types.ErrUserNotFound
// 	} else {
// 		return &userlist.Users[0], nil
// 	}
// }

func (ugw *UsersGw) GetUserById(ctx context.Context, id string) (*common.User, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	serviceToken, err := common.GetServiceToken("crmservice", &common.HeaderParams{})
	hp := common.HeaderParams{
		Client:       cfg.Db,
		ServiceToken: serviceToken,
	}
	if err != nil {
		return nil, err
	}
	reqOp := &common.RequestParams{
		Service: "users",
		Path:    fmt.Sprintf("users/all?query={\"_id\":\"%s\"}", id),
		Method:  "GET",
		Data:    map[string]interface{}{},
		Header:  &hp,
	}

	var result struct {
		Users []*common.User
	}
	res, err := common.CallService(reqOp)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&result)
	if len(result.Users) > 0 {
		return result.Users[0], err
	}
	return nil, err
}

// func (ugw *UsersGw) DeleteUser(ctx context.Context, userId string) error {
// 	cfg, err := util.GetReqAppCfg(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	newUserReqOp := &common.RequestParams{
// 		Service: "users",
// 		Path:    fmt.Sprintf("users/%s", userId),
// 		Method:  "DELETE",
// 		Data:    map[string]interface{}{},
// 		Header:  cfg.Hp,
// 	}

// 	res, err := common.CallService(newUserReqOp)

// 	if err != nil {
// 		return errors.New("something went wrong")
// 	} else if res.StatusCode != 200 {

// 		switch res.StatusCode {
// 		case 409:
// 			return types.ErrDupliateEmail
// 		case 401:
// 			return types.ErrUnauthorized
// 		case 403:
// 			return types.ErrForbidden
// 		default:
// 			return err
// 		}

// 	}

// 	return  nil
// }
