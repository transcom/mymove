package adminuser

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type requestedOfficeUserQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
}

type requestedOfficeUserFetcher struct {
	builder requestedOfficeUserQueryBuilder
}

// FetchAdminUser fetches an admin user given a slice of filters
func (o *requestedOfficeUserFetcher) FetchRequestedOfficeUser(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	err := o.builder.FetchOne(appCtx, &officeUser, filters)
	return officeUser, err
}

// NewAdminUserFetcher return an implementation of the AdminUserFetcher interface
func NewRequestedOfficeUserFetcher(builder requestedOfficeUserQueryBuilder) services.RequestedOfficeUserFetcher {
	return &requestedOfficeUserFetcher{builder}
}
