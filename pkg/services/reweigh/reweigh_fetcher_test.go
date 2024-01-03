package reweigh

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ReweighSuite) TestListReweighsByShipmentIDs() {
	reweighFetcher := NewReweighFetcher()
	appCtx := suite.AppContextForTest()
	move := factory.BuildMove(suite.DB(), nil, nil)
	parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				UsesExternalVendor:     true,
				Diversion:              true,
				DivertedFromShipmentID: nil,
			},
		},
	}, nil)
	childShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				UsesExternalVendor:     true,
				Diversion:              true,
				DivertedFromShipmentID: &parentShipment.ID,
			},
		},
	}, nil)
	grandChildShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				UsesExternalVendor:     true,
				Diversion:              true,
				DivertedFromShipmentID: &childShipment.ID,
			},
		},
	}, nil)

	parentReweighModel := &models.Reweigh{
		RequestedAt: time.Now(),
		RequestedBy: models.ReweighRequesterPrime,
		ShipmentID:  parentShipment.ID,
	}
	childReweighModel := &models.Reweigh{
		RequestedAt: time.Now(),
		RequestedBy: models.ReweighRequesterPrime,
		ShipmentID:  childShipment.ID,
	}
	reweighCreator := NewReweighCreator()
	parentReweigh, err := reweighCreator.CreateReweighCheck(appCtx, parentReweighModel)
	suite.NoError(err)
	suite.NotNil(parentReweigh)
	childReweigh, err := reweighCreator.CreateReweighCheck(appCtx, childReweighModel)
	suite.NoError(err)
	suite.NotNil(childReweigh)
	// Leave grandchild with no reweigh

	reweighsMap, err := reweighFetcher.ListReweighsByShipmentIDs(suite.AppContextForTest(), []uuid.UUID{parentShipment.ID, childShipment.ID, grandChildShipment.ID})
	suite.NoError(err)

	// Ensure reweighs are correctly fetched
	suite.Equal(parentReweigh.ID, reweighsMap[parentShipment.ID].ID)
	suite.Equal(childReweigh.ID, reweighsMap[childShipment.ID].ID)
	// Ensure the grandchild is found because it gets a reweigh created when it's found to not have any
	_, exists := reweighsMap[grandChildShipment.ID]
	suite.True(exists)
}
