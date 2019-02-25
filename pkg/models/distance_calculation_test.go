package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_DistanceCalculationCreate() {
	t := suite.T()

	address1 := testdatagen.MakeDefaultAddress(suite.DB())
	address2 := testdatagen.MakeDefaultAddress(suite.DB())

	distanceCalculation := models.DistanceCalculation{
		OriginAddress:        address1,
		OriginAddressID:      address1.ID,
		DestinationAddress:   address2,
		DestinationAddressID: address2.ID,
		DistanceMiles:        1044,
	}

	verrs, err := suite.DB().ValidateAndSave(&distanceCalculation)

	if err != nil {
		t.Fatalf("could not save DistanceCalculation: %v", err)
	}

	if verrs.Count() != 0 {
		t.Errorf("did not expect validation errors: %v", verrs)
	}
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
	planner := route.NewTestingPlanner(1044)
	address1 := testdatagen.MakeDefaultAddress(suite.DB())
	address2 := testdatagen.MakeDefaultAddress(suite.DB())
	distanceCalculation, err := models.NewDistanceCalculation(planner, address1, address2)

	suite.NoError(err)
	suite.Equal(distanceCalculation.DistanceMiles, 1044)
	suite.Equal(distanceCalculation.OriginAddressID, address1.ID)
	suite.Equal(distanceCalculation.DestinationAddressID, address2.ID)
	// And it should not have been saved
	suite.Equal(uuid.Nil, distanceCalculation.ID)
}
