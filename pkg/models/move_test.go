package models_test

import (
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicMoveInstantiation() {
	move := &Move{}

	expErrors := map[string][]string{
		"orders_id": {"OrdersID can not be blank."},
	}

	suite.verifyValidationErrors(move, expErrors)
}

func (suite *ModelSuite) TestFetchMove() {
	order1, _ := testdatagen.MakeOrder(suite.db)
	order2, _ := testdatagen.MakeOrder(suite.db)

	var selectedType = internalmessages.SelectedMoveTypeCOMBO
	move, err := order1.CreateNewMove(suite.db, &selectedType)
	suite.Nil(err)
	suite.Equal(6, len(move.Locator))

	// All correct
	fetchedMove, err := FetchMove(suite.db, order1.ServiceMember.User, move.ID)
	suite.Nil(err, "Expected to get moveResult back.")
	suite.Equal(fetchedMove.ID, move.ID, "Expected new move to match move.")

	// Bad Move
	fetchedMove, err = FetchMove(suite.db, order1.ServiceMember.User, uuid.Must(uuid.NewV4()))
	suite.Equal(ErrFetchNotFound, err, "Expected to get FetchNotFound.")

	// Bad User
	fetchedMove, err = FetchMove(suite.db, order2.ServiceMember.User, move.ID)
	suite.Equal(ErrFetchForbidden, err, "Expected to get a Forbidden back.")
}
