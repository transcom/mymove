package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchMoveTaskOrder() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	moveTaskOrder, err := FetchMoveTaskOrder(suite.DB(), mto.ID)
	if suite.NoError(err) {
		suite.Equal(moveTaskOrder.ID, mto.ID)
		suite.Equal(moveTaskOrder.MoveID, mto.MoveID)
	}

}
