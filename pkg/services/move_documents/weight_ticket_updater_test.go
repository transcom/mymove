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

func (suite *MoveDocumentServiceSuite) TestNetWeightUpdate() {
	wtu := WeightTicketUpdater{suite.DB(), moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// When: there is a move and move document
	emptyWeight1 := unit.Pound(1000)
	fullWeight1 := unit.Pound(2500)
	netWeight1 := fullWeight1 - emptyWeight1
	wtDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{PersonallyProcuredMove: models.PersonallyProcuredMove{NetWeight: &netWeight1}})
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
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
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocument.ID,
		MoveDocument:             moveDocument,
		EmptyWeight:              &emptyWeight1,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight1,
		FullWeightTicketMissing:  false,
		VehicleNickname:          "My Car",
		VehicleOptions:           "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	emptyWeight := (int64)(200)
	fullWeight := (int64)(500)
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		MoveID:           handlers.FmtUUID(move.ID),
		Title:            handlers.FmtString("super_awesome.pdf"),
		Notes:            handlers.FmtString("This document is super awesome."),
		Status:           internalmessages.MoveDocumentStatusOK,
		VehicleNickname:  "My Car",
		VehicleOptions:   "CAR",
		MoveDocumentType: internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
		EmptyWeight:      &emptyWeight,
		FullWeight:       &fullWeight,
		WeightTicketDate: handlers.FmtDate(wtDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(updateMoveDocPayload, originalMoveDocument, session)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)

	suite.Require().NotNil(md.WeightTicketSetDocument)
	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	suite.Require().Equal(models.MoveDocumentStatusOK, md.Status)
	suite.Require().Equal("My Car", md.WeightTicketSetDocument.VehicleNickname)
	suite.Require().Equal("CAR", md.WeightTicketSetDocument.VehicleOptions)
	suite.Require().Equal(unit.Pound(200), *md.WeightTicketSetDocument.EmptyWeight)
	suite.Require().Equal(unit.Pound(500), *md.WeightTicketSetDocument.FullWeight)
	actualWtDate := *md.WeightTicketSetDocument.WeightTicketDate
	suite.Require().Equal(wtDate.UTC(), actualWtDate.UTC())
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(*updateMoveDocPayload.FullWeight-*updateMoveDocPayload.EmptyWeight, int64(*updatedPpm.NetWeight))

}

func (suite *MoveDocumentServiceSuite) TestNetWeightRemovedWhenStatusNotOK() {
	wtu := WeightTicketUpdater{suite.DB(), moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	emptyWeight1 := unit.Pound(1000)
	fullWeight1 := unit.Pound(2500)
	netWeight1 := fullWeight1 - emptyWeight1
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{PersonallyProcuredMove: models.PersonallyProcuredMove{NetWeight: &netWeight1}})
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
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
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocument.ID,
		MoveDocument:             moveDocument,
		EmptyWeight:              &emptyWeight1,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight1,
		FullWeightTicketMissing:  false,
		VehicleNickname:          "My Car",
		VehicleOptions:           "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	emptyWeight := (int64)(200)
	fullWeight := (int64)(500)
	wtDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		EmptyWeight:      &emptyWeight,
		FullWeight:       &fullWeight,
		ID:               handlers.FmtUUID(moveDocument.ID),
		MoveDocumentType: internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
		MoveID:           handlers.FmtUUID(move.ID),
		Notes:            handlers.FmtString("This document is super awesome."),
		Status:           internalmessages.MoveDocumentStatusHASISSUE,
		Title:            handlers.FmtString("super_awesome.pdf"),
		VehicleNickname:  "My Car",
		VehicleOptions:   "CAR",
		WeightTicketDate: handlers.FmtDate(wtDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(updateMoveDocPayload, originalMoveDocument, session)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)

	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.WeightTicketSetDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	suite.Require().Equal(models.MoveDocumentStatusHASISSUE, md.Status)
	suite.Require().Equal("My Car", md.WeightTicketSetDocument.VehicleNickname)
	suite.Require().Equal("CAR", md.WeightTicketSetDocument.VehicleOptions)
	suite.Require().Equal(unit.Pound(200), *md.WeightTicketSetDocument.EmptyWeight)
	suite.Require().Equal(unit.Pound(500), *md.WeightTicketSetDocument.FullWeight)
	actualWtDate := *md.WeightTicketSetDocument.WeightTicketDate
	suite.Require().Equal(wtDate.UTC(), actualWtDate.UTC())
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(unit.Pound(0), *updatedPpm.NetWeight)

}

func (suite *MoveDocumentServiceSuite) TestNetWeightAfterManualOverride() {
	wtu := WeightTicketUpdater{suite.DB(), moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	emptyWeight1 := unit.Pound(1000)
	fullWeight1 := unit.Pound(2500)
	// made up net weight (as if office user overrode weight tickets)
	netWeight1 := unit.Pound(10000)
	wtDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{PersonallyProcuredMove: models.PersonallyProcuredMove{NetWeight: &netWeight1}})
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
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
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocument.ID,
		MoveDocument:             moveDocument,
		EmptyWeight:              &emptyWeight1,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight1,
		FullWeightTicketMissing:  false,
		VehicleNickname:          "My Car",
		VehicleOptions:           "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	emptyWeight := (int64)(200)
	fullWeight := (int64)(500)
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:               handlers.FmtUUID(moveDocument.ID),
		MoveID:           handlers.FmtUUID(move.ID),
		Title:            handlers.FmtString("super_awesome.pdf"),
		Notes:            handlers.FmtString("This document is super awesome."),
		Status:           internalmessages.MoveDocumentStatusOK,
		VehicleNickname:  "My Car",
		VehicleOptions:   "CAR",
		MoveDocumentType: internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
		EmptyWeight:      &emptyWeight,
		FullWeight:       &fullWeight,
		WeightTicketDate: handlers.FmtDate(wtDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(updateMoveDocPayload, originalMoveDocument, session)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID)
	suite.Nil(err)

	suite.Require().NotNil(md.WeightTicketSetDocument)
	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	suite.Require().Equal(models.MoveDocumentStatusOK, md.Status)
	suite.Require().Equal("My Car", md.WeightTicketSetDocument.VehicleNickname)
	suite.Require().Equal("CAR", md.WeightTicketSetDocument.VehicleOptions)
	suite.Require().Equal(unit.Pound(200), *md.WeightTicketSetDocument.EmptyWeight)
	suite.Require().Equal(unit.Pound(500), *md.WeightTicketSetDocument.FullWeight)
	actualWtDate := *md.WeightTicketSetDocument.WeightTicketDate
	suite.Require().Equal(wtDate.UTC(), actualWtDate.UTC())
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(*updateMoveDocPayload.FullWeight-*updateMoveDocPayload.EmptyWeight, int64(*updatedPpm.NetWeight))

}
