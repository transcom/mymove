package user

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type userQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
	DeleteOne(appCtx appcontext.AppContext, model interface{}) error
	DeleteMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
}

type userFetcher struct {
	builder userQueryBuilder
}

// FetchUser fetches an  user given a slice of filters
func (o *userFetcher) FetchUser(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.User, error) {
	var user models.User
	err := o.builder.FetchOne(appCtx, &user, filters)
	return user, err
}

// NewUserFetcher return an implementation of the UserFetcher interface
func NewUserFetcher(builder userQueryBuilder) services.UserFetcher {
	return &userFetcher{builder}
}
