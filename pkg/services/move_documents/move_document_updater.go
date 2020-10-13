package movedocument

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// moveDocumentUpdater implementation of MoveDocumentUpdater
type moveDocumentUpdater struct {
	db                    *pop.Connection
	weightTicketUpdater   Updater
	storageExpenseUpdater Updater
	ppmCompleter          Updater
	genericUpdater        Updater
	moveDocumentStatusUpdater
}

//Updater interface for individual document updaters
//go:generate mockery -name Updater
type Updater interface {
	Update(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
}

//NewMoveDocumentUpdater create NewMoveDocumentUpdater including expected updaters
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
func (m moveDocumentUpdater) Update(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDocID uuid.UUID, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	newType := models.MoveDocumentType(moveDocumentPayload.MoveDocumentType)
	newExpenseType := models.MovingExpenseType(moveDocumentPayload.MovingExpenseType)
	originalMoveDocument, err := models.FetchMoveDocument(m.db, session, moveDocID, false)
	if err != nil {
		return nil, returnVerrs, models.ErrFetchNotFound
	}
	switch {
	case newType == models.MoveDocumentTypeEXPENSE && newExpenseType == models.MovingExpenseTypeSTORAGE:
		return m.storageExpenseUpdater.Update(moveDocumentPayload, originalMoveDocument, session)
	case newType == models.MoveDocumentTypeWEIGHTTICKETSET:
		return m.weightTicketUpdater.Update(moveDocumentPayload, originalMoveDocument, session)
	case newType == models.MoveDocumentTypeSHIPMENTSUMMARY:
		return m.ppmCompleter.Update(moveDocumentPayload, originalMoveDocument, session)
	default:
		return m.genericUpdater.Update(moveDocumentPayload, originalMoveDocument, session)
	}
}

type moveDocumentStatusUpdater struct {
}

//UpdateMoveDocumentStatus attempt to transition a move document from one status to another.
// Returns and error if the status transition is invalid
func (mds moveDocumentStatusUpdater) UpdateMoveDocumentStatus(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	newStatus := models.MoveDocumentStatus(moveDocumentPayload.Status)
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

func mergeMoveDocuments(db *pop.Connection, session *auth.Session, ppmID uuid.UUID, moveDoc *models.MoveDocument, moveDocumentType models.MoveDocumentType, status models.MoveDocumentStatus) (models.MoveDocuments, error) {
	// get all documents excluding new document merge in new updated document if status is correct
	moveDocuments, err := models.FetchMoveDocuments(db, session, ppmID, &status, moveDocumentType, false)
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
