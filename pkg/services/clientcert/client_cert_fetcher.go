package clientcert

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type clientCertQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type clientCertFetcher struct {
	builder clientCertQueryBuilder
}

// FetchClientCert fetches an client cert given a slice of filters
func (o *clientCertFetcher) FetchClientCert(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.ClientCert, error) {
	var clientCert models.ClientCert
	error := o.builder.FetchOne(appCtx, &clientCert, filters)
	return clientCert, error
}

// NewClientCertFetcher return an implementation of the ClientCertFetcher interface
func NewClientCertFetcher(builder clientCertQueryBuilder) services.ClientCertFetcher {
	return &clientCertFetcher{builder}
}
