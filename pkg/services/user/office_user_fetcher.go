package user

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type officeUserFetcher struct {
	builder officeUserQueryBuilder
}

// FetchOfficeUser fetches an office user for the given a slice of filters
func (o *officeUserFetcher) FetchOfficeUser(filters []services.QueryFilter) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	error := o.builder.FetchOne(&officeUser, filters)
	return officeUser, error
}

// NewOfficeUserFetcher return an implementaion of the OfficeUserFetcher interface
func NewOfficeUserFetcher(builder officeUserQueryBuilder) services.OfficeUserFetcher {
	return &officeUserFetcher{builder}
}
