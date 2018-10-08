package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForShipmentAccessorialModels(s []models.ShipmentAccessorial) apimessages.ShipmentAccessorials {
	payloads := make(apimessages.ShipmentAccessorials, len(s))

	for i, acc := range s {
		payloads[i] = payloadForShipmentAccessorialModel(&acc)
	}

	return payloads
}

func payloadForShipmentAccessorialModel(s *models.ShipmentAccessorial) *apimessages.ShipmentAccessorial {
	if s == nil {
		return nil
	}

	return &apimessages.ShipmentAccessorial{
		ID:            handlers.FmtUUID(s.ID),
		ShipmentID:    handlers.FmtUUID(s.ShipmentID),
		Accessorial:   payloadForAccessorialModel(&s.Accessorial),
		Location:      apimessages.AccessorialLocation(s.Location),
		Notes:         s.Notes,
		Quantity1:     handlers.FmtInt64(int64(s.Quantity1)),
		Quantity2:     handlers.FmtInt64(int64(s.Quantity2)),
		Status:        apimessages.AccessorialStatus(s.Status),
		SubmittedDate: *handlers.FmtDateTime(s.SubmittedDate),
		ApprovedDate:  *handlers.FmtDateTime(s.ApprovedDate),
	}
}

// GetShipmentAccessorialsHandler returns a particular shipment
type GetShipmentAccessorialsHandler struct {
	handlers.HandlerContext
}

// Handle returns a specified shipment
func (h GetShipmentAccessorialsHandler) Handle(params accessorialop.GetShipmentAccessorialsParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID := uuid.Must(uuid.FromString(params.ShipmentID.String()))

	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, _, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else if !session.IsOfficeUser() {
		return accessorialop.NewGetShipmentAccessorialsForbidden()
	}

	shipmentAccessorials, err := models.FetchAccessorialsByShipmentID(h.DB(), &shipmentID)
	if err != nil {
		h.Logger().Error("Error fetching accessorials for shipment", zap.Error(err))
		return accessorialop.NewGetShipmentAccessorialsInternalServerError()
	}
	payload := payloadForShipmentAccessorialModels(shipmentAccessorials)
	return accessorialop.NewGetShipmentAccessorialsOK().WithPayload(payload)
}

// CreateShipmentAccessorialHandler creates a shipment_accessorial for a provided shipment_id
type CreateShipmentAccessorialHandler struct {
	handlers.HandlerContext
}

// Handle handles the request
func (h CreateShipmentAccessorialHandler) Handle(params accessorialop.CreateShipmentAccessorialParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID := uuid.Must(uuid.FromString(params.ShipmentID.String()))

	var shipment *models.Shipment
	var err error
	// If TSP user, verify TSP has shipment
	// If office user, no verification necessary
	// If myApp user, user is forbidden
	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, shipment, err = models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else if session.IsOfficeUser() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else {
		return accessorialop.NewCreateShipmentAccessorialForbidden()
	}

	accessorialID := uuid.Must(uuid.FromString(params.Payload.Accessorial.ID.String()))
	shipmentAccessorial, verrs, err := shipment.CreateShipmentAccessorial(h.DB(),
		accessorialID,
		params.Payload.Quantity1,
		params.Payload.Quantity2,
		string(params.Payload.Location),
		handlers.FmtString(params.Payload.Notes),
	)
	if verrs.HasAny() || err != nil {
		h.Logger().Error("Error fetching accessorials for shipment", zap.Error(err))
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	payload := payloadForShipmentAccessorialModel(shipmentAccessorial)
	return accessorialop.NewCreateShipmentAccessorialCreated().WithPayload(payload)
}

// UpdateShipmentAccessorialHandler updates a particular shipment accessorial
type UpdateShipmentAccessorialHandler struct {
	handlers.HandlerContext
}

// Handle updates a specified shipment accessorial
func (h UpdateShipmentAccessorialHandler) Handle(params accessorialop.UpdateShipmentAccessorialParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	shipmentID := uuid.Must(uuid.FromString(params.UpdateShipmentAccessorial.ShipmentID.String()))
	shipmentAccessorialID := uuid.Must(uuid.FromString(params.ShipmentAccessorialID.String()))

	// authorization
	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, _, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else if session.IsOfficeUser() {
		_, err := models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else {
		return accessorialop.NewUpdateShipmentAccessorialForbidden()
	}

	// Fetch shipment accessorial
	shipmentAccessorial, err := models.FetchShipmentAccessorialByID(h.DB(), &shipmentAccessorialID)
	if err != nil {
		h.Logger().Error("Error fetching shipment accessorial for shipment", zap.Error(err))
		return accessorialop.NewUpdateShipmentAccessorialInternalServerError()
	}

	accessorialID := uuid.Must(uuid.FromString(params.UpdateShipmentAccessorial.Accessorial.ID.String()))

	// update
	shipmentAccessorial.AccessorialID = accessorialID
	shipmentAccessorial.Quantity1 = unit.BaseQuantity(*params.UpdateShipmentAccessorial.Quantity1)
	shipmentAccessorial.Quantity2 = unit.BaseQuantity(*params.UpdateShipmentAccessorial.Quantity2)
	shipmentAccessorial.Location = models.ShipmentAccessorialLocation(params.UpdateShipmentAccessorial.Location)
	shipmentAccessorial.Notes = params.UpdateShipmentAccessorial.Notes

	verrs, err := h.DB().ValidateAndUpdate(&shipmentAccessorial)
	if verrs.HasAny() || err != nil {
		h.Logger().Error("Error updating shipment accessorial for shipment", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	err = h.DB().Load(&shipmentAccessorial)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := payloadForShipmentAccessorialModel(&shipmentAccessorial)
	return accessorialop.NewUpdateShipmentAccessorialOK().WithPayload(payload)
}

// DeleteShipmentAccessorialHandler deletes a particular shipment accessorial
type DeleteShipmentAccessorialHandler struct {
	handlers.HandlerContext
}

// Handle deletes a specified shipment accessorial
func (h DeleteShipmentAccessorialHandler) Handle(params accessorialop.DeleteShipmentAccessorialParams) middleware.Responder {

	// Fetch shipment accessorial first
	shipmentAccessorialID := uuid.Must(uuid.FromString(params.ShipmentAccessorialID.String()))
	shipmentAccessorial, err := models.FetchShipmentAccessorialByID(h.DB(), &shipmentAccessorialID)
	if err != nil {
		h.Logger().Error("Error fetching shipment accessorial for shipment", zap.Error(err))
		return accessorialop.NewDeleteShipmentAccessorialInternalServerError()
	}

	// authorization
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	shipmentID := uuid.Must(uuid.FromString(shipmentAccessorial.ShipmentID.String()))
	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, _, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else if session.IsOfficeUser() {
		_, err := models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else {
		return accessorialop.NewDeleteShipmentAccessorialForbidden()
	}

	// Delete the shipment accessorial
	err = h.DB().Destroy(&shipmentAccessorial)
	if err != nil {
		h.Logger().Error("Error deleting shipment accessorial for shipment", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	return accessorialop.NewDeleteShipmentAccessorialOK()
}
