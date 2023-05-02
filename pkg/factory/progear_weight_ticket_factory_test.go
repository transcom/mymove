package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildProgearWeightTicket() {
	const defaultDescription = "professional equipment"
	const defaultWeight = unit.Pound(500)
	suite.Run("Successful creation of weight ticket ", func() {
		// Under test:      BuildProgearWeightTicket
		// Mocked:          None
		// Set up:          Create a weight ticket with no customizations or traits
		// Expected outcome:progearWeightTicket should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		progearWeightTicket := BuildProgearWeightTicket(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(progearWeightTicket.Description)
		suite.Equal(defaultDescription, *progearWeightTicket.Description)

		suite.NotNil(progearWeightTicket.Weight)
		suite.Equal(defaultWeight, *progearWeightTicket.Weight)

		suite.NotNil(progearWeightTicket.HasWeightTickets)
		suite.True(*progearWeightTicket.HasWeightTickets)

		suite.NotNil(progearWeightTicket.BelongsToSelf)
		suite.True(*progearWeightTicket.BelongsToSelf)

		suite.False(progearWeightTicket.PPMShipmentID.IsNil())
		suite.False(progearWeightTicket.PPMShipment.ID.IsNil())

		suite.False(progearWeightTicket.DocumentID.IsNil())
		suite.False(progearWeightTicket.Document.ID.IsNil())
		suite.NotEmpty(progearWeightTicket.Document.UserUploads)

		serviceMemberID := progearWeightTicket.PPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID
		suite.False(serviceMemberID.IsNil())
		suite.Equal(serviceMemberID, progearWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Successful creation of customized ProgearWeightTicket", func() {
		// Under test:      BuildProgearWeightTicket
		// Mocked:          None
		// Set up:          Create a weight ticket with and pass custom fields
		// Expected outcome:progearWeightTicket should be created with custom values

		// SETUP
		customPPMShipment := models.PPMShipment{
			ExpectedDepartureDate: time.Now(),
			ActualMoveDate:        models.TimePointer(time.Now().Add(time.Duration(24 * time.Hour))),
		}
		customProgearWeightTicket := models.ProgearWeightTicket{
			Description: models.StringPointer("VIP progear"),
			Weight:      models.PoundPointer(512),
			Reason:      models.StringPointer("VIP Reason"),
		}
		customs := []Customization{
			{
				Model: customPPMShipment,
			},
			{
				Model: customProgearWeightTicket,
			},
		}
		// CALL FUNCTION UNDER TEST
		progearWeightTicket := BuildProgearWeightTicket(suite.DB(), customs, nil)

		// VALIDATE RESULTS
		suite.NotNil(progearWeightTicket.Description)
		suite.Equal(*customProgearWeightTicket.Description,
			*progearWeightTicket.Description)

		suite.NotNil(progearWeightTicket.Weight)
		suite.Equal(*customProgearWeightTicket.Weight, *progearWeightTicket.Weight)

		suite.NotNil(progearWeightTicket.Reason)
		suite.Equal(*customProgearWeightTicket.Reason,
			*progearWeightTicket.Reason)

		suite.False(progearWeightTicket.PPMShipmentID.IsNil())
		suite.False(progearWeightTicket.PPMShipment.ID.IsNil())
		suite.Equal(customPPMShipment.ExpectedDepartureDate,
			progearWeightTicket.PPMShipment.ExpectedDepartureDate)
		suite.NotNil(progearWeightTicket.PPMShipment.ActualMoveDate)
		suite.Equal(*customPPMShipment.ActualMoveDate,
			*progearWeightTicket.PPMShipment.ActualMoveDate)
	})

	suite.Run("Successful return of linkOnly ProgearWeightTicket", func() {
		// Under test:       BuildProgearWeightTicket
		// Set up:           Pass in a linkOnly progearWeightTicket
		// Expected outcome: No new ProgearWeightTicket should be created.

		// Check num ProgearWeightTicket records
		precount, err := suite.DB().Count(&models.ProgearWeightTicket{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		progearWeightTicket := BuildProgearWeightTicket(suite.DB(), []Customization{
			{
				Model: models.ProgearWeightTicket{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.ProgearWeightTicket{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, progearWeightTicket.ID)
	})

	suite.Run("Successful return of stubbed ProgearWeightTicket", func() {
		// Under test:       BuildProgearWeightTicket
		// Set up:           Pass in nil db
		// Expected outcome: No new ProgearWeightTicket should be created.

		// Check num ProgearWeightTicket records
		precount, err := suite.DB().Count(&models.ProgearWeightTicket{})
		suite.NoError(err)

		customProgearWeightTicket := models.ProgearWeightTicket{
			Weight: models.PoundPointer(9999),
		}
		// Nil passed in as db
		progearWeightTicket := BuildProgearWeightTicket(nil, []Customization{
			{
				Model: customProgearWeightTicket,
			},
		}, nil)

		count, err := suite.DB().Count(&models.ProgearWeightTicket{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customProgearWeightTicket.Weight, progearWeightTicket.Weight)
	})
}
