package user

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type userQueryBuilder interface {
	FetchOne(appCfg appconfig.AppConfig, model interface{}, filters []services.QueryFilter) error
	UpdateOne(appCfg appconfig.AppConfig, model interface{}, eTag *string) (*validate.Errors, error)
}

type userFetcher struct {
	builder userQueryBuilder
}

// FetchUser fetches an  user given a slice of filters
func (o *userFetcher) FetchUser(appCfg appconfig.AppConfig, filters []services.QueryFilter) (models.User, error) {
	var user models.User
	error := o.builder.FetchOne(appCfg, &user, filters)
	return user, error
}

// NewUserFetcher return an implementation of the UserFetcher interface
func NewUserFetcher(builder userQueryBuilder) services.UserFetcher {
	return &userFetcher{builder}
}
