package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_ServiceAgentValidations() {
	serviceAgent := &ServiceAgent{}

	expErrors := map[string][]string{
		"shipment_id": {"ShipmentID can not be blank."},
		"role":        {"Role can not be blank."},
		"company":     {"Company can not be blank."},
	}

	suite.verifyValidationErrors(serviceAgent, expErrors)
}

// Test_FetchServiceAgentsOnShipment tests that a shipment's service agents are returned when we fetch.
func (suite *ModelSuite) Test_FetchServiceAgentsOnShipment() {
	t := suite.T()
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	// Make 2 Service Agents on a shipment
	testdatagen.MakeServiceAgent(suite.DB(), testdatagen.Assertions{
		ServiceAgent: ServiceAgent{
			ShipmentID: shipment.ID,
			Shipment:   &shipment,
		},
	})
	testdatagen.MakeServiceAgent(suite.DB(), testdatagen.Assertions{
		ServiceAgent: ServiceAgent{
			Role:       RoleDESTINATION,
			ShipmentID: shipment.ID,
			Shipment:   &shipment,
		},
	})
	// And 1 Service Agent on a different shipment
	testdatagen.MakeDefaultServiceAgent(suite.DB())

	serviceAgents, err := FetchServiceAgentsOnShipment(suite.DB(), shipment.ID)

	// Expect 2 service agents returned
	if err != nil {
		t.Errorf("Failed to find Service Agents: %v", err)
	} else if len(serviceAgents) != 2 {
		t.Errorf("Returned incorrect number of service agents. Expected 2, got %v", len(serviceAgents))
	}
}
