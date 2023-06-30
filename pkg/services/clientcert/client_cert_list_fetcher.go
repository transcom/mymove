package clientcert

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type clientCertListQueryBuilder interface {
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
}

type clientCertListFetcher struct {
	builder clientCertListQueryBuilder
}

// FetchClientCertList uses the passed query builder to fetch a list of office users
func (o *clientCertListFetcher) FetchClientCertList(appCtx appcontext.AppContext, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) (models.ClientCerts, error) {
	var clientCerts models.ClientCerts
	err := o.builder.FetchMany(appCtx, &clientCerts, filters, associations, pagination, ordering)
	return clientCerts, err
}

// FetchClientCertList uses the passed query builder to fetch a list of office users
func (o *clientCertListFetcher) FetchClientCertCount(appCtx appcontext.AppContext, filters []services.QueryFilter) (int, error) {
	var clientCerts models.ClientCerts
	count, err := o.builder.Count(appCtx, &clientCerts, filters)
	return count, err
}

// NewClientCertListFetcher returns an implementation of ClientCertListFetcher
func NewClientCertListFetcher(builder clientCertListQueryBuilder) services.ClientCertListFetcher {
	return &clientCertListFetcher{builder}
}
