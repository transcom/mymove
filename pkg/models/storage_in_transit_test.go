package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestStorageInTransitValidations() {
	suite.T().Run("test valid storage in transit", func(t *testing.T) {
		validStorageInTransit := testdatagen.MakeDefaultStorageInTransit(suite.DB())
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validStorageInTransit, expErrors)
	})

	suite.T().Run("test invalid/empty storage in transit", func(t *testing.T) {
		invalidStorageInTransit := &models.StorageInTransit{}
		expErrors := map[string][]string{
			"shipment_id":          {"ShipmentID can not be blank."},
			"status":               {"Status is not in the list [REQUESTED, APPROVED, DENIED, IN_SIT, RELEASED, DELIVERED]."},
			"location":             {"Location is not in the list [ORIGIN, DESTINATION]."},
			"estimated_start_date": {"EstimatedStartDate can not be blank."},
			"warehouse_id":         {"WarehouseID can not be blank."},
			"warehouse_name":       {"WarehouseName can not be blank."},
			"warehouse_address_id": {"WarehouseAddressID can not be blank."},
		}
		suite.verifyValidationErrors(invalidStorageInTransit, expErrors)
	})

	suite.T().Run("test actual start date cannot be before the authorized start date", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultShipment(suite.DB())
		actualStartDate := time.Date(2019, time.Month(3), 25, 0, 0, 0, 0, time.UTC)
		uuid, _ := uuid.NewV4()

		sitWithNoAuthorizedStartDate := models.StorageInTransit{
			Location:            models.StorageInTransitLocationORIGIN,
			ShipmentID:          shipment.ID,
			EstimatedStartDate:  testdatagen.DateInsidePeakRateCycle,
			WarehouseID:         "000383",
			WarehouseName:       "Hercules Hauling",
			WarehouseAddressID:  uuid,
			WarehousePhone:      swag.String("(713) 868-3497"),
			WarehouseEmail:      swag.String("joe@herculeshauling.com"),
			Status:              "APPROVED",
			AuthorizedStartDate: nil,
			ActualStartDate:     &actualStartDate,
		}

		suite.verifyValidationErrors(&sitWithNoAuthorizedStartDate, map[string][]string{
			"actual_start_date": {fmt.Sprintf("cannot create this date without a no-earlier-than date")},
		})

		authorizedStartDate := time.Date(2019, time.Month(3), 26, 0, 0, 0, 0, time.UTC)
		stringAuthorizedStartDate := authorizedStartDate.Format("2006-01-02 15:04:05 -0700 UTC")
		stringActualStartDate := actualStartDate.Format("2006-01-02 15:04:05 -0700 UTC")

		sitWithInvalidActualStartDate := models.StorageInTransit{
			Location:            models.StorageInTransitLocationORIGIN,
			ShipmentID:          shipment.ID,
			EstimatedStartDate:  testdatagen.DateInsidePeakRateCycle,
			WarehouseID:         "000383",
			WarehouseName:       "Hercules Hauling",
			WarehouseAddressID:  uuid,
			WarehousePhone:      swag.String("(713) 868-3497"),
			WarehouseEmail:      swag.String("joe@herculeshauling.com"),
			Status:              "APPROVED",
			AuthorizedStartDate: &authorizedStartDate,
			ActualStartDate:     &actualStartDate,
		}
		suite.verifyValidationErrors(&sitWithInvalidActualStartDate, map[string][]string{
			"actual_start_date": {fmt.Sprintf("%s must be on or after %s", stringActualStartDate, stringAuthorizedStartDate)},
		})
	})
}

func (suite *ModelSuite) TestFetchStorageInTransitsByShipment() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	assertions := testdatagen.Assertions{
		StorageInTransit: models.StorageInTransit{
			Location:           models.StorageInTransitLocationORIGIN,
			ShipmentID:         shipment.ID,
			EstimatedStartDate: testdatagen.DateInsidePeakRateCycle,
		},
	}

	for i := 0; i < 10; i++ {
		testdatagen.MakeStorageInTransit(suite.DB(), assertions)
	}

	storageInTransits, err := models.FetchStorageInTransitsOnShipment(suite.DB(), shipment.ID)

	suite.Nil(err)
	suite.Equal(10, len(storageInTransits))

}

func (suite *ModelSuite) TestFetchStorageInTransistByID() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	assertions := testdatagen.Assertions{
		StorageInTransit: models.StorageInTransit{
			Location:           models.StorageInTransitLocationORIGIN,
			ShipmentID:         shipment.ID,
			EstimatedStartDate: testdatagen.DateInsidePeakRateCycle,
			WarehouseEmail:     swag.String("test@tester.org"),
		},
	}
	createdSIT := testdatagen.MakeStorageInTransit(suite.DB(), assertions)

	fetchedSIT, err := models.FetchStorageInTransitByID(suite.DB(), createdSIT.ID)

	suite.Nil(err)
	suite.NotEmpty(fetchedSIT)
	suite.Equal(createdSIT.ID, fetchedSIT.ID)
	suite.Equal(*createdSIT.WarehouseEmail, *createdSIT.WarehouseEmail)

	// Let's make sure we can trigger a FetchNotFound
	fakeUUID, _ := uuid.FromString("bb2de0f1-f069-4823-a4fa-bc1c89d86506")
	_, err = models.FetchStorageInTransitByID(suite.DB(), fakeUUID)
	suite.Equal(err, models.ErrFetchNotFound)

}

func (suite *ModelSuite) TestDestroyStorageInTransit() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())

	assertions := testdatagen.Assertions{
		StorageInTransit: models.StorageInTransit{
			Location:           models.StorageInTransitLocationORIGIN,
			ShipmentID:         shipment.ID,
			EstimatedStartDate: testdatagen.DateInsidePeakRateCycle,
			WarehouseEmail:     swag.String("test@tester.org"),
		},
	}
	createdSIT := testdatagen.MakeStorageInTransit(suite.DB(), assertions)

	// Let's send a zero value as the ID to ensure that fails with a ErrFetchNotFound
	err := models.DeleteStorageInTransit(suite.DB(), uuid.UUID{})
	suite.Equal(models.ErrFetchNotFound, err)

	// Make sure we can delete successfully
	err = models.DeleteStorageInTransit(suite.DB(), createdSIT.ID)
	suite.Equal(nil, err)

	// We should get ErrFetchNotFound now that the record is deleted
	_, err = models.FetchStorageInTransitByID(suite.DB(), createdSIT.ID)
	suite.Equal(models.ErrFetchNotFound, err)

}

func (suite *ModelSuite) TestSaveStorageInTransitAndAddress() {
	shipment := testdatagen.MakeDefaultShipment(suite.DB())
	address := models.Address{
		StreetAddress1: "123 Street",
		PostalCode:     "12345",
		State:          "IL",
		City:           "Chicago",
	}

	storageInTransit := models.StorageInTransit{
		Shipment:           shipment,
		ShipmentID:         shipment.ID,
		Location:           models.StorageInTransitLocationORIGIN,
		Status:             models.StorageInTransitStatusREQUESTED,
		EstimatedStartDate: testdatagen.DateInsidePeakRateCycle,
		Notes:              swag.String("This is a note"),
		WarehouseName:      "Warehouse name",
		WarehouseID:        "12345",
		WarehousePhone:     swag.String("312-111-1111"),
		WarehouseEmail:     swag.String("email@thewarehouse"),
		WarehouseAddress:   address,
	}

	verrs, err := models.SaveStorageInTransitAndAddress(suite.DB(), &storageInTransit)
	suite.Nil(err)
	suite.Equal(0, verrs.Count())

	savedStorageInTransit, err := models.FetchStorageInTransitByID(suite.DB(), storageInTransit.ID)

	suite.Equal(storageInTransit.ID, savedStorageInTransit.ID)
	suite.Equal(savedStorageInTransit.Status, storageInTransit.Status)
	suite.Equal(*savedStorageInTransit.Notes, *storageInTransit.Notes)
	suite.Equal(storageInTransit.WarehouseName, savedStorageInTransit.WarehouseName)
	suite.Equal(*storageInTransit.WarehousePhone, *savedStorageInTransit.WarehousePhone)
	suite.Equal(*storageInTransit.WarehouseEmail, *savedStorageInTransit.WarehouseEmail)
}
