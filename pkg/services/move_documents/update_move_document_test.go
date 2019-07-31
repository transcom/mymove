package movedocument

import (
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/move_documents/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func setup(suite *MoveDocumentServiceSuite) (*models.MoveDocument, uuid.UUID, *auth.Session) {
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	session := &auth.Session{}
	moveDocument := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
				Status:                   models.MoveDocumentStatusOK,
			},
			Document: models.Document{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
			},
		})
	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Require().Nil(err)
	return originalMoveDocument, moveDocument.ID, session
}

func (suite *MoveDocumentServiceSuite) TestMoveDocumentWeightTicketUpdaterWeight() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		Status:           internalmessages.MoveDocumentStatusOK,
		MoveDocumentType: internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
	}
	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocumentID.String()),
	}

	weightTicketUpdater := mocks.Updater{}
	weightTicketUpdater.On("Update", updateMoveDocParams, originalMoveDocument, session).
		Return(&models.MoveDocument{}, validate.NewErrors(), nil)
	storageExpenseUpdater := mocks.Updater{}
	ppmCompleter := mocks.Updater{}
	genericUpdater := mocks.Updater{}
	mdu := moveDocumentUpdater{
		db:                        suite.DB(),
		weightTicketUpdater:       &weightTicketUpdater,
		storageExpenseUpdater:     &storageExpenseUpdater,
		ppmCompleter:              &ppmCompleter,
		genericUpdater:            &genericUpdater,
		moveDocumentStatusUpdater: moveDocumentStatusUpdater{},
	}

	_, verrs, err := mdu.Update(updateMoveDocParams, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	weightTicketUpdater.AssertNumberOfCalls(suite.T(), "Update", 1)
}

func (suite *MoveDocumentServiceSuite) TestMoveStorageExpenseDocumentUpdater() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		Status:            internalmessages.MoveDocumentStatusOK,
		MoveDocumentType:  internalmessages.MoveDocumentTypeEXPENSE,
		MovingExpenseType: internalmessages.MovingExpenseTypeSTORAGE,
	}
	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocumentID.String()),
	}

	weightTicketUpdater := mocks.Updater{}
	storageExpenseUpdater := mocks.Updater{}
	storageExpenseUpdater.On("Update", updateMoveDocParams, originalMoveDocument, session).
		Return(&models.MoveDocument{}, validate.NewErrors(), nil)
	ppmCompleter := mocks.Updater{}
	genericUpdater := mocks.Updater{}
	mdu := moveDocumentUpdater{
		db:                        suite.DB(),
		weightTicketUpdater:       &weightTicketUpdater,
		storageExpenseUpdater:     &storageExpenseUpdater,
		ppmCompleter:              &ppmCompleter,
		genericUpdater:            &genericUpdater,
		moveDocumentStatusUpdater: moveDocumentStatusUpdater{},
	}

	_, verrs, err := mdu.Update(updateMoveDocParams, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	storageExpenseUpdater.AssertNumberOfCalls(suite.T(), "Update", 1)
}

func (suite *MoveDocumentServiceSuite) TestMoveSSWDocumentUpdater() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		MoveDocumentType: internalmessages.MoveDocumentTypeSHIPMENTSUMMARY,
		Status:           internalmessages.MoveDocumentStatusOK,
	}
	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocumentID.String()),
	}

	weightTicketUpdater := mocks.Updater{}
	storageExpenseUpdater := mocks.Updater{}
	ppmCompleter := mocks.Updater{}
	ppmCompleter.On("Update", updateMoveDocParams, originalMoveDocument, session).
		Return(&models.MoveDocument{}, validate.NewErrors(), nil)
	genericUpdater := mocks.Updater{}
	mdu := moveDocumentUpdater{
		db:                        suite.DB(),
		weightTicketUpdater:       &weightTicketUpdater,
		storageExpenseUpdater:     &storageExpenseUpdater,
		ppmCompleter:              &ppmCompleter,
		genericUpdater:            &genericUpdater,
		moveDocumentStatusUpdater: moveDocumentStatusUpdater{},
	}

	_, verrs, err := mdu.Update(updateMoveDocParams, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	ppmCompleter.AssertNumberOfCalls(suite.T(), "Update", 1)
}

func (suite *MoveDocumentServiceSuite) TestMoveGenericDocumentUpdater() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	// default case that should get called if not storage expense, ssw, or weight ticket set
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		MoveDocumentType: internalmessages.MoveDocumentTypeEXPENSE,
	}
	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocumentID.String()),
	}

	weightTicketUpdater := mocks.Updater{}
	storageExpenseUpdater := mocks.Updater{}
	ppmCompleter := mocks.Updater{}
	genericUpdater := mocks.Updater{}
	genericUpdater.On("Update", updateMoveDocParams, originalMoveDocument, session).
		Return(&models.MoveDocument{}, validate.NewErrors(), nil)
	mdu := moveDocumentUpdater{
		db:                        suite.DB(),
		weightTicketUpdater:       &weightTicketUpdater,
		storageExpenseUpdater:     &storageExpenseUpdater,
		ppmCompleter:              &ppmCompleter,
		genericUpdater:            &genericUpdater,
		moveDocumentStatusUpdater: moveDocumentStatusUpdater{},
	}

	_, verrs, err := mdu.Update(updateMoveDocParams, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	genericUpdater.AssertNumberOfCalls(suite.T(), "Update", 1)
}