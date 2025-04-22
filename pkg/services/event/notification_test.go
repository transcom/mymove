package event

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EventServiceSuite) Test_MTOServiceItemPayload() {
	now := time.Now()

	suite.Run("Success with MTOServiceItemDOFSIT", func() {
		// Under test: assembleMTOServiceItemPayload
		// Mocked:     None
		// Set up:     Create a DOFSIT in the db, assemble the webhook notification payload
		// Expected outcome: Payload should contain the DOFSIT details
		mtoServiceItemDOFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
					Name: "Destination 1st Day SIT",
				},
			},
		}, nil)
		data := &primemessages.MTOServiceItemOriginSIT{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.AppContextForTest(), mtoServiceItemDOFSIT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDOFSIT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDOFSIT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(string(mtoServiceItemDOFSIT.ReService.Code), *data.ReServiceCode)
		suite.Equal(mtoServiceItemDOFSIT.Reason, data.Reason)
	})

	suite.Run("Success with MTOServiceItemDDFSIT", func() {
		// Under test: assembleMTOServiceItemPayload
		// Set up:     Create a DDFSIT in the db, assemble the webhook notification payload
		// Expected outcome: Payload should contain the DDFSIT details

		customerContact1, err := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
			MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				DateOfContact:              time.Now(),
				TimeMilitary:               "0800Z",
				FirstAvailableDeliveryDate: time.Now(),
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
				Name: "Destination 1st Day SIT",
			},
		})
		suite.NoError(err)

		customerContact2, err := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
			MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeSecond,
				DateOfContact:              time.Now(),
				TimeMilitary:               "0400Z",
				FirstAvailableDeliveryDate: time.Now(),
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		})
		suite.NoError(err)

		mtoServiceItemDDFSIT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
					Name: "Destination 1st Day SIT",
				},
			},
			{
				Model: models.MTOServiceItem{
					CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact1, customerContact2},
				},
			},
		}, nil)

		data := &primemessages.MTOServiceItemDestSIT{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.AppContextForTest(), mtoServiceItemDDFSIT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDDFSIT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDDFSIT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(string(mtoServiceItemDDFSIT.ReService.Code), *data.ReServiceCode)
		suite.Equal(customerContact1.FirstAvailableDeliveryDate.Format("2006-01-02"), data.FirstAvailableDeliveryDate1.String())
		suite.Equal(customerContact2.FirstAvailableDeliveryDate.Format("2006-01-02"), data.FirstAvailableDeliveryDate2.String())

	})

	suite.Run("Success with MTOServiceItemDCRT", func() {
		// Under test: assembleMTOServiceItemPayload
		// Set up:     Create a DCRT in the db, assemble the webhook notification payload
		// Expected outcome: Payload should contain the DCRT details
		mtoServiceItemDCRT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDCRT,
					Name: "Dom. Crating",
				},
			},
		}, nil)

		itemDimension1 := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeItem,
					Length:    900,
					Height:    900,
					Width:     900,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			{
				Model:    mtoServiceItemDCRT,
				LinkOnly: true,
			},
		}, nil)

		crateDimension1 := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeCrate,
					Length:    2000,
					Height:    2000,
					Width:     2000,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			{
				Model:    mtoServiceItemDCRT,
				LinkOnly: true,
			},
		}, nil)
		data := &primemessages.MTOServiceItemDomesticCrating{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.AppContextForTest(), mtoServiceItemDCRT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDCRT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDCRT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(string(mtoServiceItemDCRT.ReService.Code), *data.ReServiceCode)
		suite.Equal(float32(itemDimension1.Length), float32(*data.Item.Length))
		suite.Equal(float32(crateDimension1.Length), float32(*data.Crate.Length))

	})

	suite.Run("Success with MTOServiceItemDOSHUT", func() {
		testString := "Lorem ipsum"

		mtoServiceItemDOSHUT := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &now,
					ApprovedAt:         &now,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOSHUT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Description: &testString,
					Reason:      &testString,
				},
			},
		}, nil)
		data := &primemessages.MTOServiceItemDomesticShuttle{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.AppContextForTest(), mtoServiceItemDOSHUT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDOSHUT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDOSHUT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(string(mtoServiceItemDOSHUT.ReService.Code), *data.ReServiceCode)
		suite.Equal(*mtoServiceItemDOSHUT.Reason, *data.Reason)
	})

}

func (suite *EventServiceSuite) TestAssembleOrderPayload() {

	suite.Run("Success with default Order", func() {
		order := factory.BuildOrder(suite.DB(), nil, nil)
		payload, err := assembleOrderPayload(suite.AppContextForTest(), order.ID)
		suite.FatalNoError(err)

		data := &primemessages.Order{}
		unmarshalErr := data.UnmarshalBinary(payload)

		suite.Nil(unmarshalErr)
		suite.Equal(order.ID.String(), data.ID.String())
		suite.NotNil(order.ServiceMember)
		suite.NotNil(order.Entitlement)
		suite.NotNil(order.OriginDutyLocation)
		suite.NotEqual(order.ServiceMember.ID, uuid.Nil)
		suite.NotEqual(order.Entitlement.ID, uuid.Nil)
		suite.NotEqual(order.OriginDutyLocation.ID, uuid.Nil)

		if order.OriginDutyLocation != nil {
			suite.NotNil(order.OriginDutyLocation.Address)
			suite.NotEqual(order.OriginDutyLocation.Address.ID, uuid.Nil)
		}
	})
}

func (suite *EventServiceSuite) TestAssembleMTOShipmentPayload() {
	suite.Run("Non-external shipment returns payload with all associations", func() {
		// Setup test data
		pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    secondaryPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model:    destinationAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
			{
				Model:    secondaryDeliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: models.TimePointer(time.Now()),
					ApprovedAt:         models.TimePointer(time.Now()),
				},
			},
		}, nil)

		agent := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)
		// Test the assemble function
		payload, shouldNotify, err := assembleMTOShipmentPayload(suite.AppContextForTest(), shipment.ID)
		suite.Nil(err)
		suite.True(shouldNotify)

		data := &primemessages.MTOShipment{}
		unmarshalErr := data.UnmarshalBinary(payload)
		suite.Nil(unmarshalErr)

		suite.Equal(shipment.ID.String(), data.ID.String())
		suite.Equal(shipment.PickupAddress.ID.String(), data.PickupAddress.ID.String())
		suite.Equal(shipment.SecondaryPickupAddress.ID.String(), data.SecondaryPickupAddress.ID.String())
		suite.Equal(shipment.DestinationAddress.ID.String(), data.DestinationAddress.ID.String())
		suite.Equal(shipment.SecondaryDeliveryAddress.ID.String(), data.SecondaryDeliveryAddress.ID.String())
		suite.Equal(agent.ID.String(), data.Agents[0].ID.String())
	})

	suite.Run("External shipment reports that it should not notify", func() {
		// Setup test data
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					UsesExternalVendor: true,
				},
			},
			{
				Model: models.Move{
					AvailableToPrimeAt: models.TimePointer(time.Now()),
					ApprovedAt:         models.TimePointer(time.Now()),
				},
			},
		}, nil)

		// Test the assemble function
		payload, shouldNotify, err := assembleMTOShipmentPayload(suite.AppContextForTest(), shipment.ID)
		suite.Nil(err)
		suite.False(shouldNotify)
		suite.Nil(payload)
	})
}
