package movedocument

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// moveDocumentUpdater implementation of MoveDocumentUpdater
type moveDocumentUpdater struct {
	weightTicketUpdater   Updater
	storageExpenseUpdater Updater
	ppmCompleter          Updater
	genericUpdater        Updater
	moveDocumentStatusUpdater
}

//Updater interface for individual document updaters
//go:generate mockery --name Updater --disable-version-string
type Updater interface {
	Update(appCtx appcontext.AppContext, moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error)
}

//NewMoveDocumentUpdater create NewMoveDocumentUpdater including expected updaters
func NewMoveDocumentUpdater() services.MoveDocumentUpdater {
	mdsu := moveDocumentStatusUpdater{}
	wtu := WeightTicketUpdater{mdsu}
	su := StorageExpenseUpdater{mdsu}
	ppmc := PPMCompleter{}
	gu := GenericUpdater{}
	return &moveDocumentUpdater{
		weightTicketUpdater:       wtu,
		storageExpenseUpdater:     su,
		ppmCompleter:              ppmc,
		genericUpdater:            gu,
		moveDocumentStatusUpdater: mdsu,
	}
}

//Update dispatches the various types of move documents to the appropriate Updater
func (m moveDocumentUpdater) Update(appCtx appcontext.AppContext, moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDocID uuid.UUID) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	if moveDocumentPayload.MoveDocumentType == nil {
		return nil, returnVerrs, errors.New("missing required field: MoveDocumentType")
	}
	newType := models.MoveDocumentType(*moveDocumentPayload.MoveDocumentType)
	newExpenseType := models.MovingExpenseType(moveDocumentPayload.MovingExpenseType)
	originalMoveDocument, err := models.FetchMoveDocument(appCtx.DB(), appCtx.Session(), moveDocID, false)
	if err != nil {
		return nil, returnVerrs, models.ErrFetchNotFound
	}
	switch {
	case newType == models.MoveDocumentTypeEXPENSE && newExpenseType == models.MovingExpenseTypeSTORAGE:
		return m.storageExpenseUpdater.Update(appCtx, moveDocumentPayload, originalMoveDocument)
	case newType == models.MoveDocumentTypeWEIGHTTICKETSET:
		return m.weightTicketUpdater.Update(appCtx, moveDocumentPayload, originalMoveDocument)
	case newType == models.MoveDocumentTypeSHIPMENTSUMMARY:
		return m.ppmCompleter.Update(appCtx, moveDocumentPayload, originalMoveDocument)
	default:
		return m.genericUpdater.Update(appCtx, moveDocumentPayload, originalMoveDocument)
	}
}

type moveDocumentStatusUpdater struct {
}

//UpdateMoveDocumentStatus attempt to transition a move document from one status to another.
// Returns and error if the status transition is invalid
func (mds moveDocumentStatusUpdater) UpdateMoveDocumentStatus(appCtx appcontext.AppContext, moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	if moveDocumentPayload.MoveDocumentType == nil {
		return nil, returnVerrs, errors.New("missing required field: MoveDocumentType")
	}
	newStatus := models.MoveDocumentStatus(*moveDocumentPayload.Status)
	if moveDoc == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updateMoveDocumentStatus: missing move document")
	}
	origStatus := moveDoc.Status
	// only modify the move document status if its changed
	if newStatus == origStatus {
		return moveDoc, returnVerrs, nil
	}
	err := moveDoc.AttemptTransition(newStatus)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "updateMoveDocumentStatus: transition error")
	}
	return moveDoc, returnVerrs, nil
}

func mergeMoveDocuments(appCtx appcontext.AppContext, ppmID uuid.UUID, moveDoc *models.MoveDocument, moveDocumentType models.MoveDocumentType, status models.MoveDocumentStatus) (models.MoveDocuments, error) {
	// get all documents excluding new document merge in new updated document if status is correct
	moveDocuments, err := models.FetchMoveDocuments(appCtx.DB(), appCtx.Session(), ppmID, &status, moveDocumentType, false)
	if err != nil {
		return models.MoveDocuments{}, errors.New("mergeDocuments: unable to fetch move documents")
	}
	var mergedMoveDocuments models.MoveDocuments
	for _, moveDocument := range moveDocuments {
		if moveDocument.ID != moveDoc.ID {
			mergedMoveDocuments = append(mergedMoveDocuments, moveDocument)
		}
	}
	if moveDoc.Status == status {
		mergedMoveDocuments = append(mergedMoveDocuments, *moveDoc)
	}
	return mergedMoveDocuments, nil
}
