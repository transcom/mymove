package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type weightTicketCreator struct {
	checks []weightTicketValidator
}

// NewCustomerWeightTicketCreator creates a new weightTicketCreator struct with the basic checks
func NewCustomerWeightTicketCreator() services.WeightTicketCreator {
	return &weightTicketCreator{
		checks: basicChecks(),
	}
}

func (f *weightTicketCreator) CreateWeightTicket(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.WeightTicket, error) {
	err := validateWeightTicket(appCtx, nil, nil, f.checks...)

	if err != nil {
		return nil, err
	}

	var weightTicket models.WeightTicket

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		emptyDocument, err := f.addNewDocument(txnCtx)
		if err != nil {
			return err
		}

		fullDocument, err := f.addNewDocument(txnCtx)
		if err != nil {
			return err
		}

		proofOfOwnership, err := f.addNewDocument(txnCtx)
		if err != nil {
			return err
		}

		weightTicket = models.WeightTicket{
			EmptyDocumentID:                   emptyDocument.ID,
			FullDocumentID:                    fullDocument.ID,
			ProofOfTrailerOwnershipDocumentID: proofOfOwnership.ID,
			PPMShipmentID:                     ppmShipmentID,
		}
		verrs, err := txnCtx.DB().ValidateAndCreate(&weightTicket)

		// Check validation errors.
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the WeightTicket.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("WeightTicket", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &weightTicket, nil
}

func (f *weightTicketCreator) addNewDocument(appCtx appcontext.AppContext) (*models.Document, error) {
	document := &models.Document{
		ServiceMemberID: appCtx.Session().ServiceMemberID,
	}
	verrs, err := appCtx.DB().ValidateAndCreate(document)

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the Document.")
	} else if err != nil {
		return nil, apperror.NewQueryError("Document for WeightTicket", err, "")
	}

	return document, nil
}
