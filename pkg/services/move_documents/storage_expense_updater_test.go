package movedocument

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MoveDocumentServiceSuite) TestStorageExpenseUpdate() {
	stu := StorageExpenseUpdater{suite.DB(), moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// When: there is a move and move document
	origStartDate := time.Date(2019, 05, 12, 0, 0, 0, 0, time.UTC)
	origEndDate := time.Date(2019, 05, 15, 0, 0, 0, 0, time.UTC)
	daysInStorage := int64(2)
	totalSitCost := unit.Cents(1000)
	ppm := testdatagen.MakePPM(suite.DB(),
		testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				DaysInStorage: &daysInStorage,
				TotalSITCost:  &totalSitCost,
			}})
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	moveDocument := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeEXPENSE,
				Status:                   models.MoveDocumentStatusOK,
			},
			Document: models.Document{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
			},
		})
	storageExpense := models.MovingExpenseDocument{
		MoveDocumentID:       moveDocument.ID,
		MoveDocument:         moveDocument,
		MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: totalSitCost,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       false,
		StorageStartDate:     &origStartDate,
		StorageEndDate:       &origEndDate,
	}
	verrs, err := suite.DB().ValidateAndCreate(&storageExpense)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	startDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2019, 05, 16, 0, 0, 0, 0, time.UTC)
	requestedAmount := int64(2000)
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocument.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		Status:               internalmessages.MoveDocumentStatusOK,
		MoveDocumentType:     internalmessages.MoveDocumentTypeEXPENSE,
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: requestedAmount,
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(endDate),
		StorageStartDate:     handlers.FmtDate(startDate),
	}
	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}
	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)
	umd, verrs, err := stu.Update(updateMoveDocParams, originalMoveDocument, session)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)

	suite.Require().Equal(md.ID.String(), moveDocument.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.MovingExpenseDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	mdEndDate := *md.MovingExpenseDocument.StorageEndDate
	suite.Require().Equal(endDate.UTC(), mdEndDate.UTC())
	mdStartDate := *md.MovingExpenseDocument.StorageStartDate
	suite.Require().Equal(startDate.UTC(), mdStartDate.UTC())
	suite.Require().Equal(unit.Cents(requestedAmount), md.MovingExpenseDocument.RequestedAmountCents)
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(unit.Cents(requestedAmount), *updatedPpm.TotalSITCost)
	suite.Require().Equal(int64(4), *updatedPpm.DaysInStorage)
}

func (suite *MoveDocumentServiceSuite) TestStorageCostAndDaysRemovedWhenNotOK() {
	stu := StorageExpenseUpdater{suite.DB(), moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// When: there is a move and move document
	origStartDate := time.Date(2019, 05, 12, 0, 0, 0, 0, time.UTC)
	origEndDate := time.Date(2019, 05, 15, 0, 0, 0, 0, time.UTC)
	daysInStorage := int64(2)
	totalSitCost := unit.Cents(1000)
	ppm := testdatagen.MakePPM(suite.DB(),
		testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				DaysInStorage: &daysInStorage,
				TotalSITCost:  &totalSitCost,
			}})
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	moveDocument := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeEXPENSE,
				Status:                   models.MoveDocumentStatusOK,
			},
			Document: models.Document{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
			},
		})
	storageExpense := models.MovingExpenseDocument{
		MoveDocumentID:       moveDocument.ID,
		MoveDocument:         moveDocument,
		MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: totalSitCost,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       false,
		StorageStartDate:     &origStartDate,
		StorageEndDate:       &origEndDate,
	}
	verrs, err := suite.DB().ValidateAndCreate(&storageExpense)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	startDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2019, 05, 16, 0, 0, 0, 0, time.UTC)
	requestedAmount := int64(2000)
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocument.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		MoveDocumentType:     internalmessages.MoveDocumentTypeEXPENSE,
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		Status:               internalmessages.MoveDocumentStatusHASISSUE,
		RequestedAmountCents: requestedAmount,
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(endDate),
		StorageStartDate:     handlers.FmtDate(startDate),
	}
	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}
	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)
	umd, verrs, err := stu.Update(updateMoveDocParams, originalMoveDocument, session)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)

	suite.Require().Equal(md.ID.String(), moveDocument.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.MovingExpenseDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	mdEndDate := *md.MovingExpenseDocument.StorageEndDate
	suite.Require().Equal(endDate.UTC(), mdEndDate.UTC())
	mdStartDate := *md.MovingExpenseDocument.StorageStartDate
	suite.Require().Equal(startDate.UTC(), mdStartDate.UTC())
	suite.Require().Equal(unit.Cents(requestedAmount), md.MovingExpenseDocument.RequestedAmountCents)
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(unit.Cents(0), *updatedPpm.TotalSITCost)
	suite.Require().Equal(int64(0), *updatedPpm.DaysInStorage)
}

func (suite *MoveDocumentServiceSuite) TestStorageCostAndDaysAfterManualOverride() {
	stu := StorageExpenseUpdater{suite.DB(), moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// When: there is a move and move document
	origStartDate := time.Date(2019, 05, 12, 0, 0, 0, 0, time.UTC)
	origEndDate := time.Date(2019, 05, 15, 0, 0, 0, 0, time.UTC)
	// made up number as if office user overrides
	daysInStorage := int64(4)
	totalSitCost := unit.Cents(2000)
	ppm := testdatagen.MakePPM(suite.DB(),
		testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				DaysInStorage: &daysInStorage,
				TotalSITCost:  &totalSitCost,
			}})
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	moveDocument := testdatagen.MakeMoveDocument(suite.DB(),
		testdatagen.Assertions{
			MoveDocument: models.MoveDocument{
				MoveID:                   move.ID,
				Move:                     move,
				PersonallyProcuredMoveID: &ppm.ID,
				MoveDocumentType:         models.MoveDocumentTypeEXPENSE,
				Status:                   models.MoveDocumentStatusOK,
			},
			Document: models.Document{
				ServiceMemberID: sm.ID,
				ServiceMember:   sm,
			},
		})
	storageExpense := models.MovingExpenseDocument{
		MoveDocumentID:       moveDocument.ID,
		MoveDocument:         moveDocument,
		MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: totalSitCost,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       false,
		StorageStartDate:     &origStartDate,
		StorageEndDate:       &origEndDate,
	}
	verrs, err := suite.DB().ValidateAndCreate(&storageExpense)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	startDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2019, 05, 16, 0, 0, 0, 0, time.UTC)
	requestedAmount := int64(2000)
	updateMoveDocPayload := internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocument.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		Status:               internalmessages.MoveDocumentStatusOK,
		MoveDocumentType:     internalmessages.MoveDocumentTypeEXPENSE,
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: requestedAmount,
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(endDate),
		StorageStartDate:     handlers.FmtDate(startDate),
	}
	updateMoveDocParams := movedocop.UpdateMoveDocumentParams{
		UpdateMoveDocument: &updateMoveDocPayload,
		MoveDocumentID:     strfmt.UUID(moveDocument.ID.String()),
	}
	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)
	umd, verrs, err := stu.Update(updateMoveDocParams, originalMoveDocument, session)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)

	suite.Require().Equal(md.ID.String(), moveDocument.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.MovingExpenseDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	mdEndDate := *md.MovingExpenseDocument.StorageEndDate
	suite.Require().Equal(endDate.UTC(), mdEndDate.UTC())
	mdStartDate := *md.MovingExpenseDocument.StorageStartDate
	suite.Require().Equal(startDate.UTC(), mdStartDate.UTC())
	suite.Require().Equal(unit.Cents(requestedAmount), md.MovingExpenseDocument.RequestedAmountCents)
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(unit.Cents(2000), *updatedPpm.TotalSITCost)
	suite.Require().Equal(int64(4), *updatedPpm.DaysInStorage)
}
