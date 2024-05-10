package adminuser

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type requestedOfficeUserQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type requestedOfficeUserFetcher struct {
	builder requestedOfficeUserQueryBuilder
}

// FetchRequestedOfficeUser fetches an office user given a slice of filters
func (o *requestedOfficeUserFetcher) FetchRequestedOfficeUser(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.OfficeUser, error) {
	var requestedOfficeUser models.OfficeUser
	err := o.builder.FetchOne(appCtx, &requestedOfficeUser, filters)
	return requestedOfficeUser, err
}

// NewAdminUserFetcher return an implementation of the AdminUserFetcher interface
func NewRequestedOfficeUserFetcher(builder requestedOfficeUserQueryBuilder) services.RequestedOfficeUserFetcher {
	return &requestedOfficeUserFetcher{builder}
}
