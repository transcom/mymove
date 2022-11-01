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
	newProgearWeightTicket := &models.ProgearWeightTicket{
		PPMShipmentID: ppmShipmentID,
		FullDocument: models.Document{
			ServiceMemberID: appCtx.Session().ServiceMemberID,
		},
		ConstructedWeightDocument: models.Document{
			ServiceMemberID: appCtx.Session().ServiceMemberID,
		},
	}

	err := validateProgearWeightTicket(appCtx, newProgearWeightTicket, nil, f.checks...)

	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().Eager().ValidateAndCreate(newProgearWeightTicket)

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

	return newProgearWeightTicket, nil
}
