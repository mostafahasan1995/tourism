package ourservice

import (
	"context"
	"errors"
	"math"
	"time"

	"larsa-tourism-microservices/pkg/services/member"
	membermodels "larsa-tourism-microservices/pkg/services/member/models"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/transl"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	repo                  repo.AgentFinancialRepo
	travelrequestrepo     repo.TravelRequestRepo
	agentsvcs             member.AgentSvcs
	financialsettingssvcs FinancialSettingsSvcs
}

func NewAgentFinancialSvcs(i *do.Injector) (AgentFinancialSvcs, error) {
	return &agentfinancialsvcs{
		repo:                  do.MustInvoke[repo.AgentFinancialRepo](i),
		travelrequestrepo:     do.MustInvoke[repo.TravelRequestRepo](i),
		agentsvcs:             do.MustInvoke[member.AgentSvcs](i),
		financialsettingssvcs: do.MustInvoke[FinancialSettingsSvcs](i),
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

	// Recalculate balance based on actual transactions (replicate GetAgentTransactions logic)
	// This avoids circular dependency with TravelRequestSvcs
	totalProfitFromTransactions, err := s.calculateTotalCommission(ctx, agentObjectID)
	if err != nil {
		// If error calculating, return account with stored values for backward compatibility
		s.limitWithdrawals(account, limit)
		return account, nil
	}

	// Calculate total withdrawn (only approved + pending withdrawals reduce balance)
	// Rejected withdrawals don't reduce balance as they are returned
	var actualWithdrawn float64
	for _, withdrawal := range account.Withdrawals {
		if withdrawal.Status == enums.WithdrawalStatusApproved || withdrawal.Status == enums.WithdrawalStatusPending {
			actualWithdrawn += withdrawal.Amount
		}
	}

	// Update account with calculated values based on actual transactions
	account.TotalProfit = roundTo3Decimals(totalProfitFromTransactions)
	account.TotalWithdrawn = roundTo3Decimals(actualWithdrawn)
	account.Balance = roundTo3Decimals(totalProfitFromTransactions - actualWithdrawn)

	// Limit withdrawals for performance
	s.limitWithdrawals(account, limit)

	return account, nil
}

// calculateTotalCommission replicates the commission calculation from GetAgentTransactions
// to avoid circular dependency
func (s *agentfinancialsvcs) calculateTotalCommission(ctx context.Context, agentObjectID primitive.ObjectID) (float64, error) {
	match := bson.M{
		"status": enums.TravelReqStatusCompleted,
		"trash":  false,
		"$or": []bson.M{
			{"departureAgent": agentObjectID},
			{"tripCoordinator": agentObjectID},
		},
	}

	pipeline := []bson.M{
		{"$match": match},
	}

	// Lookup invoices
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

	// Lookup agents
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

	pipeline = append(pipeline, invoiceLookup...)
	pipeline = append(pipeline, departureAgentLookup...)
	pipeline = append(pipeline, tripCoordinatorAgentLookup...)

	type aux struct {
		models.TravelRequest `bson:",inline"`
		Invoice              models.Invoice     `bson:"invoice" json:"invoice"`
		DepartureAgentData   membermodels.Agent `bson:"departureAgentData" json:"departureAgentData"`
		TripCoordinatorData  membermodels.Agent `bson:"tripCoordinatorData" json:"tripCoordinatorData"`
	}

	var result []aux
	err := s.travelrequestrepo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
		return cur.All(ctx, &result)
	})

	if err != nil {
		return 0, err
	}

	financialSettings, err := s.financialsettingssvcs.Get(ctx)
	if err != nil {
		return 0, errors.New("error get settings")
	}

	profitRatio := financialSettings.ProfitRatio // platform profit ratio
	var totalCommission float64

	loadAgentFinancial := func(agentData membermodels.Agent, agentID primitive.ObjectID) (membermodels.AgentFinancial, error) {
		financial := agentData.Financial
		if !financial.ProfitOfTourismProgram || financial.Ratio == 0 {
			agentDoc, err := s.agentsvcs.GetByFilter(ctx, bson.M{"_id": agentID, "trash": false})
			if err != nil {
				return financial, err
			}
			financial = agentDoc.Financial
		}
		return financial, nil
	}

	for _, r := range result {
		if r.Invoice.Id == primitive.NilObjectID {
			continue
		}

		// Use SubTotal for profit calculation (before fees)
		invoiceSubTotal := r.Invoice.SubTotal
		if invoiceSubTotal == 0 {
			invoiceSubTotal = r.Invoice.Total
		}

		// Calculate client profit: invoiceSubTotal * (profitRatio / 100)
		clientProfit := invoiceSubTotal * (profitRatio / 100)

		// Determine which agent to calculate commission for
		var targetFinancial membermodels.AgentFinancial
		var hasValidFinancial bool

		if r.DepartureAgent == r.TripCoordinator {
			// Same agent for both roles
			if r.DepartureAgent == agentObjectID {
				financial, err := loadAgentFinancial(r.DepartureAgentData, r.DepartureAgent)
				if err == nil && financial.ProfitOfTourismProgram && financial.Ratio > 0 {
					targetFinancial = financial
					hasValidFinancial = true
				}
			}
		} else {
			// Different agents
			if r.DepartureAgent == agentObjectID {
				financial, err := loadAgentFinancial(r.DepartureAgentData, r.DepartureAgent)
				if err == nil && financial.ProfitOfTourismProgram && financial.Ratio > 0 {
					targetFinancial = financial
					hasValidFinancial = true
				}
			} else if r.TripCoordinator == agentObjectID {
				financial, err := loadAgentFinancial(r.TripCoordinatorData, r.TripCoordinator)
				if err == nil && financial.ProfitOfTourismProgram && financial.Ratio > 0 {
					targetFinancial = financial
					hasValidFinancial = true
				}
			}
		}

		// Calculate commission: clientProfit * (agentRatio / 100)
		if hasValidFinancial {
			commission := clientProfit * (targetFinancial.Ratio / 100)
			totalCommission += commission
		}
	}

	return totalCommission, nil
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

	account.TotalProfit = roundTo3Decimals(account.TotalProfit + amount)
	account.Balance = roundTo3Decimals(account.Balance + amount)
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

	// Set default date if not provided
	withdrawalDate := req.Date
	if withdrawalDate.IsZero() {
		withdrawalDate = time.Now()
	}

	withdrawal := models.AgentWithdrawal{
		Id:      primitive.NewObjectID(),
		Amount:  req.Amount,
		Method:  req.Method,
		Note:    req.Note,
		Status:  enums.WithdrawalStatusPending,
		Name:    req.Name,
		Date:    withdrawalDate,
		Email:   req.Email,
		Receipt: req.Receipt,
	}

	filter := bson.M{"_id": account.Id}
	update := bson.M{
		"$push": bson.M{
			"withdrawals": bson.M{
				"$each":     []models.AgentWithdrawal{withdrawal},
				"$position": 0, // Insert at the beginning
			},
		},
		"$inc": bson.M{
			"totalWithdrawn": req.Amount,
			"balance":        -req.Amount,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	_, err = s.repo.Patch(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Reload the account to ensure all fields are properly retrieved
	updated, err := s.repo.GetByFilter(ctx, bson.M{"_id": account.Id})
	if err != nil {
		return nil, err
	}

	// Round balance, totalProfit, and totalWithdrawn to 3 decimal places
	updated.Balance = roundTo3Decimals(updated.Balance)
	updated.TotalProfit = roundTo3Decimals(updated.TotalProfit)
	updated.TotalWithdrawn = roundTo3Decimals(updated.TotalWithdrawn)

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

	// Round balance, totalProfit, and totalWithdrawn to 3 decimal places
	updated.Balance = roundTo3Decimals(updated.Balance)
	updated.TotalProfit = roundTo3Decimals(updated.TotalProfit)
	updated.TotalWithdrawn = roundTo3Decimals(updated.TotalWithdrawn)

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
	account.Balance = roundTo3Decimals(account.Balance + withdrawal.Amount)
	if account.TotalWithdrawn >= withdrawal.Amount {
		account.TotalWithdrawn = roundTo3Decimals(account.TotalWithdrawn - withdrawal.Amount)
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

	// Round balance, totalProfit, and totalWithdrawn to 3 decimal places
	updated.Balance = roundTo3Decimals(updated.Balance)
	updated.TotalProfit = roundTo3Decimals(updated.TotalProfit)
	updated.TotalWithdrawn = roundTo3Decimals(updated.TotalWithdrawn)

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

// roundTo3Decimals rounds a float64 to 3 decimal places
func roundTo3Decimals(value float64) float64 {
	return math.Round(value*1000) / 1000
}
