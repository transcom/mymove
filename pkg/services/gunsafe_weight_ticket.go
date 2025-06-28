package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// GunSafeWeightTickerCreator creates a GunSafeWeightTicket that is associated with a PPMShipment
//
//go:generate mockery --name GunSafeWeightTicketCreator
type GunSafeWeightTicketCreator interface {
	CreateGunSafeWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.GunSafeWeightTicket, error)
}

// GunSafeWeightTicketUpdater updates a GunSafeWeightTicket
//
//go:generate mockery --name GunSafeWeightTicketUpdater
type GunSafeWeightTicketUpdater interface {
	UpdateGunSafeWeightTicket(appCtx appcontext.AppContext, gunsafeWeightTicket models.GunSafeWeightTicket, eTag string) (*models.GunSafeWeightTicket, error)
}

// GunSafeWeightTicketDeleter deletes a GunSafeWeightTicket
//
//go:generate mockery --name GunSafeWeightTicketDeleter
type GunSafeWeightTicketDeleter interface {
	DeleteGunSafeWeightTicket(appCtx appcontext.AppContext, ppmID uuid.UUID, gunsafeWeightTicketID uuid.UUID) error
}
