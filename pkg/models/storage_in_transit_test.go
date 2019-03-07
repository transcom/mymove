package models_test

import (
	"testing"

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
