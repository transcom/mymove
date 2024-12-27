package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type AddressCreator interface {
	CreateAddress(appCtx appcontext.AppContext, address *models.Address) (*models.Address, error)
}

type AddressUpdater interface {
	UpdateAddress(appCtx appcontext.AppContext, address *models.Address, eTag string) (*models.Address, error)
}
