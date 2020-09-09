package primeapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

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
}

// MTOShipmentAddressUpdater handles the db connection
type MTOShipmentAddressUpdater struct {
	db *pop.Connection
}

// NewMTOShipmentAddressUpdater updates the address for an MTO Shipment
func NewMTOShipmentAddressUpdater(db *pop.Connection) MTOShipmentAddressUpdater {
	return MTOShipmentAddressUpdater{
		db: db}
}

// isAddressOnShipment returns true if address is associated with the shipment, false if not
func isAddressOnShipment(address *models.Address, mtoShipment *models.MTOShipment) bool {
	addressIDs := []*uuid.UUID{
		mtoShipment.PickupAddressID,
		mtoShipment.DestinationAddressID,
		mtoShipment.SecondaryDeliveryAddressID,
		mtoShipment.SecondaryPickupAddressID,
	}

	for _, id := range addressIDs {
		if id != nil {
			if *id == address.ID {
				return true
			}
		}
	}
	return false
}

// ValidateAddress does some validation.
// MYTODO: Should the receiver be the handler instead?
func (f MTOShipmentAddressUpdater) ValidateAddress(newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string) (bool, error) {

	// Find the mtoShipment based on id, so we can pull the uuid for the move
	mtoShipment := models.MTOShipment{}
	oldAddress := models.Address{}

	// Find the shipment, return error if not found
	err := f.db.Find(&mtoShipment, mtoShipmentID)
	if err != nil {
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return false, services.NewNotFoundError(mtoShipmentID, "looking for mtoShipment")
		}
	}

	// Make sure the associated move is available to the prime
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(f.db)
	mtoAvailableToPrime, _ := mtoChecker.MTOAvailableToPrime(mtoShipment.MoveTaskOrderID)
	if !mtoAvailableToPrime {
		return false, services.NewNotFoundError(mtoShipment.MoveTaskOrderID, "looking for moveTaskOrder")
	}

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldAddress.UpdatedAt)
	fmt.Println(encodedUpdatedAt)

	// MYTODO: Revert etag check!
	// if encodedUpdatedAt != eTag {
	// 	return false, services.NewPreconditionFailedError(newAddress.ID, err)
	// }

	// Find the address, return error if not found
	err = f.db.Find(&oldAddress, newAddress.ID)
	if err != nil {
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return false, services.NewNotFoundError(newAddress.ID, "looking for address")
		}
	}

	// Check that address is associated with this shipment
	if isAddressOnShipment(newAddress, &mtoShipment) {
		return true, nil
	}
	err = services.NewConflictError(newAddress.ID, ": Address is not associated with the provided MTOShipmentID.")
	return false, err
}

// Handle updates the shipment
func (h UpdateMTOShipmentAddressHandler) Handle(params mtoshipmentops.UpdateMTOShipmentAddressParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	// Get the params and payload
	payload := params.Body
	eTag := params.IfMatch
	mtoShipmentID := params.MtoShipmentID.String()
	addressID := params.AddressID.String()

	// Get the new address model
	newAddress := payloads.AddressModel(payload)
	newAddress.ID = uuid.FromStringOrNil(addressID)

	// Validate the issues
	isValidated, err := h.MTOShipmentAddressUpdater.ValidateAddress(newAddress, uuid.FromStringOrNil(mtoShipmentID), eTag)

	var verrs *validate.Errors
	if isValidated {
		// Make the update and create a InvalidInput Error if there were validation issues
		verrs, err = h.db.ValidateAndSave(newAddress)
		// If there were validation errors create an InvalidInputError type
		if verrs != nil && verrs.HasAny() {
			logger.Error("Error validatating address: ", zap.Error(verrs))
			err = services.NewInvalidInputError(newAddress.ID, err, verrs, "")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			// This wraps the original error and gets handled properly in the switch statement
			err = services.NewQueryError("Address", err, "")
		}
	}

	if err != nil {
		logger.Error("primeapi.UpdateMTOShipmentAddressHandler", zap.Error(err))

		switch e := err.(type) {
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		// Not Found Error -> Not Found Response
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		// InvalidInputError -> Unprocessable Entity Response
		case services.InvalidInputError:
			return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(
				payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		// ConflictError -> ConflictError Response
		case services.ConflictError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressConflict().WithPayload(
				payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		// QueryError -> Internal Server Error
		case services.QueryError:
			if e.Unwrap() != nil {
				logger.Error("primeapi.UpdateMTOShipmentAddressHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewUpdateMTOShipmentAddressInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		// Unknown -> Internal Server Error
		default:
			return mtoshipmentops.NewUpdateMTOShipmentAddressInternalServerError().
				WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}

	}

	mtoShipmentAddressPayload := payloads.Address(newAddress)
	return mtoshipmentops.NewUpdateMTOShipmentAddressOK().WithPayload(mtoShipmentAddressPayload)

}
