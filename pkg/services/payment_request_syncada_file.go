package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type PaymentRequestSyncadaFileFetcher interface {
	FetchPaymentRequestSyncadaFile(appCtx appcontext.AppContext, filters []QueryFilter) (models.PaymentRequestEdiFile, error)
}
