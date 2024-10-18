package payloads

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestFetchPPMShipment() {

	ppmShipmentID, _ := uuid.NewV4()
	streetAddress1 := "MacDill AFB"
	streetAddress2, streetAddress3 := "", ""
	city := "Tampa"
	state := "FL"
	postalcode := "33621"
	country := "US"
	county := "HILLSBOROUGH"

	expectedAddress := models.Address{
		StreetAddress1: streetAddress1,
		StreetAddress2: &streetAddress2,
		StreetAddress3: &streetAddress3,
		City:           city,
		State:          state,
		PostalCode:     postalcode,
		Country:        &country,
		County:         county,
	}

	isActualExpenseReimbursement := true

	expectedPPMShipment := models.PPMShipment{
		ID:                           ppmShipmentID,
		PickupAddress:                &expectedAddress,
		DestinationAddress:           &expectedAddress,
		IsActualExpenseReimbursement: &isActualExpenseReimbursement,
	}

	suite.Run("Success -", func() {
		returnedPPMShipment := PPMShipment(nil, &expectedPPMShipment)

		suite.IsType(&internalmessages.PPMShipment{}, returnedPPMShipment)
		suite.Equal(&streetAddress1, returnedPPMShipment.PickupAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress2, returnedPPMShipment.PickupAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.PickupAddress.StreetAddress3, returnedPPMShipment.PickupAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.PickupAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.PickupAddress.City)
		suite.Equal(&state, returnedPPMShipment.PickupAddress.State)
		suite.Equal(&country, returnedPPMShipment.PickupAddress.Country)
		suite.Equal(&county, returnedPPMShipment.PickupAddress.County)

		suite.Equal(&streetAddress1, returnedPPMShipment.DestinationAddress.StreetAddress1)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress2, returnedPPMShipment.DestinationAddress.StreetAddress2)
		suite.Equal(expectedPPMShipment.DestinationAddress.StreetAddress3, returnedPPMShipment.DestinationAddress.StreetAddress3)
		suite.Equal(&postalcode, returnedPPMShipment.DestinationAddress.PostalCode)
		suite.Equal(&city, returnedPPMShipment.DestinationAddress.City)
		suite.Equal(&state, returnedPPMShipment.DestinationAddress.State)
		suite.Equal(&country, returnedPPMShipment.DestinationAddress.Country)
		suite.Equal(&county, returnedPPMShipment.DestinationAddress.County)

		suite.True(*returnedPPMShipment.IsActualExpenseReimbursement)
	})
}
