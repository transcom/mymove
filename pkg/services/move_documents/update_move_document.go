package movedocument

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// createForm is a service object to create a form with data
type moveDocumentUpdater struct {
	db                    *pop.Connection
	weightTicketUpdater   Updater
	storageExpenseUpdater Updater
	ppmCompleter          Updater
	genericUpdater        Updater
	moveDocumentStatusUpdater
}

type Updater interface {
	Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
}

func NewMoveDocumentUpdater(db *pop.Connection) services.MoveDocumentUpdater {
	mdsu := moveDocumentStatusUpdater{}
	wtu := WeightTicketUpdater{db, mdsu}
	su := StorageExpenseUpdater{db, mdsu}
	ppmc := PPMCompleter{db, mdsu}
	gu := GenericUpdater{db, mdsu}
	return &moveDocumentUpdater{
		db:                        db,
		weightTicketUpdater:       wtu,
		storageExpenseUpdater:     su,
		ppmCompleter:              ppmc,
		genericUpdater:            gu,
		moveDocumentStatusUpdater: mdsu,
	}
}

//Update dispatches the various types of move documents to the appropriate Updater
func (m moveDocumentUpdater) Update(params movedocop.UpdateMoveDocumentParams, moveDocId uuid.UUID, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newStatus := models.MoveDocumentStatus(payload.Status)
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	newExpenseType := models.MovingExpenseType(payload.MovingExpenseType)
	originalMoveDocument, err := models.FetchMoveDocument(m.db, session, moveDocId)
	if err != nil {
		return nil, returnVerrs, models.ErrFetchNotFound
	}
	switch {
	case newType == models.MoveDocumentTypeEXPENSE && newExpenseType == models.MovingExpenseTypeSTORAGE:
		return m.storageExpenseUpdater.Update(params, originalMoveDocument, session)
	case newType == models.MoveDocumentTypeWEIGHTTICKETSET:
		return m.weightTicketUpdater.Update(params, originalMoveDocument, session)
	case newType == models.MoveDocumentTypeSHIPMENTSUMMARY && newStatus == models.MoveDocumentStatusOK:
		return m.ppmCompleter.Update(params, originalMoveDocument, session)
	default:
		return m.genericUpdater.Update(params, originalMoveDocument, session)
	}
}

type moveDocumentStatusUpdater struct {
}

//UpdateMoveDocumentStatus attempt to transition a move document from one status to another. Returns and error if
// the status transition is invalid
func (mds moveDocumentStatusUpdater) UpdateMoveDocumentStatus(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newStatus := models.MoveDocumentStatus(payload.Status)
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

func mergeMoveDocuments(moveDocuments models.MoveDocuments, moveDoc models.MoveDocument) models.MoveDocuments {
	var oldMoveDocuments models.MoveDocuments
	for _, v := range moveDocuments {
		if v.ID != moveDoc.ID {
			oldMoveDocuments = append(oldMoveDocuments, v)
		}
		if v.ID == moveDoc.ID && moveDoc.Status == models.MoveDocumentStatusOK {
			oldMoveDocuments = append(oldMoveDocuments, moveDoc)
		}
	}
	return oldMoveDocuments
}


