package accesscode

import (
	"fmt"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"strings"
)

type accessCodeListQueryBuilder interface {
	QueryAssociations(model interface{}, associations []services.QueryAssociation, filters []services.QueryFilter) error
}

type accessCodeListFetcher struct {
	builder accessCodeListQueryBuilder
}

// FetchAccessCodeList uses the passed query builder to fetch a list of access codes
func (o *accessCodeListFetcher) FetchAccessCodeList(filters []services.QueryFilter, associations []services.QueryAssociation) (models.AccessCodes, error) {
	var accessCodes models.AccessCodes

	fmt.Println(strings.Repeat("*", 100))
	fmt.Println(accessCodes)
	fmt.Println(len(accessCodes))
	error := o.builder.QueryAssociations(&accessCodes, associations, filters)

	return accessCodes, error
}

// NewAccessCodeListFetcher returns an implementation of OfficeUserListFetcher
func NewAccessCodeListFetcher(builder accessCodeListQueryBuilder) services.AccessCodeListFetcher {
	return &accessCodeListFetcher{builder}
}
