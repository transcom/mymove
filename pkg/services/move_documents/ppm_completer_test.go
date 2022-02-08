package movedocument

//
//import (
//	"github.com/transcom/mymove/pkg/auth"
//	"github.com/transcom/mymove/pkg/gen/internalmessages"
//	"github.com/transcom/mymove/pkg/handlers"
//	"github.com/transcom/mymove/pkg/models"
//	"github.com/transcom/mymove/pkg/testdatagen"
//)
//
//func (suite *MoveDocumentServiceSuite) TestPPMCompleteWhenSSWOK() {
//	ppmc := PPMCompleter{moveDocumentStatusUpdater{}}
//	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
//	session := &auth.Session{
//		ApplicationName: auth.OfficeApp,
//		UserID:          *officeUser.UserID,
//		OfficeUserID:    officeUser.ID,
//	}
//
//	// When: there is a move and move document
//	ppm := testdatagen.MakePPM(suite.DB(),
//		testdatagen.Assertions{
//			PersonallyProcuredMove: models.PersonallyProcuredMove{
//				Status: models.PPMStatusPAYMENTREQUESTED,
//			}})
//	move := ppm.Move
//	sm := ppm.Move.Orders.ServiceMember
//	moveDocument := testdatagen.MakeMoveDocument(suite.DB(),
//		testdatagen.Assertions{
//			MoveDocument: models.MoveDocument{
//				MoveID:                   move.ID,
//				Move:                     move,
//				PersonallyProcuredMoveID: &ppm.ID,
//				MoveDocumentType:         models.MoveDocumentTypeSHIPMENTSUMMARY,
//				Status:                   models.MoveDocumentStatusAWAITINGREVIEW,
//			},
//			Document: models.Document{
//				ServiceMemberID: sm.ID,
//				ServiceMember:   sm,
//			},
//		})
//	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
//		ID:               handlers.FmtUUID(moveDocument.ID),
//		MoveID:           handlers.FmtUUID(move.ID),
//		Title:            handlers.FmtString("super_awesome.pdf"),
//		Notes:            handlers.FmtString("This document is super awesome."),
//		Status:           internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
//		MoveDocumentType: internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeSHIPMENTSUMMARY),
//	}
//
//	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
//	suite.Nil(err)
//	umd, verrs, err := ppmc.Update(suite.AppContextWithSessionForTest(session), updateMoveDocPayload, originalMoveDocument)
//	suite.NotNil(umd)
//	suite.NoVerrs(verrs)
//	suite.Nil(err)
//	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
//	suite.Nil(err)
//
//	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
//	suite.Require().Equal("super_awesome.pdf", md.Title)
//	suite.Require().Equal("This document is super awesome.", *md.Notes)
//	updatedPpm := models.PersonallyProcuredMove{}
//	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
//	suite.Require().Nil(err)
//	suite.Require().Equal(models.PPMStatusCOMPLETED, updatedPpm.Status)
//}
//
//func (suite *MoveDocumentServiceSuite) TestPPMNothingHappensWhenPPMAlreadyCompleted() {
//	ppmc := PPMCompleter{moveDocumentStatusUpdater{}}
//	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
//	session := &auth.Session{
//		ApplicationName: auth.OfficeApp,
//		UserID:          *officeUser.UserID,
//		OfficeUserID:    officeUser.ID,
//	}
//
//	ppm := testdatagen.MakePPM(suite.DB(),
//		testdatagen.Assertions{
//			PersonallyProcuredMove: models.PersonallyProcuredMove{
//				Status: models.PPMStatusCOMPLETED,
//			}})
//	move := ppm.Move
//	sm := ppm.Move.Orders.ServiceMember
//	moveDocument := testdatagen.MakeMoveDocument(suite.DB(),
//		testdatagen.Assertions{
//			MoveDocument: models.MoveDocument{
//				MoveID:                   move.ID,
//				Move:                     move,
//				PersonallyProcuredMoveID: &ppm.ID,
//				MoveDocumentType:         models.MoveDocumentTypeSHIPMENTSUMMARY,
//				Status:                   models.MoveDocumentStatusHASISSUE,
//			},
//			Document: models.Document{
//				ServiceMemberID: sm.ID,
//				ServiceMember:   sm,
//			},
//		})
//	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
//		ID:               handlers.FmtUUID(moveDocument.ID),
//		MoveID:           handlers.FmtUUID(move.ID),
//		Title:            handlers.FmtString("super_awesome.pdf"),
//		Notes:            handlers.FmtString("This document is super awesome."),
//		MoveDocumentType: internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeSHIPMENTSUMMARY),
//		Status:           internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusHASISSUE),
//	}
//
//	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
//	suite.Nil(err)
//	_, verrs, err := ppmc.Update(suite.AppContextWithSessionForTest(session), updateMoveDocPayload, originalMoveDocument)
//	suite.Nil(err)
//	suite.NoVerrs(verrs)
//	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
//	suite.Nil(err)
//
//	suite.Require().Equal(md.ID.String(), moveDocument.ID.String(), "expected move doc ids to match")
//	suite.Require().Equal(md.Title, "super_awesome.pdf")
//	suite.Require().Equal(*md.Notes, "This document is super awesome.")
//	updatedPpm := models.PersonallyProcuredMove{}
//	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
//	suite.Require().Nil(err)
//	suite.Require().Equal(models.PPMStatusCOMPLETED, updatedPpm.Status)
//}
