package officeuser

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
	UpdateOne(model interface{}) (*validate.Errors, error)
}

type fetchOfficeUserQueryBuilder interface {
	query.FetchOne
}

type officeUserFetcher struct {
	builder fetchOfficeUserQueryBuilder
}

// FetchOfficeUser fetches an office user given a slice of filters
func (o *officeUserFetcher) FetchOfficeUser(filters []services.QueryFilter) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	error := o.builder.Filters(filters).Execute(&officeUser)
	return officeUser, error
}

// NewOfficeUserFetcher return an implementation of the OfficeUserFetcher interface
func NewOfficeUserFetcher(builder fetchOfficeUserQueryBuilder) services.OfficeUserFetcher {
	return &officeUserFetcher{builder}
}
