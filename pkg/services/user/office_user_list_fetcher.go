package user

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserListQueryBuilder interface {
	FetchMany(model interface{}, fitlers map[string]interface{}) error
}

type officeUserListFetcher struct {
	builder officeUserListQueryBuilder
}

// FetchOfficeUserList is uses the passed query builder to fetch a list of office users
func (o officeUserListFetcher) FetchOfficeUserList(filters map[string]interface{}) (models.OfficeUsers, error) {
	var officeUsers models.OfficeUsers
	error := o.builder.FetchMany(&officeUsers, filters)
	return officeUsers, error
}

// NewOfficeUserListFetcher returns an implementation of OfficeUserListFetcher
func NewOfficeUserListFetcher(builder officeUserListQueryBuilder) services.OfficeUserListFetcher {
	return officeUserListFetcher{builder}
}
