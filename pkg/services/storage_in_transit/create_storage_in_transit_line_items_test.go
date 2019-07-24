package storageintransit

import (
	"testing"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *StorageInTransitServiceSuite) helperSetup() {
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents = 320700
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)
}

func (suite *StorageInTransitServiceSuite) helperCreateShipment(
	originSITAddress *models.Address,
	destinationSITAddress *models.Address) (rateengine.CostByShipment, uuid.UUID) {

	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), 1, 1, []int{1}, []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	authorizedStartDate := shipment.ActualPickupDate

	if originSITAddress != nil {
		makeOriginSITAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: *originSITAddress,
		})
		sitOriginID := uuid.Must(uuid.NewV4())
		sitOrigin := models.StorageInTransit{
			ID:                  sitOriginID,
			ShipmentID:          shipment.ID,
			Shipment:            shipment,
			Location:            models.StorageInTransitLocationORIGIN,
			Status:              models.StorageInTransitStatusRELEASED,
			EstimatedStartDate:  *authorizedStartDate,
			AuthorizedStartDate: authorizedStartDate,
			ActualStartDate:     authorizedStartDate,
			WarehouseID:         "450383",
			WarehouseName:       "Extra Space Storage",
			WarehouseAddress:    makeOriginSITAddress,
		}
		testdatagen.MakeStorageInTransit(suite.DB(), testdatagen.Assertions{
			StorageInTransit: sitOrigin,
		})
	}

	if destinationSITAddress != nil {
		makeDestinationSITAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: *destinationSITAddress,
		})

		sitDestinationID := uuid.Must(uuid.NewV4())
		sitDestination := models.StorageInTransit{
			ID:                  sitDestinationID,
			ShipmentID:          shipment.ID,
			Shipment:            shipment,
			Location:            models.StorageInTransitLocationDESTINATION,
			Status:              models.StorageInTransitStatusDELIVERED,
			EstimatedStartDate:  *authorizedStartDate,
			AuthorizedStartDate: authorizedStartDate,
			ActualStartDate:     authorizedStartDate,
			WarehouseID:         "450384",
			WarehouseName:       "Iron Guard Storage",
			WarehouseAddress:    makeDestinationSITAddress,
		}
		testdatagen.MakeStorageInTransit(suite.DB(), testdatagen.Assertions{
			StorageInTransit: sitDestination,
		})
	}

	// Refetching shipments from database to get all needed eagerly fetched relationships.
	dbShipment, err := models.FetchShipmentByTSP(suite.DB(), tspUser.TransportationServiceProviderID, shipment.ID)
	suite.FatalNoError(err)

	logger, err := zap.NewDevelopment()
	suite.FatalNoError(err)

	engine := rateengine.NewRateEngine(suite.DB(), logger)
	shipmentCost, err := engine.HandleRunOnShipment(*dbShipment, dbShipment.ShippingDistance)
	suite.FatalNoError(err)

	return shipmentCost, shipment.ID
}

func (suite *StorageInTransitServiceSuite) TestCreateStorageInTransitLineItems() {

	suite.helperSetup()

	suite.T().Run("Create Storage In Transit Has 1044 distance miles", func(t *testing.T) {

		// Because of how the planner is setup the distance from storage warehouse will always be 1044 mi
		// to test fully the real distances will be tested in e2e using Cypress.io
		// 9.2 mi from Origin: Saf Keep Storage, 4996 Melrose Ave, Los Angeles, CA 90029
		sitOriginAddress := models.Address{
			StreetAddress1: "4996 Melrose Ave",
			StreetAddress2: swag.String(""),
			StreetAddress3: swag.String(""),
			City:           "Los Angeles",
			State:          "CA",
			PostalCode:     "90029",
			Country:        swag.String("US"),
		}

		shipmentCost, shipmentID := suite.helperCreateShipment(&sitOriginAddress, nil)
		suite.NotEqual(shipmentCost.Shipment.ID, uuid.Nil, "shipmentCost.Shipment.ID not uuid.Nil")

		storageInTransits, err := models.FetchStorageInTransitsOnShipment(suite.DB(), shipmentCost.Shipment.ID)
		suite.FatalNoError(err)
		suite.Len(storageInTransits, 1)

		suite.Equal(sitOriginAddress.StreetAddress1, storageInTransits[0].WarehouseAddress.StreetAddress1, "Origin SIT Address is what we expect")
		suite.Equal(sitOriginAddress.PostalCode, storageInTransits[0].WarehouseAddress.PostalCode, "Origin SIT Zip is what we expect")

		suite.Equal(shipmentID, shipmentCost.Shipment.ID,
			"shipmentID and shipmentCost.Shipment.ID IDs are the same")

		suite.Equal(shipmentCost.Shipment.ID, storageInTransits[0].ShipmentID,
			"shipmentCost.Shipment.ID and storageInTransits[0].ShipmentID IDs are the same")

		planner := route.NewTestingPlanner(shipmentCost.Shipment.ShippingDistance.DistanceMiles)

		// Create Storage in Transit (SIT) line items for Shipment
		createStorageInTransitLineItems := CreateStorageInTransitLineItems{
			DB:      suite.DB(),
			Planner: planner,
		}
		storageInTransitLineItems, err := createStorageInTransitLineItems.CreateStorageInTransitLineItems(shipmentCost)
		suite.FatalNoError(err)

		for _, sit := range storageInTransits {
			suite.Equal(shipmentCost.Shipment.ID, sit.ShipmentID, "shipmentCost.Shipment.ID, sit.ShipmentID are the same")
			suite.Equal(sit.ShipmentID, shipmentID, "sit.ShipmentID, shipmentID are the same")
		}

		item210C := suite.findLineItem(storageInTransitLineItems, "210C")
		if item210C != nil {
			suite.validateLineItemFields(*item210C, unit.BaseQuantityFromInt(1044), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationORIGIN)
		}

	})
}

func (suite *StorageInTransitServiceSuite) findLineItem(lineItems []models.ShipmentLineItem, itemCode string) *models.ShipmentLineItem {
	for _, lineItem := range lineItems {
		if itemCode == lineItem.Tariff400ngItem.Code {
			return &lineItem
		}
	}

	suite.T().Errorf("Could not find shipment line item for %s", itemCode)
	return nil
}

func (suite *StorageInTransitServiceSuite) validateLineItemFields(lineItem models.ShipmentLineItem, quantity1 unit.BaseQuantity, quantity2 unit.BaseQuantity, location models.ShipmentLineItemLocation) {
	suite.Equal(quantity1, lineItem.Quantity1)
	suite.Equal(quantity2, lineItem.Quantity2)
	suite.Equal(location, lineItem.Location)
	suite.Equal(models.ShipmentLineItemStatusAPPROVED, lineItem.Status)
}