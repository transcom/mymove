package storageintransit

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type createStorageInTransit struct {
	db *pop.Connection
}

// CreateStorageInTransit creates a new Storage In Transit for a shipment and returns the newly created object.
func (c *createStorageInTransit) CreateStorageInTransit(storageInTransitPayload apimessages.StorageInTransit, shipmentID uuid.UUID, session *auth.Session) (*models.StorageInTransit, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()

	isUserAuthorized, err := authorizeStorageInTransitHTTPRequest(c.db, session, shipmentID, false)

	if err != nil {
		return nil, returnVerrs, err
	}

	if !isUserAuthorized {
		return nil, returnVerrs, err
	}

	storageInTransit, err := processStorageInTransitInput(shipmentID, storageInTransitPayload)

	if err != nil {
		return nil, returnVerrs, err
	}

	verrs, err := models.SaveStorageInTransitAndAddress(c.db, &storageInTransit)

	if err != nil || verrs.HasAny() {
		returnVerrs.Append(verrs)
		return nil, returnVerrs, err
	}
	return &storageInTransit, returnVerrs, nil
}

func processStorageInTransitInput(shipmentID uuid.UUID, payload apimessages.StorageInTransit) (models.StorageInTransit, error) {
	incomingLocation := *payload.Location
	var savedLocation models.StorageInTransitLocation

	if incomingLocation == "ORIGIN" {
		savedLocation = models.StorageInTransitLocationORIGIN
	} else {
		savedLocation = models.StorageInTransitLocationDESTINATION
	}

	status := models.StorageInTransitStatusREQUESTED

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
		Status:             status,
	}

	return newStorageInTransit, nil

}

// NewStorageInTransitCreator is the public constructor for a `StorageInTransitCreator`
// using Pop
func NewStorageInTransitCreator(db *pop.Connection) services.StorageInTransitCreator {
	return &createStorageInTransit{db}
}
