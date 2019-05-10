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

func (o officeUserListFetcher) FetchOfficeUserList(filters map[string]interface{}) (models.OfficeUsers, error) {
	var officeUsers models.OfficeUsers
	error := o.builder.FetchMany(&officeUsers, filters)
	return officeUsers, error
}

func NewOfficeUserListFetcher(builder officeUserListQueryBuilder) services.OfficeUserListFetcher {
	return officeUserListFetcher{builder}
}
