package storageintransit

import (
	"fmt"
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

func (suite *StorageInTransitServiceSuite) helperCreateShipment(originSITAddress models.Address, destinationSITAddress models.Address) rateengine.CostByShipment {

	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), 1, 1, []int{1}, []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	tspUser := tspUsers[0]
	shipment := shipments[0]

	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	assertions.FuelEIADieselPrice.EIAPricePerGallonMillicents = 320700
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	pickupAddress := models.Address{
		StreetAddress1: "9611 Highridge Dr",
		StreetAddress2: swag.String("P.O. Box 12345"),
		StreetAddress3: swag.String("c/o Some Person"),
		City:           "Beverly Hills",
		State:          "CA",
		PostalCode:     "90210",
		Country:        swag.String("US"),
	}
	pickupAddress = testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: pickupAddress,
	})

	destAddress := models.Address{
		StreetAddress1: "2157 Willhaven Dr ",
		StreetAddress2: swag.String(""),
		StreetAddress3: swag.String(""),
		City:           "Augusta",
		State:          "GA",
		PostalCode:     "30909",
		Country:        swag.String("US"),
	}
	destAddress = testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: destAddress,
	})

	originSITAddress = testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: originSITAddress,
	})

	destinationSITAddress = testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: destinationSITAddress,
	})

	authorizedStartDate := shipment.ActualPickupDate
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
		WarehouseAddress:    originSITAddress,
	}
	testdatagen.MakeStorageInTransit(suite.DB(), testdatagen.Assertions{
		StorageInTransit: sitOrigin,
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
		WarehouseAddress:    destinationSITAddress,
	}
	testdatagen.MakeStorageInTransit(suite.DB(), testdatagen.Assertions{
		StorageInTransit: sitDestination,
	})

	// Refetching shipments from database to get all needed eagerly fetched relationships.
	dbShipment, err := models.FetchShipmentByTSP(suite.DB(), tspUser.TransportationServiceProviderID, shipment.ID)
	suite.FatalNoError(err)

	logger, err := zap.NewDevelopment()
	suite.FatalNoError(err)
	logger.Debug("^^^^^^^^^^^^^^^^TESTING NON TEST SUITE DEBUGGER")
	//suite.logger.Debug("$$$$$$$$$$$$$$$TEST TEST SUITE DEBUGGER")

	//engine := rateengine.NewRateEngine(suite.DB(), suite.logger)

	engine := rateengine.NewRateEngine(suite.DB(), logger)
	//storageInTransitCompare(suite, *actualStorageInTransit, sit)
	//engine := rateengine.NewRateEngine(suite.DB(), suite.logger)
	shipmentCost, err := engine.HandleRunOnShipment(*dbShipment, dbShipment.ShippingDistance)
	suite.FatalNoError(err)

	return shipmentCost
}

func (suite *StorageInTransitServiceSuite) TestCreateStorageInTransitLineItems() {

	suite.T().Run("Create Storage In Transit Less Than 30 mi", func(t *testing.T) {
		sitOriginAddress := models.Address{
			StreetAddress1: "1860 Vine St",
			StreetAddress2: swag.String(""),
			StreetAddress3: swag.String(""),
			City:           "Los Angeles",
			State:          "CA",
			PostalCode:     "90028",
			Country:        swag.String("US"),
		}

		sitDestinationAddress := models.Address{
			StreetAddress1: "1045 Bertram Rd",
			StreetAddress2: swag.String(""),
			StreetAddress3: swag.String(""),
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30909",
			Country:        swag.String("US"),
		}

		shipmentCost := suite.helperCreateShipment(sitOriginAddress, sitDestinationAddress)

		//planner := route.NewTestingPlanner(1100)
		//planner := route.NewTestingPlanner(shipmentCost.Cost.Mileage)
		planner := route.NewTestingPlanner(shipmentCost.Shipment.ShippingDistance.DistanceMiles)
		// Create Storage in Transit (SIT) line items for Shipment

		createStorageInTransitLineItems := CreateStorageInTransitLineItems{
			DB:      suite.DB(),
			Planner: planner,
		}
		storageInTransitLineItems, err := createStorageInTransitLineItems.CreateStorageInTransitLineItems(shipmentCost)
		suite.FatalNoError(err)

		// There are 6 Base Shipment line items:
		// origin fee, destination fee, linehaul, pack, unpack, fuel surcharge
		//suite.Len(lineItems, 6)

		/*

			itemLHS := suite.findLineItem(storageInTransitLineItems, "210A")
			if itemLHS != nil {
				//suite.validateLineItemFields(*itemLHS, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(1044), models.ShipmentLineItemLocationORIGIN, unit.Cents(260858), unit.Millicents(0))
			}

			item135A := suite.findLineItem(storageInTransitLineItems, "210B")
			if item135A != nil {
				//suite.validateLineItemFields(*item135A, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationORIGIN, unit.Cents(10230), unit.Millicents(511000))
			}
		*/

		item135B := suite.findLineItem(storageInTransitLineItems, "210C")
		if item135B != nil {
			//suite.validateLineItemFields(*item135B, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationDESTINATION, unit.Cents(11524), unit.Millicents(576000))
		}

	})

	suite.T().Run("Create Storage In Transit Less Than 50 mi", func(t *testing.T) {

	})

	suite.T().Run("Create Storage In Transit At Least 50 mi", func(t *testing.T) {
		sitOriginAddress := models.Address{
			StreetAddress1: "1860 Vine St",
			StreetAddress2: swag.String(""),
			StreetAddress3: swag.String(""),
			City:           "Los Angeles",
			State:          "CA",
			PostalCode:     "90028",
			Country:        swag.String("US"),
		}

		sitDestinationAddress := models.Address{
			StreetAddress1: "1045 Bertram Rd",
			StreetAddress2: swag.String(""),
			StreetAddress3: swag.String(""),
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30909",
			Country:        swag.String("US"),
		}

		shipmentCost := suite.helperCreateShipment(sitOriginAddress, sitDestinationAddress)

		//planner := route.NewTestingPlanner(1100)
		//planner := route.NewTestingPlanner(shipmentCost.Cost.Mileage)
		planner := route.NewTestingPlanner(shipmentCost.Shipment.ShippingDistance.DistanceMiles)
		// Create Storage in Transit (SIT) line items for Shipment

		createStorageInTransitLineItems := CreateStorageInTransitLineItems{
			DB:      suite.DB(),
			Planner: planner,
		}
		storageInTransitLineItems, err := createStorageInTransitLineItems.CreateStorageInTransitLineItems(shipmentCost)
		suite.FatalNoError(err)

		// There are 6 Base Shipment line items:
		// origin fee, destination fee, linehaul, pack, unpack, fuel surcharge
		//suite.Len(lineItems, 6)

		/*

			item210A := suite.findLineItem(storageInTransitLineItems, "210A")
			if item210A != nil {
				//suite.validateLineItemFields(*itemLHS, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(1044), models.ShipmentLineItemLocationORIGIN, unit.Cents(260858), unit.Millicents(0))
			}

			item210B := suite.findLineItem(storageInTransitLineItems, "210B")
			if item210B != nil {
				//suite.validateLineItemFields(*item135A, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationORIGIN, unit.Cents(10230), unit.Millicents(511000))
			}
		*/

		item210C := suite.findLineItem(storageInTransitLineItems, "210C")
		if item210C != nil {
			//suite.validateLineItemFields(*item135B, unit.BaseQuantityFromInt(2000), unit.BaseQuantityFromInt(0), models.ShipmentLineItemLocationDESTINATION, unit.Cents(11524), unit.Millicents(576000))
		}
	})
}

func (suite *StorageInTransitServiceSuite) findLineItem(lineItems []models.ShipmentLineItem, itemCode string) *models.ShipmentLineItem {
	for _, lineItem := range lineItems {
		if itemCode == lineItem.Tariff400ngItem.Code {
			fmt.Println("DEBUG line item",
				zap.Any("code", itemCode),
				zap.Any("quantity1", lineItem.Quantity1),
				zap.Any("quantity2", lineItem.Quantity2),
				zap.Any("location", lineItem.Location),
				zap.Any("status", lineItem.Status),
			//zap.Any("amountCents", *lineItem.AmountCents),
			//zap.Any("appliedRate", *lineItem.AppliedRate),
			)
			return &lineItem
		}
	}

	suite.T().Errorf("Could not find shipment line item for %s", itemCode)
	return nil
}

func (suite *StorageInTransitServiceSuite) validateLineItemFields(lineItem models.ShipmentLineItem, quantity1 unit.BaseQuantity, quantity2 unit.BaseQuantity, location models.ShipmentLineItemLocation, amountCents unit.Cents, appliedRate unit.Millicents) {
	suite.Equal(quantity1, lineItem.Quantity1)
	suite.Equal(quantity2, lineItem.Quantity2)
	suite.Equal(location, lineItem.Location)
	suite.Equal(amountCents, *lineItem.AmountCents)
	suite.Equal(appliedRate, *lineItem.AppliedRate)

	suite.Equal(models.ShipmentLineItemStatusSUBMITTED, lineItem.Status)
}