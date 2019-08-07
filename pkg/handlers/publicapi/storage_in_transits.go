package publicapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	sitop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/storage_in_transits"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForStorageInTransitModel(s *models.StorageInTransit) *apimessages.StorageInTransit {
	if s == nil {
		return nil
	}

	location := string(s.Location)
	status := string(s.Status)

	return &apimessages.StorageInTransit{
		ID:                  *handlers.FmtUUID(s.ID),
		ShipmentID:          *handlers.FmtUUID(s.ShipmentID),
		EstimatedStartDate:  handlers.FmtDate(s.EstimatedStartDate),
		Notes:               handlers.FmtStringPtr(s.Notes),
		WarehouseAddress:    payloadForAddressModel(&s.WarehouseAddress),
		WarehouseEmail:      handlers.FmtStringPtr(s.WarehouseEmail),
		WarehouseID:         handlers.FmtString(s.WarehouseID),
		WarehouseName:       handlers.FmtString(s.WarehouseName),
		WarehousePhone:      handlers.FmtStringPtr(s.WarehousePhone),
		Location:            &location,
		Status:              *handlers.FmtString(status),
		AuthorizationNotes:  handlers.FmtStringPtr(s.AuthorizationNotes),
		AuthorizedStartDate: handlers.FmtDatePtr(s.AuthorizedStartDate),
		ActualStartDate:     handlers.FmtDatePtr(s.ActualStartDate),
		OutDate:             handlers.FmtDatePtr(s.OutDate),
		SitNumber:           s.SITNumber,
	}
}

// IndexStorageInTransitHandler returns a list of Storage In Transit entries
type IndexStorageInTransitHandler struct {
	handlers.HandlerContext
	storageInTransitIndexer services.StorageInTransitsIndexer
}

// Handle handles the handling
// This is meant to return a list of storage in transits using their shipment ID.
func (h IndexStorageInTransitHandler) Handle(params sitop.IndexStorageInTransitsParams) middleware.Responder {
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	if err != nil {
		logger.Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransits, err := h.storageInTransitIndexer.IndexStorageInTransits(shipmentID, session)

	if err != nil {
		logger.Error(fmt.Sprintf("SITs Retrieval failed for shipment: %s", shipmentID), zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransitsList := make(apimessages.StorageInTransits, len(storageInTransits))

	for i, storageInTransit := range storageInTransits {
		storageInTransitsList[i] = payloadForStorageInTransitModel(&storageInTransit)
	}

	return sitop.NewIndexStorageInTransitsOK().WithPayload(storageInTransitsList)
}

// CreateStorageInTransitHandler creates a storage in transit entry and returns a payload.
type CreateStorageInTransitHandler struct {
	handlers.HandlerContext
	storageInTransitCreator services.StorageInTransitCreator
}

// Handle handles the handling
// This is meant to create a storage in transit and return it in a payload
func (h CreateStorageInTransitHandler) Handle(params sitop.CreateStorageInTransitParams) middleware.Responder {
	payload := params.StorageInTransit
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error(fmt.Sprintf("invalid shipment id: %s", params.ShipmentID.String()), zap.Error(err))
	}

	newStorageInTransit, verrs, err := h.storageInTransitCreator.CreateStorageInTransit(*payload, shipmentID, session)

	if verrs.HasAny() || err != nil {
		logger.Error(fmt.Sprintf("SIT Creation failed for shipment: %s", shipmentID), zap.Error(err), zap.Error(verrs))
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(newStorageInTransit)

	return sitop.NewCreateStorageInTransitCreated().WithPayload(storageInTransitPayload)
}

// ApproveStorageInTransitHandler approves an existing Storage In Transit
type ApproveStorageInTransitHandler struct {
	handlers.HandlerContext
	storageInTransitApprover services.StorageInTransitApprover
}

// Handle handles the handling
// This is meant to set the status for a storage in transit to approved, save the authorization notes that
// support that status, save the authorization date, and return the saved object in a payload.
func (h ApproveStorageInTransitHandler) Handle(params sitop.ApproveStorageInTransitParams) middleware.Responder {

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	payload := params.StorageInTransitApprovalPayload
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error("invalid shipment id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())

	if err != nil {
		logger.Error("invalid storage in transit id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransit, verrs, err := h.storageInTransitApprover.ApproveStorageInTransit(*payload, shipmentID, session, storageInTransitID)

	if err != nil || verrs.HasAny() {
		logger.Error(fmt.Sprintf("SIT approval failed for ID: %s on shipment: %s", storageInTransitID, shipmentID), zap.Error(err), zap.Error(verrs))
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewApproveStorageInTransitOK().WithPayload(returnPayload)

}

// DenyStorageInTransitHandler denies an existing storage in transit
type DenyStorageInTransitHandler struct {
	handlers.HandlerContext
	storageInTransitDenier services.StorageInTransitDenier
}

// Handle handles the handling
// This is meant to set the status for a storage in transit to denied, save the supporting authorization notes,
// and return the saved object in a payload.
func (h DenyStorageInTransitHandler) Handle(params sitop.DenyStorageInTransitParams) middleware.Responder {

	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	payload := params.StorageInTransitDenyPayload
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error("Invalid shipment id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())

	if err != nil {
		logger.Error("Invalid storage in transit id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransit, verrs, err := h.storageInTransitDenier.DenyStorageInTransit(*payload, shipmentID, session, storageInTransitID)

	if err != nil || verrs.HasAny() {
		logger.Error(fmt.Sprintf("SIT denial failed for ID: %s on shipment: %s", storageInTransitID, shipmentID), zap.Error(err), zap.Error(verrs))
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewDenyStorageInTransitOK().WithPayload(returnPayload)

}

// InSitStorageInTransitHandler places storage in transit into in sit
type InSitStorageInTransitHandler struct {
	handlers.HandlerContext
	storageInTransitInSITPlacer services.StorageInTransitInSITPlacer
}

// Handle handles the handling
// This is meant to set the status for a storage in transit to 'in SIT' and return the saved object in a payload.
func (h InSitStorageInTransitHandler) Handle(params sitop.InSitStorageInTransitParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error("Invalid shipment id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	if err != nil {
		logger.Error("Invalid storage in transit id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	inSitPayload := params.StorageInTransitInSitPayload

	storageInTransit, verrs, err := h.storageInTransitInSITPlacer.PlaceIntoSITStorageInTransit(*inSitPayload, shipmentID, session, storageInTransitID)

	if err != nil || verrs.HasAny() {
		logger.Error(fmt.Sprintf("Place into SIT failed for ID: %s on shipment: %s", storageInTransitID, shipmentID), zap.Error(err), zap.Error(verrs))
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewInSitStorageInTransitOK().WithPayload(returnPayload)

}

// ReleaseStorageInTransitHandler releases an existing storage in transit
type ReleaseStorageInTransitHandler struct {
	handlers.HandlerContext
	releaseStorageInTransit services.StorageInTransitReleaser
}

// Handle Handles the handling
// This is meant to set the status of a storage in transit to released, save the actual date that supports that,
// and return the saved object in a payload.
func (h ReleaseStorageInTransitHandler) Handle(params sitop.ReleaseStorageInTransitParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error("Invalid shipment id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	if err != nil {
		logger.Error("Invalid storage in transit id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	payload := params.StorageInTransitOnReleasePayload

	storageInTransit, verrs, err := h.releaseStorageInTransit.ReleaseStorageInTransit(*payload, shipmentID, session, storageInTransitID)

	if err != nil || verrs.HasAny() {
		logger.Error(fmt.Sprintf("Release SIT failed for ID: %s on shipment: %s", storageInTransitID, shipmentID), zap.Error(err), zap.Error(verrs))
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewReleaseStorageInTransitOK().WithPayload(returnPayload)

}

// PatchStorageInTransitHandler updates an existing Storage In Transit entry
type PatchStorageInTransitHandler struct {
	handlers.HandlerContext
	patchStorageInTransit services.StorageInTransitPatcher
}

// Handle handles the handling
// This is meant to edit a storage in transit object based on provided parameters and return the saved object
// in a payload
func (h PatchStorageInTransitHandler) Handle(params sitop.PatchStorageInTransitParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	payload := params.StorageInTransit
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error("Invalid shipment id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	if err != nil {
		logger.Error("Invalid storage in transit id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransit, verrs, err := h.patchStorageInTransit.PatchStorageInTransit(*payload, shipmentID, storageInTransitID, session)

	if err != nil || verrs.HasAny() {
		logger.Error(fmt.Sprintf("Patch SIT failed for ID: %s on shipment: %s", storageInTransitID, shipmentID), zap.Error(err), zap.Error(verrs))
		return handlers.ResponseForVErrors(logger, verrs, err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(storageInTransit)

	return sitop.NewPatchStorageInTransitOK().WithPayload(storageInTransitPayload)
}

// GetStorageInTransitHandler gets a single Storage In Transit based on its own ID
type GetStorageInTransitHandler struct {
	handlers.HandlerContext
	storageInTransitFetcher services.StorageInTransitByIDFetcher
}

// Handle handles the handling
// This is meant to fetch a single storage in transit using its shipment and object IDs
func (h GetStorageInTransitHandler) Handle(params sitop.GetStorageInTransitParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	if err != nil {
		logger.Error("Invalid storage in transit id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error("Invalid shipment id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransit, err := h.storageInTransitFetcher.FetchStorageInTransitByID(storageInTransitID, shipmentID, session)
	if err != nil {
		logger.Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewGetStorageInTransitOK().WithPayload(storageInTransitPayload)
}

// DeleteStorageInTransitHandler deletes a Storage in Transit based on the provided ID
type DeleteStorageInTransitHandler struct {
	handlers.HandlerContext
	deleteStorageInTransit services.StorageInTransitDeleter
}

// Handle handles the handling
// This is meant to delete a storage in transit object using its own shipment and object IDs
func (h DeleteStorageInTransitHandler) Handle(params sitop.DeleteStorageInTransitParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	if err != nil {
		logger.Error("Invalid storage in transit id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	if err != nil {
		logger.Error("Invalid shipment id", zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	storageInTransit, err := h.deleteStorageInTransit.DeleteStorageInTransit(shipmentID, storageInTransitID, session)

	if err != nil {
		logger.Error(fmt.Sprintf("Deleting SIT failed for id: %s on shipment: %s", storageInTransitID, shipmentID), zap.Error(err))
		return handlers.ResponseForError(logger, err)
	}

	payload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewDeleteStorageInTransitOK().WithPayload(payload)
}
