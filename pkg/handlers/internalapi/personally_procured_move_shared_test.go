package internalapi

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) setupPersonallyProcuredMoveSharedTest(orderID uuid.UUID) {
	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}
	suite.MustSave(&address)

	locationName := "New Duty Location"
	location := models.DutyLocation{
		Name:      locationName,
		AddressID: address.ID,
		Address:   address,
	}
	suite.MustSave(&location)

	_ = testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:                orderID,
			NewDutyLocationID: location.ID,
		},
	})
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
