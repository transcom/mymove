package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceAgent finds or makes a single service_agent record
func MakeServiceAgent(db *pop.Connection, assertions Assertions) models.ServiceAgent {

	// Create a shipment if one wasn't already created
	shipmentID := assertions.ServiceAgent.ShipmentID
	if isZeroUUID(shipmentID) {
		shipment := MakeDefaultShipment(db)
		shipmentID = shipment.ID
	}

	// Manage the role
	role := assertions.ServiceAgent.Role
	if role == models.Role("") {
		role = models.RoleORIGIN
	}

	poc := assertions.ServiceAgent.PointOfContact
	if poc == "" {
		poc = "Jenny at ACME Movers"
	}

	phone := assertions.ServiceAgent.PhoneNumber
	if phone == nil {
		phone = stringPointer("303-867-5309")
	}

	email := assertions.ServiceAgent.Email
	if email == nil {
		email = stringPointer("jenny_acme@example.com")
	}

	serviceAgent := models.ServiceAgent{
		ShipmentID:     shipmentID,
		Role:           role,
		PointOfContact: poc,
		PhoneNumber:    phone,
		Email:          email,
	}

	verrs, err := db.ValidateAndSave(&serviceAgent)
	if verrs.HasAny() {
		err = fmt.Errorf("serviceAgent validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return serviceAgent
}

// MakeDefaultServiceAgent makes a Service Agent with default values
func MakeDefaultServiceAgent(db *pop.Connection) models.ServiceAgent {
	return MakeServiceAgent(db, Assertions{})
}
