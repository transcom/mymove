package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_DistanceCalculationCreate() {
	address1 := testdatagen.MakeStubbedAddress(suite.DB())
	address2 := testdatagen.MakeStubbedAddress(suite.DB())

	distanceCalculation := models.DistanceCalculation{
		OriginAddress:        address1,
		OriginAddressID:      address1.ID,
		DestinationAddress:   address2,
		DestinationAddressID: address2.ID,
		DistanceMiles:        1044,
	}

	verrs, err := distanceCalculation.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) Test_DistanceCalculationValidations() {
	distanceCalculation := &models.DistanceCalculation{}

	var expErrors = map[string][]string{
		"origin_address_id":      {"OriginAddressID can not be blank."},
		"destination_address_id": {"DestinationAddressID can not be blank."},
		"distance_miles":         {"DistanceMiles can not be blank."},
	}

	suite.verifyValidationErrors(distanceCalculation, expErrors)
}

func (suite *ModelSuite) Test_NewDistanceCalculationCallsPlanner() {
	address1 := testdatagen.MakeStubbedAddress(suite.DB())
	address2 := testdatagen.MakeStubbedAddress(suite.DB())
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		address1.PostalCode,
		address2.PostalCode,
	).Return(1044, nil)
	useZipOnlyForDistance := true
	distanceCalculation, err := models.NewDistanceCalculation(planner, address1, address2, useZipOnlyForDistance)

	suite.NoError(err)
	suite.Equal(distanceCalculation.DistanceMiles, 1044)
	suite.Equal(distanceCalculation.OriginAddressID, address1.ID)
	suite.Equal(distanceCalculation.DestinationAddressID, address2.ID)
	// And it should not have been saved
	suite.Equal(uuid.Nil, distanceCalculation.ID)
}
