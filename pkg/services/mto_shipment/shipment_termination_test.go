package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MTOShipmentServiceSuite) TestTerminateShipment() {
	terminator := NewShipmentTermination()

	suite.Run("If the shipment is terminated successfully, it should update the shipment status, terminated_at, and termination_comments", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		terminatedShipment, err := terminator.TerminateShipment(session, shipment.ID, "get in the choppuh")
		suite.NoError(err)
		suite.Equal(shipment.ID, terminatedShipment.ID)

		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusTerminatedForCause, fetchedShipment.Status)
		suite.Equal("TERMINATED FOR CAUSE - get in the choppuh", *fetchedShipment.TerminationComments)
		suite.NotNil(fetchedShipment.TerminatedAt)
	})

	suite.Run("Returns NotFoundError if shipment does not exist", func() {
		invalidShipmentID := uuid.Must(uuid.NewV4())
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		terminatedShipment, err := terminator.TerminateShipment(session, invalidShipmentID, "doesn't matter")
		suite.Error(err)
		suite.Nil(terminatedShipment)

		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns invalid input error if shipment has an actual pickup date and a termination is attempted", func() {
		now := time.Now()
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					ActualPickupDate: &now,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		terminatedShipment, err := terminator.TerminateShipment(session, shipment.ID, "get in the choppuh")
		suite.Error(err)
		suite.Nil(terminatedShipment)

		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Shipment cannot have an actual pickup date set in order to terminate for cause", err.Error())
	})

	suite.Run("Returns invalid input error if shipment is in a non-approved status and a termination is attempted", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		terminatedShipment, err := terminator.TerminateShipment(session, shipment.ID, "get in the choppuh")
		suite.Error(err)
		suite.Nil(terminatedShipment)

		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Shipment must be in APPROVED status in order to terminate for cause", err.Error())
	})

	suite.Run("Returns invalid input error if shipment is already in TERMINATED_FOR_CAUSE status and a termination is attempted", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		// Append terminated after the factory is done because the factory does lots of saves, triggering
		// the db trigger protection early
		shipment.Status = models.MTOShipmentStatusTerminatedForCause
		err := suite.DB().Save(&shipment)
		suite.FatalNoError(err)

		// Now try to terminate while already terminated
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		terminatedShipment, err := terminator.TerminateShipment(session, shipment.ID, "get in the choppuh")
		suite.Error(err)
		suite.Nil(terminatedShipment)

		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Shipment in TERMINATED FOR CAUSE status cannot be terminated for cause again", err.Error())
	})

	suite.Run(("Won't allow termination of a shipment tied to a PPM"), func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		// Fetch `mto_shipments` entry
		var mtoShipment models.MTOShipment
		err := suite.DB().Where("id = ?", ppmShipment.ShipmentID).First(&mtoShipment)
		suite.NoError(err)
		suite.NotEmpty(mtoShipment)

		// Make sure it's approved
		mtoShipment.Status = models.MTOShipmentStatusApproved
		err = suite.DB().Save(&mtoShipment)
		suite.FatalNoError(err)

		// Attempt to terminate the parent mto_shipment
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		terminatedShipment, err := terminator.TerminateShipment(session, mtoShipment.ID, "this will fail")
		suite.Error(err)
		suite.Nil(terminatedShipment)
		suite.EqualError(err, "Shipments tied to PPMs do not qualify for termination")
	})
}
