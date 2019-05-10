package user

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserQueryBuilder interface {
	FetchOne(model interface{}, field string, value interface{}) error
}

type officeUserFetcher struct {
	builder officeUserQueryBuilder
}

// FetchOfficeUser fetches an office user for the given field/value pair
func (o officeUserFetcher) FetchOfficeUser(field string, value interface{}) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	error := o.builder.FetchOne(&officeUser, field, value)
	return officeUser, error
}

// NewOfficeUserFetcher return an implementaion of the OfficeUserFetcher interface
func NewOfficeUserFetcher(builder officeUserQueryBuilder) services.OfficeUserFetcher {
	return officeUserFetcher{builder}
}
