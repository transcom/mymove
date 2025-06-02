package gunsafeweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
)

type gunSafeWeightTicketCreator struct {
	checks []gunSafeWeightTicketValidator
}

// NewOfficeGunSafeWeightTicketCreator creates a new gunSafeWeightTicketCreator struct with the basic checks
func NewOfficeGunSafeWeightTicketCreator() services.GunSafeWeightTicketCreator {
	return &gunSafeWeightTicketCreator{
		checks: basicChecksForCreate(),
	}
}

// NewCustomerGunSafeWeightTicketCreator creates a new gunSafeWeightTicketCreator struct with the basic checks
func NewCustomerGunSafeWeightTicketCreator() services.GunSafeWeightTicketCreator {
	return &gunSafeWeightTicketCreator{
		checks: basicChecksForCreate(),
	}
}

func (f *gunSafeWeightTicketCreator) CreateGunSafeWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.GunSafeWeightTicket, error) {
	err := validateGunSafeWeightTicket(appCtx, nil, nil, f.checks...)

	if err != nil {
		return nil, err
	}

	shipmentFetcher := ppmshipment.NewPPMShipmentFetcher()

	// This serves as a way of ensuring that the PPM shipment exists. It also ensures a shipment belongs to the logged
	//  in user, for customer app requests.
	eagerPreloadAssociations := []string{"Shipment.MoveTaskOrder.Orders.ServiceMember"}
	ppmShipment, ppmShipmentErr := shipmentFetcher.GetPPMShipment(appCtx, ppmShipmentID, eagerPreloadAssociations, nil)

	if ppmShipmentErr != nil {
		return nil, ppmShipmentErr
	}

	var gunSafeWeightTicket models.GunSafeWeightTicket

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {

		document := &models.Document{
			ServiceMemberID: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID,
		}

		verrs, err := appCtx.DB().ValidateAndCreate(document)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the Document.")
		} else if err != nil {
			return apperror.NewQueryError("Document for GunSafeWeightTicket", err, "")
		}

		gunSafeWeightTicket = models.GunSafeWeightTicket{
			Document:      *document,
			DocumentID:    document.ID,
			PPMShipmentID: ppmShipment.ID,
		}
		verrs, err = txnCtx.DB().ValidateAndCreate(&gunSafeWeightTicket)

		// Check validation errors.
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the GunSafeWeightTicket.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("GunSafeWeightTicket", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &gunSafeWeightTicket, nil
}
