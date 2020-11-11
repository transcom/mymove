package event

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *EventServiceSuite) Test_MTOServiceItemPayload() {
	now := time.Now()

	mtoServiceItemDOFSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
		ReService: models.ReService{
			Code: "DOFSIT",
			Name: "Destination 1st Day SIT",
		},
	})

	mtoServiceItemDDFSIT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
		ReService: models.ReService{
			Code: "DDFSIT",
			Name: "Destination 1st Day SIT",
		},
	})

	customerContact1 := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			MTOServiceItemID:           mtoServiceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeFirst,
			TimeMilitary:               "0800Z",
			FirstAvailableDeliveryDate: time.Now(),
		},
		ReService: models.ReService{
			Code: "DDFSIT",
			Name: "Destination 1st Day SIT",
		},
	})

	customerContact2 := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			MTOServiceItemID:           mtoServiceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeSecond,
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: time.Now(),
		},
		ReService: models.ReService{
			Code: "DDFSIT",
		},
	})

	mtoServiceItemDCRT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
		ReService: models.ReService{
			Code: "DCRT",
			Name: "Dom. Crating",
		},
	})

	itemDimension1 := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
		MTOServiceItemDimension: models.MTOServiceItemDimension{
			Type:      models.DimensionTypeItem,
			Length:    900,
			Height:    900,
			Width:     900,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		MTOServiceItem: mtoServiceItemDCRT,
	})

	crateDimension1 := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
		MTOServiceItemDimension: models.MTOServiceItemDimension{
			MTOServiceItemID: mtoServiceItemDCRT.ID,
			Type:             models.DimensionTypeCrate,
			Length:           2000,
			Height:           2000,
			Width:            2000,
			CreatedAt:        time.Time{},
			UpdatedAt:        time.Time{},
		},
	})

	testString := "Lorem ipsum"

	mtoServiceItemDOSHUT := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
		},
		ReService: models.ReService{
			Code: "DOSHUT",
		},
		MTOServiceItem: models.MTOServiceItem{
			Description: &testString,
			Reason:      &testString,
		},
	})

	suite.T().Run("Success with MTOServiceItemDOFSIT", func(t *testing.T) {
		data := &primemessages.MTOServiceItemDOFSIT{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.DB(), mtoServiceItemDOFSIT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDOFSIT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDOFSIT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(fmt.Sprintf("%s", mtoServiceItemDOFSIT.ReService.Code), *data.ReServiceCode)
		suite.Equal(mtoServiceItemDOFSIT.Reason, data.Reason)
	})

	suite.T().Run("Success with MTOServiceItemDDFSIT", func(t *testing.T) {
		data := &primemessages.MTOServiceItemDDFSIT{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.DB(), mtoServiceItemDDFSIT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDDFSIT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDDFSIT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(fmt.Sprintf("%s", mtoServiceItemDDFSIT.ReService.Code), *data.ReServiceCode)
		suite.Equal(customerContact1.FirstAvailableDeliveryDate.Format("2006-01-02"), data.FirstAvailableDeliveryDate1.String())
		suite.Equal(customerContact2.FirstAvailableDeliveryDate.Format("2006-01-02"), data.FirstAvailableDeliveryDate2.String())

	})

	suite.T().Run("Success with MTOServiceItemDCRT", func(t *testing.T) {
		data := &primemessages.MTOServiceItemDomesticCrating{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.DB(), mtoServiceItemDCRT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDCRT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDCRT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(fmt.Sprintf("%s", mtoServiceItemDCRT.ReService.Code), *data.ReServiceCode)
		suite.Equal(float32(itemDimension1.Length), float32(*data.Item.Length))
		suite.Equal(float32(crateDimension1.Length), float32(*data.Crate.Length))

	})

	suite.T().Run("Success with MTOServiceItemDOSHUT", func(t *testing.T) {
		data := &primemessages.MTOServiceItemShuttle{}

		payload, assemblePayloadErr := assembleMTOServiceItemPayload(suite.DB(), mtoServiceItemDOSHUT.ID)

		unmarshalErr := data.UnmarshalJSON(payload)

		suite.Nil(assemblePayloadErr)
		suite.Nil(unmarshalErr)
		suite.Equal(mtoServiceItemDOSHUT.ID.String(), data.ID().String())
		suite.Equal(mtoServiceItemDOSHUT.MTOShipmentID.String(), data.MtoShipmentID().String())
		suite.Equal(fmt.Sprintf("%s", mtoServiceItemDOSHUT.ReService.Code), *data.ReServiceCode)
		suite.Equal(*mtoServiceItemDOSHUT.Description, *data.Description)
		suite.Equal(*mtoServiceItemDOSHUT.Reason, *data.Reason)
	})

}

func (suite *EventServiceSuite) TestAssembleMoveOrderPayload() {
	moveOrder := testdatagen.MakeDefaultOrder(suite.DB())

	suite.T().Run("Success with default MoveOrder", func(t *testing.T) {
		payload, err := assembleMoveOrderPayload(suite.DB(), moveOrder.ID)

		data := &primemessages.MoveOrder{}
		unmarshalErr := data.UnmarshalBinary(payload)

		suite.Nil(err)
		suite.Nil(unmarshalErr)
		suite.Equal(moveOrder.ID.String(), data.ID.String())
		suite.NotNil(moveOrder.ServiceMember)
		suite.NotNil(moveOrder.Entitlement)
		suite.NotNil(moveOrder.OriginDutyStation)
		suite.NotEqual(moveOrder.ServiceMember.ID, uuid.Nil)
		suite.NotEqual(moveOrder.Entitlement.ID, uuid.Nil)
		suite.NotEqual(moveOrder.OriginDutyStation.ID, uuid.Nil)

		if moveOrder.OriginDutyStation != nil {
			suite.NotNil(moveOrder.OriginDutyStation.Address)
			suite.NotEqual(moveOrder.OriginDutyStation.Address.ID, uuid.Nil)
		}
	})
}

// Check Success that if thre is an MTO, a notification is STORED, but if there isn't no error is returned.
