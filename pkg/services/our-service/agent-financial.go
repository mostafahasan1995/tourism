package ourservice

import (
	"context"
	"errors"
	"time"

	membermodels "larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/transl"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentFinancialSvcs interface {
	GetAccount(ctx context.Context, agentId string, limit int) (*models.AgentFinancialAccount, error)
	CreateOrUpdateAccount(ctx context.Context, agentId string, dto *models.AgentFinancialAccountDto) (*models.AgentFinancialAccount, error)
	AddProfit(ctx context.Context, agent *membermodels.Agent, amount float64) (*models.AgentFinancialAccount, error)
	Withdraw(ctx context.Context, agentId string, req *models.AgentWithdrawRequest, limit int) (*models.AgentFinancialAccount, error)
	ApproveWithdrawal(ctx context.Context, agentId string, withdrawalId string) (*models.AgentFinancialAccount, error)
	RejectWithdrawal(ctx context.Context, agentId string, withdrawalId string) (*models.AgentFinancialAccount, error)
}

type agentfinancialsvcs struct {
	repo repo.AgentFinancialRepo
}

func NewAgentFinancialSvcs(i *do.Injector) (AgentFinancialSvcs, error) {
	return &agentfinancialsvcs{
		repo: do.MustInvoke[repo.AgentFinancialRepo](i),
	}, nil
}

func (s *agentfinancialsvcs) ensureAccount(ctx context.Context, agentObjectID primitive.ObjectID, agentName transl.Localizable[string]) (*models.AgentFinancialAccount, error) {
	account, err := s.repo.GetByFilter(ctx, bson.M{"agentId": agentObjectID})
	if err == nil {
		// Account exists, update agent name if provided and different
		if len(agentName) > 0 {
			if account.AgentName == nil || !equalLocalizable(account.AgentName, agentName) {
				filter := bson.M{"_id": account.Id}
				update := bson.M{
					"$set": bson.M{
						"agentName": agentName,
						"updatedAt": time.Now(),
					},
				}
				updated, err := s.repo.Patch(ctx, filter, update)
				if err != nil {
					return nil, err
				}
				return updated, nil
			}
		}
		return account, nil
	}

	// Account doesn't exist, create a new one
	now := time.Now()
	account = &models.AgentFinancialAccount{
		Id:             primitive.NewObjectID(),
		AgentId:        agentObjectID,
		AgentName:      agentName,
		TotalProfit:    0,
		TotalWithdrawn: 0,
		Balance:        0,
		Withdrawals:    []models.AgentWithdrawal{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.Add(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *agentfinancialsvcs) GetAccount(ctx context.Context, agentId string, limit int) (*models.AgentFinancialAccount, error) {
	agentObjectID, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	account, err := s.ensureAccount(ctx, agentObjectID, nil)
	if err != nil {
		return nil, err
	}

	// Limit withdrawals for performance
	s.limitWithdrawals(account, limit)

	return account, nil
}

func (s *agentfinancialsvcs) CreateOrUpdateAccount(ctx context.Context, agentId string, dto *models.AgentFinancialAccountDto) (*models.AgentFinancialAccount, error) {
	if dto == nil {
		return nil, errors.New("dto is required")
	}

	agentObjectID, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	account, err := s.repo.GetByFilter(ctx, bson.M{"agentId": agentObjectID})
	if err != nil {
		// Account doesn't exist, create a new one
		now := time.Now()
		account = &models.AgentFinancialAccount{
			Id:              primitive.NewObjectID(),
			AgentId:         agentObjectID,
			PaymentMethod:   dto.PaymentMethod,
			BankAccountInfo: dto.BankAccountInfo,
			TotalProfit:     0,
			TotalWithdrawn:  0,
			Balance:         0,
			Withdrawals:     []models.AgentWithdrawal{},
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if err := s.repo.Add(ctx, account); err != nil {
			return nil, err
		}

		return account, nil
	}

	// Account exists, update payment info
	account.PaymentMethod = dto.PaymentMethod
	account.BankAccountInfo = dto.BankAccountInfo
	account.UpdatedAt = time.Now()

	filter := bson.M{"_id": account.Id}
	update := bson.M{
		"$set": bson.M{
			"paymentMethod":   account.PaymentMethod,
			"bankAccountInfo": account.BankAccountInfo,
			"updatedAt":       account.UpdatedAt,
		},
	}

	updated, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *agentfinancialsvcs) AddProfit(ctx context.Context, agent *membermodels.Agent, amount float64) (*models.AgentFinancialAccount, error) {
	if agent == nil {
		return nil, errors.New("agent is required")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	account, err := s.ensureAccount(ctx, agent.Id, agent.Name)
	if err != nil {
		return nil, err
	}

	account.TotalProfit += amount
	account.Balance += amount
	account.UpdatedAt = time.Now()

	filter := bson.M{"_id": account.Id}
	update := bson.M{"$set": account}

	updated, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *agentfinancialsvcs) Withdraw(ctx context.Context, agentId string, req *models.AgentWithdrawRequest, limit int) (*models.AgentFinancialAccount, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	agentObjectID, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	account, err := s.ensureAccount(ctx, agentObjectID, nil)
	if err != nil {
		return nil, err
	}

	if req.Amount > account.Balance {
		return nil, errors.New("insufficient balance")
	}

	withdrawal := models.AgentWithdrawal{
		Id:      primitive.NewObjectID(),
		Amount:  req.Amount,
		Method:  req.Method,
		Note:    req.Note,
		Status:  enums.WithdrawalStatusPending,
		Name:    req.Name,
		Date:    req.Date,
		Email:   req.Email,
		Receipt: req.Receipt,
	}

	account.TotalWithdrawn += req.Amount
	account.Balance -= req.Amount
	account.Withdrawals = append([]models.AgentWithdrawal{withdrawal}, account.Withdrawals...)
	account.UpdatedAt = time.Now()

	filter := bson.M{"_id": account.Id}
	update := bson.M{"$set": account}

	updated, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Limit withdrawals for performance
	s.limitWithdrawals(updated, limit)

	return updated, nil
}

func (s *agentfinancialsvcs) ApproveWithdrawal(ctx context.Context, agentId string, withdrawalId string) (*models.AgentFinancialAccount, error) {
	agentObjectID, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	withdrawalObjectID, err := primitive.ObjectIDFromHex(withdrawalId)
	if err != nil {
		return nil, err
	}

	account, err := s.ensureAccount(ctx, agentObjectID, nil)
	if err != nil {
		return nil, err
	}

	// Find the withdrawal in the account
	var withdrawalIndex = -1
	var withdrawal *models.AgentWithdrawal
	for i := range account.Withdrawals {
		if account.Withdrawals[i].Id == withdrawalObjectID {
			withdrawalIndex = i
			withdrawal = &account.Withdrawals[i]
			break
		}
	}

	if withdrawal == nil {
		return nil, errors.New("withdrawal not found")
	}

	if withdrawal.Status != enums.WithdrawalStatusPending {
		return nil, errors.New("withdrawal is not pending, cannot approve")
	}

	// Update withdrawal status to approved
	account.Withdrawals[withdrawalIndex].Status = enums.WithdrawalStatusApproved
	account.UpdatedAt = time.Now()

	filter := bson.M{"_id": account.Id}
	update := bson.M{"$set": account}

	updated, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *agentfinancialsvcs) RejectWithdrawal(ctx context.Context, agentId string, withdrawalId string) (*models.AgentFinancialAccount, error) {
	agentObjectID, err := primitive.ObjectIDFromHex(agentId)
	if err != nil {
		return nil, err
	}

	withdrawalObjectID, err := primitive.ObjectIDFromHex(withdrawalId)
	if err != nil {
		return nil, err
	}

	account, err := s.ensureAccount(ctx, agentObjectID, nil)
	if err != nil {
		return nil, err
	}

	// Find the withdrawal in the account
	var withdrawalIndex = -1
	var withdrawal *models.AgentWithdrawal
	for i := range account.Withdrawals {
		if account.Withdrawals[i].Id == withdrawalObjectID {
			withdrawalIndex = i
			withdrawal = &account.Withdrawals[i]
			break
		}
	}

	if withdrawal == nil {
		return nil, errors.New("withdrawal not found")
	}

	if withdrawal.Status != enums.WithdrawalStatusPending {
		return nil, errors.New("withdrawal is not pending, cannot reject")
	}

	// Reverse the withdrawal: return amount to balance and adjust totals
	account.Balance += withdrawal.Amount
	if account.TotalWithdrawn >= withdrawal.Amount {
		account.TotalWithdrawn -= withdrawal.Amount
	} else {
		account.TotalWithdrawn = 0
	}

	account.Withdrawals[withdrawalIndex].Status = enums.WithdrawalStatusRejected
	account.UpdatedAt = time.Now()

	filter := bson.M{"_id": account.Id}
	update := bson.M{"$set": account}

	updated, err := s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// limitWithdrawals limits the withdrawals array to the specified limit for performance
// Since withdrawals are stored with newest first (prepended), we take the first 'limit' items
func (s *agentfinancialsvcs) limitWithdrawals(account *models.AgentFinancialAccount, limit int) {
	if limit <= 0 {
		return // Don't limit if limit is 0 or negative
	}
	if len(account.Withdrawals) > limit {
		account.Withdrawals = account.Withdrawals[:limit]
	}
}

func equalLocalizable(a transl.Localizable[string], b transl.Localizable[string]) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
