package internalapi

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) setupPersonallyProcuredMoveSharedTest(orderID uuid.UUID) {
	_ = factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				ID: orderID,
			},
		},
		{
			Model: models.DutyLocation{
				Name: "New Duty Location",
			},
			Type: &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Address{
				StreetAddress1: "some address",
				City:           "city",
				State:          "state",
				PostalCode:     "12345",
			},
			Type: &factory.Addresses.DutyLocationAddress,
		},
	}, nil)
}

func (suite *HandlerSuite) GetDestinationDutyLocationPostalCode() {
	orderID := uuid.Must(uuid.NewV4())
	invalidID, _ := uuid.FromString("00000000-0000-0000-0000-000000000000")
	suite.setupPersonallyProcuredMoveSharedTest(orderID)

	tests := []struct {
		lookupID  uuid.UUID
		resultZip string
		resultErr error
	}{
		{lookupID: orderID, resultZip: "12345", resultErr: nil},
		{lookupID: invalidID, resultZip: "", resultErr: models.ErrFetchNotFound},
	}

	for _, ts := range tests {
		destinationZip, err := GetDestinationDutyLocationPostalCode(suite.AppContextForTest(), ts.lookupID)
		suite.Equal(ts.resultErr, err, "Wrong resultErr: %s", ts.lookupID)
		suite.Equal(ts.resultZip, destinationZip, "Wrong moveID: %s", ts.lookupID)
	}
}
