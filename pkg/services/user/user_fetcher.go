package user

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type userQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type userFetcher struct {
	builder userQueryBuilder
}

// FetchUser fetches an  user given a slice of filters
func (o *userFetcher) FetchUser(filters []services.QueryFilter) (models.User, error) {
	var user models.User
	error := o.builder.FetchOne(&user, filters)
	return user, error
}

// NewUserFetcher return an implementation of the UserFetcher interface
func NewUserFetcher(builder userQueryBuilder) services.UserFetcher {
	return &userFetcher{builder}
}
