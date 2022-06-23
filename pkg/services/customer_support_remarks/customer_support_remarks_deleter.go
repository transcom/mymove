package customersupportremarks

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type customerSupportRemarkDeleter struct {
}

func NewCustomerSupportRemarkDeleter() services.CustomerSupportRemarkDeleter {
	return &customerSupportRemarkDeleter{}
}

func (o customerSupportRemarkDeleter) DeleteCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemarkID uuid.UUID) error {
	return nil
}
