package publicapi

import (
	"database/sql"
	"github.com/transcom/mymove/pkg/server"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForShipmentLineItemModels(s []models.ShipmentLineItem) apimessages.ShipmentLineItems {
	payloads := make(apimessages.ShipmentLineItems, len(s))

	for i, acc := range s {
		payloads[i] = payloadForShipmentLineItemModel(&acc)
	}

	return payloads
}

func payloadForShipmentLineItemModel(s *models.ShipmentLineItem) *apimessages.ShipmentLineItem {
	if s == nil {
		return nil
	}

	var amt *int64
	if s.AmountCents != nil {
		int := s.AmountCents.Int64()
		amt = &int
	}

	return &apimessages.ShipmentLineItem{
		ID:                *handlers.FmtUUID(s.ID),
		ShipmentID:        *handlers.FmtUUID(s.ShipmentID),
		Tariff400ngItem:   payloadForTariff400ngItemModel(&s.Tariff400ngItem),
		Tariff400ngItemID: handlers.FmtUUID(s.Tariff400ngItemID),
		Location:          apimessages.ShipmentLineItemLocation(s.Location),
		Notes:             s.Notes,
		Quantity1:         handlers.FmtInt64(int64(s.Quantity1)),
		Quantity2:         handlers.FmtInt64(int64(s.Quantity2)),
		Status:            apimessages.ShipmentLineItemStatus(s.Status),
		AmountCents:       amt,
		SubmittedDate:     *handlers.FmtDateTime(s.SubmittedDate),
		ApprovedDate:      *handlers.FmtDateTime(s.ApprovedDate),
	}
}

// GetShipmentLineItemsHandler returns a particular shipment line item
type GetShipmentLineItemsHandler struct {
	handlers.HandlerContext
}

// Handle returns a specified shipment line item
func (h GetShipmentLineItemsHandler) Handle(params accessorialop.GetShipmentLineItemsParams) middleware.Responder {

	session := server.SessionFromRequestContext(params.HTTPRequest)

	shipmentID := uuid.Must(uuid.FromString(params.ShipmentID.String()))

	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, _, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else if !session.IsOfficeUser() {
		return accessorialop.NewGetShipmentLineItemsForbidden()
	}

	shipmentLineItems, err := models.FetchLineItemsByShipmentID(h.DB(), &shipmentID)
	if err != nil {
		h.Logger().Error("Error fetching line items for shipment", zap.Error(err))
		return accessorialop.NewGetShipmentLineItemsInternalServerError()
	}
	payload := payloadForShipmentLineItemModels(shipmentLineItems)
	return accessorialop.NewGetShipmentLineItemsOK().WithPayload(payload)
}

// CreateShipmentLineItemHandler creates a shipment_line_item for a provided shipment_id
type CreateShipmentLineItemHandler struct {
	handlers.HandlerContext
}

// Handle handles the request
func (h CreateShipmentLineItemHandler) Handle(params accessorialop.CreateShipmentLineItemParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)

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
		return accessorialop.NewCreateShipmentLineItemForbidden()
	}

	tariff400ngItemID := uuid.Must(uuid.FromString(params.Payload.Tariff400ngItemID.String()))
	tariff400ngItem, err := models.FetchTariff400ngItem(h.DB(), tariff400ngItemID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	if !tariff400ngItem.RequiresPreApproval {
		return accessorialop.NewCreateShipmentLineItemForbidden()
	}

	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(h.DB(),
		tariff400ngItemID,
		params.Payload.Quantity1,
		params.Payload.Quantity2,
		string(params.Payload.Location),
		handlers.FmtString(params.Payload.Notes),
	)
	if verrs.HasAny() || err != nil {
		h.Logger().Error("Error fetching shipment line items for shipment", zap.Error(err))
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	payload := payloadForShipmentLineItemModel(shipmentLineItem)
	return accessorialop.NewCreateShipmentLineItemCreated().WithPayload(payload)
}

// UpdateShipmentLineItemHandler updates a particular shipment line item
type UpdateShipmentLineItemHandler struct {
	handlers.HandlerContext
}

// Handle updates a specified shipment line item
func (h UpdateShipmentLineItemHandler) Handle(params accessorialop.UpdateShipmentLineItemParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)
	shipmentLineItemID := uuid.Must(uuid.FromString(params.ShipmentLineItemID.String()))

	// Fetch shipment line item
	shipmentLineItem, err := models.FetchShipmentLineItemByID(h.DB(), &shipmentLineItemID)
	if err != nil {
		h.Logger().Error("Error fetching shipment line item for shipment", zap.Error(err))
		return accessorialop.NewUpdateShipmentLineItemInternalServerError()
	}

	// authorization
	if session.IsTspUser() {
		// Check that the TSP user can access the shipment
		_, _, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentLineItem.ShipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for TSP user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else if session.IsOfficeUser() {
		_, err := models.FetchShipment(h.DB(), session, shipmentLineItem.ShipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else {
		return accessorialop.NewUpdateShipmentLineItemForbidden()
	}

	tariff400ngItemID := uuid.Must(uuid.FromString(params.Payload.Tariff400ngItemID.String()))
	tariff400ngItem, err := models.FetchTariff400ngItem(h.DB(), tariff400ngItemID)

	if !tariff400ngItem.RequiresPreApproval {
		return accessorialop.NewUpdateShipmentLineItemForbidden()
	}

	// update
	shipmentLineItem.Tariff400ngItemID = tariff400ngItemID
	shipmentLineItem.Quantity1 = unit.BaseQuantity(*params.Payload.Quantity1)
	if params.Payload.Quantity2 != nil {
		shipmentLineItem.Quantity2 = unit.BaseQuantity(*params.Payload.Quantity2)
	}
	shipmentLineItem.Location = models.ShipmentLineItemLocation(params.Payload.Location)
	shipmentLineItem.Notes = params.Payload.Notes

	verrs, err := h.DB().ValidateAndUpdate(&shipmentLineItem)
	if verrs.HasAny() || err != nil {
		h.Logger().Error("Error updating shipment line item for shipment", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	err = h.DB().Load(&shipmentLineItem)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := payloadForShipmentLineItemModel(&shipmentLineItem)
	return accessorialop.NewUpdateShipmentLineItemOK().WithPayload(payload)
}

// DeleteShipmentLineItemHandler deletes a particular shipment line item
type DeleteShipmentLineItemHandler struct {
	handlers.HandlerContext
}

// Handle deletes a specified shipment line item
func (h DeleteShipmentLineItemHandler) Handle(params accessorialop.DeleteShipmentLineItemParams) middleware.Responder {

	// Fetch shipment line item first
	shipmentLineItemID := uuid.Must(uuid.FromString(params.ShipmentLineItemID.String()))
	shipmentLineItem, err := models.FetchShipmentLineItemByID(h.DB(), &shipmentLineItemID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			h.Logger().Error("Error shipment line item for shipment not found", zap.Error(err))
			return accessorialop.NewDeleteShipmentLineItemNotFound()
		}

		h.Logger().Error("Error fetching shipment line item for shipment", zap.Error(err))
		return accessorialop.NewDeleteShipmentLineItemInternalServerError()
	}

	if !shipmentLineItem.Tariff400ngItem.RequiresPreApproval {
		return accessorialop.NewDeleteShipmentLineItemForbidden()
	}

	// authorization
	session := server.SessionFromRequestContext(params.HTTPRequest)
	shipmentID := uuid.Must(uuid.FromString(shipmentLineItem.ShipmentID.String()))
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
		return accessorialop.NewDeleteShipmentLineItemForbidden()
	}

	// Delete the shipment line item
	err = h.DB().Destroy(&shipmentLineItem)
	if err != nil {
		h.Logger().Error("Error deleting shipment line item for shipment", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := payloadForShipmentLineItemModel(&shipmentLineItem)
	return accessorialop.NewDeleteShipmentLineItemOK().WithPayload(payload)
}

// ApproveShipmentLineItemHandler returns a particular shipment
type ApproveShipmentLineItemHandler struct {
	handlers.HandlerContext
}

// Handle returns a specified shipment
func (h ApproveShipmentLineItemHandler) Handle(params accessorialop.ApproveShipmentLineItemParams) middleware.Responder {
	session := server.SessionFromRequestContext(params.HTTPRequest)

	shipmentLineItemID := uuid.Must(uuid.FromString(params.ShipmentLineItemID.String()))

	shipmentLineItem, err := models.FetchShipmentLineItemByID(h.DB(), &shipmentLineItemID)
	if err != nil {
		h.Logger().Error("Error fetching line items for shipment", zap.Error(err))
		return accessorialop.NewApproveShipmentLineItemInternalServerError()
	}

	// Non-accessorial line items shouldn't require approval
	// Only office users can approve a shipment line item
	if shipmentLineItem.Tariff400ngItem.RequiresPreApproval && session.IsOfficeUser() {
		_, err := models.FetchShipment(h.DB(), session, shipmentLineItem.ShipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else {
		h.Logger().Error("Error does not require pre-approval for shipment")
		return accessorialop.NewApproveShipmentLineItemForbidden()
	}

	// Approve and save the shipment line item
	err = shipmentLineItem.Approve()
	if err != nil {
		h.Logger().Error("Error approving shipment line item for shipment", zap.Error(err))
		return accessorialop.NewApproveShipmentLineItemForbidden()
	}
	h.DB().ValidateAndUpdate(&shipmentLineItem)

	payload := payloadForShipmentLineItemModel(&shipmentLineItem)
	return accessorialop.NewApproveShipmentLineItemOK().WithPayload(payload)
}
