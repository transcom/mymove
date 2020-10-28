package movedocument

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
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
	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Require().Nil(err)
	return originalMoveDocument, moveDocument.ID, session
}

func (suite *MoveDocumentServiceSuite) TestMoveDocumentWeightTicketUpdaterWeight() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		Status:           internalmessages.MoveDocumentStatusOK,
		MoveDocumentType: internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
	}

	weightTicketUpdater := mocks.Updater{}
	weightTicketUpdater.On("Update", updateMoveDocPayload, originalMoveDocument, session).
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

	_, verrs, err := mdu.Update(updateMoveDocPayload, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	weightTicketUpdater.AssertNumberOfCalls(suite.T(), "Update", 1)
}

func (suite *MoveDocumentServiceSuite) TestMoveStorageExpenseDocumentUpdater() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		Status:            internalmessages.MoveDocumentStatusOK,
		MoveDocumentType:  internalmessages.MoveDocumentTypeEXPENSE,
		MovingExpenseType: internalmessages.MovingExpenseTypeSTORAGE,
	}

	weightTicketUpdater := mocks.Updater{}
	storageExpenseUpdater := mocks.Updater{}
	storageExpenseUpdater.On("Update", updateMoveDocPayload, originalMoveDocument, session).
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

	_, verrs, err := mdu.Update(updateMoveDocPayload, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	storageExpenseUpdater.AssertNumberOfCalls(suite.T(), "Update", 1)
}

func (suite *MoveDocumentServiceSuite) TestMoveSSWDocumentUpdater() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		MoveDocumentType: internalmessages.MoveDocumentTypeSHIPMENTSUMMARY,
		Status:           internalmessages.MoveDocumentStatusOK,
	}

	weightTicketUpdater := mocks.Updater{}
	storageExpenseUpdater := mocks.Updater{}
	ppmCompleter := mocks.Updater{}
	ppmCompleter.On("Update", updateMoveDocPayload, originalMoveDocument, session).
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

	_, verrs, err := mdu.Update(updateMoveDocPayload, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	ppmCompleter.AssertNumberOfCalls(suite.T(), "Update", 1)
}

func (suite *MoveDocumentServiceSuite) TestMoveGenericDocumentUpdater() {
	originalMoveDocument, moveDocumentID, session := setup(suite)

	// default case that should get called if not storage expense, ssw, or weight ticket set
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		MoveDocumentType: internalmessages.MoveDocumentTypeEXPENSE,
	}

	weightTicketUpdater := mocks.Updater{}
	storageExpenseUpdater := mocks.Updater{}
	ppmCompleter := mocks.Updater{}
	genericUpdater := mocks.Updater{}
	genericUpdater.On("Update", updateMoveDocPayload, originalMoveDocument, session).
		Return(&models.MoveDocument{}, validate.NewErrors(), nil)
	mdu := moveDocumentUpdater{
		db:                        suite.DB(),
		weightTicketUpdater:       &weightTicketUpdater,
		storageExpenseUpdater:     &storageExpenseUpdater,
		ppmCompleter:              &ppmCompleter,
		genericUpdater:            &genericUpdater,
		moveDocumentStatusUpdater: moveDocumentStatusUpdater{},
	}

	_, verrs, err := mdu.Update(updateMoveDocPayload, moveDocumentID, session)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	genericUpdater.AssertNumberOfCalls(suite.T(), "Update", 1)
}
