package mtoagent

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMTOAgentBasicValidator is the key for generic validation on the MTO Agent
const UpdateMTOAgentBasicValidator string = "UpdateMTOAgentBasicValidator"

// UpdateMTOAgentPrimeValidator is the key for validating the MTO Agent for the Prime contractor
const UpdateMTOAgentPrimeValidator string = "UpdateMTOAgentPrimeValidator"

// CreateMTOAgentPrimeValidator is the key for validating the MTO Agent for the Prime
const CreateMTOAgentPrimeValidator = "CreateMTOAgentPrimeValidator"

// UpdateMTOAgentValidators is the map connecting the constant keys to the correct validator
var UpdateMTOAgentValidators = map[string]updateMTOAgentValidator{
	UpdateMTOAgentBasicValidator: new(basicUpdateMTOAgentValidator),
	UpdateMTOAgentPrimeValidator: new(primeUpdateMTOAgentValidator),
}

type updateMTOAgentValidator interface {
	validate(agentData *updateMTOAgentData) error
}

// basicUpdateMTOAgentValidator is the type for validation that should happen no matter who uses this service object
type basicUpdateMTOAgentValidator struct{}

func (v *basicUpdateMTOAgentValidator) validate(agentData *updateMTOAgentData) error {
	err := agentData.checkShipmentID()
	if err != nil {
		return err
	}

	err = agentData.getVerrs()
	if err != nil {
		return err
	}

	return nil
}

// primeUpdateMTOAgentValidator is the type for validation that is just for updates from the Prime contractor
type primeUpdateMTOAgentValidator struct{}

func (v *primeUpdateMTOAgentValidator) validate(agentData *updateMTOAgentData) error {
	err := agentData.checkShipmentID()
	if err != nil {
		return err
	}

	err = agentData.checkPrimeAvailability()
	if err != nil {
		return err
	}

	err = agentData.checkContactInfo()
	if err != nil {
		return err
	}

	err = agentData.getVerrs()
	if err != nil {
		return err
	}

	return nil
}

// updateMTOAgentData represents the data needed to validate an update on an MTOAgent
type updateMTOAgentData struct {
	updatedAgent        models.MTOAgent
	oldAgent            models.MTOAgent
	availabilityChecker services.MoveTaskOrderChecker
	verrs               *validate.Errors
}

// checkShipmentID checks that the user didn't attempt to change the agent's shipment ID
func (v *updateMTOAgentData) checkShipmentID() error {
	if v.updatedAgent.MTOShipmentID != uuid.Nil && v.updatedAgent.MTOShipmentID != v.oldAgent.MTOShipmentID {
		v.verrs.Add("mtoShipmentID", "cannot be updated")
	}

	return nil
}

// checkPrimeAvailability checks that agent is connected to a Prime-available shipment
func (v *updateMTOAgentData) checkPrimeAvailability() error {
	isAvailable, err := v.availabilityChecker.MTOAvailableToPrime(v.oldAgent.MTOShipment.MoveTaskOrderID)

	if !isAvailable || err != nil {
		return services.NewNotFoundError(v.updatedAgent.ID, "while looking for Prime-available MTOAgent")
	}

	return nil
}

// checkContactInfo checks that the new agent has the minimum required contact info: First Name and one of Email or Phone
func (v *updateMTOAgentData) checkContactInfo() error {
	firstName := v.oldAgent.FirstName
	if v.updatedAgent.FirstName != nil {
		firstName = v.updatedAgent.FirstName
	}

	// Check that we have something in the FirstName field:
	if firstName == nil || *firstName == "" {
		v.verrs.Add("firstName", "cannot be blank")
	}

	email := v.oldAgent.Email
	if v.updatedAgent.Email != nil {
		email = v.updatedAgent.Email
	}

	phone := v.oldAgent.Phone
	if v.updatedAgent.Phone != nil {
		phone = v.updatedAgent.Phone
	}

	// Check that we have one method of contacting the agent:
	if (email == nil || *email == "") && (phone == nil || *phone == "") {
		v.verrs.Add("contactInfo", "agent must have at least one contact method provided")
	}

	return nil
}

// getVerrs looks for any validation errors and returns a formatted InvalidInputError if any are found.
// Should only be called after the other check methods have been called.
func (v *updateMTOAgentData) getVerrs() error {
	if v.verrs.HasAny() {
		return services.NewInvalidInputError(v.updatedAgent.ID, nil, v.verrs, "Invalid input found while validating the agent.")
	}

	return nil
}

// setNewMTOAgent compares updatedAgent and oldAgent and updates a new MTOAgent instance with all data
// (changed and unchanged) filled in. Does not return an error, data must be checked for validation before this step.
func (v *updateMTOAgentData) setNewMTOAgent() *models.MTOAgent {
	newAgent := v.oldAgent

	if v.updatedAgent.MTOAgentType != "" {
		newAgent.MTOAgentType = v.updatedAgent.MTOAgentType
	}
	if v.updatedAgent.FirstName != nil {
		newAgent.FirstName = v.updatedAgent.FirstName

		if *v.updatedAgent.FirstName == "" {
			newAgent.FirstName = nil
		}
	}
	if v.updatedAgent.LastName != nil {
		newAgent.LastName = v.updatedAgent.LastName

		if *v.updatedAgent.LastName == "" {
			newAgent.LastName = nil
		}
	}
	if v.updatedAgent.Email != nil {
		newAgent.Email = v.updatedAgent.Email

		if *v.updatedAgent.Email == "" {
			newAgent.Email = nil
		}
	}
	if v.updatedAgent.Phone != nil {
		newAgent.Phone = v.updatedAgent.Phone

		if *v.updatedAgent.Phone == "" {
			newAgent.Phone = nil
		}
	}

	return &newAgent
}
