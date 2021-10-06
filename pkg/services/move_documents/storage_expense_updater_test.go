package movedocument

import (
	"time"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MoveDocumentServiceSuite) TestStorageExpenseUpdate() {
	stu := StorageExpenseUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	origStartDate := time.Date(2019, 05, 12, 0, 0, 0, 0, time.UTC)
	origEndDate := time.Date(2019, 05, 15, 0, 0, 0, 0, time.UTC)
	origDaysInStorage := int64(2)
	origTotalSitCost := unit.Cents(1000)
	ppm := testdatagen.MakePPM(suite.DB(),
		testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				DaysInStorage: &origDaysInStorage,
				TotalSITCost:  &origTotalSitCost,
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
		RequestedAmountCents: origTotalSitCost,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       false,
		StorageStartDate:     &origStartDate,
		StorageEndDate:       &origEndDate,
	}
	verrs, err := suite.DB().ValidateAndCreate(&storageExpense)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	newStartDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	newEndDate := time.Date(2019, 05, 16, 0, 0, 0, 0, time.UTC)
	newRequestedAmount := int64(2000)
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocument.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		Status:               internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		MoveDocumentType:     internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeEXPENSE),
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: newRequestedAmount,
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(newEndDate),
		StorageStartDate:     handlers.FmtDate(newStartDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)
	umd, verrs, err := stu.Update(suite.TestAppContext(), updateMoveDocPayload, originalMoveDocument)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)

	suite.Require().Equal(md.ID.String(), moveDocument.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.MovingExpenseDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	mdEndDate := *md.MovingExpenseDocument.StorageEndDate
	suite.Require().Equal(newEndDate.UTC(), mdEndDate.UTC())
	mdStartDate := *md.MovingExpenseDocument.StorageStartDate
	suite.Require().Equal(newStartDate.UTC(), mdStartDate.UTC())
	suite.Require().Equal(unit.Cents(newRequestedAmount), md.MovingExpenseDocument.RequestedAmountCents)
	updatedPpm := models.PersonallyProcuredMove{}
	// ppm is updated to reflect new sit total cost and days in storage
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(unit.Cents(newRequestedAmount), *updatedPpm.TotalSITCost)
	suite.Require().Equal(int64(5), *updatedPpm.DaysInStorage)
}

func (suite *MoveDocumentServiceSuite) TestStorageCostAndDaysRemovedWhenNotOK() {
	stu := StorageExpenseUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

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

	newStartDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	newEndDate := time.Date(2019, 05, 16, 0, 0, 0, 0, time.UTC)
	newRequestedAmount := int64(2000)
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocument.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		MoveDocumentType:     internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeEXPENSE),
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		Status:               internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusHASISSUE),
		RequestedAmountCents: newRequestedAmount,
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(newEndDate),
		StorageStartDate:     handlers.FmtDate(newStartDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)
	umd, verrs, err := stu.Update(suite.TestAppContext(), updateMoveDocPayload, originalMoveDocument)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)

	suite.Require().Equal(md.ID.String(), moveDocument.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.MovingExpenseDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	mdEndDate := *md.MovingExpenseDocument.StorageEndDate
	suite.Require().Equal(newEndDate.UTC(), mdEndDate.UTC())
	mdStartDate := *md.MovingExpenseDocument.StorageStartDate
	suite.Require().Equal(newStartDate.UTC(), mdStartDate.UTC())
	suite.Require().Equal(unit.Cents(newRequestedAmount), md.MovingExpenseDocument.RequestedAmountCents)
	updatedPpm := models.PersonallyProcuredMove{}
	// ppm is updated to reflect exlusion of this sit expense from total cost and days in storage
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(unit.Cents(0), *updatedPpm.TotalSITCost)
	suite.Require().Equal(int64(0), *updatedPpm.DaysInStorage)
}

func (suite *MoveDocumentServiceSuite) TestStorageDaysTotalCostMultipleReceipts() {
	stu := StorageExpenseUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	startDateOne := time.Date(2019, 05, 12, 0, 0, 0, 0, time.UTC)
	endDateOne := time.Date(2019, 05, 15, 0, 0, 0, 0, time.UTC)
	totalDaysOne := int64(endDateOne.Sub(startDateOne).Hours() / 24)
	totalSitCostOne := unit.Cents(2000)
	ppm := testdatagen.MakeDefaultPPM(suite.DB())
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	moveDocumentOne := testdatagen.MakeMoveDocument(suite.DB(),
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
	storageExpenseOne := models.MovingExpenseDocument{
		MoveDocumentID:       moveDocumentOne.ID,
		MoveDocument:         moveDocumentOne,
		MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: totalSitCostOne,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       false,
		StorageStartDate:     &startDateOne,
		StorageEndDate:       &endDateOne,
	}
	verrs, err := suite.DB().ValidateAndCreate(&storageExpenseOne)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	startDateTwo := time.Date(2019, 05, 20, 0, 0, 0, 0, time.UTC)
	endDateTwo := time.Date(2019, 05, 25, 0, 0, 0, 0, time.UTC)
	totalDaysTwo := int64(endDateTwo.Sub(startDateTwo).Hours() / 24)
	totalSitCostTwo := unit.Cents(1000)
	moveDocumentTwo := testdatagen.MakeMoveDocument(suite.DB(),
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
	storageExpenseTwo := models.MovingExpenseDocument{
		MoveDocumentID:       moveDocumentTwo.ID,
		MoveDocument:         moveDocumentTwo,
		MovingExpenseType:    models.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: totalSitCostTwo,
		PaymentMethod:        "GTCC",
		ReceiptMissing:       false,
		StorageStartDate:     &startDateTwo,
		StorageEndDate:       &endDateTwo,
	}
	verrs, err = suite.DB().ValidateAndCreate(&storageExpenseTwo)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	updateMoveDocOnePayload := &internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocumentOne.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		MoveDocumentType:     internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeEXPENSE),
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		Status:               internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		RequestedAmountCents: int64(totalSitCostOne),
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(endDateOne),
		StorageStartDate:     handlers.FmtDate(startDateOne),
	}
	updateMoveDocTwoPayload := &internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocumentTwo.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		MoveDocumentType:     internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeEXPENSE),
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		Status:               internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		RequestedAmountCents: int64(totalSitCostTwo),
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(endDateTwo),
		StorageStartDate:     handlers.FmtDate(startDateTwo),
	}

	originalMoveDocumentOne, err := models.FetchMoveDocument(suite.DB(), session, moveDocumentOne.ID, false)
	suite.Nil(err)
	umd, verrs, err := stu.Update(suite.TestAppContext(), updateMoveDocOnePayload, originalMoveDocumentOne)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)

	originalMoveDocumentTwo, err := models.FetchMoveDocument(suite.DB(), session, moveDocumentTwo.ID, false)
	suite.Nil(err)
	umd, verrs, err = stu.Update(suite.TestAppContext(), updateMoveDocTwoPayload, originalMoveDocumentTwo)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)

	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(totalSitCostOne+totalSitCostTwo, *updatedPpm.TotalSITCost)
	suite.Require().Equal(totalDaysOne+totalDaysTwo, *updatedPpm.DaysInStorage)
}

func (suite *MoveDocumentServiceSuite) TestStorageCostAndDaysAfterManualOverride() {
	stu := StorageExpenseUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	origStartDate := time.Date(2019, 05, 12, 0, 0, 0, 0, time.UTC)
	origEndDate := time.Date(2019, 05, 15, 0, 0, 0, 0, time.UTC)
	// made up daysInStorage and totalSitCost (as if office user overrides)
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

	newStartDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	newEndDate := time.Date(2019, 05, 16, 0, 0, 0, 0, time.UTC)
	newRequestedAmount := int64(2000)
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:                   handlers.FmtUUID(moveDocument.ID),
		MoveID:               handlers.FmtUUID(move.ID),
		Title:                handlers.FmtString("super_awesome.pdf"),
		Notes:                handlers.FmtString("This document is super awesome."),
		Status:               internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		MoveDocumentType:     internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeEXPENSE),
		MovingExpenseType:    internalmessages.MovingExpenseTypeSTORAGE,
		RequestedAmountCents: newRequestedAmount,
		PaymentMethod:        "GTCC",
		StorageEndDate:       handlers.FmtDate(newEndDate),
		StorageStartDate:     handlers.FmtDate(newStartDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)
	umd, verrs, err := stu.Update(suite.TestAppContext(), updateMoveDocPayload, originalMoveDocument)
	suite.NotNil(umd)
	suite.Nil(err)
	suite.NoVerrs(verrs)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)

	suite.Require().Equal(md.ID.String(), moveDocument.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.MovingExpenseDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	mdEndDate := *md.MovingExpenseDocument.StorageEndDate
	suite.Require().Equal(newEndDate.UTC(), mdEndDate.UTC())
	mdStartDate := *md.MovingExpenseDocument.StorageStartDate
	suite.Require().Equal(newStartDate.UTC(), mdStartDate.UTC())
	suite.Require().Equal(unit.Cents(newRequestedAmount), md.MovingExpenseDocument.RequestedAmountCents)
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(unit.Cents(newRequestedAmount), *updatedPpm.TotalSITCost)
	newDaysInStorage := int64(newEndDate.Sub(newStartDate).Hours() / 24)
	suite.Require().Equal(newDaysInStorage, *updatedPpm.DaysInStorage)
}
