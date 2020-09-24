package mtoagent

import (
	"fmt"

	"github.com/getlantern/deepcopy"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// mtoAgentUpdater handles the db connection
type mtoAgentUpdater struct {
	db *pop.Connection
}

// NewMTOAgentUpdater creates a new struct with the service dependencies
func NewMTOAgentUpdater(db *pop.Connection) services.MTOAgentUpdater {
	return &mtoAgentUpdater{
		db: db,
	}
}

// UpdateMTOAgent updates the MTO Agent
func (f *mtoAgentUpdater) UpdateMTOAgent(mtoAgent *models.MTOAgent, eTag string, validatorKey string) (*models.MTOAgent, error) {
	oldAgent := models.MTOAgent{}

	// Find the agent, return error if not found
	err := f.db.Eager("MTOShipment").Find(&oldAgent, mtoAgent.ID)
	if err != nil {
		return nil, services.NewNotFoundError(mtoAgent.ID, "while looking for MTOAgent")
	}

	checker := movetaskorder.NewMoveTaskOrderChecker(f.db)
	agentData := updateMTOAgentData{
		updatedAgent:        *mtoAgent,
		oldAgent:            oldAgent,
		availabilityChecker: checker,
	}

	newAgent, err := validateUpdateMTOAgent(&agentData, validatorKey)
	if err != nil {
		return nil, err
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldAgent.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, services.NewPreconditionFailedError(newAgent.ID, nil)
	}

	// Make the update and create a InvalidInputError if there were validation issues
	verrs, err := f.db.ValidateAndSave(newAgent)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(newAgent.ID, err, verrs, "Invalid input found while updating the agent.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("MTOAgent", err, "")
	}

	// Get the updated address and return
	updatedAgent := models.MTOAgent{}
	err = f.db.Find(&updatedAgent, newAgent.ID)
	if err != nil {
		return nil, services.NewQueryError("MTOAgent", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}
	return &updatedAgent, nil
}

func validateUpdateMTOAgent(agentData *updateMTOAgentData, validatorKey string) (*models.MTOAgent, error) {
	var newAgent models.MTOAgent

	if validatorKey == "" {
		validatorKey = UpdateMTOAgentBaseValidator
	}
	validator, ok := validators[validatorKey]
	if !ok {
		err := fmt.Errorf("validator key %s was not found in update MTO Agent validators", validatorKey)
		return nil, err
	}
	err := validator.validate(agentData)
	if err != nil {
		return nil, err
	}
	err = agentData.setNewMTOAgent(&newAgent)
	if err != nil {
		return nil, err
	}

	return &newAgent, nil
}

// UpdateMTOAgentBaseValidator is the key for generic validation on the MTO Agent
const UpdateMTOAgentBaseValidator string = "UpdateMTOAgentBaseValidator"

// UpdateMTOAgentPrimeValidator is the key for validating the MTO Agent for the Prime contractor
const UpdateMTOAgentPrimeValidator string = "UpdateMTOAgentPrimeValidator"

var validators = map[string]updateMTOAgentValidator{
	UpdateMTOAgentBaseValidator:  new(baseUpdateMTOAgentValidator),
	UpdateMTOAgentPrimeValidator: new(primeUpdateMTOAgentValidator),
}

type updateMTOAgentValidator interface {
	validate(agentData *updateMTOAgentData) error
}

// baseUpdateMTOAgentValidator is the type for validation that should happen no matter who uses this service object
type baseUpdateMTOAgentValidator struct{}

func (v *baseUpdateMTOAgentValidator) validate(agentData *updateMTOAgentData) error {
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

// updateMTOAgentData represents the data needed to validate the update on the MTOAgent
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
	if *firstName == "" || firstName == nil {
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
	if (*email == "" || email == nil) && (*phone == "" || phone == nil) {
		v.verrs.Add("contactInfo", "agent must have at least one contact method provided")
	}

	return nil
}

// getVerrs looks for any validation errors and returns a formatted InvalidInputError if any are found.
// Should only be called after the other check methods have been called.
func (v *updateMTOAgentData) getVerrs() error {
	if v.verrs.HasAny() {
		return services.NewInvalidInputError(v.updatedAgent.ID, nil, v.verrs, "Invalid input found while updating the agent.")
	}

	return nil
}

// setNewMTOAgent compares updatedAgent and oldAgent and updates a new MTOAgent instance with all data
// (changed and unchanged) filled in
func (v *updateMTOAgentData) setNewMTOAgent(newAgent *models.MTOAgent) error {
	err := deepcopy.Copy(&newAgent, &v.oldAgent)
	if err != nil {
		return fmt.Errorf("error copying agent data %w", err)
	}

	if v.updatedAgent.MTOAgentType != "" {
		newAgent.MTOAgentType = v.updatedAgent.MTOAgentType
	}
	if v.updatedAgent.FirstName != nil {
		newAgent.FirstName = v.updatedAgent.FirstName
	}
	if v.updatedAgent.LastName != nil {
		newAgent.LastName = v.updatedAgent.LastName
	}
	if v.updatedAgent.Email != nil {
		newAgent.Email = v.updatedAgent.Email
	}
	if v.updatedAgent.Phone != nil {
		newAgent.Phone = v.updatedAgent.Phone
	}

	return nil
}
