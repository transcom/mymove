package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceAgent finds or makes a single service_agent record
func MakeServiceAgent(db *pop.Connection, assertions Assertions) models.ServiceAgent {

	// Create a shipment if one wasn't already created
	shipment := assertions.ServiceAgent.Shipment
	if shipment == nil {
		s := MakeDefaultShipment(db)
		shipment = &s
	}

	company := assertions.ServiceAgent.Company
	if company == "" {
		company = "ACME Movers"
	}

	role := assertions.ServiceAgent.Role
	if role == "" {
		role = models.RoleORIGIN
	}

	serviceAgent := models.ServiceAgent{
		ShipmentID:  shipment.ID,
		Shipment:    shipment,
		Role:        role,
		Company:     company,
		PhoneNumber: stringPointer("303-867-5309"),
		Email:       stringPointer("acme@example.com"),
	}

	mergeModels(&serviceAgent, assertions.ServiceAgent)

	mustCreate(db, &serviceAgent)

	return serviceAgent
}

// MakeDefaultServiceAgent makes a Service Agent with default values
func MakeDefaultServiceAgent(db *pop.Connection) models.ServiceAgent {
	return MakeServiceAgent(db, Assertions{})
}
