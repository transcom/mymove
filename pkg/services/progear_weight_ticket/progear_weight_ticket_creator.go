package progearweightticket

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
)

type progearWeightTicketCreator struct {
	checks []progearWeightTicketValidator
}

// NewCustomerProgearWeightTicketCreator creates a new progearWeightTicketCreator struct with the basic checks
func NewCustomerProgearWeightTicketCreator() services.ProgearWeightTicketCreator {
	return &progearWeightTicketCreator{
		checks: basicChecksForCreate(),
	}
}

func (f *progearWeightTicketCreator) CreateProgearWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.ProgearWeightTicket, error) {
	err := validateProgearWeightTicket(appCtx, nil, nil, f.checks...)

	if err != nil {
		return nil, err
	}

	// TODO: Ideally this service would be passed in as a dependency to the `NewMovingExpenseCreator` function.
	//  Our docs have an example, though instead of using the dependency in the service function, it is being used in
	//  the check functions, but the idea is similar:
	//  https://transcom.github.io/mymove-docs/docs/backend/guides/service-objects/implementation#creating-an-instance-of-our-service-object
	shipmentFetcher := ppmshipment.NewPPMShipmentFetcher()

	// This serves as a way of ensuring that the PPM shipment exists. It also ensures a shipment belongs to the logged
	//  in user, for customer app requests.
	ppmShipment, ppmShipmentErr := shipmentFetcher.GetPPMShipment(appCtx, ppmShipmentID, nil, nil)

	if ppmShipmentErr != nil {
		return nil, ppmShipmentErr
	}

	var progearWeightTicket models.ProgearWeightTicket

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		var document models.Document

		// If this is an office user request, get the service member ID from PPM Shipment instead
		if appCtx.Session().IsOfficeApp() {
			serviceMember, err := models.GetCustomerFromPPMShipment(appCtx.DB(), ppmShipmentID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return apperror.NewNotFoundError(ppmShipmentID, "PPM Shipment not found")
				default:
					return apperror.NewQueryError("PPM Shipment", err, "")
				}
			}
			document = models.Document{
				ServiceMemberID: serviceMember.ID,
			}
		} else {
			document = models.Document{
				ServiceMemberID: appCtx.Session().ServiceMemberID,
			}
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&document)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the Document.")
		} else if err != nil {
			return apperror.NewQueryError("Document for ProgearWeightTicket", err, "")
		}

		progearWeightTicket = models.ProgearWeightTicket{
			Document:      document,
			DocumentID:    document.ID,
			PPMShipmentID: ppmShipment.ID,
		}
		verrs, err = txnCtx.DB().ValidateAndCreate(&progearWeightTicket)

		// Check validation errors.
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the ProgearWeightTicket.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("ProgearWeightTicket", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &progearWeightTicket, nil
}
