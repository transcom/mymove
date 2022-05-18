package move

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestMoveSearch() {
	firstMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
		Locator: "AAAAAA",
	}})
	secondMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
		Locator: "BBBBBB",
	}})
	searcher := NewMoveSearcher()

	suite.Run("search with no filters should fail", func() {
		_, err := searcher.SearchMoves(suite.AppContextForTest(), nil, nil)
		suite.Error(err)
	})
	suite.Run("search with valid locator", func() {
		moves, err := searcher.SearchMoves(suite.AppContextForTest(), &firstMove.Locator, nil)
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(firstMove.Locator, moves[0].Locator)
	})
	suite.Run("search with valid DOD ID", func() {
		moves, err := searcher.SearchMoves(suite.AppContextForTest(), nil, secondMove.Orders.ServiceMember.Edipi)
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(secondMove.Locator, moves[0].Locator)
	})
	suite.Run("search with both DOD ID and locator filters should fail", func() {
		// Search for Locator of one move and DOD ID of another move
		_, err := searcher.SearchMoves(suite.AppContextForTest(), &firstMove.Locator, secondMove.Orders.ServiceMember.Edipi)
		suite.Error(err)
	})
	suite.Run("search with no results", func() {
		nonexistantLocator := "CCCCCC"
		moves, err := searcher.SearchMoves(suite.AppContextForTest(), &nonexistantLocator, nil)
		suite.NoError(err)
		suite.Len(moves, 0)
	})
}
