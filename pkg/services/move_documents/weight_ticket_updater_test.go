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
	wtu := WeightTicketUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// When: there is a move and move document
	vehicleNickname := "My Car"
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
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      models.WeightTicketSetTypeCAR,
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	vehicleNickname2 := "My Car2"
	emptyWeight := int64(200)
	fullWeight := int64(500)
	weightTicketSetType := internalmessages.WeightTicketSetTypeCAR
	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:                  handlers.FmtUUID(moveDocument.ID),
		MoveID:              handlers.FmtUUID(move.ID),
		Title:               handlers.FmtString("super_awesome.pdf"),
		Notes:               handlers.FmtString("This document is super awesome."),
		Status:              internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		VehicleNickname:     &vehicleNickname2,
		WeightTicketSetType: &weightTicketSetType,
		MoveDocumentType:    internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeWEIGHTTICKETSET),
		EmptyWeight:         &emptyWeight,
		FullWeight:          &fullWeight,
		WeightTicketDate:    handlers.FmtDate(wtDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(suite.AppContextWithSessionForTest(session), updateMoveDocPayload, originalMoveDocument)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)

	suite.Require().NotNil(md.WeightTicketSetDocument)
	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	suite.Require().Equal(models.MoveDocumentStatusOK, md.Status)
	suite.Require().Equal(vehicleNickname2, *md.WeightTicketSetDocument.VehicleNickname)
	suite.Require().Equal(models.WeightTicketSetTypeCAR, md.WeightTicketSetDocument.WeightTicketSetType)
	suite.Require().Equal(unit.Pound(200), *md.WeightTicketSetDocument.EmptyWeight)
	suite.Require().Equal(unit.Pound(500), *md.WeightTicketSetDocument.FullWeight)
	actualWtDate := *md.WeightTicketSetDocument.WeightTicketDate
	suite.Require().Equal(wtDate.UTC(), actualWtDate.UTC())
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(*updateMoveDocPayload.FullWeight-*updateMoveDocPayload.EmptyWeight, int64(*updatedPpm.NetWeight))

}

func (suite *MoveDocumentServiceSuite) TestNetWeightWhenMultipleWeightTickets() {
	wtu := WeightTicketUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	vehicleNickname := "My Car"
	emptyWeight1 := unit.Pound(1000)
	fullWeight1 := unit.Pound(2500)
	netWeight1 := fullWeight1 - emptyWeight1
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{PersonallyProcuredMove: models.PersonallyProcuredMove{NetWeight: &netWeight1}})
	move := ppm.Move
	sm := ppm.Move.Orders.ServiceMember
	moveDocumentOne := testdatagen.MakeMoveDocument(suite.DB(),
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
	weightTicketSetDocumentOne := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocumentOne.ID,
		MoveDocument:             moveDocumentOne,
		EmptyWeight:              &emptyWeight1,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight1,
		FullWeightTicketMissing:  false,
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocumentOne)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	emptyWeight2 := unit.Pound(1000)
	fullWeight2 := unit.Pound(5000)
	netWeight2 := fullWeight2 - emptyWeight2
	wtDateTwo := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	moveDocumentTwo := testdatagen.MakeMoveDocument(suite.DB(),
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

	weightTicketSetDocumentTwo := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocumentTwo.ID,
		MoveDocument:             moveDocumentTwo,
		EmptyWeight:              &emptyWeight2,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight2,
		FullWeightTicketMissing:  false,
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err = suite.DB().ValidateAndCreate(&weightTicketSetDocumentTwo)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	vehicleNickname2 := "My Car2"
	weightTicketSetType := internalmessages.WeightTicketSetTypeCAR

	updateMoveDocOnePayload := &internalmessages.MoveDocumentPayload{
		ID:                  handlers.FmtUUID(moveDocumentOne.ID),
		Status:              internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		Title:               handlers.FmtString("super_awesome.pdf"),
		Notes:               handlers.FmtString("This document is super awesome."),
		VehicleNickname:     &vehicleNickname2,
		WeightTicketSetType: &weightTicketSetType,
		MoveDocumentType:    internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeWEIGHTTICKETSET),
		EmptyWeight:         handlers.FmtInt64((int64)(emptyWeight1)),
		FullWeight:          handlers.FmtInt64((int64)(fullWeight1)),
		WeightTicketDate:    handlers.FmtDate(wtDateTwo),
	}
	updateMoveDocTwoPayload := &internalmessages.MoveDocumentPayload{
		ID:                  handlers.FmtUUID(moveDocumentOne.ID),
		Status:              internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		Title:               handlers.FmtString("super_awesome.pdf"),
		Notes:               handlers.FmtString("This document is super awesome."),
		VehicleNickname:     &vehicleNickname2,
		WeightTicketSetType: &weightTicketSetType,
		MoveDocumentType:    internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeWEIGHTTICKETSET),
		EmptyWeight:         handlers.FmtInt64((int64)(emptyWeight2)),
		FullWeight:          handlers.FmtInt64((int64)(fullWeight2)),
		WeightTicketDate:    handlers.FmtDate(wtDateTwo),
	}
	originalMoveDocumentOne, err := models.FetchMoveDocument(suite.DB(), session, moveDocumentOne.ID, false)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(suite.AppContextWithSessionForTest(session), updateMoveDocOnePayload, originalMoveDocumentOne)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)

	originalMoveDocumentTwo, err := models.FetchMoveDocument(suite.DB(), session, moveDocumentTwo.ID, false)
	suite.Nil(err)
	umd, verrs, err = wtu.Update(suite.AppContextWithSessionForTest(session), updateMoveDocTwoPayload, originalMoveDocumentTwo)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)

	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(netWeight1+netWeight2, *updatedPpm.NetWeight)

}

func (suite *MoveDocumentServiceSuite) TestNetWeightRemovedWhenStatusNotOK() {
	wtu := WeightTicketUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	vehicleNickname := "My Car"
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
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      "CAR",
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	vehicleNickname2 := "My Car2"
	emptyWeight := int64(200)
	fullWeight := int64(500)
	wtDate := time.Date(2019, 05, 11, 0, 0, 0, 0, time.UTC)
	weightTicketSetType := internalmessages.WeightTicketSetTypeCAR

	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		EmptyWeight:         &emptyWeight,
		FullWeight:          &fullWeight,
		ID:                  handlers.FmtUUID(moveDocument.ID),
		MoveDocumentType:    internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeWEIGHTTICKETSET),
		MoveID:              handlers.FmtUUID(move.ID),
		Notes:               handlers.FmtString("This document is super awesome."),
		Status:              internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusHASISSUE),
		Title:               handlers.FmtString("super_awesome.pdf"),
		VehicleNickname:     &vehicleNickname2,
		WeightTicketSetType: &weightTicketSetType,
		WeightTicketDate:    handlers.FmtDate(wtDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(suite.AppContextWithSessionForTest(session), updateMoveDocPayload, originalMoveDocument)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)

	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
	suite.Require().NotNil(md.WeightTicketSetDocument)
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	suite.Require().Equal(models.MoveDocumentStatusHASISSUE, md.Status)
	suite.Require().Equal(vehicleNickname2, *md.WeightTicketSetDocument.VehicleNickname)
	suite.Require().Equal(models.WeightTicketSetTypeCAR, md.WeightTicketSetDocument.WeightTicketSetType)
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
	wtu := WeightTicketUpdater{moveDocumentStatusUpdater{}}
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	vehicleNickname := "My Car"
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
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      models.WeightTicketSetTypeCAR,
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	vehicleNickname2 := "My Car2"
	emptyWeight := (int64)(200)
	fullWeight := (int64)(500)
	weightTicketSetType := internalmessages.WeightTicketSetTypeCAR

	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:                  handlers.FmtUUID(moveDocument.ID),
		MoveID:              handlers.FmtUUID(move.ID),
		Title:               handlers.FmtString("super_awesome.pdf"),
		Notes:               handlers.FmtString("This document is super awesome."),
		Status:              internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		VehicleNickname:     &vehicleNickname2,
		WeightTicketSetType: &weightTicketSetType,
		MoveDocumentType:    internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeWEIGHTTICKETSET),
		EmptyWeight:         &emptyWeight,
		FullWeight:          &fullWeight,
		WeightTicketDate:    handlers.FmtDate(wtDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(suite.AppContextWithSessionForTest(session), updateMoveDocPayload, originalMoveDocument)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)

	suite.Require().NotNil(md.WeightTicketSetDocument)
	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	suite.Require().Equal(models.MoveDocumentStatusOK, md.Status)
	suite.Require().Equal(vehicleNickname2, *md.WeightTicketSetDocument.VehicleNickname)
	suite.Require().Equal(models.WeightTicketSetTypeCAR, md.WeightTicketSetDocument.WeightTicketSetType)
	suite.Require().Equal(unit.Pound(200), *md.WeightTicketSetDocument.EmptyWeight)
	suite.Require().Equal(unit.Pound(500), *md.WeightTicketSetDocument.FullWeight)
	actualWtDate := *md.WeightTicketSetDocument.WeightTicketDate
	suite.Require().Equal(wtDate.UTC(), actualWtDate.UTC())
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(*updateMoveDocPayload.FullWeight-*updateMoveDocPayload.EmptyWeight, int64(*updatedPpm.NetWeight))

}

func (suite *MoveDocumentServiceSuite) TestMakeAndModelUpdate() {
	wtu := WeightTicketUpdater{moveDocumentStatusUpdater{}}
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

	vehicleNickname := "My Box Truck"
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocument.ID,
		MoveDocument:             moveDocument,
		EmptyWeight:              &emptyWeight1,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight1,
		FullWeightTicketMissing:  false,
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      models.WeightTicketSetTypeBOXTRUCK,
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	vehicleMake := "Radio Flyer"
	vehicleModel := "Wagon"
	emptyWeight := (int64)(1000)
	fullWeight := (int64)(2500)
	weightTicketSetType := internalmessages.WeightTicketSetTypeCAR

	updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
		ID:                  handlers.FmtUUID(moveDocument.ID),
		MoveID:              handlers.FmtUUID(move.ID),
		Title:               handlers.FmtString("super_awesome.pdf"),
		Notes:               handlers.FmtString("This document is super awesome."),
		Status:              internalmessages.NewMoveDocumentStatus(internalmessages.MoveDocumentStatusOK),
		VehicleNickname:     &vehicleNickname,
		VehicleMake:         &vehicleMake,
		VehicleModel:        &vehicleModel,
		EmptyWeight:         &emptyWeight,
		FullWeight:          &fullWeight,
		WeightTicketSetType: &weightTicketSetType,
		MoveDocumentType:    internalmessages.NewMoveDocumentType(internalmessages.MoveDocumentTypeWEIGHTTICKETSET),
		WeightTicketDate:    handlers.FmtDate(wtDate),
	}

	originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)
	umd, verrs, err := wtu.Update(suite.AppContextWithSessionForTest(session), updateMoveDocPayload, originalMoveDocument)
	suite.NotNil(umd)
	suite.NoVerrs(verrs)
	suite.Nil(err)
	md, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
	suite.Nil(err)

	suite.Require().NotNil(md.WeightTicketSetDocument)
	suite.Require().Equal(moveDocument.ID.String(), md.ID.String(), "expected move doc ids to match")
	suite.Require().Equal("super_awesome.pdf", md.Title)
	suite.Require().Equal("This document is super awesome.", *md.Notes)
	suite.Require().Equal(models.MoveDocumentStatusOK, md.Status)
	suite.Require().Equal(vehicleNickname, *md.WeightTicketSetDocument.VehicleNickname)
	suite.Require().Equal(vehicleMake, *md.WeightTicketSetDocument.VehicleMake)
	suite.Require().Equal(vehicleModel, *md.WeightTicketSetDocument.VehicleModel)
	suite.Require().Equal(models.WeightTicketSetTypeCAR, md.WeightTicketSetDocument.WeightTicketSetType)
	suite.Require().Equal(unit.Pound(1000), *md.WeightTicketSetDocument.EmptyWeight)
	suite.Require().Equal(unit.Pound(2500), *md.WeightTicketSetDocument.FullWeight)
	actualWtDate := *md.WeightTicketSetDocument.WeightTicketDate
	suite.Require().Equal(wtDate.UTC(), actualWtDate.UTC())
	updatedPpm := models.PersonallyProcuredMove{}
	err = suite.DB().Where(`id = $1`, ppm.ID).First(&updatedPpm)
	suite.Require().Nil(err)
	suite.Require().Equal(*updateMoveDocPayload.FullWeight-*updateMoveDocPayload.EmptyWeight, int64(*updatedPpm.NetWeight))

}

// TODO: Fix now that we capture transacation error
/*
func (suite *MoveDocumentServiceSuite) TestValueForEitherMakeOrModelFails() {
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

	vehicleNickname := "My Box Truck"
	weightTicketSetDocument := models.WeightTicketSetDocument{
		MoveDocumentID:           moveDocument.ID,
		MoveDocument:             moveDocument,
		EmptyWeight:              &emptyWeight1,
		EmptyWeightTicketMissing: false,
		FullWeight:               &fullWeight1,
		FullWeightTicketMissing:  false,
		VehicleNickname:          &vehicleNickname,
		WeightTicketSetType:      models.WeightTicketSetTypeBOXTRUCK,
		WeightTicketDate:         &testdatagen.NextValidMoveDate,
		TrailerOwnershipMissing:  false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument)
	suite.NoVerrs(verrs)
	suite.NoError(err)

	suite.Run("weight ticket set has model but not make fails", func() {
		vehicleModel := "Wagon"
		emptyWeight := (int64)(1000)
		fullWeight := (int64)(2500)
		weightTicketSetType := internalmessages.WeightTicketSetTypeCAR

		updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
			ID:                  handlers.FmtUUID(moveDocument.ID),
			MoveID:              handlers.FmtUUID(move.ID),
			Title:               handlers.FmtString("super_awesome.pdf"),
			Notes:               handlers.FmtString("This document is super awesome."),
			Status:              internalmessages.MoveDocumentStatusOK,
			VehicleNickname:     &vehicleNickname,
			VehicleMake:         nil,
			VehicleModel:        &vehicleModel,
			EmptyWeight:         &emptyWeight,
			FullWeight:          &fullWeight,
			WeightTicketSetType: &weightTicketSetType,
			MoveDocumentType:    internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
			WeightTicketDate:    handlers.FmtDate(wtDate),
		}

		originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
		suite.Nil(err)
		_, verrs, err = wtu.Update(updateMoveDocPayload, originalMoveDocument, session)
		suite.NoError(err)
		suite.NotEmpty(verrs)
	})

	suite.Run("weight ticket set has make but not model fails", func() {
		vehicleMake := "Radio Flyer"
		emptyWeight := (int64)(1000)
		fullWeight := (int64)(2500)
		weightTicketSetType := internalmessages.WeightTicketSetTypeCAR

		updateMoveDocPayload := &internalmessages.MoveDocumentPayload{
			ID:                  handlers.FmtUUID(moveDocument.ID),
			MoveID:              handlers.FmtUUID(move.ID),
			Title:               handlers.FmtString("super_awesome.pdf"),
			Notes:               handlers.FmtString("This document is super awesome."),
			Status:              internalmessages.MoveDocumentStatusOK,
			VehicleNickname:     &vehicleNickname,
			VehicleMake:         &vehicleMake,
			VehicleModel:        nil,
			EmptyWeight:         &emptyWeight,
			FullWeight:          &fullWeight,
			WeightTicketSetType: &weightTicketSetType,
			MoveDocumentType:    internalmessages.MoveDocumentTypeWEIGHTTICKETSET,
			WeightTicketDate:    handlers.FmtDate(wtDate),
		}

		originalMoveDocument, err := models.FetchMoveDocument(suite.DB(), session, moveDocument.ID, false)
		suite.Nil(err)
		_, verrs, err = wtu.Update(updateMoveDocPayload, originalMoveDocument, session)
		suite.NoError(err)
		suite.NotEmpty(verrs)
	})
} */
