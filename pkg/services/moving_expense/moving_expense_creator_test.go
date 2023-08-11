package movingexpense

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *MovingExpenseSuite) TestMovingExpenseCreator() {
	suite.Run("Successfully creates a MovingExpense", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID

		session := &auth.Session{
			ServiceMemberID: serviceMemberID,
		}

		movingExpenseCreator := NewMovingExpenseCreator()
		movingExpense, err := movingExpenseCreator.CreateMovingExpense(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(movingExpense)
		suite.Equal(ppmShipment.ID, movingExpense.PPMShipmentID)
		suite.NotNil(movingExpense.DocumentID)
		suite.Equal(serviceMemberID, movingExpense.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
		}

		movingExpenseCreator := NewMovingExpenseCreator()
		movingExpense, err := movingExpenseCreator.CreateMovingExpense(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(movingExpense)
		suite.ErrorContains(err, fmt.Sprintf("Error fetching PPM with ID %s", uuid.Nil))
	})

	suite.Run("Fails when session has invalid serviceMemberID", func() {
		session := &auth.Session{
			ServiceMemberID: uuid.Nil,
		}
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		movingExpenseCreator := NewMovingExpenseCreator()
		movingExpense, err := movingExpenseCreator.CreateMovingExpense(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(movingExpense)
		suite.NotNil(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "No such shipment found for this service member")
	})

}
