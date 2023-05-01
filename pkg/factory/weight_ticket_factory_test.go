package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildWeightTicket() {
	const defaultVehicleDescription = "2022 Honda CR-V Hybrid"
	const defaultEmptyWeight = unit.Pound(14500)
	const defaultFullWeight = defaultEmptyWeight + unit.Pound(4000)
	suite.Run("Successful creation of weight ticket ", func() {
		// Under test:      BuildWeightTicket
		// Mocked:          None
		// Set up:          Create a weight ticket with no customizations or traits
		// Expected outcome:weightTicket should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		weightTicket := BuildWeightTicket(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(weightTicket.VehicleDescription)
		suite.Equal(defaultVehicleDescription, *weightTicket.VehicleDescription)

		suite.NotNil(weightTicket.EmptyWeight)
		suite.Equal(defaultEmptyWeight, *weightTicket.EmptyWeight)

		suite.NotNil(weightTicket.MissingEmptyWeightTicket)
		suite.False(*weightTicket.MissingEmptyWeightTicket)

		suite.NotNil(weightTicket.FullWeight)
		suite.Equal(defaultFullWeight, *weightTicket.FullWeight)

		suite.NotNil(weightTicket.MissingFullWeightTicket)
		suite.False(*weightTicket.MissingFullWeightTicket)

		suite.NotNil(weightTicket.OwnsTrailer)
		suite.False(*weightTicket.OwnsTrailer)

		suite.NotNil(weightTicket.TrailerMeetsCriteria)
		suite.False(*weightTicket.TrailerMeetsCriteria)

		suite.False(weightTicket.PPMShipmentID.IsNil())
		suite.False(weightTicket.PPMShipment.ID.IsNil())

		suite.False(weightTicket.EmptyDocumentID.IsNil())
		suite.False(weightTicket.EmptyDocument.ID.IsNil())
		suite.NotEmpty(weightTicket.EmptyDocument.UserUploads)

		suite.False(weightTicket.FullDocumentID.IsNil())
		suite.False(weightTicket.FullDocument.ID.IsNil())
		suite.NotEmpty(weightTicket.FullDocument.UserUploads)

		suite.False(weightTicket.ProofOfTrailerOwnershipDocumentID.IsNil())
		suite.False(weightTicket.ProofOfTrailerOwnershipDocument.ID.IsNil())
		suite.NotEmpty(weightTicket.ProofOfTrailerOwnershipDocument.UserUploads)

		serviceMemberID := weightTicket.PPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID
		suite.False(serviceMemberID.IsNil())
		suite.Equal(serviceMemberID, weightTicket.EmptyDocument.ServiceMemberID)
		suite.Equal(serviceMemberID, weightTicket.FullDocument.ServiceMemberID)
		suite.Equal(serviceMemberID, weightTicket.ProofOfTrailerOwnershipDocument.ServiceMemberID)
	})

	suite.Run("Successful creation of customized WeightTicket", func() {
		// Under test:      BuildWeightTicket
		// Mocked:          None
		// Set up:          Create a weight ticket with and pass custom fields
		// Expected outcome:weightTicket should be created with custom values

		// SETUP
		customPPMShipment := models.PPMShipment{
			ExpectedDepartureDate: time.Now(),
			ActualMoveDate:        models.TimePointer(time.Now().Add(time.Duration(24 * time.Hour))),
		}
		customWeightTicket := models.WeightTicket{
			VehicleDescription: models.StringPointer("1975 Ferrari Dino"),
			EmptyWeight:        models.PoundPointer(12345),
			FullWeight:         models.PoundPointer(16543),
		}
		customs := []Customization{
			{
				Model: customPPMShipment,
			},
			{
				Model: customWeightTicket,
			},
		}
		// CALL FUNCTION UNDER TEST
		weightTicket := BuildWeightTicket(suite.DB(), customs, nil)

		// VALIDATE RESULTS
		suite.NotNil(weightTicket.VehicleDescription)
		suite.Equal(*customWeightTicket.VehicleDescription,
			*weightTicket.VehicleDescription)

		suite.NotNil(weightTicket.EmptyWeight)
		suite.Equal(*customWeightTicket.EmptyWeight, *weightTicket.EmptyWeight)
		suite.NotNil(weightTicket.FullWeight)
		suite.Equal(*customWeightTicket.FullWeight, *weightTicket.FullWeight)

		suite.False(weightTicket.PPMShipmentID.IsNil())
		suite.False(weightTicket.PPMShipment.ID.IsNil())
		suite.Equal(customPPMShipment.ExpectedDepartureDate,
			weightTicket.PPMShipment.ExpectedDepartureDate)
		suite.NotNil(weightTicket.PPMShipment.ActualMoveDate)
		suite.Equal(*customPPMShipment.ActualMoveDate,
			*weightTicket.PPMShipment.ActualMoveDate)
	})

	suite.Run("Successful return of linkOnly WeightTicket", func() {
		// Under test:       BuildWeightTicket
		// Set up:           Pass in a linkOnly weightTicket
		// Expected outcome: No new WeightTicket should be created.

		// Check num WeightTicket records
		precount, err := suite.DB().Count(&models.WeightTicket{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		weightTicket := BuildWeightTicket(suite.DB(), []Customization{
			{
				Model: models.WeightTicket{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.WeightTicket{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, weightTicket.ID)
	})

	suite.Run("Successful return of stubbed WeightTicket", func() {
		// Under test:       BuildWeightTicket
		// Set up:           Pass in nil db
		// Expected outcome: No new WeightTicket should be created.

		// Check num WeightTicket records
		precount, err := suite.DB().Count(&models.WeightTicket{})
		suite.NoError(err)

		customWeightTicket := models.WeightTicket{
			EmptyWeight: models.PoundPointer(9999),
		}
		// Nil passed in as db
		weightTicket := BuildWeightTicket(nil, []Customization{
			{
				Model: customWeightTicket,
			},
		}, nil)

		count, err := suite.DB().Count(&models.WeightTicket{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customWeightTicket.EmptyWeight, weightTicket.EmptyWeight)
	})
	suite.Run("Successful creation of weight ticket with constructed weight ", func() {
		// Under test:      BuildWeightTicketWithConstructedWeight
		// Mocked:          None
		// Set up:          Create a weight ticket with constructed weight with no customizations or traits
		// Expected outcome:weightTicket should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		weightTicket := BuildWeightTicketWithConstructedWeight(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(weightTicket.MissingEmptyWeightTicket)
		suite.True(*weightTicket.MissingEmptyWeightTicket)

		suite.NotNil(weightTicket.MissingFullWeightTicket)
		suite.True(*weightTicket.MissingFullWeightTicket)

		suite.False(weightTicket.ProofOfTrailerOwnershipDocumentID.IsNil())
		suite.False(weightTicket.ProofOfTrailerOwnershipDocument.ID.IsNil())
		suite.NotEmpty(weightTicket.ProofOfTrailerOwnershipDocument.UserUploads)

		serviceMemberID := weightTicket.PPMShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.ID
		suite.False(serviceMemberID.IsNil())
		suite.Equal(serviceMemberID, weightTicket.EmptyDocument.ServiceMemberID)
		suite.Equal(serviceMemberID, weightTicket.FullDocument.ServiceMemberID)
		suite.Equal(serviceMemberID, weightTicket.ProofOfTrailerOwnershipDocument.ServiceMemberID)
	})
}
