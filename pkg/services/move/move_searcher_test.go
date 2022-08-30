package move

import (
	"fmt"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveServiceSuite) TestMoveSearch() {
	searcher := NewMoveSearcher()

	suite.Run("search with no filters should fail", func() {
		testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "AAAAAA",
		}})
		testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "BBBBBB",
		}})

		_, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{})
		suite.Error(err)
	})
	suite.Run("search with valid locator", func() {
		firstMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "AAAAAA",
		}})
		testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "BBBBBB",
		}})

		moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{Locator: &firstMove.Locator})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(firstMove.Locator, moves[0].Locator)
	})
	suite.Run("search with valid DOD ID", func() {
		testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "AAAAAA",
		}})
		secondMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "BBBBBB",
		}})

		moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{DodID: secondMove.Orders.ServiceMember.Edipi})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(secondMove.Locator, moves[0].Locator)
	})
	suite.Run("search with customer name", func() {
		firstMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Locator: "AAAAAA",
			},
			ServiceMember: models.ServiceMember{
				FirstName: swag.String("Grace"),
				LastName:  swag.String("Griffin"),
			},
		})
		_ = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "BBBBBB",
		}})

		moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{CustomerName: swag.String("Grace Griffin")})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(firstMove.Locator, moves[0].Locator)
	})
	suite.Run("search with both DOD ID and locator filters should fail", func() {
		firstMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "AAAAAA",
		}})
		secondMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{
			Locator: "BBBBBB",
		}})

		// Search for Locator of one move and DOD ID of another move
		_, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{
			Locator: &firstMove.Locator,
			DodID:   secondMove.Orders.ServiceMember.Edipi,
		})
		suite.Error(err)
	})
	suite.Run("search with no results", func() {
		nonexistantLocator := "CCCCCC"
		moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{Locator: &nonexistantLocator})
		suite.NoError(err)
		suite.Len(moves, 0)
	})

	suite.Run("test pagination", func() {
		firstMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Locator: "AAAAAA",
			},
			ServiceMember: models.ServiceMember{
				FirstName: swag.String("Grace"),
				LastName:  swag.String("Griffin"),
			},
		})
		secondMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Locator: "BBBBBB",
			},
			ServiceMember: models.ServiceMember{
				FirstName: swag.String("Grace"),
				LastName:  swag.String("Groffin"),
			},
		})
		// get first page
		moves, totalCount, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{
			CustomerName: swag.String("grace griffin"),
			PerPage:      1,
			Page:         1,
		})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(firstMove.Locator, moves[0].Locator)
		suite.Equal(2, totalCount)

		// get second page
		moves, totalCount, err = searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{
			CustomerName: swag.String("grace griffin"),
			PerPage:      1,
			Page:         2,
		})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(secondMove.Locator, moves[0].Locator)
		suite.Equal(2, totalCount)
	})
}
func setupTestData(suite *MoveServiceSuite) (models.Move, models.Move) {
	armyAffiliation := models.AffiliationARMY
	navyAffiliation := models.AffiliationNAVY
	firstMoveOriginDutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		Address: models.Address{PostalCode: "89523"},
	})
	firstMoveNewDutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		// For some reason, I need to include both of these. If either one is omitted,
		// the postal code in the Go model returned by MakeDutyLocation won't match what gets
		// saved in the database.
		Address:      models.Address{PostalCode: "11111"},
		DutyLocation: models.DutyLocation{Address: models.Address{PostalCode: "11111"}},
	})
	firstMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:   swag.String("María"),
			LastName:    swag.String("Johnson"),
			Affiliation: &armyAffiliation,
		},
		Move: models.Move{
			Locator: "MOVE01",
			Status:  models.MoveStatusDRAFT,
		},
		OriginDutyLocation: firstMoveOriginDutyLocation,
		Order: models.Order{
			NewDutyLocationID: firstMoveNewDutyLocation.ID,
			NewDutyLocation:   firstMoveNewDutyLocation,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: firstMove,
	})
	secondMoveOriginDutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		Address: models.Address{PostalCode: "90211"},
	})
	secondMoveNewDutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		// For some reason, I need to include both of these. If I omit either one,
		// the postal code in the Go model returned by MakeDutyLocation won't match what gets
		// saved in the database.
		Address:      models.Address{PostalCode: "22222"},
		DutyLocation: models.DutyLocation{Address: models.Address{PostalCode: "22222"}},
	})
	secondMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:   swag.String("Mariah"),
			LastName:    swag.String("Johnson"),
			Affiliation: &navyAffiliation,
		},
		Move: models.Move{
			Locator: "MOVE02",
			Status:  models.MoveStatusNeedsServiceCounseling,
		},
		OriginDutyLocation: secondMoveOriginDutyLocation,
		Order:              models.Order{NewDutyLocationID: secondMoveNewDutyLocation.ID, NewDutyLocation: secondMoveNewDutyLocation},
	})
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        secondMove,
		MTOShipment: models.MTOShipment{Status: models.MTOShipmentStatusSubmitted},
	})
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        secondMove,
		MTOShipment: models.MTOShipment{Status: models.MTOShipmentStatusApproved},
	})

	return firstMove, secondMove
}
func (suite *MoveServiceSuite) TestMoveSearchOrdering() {
	suite.Run("search results ordering", func() {
		firstMove, secondMove := setupTestData(suite)
		testMoves := models.Moves{}
		suite.NoError(suite.DB().EagerPreload("Orders", "Orders.NewDutyLocation", "Orders.NewDutyLocation.Address").All(&testMoves))

		searcher := NewMoveSearcher()
		columns := []string{"status", "originPostalCode", "destinationPostalCode", "branch", "shipmentsCount"}
		for _, order := range []string{"asc", "desc"} {
			order := order
			for ci, col := range columns {
				params := services.SearchMovesParams{
					CustomerName: swag.String("maria johnson"),
					Sort:         &columns[ci],
					Order:        &order,
				}
				moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &params)
				suite.NoError(err)
				suite.Len(moves, 2)
				message := fmt.Sprintf("Sort by %s, %s failed", col, order)
				if order == "asc" {
					suite.Equal(firstMove.Locator, moves[0].Locator, message)
					suite.Equal(secondMove.Locator, moves[1].Locator, message)
				} else {
					suite.Equal(firstMove.Locator, moves[1].Locator, message)
					suite.Equal(secondMove.Locator, moves[0].Locator, message)
				}
			}
		}
	})
	suite.Run("search results filtering", func() {
		_, secondMove := setupTestData(suite)
		nameToSearch := "maria johnson"
		searcher := NewMoveSearcher()

		cases := []struct {
			column string
			value  string
			services.SearchMovesParams
		}{
			{column: "Status", value: fmt.Sprintf("[%s]", string(secondMove.Status)), SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, Status: []string{string(secondMove.Status)}}},
			{column: "OriginPostalCode", value: secondMove.Orders.OriginDutyLocation.Address.PostalCode, SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, OriginPostalCode: &secondMove.Orders.OriginDutyLocation.Address.PostalCode}},
			{column: "Branch", value: string(*secondMove.Orders.ServiceMember.Affiliation), SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, Branch: swag.String(secondMove.Orders.ServiceMember.Affiliation.String())}},
			{column: "ShipmentsCount", value: "2", SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, ShipmentsCount: swag.Int64(2)}},
			{column: "DestinationPostalCode", value: secondMove.Orders.NewDutyLocation.Address.PostalCode, SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, DestinationPostalCode: &secondMove.Orders.NewDutyLocation.Address.PostalCode}},
		}
		for _, testCase := range cases {
			message := fmt.Sprintf("Filtering results of search by column %s = %s has failed", testCase.column, testCase.value)
			moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &testCase.SearchMovesParams)
			suite.NoError(err)
			suite.Len(moves, 1, message)
			suite.Equal(secondMove.Locator, moves[0].Locator, message)
		}
	})
}
