package ourservice

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	"larsa-tourism-microservices/pkg/query"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/messaging"
	"larsa-tourism-microservices/pkg/services/our-service/enums"
	"larsa-tourism-microservices/pkg/services/our-service/filter"
	"larsa-tourism-microservices/pkg/services/our-service/models"
	"larsa-tourism-microservices/pkg/services/our-service/repo"
	"larsa-tourism-microservices/pkg/types"
	"larsa-tourism-microservices/pkg/util"
	"math"
	"time"

	messagingenums "larsa-tourism-microservices/pkg/services/messaging/enums"
	messagingmodels "larsa-tourism-microservices/pkg/services/messaging/models"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Invoice, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.InvoicePagination, error)
	Add(ctx context.Context, data *models.InvoiceDto) (*models.Invoice, error)
	Update(ctx context.Context, id string, data *models.InvoiceDto) (*models.Invoice, error)
	UpdateStatus(ctx context.Context, id string, status enums.InvoiceStatus) (*models.Invoice, error)
	AddPayment(ctx context.Context, invoiceId string, data *models.PaymentDto) (*models.Invoice, error)
	UpdatePayment(ctx context.Context, invoiceId, paymentId string, data *models.PaymentDto) (*models.Invoice, error)
	DeletePayment(ctx context.Context, invoiceId, paymentId string) (*models.Invoice, error)
	PayOrder(ctx context.Context, invoiceId string, data *models.PayOrder) (*models.Invoice, error)
	//
	AddInvoiceForTravelRequest(ctx context.Context, travelReqId primitive.ObjectID, data *models.InvoiceDto) (*models.Invoice, error)
	SendInvoice(ctx context.Context, data *models.SendInvoiceDto) error
	//v2
	GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.InvoicePagination, error)
}

type invoiceSvcs struct {
	repo        repo.InvoiceRepo
	sortingsvcs dbsvcs.SortingSvcs
	messagesvcs messaging.MessageSvcs
	withtxn     *db.WithTxn
}

func NewInvoiceSvcs(i *do.Injector) (InvoiceSvcs, error) {
	return &invoiceSvcs{
		repo:        do.MustInvoke[repo.InvoiceRepo](i),
		sortingsvcs: do.MustInvoke[dbsvcs.SortingSvcs](i),
		messagesvcs: do.MustInvoke[messaging.MessageSvcs](i),
		withtxn:     do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (i *invoiceSvcs) GetOne(ctx context.Context, id string) (*models.Invoice, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return i.repo.GetByFilter(ctx, bson.M{"_id": _id})
}

func (i *invoiceSvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.InvoicePagination, error) {
	match := bson.M{}

	filters, err := filter.NewInvoiceFilter(query)
	if err != nil {
		return nil, errors.New("invalid query")
	}

	pipeline := filters.BuildPipeline(match)

	countPipeline := make([]bson.M, len(pipeline))
	copy(countPipeline, pipeline)

	count, err := i.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Invoice
	errAg := i.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.InvoicePagination{
		Invoices:   result,
		Pagination: pagination,
	}, nil
}

func (i *invoiceSvcs) Add(ctx context.Context, data *models.InvoiceDto) (*models.Invoice, error) {
	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice := &models.Invoice{
			Id:         primitive.NewObjectID(),
			InvoiceDto: *data,
			Payments:   []models.Payment{},
		
		}
		invoice.Status = enums.InvoiceStatusPending
		if err := invoice.SetTotals(); err != nil {
			return nil, err
		}

		seq, err := i.sortingsvcs.GetAndUpdateSourceSeq(ctx, "invoice")
		if err != nil {
			return nil, err
		}

		invoice.InvoiceId = fmt.Sprintf("INV-%d-%d", time.Now().Year(), seq)

		if err := i.repo.Add(ctx, invoice); err != nil {
			return nil, err
		}

		return invoice, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Invoice), nil

}

func (i *invoiceSvcs) Update(ctx context.Context, id string, data *models.InvoiceDto) (*models.Invoice, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice, err := i.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
		if err != nil {
			return nil, errors.New("invoice not found")
		}

		invoice.InvoiceDto = *data

		if err := invoice.SetTotals(); err != nil {
			return nil, err
		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": invoice}

		updatedInvoice, err := i.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedInvoice, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Invoice), nil

}

func (i *invoiceSvcs) UpdateStatus(ctx context.Context, id string, status enums.InvoiceStatus) (*models.Invoice, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice, err := i.repo.GetByFilter(ctx, bson.M{"_id": _id, "trash": false})
		if err != nil {
			return nil, errors.New("invoice not found")
		}

		now := time.Now()

		filter := bson.M{"_id": _id}
		update := bson.M{
			"$set": bson.M{
				"status":    status,
				"updatedAt": now,
				"updatedBy": cfg.User.Id,
			},
		}

		if _, err := i.repo.Patch(ctx, filter, update); err != nil {
			return nil, err
		}

		invoice.Status = status
		invoice.UpdatedAt = now
		invoice.UpdatedBy = cfg.User.Id

		return invoice, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Invoice), nil
}

// when add payment by admin
func (i *invoiceSvcs) AddPayment(ctx context.Context, invoiceId string, data *models.PaymentDto) (*models.Invoice, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(invoiceId)
	if err != nil {
		return nil, err
	}

	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice, err := i.repo.GetByFilter(ctx, bson.M{"_id": _id})
		if err != nil {
			return nil, errors.New("invoice not found")
		}

		payments := invoice.Payments

		newPayment := models.Payment{
			Id:         primitive.NewObjectID(),
			PaymentDto: *data,
			CreatedAt:  time.Now(),
			CreatedBy:  cfg.User.Id,
		}

		payments = append(payments, newPayment)
		invoice.Payments = payments

		if err := invoice.SetTotals(); err != nil {
			return nil, err
		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": invoice}

		updatedInvoice, err := i.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedInvoice, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Invoice), nil

}

func (i *invoiceSvcs) UpdatePayment(ctx context.Context, invoiceId, paymentId string, data *models.PaymentDto) (*models.Invoice, error) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(invoiceId)
	if err != nil {
		return nil, err
	}

	pId, err := primitive.ObjectIDFromHex(paymentId)
	if err != nil {
		return nil, err
	}

	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice, err := i.repo.GetByFilter(ctx, bson.M{"_id": _id})
		if err != nil {
			return nil, errors.New("invoice not found")
		}

		var found bool
		var payments []models.Payment
		for _, payment := range invoice.Payments {
			if payment.Id == pId {
				payment.PaymentDto = *data
				payment.UpdatedAt = time.Now()
				payment.UpdatedBy = cfg.User.Id
				found = true
			}
			payments = append(payments, payment)
		}

		if !found {
			return nil, errors.New("payment not found")
		}

		invoice.Payments = payments

		if err := invoice.SetTotals(); err != nil {
			return nil, err
		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": invoice}

		updatedInvoice, err := i.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedInvoice, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Invoice), nil

}

func (i *invoiceSvcs) DeletePayment(ctx context.Context, invoiceId, paymentId string) (*models.Invoice, error) {
	_id, err := primitive.ObjectIDFromHex(invoiceId)
	if err != nil {
		return nil, err
	}

	pId, err := primitive.ObjectIDFromHex(paymentId)
	if err != nil {
		return nil, err
	}

	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice, err := i.repo.GetByFilter(ctx, bson.M{"_id": _id})
		if err != nil {
			return nil, errors.New("invoice not found")
		}

		payments := util.SliceFilter(invoice.Payments, func(payment models.Payment) bool {
			return payment.Id != pId
		})

		invoice.Payments = payments
		if err := invoice.SetTotals(); err != nil {
			return nil, err
		}

		filter := bson.M{"_id": _id}
		update := bson.M{"$set": invoice}

		updatedInvoice, err := i.repo.Patch(ctx, filter, update)
		if err != nil {
			return nil, err
		}

		return updatedInvoice, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Invoice), nil

}

// when pay by customer
func (i *invoiceSvcs) PayOrder(ctx context.Context, invoiceId string, data *models.PayOrder) (*models.Invoice, error) {
	newPaymentDto := models.PaymentDto{
		Date:    time.Now(),
		Method:  data.Method,
		Amount:  data.Amount,
		Status:  enums.PaymentStatusUnpaid,
		Receipt: data.Receipt,
	}

	return i.AddPayment(ctx, invoiceId, &newPaymentDto)

}

// add invoice for travel request
func (i *invoiceSvcs) AddInvoiceForTravelRequest(ctx context.Context, travelReqId primitive.ObjectID, data *models.InvoiceDto) (*models.Invoice, error) {
	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice := &models.Invoice{
			Id:          primitive.NewObjectID(),
			InvoiceDto:  *data,
			TravelReqId: travelReqId,
			//DepartureAgent:   invoiceTravelRequestData.DepartureAgent,
			// DestinationAgent: invoiceTravelRequestData.DestinationAgent,
			Payments: []models.Payment{},
			invoice.Status = enums.InvoiceStatusPending
		}

		if err := invoice.SetTotals(); err != nil {
			return nil, err
		}

		seq, err := i.sortingsvcs.GetAndUpdateSourceSeq(ctx, "invoice")
		if err != nil {
			return nil, err
		}

		invoice.InvoiceId = fmt.Sprintf("INV-%d-%d", time.Now().Year(), seq)

		if err := i.repo.Add(ctx, invoice); err != nil {
			return nil, err
		}

		return invoice, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Invoice), nil
}

func (i *invoiceSvcs) SendInvoice(ctx context.Context, data *models.SendInvoiceDto) error {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return err
	}

	msg := messagingmodels.Message{
		Type:        messagingenums.None,
		Email:       data.To,
		Subject:     data.Subject,
		Message:     data.Message,
		MessageHtml: data.Message,
		Attachments: data.Files,
		SenderId:    cfg.User.Id,
		CreatedBy:   cfg.User.Id,
		CreatedAt:   time.Now(),
	}

	type Attachment struct {
		FileName string `json:"filename"`
		Href     string `json:"href"`
	}

	attachments := []Attachment{}

	for _, file := range data.Files {
		attch := Attachment{
			FileName: file.OriginalName,
			Href:     file.Path,
		}

		attachments = append(attachments, attch)

	}

	msg.Others = map[string]interface{}{
		"attachments": attachments,
	}

	if err := i.messagesvcs.SendEmail(ctx, &msg); err != nil {
		return err
	}

	return nil
}

// v2
func (i *invoiceSvcs) GetV2(ctx context.Context, skip, limit int64, query *query.Conditions) (*models.InvoicePagination, error) {
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

	count, err := i.repo.Count(ctx, countPipeline)
	if err != nil {
		return nil, err
	}

	pipeline = append(pipeline, bson.M{"$sort": bson.M{"_id": -1}})
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})

	var result []models.Invoice
	errAg := i.repo.Aggregate(ctx, pipeline, func(cur *mongo.Cursor) error {
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

	return &models.InvoicePagination{
		Invoices:   result,
		Pagination: pagination,
	}, nil
}
