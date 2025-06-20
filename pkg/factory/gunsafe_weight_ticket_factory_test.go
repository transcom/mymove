package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildGunSafeWeightTicket() {
	const defaultDescription = "gun safe"
	const defaultWeight = unit.Pound(500)
	suite.Run("Successful creation of weight ticket ", func() {
		// Under test:      BuildGunSafeWeightTicket
		// Mocked:          None
		// Set up:          Create a weight ticket with no customizations or traits
		// Expected outcome:gunSafeWeightTicket should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		gunSafeWeightTicket := BuildGunSafeWeightTicket(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(gunSafeWeightTicket.Description)
		suite.Equal(defaultDescription, *gunSafeWeightTicket.Description)

		suite.NotNil(gunSafeWeightTicket.Weight)
		suite.Equal(defaultWeight, *gunSafeWeightTicket.Weight)

		suite.NotNil(gunSafeWeightTicket.HasWeightTickets)
		suite.True(*gunSafeWeightTicket.HasWeightTickets)

		suite.False(gunSafeWeightTicket.PPMShipmentID.IsNil())
		suite.False(gunSafeWeightTicket.PPMShipment.ID.IsNil())

		suite.False(gunSafeWeightTicket.DocumentID.IsNil())
		suite.False(gunSafeWeightTicket.Document.ID.IsNil())
		suite.NotEmpty(gunSafeWeightTicket.Document.UserUploads)

		serviceMemberID := gunSafeWeightTicket.PPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID
		suite.False(serviceMemberID.IsNil())
		suite.Equal(serviceMemberID, gunSafeWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Successful creation of customized GunSafeWeightTicket", func() {
		// Under test:      BuildGunSafeWeightTicket
		// Mocked:          None
		// Set up:          Create a weight ticket with and pass custom fields
		// Expected outcome:gunSafeWeightTicket should be created with custom values

		// SETUP
		customPPMShipment := models.PPMShipment{
			ExpectedDepartureDate: time.Now(),
			ActualMoveDate:        models.TimePointer(time.Now().Add(time.Duration(24 * time.Hour))),
		}
		customGunSafeWeightTicket := models.GunSafeWeightTicket{
			Description: models.StringPointer("VIP gunSafe"),
			Weight:      models.PoundPointer(512),
			Reason:      models.StringPointer("VIP Reason"),
		}
		customs := []Customization{
			{
				Model: customPPMShipment,
			},
			{
				Model: customGunSafeWeightTicket,
			},
		}
		// CALL FUNCTION UNDER TEST
		gunSafeWeightTicket := BuildGunSafeWeightTicket(suite.DB(), customs, nil)

		// VALIDATE RESULTS
		suite.NotNil(gunSafeWeightTicket.Description)
		suite.Equal(*customGunSafeWeightTicket.Description,
			*gunSafeWeightTicket.Description)

		suite.NotNil(gunSafeWeightTicket.Weight)
		suite.Equal(*customGunSafeWeightTicket.Weight, *gunSafeWeightTicket.Weight)

		suite.NotNil(gunSafeWeightTicket.Reason)
		suite.Equal(*customGunSafeWeightTicket.Reason,
			*gunSafeWeightTicket.Reason)

		suite.False(gunSafeWeightTicket.PPMShipmentID.IsNil())
		suite.False(gunSafeWeightTicket.PPMShipment.ID.IsNil())
		suite.Equal(customPPMShipment.ExpectedDepartureDate,
			gunSafeWeightTicket.PPMShipment.ExpectedDepartureDate)
		suite.NotNil(gunSafeWeightTicket.PPMShipment.ActualMoveDate)
		suite.Equal(*customPPMShipment.ActualMoveDate,
			*gunSafeWeightTicket.PPMShipment.ActualMoveDate)
	})

	suite.Run("Successful return of linkOnly GunSafeWeightTicket", func() {
		// Under test:       BuildGunSafeWeightTicket
		// Set up:           Pass in a linkOnly gunSafeWeightTicket
		// Expected outcome: No new GunSafeWeightTicket should be created.

		// Check num GunSafeWeightTicket records
		precount, err := suite.DB().Count(&models.GunSafeWeightTicket{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		gunSafeWeightTicket := BuildGunSafeWeightTicket(suite.DB(), []Customization{
			{
				Model: models.GunSafeWeightTicket{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.GunSafeWeightTicket{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, gunSafeWeightTicket.ID)
	})

	suite.Run("Successful return of stubbed GunSafeWeightTicket", func() {
		// Under test:       BuildGunSafeWeightTicket
		// Set up:           Pass in nil db
		// Expected outcome: No new GunSafeWeightTicket should be created.

		// Check num GunSafeWeightTicket records
		precount, err := suite.DB().Count(&models.GunSafeWeightTicket{})
		suite.NoError(err)

		customGunSafeWeightTicket := models.GunSafeWeightTicket{
			Weight: models.PoundPointer(9999),
		}
		// Nil passed in as db
		gunSafeWeightTicket := BuildGunSafeWeightTicket(nil, []Customization{
			{
				Model: customGunSafeWeightTicket,
			},
		}, nil)

		count, err := suite.DB().Count(&models.GunSafeWeightTicket{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customGunSafeWeightTicket.Weight, gunSafeWeightTicket.Weight)
	})
}
