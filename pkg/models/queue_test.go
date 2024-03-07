package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCreateMoveWithPPMShow() {

	factory.BuildMoveWithPPMShipment(suite.DB(), nil, nil)

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Len(moves, 1)
}

func (suite *ModelSuite) TestCreateMoveWithPPMNoShow() {
	moveTemplate := Move{
		Show: BoolPointer(false),
	}
	factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: moveTemplate,
		},
	}, nil)

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Empty(moves)

}

func (suite *ModelSuite) TestCreateNewMoveWithNoPPMShow() {
	orders := factory.BuildOrder(suite.DB(), nil, nil)
	factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	moveOptions := MoveOptions{
		Show: BoolPointer(true),
	}
	_, verrs, err := orders.CreateNewMove(suite.DB(), moveOptions)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "failed to validate move")

	moves, moveErrs := GetMoveQueueItems(suite.DB(), "all")
	suite.Nil(moveErrs)
	suite.Empty(moves)
}

func (suite *ModelSuite) TestQueueNotFound() {
	moves, moveErrs := GetMoveQueueItems(suite.DB(), "queue_not_found")
	suite.Equal(ErrFetchNotFound, moveErrs, "Expected not to find move queue items")
	suite.Empty(moves)
}
