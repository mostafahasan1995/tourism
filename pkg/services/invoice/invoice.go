package invoice

import (
	"context"
	"errors"
	"fmt"
	"larsa-tourism-microservices/pkg/db"
	dbsvcs "larsa-tourism-microservices/pkg/services/db"
	"larsa-tourism-microservices/pkg/services/invoice/filters"
	"larsa-tourism-microservices/pkg/services/invoice/models"
	"larsa-tourism-microservices/pkg/services/invoice/repo"
	"larsa-tourism-microservices/pkg/types"
	"math"
	"time"

	"github.com/samber/do"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceSvcs interface {
	GetOne(ctx context.Context, id string) (*models.Invoice, error)
	Get(ctx context.Context, skip, limit int64, query string) (*models.InvoicePagination, error)
	Add(ctx context.Context, data *models.InvoiceDto) (*models.Invoice, error)
}

type invoiceSvcs struct {
	repo        repo.InvoiceRepo
	sortingsvcs dbsvcs.SortingSvcs
	withtxn     *db.WithTxn
}

func NewInvoiceRepo(i *do.Injector) (InvoiceSvcs, error) {
	return &invoiceSvcs{
		repo:        do.MustInvoke[repo.InvoiceRepo](i),
		sortingsvcs: do.MustInvoke[dbsvcs.SortingSvcs](i),
		withtxn:     do.MustInvoke[*db.WithTxn](i),
	}, nil
}

func (i *invoiceSvcs) Add(ctx context.Context, data *models.InvoiceDto) (*models.Invoice, error) {
	result, err := i.withtxn.Exec(ctx, func(ctx mongo.SessionContext) (any, error) {
		invoice := &models.Invoice{
			Id:         primitive.NewObjectID(),
			InvoiceDto: *data,
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

func (i *invoiceSvcs) GetOne(ctx context.Context, id string) (*models.Invoice, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return i.repo.GetByFilter(ctx, bson.M{"_id": _id})
}

func (i *invoiceSvcs) Get(ctx context.Context, skip, limit int64, query string) (*models.InvoicePagination, error) {
	match := bson.M{}

	filters, err := filters.NewInvoiceFilter(query)
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
