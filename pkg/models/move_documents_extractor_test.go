package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchAllMoveDocumentsForMove() {
	// When: there is a move and move document
	move := testdatagen.MakeDefaultMove(suite.DB())
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

	testdatagen.MakeMoveDocument(suite.DB(), assertions)
	testdatagen.MakeMovingExpenseDocument(suite.DB(), assertions)
	testdatagen.MakeWeightTicketSetDocument(suite.DB(), assertions)

	deletedAt := time.Date(2019, 8, 7, 0, 0, 0, 0, time.UTC)
	deleteAssertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:    move.ID,
			Move:      move,
			DeletedAt: &deletedAt,
		},
		Document: models.Document{
			ServiceMemberID: sm.ID,
			ServiceMember:   sm,
			DeletedAt:       &deletedAt,
		},
	}
	testdatagen.MakeMoveDocument(suite.DB(), deleteAssertions)

	docs, err := move.FetchAllMoveDocumentsForMove(suite.DB())
	if suite.NoError(err) {
		suite.Len(docs, 3)
	}
}
