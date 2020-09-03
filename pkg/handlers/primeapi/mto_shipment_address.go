package primeapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
)

// UpdateMTOShipmentAddressHandler is the handler to update an address
type UpdateMTOShipmentAddressHandler struct {
	handlers.HandlerContext
	MTOShipmentAddressUpdater
	// mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// MTOShipmentAddressUpdater handles the db connection
type MTOShipmentAddressUpdater struct {
	db                     *pop.Connection
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// NewMTOShipmentAddressUpdater updates the address for an MTO Shipment
func NewMTOShipmentAddressUpdater(db *pop.Connection) *MTOShipmentAddressUpdater {
	return &MTOShipmentAddressUpdater{
		db: db}
}

// StaleIdentifierError is used when optimistic locking determines that the identifier refers to stale data
type StaleIdentifierError struct {
	StaleIdentifier string
}

func (e StaleIdentifierError) Error() string {
	return fmt.Sprintf("stale identifier: %s", e.StaleIdentifier)
}

// ValidateAddress does some validation #TODO rename to validate shipment. Should the receiver be the handler instead?
func (f MTOShipmentAddressUpdater) ValidateAddress(newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string) (bool, error) {
	// Use the existing service for checking mtoAvailableToPrime.
	// Instantiate on MTOShipmentAddressUpdater
	f.mtoAvailabilityChecker = movetaskorder.NewMoveTaskOrderChecker(f.db)
	// Find the mtoShipment based on id, so we can pull the uuid for the move
	mtoShipment := models.MTOShipment{}
	err := f.db.Find(&mtoShipment, mtoShipmentID)
	if err != nil {
		return false, err
	}

	oldAddress := models.Address{}
	err = f.db.Find(&oldAddress, newAddress.ID)

	if err != nil {
		return false, err
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldAddress.UpdatedAt)

	if encodedUpdatedAt != eTag {
		return false, services.NewPreconditionFailedError(newAddress.ID, err)
	}

	// Find the move associated with the mtoShipment
	move := &models.Move{}
	err = f.db.Find(move, mtoShipment.MoveTaskOrderID)
	if err != nil {
		return false, err
	}

	// Make sure the associated move is available to the prime, otherwise
	// they should not be updating anything
	mtoAvailableToPrime, err := f.mtoAvailabilityChecker.MTOAvailableToPrime(move.ID)
	if !mtoAvailableToPrime {
		return false, err // #TODO: return proper error msg
	}

	// Gather existing addressIDs for the shipment and see if our ID
	// matches one of them
	addressIDs := []*uuid.UUID{
		mtoShipment.PickupAddressID,
		mtoShipment.DestinationAddressID,
		mtoShipment.SecondaryDeliveryAddressID,
		mtoShipment.SecondaryPickupAddressID,
	}

	for _, id := range addressIDs {
		if id != nil {
			if *id == newAddress.ID {
				return true, nil
			}
		}

	}

	return false, err
}

// Handle updates the shipment
func (h UpdateMTOShipmentAddressHandler) Handle(params mtoshipmentops.UpdateMTOShipmentAddressParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	payload := params.Body
	eTag := params.IfMatch
	mtoShipmentID := params.MtoShipmentID.String()
	addressID := params.AddressID.String()

	if payload == nil {
		logger.Error("Error")
		return mtoshipmentops.NewUpdateMTOShipmentAddressBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"The MTO Shipment Address request body cannot be empty.", h.GetTraceID()))
	}

	newAddress := payloads.AddressModel(payload)
	newAddress.ID = uuid.FromStringOrNil(addressID)

	isValidated, err := h.MTOShipmentAddressUpdater.ValidateAddress(newAddress, uuid.FromStringOrNil(mtoShipmentID), eTag)

	if !isValidated {
		fmt.Println(err) //TODO: replace with error msg
		return mtoshipmentops.NewUpdateMTOShipmentAddressBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"The MTO Shipment Address cannot be updated.", h.GetTraceID())) // #TODO: update this error msg
	}

	// Make the update
	err = h.db.Save(newAddress) //#TODO: Add validation?

	if err != nil {
		fmt.Println("There was err")
		fmt.Println(err)
	} else {

		mtoShipmentAddressPayload := payloads.Address(newAddress)
		return mtoshipmentops.NewUpdateMTOShipmentAddressOK().WithPayload(mtoShipmentAddressPayload)
	}

	return nil
}
