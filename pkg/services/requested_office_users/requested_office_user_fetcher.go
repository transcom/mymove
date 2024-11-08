package adminuser

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
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

type requestedOfficeUserFetcherPop struct {
}

// FetchOfficeUserByID fetches an office user given an ID
func (o *requestedOfficeUserFetcherPop) FetchRequestedOfficeUserByID(appCtx appcontext.AppContext, id uuid.UUID) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	err := appCtx.DB().Eager("TransportationOffice").Find(&officeUser, id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.OfficeUser{}, apperror.NewNotFoundError(id, "looking for OfficeUser")
		default:
			return models.OfficeUser{}, apperror.NewQueryError("OfficeUser", err, "")
		}
	}

	return officeUser, err
}

// NewOfficeUserFetcherPop return an implementation of the OfficeUserFetcherPop interface
func NewRequestedOfficeUserFetcherPop() services.RequestedOfficeUserFetcherPop {
	return &requestedOfficeUserFetcherPop{}
}
