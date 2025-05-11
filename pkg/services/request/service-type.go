package request

import (
	"context"
	"encoding/json"

	"github.com/samber/do"
)

type TypeSvcs interface {
	Add(ctx context.Context, data json.RawMessage) error // todo: return uniform response
}

type ServiceType struct {
	Svcs     TypeSvcs
	CollName string
}

type ServiceTypes map[string]ServiceType

func NewServiceTypes(i *do.Injector) (ServiceTypes, error) {
	return map[string]ServiceType{
		"vipcar": {
			Svcs:     do.MustInvoke[VipCarRequestSvcs](i),
			CollName: "tourismVipCarRequest",
		},
	}, nil
}
