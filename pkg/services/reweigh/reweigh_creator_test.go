package reweigh

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ReweighSuite) TestReweighCreator() {
	// Create new mtoShipment

	suite.Run("CreateReweigh - Success", func() {
		// Under test:	CreateReweigh
		// Set up:		Established valid shipment and valid reweigh
		// Expected:	New reweigh successfully created
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		// Create a valid reweigh for the move
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  mtoShipment.ID,
		}
		reweighCreator := NewReweighCreator()
		createdReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)

		suite.Nil(err)
		suite.NotNil(createdReweigh)
		suite.Equal(mtoShipment.ID, createdReweigh.ShipmentID)

	})

	// InvalidInputError
	suite.Run("Reweigh with validation errors returns an InvalidInputError", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		// Create a reweigh with a bad requester
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequester("not requested by anyone"),
			ShipmentID:  mtoShipment.ID,
		}
		reweighCreator := NewReweighCreator()
		createReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)

		suite.Error(err)
		suite.Nil(createReweigh)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("Not Found Error", func() {

		// Create a reweigh with a shipment that doesn't exist
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  uuid.Must(uuid.NewV4()),
		}
		reweighCreator := NewReweighCreator()
		createdReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)

		suite.Nil(createdReweigh)
		suite.IsType(apperror.NotFoundError{}, err)
	})
	suite.Run("Create reweigh for a single diverted shipment", func() {
		parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: nil,
				},
			},
		}, nil)
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  parentShipment.ID,
		}
		reweighCreator := NewReweighCreator()
		createdReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)
		suite.NoError(err)
		suite.NotNil(createdReweigh)
		suite.Equal(parentShipment.ID, createdReweigh.ShipmentID)
	})
	suite.Run("Create reweighs for diverted shipment chain, starting with the child", func() {
		parentShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)
		childShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: &parentShipment.ID,
				},
			},
		}, nil)
		newReweigh := &models.Reweigh{
			RequestedAt: time.Now(),
			RequestedBy: models.ReweighRequesterPrime,
			ShipmentID:  childShipment.ID,
		}
		reweighCreator := NewReweighCreator()
		// Create reweigh for the child
		createdReweigh, err := reweighCreator.CreateReweighCheck(suite.AppContextForTest(), newReweigh)
		suite.NoError(err)
		suite.NotNil(createdReweigh)
		suite.Equal(childShipment.ID, createdReweigh.ShipmentID)
		// Verify that a reweigh was also created for the parent shipment
		reweighFetcher := NewReweighFetcher()
		reweighMap, err := reweighFetcher.ListReweighsByShipmentIDs(suite.AppContextForTest(), []uuid.UUID{parentShipment.ID})
		suite.NoError(err)
		suite.NotNil(reweighMap)
		parentReweigh, exists := reweighMap[parentShipment.ID]
		suite.True(exists)
		suite.Equal(parentShipment.ID, parentReweigh.ShipmentID)
	})
}
