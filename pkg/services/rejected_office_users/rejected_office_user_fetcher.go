package adminuser

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type rejectedOfficeUserQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type rejectedOfficeUserFetcher struct {
	builder rejectedOfficeUserQueryBuilder
}

// FetchRejectedOfficeUser fetches an office user given a slice of filters
func (o *rejectedOfficeUserFetcher) FetchRejectedOfficeUser(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.OfficeUser, error) {
	var rejectedOfficeUser models.OfficeUser
	err := o.builder.FetchOne(appCtx, &rejectedOfficeUser, filters)
	return rejectedOfficeUser, err
}

// NewRejectedUserFetcher return an implementation of the RejectedUserFetcher interface
func NewRejectedOfficeUserFetcher(builder rejectedOfficeUserQueryBuilder) services.RejectedOfficeUserFetcher {
	return &rejectedOfficeUserFetcher{builder}
}
