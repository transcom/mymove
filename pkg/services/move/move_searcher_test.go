package move

import (
	"fmt"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *MoveServiceSuite) TestMoveSearch() {
	searcher := NewMoveSearcher()

	suite.Run("search with no filters should fail", func() {
		factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "AAAAAA",
				},
			},
		}, nil)

		factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "BBBBBB",
				},
			},
		}, nil)

		_, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{})
		suite.Error(err)
	})
	suite.Run("search with valid locator", func() {
		firstMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "AAAAAA",
				},
			},
		}, nil)

		factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "BBBBBB",
				},
			},
		}, nil)

		moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{Locator: &firstMove.Locator})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(firstMove.Locator, moves[0].Locator)
	})
	suite.Run("search with valid DOD ID", func() {
		factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "AAAAAA",
				},
			},
		}, nil)

		secondMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "BBBBBB",
				},
			},
		}, nil)

		moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{DodID: secondMove.Orders.ServiceMember.Edipi})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(secondMove.Locator, moves[0].Locator)
	})
	suite.Run("search with customer name", func() {
		firstMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "AAAAAA",
				},
			},
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Grace"),
					LastName:  models.StringPointer("Griffin"),
				},
			},
		}, nil)

		_ = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "BBBBBB",
				},
			},
		}, nil)

		moves, _, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{CustomerName: models.StringPointer("Grace Griffin")})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(firstMove.Locator, moves[0].Locator)
	})
	suite.Run("search with both DOD ID and locator filters should fail", func() {

		firstMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "AAAAAA",
				},
			},
		}, nil)

		secondMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "BBBBBB",
				},
			},
		}, nil)

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
		firstMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "AAAAAA",
				},
			},
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Grace"),
					LastName:  models.StringPointer("Griffin"),
				},
			},
		}, nil)

		secondMove := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "BBBBBB",
				},
			},
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Grace"),
					LastName:  models.StringPointer("Groffin"),
				},
			},
		}, nil)
		// get first page
		moves, totalCount, err := searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{
			CustomerName: models.StringPointer("grace griffin"),
			PerPage:      1,
			Page:         1,
		})
		suite.NoError(err)
		suite.Len(moves, 1)
		suite.Equal(firstMove.Locator, moves[0].Locator)
		suite.Equal(2, totalCount)

		// get second page
		moves, totalCount, err = searcher.SearchMoves(suite.AppContextForTest(), &services.SearchMovesParams{
			CustomerName: models.StringPointer("grace griffin"),
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
	firstMoveOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{PostalCode: "89523"},
		},
	}, nil)
	firstMoveNewDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{PostalCode: "11111"},
		},
	}, nil)

	firstMove := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:   models.StringPointer("María"),
				LastName:    models.StringPointer("Johnson"),
				Affiliation: &armyAffiliation,
			},
		},
		{
			Model: models.Move{
				Locator: "MOVE01",
				Status:  models.MoveStatusDRAFT,
			},
		},
		{
			Model:    firstMoveOriginDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    firstMoveNewDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    firstMove,
			LinkOnly: true,
		},
	}, nil)
	secondMoveOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{PostalCode: "90211"},
		},
	}, nil)
	secondMoveNewDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{PostalCode: "22222"},
		},
	}, nil)

	secondMove := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName:   models.StringPointer("Mariah"),
				LastName:    models.StringPointer("Johnson"),
				Affiliation: &navyAffiliation,
			},
		},
		{
			Model: models.Move{
				Locator: "MOVE02",
				Status:  models.MoveStatusNeedsServiceCounseling,
			},
		},
		{
			Model:    secondMoveOriginDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    secondMoveNewDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
	}, nil)
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    secondMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    secondMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	return firstMove, secondMove
}
func (suite *MoveServiceSuite) TestMoveSearchOrdering() {
	suite.Run("search results ordering", func() {
		firstMove, secondMove := setupTestData(suite)
		testMoves := models.Moves{}
		suite.NoError(suite.DB().EagerPreload("Orders", "Orders.NewDutyLocation").All(&testMoves))

		searcher := NewMoveSearcher()
		columns := []string{"status", "originPostalCode", "destinationPostalCode", "branch", "shipmentsCount"}
		for _, order := range []string{"asc", "desc"} {
			order := order
			for ci, col := range columns {
				params := services.SearchMovesParams{
					CustomerName: models.StringPointer("maria johnson"),
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
			{column: "OriginPostalCode", value: secondMove.Orders.OriginDutyLocation.PostalCode, SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, OriginPostalCode: &secondMove.Orders.OriginDutyLocation.PostalCode}},
			{column: "Branch", value: string(*secondMove.Orders.ServiceMember.Affiliation), SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, Branch: models.StringPointer(secondMove.Orders.ServiceMember.Affiliation.String())}},
			{column: "ShipmentsCount", value: "2", SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, ShipmentsCount: models.Int64Pointer(2)}},
			{column: "DestinationPostalCode", value: secondMove.Orders.NewDutyLocation.PostalCode, SearchMovesParams: services.SearchMovesParams{CustomerName: &nameToSearch, DestinationPostalCode: &secondMove.Orders.NewDutyLocation.PostalCode}},
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
