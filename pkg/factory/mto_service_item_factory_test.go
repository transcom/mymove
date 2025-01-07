package factory

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildMTOServiceItem() {
	suite.Run("Successful creation of default extended MTOServiceItem", func() {
		// Under test:      BuildMTOServiceItem
		// Mocked:          None
		// Set up:          Create a service item with no customizations or traits
		// Expected outcome:mtoServiceItem should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		mtoServiceItem := BuildMTOServiceItem(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.False(mtoServiceItem.MoveTaskOrderID.IsNil())
		suite.False(mtoServiceItem.MoveTaskOrder.ID.IsNil())
		suite.NotNil(mtoServiceItem.MTOShipmentID)
		suite.False(mtoServiceItem.MTOShipmentID.IsNil())
		suite.False(mtoServiceItem.MTOShipment.ID.IsNil())
		suite.False(mtoServiceItem.ReServiceID.IsNil())
		suite.False(mtoServiceItem.ReService.ID.IsNil())
		suite.Equal(models.MTOServiceItemStatusSubmitted, mtoServiceItem.Status)
	})

	suite.Run("Successful creation of customized MTOServiceItem", func() {
		// Under test:      BuildMTOServiceItem
		// Mocked:          None
		// Set up:          Create a service item with and pass custom fields
		// Expected outcome:mtoServiceItem should be created with custom values

		// SETUP
		customMove := models.Move{
			Locator: "ABC123",
			Show:    models.BoolPointer(true),
		}
		customMTOShipment := models.MTOShipment{
			Status: models.MTOShipmentStatusDraft,
		}
		customReService := models.ReService{
			Name: "Domestic linehaul",
			Code: models.ReServiceCode("DLH"),
		}
		customMtoServiceItem := models.MTOServiceItem{
			Status: models.MTOServiceItemStatusRejected,
		}
		sitOriginalAddress := BuildAddress(suite.DB(), nil, nil)
		sitActualAddress := BuildAddress(suite.DB(), nil, nil)
		sitDestinationOriginalAddress := BuildAddress(suite.DB(), nil, nil)
		sitDestinationAddress := BuildAddress(suite.DB(), nil, nil)
		customs := []Customization{
			{
				Model: customMove,
			},
			{
				Model: customMTOShipment,
			},
			{
				Model: customReService,
			},
			{
				Model: customMtoServiceItem,
			},
			{
				Model:    sitOriginalAddress,
				LinkOnly: true,
				Type:     &Addresses.SITOriginHHGOriginalAddress,
			},
			{
				Model:    sitActualAddress,
				LinkOnly: true,
				Type:     &Addresses.SITOriginHHGActualAddress,
			},
			{
				Model:    sitDestinationAddress,
				LinkOnly: true,
				Type:     &Addresses.SITDestinationFinalAddress,
			},
			{
				Model:    sitDestinationOriginalAddress,
				LinkOnly: true,
				Type:     &Addresses.SITDestinationOriginalAddress,
			},
		}
		// CALL FUNCTION UNDER TEST
		mtoServiceItem := BuildMTOServiceItem(suite.DB(), customs, nil)

		// VALIDATE RESULTS
		suite.False(mtoServiceItem.MoveTaskOrderID.IsNil())
		suite.False(mtoServiceItem.MoveTaskOrder.ID.IsNil())
		suite.Equal(customMove.Locator, mtoServiceItem.MoveTaskOrder.Locator)
		suite.NotNil(mtoServiceItem.MoveTaskOrder.Show)
		suite.True(*mtoServiceItem.MoveTaskOrder.Show)

		suite.NotNil(mtoServiceItem.MTOShipmentID)
		suite.False(mtoServiceItem.MTOShipmentID.IsNil())
		suite.False(mtoServiceItem.MTOShipment.ID.IsNil())
		suite.Equal(customMTOShipment.Status, mtoServiceItem.MTOShipment.Status)

		suite.False(mtoServiceItem.ReServiceID.IsNil())
		suite.False(mtoServiceItem.ReService.ID.IsNil())
		suite.Equal(customReService.Name, mtoServiceItem.ReService.Name)
		suite.Equal(customReService.Code, mtoServiceItem.ReService.Code)

		suite.Equal(customMtoServiceItem.Status, mtoServiceItem.Status)

		suite.NotNil(mtoServiceItem.SITOriginHHGOriginalAddressID)
		suite.Equal(sitOriginalAddress.ID, *mtoServiceItem.SITOriginHHGOriginalAddressID)
		suite.Equal(sitOriginalAddress.ID, mtoServiceItem.SITOriginHHGOriginalAddress.ID)
		suite.NotNil(mtoServiceItem.SITOriginHHGActualAddressID)
		suite.Equal(sitActualAddress.ID, *mtoServiceItem.SITOriginHHGActualAddressID)
		suite.Equal(sitActualAddress.ID, mtoServiceItem.SITOriginHHGActualAddress.ID)
		suite.NotNil(mtoServiceItem.SITDestinationFinalAddressID)
		suite.Equal(sitDestinationAddress.ID, *mtoServiceItem.SITDestinationFinalAddressID)
		suite.Equal(sitDestinationAddress.ID, mtoServiceItem.SITDestinationFinalAddress.ID)
		suite.NotNil(mtoServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(sitDestinationOriginalAddress.ID, *mtoServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(sitDestinationOriginalAddress.ID, mtoServiceItem.SITDestinationOriginalAddress.ID)
	})

	suite.Run("Successful return of linkOnly MTOServiceItem", func() {
		// Under test:       BuildMTOServiceItem
		// Set up:           Pass in a linkOnly mtoServiceItem
		// Expected outcome: No new MTOServiceItem should be created.

		// Check num MTOServiceItem records
		precount, err := suite.DB().Count(&models.MTOServiceItem{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		mtoServiceItem := BuildMTOServiceItem(suite.DB(), []Customization{
			{
				Model: models.MTOServiceItem{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.MTOServiceItem{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, mtoServiceItem.ID)
	})

	suite.Run("Successful return of stubbed MTOServiceItem", func() {
		// Under test:       BuildMTOServiceItem
		// Set up:           Pass in nil db
		// Expected outcome: No new MTOServiceItem should be created.

		// Check num MTOServiceItem records
		precount, err := suite.DB().Count(&models.MTOServiceItem{})
		suite.NoError(err)

		customMtoServiceItem := models.MTOServiceItem{
			Status: models.MTOServiceItemStatusRejected,
		}
		// Nil passed in as db
		mtoServiceItem := BuildMTOServiceItem(nil, []Customization{
			{
				Model: customMtoServiceItem,
			},
		}, nil)

		count, err := suite.DB().Count(&models.MTOServiceItem{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customMtoServiceItem.Status, mtoServiceItem.Status)
	})

	suite.Run("Successful creation of basic MTOServiceItem", func() {
		// Under test:      BuildMTOServiceItemBasic
		// Mocked:          None
		// Set up:          Create a basic service item with no customizations or traits
		// Expected outcome:mtoServiceItem should be created with
		// default values and no shipment

		// SETUP
		// CALL FUNCTION UNDER TEST
		mtoServiceItem := BuildMTOServiceItemBasic(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.False(mtoServiceItem.MoveTaskOrderID.IsNil())
		suite.False(mtoServiceItem.MoveTaskOrder.ID.IsNil())
		suite.Nil(mtoServiceItem.MTOShipmentID)
		suite.True(mtoServiceItem.MTOShipment.ID.IsNil())
		suite.False(mtoServiceItem.ReServiceID.IsNil())
		suite.False(mtoServiceItem.ReService.ID.IsNil())
		suite.Equal(models.MTOServiceItemStatusSubmitted, mtoServiceItem.Status)
	})

	suite.Run("Find or create ReService by code", func() {
		// build 2 service items, they will have the same ReService Code
		si1 := BuildMTOServiceItem(suite.DB(), nil, nil)
		si2 := BuildMTOServiceItem(suite.DB(), nil, nil)
		suite.Equal(si1.ReService.Code, si2.ReService.Code)
	})
	suite.Run("Find or create ReService by custom code", func() {
		// build 2 service items with the same code
		customReService := models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		}
		si1 := BuildMTOServiceItem(suite.DB(), []Customization{
			{
				Model: customReService,
			},
		}, nil)
		si2 := BuildMTOServiceItem(suite.DB(), []Customization{
			{
				Model: customReService,
			},
		}, nil)

		suite.Equal(customReService.Code, si1.ReService.Code)
		suite.Equal(customReService.Code, si2.ReService.Code)
	})

	suite.Run("build real MTOServiceItem with all deps valid code", func() {
		move := BuildMove(suite.DB(), nil, nil)
		shipment := BuildMTOShipment(suite.DB(), []Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		mtoServiceItem := BuildRealMTOServiceItemWithAllDeps(suite.DB(),
			models.ReServiceCodeMS, move, shipment, nil, nil)

		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeMS)

		suite.Equal(move.ID, mtoServiceItem.MoveTaskOrderID)
		suite.NotNil(mtoServiceItem.MTOShipmentID)
		suite.Equal(shipment.ID, *mtoServiceItem.MTOShipmentID)
		suite.Equal(reService.ID, mtoServiceItem.ReServiceID)
		suite.Equal(models.MTOServiceItemStatusApproved, mtoServiceItem.Status)
	})

	suite.Run("build real MTOServiceItem with all deps invalid code", func() {
		move := BuildMove(suite.DB(), nil, nil)
		shipment := BuildMTOShipment(suite.DB(), []Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		buf := bytes.NewBuffer([]byte{})
		log.SetOutput(buf)
		// ensure log output is set back to default
		defer log.SetOutput(os.Stderr)
		suite.Panics(func() {
			// this code is not a supported one
			BuildRealMTOServiceItemWithAllDeps(suite.DB(),
				models.ReServiceCodeDUCRT, move, shipment, nil, nil)
		})
		suite.Contains(buf.String(), "couldn't create service item service code")
	})

	suite.Run("build full DLH MTOServiceItem", func() {
		// the original MakeFullDHLMTOServiceItem did an override of
		// the MTOShipment Status to ensure it was Submitted, but in
		// our conversion, we will require the caller to make that
		// customization
		move, mtoServiceItems := BuildFullDLHMTOServiceItems(suite.DB(), []Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		suite.Equal(1, len(move.MTOShipments))
		suite.Equal(models.MTOShipmentStatusSubmitted, move.MTOShipments[0].Status)
		reServiceCodes := []models.ReServiceCode{}
		for i := range mtoServiceItems {
			reServiceCodes = append(reServiceCodes, mtoServiceItems[i].ReService.Code)
		}
		expectedCodes := []models.ReServiceCode{
			models.ReServiceCodeMS,
			models.ReServiceCodeCS,
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
		}
		suite.Equal(expectedCodes, reServiceCodes)
	})

	suite.Run("build full origin MTOServiceItem", func() {
		// the original MakeFullOriginMTOServiceItem did an override of
		// the MTOShipment Status to ensure it was Submitted, but in
		// our conversion, we will require the caller to make that
		// customization
		move, mtoServiceItems := BuildFullOriginMTOServiceItems(suite.DB(), []Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		suite.Equal(1, len(move.MTOShipments))
		suite.Equal(models.MTOShipmentStatusSubmitted, move.MTOShipments[0].Status)
		reServiceCodes := []models.ReServiceCode{}
		for i := range mtoServiceItems {
			reServiceCodes = append(reServiceCodes, mtoServiceItems[i].ReService.Code)
		}
		expectedCodes := []models.ReServiceCode{
			models.ReServiceCodeDPK,
			models.ReServiceCodeDOP,
		}
		suite.Equal(expectedCodes, reServiceCodes)
	})

	suite.Run("Build SIT service items without SIT out Date", func() {

		customMove := BuildMove(suite.DB(), nil, nil)
		customMTOShipment := BuildMTOShipment(suite.DB(), []Customization{
			{
				Model:    customMove,
				LinkOnly: true,
			},
		}, nil)

		defaultEntryDate := time.Now().AddDate(0, 0, -45)
		sitServiceItems := BuildDestSITServiceItemsNoSITDepartureDate(suite.DB(), customMove, customMTOShipment, &defaultEntryDate)
		reServiceCodes := []models.ReServiceCode{}

		for i := range sitServiceItems {
			reServiceCodes = append(reServiceCodes, sitServiceItems[i].ReService.Code)
			suite.Nil(sitServiceItems[i].SITDepartureDate)
		}
		expectedCodes := []models.ReServiceCode{
			models.ReServiceCodeDDFSIT,
			models.ReServiceCodeDDASIT,
			models.ReServiceCodeDDDSIT,
			models.ReServiceCodeDDSFSC,
		}
		suite.Equal(expectedCodes, reServiceCodes)
	})

	suite.Run("Port Locations not populated by default", func() {

		serviceItem := BuildMTOServiceItem(suite.DB(), nil, nil)
		suite.Nil(serviceItem.POELocation, nil)
		suite.Nil(serviceItem.PODLocation, nil)
	})

	suite.Run("PODLocation populated for PortOfDebarkation type", func() {

		portLocation := FetchPortLocation(suite.DB(), []Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)
		serviceItem := BuildMTOServiceItem(suite.DB(), []Customization{
			{
				Model:    portLocation,
				LinkOnly: true,
				Type:     &PortLocations.PortOfDebarkation,
			},
		}, nil)
		suite.Equal(portLocation.Port.PortCode, serviceItem.PODLocation.Port.PortCode)
		suite.Nil(serviceItem.POELocation)
	})

	suite.Run("POELocation populated for PortOfEmbarkation type", func() {

		portLocation := FetchPortLocation(suite.DB(), []Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)
		serviceItem := BuildMTOServiceItem(suite.DB(), []Customization{
			{
				Model:    portLocation,
				LinkOnly: true,
				Type:     &PortLocations.PortOfEmbarkation,
			},
		}, nil)
		suite.Nil(serviceItem.PODLocation)
		suite.Equal(portLocation.Port.PortCode, serviceItem.POELocation.Port.PortCode)
	})
}
