package paymentrequest

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//go:generate mockery --name PaymentRequestHelper --disable-version-string
type Helper interface {
	FetchServiceParamList(appCtx appcontext.AppContext, mtoServiceItem models.MTOServiceItem) (models.ServiceParams, error)
	FetchDistinctSystemServiceParamList(appCtx appcontext.AppContext, mtoServiceItems []models.MTOServiceItem) ([]models.ServiceItemParamKey, error)
}

// RequestPaymentHelper is a helper to connect to the DB
type RequestPaymentHelper struct{}
