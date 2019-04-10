package publicapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
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
		AuthorizedStartDate: handlers.FmtDatePtr(s.AuthorizedStartDate),
		ActualStartDate:     handlers.FmtDatePtr(s.ActualStartDate),
		Notes:               handlers.FmtStringPtr(s.Notes),
		WarehouseAddress:    payloadForAddressModel(&s.WarehouseAddress),
		WarehouseEmail:      handlers.FmtStringPtr(s.WarehouseEmail),
		WarehouseID:         handlers.FmtString(s.WarehouseID),
		WarehouseName:       handlers.FmtString(s.WarehouseName),
		WarehousePhone:      handlers.FmtStringPtr(s.WarehousePhone),
		Location:            &location,
		Status:              *handlers.FmtString(status),
	}
}

func authorizeStorageInTransitRequest(db *pop.Connection, session *auth.Session, shipmentID uuid.UUID, allowOffice bool) (isUserAuthorized bool, err error) {
	if session.IsTspUser() {
		_, _, err := models.FetchShipmentForVerifiedTSPUser(db, session.TspUserID, shipmentID)

		if err != nil {
			return false, err
		}
		return true, nil
	} else if session.IsOfficeUser() {
		if allowOffice {
			return true, nil
		}
	} else {
		return false, models.ErrFetchForbidden
	}
	return false, models.ErrFetchForbidden
}

func processStorageInTransitInput(h handlers.HandlerContext, shipmentID uuid.UUID, payload apimessages.StorageInTransit) (models.StorageInTransit, error) {
	incomingLocation := *payload.Location
	var savedLocation models.StorageInTransitLocation

	if incomingLocation == "ORIGIN" {
		savedLocation = models.StorageInTransitLocationORIGIN
	} else {
		savedLocation = models.StorageInTransitLocationDESTINATION
	}

	var estimatedStartDate time.Time
	if payload.EstimatedStartDate != nil {
		estimatedStartDate = time.Time(*payload.EstimatedStartDate)
	}

	var warehouseName string
	if payload.WarehouseName != nil {
		warehouseName = *payload.WarehouseName
	}

	var warehouseAddress models.Address
	if payload.WarehouseAddress != nil {
		warehouseAddress = models.Address{
			StreetAddress1: *payload.WarehouseAddress.StreetAddress1,
			StreetAddress2: payload.WarehouseAddress.StreetAddress2,
			StreetAddress3: payload.WarehouseAddress.StreetAddress3,
			City:           *payload.WarehouseAddress.City,
			State:          *payload.WarehouseAddress.State,
			PostalCode:     *payload.WarehouseAddress.PostalCode,
			Country:        payload.WarehouseAddress.Country,
		}
	}

	newStorageInTransit := models.StorageInTransit{
		ShipmentID:         shipmentID,
		Location:           savedLocation,
		EstimatedStartDate: estimatedStartDate,
		Notes:              payload.Notes,
		WarehouseID:        *payload.WarehouseID,
		WarehouseName:      warehouseName,
		WarehouseAddressID: warehouseAddress.ID,
		WarehouseAddress:   warehouseAddress,
		WarehouseEmail:     payload.WarehouseEmail,
		WarehousePhone:     payload.WarehousePhone,
	}

	return newStorageInTransit, nil

}

func patchStorageInTransitWithPayload(storageInTransit *models.StorageInTransit, payload *apimessages.StorageInTransit) {
	if *payload.Location == "ORIGIN" {
		storageInTransit.Location = models.StorageInTransitLocationORIGIN
	} else {
		storageInTransit.Location = models.StorageInTransitLocationDESTINATION
	}

	if payload.EstimatedStartDate != nil {
		storageInTransit.EstimatedStartDate = *(*time.Time)(payload.EstimatedStartDate)
	}

	storageInTransit.Notes = handlers.FmtStringPtrNonEmpty(payload.Notes)

	if payload.WarehouseID != nil {
		storageInTransit.WarehouseID = *payload.WarehouseID
	}

	if payload.WarehouseName != nil {
		storageInTransit.WarehouseName = *payload.WarehouseName
	}

	if payload.WarehouseAddress != nil {
		updateAddressWithPayload(&storageInTransit.WarehouseAddress, payload.WarehouseAddress)
	}

	storageInTransit.WarehousePhone = handlers.FmtStringPtrNonEmpty(payload.WarehousePhone)
	storageInTransit.WarehouseEmail = handlers.FmtStringPtrNonEmpty(payload.WarehouseEmail)
}

// IndexStorageInTransitHandler returns a list of Storage In Transit entries
type IndexStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of Storage In Transit entries
func (h IndexStorageInTransitHandler) Handle(params sitop.IndexStorageInTransitsParams) middleware.Responder {
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, true)

	if isUserAuthorized == false {
		h.Logger().Error("Unauthorized User", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransits, err := models.FetchStorageInTransitsOnShipment(h.DB(), shipmentID)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
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
}

// Handle handles the handling
func (h CreateStorageInTransitHandler) Handle(params sitop.CreateStorageInTransitParams) middleware.Responder {
	payload := params.StorageInTransit
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, false)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	newStorageInTransit, err := processStorageInTransitInput(h, shipmentID, *payload)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	newStorageInTransit.Status = models.StorageInTransitStatusREQUESTED

	verrs, err := models.SaveStorageInTransitAndAddress(h.DB(), &newStorageInTransit)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(&newStorageInTransit)

	return sitop.NewCreateStorageInTransitCreated().WithPayload(storageInTransitPayload)

}

// PatchStorageInTransitHandler updates an existing Storage In Transit entry
type PatchStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
func (h PatchStorageInTransitHandler) Handle(params sitop.PatchStorageInTransitParams) middleware.Responder {
	payload := params.StorageInTransit
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	storageInTransitID, _ := uuid.FromString(params.StorageInTransitID.String())

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, true)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)
	if err != nil {
		h.Logger().Error("Could not find existing SIT record", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if storageInTransit.ShipmentID != shipmentID {
		h.Logger().Error("Shipment ID clash between endpoint URL and SIT record")
		return sitop.NewPatchStorageInTransitForbidden()
	}

	patchStorageInTransitWithPayload(storageInTransit, payload)

	verrs, err := models.SaveStorageInTransitAndAddress(h.DB(), storageInTransit)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(storageInTransit)

	return sitop.NewPatchStorageInTransitOK().WithPayload(storageInTransitPayload)
}

// GetStorageInTransitHandler gets a single Storage In Transit based on its own ID
type GetStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
func (h GetStorageInTransitHandler) Handle(params sitop.GetStorageInTransitParams) middleware.Responder {
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	shipmentID, err := uuid.FromString(params.ShipmentID.String())

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, true)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewGetStorageInTransitOK().WithPayload(storageInTransitPayload)

}

// DeleteStorageInTransitHandler deletes a Storage in Transit based on the provided ID
type DeleteStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
func (h DeleteStorageInTransitHandler) Handle(params sitop.DeleteStorageInTransitParams) middleware.Responder {
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	shipmentID, err := uuid.FromString(params.ShipmentID.String())

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, false)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	err = models.DeleteStorageInTransit(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	return sitop.NewDeleteStorageInTransitOK()

}
