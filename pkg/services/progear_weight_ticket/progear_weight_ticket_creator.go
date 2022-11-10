package progearweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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

	var progearWeightTicket models.ProgearWeightTicket

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {

		document := &models.Document{
			ServiceMemberID: appCtx.Session().ServiceMemberID,
		}

		verrs, err := appCtx.DB().ValidateAndCreate(document)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the Document.")
		} else if err != nil {
			return apperror.NewQueryError("Document for ProgearWeightTicket", err, "")
		}

		progearWeightTicket = models.ProgearWeightTicket{
			Document:      *document,
			DocumentID:    document.ID,
			PPMShipmentID: ppmShipmentID,
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
