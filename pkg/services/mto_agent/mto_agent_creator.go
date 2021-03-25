package mtoagent

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mtoAgentCreator sets up the service object
type mtoAgentCreator struct {
	db                  *pop.Connection
	availabilityChecker services.MoveTaskOrderChecker
}

// NewMTOAgentCreator creates a new struct with the service dependencies
func NewMTOAgentCreator(db *pop.Connection) services.MTOAgentCreator {
	return &mtoAgentCreator{
		db: db,
	}
}

// #TODO: Will this only be used by Prime?
// I believe so, which means we only need to
// worry about validating for Prime
func (f *mtoAgentCreator) CreateMTOAgentPrime(mtoAgent *models.MTOAgent) (*models.MTOAgent, error) {
	return f.CreateMTOAgent(mtoAgent, "prime")
}

func (f *mtoAgentCreator) CreateMTOAgent(mtoAgent *models.MTOAgent, validatorKey string) (*models.MTOAgent, error) {
	fmt.Println("🧩Inside CreateMTOAgent")
	fmt.Println(mtoAgent.MTOShipmentID)

	mtoShipment := &models.MTOShipment{}

	err := f.db.Find(mtoShipment, mtoAgent.MTOShipmentID)

	if err != nil {
		fmt.Println("No valid shipment found")
		return nil, nil
	}

	// #TODO: Add validation with udpated validation pattern
	// Existing fn to check if move is available to prime
	// #TODO Should not be able to create a new agent unless none exists
	// If one exists, return error message saying to use Update endpoint

	// Create MTOAgent
	err = f.db.Create(mtoAgent)

	if err != nil {
		fmt.Println("Error creating Agent")
	}

	fmt.Println("Created???")

	return mtoAgent, nil

}
