package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildMovingExpense() {
	const defaultDescription = "Packing Peanuts"
	const defaultAmount = unit.Cents(2345)
	const defaultMovingExpenseType = models.MovingExpenseReceiptTypePackingMaterials
	suite.Run("Successful creation of weight ticket ", func() {
		// Under test:      BuildMovingExpense
		// Mocked:          None
		// Set up:          Create a weight ticket with no customizations or traits
		// Expected outcome:movingExpense should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		movingExpense := BuildMovingExpense(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(movingExpense.Description)
		suite.Equal(defaultDescription, *movingExpense.Description)

		suite.NotNil(movingExpense.MovingExpenseType)
		suite.Equal(defaultMovingExpenseType, *movingExpense.MovingExpenseType)

		suite.NotNil(movingExpense.Amount)
		suite.Equal(defaultAmount, *movingExpense.Amount)

		suite.NotNil(movingExpense.PaidWithGTCC)
		suite.True(*movingExpense.PaidWithGTCC)

		suite.False(movingExpense.PPMShipmentID.IsNil())
		suite.False(movingExpense.PPMShipment.ID.IsNil())

		suite.False(movingExpense.DocumentID.IsNil())
		suite.False(movingExpense.Document.ID.IsNil())
		suite.NotEmpty(movingExpense.Document.UserUploads)

		serviceMemberID := movingExpense.PPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID
		suite.False(serviceMemberID.IsNil())
		suite.Equal(serviceMemberID, movingExpense.Document.ServiceMemberID)
	})

	suite.Run("Successful creation of customized MovingExpense", func() {
		// Under test:      BuildMovingExpense
		// Mocked:          None
		// Set up:          Create a weight ticket with and pass custom fields
		// Expected outcome:movingExpense should be created with custom values

		// SETUP
		customPPMShipment := models.PPMShipment{
			ExpectedDepartureDate: time.Now(),
			ActualMoveDate:        models.TimePointer(time.Now().Add(time.Duration(24 * time.Hour))),
		}
		customMovingExpense := models.MovingExpense{
			Description: models.StringPointer("VIP expense"),
			Reason:      models.StringPointer("VIP reason"),
			Amount:      models.CentPointer(9999),
		}
		customs := []Customization{
			{
				Model: customPPMShipment,
			},
			{
				Model: customMovingExpense,
			},
		}
		// CALL FUNCTION UNDER TEST
		movingExpense := BuildMovingExpense(suite.DB(), customs, nil)

		// VALIDATE RESULTS
		suite.NotNil(movingExpense.Description)
		suite.Equal(*customMovingExpense.Description,
			*movingExpense.Description)

		suite.NotNil(movingExpense.Amount)
		suite.Equal(*customMovingExpense.Amount, *movingExpense.Amount)

		suite.NotNil(movingExpense.Reason)
		suite.Equal(*customMovingExpense.Reason,
			*movingExpense.Reason)

		suite.False(movingExpense.PPMShipmentID.IsNil())
		suite.False(movingExpense.PPMShipment.ID.IsNil())
		suite.Equal(customPPMShipment.ExpectedDepartureDate,
			movingExpense.PPMShipment.ExpectedDepartureDate)
		suite.NotNil(movingExpense.PPMShipment.ActualMoveDate)
		suite.Equal(*customPPMShipment.ActualMoveDate,
			*movingExpense.PPMShipment.ActualMoveDate)
	})

	suite.Run("Successful return of linkOnly MovingExpense", func() {
		// Under test:       BuildMovingExpense
		// Set up:           Pass in a linkOnly movingExpense
		// Expected outcome: No new MovingExpense should be created.

		// Check num MovingExpense records
		precount, err := suite.DB().Count(&models.MovingExpense{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		movingExpense := BuildMovingExpense(suite.DB(), []Customization{
			{
				Model: models.MovingExpense{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.MovingExpense{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, movingExpense.ID)
	})

	suite.Run("Successful return of stubbed MovingExpense", func() {
		// Under test:       BuildMovingExpense
		// Set up:           Pass in nil db
		// Expected outcome: No new MovingExpense should be created.

		// Check num MovingExpense records
		precount, err := suite.DB().Count(&models.MovingExpense{})
		suite.NoError(err)

		customMovingExpense := models.MovingExpense{
			Amount: models.CentPointer(888),
		}
		// Nil passed in as db
		movingExpense := BuildMovingExpense(nil, []Customization{
			{
				Model: customMovingExpense,
			},
		}, nil)

		count, err := suite.DB().Count(&models.MovingExpense{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(*customMovingExpense.Amount, *movingExpense.Amount)
	})
}
