package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildSITDurationUpdate() {
	suite.Run("Successful creation of default SITDurationUpdate", func() {
		// Under test:      BuildSITDurationUpdate
		// Mocked:          None
		// Set up:          Create an SITDurationUpdate with no customizations or traits
		// Expected outcome:SITDurationUpdate should be created with default values

		// SETUP
		// Create a default SITDurationUpdate to compare values
		defaultSIT := models.SITDurationUpdate{
			RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
			RequestedDays: 45,
			Status:        models.SITExtensionStatusPending,
		}

		// FUNCTION UNDER TEST
		sitDurationUpdate := BuildSITDurationUpdate(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultSIT.RequestReason, sitDurationUpdate.RequestReason)
		suite.Equal(defaultSIT.RequestedDays, sitDurationUpdate.RequestedDays)
		suite.Equal(defaultSIT.Status, sitDurationUpdate.Status)
		suite.NotNil(sitDurationUpdate.MTOShipment)
		suite.False(sitDurationUpdate.MTOShipment.ID.IsNil())
		suite.False(sitDurationUpdate.MTOShipmentID.IsNil())
	})

	suite.Run("Successful creation of customized SITDurationUpdate", func() {
		// Under test:      BuildSITDurationUpdate
		// Mocked:          None
		// Set up:          Create SITDurationUpdate with customization
		// Expected outcome:SITDurationUpdate should be created with customized values

		// SETUP
		customSIT := models.SITDurationUpdate{
			ID:                uuid.Must(uuid.NewV4()),
			RequestReason:     models.SITExtensionRequestReasonNonavailabilityOfCivilianHousing,
			ContractorRemarks: models.StringPointer("contractor remarks"),
			RequestedDays:     90,
			Status:            models.SITExtensionStatusDenied,
			DecisionDate:      models.TimePointer(time.Now()),
			OfficeRemarks:     models.StringPointer("office remarks"),
		}

		customShipment := models.MTOShipment{
			ID: uuid.Must(uuid.NewV4()),
		}

		// FUNCTION UNDER TEST
		sitDurationUpdate := BuildSITDurationUpdate(suite.DB(), []Customization{
			{Model: customSIT},
			{Model: customShipment},
		}, nil)

		// VALIDATE RESULTS
		suite.NotNil(sitDurationUpdate.MTOShipment)
		suite.Equal(customShipment.ID, sitDurationUpdate.MTOShipment.ID)
		suite.Equal(customShipment.ID, sitDurationUpdate.MTOShipmentID)
		suite.Equal(customSIT.ID, sitDurationUpdate.ID)
		suite.Equal(customSIT.RequestReason, sitDurationUpdate.RequestReason)
		suite.Equal(*customSIT.ContractorRemarks, *sitDurationUpdate.ContractorRemarks)
		suite.Equal(customSIT.RequestedDays, sitDurationUpdate.RequestedDays)
		suite.Equal(customSIT.Status, sitDurationUpdate.Status)
		suite.Equal(*customSIT.DecisionDate, *sitDurationUpdate.DecisionDate)
		suite.Equal(*customSIT.OfficeRemarks, *sitDurationUpdate.OfficeRemarks)
	})

	suite.Run("Successful return of linkOnly SITDurationUpdate", func() {
		// Under test:       BuildSITDurationUpdate
		// Set up:           Pass in a linkOnly SITDurationUpdate
		// Expected outcome: No new SITDurationUpdate should be created.

		// Check num SITDurationUpdates
		precount, err := suite.DB().Count(&models.SITDurationUpdate{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		sitDurationUpdate := BuildSITDurationUpdate(suite.DB(), []Customization{
			{
				Model: models.SITDurationUpdate{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(id, sitDurationUpdate.ID)

		// Count how many notification are in the DB, no new
		// SITDurationUpdate should have been created
		count, err := suite.DB().Count(&models.SITDurationUpdate{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of stubbed SITDurationUpdate", func() {
		// Under test:      BuildSITDurationUpdate
		// Set up:          Create a stubbed SITDurationUpdate
		// Expected outcome:No new SITDurationUpdate should be created

		// Check num SITDurationUpdates
		precount, err := suite.DB().Count(&models.SITDurationUpdate{})
		suite.NoError(err)

		sitDurationUpdate := BuildSITDurationUpdate(nil, nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(sitDurationUpdate.MTOShipment)
		suite.True(sitDurationUpdate.MTOShipment.ID.IsNil())
		suite.True(sitDurationUpdate.MTOShipmentID.IsNil())

		// Count how many notification are in the DB, no new
		// SITDurationUpdate should have been created
		count, err := suite.DB().Count(&models.SITDurationUpdate{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}
