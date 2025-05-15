package di

import (
	"larsa-tourism-microservices/pkg/gateway"

	"github.com/samber/do"
)

func Init(i *do.Injector) *do.Injector {

	do.Provide(i, gateway.NewUsersGw)

	return i
}
