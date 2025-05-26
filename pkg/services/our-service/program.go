package ourservice

import (
	"larsa-tourism-microservices/pkg/services/our-service/repo"

	"github.com/samber/do"
)

type ProgramSvcs interface{}

type programsvcs struct {
	generalProgramRepo repo.GeneralProgramRepo
	customProgramRepo  repo.CustomProgramRepo
}

func NewProgramSvcs(i *do.Injector) (ProgramSvcs, error) {
	return &programsvcs{
		generalProgramRepo: do.MustInvoke[repo.GeneralProgramRepo](i),
		customProgramRepo:  do.MustInvoke[repo.CustomProgramRepo](i),
	}, nil
}

//
