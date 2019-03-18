package publicapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/services/invoice"
	shipmentop "github.com/transcom/mymove/pkg/services/shipment"
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
		intVal := s.AmountCents.Int64()
		amt = &intVal
	}

	var rate *int64
	if s.AppliedRate != nil {
		intVal := s.AppliedRate.Int64()
		rate = &intVal
	}

	var estAmt *int64
	if s.EstimateAmountCents != nil {
		intVal := s.EstimateAmountCents.Int64()
		estAmt = &intVal
	}

	var actAmt *int64
	if s.ActualAmountCents != nil {
		intVal := s.ActualAmountCents.Int64()
		actAmt = &intVal
	}

	return &apimessages.ShipmentLineItem{
		ID:                  *handlers.FmtUUID(s.ID),
		ShipmentID:          *handlers.FmtUUID(s.ShipmentID),
		Tariff400ngItem:     payloadForTariff400ngItemModel(&s.Tariff400ngItem),
		Tariff400ngItemID:   handlers.FmtUUID(s.Tariff400ngItemID),
		Location:            apimessages.ShipmentLineItemLocation(s.Location),
		Notes:               s.Notes,
		Description:         s.Description,
		Reason:              s.Reason,
		Quantity1:           handlers.FmtInt64(int64(s.Quantity1)),
		Quantity2:           handlers.FmtInt64(int64(s.Quantity2)),
		Status:              apimessages.ShipmentLineItemStatus(s.Status),
		InvoiceID:           handlers.FmtUUIDPtr(s.InvoiceID),
		ItemDimensions:      payloadForDimensionsModel(&s.ItemDimensions),
		CrateDimensions:     payloadForDimensionsModel(&s.CrateDimensions),
		EstimateAmountCents: estAmt,
		ActualAmountCents:   actAmt,
		AmountCents:         amt,
		AppliedRate:         rate,
		SubmittedDate:       *handlers.FmtDateTime(s.SubmittedDate),
		ApprovedDate:        handlers.FmtDateTime(s.ApprovedDate),
	}
}

func payloadForDimensionsModel(a *models.ShipmentLineItemDimensions) *apimessages.Dimensions {
	if a == nil {
		return nil
	}
	if a.ID == uuid.Nil {
		return nil
	}

	return &apimessages.Dimensions{
		ID:     *handlers.FmtUUID(a.ID),
		Length: handlers.FmtInt64(int64(a.Length)),
		Width:  handlers.FmtInt64(int64(a.Width)),
		Height: handlers.FmtInt64(int64(a.Height)),
	}
}

// GetShipmentLineItemsHandler returns a particular shipment line item
type GetShipmentLineItemsHandler struct {
	handlers.HandlerContext
}

func (h GetShipmentLineItemsHandler) recalculateShipmentLineItems(shipmentLineItems models.ShipmentLineItems, shipmentID uuid.UUID, session *auth.Session) (bool, middleware.Responder) {
	update := false

	// If there is a shipment line item with an invoice do not run the recalculate function
	// the system is currently not setup to re-price a shipment with an existing invoice
	// and currently the system does not expect to have multiple invoices per shipment
	for _, item := range shipmentLineItems {
		if item.InvoiceID != nil {
			return update, nil
		}
	}

	// Need to fetch Shipment to get the Accepted Offer and the ShipmentLineItems
	// Only returning ShipmentLineItems that are approved and have no InvoiceID
	shipment, err := invoice.FetchShipmentForInvoice{DB: h.DB()}.Call(shipmentID)
	if err != nil {
		h.Logger().Error("Error fetching Shipment for re-pricing line items for shipment", zap.Error(err))
		return update, accessorialop.NewGetShipmentLineItemsInternalServerError()
	}

	// Run re-calculation process
	update, err = shipmentop.ProcessRecalculateShipment{
		DB:     h.DB(),
		Logger: h.Logger(),
	}.Call(&shipment, shipmentLineItems, h.Planner())

	if err != nil {
		h.Logger().Error("Error re-pricing line items for shipment", zap.Error(err))
		return update, accessorialop.NewGetShipmentLineItemsInternalServerError()
	}

	return update, nil
}

// Handle returns a specified shipment line item
func (h GetShipmentLineItemsHandler) Handle(params accessorialop.GetShipmentLineItemsParams) middleware.Responder {

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
		return accessorialop.NewGetShipmentLineItemsForbidden()
	}

	shipmentLineItems, err := models.FetchLineItemsByShipmentID(h.DB(), &shipmentID)
	if err != nil {
		h.Logger().Error("Error fetching line items for shipment", zap.Error(err))
		return accessorialop.NewGetShipmentLineItemsInternalServerError()
	}

	update, recalculateError := h.recalculateShipmentLineItems(shipmentLineItems, shipmentID, session)
	if recalculateError != nil {
		return recalculateError
	}
	if update {
		shipmentLineItems, err = models.FetchLineItemsByShipmentID(h.DB(), &shipmentID)
		if err != nil {
			h.Logger().Error("Error fetching line items for shipment after re-calculation",
				zap.Error(err))
			return accessorialop.NewGetShipmentLineItemsInternalServerError()
		}
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

	baseParams := models.BaseShipmentLineItemParams{
		Tariff400ngItemID:   tariff400ngItemID,
		Tariff400ngItemCode: tariff400ngItem.Code,
		Quantity1:           unit.IntToBaseQuantity(params.Payload.Quantity1),
		Quantity2:           unit.IntToBaseQuantity(params.Payload.Quantity2),
		Location:            string(params.Payload.Location),
		Notes:               handlers.FmtString(params.Payload.Notes),
	}

	var itemDimensions, crateDimensions *models.AdditionalLineItemDimensions
	if params.Payload.ItemDimensions != nil {
		itemDimensions = &models.AdditionalLineItemDimensions{
			Length: unit.ThousandthInches(*params.Payload.ItemDimensions.Length),
			Width:  unit.ThousandthInches(*params.Payload.ItemDimensions.Width),
			Height: unit.ThousandthInches(*params.Payload.ItemDimensions.Height),
		}
	}
	if params.Payload.CrateDimensions != nil {
		crateDimensions = &models.AdditionalLineItemDimensions{
			Length: unit.ThousandthInches(*params.Payload.CrateDimensions.Length),
			Width:  unit.ThousandthInches(*params.Payload.CrateDimensions.Width),
			Height: unit.ThousandthInches(*params.Payload.CrateDimensions.Height),
		}
	}

	var estAmtCents *unit.Cents
	var actAmtCents *unit.Cents
	if params.Payload.EstimateAmountCents != nil {
		centsValue := unit.Cents(*params.Payload.EstimateAmountCents)
		estAmtCents = &centsValue
	}
	if params.Payload.ActualAmountCents != nil {
		centsValue := unit.Cents(*params.Payload.ActualAmountCents)
		actAmtCents = &centsValue
	}

	additionalParams := models.AdditionalShipmentLineItemParams{
		ItemDimensions:      itemDimensions,
		CrateDimensions:     crateDimensions,
		Description:         params.Payload.Description,
		Reason:              params.Payload.Reason,
		EstimateAmountCents: estAmtCents,
		ActualAmountCents:   actAmtCents,
	}

	shipmentLineItem, verrs, err := shipment.CreateShipmentLineItem(h.DB(),
		baseParams,
		additionalParams,
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
	session := auth.SessionFromRequestContext(params.HTTPRequest)
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
	shipment := shipmentLineItem.Shipment

	if !tariff400ngItem.RequiresPreApproval {
		return accessorialop.NewUpdateShipmentLineItemForbidden()
	}

	baseParams := models.BaseShipmentLineItemParams{
		Tariff400ngItemID:   tariff400ngItemID,
		Tariff400ngItemCode: tariff400ngItem.Code,
		Quantity1:           unit.IntToBaseQuantity(params.Payload.Quantity1),
		Quantity2:           unit.IntToBaseQuantity(params.Payload.Quantity2),
		Location:            string(params.Payload.Location),
		Notes:               handlers.FmtString(params.Payload.Notes),
	}

	var itemDimensions, crateDimensions *models.AdditionalLineItemDimensions
	if params.Payload.ItemDimensions != nil {
		itemDimensions = &models.AdditionalLineItemDimensions{
			Length: unit.ThousandthInches(*params.Payload.ItemDimensions.Length),
			Width:  unit.ThousandthInches(*params.Payload.ItemDimensions.Width),
			Height: unit.ThousandthInches(*params.Payload.ItemDimensions.Height),
		}
	}
	if params.Payload.CrateDimensions != nil {
		crateDimensions = &models.AdditionalLineItemDimensions{
			Length: unit.ThousandthInches(*params.Payload.CrateDimensions.Length),
			Width:  unit.ThousandthInches(*params.Payload.CrateDimensions.Width),
			Height: unit.ThousandthInches(*params.Payload.CrateDimensions.Height),
		}
	}

	var estAmtCents *unit.Cents
	var actAmtCents *unit.Cents
	if params.Payload.EstimateAmountCents != nil {
		centsValue := unit.Cents(*params.Payload.EstimateAmountCents)
		estAmtCents = &centsValue
	}
	if params.Payload.ActualAmountCents != nil {
		centsValue := unit.Cents(*params.Payload.ActualAmountCents)
		actAmtCents = &centsValue
	}
	additionalParams := models.AdditionalShipmentLineItemParams{
		ItemDimensions:      itemDimensions,
		CrateDimensions:     crateDimensions,
		Description:         params.Payload.Description,
		Reason:              params.Payload.Reason,
		EstimateAmountCents: estAmtCents,
		ActualAmountCents:   actAmtCents,
	}

	verrs, err := shipment.UpdateShipmentLineItem(h.DB(),
		baseParams,
		additionalParams,
		&shipmentLineItem,
	)
	if verrs.HasAny() || err != nil {
		h.Logger().Error("Error fetching shipment line items for shipment", zap.Error(err))
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
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
	session := auth.SessionFromRequestContext(params.HTTPRequest)
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
	var shipment *models.Shipment
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentLineItemID := uuid.Must(uuid.FromString(params.ShipmentLineItemID.String()))

	shipmentLineItem, err := models.FetchShipmentLineItemByID(h.DB(), &shipmentLineItemID)
	if err != nil {
		h.Logger().Error("Error fetching line items for shipment", zap.Error(err))
		return accessorialop.NewApproveShipmentLineItemInternalServerError()
	}

	// Non-accessorial line items shouldn't require approval
	// Only office users can approve a shipment line item
	if shipmentLineItem.Tariff400ngItem.RequiresPreApproval && session.IsOfficeUser() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentLineItem.ShipmentID)
		if err != nil {
			h.Logger().Error("Error fetching shipment for office user", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else {
		h.Logger().Error("Error does not require pre-approval for shipment")
		return accessorialop.NewApproveShipmentLineItemForbidden()
	}

	// If shipment is delivered, price single shipment line item
	if shipmentLineItem.Shipment.Status == models.ShipmentStatusDELIVERED {
		shipmentLineItem.Shipment = *shipment
		engine := rateengine.NewRateEngine(h.DB(), h.Logger())
		err = engine.PricePreapprovalRequest(&shipmentLineItem)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
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
