package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchAllMoveDocumentsForMove() {
	// When: there is a move and move document
	move := testdatagen.MakeDefaultMove(suite.db)
	sm := move.Orders.ServiceMember

	assertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID: move.ID,
			Move:   move,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
		},
	}

	testdatagen.MakeMoveDocument(suite.db, assertions)
	testdatagen.MakeMovingExpenseDocument(suite.db, assertions)

	docs, err := move.FetchAllMoveDocumentsForMove(suite.db)
	if suite.NoError(err) {
		suite.Len(docs, 2)
	}
}
