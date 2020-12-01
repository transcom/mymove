package officeuser

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	QueryForAssociations(model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error
	CreateOne(model interface{}) (*validate.Errors, error)
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type officeUserFetcher struct {
	builder officeUserQueryBuilder
}

// FetchOfficeUser fetches an office user given a slice of filters
func (o *officeUserFetcher) FetchOfficeUser(filters []services.QueryFilter) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	err := o.builder.FetchOne(&officeUser, filters)
	return officeUser, err
}

// NewOfficeUserFetcher return an implementation of the OfficeUserFetcher interface
func NewOfficeUserFetcher(builder officeUserQueryBuilder) services.OfficeUserFetcher {
	return &officeUserFetcher{builder}
}

// TODO - Eventually move away from the query builder and back to pop
type officeUserFetcherPop struct {
	db *pop.Connection
}

// FetchOfficeUserByID fetches an office user given a slice of filters
func (o *officeUserFetcherPop) FetchOfficeUserByID(id uuid.UUID) (models.OfficeUser, error) {
	var officeUser models.OfficeUser
	err := o.db.Eager("TransportationOffice").Find(&officeUser, id)
	return officeUser, err
}

// NewOfficeUserFetcherPop return an implementation of the OfficeUserFetcherPop interface
func NewOfficeUserFetcherPop(db *pop.Connection) services.OfficeUserFetcherPop {
	return &officeUserFetcherPop{db}
}
