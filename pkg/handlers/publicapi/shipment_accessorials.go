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
		Quantity1:     handlers.FmtInt64(int64(s.Quantity1.ToUnitInt())),
		Quantity2:     handlers.FmtInt64(int64(s.Quantity2.ToUnitInt())),
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
