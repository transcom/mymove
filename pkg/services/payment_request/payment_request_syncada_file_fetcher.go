package paymentrequest

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentReqeustSyncadaFileQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
}

type paymentRequestSyncadaFileFetcher struct {
	builder paymentReqeustSyncadaFileQueryBuilder
}

func (p *paymentRequestSyncadaFileFetcher) FetchPaymentRequestSyncadaFile(appCtx appcontext.AppContext, filters []services.QueryFilter) (models.PaymentRequestEdiFile, error) {
	var paymentRequestEdiFile models.PaymentRequestEdiFile
	err := p.builder.FetchOne(appCtx, &paymentRequestEdiFile, filters)
	return paymentRequestEdiFile, err
}

func NewPaymentRequestSyncadaFileFetcher(builder paymentReqeustSyncadaFileQueryBuilder) services.PaymentRequestSyncadaFileFetcher {
	return &paymentRequestSyncadaFileFetcher{builder}
}
