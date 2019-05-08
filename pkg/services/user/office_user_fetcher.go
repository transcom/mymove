package user

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserQueryBuilder interface {
	FetchOne(model interface{}, field string, value string) error
}

type officeUserFetcher struct {
	builder officeUserQueryBuilder
}

func (o officeUserFetcher) FetchOfficeUser(field string, value string) (models.OfficeUser, error) {
	officeUser := models.OfficeUser{}
	error := o.builder.FetchOne(&officeUser, field, value)
	return officeUser, error
}

func NewOfficeUserFetcher(builder officeUserQueryBuilder) services.OfficeUserFetcher {
	return officeUserFetcher{builder}
}
