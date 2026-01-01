package db

import (
	"context"
	"larsa-tourism-microservices/pkg/services/db/repo"

	"github.com/samber/do"
)

type SortingSvcs interface {
	GetAndUpdateSourceSeq(ctx context.Context, name string) (int64, error)
	Reorder(ctx context.Context, id, collName string, newSeq int64) error
}

type sortingsvcs struct {
	repo repo.SortingRepo
}

func NewSortingSvcs(i *do.Injector) (SortingSvcs, error) {
	return &sortingsvcs{
		repo: do.MustInvoke[repo.SortingRepo](i),
	}, nil
}

func (s *sortingsvcs) GetAndUpdateSourceSeq(ctx context.Context, name string) (int64, error) {
	return s.repo.GetAndUpdateSourceSeq(ctx, name)
}

func (s *sortingsvcs) Reorder(ctx context.Context, id, collName string, newSeq int64) error {
	return s.repo.Reorder(ctx, id, collName, newSeq)
}
