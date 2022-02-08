package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicWeightTicketSetDocumentInstantiation() {
	expenseDoc := &models.WeightTicketSetDocument{}

	expErrors := map[string][]string{
		"move_document_id":       {"MoveDocumentID can not be blank."},
		"weight_ticket_set_type": {"WeightTicketSetType can not be blank."},
	}

	suite.verifyValidationErrors(expenseDoc, expErrors)
}

func (suite *ModelSuite) TestBasicWeightTicketSetDocumentMustHaveMakeAndModel() {
	expenseDoc := &models.WeightTicketSetDocument{
		MoveDocumentID:      uuid.FromStringOrNil("d5602d02-45f5-412b-893b-aacf130051c9"),
		WeightTicketSetType: models.WeightTicketSetTypeCAR,
		VehicleMake:         models.StringPointer("Honda")}

	expErrors := map[string][]string{
		"vehicle_make": {"VehicleMake can not be nil if VehicleModel has a value and vice versa"},
	}

	suite.verifyValidationErrors(expenseDoc, expErrors)
}

//func (suite *ModelSuite) TestCalculateNetWeightWeightTicketAwaitingReview() {
//	// When: there is a move and move document
//	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
//		PersonallyProcuredMove: models.PersonallyProcuredMove{
//			Status: models.PPMStatusPAYMENTREQUESTED,
//		},
//	})
//	move := ppm.Move
//	sm := move.Orders.ServiceMember
//	session := &auth.Session{
//		ApplicationName: auth.MilApp,
//		UserID:          sm.UserID,
//		ServiceMemberID: sm.ID,
//	}
//	moveDoc1 := testdatagen.MakeMoveDocument(suite.DB(),
//		testdatagen.Assertions{
//			MoveDocument: models.MoveDocument{
//				MoveID:                   move.ID,
//				Move:                     move,
//				PersonallyProcuredMoveID: &ppm.ID,
//				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
//				Status:                   models.MoveDocumentStatusOK,
//			},
//		})
//
//	vehicleNickname := "My Car"
//	emptyWeight1 := unit.Pound(1000)
//	fullWeight1 := unit.Pound(2500)
//	weightTicketSetDocument1 := models.WeightTicketSetDocument{
//		MoveDocumentID:           moveDoc1.ID,
//		MoveDocument:             moveDoc1,
//		EmptyWeight:              &emptyWeight1,
//		EmptyWeightTicketMissing: false,
//		FullWeight:               &fullWeight1,
//		FullWeightTicketMissing:  false,
//		VehicleNickname:          &vehicleNickname,
//		WeightTicketSetType:      "CAR",
//		WeightTicketDate:         &testdatagen.NextValidMoveDate,
//		TrailerOwnershipMissing:  false,
//	}
//	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument1)
//	suite.NoVerrs(verrs)
//	suite.NoError(err)
//	moveDoc2 := testdatagen.MakeMoveDocument(suite.DB(),
//		testdatagen.Assertions{
//			MoveDocument: models.MoveDocument{
//				MoveID:                   move.ID,
//				Move:                     move,
//				PersonallyProcuredMoveID: &ppm.ID,
//				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
//				Status:                   models.MoveDocumentStatusOK,
//			},
//		})
//	emptyWeight2 := unit.Pound(1000)
//	fullWeight2 := unit.Pound(2500)
//	weightTicketSetDocument2 := models.WeightTicketSetDocument{
//		MoveDocumentID:           moveDoc2.ID,
//		MoveDocument:             moveDoc2,
//		EmptyWeight:              &emptyWeight2,
//		EmptyWeightTicketMissing: false,
//		FullWeight:               &fullWeight2,
//		FullWeightTicketMissing:  false,
//		VehicleNickname:          &vehicleNickname,
//		WeightTicketSetType:      "CAR",
//		WeightTicketDate:         &testdatagen.NextValidMoveDate,
//		TrailerOwnershipMissing:  false,
//	}
//	verrs, err = suite.DB().ValidateAndCreate(&weightTicketSetDocument2)
//	suite.NoVerrs(verrs)
//	suite.NoError(err)
//
//	status := models.MoveDocumentStatusOK
//	wts, err := models.FetchMoveDocuments(suite.DB(), session, ppm.ID, &status, models.MoveDocumentTypeWEIGHTTICKETSET, false)
//	suite.NoError(err)
//	suite.Len(wts, 2)
//
//	total, err := models.SumWeightTicketSetsForPPM(suite.DB(), session, ppm.ID)
//	suite.NoError(err)
//	expectedTotal := (fullWeight1 + fullWeight2) - (emptyWeight1 + emptyWeight2)
//	suite.Equal(&expectedTotal, total)
//
//}
//
//func (suite *ModelSuite) TestCalculateNetWeightNoWeightTicket() {
//	// When: there is a move and move document
//	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
//		PersonallyProcuredMove: models.PersonallyProcuredMove{
//			Status: models.PPMStatusPAYMENTREQUESTED,
//		},
//	})
//	move := ppm.Move
//	sm := move.Orders.ServiceMember
//	session := &auth.Session{
//		ApplicationName: auth.MilApp,
//		UserID:          sm.UserID,
//		ServiceMemberID: sm.ID,
//	}
//	total, err := models.SumWeightTicketSetsForPPM(suite.DB(), session, ppm.ID)
//	expectedTotal := unit.Pound(0)
//	suite.NoError(err)
//	suite.Equal(&expectedTotal, total)
//
//}
//
//func (suite *ModelSuite) TestCalculateNetWeight() {
//	// When: there is a move and move document
//	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
//		PersonallyProcuredMove: models.PersonallyProcuredMove{
//			Status: models.PPMStatusPAYMENTREQUESTED,
//		},
//	})
//	move := ppm.Move
//	sm := move.Orders.ServiceMember
//	session := &auth.Session{
//		ApplicationName: auth.MilApp,
//		UserID:          sm.UserID,
//		ServiceMemberID: sm.ID,
//	}
//	moveDoc1 := testdatagen.MakeMoveDocument(suite.DB(),
//		testdatagen.Assertions{
//			MoveDocument: models.MoveDocument{
//				MoveID:                   move.ID,
//				Move:                     move,
//				PersonallyProcuredMoveID: &ppm.ID,
//				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
//				Status:                   models.MoveDocumentStatusOK,
//			},
//		})
//
//	vehicleNickname := "My Car"
//	emptyWeight1 := unit.Pound(1000)
//	fullWeight1 := unit.Pound(2500)
//	weightTicketSetDocument1 := models.WeightTicketSetDocument{
//		MoveDocumentID:           moveDoc1.ID,
//		MoveDocument:             moveDoc1,
//		EmptyWeight:              &emptyWeight1,
//		EmptyWeightTicketMissing: false,
//		FullWeight:               &fullWeight1,
//		FullWeightTicketMissing:  false,
//		VehicleNickname:          &vehicleNickname,
//		WeightTicketSetType:      "CAR",
//		WeightTicketDate:         &testdatagen.NextValidMoveDate,
//		TrailerOwnershipMissing:  false,
//	}
//	verrs, err := suite.DB().ValidateAndCreate(&weightTicketSetDocument1)
//	suite.NoVerrs(verrs)
//	suite.NoError(err)
//	moveDoc2 := testdatagen.MakeMoveDocument(suite.DB(),
//		testdatagen.Assertions{
//			MoveDocument: models.MoveDocument{
//				MoveID:                   move.ID,
//				Move:                     move,
//				PersonallyProcuredMoveID: &ppm.ID,
//				MoveDocumentType:         models.MoveDocumentTypeWEIGHTTICKETSET,
//				Status:                   models.MoveDocumentStatusAWAITINGREVIEW,
//			},
//		})
//
//	emptyWeight2 := unit.Pound(1000)
//	fullWeight2 := unit.Pound(2500)
//	weightTicketSetDocument2 := models.WeightTicketSetDocument{
//		MoveDocumentID:           moveDoc2.ID,
//		MoveDocument:             moveDoc2,
//		EmptyWeight:              &emptyWeight2,
//		EmptyWeightTicketMissing: false,
//		FullWeight:               &fullWeight2,
//		FullWeightTicketMissing:  false,
//		VehicleNickname:          &vehicleNickname,
//		WeightTicketSetType:      "CAR",
//		WeightTicketDate:         &testdatagen.NextValidMoveDate,
//		TrailerOwnershipMissing:  false,
//	}
//	verrs, err = suite.DB().ValidateAndCreate(&weightTicketSetDocument2)
//	suite.NoVerrs(verrs)
//	suite.NoError(err)
//
//	total, err := models.SumWeightTicketSetsForPPM(suite.DB(), session, ppm.ID)
//	suite.NoError(err)
//	expectedTotal := fullWeight1 - emptyWeight1
//	suite.Equal(&expectedTotal, total)
//
//}
