package serviceparamvaluelookups

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestMTOEarliestRequestedPickup() {
	key := models.ServiceItemParamNameMTOEarliestRequestedPickup

	earliestRequestedPickup := time.Date(2024, time.March, 15, 0, 0, 0, 452487000, time.Local)
	laterRequestedPickup := time.Date(2025, time.November, 1, 0, 0, 0, 0, time.Local)
	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	var paramLookup *ServiceItemParamKeyData

	setupTestData := func() {
		shipment1 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &earliestRequestedPickup,
				},
			},
		}, nil)
		shipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &laterRequestedPickup,
				},
			},
		}, nil)
		shipment3 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &earliestRequestedPickup,
					DeletedAt:           models.TimePointer(time.Now()),
				},
			},
		}, nil)

		shipments := models.MTOShipments{shipment1, shipment2, shipment3}

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: models.TimePointer(time.Now()),
					MTOShipments:       shipments,
				},
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model:    mtoServiceItem.MTOShipment,
				LinkOnly: true,
			},
		}, nil)

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().Add(time.Hour * 12),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})

		var err error
		paramLookup, err = ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
	}

	suite.Run("golden path", func() {
		setupTestData()

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := earliestRequestedPickup.Format(ghcrateengine.TimestampParamFormat)
		suite.Equal(expected, valueStr)
	})

	suite.Run("bogus MoveTaskOrderID", func() {
		setupTestData()

		// Pass in a non-existent MoveTaskOrderID
		invalidMoveTaskOrderID := uuid.Must(uuid.NewV4())
		badParamLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, invalidMoveTaskOrderID, nil)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Nil(badParamLookup)
	})

	suite.Run("no valid shipments", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &laterRequestedPickup,
					ShipmentType:        models.MTOShipmentTypeHHG,
					Status:              models.MTOShipmentStatusSubmitted,
					DeletedAt:           models.TimePointer(time.Now()),
				},
			},
		}, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: models.TimePointer(time.Now()),
				},
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)

		_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		print(err.Error())
		suite.IsType(apperror.ConflictError{}, err)
	})
}
