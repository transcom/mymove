package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
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
	db                     *pop.Connection
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// NewMTOShipmentAddressUpdater updates the address for an MTO Shipment
func NewMTOShipmentAddressUpdater(db *pop.Connection) *MTOShipmentAddressUpdater {
	return &MTOShipmentAddressUpdater{
		db: db}
}

// ValidateAddress does some validation. #TODO: Should the receiver be the handler instead?
func (f MTOShipmentAddressUpdater) ValidateAddress(newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string) (bool, error) {
	// Use the existing service for checking mtoAvailableToPrime.
	// Instantiate on MTOShipmentAddressUpdater
	f.mtoAvailabilityChecker = movetaskorder.NewMoveTaskOrderChecker(f.db)
	// Find the mtoShipment based on id, so we can pull the uuid for the move
	mtoShipment := models.MTOShipment{}
	oldAddress := models.Address{}

	err := f.db.Find(&mtoShipment, mtoShipmentID)

	if err != nil {
		//#TODO: What other types of errors need to be handled here?
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return false, services.NewNotFoundError(mtoShipmentID, "")
		}
	}

	err = f.db.Find(&oldAddress, newAddress.ID)

	if err != nil {
		//#TODO: What other types of errors need to be handled here?
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return false, services.NewNotFoundError(mtoShipmentID, "")
		}
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
		//#TODO: What other types of errors need to be handled here?
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return false, services.NewNotFoundError(mtoShipmentID, "")
		}
	}

	// Make sure the associated move is available to the prime, otherwise
	// they should not be updating anything
	mtoAvailableToPrime, _ := f.mtoAvailabilityChecker.MTOAvailableToPrime(move.ID)
	if !mtoAvailableToPrime {
		return false, services.NewNotFoundError(newAddress.ID, "")

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

	err = services.NewConflictError(newAddress.ID, "Address is not associated with the provided MTOShipmentID.")
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
		// #TODO: Confirm how to handle this validation error. Swagger validation returns the error,
		// so this would be backup validation.
		logger.Error("Invalid request: params Body is nil", zap.Any("payload", payload))
		return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), nil))
	}

	newAddress := payloads.AddressModel(payload)
	newAddress.ID = uuid.FromStringOrNil(addressID)

	isValidated, err := h.MTOShipmentAddressUpdater.ValidateAddress(newAddress, uuid.FromStringOrNil(mtoShipmentID), eTag)

	if isValidated {
		// Make the update
		err = h.db.Save(newAddress)
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
		case services.ConflictError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressConflict().WithPayload(
				payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
			// QueryError -> Internal Server Error
		case services.QueryError: // #TODO: When would this be used?
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
