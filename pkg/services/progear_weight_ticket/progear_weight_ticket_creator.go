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

func NewProgearWeightTicketCreator() services.ProgearWeightTicketCreator {
	return &progearWeightTicketCreator{
		checks: createChecks(),
	}
}

func (f *progearWeightTicketCreator) CreateProgearWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.ProgearWeightTicket, error) {
	err := validateProgearWeightTicket(appCtx, nil, nil, f.checks...)
	if err != nil {
		return nil, err
	}

	var progearWeightTicket models.ProgearWeightTicket

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		document := models.Document{
			ServiceMemberID: appCtx.Session().ServiceMemberID,
		}
		allDocs := models.Documents{document}

		verrs, err := appCtx.DB().ValidateAndCreate(allDocs)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the Document.")
		} else if err != nil {
			return apperror.NewQueryError("Document for ProgearWeightTicket", err, "")
		}

		progearWeightTicket = models.ProgearWeightTicket{
			Document:      allDocs[0],
			DocumentID:    allDocs[0].ID,
			PPMShipmentID: ppmShipmentID,
		}

		verrs, err = txnCtx.DB().ValidateAndCreate(&progearWeightTicket)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "")
		} else if err != nil {
			return apperror.NewQueryError("Progear Weight Ticket", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &progearWeightTicket, nil
}
