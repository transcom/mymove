package movingexpense

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *MovingExpenseSuite) TestMovingExpenseCreator() {
	suite.Run("Successfully creates a MovingExpense - Customer", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID

		session := &auth.Session{
			ApplicationName: auth.MilApp,
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

	suite.Run("Fails when an invalid ppmShipmentID is used - Customer", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		}

		movingExpenseCreator := NewMovingExpenseCreator()
		movingExpense, err := movingExpenseCreator.CreateMovingExpense(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(movingExpense)

		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})

	suite.Run("Fails when session has invalid serviceMemberID", func() {
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		}
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		movingExpenseCreator := NewMovingExpenseCreator()
		movingExpense, err := movingExpenseCreator.CreateMovingExpense(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(movingExpense)

		expectedErr := apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})

	suite.Run("Successfully creates a MovingExpense - Office", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID

		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		session := &auth.Session{
			OfficeUserID:    officeUser.ID,
			ApplicationName: auth.OfficeApp,
		}

		movingExpenseCreator := NewMovingExpenseCreator()
		movingExpense, err := movingExpenseCreator.CreateMovingExpense(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(movingExpense)
		suite.Equal(ppmShipment.ID, movingExpense.PPMShipmentID)
		suite.NotNil(movingExpense.DocumentID)
		suite.Equal(serviceMemberID, movingExpense.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used - Office", func() {
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		session := &auth.Session{
			OfficeUserID:    officeUser.ID,
			ApplicationName: auth.OfficeApp,
		}

		movingExpenseCreator := NewMovingExpenseCreator()
		movingExpense, err := movingExpenseCreator.CreateMovingExpense(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(movingExpense)

		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})
}
