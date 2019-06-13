package models_test

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestCreateNewMoveShow() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())

	selectedMoveType := SelectedMoveTypeHHG

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(true),
	}
	_, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Len(moves, 1)
}

func (suite *ModelSuite) TestCreateNewMoveShowFalse() {
	orders := testdatagen.MakeDefaultOrder(suite.DB())

	selectedMoveType := SelectedMoveTypeHHG

	moveOptions := MoveOptions{
		SelectedType: &selectedMoveType,
		Show:         swag.Bool(false),
	}
	_, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.Nil(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Empty(moves)
}

func (suite *ModelSuite) TestShowMovesDraftSubmittedApprovedPPMQueue() {
	// Make three PPMs with different move statuses
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusDRAFT,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
	})
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
	})

	// Move canceled should not return, for testing
	testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status: models.MoveStatusCANCELED,
		},
	})

	// Expected 3 moves for PPM queue returned
	moves, err := GetMoveQueueItems(suite.DB(), "ppm")
	suite.Nil(err)
	suite.Len(moves, 3)
}
