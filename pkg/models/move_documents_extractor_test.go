package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

//func (suite *ModelSuite) TestFetchAllMoveDocumentsForMove() {
//	// When: there is a move and move document
//	move := testdatagen.MakeDefaultMove(suite.DB())
//	sm := move.Orders.ServiceMember
//
//	assertions := testdatagen.Assertions{
//		MoveDocument: models.MoveDocument{
//			MoveID: move.ID,
//			Move:   move,
//		},
//		Document: models.Document{
//			ServiceMemberID: sm.ID,
//			ServiceMember:   sm,
//		},
//	}
//
//	testdatagen.MakeMoveDocument(suite.DB(), assertions)
//	testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions)
//
//	moveDocument2 := testdatagen.MakeMoveDocument(suite.DB(), assertions)
//	weightTicketCarAssertions := testdatagen.Assertions{
//		WeightTicketSetDocument: models.WeightTicketSetDocument{
//			VehicleMake:              models.StringPointer("Honda"),
//			VehicleModel:             models.StringPointer("Civic"),
//			WeightTicketSetType:      models.WeightTicketSetTypeCAR,
//			MoveDocument: models.MoveDocument{
//				ID: moveDocument2.ID,
//				MoveID: move.ID,
//				Move: move,
//				MoveDocumentType: models.MoveDocumentTypeWEIGHTTICKET,
//				Document: models.Document{
//					ServiceMemberID: sm.ID,
//					ServiceMember:   sm,
//				},
//			},
//
//		},
//	}
//	carWeightTicketSetDocument := testdatagen.MakeWeightTicketSetDocument(suite.DB(), weightTicketCarAssertions)
//	moveDocument2.WeightTicketSetDocument = &carWeightTicketSetDocument
//
//	moveDocument3 := testdatagen.MakeMoveDocument(suite.DB(), assertions)
//	weightTicketTruckAssertions := testdatagen.Assertions{
//		WeightTicketSetDocument: models.WeightTicketSetDocument{
//			MoveDocument: models.MoveDocument{
//				ID:               moveDocument3.ID,
//				MoveID: move.ID,
//				Move: move,
//				MoveDocumentType: models.MoveDocumentTypeWEIGHTTICKET,
//				Document: models.Document{
//					ServiceMemberID: sm.ID,
//					ServiceMember:   sm,
//				},
//			},
//		},
//	}
//	truckWeightTicketSetDocument := testdatagen.MakeWeightTicketSetDocument(suite.DB(), weightTicketTruckAssertions)
//	moveDocument3.WeightTicketSetDocument = &truckWeightTicketSetDocument
//
//	deletedAt := time.Date(2019, 8, 7, 0, 0, 0, 0, time.UTC)
//	deleteAssertions := testdatagen.Assertions{
//		MoveDocument: models.MoveDocument{
//			MoveID:    move.ID,
//			Move:      move,
//			DeletedAt: &deletedAt,
//		},
//		Document: models.Document{
//			ServiceMemberID: sm.ID,
//			ServiceMember:   sm,
//			DeletedAt:       &deletedAt,
//		},
//	}
//	testdatagen.MakeMoveDocument(suite.DB(), deleteAssertions)
//
//	docs, err := move.FetchAllMoveDocumentsForMove(suite.DB(), false)
//	if suite.NoError(err) {
//		suite.Len(docs, 3)
//	}
//}

func (suite *ModelSuite) TestFetchAllMoveDocumentsForMove2() {
	// Create move an SM on which to attach move docs
	move := testdatagen.MakeDefaultMove(suite.DB())
	sm := move.Orders.ServiceMember

	moveDoc1 := suite.documentMakerHelper(move, sm, models.MoveDocumentTypeWEIGHTTICKET)
	carWeightTicketSetDocument := models.WeightTicketSetDocument{
		VehicleMake:         models.StringPointer("Honda"),
		VehicleModel:        models.StringPointer("Civic"),
		WeightTicketSetType: models.WeightTicketSetTypeBOXTRUCK,
		MoveDocumentID:      moveDoc1.ID,
		MoveDocument:        moveDoc1,
	}

	testdatagen.MakeWeightTicketSetDocument(suite.DB(), testdatagen.Assertions{
		WeightTicketSetDocument: carWeightTicketSetDocument,
	})

	moveDoc2 := suite.documentMakerHelper(move, sm, models.MoveDocumentTypeWEIGHTTICKET)
	truckWeightTicketSetDocument := models.WeightTicketSetDocument{
		VehicleNickname:     models.StringPointer("Hank the Tank"),
		WeightTicketSetType: models.WeightTicketSetTypeBOXTRUCK,
		MoveDocumentID:      moveDoc2.ID,
		MoveDocument:        moveDoc2,
	}

	testdatagen.MakeWeightTicketSetDocument(suite.DB(), testdatagen.Assertions{
		WeightTicketSetDocument: truckWeightTicketSetDocument,
	})

	docs, err := move.FetchAllMoveDocumentsForMove(suite.DB(), false)
	suite.NoError(err)

	// Check car weight ticket values
	carDoc := docs[0]
	suite.Equal(models.MoveDocumentTypeWEIGHTTICKET, carDoc.MoveDocumentType)
	suite.Equal(carWeightTicketSetDocument.VehicleMake, carDoc.VehicleMake)
	suite.Equal(carWeightTicketSetDocument.VehicleModel, carDoc.VehicleModel)

	truckDoc := docs[1]
	suite.Equal(models.MoveDocumentTypeWEIGHTTICKET, truckDoc.MoveDocumentType)
	suite.Equal(truckWeightTicketSetDocument.VehicleNickname, truckDoc.VehicleNickname)

	// Create move documents with expense docs

}

func (suite *ModelSuite) documentMakerHelper(move models.Move, sm models.ServiceMember, moveDocumentType models.MoveDocumentType) models.MoveDocument {
	// Create a Document and MoveDocument for a given move
	moveDocAssertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:           move.ID,
			Move:             move,
			MoveDocumentType: models.MoveDocumentTypeWEIGHTTICKET,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}
	// Create move document and weight ticket
	moveDoc := testdatagen.MakeMoveDocument(suite.DB(), moveDocAssertions)

	return moveDoc
}