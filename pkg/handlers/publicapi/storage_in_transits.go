package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	sitop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/storage_in_transits"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
	"time"
)

func payloadForStorageInTransitModel(s *models.StorageInTransit) *apimessages.StorageInTransit {
	if s == nil {
		return nil
	}

	return &apimessages.StorageInTransit{
		EstimatedStartDate: handlers.FmtDate(s.EstimatedStartDate),
		Notes:              handlers.FmtStringPtr(s.Notes),
		WarehouseAddress:   payloadForAddressModel(&s.WarehouseAddress),
		WarehouseEmail:     handlers.FmtStringPtr(s.WarehouseEmail),
		WarehouseID:        handlers.FmtString(s.WarehouseID),
		WarehouseName:      handlers.FmtString(s.WarehouseName),
		WarehousePhone:     handlers.FmtStringPtr(s.WarehousePhone),
	}
}

// IndexStorageInTransitHandler returns a list of Storage In Transit entries
type IndexStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of Storage In Transit entries
func (h IndexStorageInTransitHandler) Handle(params sitop.CreateStorageInTransitParams) middleware.Responder {
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

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
	//session := auth.SessionFromRequestContext(params.HTTPRequest)
	payload := params.StorageInTransit

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	incomingLocation := *payload.Location
	var savedLocation models.StorageInTransitLocation

	if incomingLocation == "ORIGIN" {
		savedLocation = models.StorageInTransitLocationORIGIN
	}

	if incomingLocation == "DESTINATION" {
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
		_, err := h.DB().ValidateAndCreate(&warehouseAddress)

		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
	}

	newStorageInTransit := models.StorageInTransit{
		ShipmentID:         shipmentID,
		Status:             models.StorageInTransitStatusREQUESTED,
		Location:           savedLocation,
		EstimatedStartDate: estimatedStartDate,
		Notes:              payload.Notes,
		WarehouseID:        *payload.WarehouseID,
		WarehouseName:      warehouseName,
		WarehouseAddressID: warehouseAddress.ID,
		WarehouseAddress:   warehouseAddress,
		WarehouseEmail:     payload.WarehouseEmail,
	}

	sitPayload := payloadForStorageInTransitModel(&newStorageInTransit)

	return sitop.NewCreateStorageInTransitCreated().WithPayload(sitPayload)

}

// GetStorageInTransitHandler gets a single Storage In Transit based on its own ID
type GetStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
func (h GetStorageInTransitHandler) Handle(params sitop.GetStorageInTransitParams) middleware.Responder {
	storageInTransitID, _ := uuid.FromString(params.StorageInTransitID.String())
	storageInTransit, err := models.FetchStorageInTransitsByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(&storageInTransit)
	return sitop.NewGetStorageInTransitOK().WithPayload(storageInTransitPayload)

}
