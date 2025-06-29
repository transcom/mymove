// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package mtoserviceitem

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	mocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	portlocation "github.com/transcom/mymove/pkg/services/port_location"
	"github.com/transcom/mymove/pkg/services/query"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *MTOServiceItemServiceSuite) TestMTOServiceItemUpdater() {

	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	shipmentRouter := mtoshipment.NewShipmentRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	sitStatusService := sitstatus.NewShipmentSITStatus()
	portLocationFetcher := portlocation.NewPortLocationFetcher()

	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	updater := NewMTOServiceItemUpdater(planner, builder, moveRouter, shipmentRouter, shipmentFetcher, addressCreator, portLocationFetcher, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	setupServiceItem := func() (models.MTOServiceItem, string) {
		serviceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		return serviceItem, eTag
	}

	// Test not found error
	suite.Run("Not Found Error", func() {
		serviceItem, eTag := setupServiceItem()
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundServiceItem := serviceItem
		notFoundServiceItem.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &notFoundServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})

	// Test validation error
	suite.Run("Validation Error", func() {
		serviceItem, eTag := setupServiceItem()
		invalidServiceItem := serviceItem
		invalidServiceItem.MoveTaskOrderID = serviceItem.ID // invalid Move ID

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &invalidServiceItem, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "moveTaskOrderID")
	})

	// Test precondition failed (stale eTag)
	suite.Run("Precondition Failed", func() {
		serviceItem, _ := setupServiceItem()
		newServiceItem := serviceItem
		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &newServiceItem, "bloop")

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	// Test successful update
	suite.Run("Successful update of service item ", func() {
		serviceItem, eTag := setupServiceItem()
		reason := "because we did this service"
		sitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)

		country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
		newAddress := factory.BuildAddress(nil, nil, nil)
		newAddress.Country = &country
		newAddress.CountryId = &country.ID
		newServiceItem := serviceItem
		newServiceItem.Reason = &reason
		newServiceItem.SITEntryDate = &sitEntryDate
		newServiceItem.Status = "" // should keep the status from the original service item
		newServiceItem.SITDestinationFinalAddress = &newAddress
		actualWeight := int64(4000)
		estimatedWeight := int64(4200)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&actualWeight)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&estimatedWeight)

		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &newServiceItem, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(serviceItem.ID, updatedServiceItem.ID)
		suite.Equal(serviceItem.MTOShipmentID, updatedServiceItem.MTOShipmentID)
		suite.Equal(serviceItem.MoveTaskOrderID, updatedServiceItem.MoveTaskOrderID)
		suite.Equal(newServiceItem.Reason, updatedServiceItem.Reason)
		suite.Equal(newServiceItem.SITEntryDate.Local(), updatedServiceItem.SITEntryDate.Local())
		suite.Equal(serviceItem.Status, updatedServiceItem.Status) // should not have been updated
		suite.Equal(newAddress.StreetAddress1, updatedServiceItem.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(newAddress.City, updatedServiceItem.SITDestinationFinalAddress.City)
		suite.Equal(newAddress.State, updatedServiceItem.SITDestinationFinalAddress.State)
		suite.Equal(newAddress.Country, updatedServiceItem.SITDestinationFinalAddress.Country)
		suite.Equal(newAddress.PostalCode, updatedServiceItem.SITDestinationFinalAddress.PostalCode)
		suite.Equal(newServiceItem.ActualWeight, updatedServiceItem.ActualWeight)
		suite.Equal(newServiceItem.EstimatedWeight, updatedServiceItem.EstimatedWeight)
		suite.NotEqual(newServiceItem.Status, updatedServiceItem.Status)
	})

	// Success for DDDSIT
	suite.Run("Successful update of DDDSIT service item", func() {
		serviceItem, eTag := setupServiceItem()
		serviceItem.ReService.Code = models.ReServiceCodeDDDSIT
		reason := "because we did this service"
		sitEntryDate := time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC)

		country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
		newAddress := factory.BuildAddress(nil, nil, nil)
		newAddress.Country = &country
		newAddress.CountryId = &country.ID
		newServiceItem := serviceItem
		newServiceItem.Reason = &reason
		newServiceItem.SITEntryDate = &sitEntryDate
		newServiceItem.Status = "" // should keep the status from the original service item
		newServiceItem.SITDestinationFinalAddress = &newAddress
		actualWeight := int64(4000)
		estimatedWeight := int64(4200)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&actualWeight)
		newServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&estimatedWeight)
		newServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              time.Date(2020, time.December, 04, 0, 0, 0, 0, time.UTC),
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC),
				Type:                       models.CustomerContactTypeFirst,
			},
		}
		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &newServiceItem, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(serviceItem.ID, updatedServiceItem.ID)
		suite.Equal(serviceItem.MTOShipmentID, updatedServiceItem.MTOShipmentID)
		suite.Equal(serviceItem.MoveTaskOrderID, updatedServiceItem.MoveTaskOrderID)
		suite.Equal(newServiceItem.Reason, updatedServiceItem.Reason)
		suite.Equal(newServiceItem.SITEntryDate.Local(), updatedServiceItem.SITEntryDate.Local())
		suite.Equal(serviceItem.Status, updatedServiceItem.Status) // should not have been updated
		suite.Equal(newAddress.StreetAddress1, updatedServiceItem.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(newAddress.City, updatedServiceItem.SITDestinationFinalAddress.City)
		suite.Equal(newAddress.State, updatedServiceItem.SITDestinationFinalAddress.State)
		suite.Equal(newAddress.Country, updatedServiceItem.SITDestinationFinalAddress.Country)
		suite.Equal(newAddress.PostalCode, updatedServiceItem.SITDestinationFinalAddress.PostalCode)
		suite.Equal(newServiceItem.ActualWeight, updatedServiceItem.ActualWeight)
		suite.Equal(newServiceItem.EstimatedWeight, updatedServiceItem.EstimatedWeight)
		suite.Equal(newServiceItem.CustomerContacts[0].TimeMilitary, updatedServiceItem.CustomerContacts[0].TimeMilitary)
		suite.Equal(newServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate, updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
		suite.Equal(newServiceItem.CustomerContacts[0].DateOfContact, updatedServiceItem.CustomerContacts[0].DateOfContact)
		suite.NotEqual(newServiceItem.Status, updatedServiceItem.Status)
	})

	// Success for DDDSIT with an existing customer contact
	suite.Run("Successful update of DDDSIT service item that already has Customer Contacts", func() {
		customerContact := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
			MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				DateOfContact:              time.Date(1984, time.March, 24, 0, 0, 0, 0, time.UTC),
				TimeMilitary:               "0400Z",
				FirstAvailableDeliveryDate: time.Date(1984, time.March, 20, 0, 0, 0, 0, time.UTC),
			},
		})
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact},
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		newServiceItem := serviceItem
		newServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              time.Date(2020, time.December, 04, 0, 0, 0, 0, time.UTC),
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC),
				Type:                       models.CustomerContactTypeFirst,
			},
		}
		updatedServiceItem, err := updater.UpdateMTOServiceItemBasic(suite.AppContextForTest(), &newServiceItem, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.Equal(serviceItem.ID, updatedServiceItem.ID)
		suite.Equal(serviceItem.MTOShipmentID, updatedServiceItem.MTOShipmentID)
		suite.Equal(serviceItem.MoveTaskOrderID, updatedServiceItem.MoveTaskOrderID)

		// We updated the old customer contact, so the ID should be the same
		suite.Equal(customerContact.ID, updatedServiceItem.CustomerContacts[0].ID)

		// And the new values should be reflected in the updated customer contact
		suite.NotEqual(customerContact.TimeMilitary, updatedServiceItem.CustomerContacts[0].TimeMilitary)
		suite.NotEqual(customerContact.DateOfContact, updatedServiceItem.CustomerContacts[0].DateOfContact)
		suite.NotEqual(customerContact.FirstAvailableDeliveryDate, updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
		suite.Equal(newServiceItem.CustomerContacts[0].TimeMilitary, updatedServiceItem.CustomerContacts[0].TimeMilitary)
		suite.Equal(newServiceItem.CustomerContacts[0].DateOfContact, updatedServiceItem.CustomerContacts[0].DateOfContact)
		suite.Equal(newServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate, updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
	})

	suite.Run("Successful Prime update - adding SITDestinationFinalAddress", func() {
		now := time.Now()
		requestApproavalsRequestedStatus := false
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDDFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "APPROVED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "REJECTED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to add SITDestinationFinalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDDFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item
		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipment, eTag)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.NotNil(updatedServiceItem.SITDestinationFinalAddress)
		suite.Equal(newAddress.StreetAddress1, updatedServiceItem.SITDestinationFinalAddress.StreetAddress1)
		suite.Equal(newAddress.StreetAddress2, updatedServiceItem.SITDestinationFinalAddress.StreetAddress2)
		suite.Equal(newAddress.StreetAddress3, updatedServiceItem.SITDestinationFinalAddress.StreetAddress3)
		suite.Equal(newAddress.City, updatedServiceItem.SITDestinationFinalAddress.City)
		suite.Equal(newAddress.State, updatedServiceItem.SITDestinationFinalAddress.State)
		suite.Equal(newAddress.PostalCode, updatedServiceItem.SITDestinationFinalAddress.PostalCode)
	})

	// Test that if a SITDepartureDate is provided successfully and it is a date before the shipments
	// authorized end date then the shipment's end date will be adjusted to be equal to the SITDepartureDate
	// DESTINATION
	suite.Run("Successful Prime update - adding SITDepartureDate adjusts shipment's Destination SIT authorized end date", func() {
		now := time.Now()
		requestApproavalsRequestedStatus := false
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		departureDate := contactDatePlusGracePeriod.Add(time.Hour * 24)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
		}, []factory.Trait{factory.GetTraitAvailableToPrimeMove})
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)

		requestedDays := 90
		officeRemarks := "TESTING REMARKS"
		sitExtension := models.SITDurationUpdate{
			RequestedDays: requestedDays,
			RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
			Status:        models.SITExtensionStatusPending,
			OfficeRemarks: &officeRemarks,
		}

		populatesitExtensions := []models.SITDurationUpdate{sitExtension}

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:                    models.MTOShipmentStatusApproved,
					SITDaysAllowance:          &shipmentSITAllowance,
					PrimeEstimatedWeight:      &estimatedWeight,
					RequiredDeliveryDate:      &aMonthAgo,
					UpdatedAt:                 aMonthAgo,
					SITDurationUpdates:        populatesitExtensions,
					DestinationSITAuthEndDate: &departureDate,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// Link sitExtension for existing shipment
		factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)

		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDDFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "APPROVED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "REJECTED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to add SITDestinationFinalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDDFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Confirm sitExtension exists for the shipment
		var sitExtensions []models.SITDurationUpdate
		suite.DB().Q().All(&sitExtensions)
		suite.Equal(1, len(sitExtensions))
		suite.Equal(models.SITExtensionStatusPending, sitExtensions[0].Status)
		suite.Equal(shipment.ID, sitExtensions[0].MTOShipmentID)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, shipment.MoveTaskOrder.Status)

		// Confirm move status is APPROVALS REQUESTED before sit extension removal
		var moves []models.Move
		suite.DB().Q().All(&moves)
		suite.Equal(1, len(moves))
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moves[0].Status)

		// Update MTO service item
		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipmentWithCalculatedStatus, eTag)
		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)

		// Confirm sitExtension status was updated for the shipment
		suite.DB().Q().All(&sitExtensions)
		suite.Equal(1, len(sitExtensions))
		suite.Equal(models.SITExtensionStatusRemoved, sitExtensions[0].Status)
		// Confirm decision date is set to today
		suite.Equal(time.Now().Truncate(time.Hour*24), sitExtensions[0].DecisionDate.Truncate(time.Hour*24).Local())
		suite.Equal(shipment.ID, sitExtensions[0].MTOShipmentID)

		// Confirm move status is APPROVED after sit extension removal
		suite.DB().Q().All(&moves)
		suite.Equal(1, len(moves))
		suite.Equal(models.MoveStatusAPPROVED, moves[0].Status)

		// Verify that the shipment's SIT authorized end date has been adjusted to be equal
		// to the SIT departure date
		var postUpdatedServiceItemShipment models.MTOShipment
		suite.DB().Q().Find(&postUpdatedServiceItemShipment, shipment.ID)
		suite.NotNil(postUpdatedServiceItemShipment)
		// Verify the departure date is equal to the shipment SIT status departure date (Previously shipment SIT status would have an improper end date due to calc issues. This was fixed in B-20967)
		suite.True(updatedServiceItem.SITDepartureDate.Equal(*shipmentWithCalculatedStatus.DestinationSITAuthEndDate))
		// Verify the updated shipment authorized end date is equal to the departure date
		// Truncate to the nearest day. This is because the shipment only inherits the day, month, year from the service item, not the hour, minute, or second
		suite.True(updatedServiceItem.SITDepartureDate.Truncate(24 * time.Hour).Equal(postUpdatedServiceItemShipment.DestinationSITAuthEndDate.Truncate(24 * time.Hour)))
	})

	// Test that if a SITDepartureDate is provided successfully and it is a date before the shipments
	// authorized end date then the shipment's end date will be adjusted to be equal to the SITDepartureDate
	// ORIGIN
	suite.Run("Successful Prime update - adding SITDepartureDate adjusts shipment's Origin SIT authorized end date", func() {
		now := time.Now()
		requestApproavalsRequestedStatus := false
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDOFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "APPROVED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOPSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "REJECTED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to add SITDestinationFinalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDOFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item
		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipmentWithCalculatedStatus, eTag)
		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)

		// Verify that the shipment's SIT authorized end date has been adjusted to be equal
		// to the SIT departure date
		var postUpdatedServiceItemShipment models.MTOShipment
		suite.DB().Q().Find(&postUpdatedServiceItemShipment, shipment.ID)
		suite.NotNil(postUpdatedServiceItemShipment)
		// Verify the departure date is equal to the shipment SIT status departure date (Previously shipment SIT status would have an improper end date due to calc issues. This was fixed in B-20967)
		suite.True(updatedServiceItem.SITDepartureDate.Equal(*shipmentWithCalculatedStatus.OriginSITAuthEndDate))
		// Verify the updated shipment authorized end date is equal to the departure date
		// Truncate to the nearest day. This is because the shipment only inherits the day, month, year from the service item, not the hour, minute, or second
		suite.True(updatedServiceItem.SITDepartureDate.Truncate(24 * time.Hour).Equal(postUpdatedServiceItemShipment.OriginSITAuthEndDate.Truncate(24 * time.Hour)))

	})

	suite.Run("Unsuccessful Prime update - updating existing SITDestinationFinalAddres", func() {
		now := time.Now()
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
		}, nil)
		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDDFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
					Status:               "APPROVED",
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
					Status:               models.MTOServiceItemStatusRejected,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to update SITDestinationFinalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDDFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item
		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipment, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "SITDestinationFinalAddress")
	})

	suite.Run("Successful Prime update - resubmitting all rejected origin and destination SIT service item", func() {
		now := time.Now()
		requestApprovalsRequestedStatus := false
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		reason := "this is why the service item was created"

		// going to create and test all of these service items
		serviceItemCodes := []models.ReServiceCode{
			models.ReServiceCodeDDFSIT,
			models.ReServiceCodeDDASIT,
			models.ReServiceCodeDDDSIT,
			models.ReServiceCodeDDSFSC,
			models.ReServiceCodeDOASIT,
			models.ReServiceCodeDOPSIT,
			models.ReServiceCodeDOFSIT,
			models.ReServiceCodeDOSFSC,
			models.ReServiceCodeIDFSIT,
			models.ReServiceCodeIDASIT,
			models.ReServiceCodeIDDSIT,
			models.ReServiceCodeIDSFSC,
			models.ReServiceCodeIOASIT,
			models.ReServiceCodeIOPSIT,
			models.ReServiceCodeIOFSIT,
			models.ReServiceCodeIOSFSC,
		}

		shipmentSITAllowance := 90
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.NoError(err)

		// build rejected SIT service items & update them with new reasons or else we will get an error
		for _, code := range serviceItemCodes {
			serviceItem := buildRejectedServiceItem(suite, code, reason, contactDatePlusGracePeriod, aMonthAgo, now, sitRequestedDelivery, requestApprovalsRequestedStatus)
			eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

			updatedServiceItem := serviceItem
			updatedServiceItem.Reason = models.StringPointer("this is a new reason")
			updatedServiceItem.RequestedApprovalsRequestedStatus = models.BoolPointer(true)
			updatedServiceItem.Status = models.MTOServiceItemStatusSubmitted

			updatedServiceItemResult, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &updatedServiceItem, planner, shipment, eTag)

			suite.NoError(err)
			suite.NotNil(updatedServiceItemResult)
			suite.IsType(models.MTOServiceItem{}, *updatedServiceItemResult)
		}
	})

	suite.Run("Unsuccessful Prime update - resubmitting all rejected origin and destination SIT service without updating the reason", func() {
		now := time.Now()
		requestApprovalsRequestedStatus := false
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		reason := "this is why the service item was created"

		// going to create and test all of these service items
		serviceItemCodes := []models.ReServiceCode{
			models.ReServiceCodeDDFSIT,
			models.ReServiceCodeDDASIT,
			models.ReServiceCodeDDDSIT,
			models.ReServiceCodeDDSFSC,
			models.ReServiceCodeDOASIT,
			models.ReServiceCodeDOPSIT,
			models.ReServiceCodeDOFSIT,
			models.ReServiceCodeDOSFSC,
		}

		shipmentSITAllowance := 90
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		suite.NoError(err)

		// build rejected SIT service items & update them with new reasons or else we will get an error
		for _, code := range serviceItemCodes {
			serviceItem := buildRejectedServiceItem(suite, code, reason, contactDatePlusGracePeriod, aMonthAgo, now, sitRequestedDelivery, requestApprovalsRequestedStatus)
			eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

			updatedServiceItem := serviceItem
			updatedServiceItem.RequestedApprovalsRequestedStatus = models.BoolPointer(true)
			updatedServiceItem.Status = models.MTOServiceItemStatusSubmitted

			updatedServiceItemResult, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &updatedServiceItem, planner, shipment, eTag)

			// we should get an error back since the reason MUST be changed
			suite.Nil(updatedServiceItemResult)
			suite.Error(err)
			suite.IsType(apperror.ConflictError{}, err)
		}
	})

	suite.Run("Unsuccessful basic update - adding SITDestinationOriginalAddress", func() {
		now := time.Now()
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDDFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
					Status:               "APPROVED",
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to update SITDestinationOriginalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationOriginalAddress = &newAddress
		newServiceItemPrime.SITDestinationOriginalAddressID = &newAddress.ID

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDDFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item
		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipment, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "SITDestinationOriginalAddress")
	})

	suite.Run("Unsuccessful prime update - adding SITDestinationOriginalAddress", func() {
		now := time.Now()
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDDFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
					Status:               "APPROVED",
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		// Try to update SITDestinationOriginalAddress
		newServiceItemPrime := oldServiceItemPrime
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationOriginalAddress = &newAddress
		newServiceItemPrime.SITDestinationOriginalAddressID = &newAddress.ID

		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDDFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item
		updatedServiceItem, err := updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipment, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "SITDestinationOriginalAddress")
	})
	suite.Run("When TOO converts a SIT to customer expense", func() {
		// Build shipment with SIT
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dofsit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aMonthAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{dofsit}

		// Update ConvertToCustomerExpense and CustomerExpenseReason
		updatedServiceItem, err := updater.ConvertItemToCustomerExpense(
			suite.AppContextForTest(), &approvedShipment, models.StringPointer("test"), true)
		suite.NoError(err)

		// Check the SIT for updated value
		suite.Equal(true, updatedServiceItem.CustomerExpense)
		suite.Equal(models.StringPointer("test"), updatedServiceItem.CustomerExpenseReason)
	})

	suite.Run("failure test for ghc transit time query", func() {
		now := time.Now()
		requestApproavalsRequestedStatus := false
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Now().AddDate(0, 0, 10)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		// Do not provide a custom prime estimated weight
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// We need to create a origin first day sit in order to properly calculate authorized end date
		oldDOFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
					Status:               "APPROVED",
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "REJECTED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		newServiceItemPrime := oldServiceItemPrime
		newServiceItemPrime.Status = models.MTOServiceItemStatusApproved
		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDOFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, newServiceItemPrime)
		_, err = updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipment, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("failure test for ZipTransitDistance", func() {
		now := time.Now()
		requestApproavalsRequestedStatus := false
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Now().AddDate(0, 0, 10)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDOFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
					Status:               "APPROVED",
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:                  &contactDatePlusGracePeriod,
					SITEntryDate:                      &aMonthAgo,
					SITCustomerContacted:              &now,
					SITRequestedDelivery:              &sitRequestedDelivery,
					Status:                            "REJECTED",
					RequestedApprovalsRequestedStatus: &requestApproavalsRequestedStatus,
				},
			},
		}, nil)

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(1234, apperror.UnprocessableEntityError{})

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 1,
			DistanceMilesUpper: 2000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
		eTag := etag.GenerateEtag(oldServiceItemPrime.UpdatedAt)

		newServiceItemPrime := oldServiceItemPrime
		newServiceItemPrime.Status = models.MTOServiceItemStatusApproved
		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDOFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item

		_, err = updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipment, eTag)

		suite.Error(err)
		suite.IsType(apperror.UnprocessableEntityError{}, err)
	})

	suite.Run("Successful update of port service item with updated pricing estimates of basic iHHG service items ", func() {
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"50314",
			"98158",
		).Return(1000, nil)

		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		pickupUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "50314")
		suite.FatalNoError(err)
		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "Tester Address",
					City:               "Des Moines",
					State:              "IA",
					PostalCode:         "50314",
					IsOconus:           models.BoolPointer(false),
					UsPostRegionCityID: &pickupUSPRC.ID,
				},
			},
		}, nil)

		destUSPRC, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99505")
		suite.FatalNoError(err)
		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "Anchorage",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		pickupDate := time.Now()
		requestedPickup := time.Now()
		estimatedWeight := unit.Pound(1212)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
					ScheduledPickupDate:  &pickupDate,
					RequestedPickupDate:  &requestedPickup,
					PrimeEstimatedWeight: &estimatedWeight,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// building service items with NO pricing estimates
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeISLH,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIHUPK,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)
		portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "PDX",
				},
			},
		}, nil)
		poefsc := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
			{
				Model:    portLocation,
				LinkOnly: true,
				Type:     &factory.PortLocations.PortOfEmbarkation,
			},
		}, nil)

		eTag := etag.GenerateEtag(poefsc.UpdatedAt)

		// update the port
		newServiceItemPrime := poefsc
		newServiceItemPrime.POELocation.Port.PortCode = "SEA"

		// Update MTO service item
		_, err = updater.UpdateMTOServiceItemPrime(suite.AppContextForTest(), &newServiceItemPrime, planner, shipment, eTag)
		suite.NoError(err)

		// checking the service item data
		var serviceItems []models.MTOServiceItem
		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		suite.Equal(4, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			// because the estimated weight is provided & POEFSC has a port location now, estimated pricing should be updated
			suite.NotNil(serviceItems[i].PricingEstimate)
		}
	})
}

func (suite *MTOServiceItemServiceSuite) TestValidateUpdateMTOServiceItem() {
	// Set up the data needed for updateMTOServiceItemData obj
	checker := movetaskorder.NewMoveTaskOrderChecker()
	before := time.Now().AddDate(0, 0, -3)
	now := time.Now()
	sitStatusService := sitstatus.NewShipmentSITStatus()

	// Test with bad string key
	suite.Run("bad validatorKey - failure", func() {
		serviceItemData := updateMTOServiceItemData{}
		fakeKey := "FakeKey"
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, fakeKey)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.Contains(err.Error(), fakeKey)
	})

	// Test successful Basic validation
	suite.Run("UpdateMTOServiceItemBasicValidator - success", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		newServiceItem := models.MTOServiceItem{
			ID:              oldServiceItem.ID,
			MTOShipmentID:   oldServiceItem.MTOShipmentID,
			MoveTaskOrderID: oldServiceItem.MoveTaskOrderID,
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemBasicValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})

	// Test unsuccessful Basic validation
	suite.Run("UpdateMTOServiceItemBasicValidator - failure", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		newServiceItem := models.MTOServiceItem{
			ID:            oldServiceItem.ID,
			MTOShipmentID: &oldServiceItem.ID, // bad value
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemBasicValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	// Test successful Prime validation for Port of Embarkation
	suite.Run("UpdateMTOServiceItemPrimeValidator - Update Port of Embarkation - success", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime
		poeId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		newServiceItemPrime.POELocationID = &poeId
		newServiceItemPrime.POELocation = &models.PortLocation{}

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.Equal(updatedServiceItem.POELocationID, newServiceItemPrime.POELocationID)
	})

	// Test success Prime validation for Port of Embarkation
	suite.Run("UpdateMTOServiceItemPrimeValidator - Update Port of Embarkation - Port not updated when new port ID is nil", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
		}, nil)
		poeId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		oldServiceItemPrime.POELocationID = &poeId

		newServiceItemPrime := oldServiceItemPrime
		newServiceItemPrime.POELocationID = nil
		newServiceItemPrime.POELocation = &models.PortLocation{}

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.NotEqual(oldServiceItemPrime.POELocationID, newServiceItemPrime.POELocationID)
		suite.Equal(oldServiceItemPrime.POELocationID, updatedServiceItem.POELocationID)
	})

	// Test failure Prime validation for Port of Embarkation
	suite.Run("UpdateMTOServiceItemPrimeValidator - Update Port of Embarkation - Port not updated wrong service code is supplied", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePODFSC,
				},
			},
		}, nil)
		poeId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		oldServiceItemPrime.POELocationID = &poeId

		newServiceItemPrime := oldServiceItemPrime
		newServiceItemPrime.POELocationID = nil
		newServiceItemPrime.POELocation = &models.PortLocation{}

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Error(err)
		suite.Empty(updatedServiceItem)
		suite.Contains(err.Error(), "is in a conflicting state POE Location can only be updated for service item POEFSC")
	})

	// Test successful Prime validation for Port of Debarkation
	suite.Run("UpdateMTOServiceItemPrimeValidator - Update Port of Debarkation - success", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePODFSC,
				},
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime
		podId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		newServiceItemPrime.PODLocationID = &podId
		newServiceItemPrime.PODLocation = &models.PortLocation{}

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.Equal(updatedServiceItem.PODLocationID, newServiceItemPrime.PODLocationID)
	})

	// Test successful Prime validation for Port of Debarkation
	suite.Run("UpdateMTOServiceItemPrimeValidator - Update Port of Debarkation - Port not updated when new port ID is nil", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePODFSC,
				},
			},
		}, nil)
		podId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		oldServiceItemPrime.PODLocationID = &podId

		newServiceItemPrime := oldServiceItemPrime
		newServiceItemPrime.PODLocationID = nil
		newServiceItemPrime.PODLocation = &models.PortLocation{}

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.NotEqual(oldServiceItemPrime.PODLocationID, newServiceItemPrime.PODLocationID)
		suite.Equal(oldServiceItemPrime.PODLocationID, updatedServiceItem.PODLocationID)
	})

	// Test failure Prime validation for Port of Debarkation
	suite.Run("UpdateMTOServiceItemPrimeValidator - Update Port of Debarkation - Port not updated wrong service code is supplied", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodePOEFSC,
				},
			},
		}, nil)
		podId := uuid.FromStringOrNil("b6e94f5b-33c0-43f3-b960-7c7b2a4ee5fc")
		oldServiceItemPrime.PODLocationID = &podId

		newServiceItemPrime := oldServiceItemPrime
		newServiceItemPrime.PODLocationID = nil
		newServiceItemPrime.PODLocation = &models.PortLocation{}

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Error(err)
		suite.Empty(updatedServiceItem)
		suite.Contains(err.Error(), "is in a conflicting state POD Location can only be updated for service item PODFSC")
	})

	// Test successful Prime validation
	suite.Run("UpdateMTOServiceItemPrimeValidator - success", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Change something allowed by Prime:
		reason := "because"
		newServiceItemPrime.Reason = &reason
		newServiceItemPrime.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				TimeMilitary:               "1300Z",
				DateOfContact:              time.Date(2020, time.December, 04, 0, 0, 0, 0, time.UTC),
				FirstAvailableDeliveryDate: time.Date(2020, time.December, 02, 0, 0, 0, 0, time.UTC),
			},
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.Equal(updatedServiceItem.CustomerContacts[0].TimeMilitary, newServiceItemPrime.CustomerContacts[0].TimeMilitary)
		suite.Equal(updatedServiceItem.CustomerContacts[0].DateOfContact, newServiceItemPrime.CustomerContacts[0].DateOfContact)
		suite.Equal(updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate, newServiceItemPrime.CustomerContacts[0].FirstAvailableDeliveryDate)
	})

	// Test unsuccessful Prime validation - Not available to Prime
	suite.Run("UpdateMTOServiceItemPrimeValidator - not available failure", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		newServiceItemNotPrime := oldServiceItem // this service item should not be Prime-available

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemNotPrime,
			oldServiceItem:      oldServiceItem,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	// Test unsuccessful Prime validation - Invalid input
	suite.Run("UpdateMTOServiceItemPrimeValidator - invalid input failure", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Change something unavailable to Prime:
		newServiceItemPrime.Status = models.MTOServiceItemStatusApproved
		newServiceItemPrime.ApprovedAt = &now

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "status")
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "approvedAt")
	})

	// Test unsuccessful Prime validation - Payment requests
	suite.Run("UpdateMTOServiceItemPrimeValidator - payment request failure", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Create payment requests for service item:
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: paymentRequest, LinkOnly: true,
			}, {
				Model: oldServiceItemPrime, LinkOnly: true,
			},
		}, nil)

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	// Test with empty string key (successful Base validation)
	suite.Run("empty validatorKey - success", func() {
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		newServiceItem := oldServiceItem
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, "")

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
	})

	// Test that when an approved DDDSIT sitDestination is updated the serviceItem stays approved
	suite.Run("UpdateMTOServiceItemPrimeValidator - Successfully Update Approved ServiceItem sitDepartureDate", func() {
		year, month, day := now.Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		contactDatePlusGracePeriod := now.AddDate(0, 0, GracePeriodDays)
		sitRequestedDelivery := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipmentSITAllowance := int(90)
		estimatedWeight := unit.Pound(1400)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:               models.MTOShipmentStatusApproved,
					SITDaysAllowance:     &shipmentSITAllowance,
					PrimeEstimatedWeight: &estimatedWeight,
					RequiredDeliveryDate: &aMonthAgo,
					UpdatedAt:            aMonthAgo,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		// We need to create a destination first day sit in order to properly calculate authorized end date
		oldDDFSITServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate:     &contactDatePlusGracePeriod,
					SITEntryDate:         &aMonthAgo,
					SITCustomerContacted: &now,
					SITRequestedDelivery: &sitRequestedDelivery,
					Status:               "APPROVED",
				},
			},
		}, nil)
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					SITDepartureDate: &now,
					SITEntryDate:     &before,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)

		newServiceItemPrime := oldServiceItemPrime
		newServiceItemPrime.SITEntryDate = nil
		newServiceItemPrime.RequestedApprovalsRequestedStatus = nil

		// Change sitDepartureDate:
		newDate := time.Now().AddDate(0, 0, 5)
		newServiceItemPrime.SITDepartureDate = &newDate

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		// Set shipment SIT status
		shipment.MTOServiceItems = append(shipment.MTOServiceItems, oldServiceItemPrime, oldDDFSITServiceItemPrime)
		sitStatus, shipmentWithCalculatedStatus, err := sitStatusService.CalculateShipmentSITStatus(suite.AppContextForTest(), shipment)
		suite.MustSave(&shipmentWithCalculatedStatus)
		suite.NoError(err)
		suite.NotNil(sitStatus)

		// Update MTO service item
		updatedServiceItem, err := ValidateUpdateMTOServiceItem(suite.AppContextForTest(), &serviceItemData, UpdateMTOServiceItemPrimeValidator)

		suite.NoError(err)
		suite.NotNil(updatedServiceItem)
		suite.IsType(models.MTOServiceItem{}, *updatedServiceItem)
		suite.Equal(updatedServiceItem.Status, models.MTOServiceItemStatusApproved)
	})
}

func (suite *MTOServiceItemServiceSuite) createServiceItem() (string, models.MTOServiceItem, models.Move) {
	move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)

	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) createServiceItemForUnapprovedMove() (string, models.MTOServiceItem, models.Move) {
	move := factory.BuildMove(suite.DB(), nil, nil)

	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) createServiceItemForMoveWithUnacknowledgedAmendedOrders() (string, models.MTOServiceItem, models.Move) {
	storer := storageTest.NewFakeS3Storage(true)
	userUploader, err := uploader.NewUserUploader(storer, 100*uploader.MB)
	suite.NoError(err)
	amendedDocument := factory.BuildDocument(suite.DB(), nil, nil)
	amendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
		{
			Model:    amendedDocument,
			LinkOnly: true,
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   suite.AppContextForTest(),
			},
		},
	}, nil)

	amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
	now := time.Now()
	move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				ExcessWeightQualifiedAt: &now,
			},
		},
		{
			Model:    amendedDocument,
			LinkOnly: true,
			Type:     &factory.Documents.UploadedAmendedOrders,
		},
		{
			Model:    amendedDocument.ServiceMember,
			LinkOnly: true,
		},
	}, nil)

	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

	return eTag, serviceItem, move
}

func (suite *MTOServiceItemServiceSuite) setupAssignmentTestData() (models.MTOServiceItems, []models.OfficeUser, models.Move) {
	officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	officeUser2 := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVALSREQUESTED,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    officeUser1,
			LinkOnly: true,
			Type:     &factory.OfficeUsers.TOOTaskOrderAssignedUser,
		},
		{
			Model:    officeUser2,
			LinkOnly: true,
			Type:     &factory.OfficeUsers.TOODestinationAssignedUser,
		},
	}, nil)

	now := time.Now()
	originServiceItem1 := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDCRT,
			},
		},
		{
			Model: models.MTOServiceItem{
				Status:     models.MTOServiceItemStatusSubmitted,
				ApprovedAt: &now,
			},
		},
	}, nil)
	originServiceItem2 := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDUCRT,
			},
		},
		{
			Model: models.MTOServiceItem{
				Status:     models.MTOServiceItemStatusSubmitted,
				ApprovedAt: &now,
			},
		},
	}, nil)
	destinationServiceItem1 := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		},
		{
			Model: models.MTOServiceItem{
				Status:     models.MTOServiceItemStatusSubmitted,
				ApprovedAt: &now,
			},
		},
	}, nil)
	destinationServiceItem2 := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDSHUT,
			},
		},
		{
			Model: models.MTOServiceItem{
				Status:     models.MTOServiceItemStatusSubmitted,
				ApprovedAt: &now,
			},
		},
	}, nil)

	serviceItems := models.MTOServiceItems{
		originServiceItem1,
		originServiceItem2,
		destinationServiceItem1,
		destinationServiceItem2,
	}
	officeUsers := models.OfficeUsers{
		officeUser1,
		officeUser2,
	}

	return serviceItems, officeUsers, move
}

func (suite *MTOServiceItemServiceSuite) TestUpdateMTOServiceItemStatus() {
	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	shipmentRouter := mtoshipment.NewShipmentRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	portLocationFetcher := portlocation.NewPortLocationFetcher()
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	updater := NewMTOServiceItemUpdater(planner, builder, moveRouter, shipmentRouter, shipmentFetcher, addressCreator, portLocationFetcher, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	rejectionReason := models.StringPointer("")

	// Test that the move's status changes to Approved when the service item's
	// status is no longer SUBMITTED
	suite.Run("When TOO reviews move and approves service item", func() {
		eTag, serviceItem, move := suite.createServiceItem()

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)
		var shipment models.MTOShipment
		err = suite.DB().Find(&shipment, serviceItem.MTOShipmentID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.Equal(models.MTOServiceItemStatusApproved, serviceItem.Status)
		suite.NotNil(serviceItem.ApprovedAt)
		suite.Nil(serviceItem.RejectionReason)
		suite.Nil(serviceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
	})
	suite.Run("Handling assigned user When TOO reviews move and approves service item", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					ProvidesCloseout: true,
				},
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CloseoutOffice,
			},
		}, []roles.RoleType{roles.RoleTypeTOO})

		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &factory.OfficeUsers.TOOTaskOrderAssignedUser,
			},
		}, nil)

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(move.TOOTaskOrderAssignedUser)
		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)
		suite.Nil(move.TOOTaskOrderAssignedID)
		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
	})

	suite.Run("When TOO approves a DDDSIT service item with an existing SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		sitDestinationFinalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
			{
				Model:    sitDestinationFinalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)

		destinationAddress := serviceItem.MTOShipment.DestinationAddress
		suite.Equal(destinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(destinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(destinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(destinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a DDDSIT service item without a SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		suite.NotEqual(shipment.DestinationAddressID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.NotEqual(shipment.DestinationAddress.ID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(shipment.DestinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(shipment.DestinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(shipment.DestinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(shipment.DestinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a DDSFSC service item with an existing SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		sitDestinationFinalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
			{
				Model:    sitDestinationFinalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		destinationAddress := serviceItem.MTOShipment.DestinationAddress
		suite.Equal(destinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(destinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(destinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(destinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a IDSFSC service item with an existing SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		destUSPRC, _ := models.FindByZipCode(suite.AppContextForTest().DB(), "99505")
		sitDestinationFinalAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1:     "JBER",
					City:               "Anchorage",
					State:              "AK",
					PostalCode:         "99505",
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &destUSPRC.ID,
				},
			},
		}, nil)

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDSFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
			{
				Model:    sitDestinationFinalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		destinationAddress := serviceItem.MTOShipment.DestinationAddress
		suite.Equal(destinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(destinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(destinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(destinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a DDSFSC service item without a SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSFSC,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		suite.NotEqual(shipment.DestinationAddressID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.NotEqual(shipment.DestinationAddress.ID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(shipment.DestinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(shipment.DestinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(shipment.DestinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(shipment.DestinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a DDASIT service item with an existing SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		sitDestinationFinalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
			{
				Model:    sitDestinationFinalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		destinationAddress := serviceItem.MTOShipment.DestinationAddress
		suite.Equal(destinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(destinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(destinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(destinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a IDASIT service item with an existing SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		sitDestinationFinalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
			{
				Model:    sitDestinationFinalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		destinationAddress := serviceItem.MTOShipment.DestinationAddress
		suite.Equal(destinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(destinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(destinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(destinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a DDASIT service item without a SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		suite.NotEqual(shipment.DestinationAddressID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.NotEqual(shipment.DestinationAddress.ID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(shipment.DestinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(shipment.DestinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(shipment.DestinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(shipment.DestinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a IDASIT service item without a SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDASIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		suite.NotEqual(shipment.DestinationAddressID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.NotEqual(shipment.DestinationAddress.ID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(shipment.DestinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(shipment.DestinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(shipment.DestinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(shipment.DestinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a DDFSIT service item with an existing SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		sitDestinationFinalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
			{
				Model:    sitDestinationFinalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		destinationAddress := serviceItem.MTOShipment.DestinationAddress
		suite.Equal(destinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(destinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(destinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(destinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a IDFSIT service item with an existing SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		sitDestinationFinalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
			{
				Model:    sitDestinationFinalAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		destinationAddress := serviceItem.MTOShipment.DestinationAddress
		suite.Equal(destinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(destinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(destinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(destinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a DDFSIT service item without a SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		suite.NotEqual(shipment.DestinationAddressID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.NotEqual(shipment.DestinationAddress.ID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(shipment.DestinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(shipment.DestinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(shipment.DestinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(shipment.DestinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	suite.Run("When TOO approves a IDFSIT service item without a SITDestinationFinalAddress", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIDFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusSubmitted,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
		suite.NotEqual(shipment.DestinationAddressID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.NotEqual(shipment.DestinationAddress.ID, *updatedServiceItem.SITDestinationOriginalAddressID)
		suite.Equal(shipment.DestinationAddress.StreetAddress1, updatedServiceItem.SITDestinationOriginalAddress.StreetAddress1)
		suite.Equal(shipment.DestinationAddress.City, updatedServiceItem.SITDestinationOriginalAddress.City)
		suite.Equal(shipment.DestinationAddress.State, updatedServiceItem.SITDestinationOriginalAddress.State)
		suite.Equal(shipment.DestinationAddress.PostalCode, updatedServiceItem.SITDestinationOriginalAddress.PostalCode)
	})

	// Test that the move's status changes to Approvals Requested if any of its service
	// items' status is SUBMITTED
	suite.Run("When move is approved and service item is submitted", func() {
		eTag, serviceItem, move := suite.createServiceItem()
		move.Status = models.MoveStatusAPPROVED
		suite.MustSave(&move)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusSubmitted, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItem.Status)
		suite.Nil(serviceItem.RejectionReason)
		suite.Nil(serviceItem.RejectedAt)
		suite.Nil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
	})

	// Test that the move's status changes to Approved if the service item is
	// rejected
	suite.Run("When TOO reviews move and rejects service item", func() {
		eTag, serviceItem, move := suite.createServiceItem()
		rejectionReason = models.StringPointer("incomplete")

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusRejected, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)
		var shipment models.MTOShipment
		err = suite.DB().Find(&shipment, serviceItem.MTOShipmentID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
		suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
		suite.Equal(models.MTOServiceItemStatusRejected, serviceItem.Status)
		suite.Equal(rejectionReason, serviceItem.RejectionReason)
		suite.NotNil(serviceItem.RejectedAt)
		suite.Nil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
	})

	// Test that a service item's status can only be updated if the Move's status
	// is either Approved or Approvals Requested. Neither the Move's status nor
	// the service item's status should be changed if the requirements aren't met.
	suite.Run("When the Move has not been approved yet", func() {
		eTag, serviceItem, move := suite.createServiceItemForUnapprovedMove()

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)

		suite.Error(err)
		suite.Contains(err.Error(), "Cannot approve or reject a service item if the move's status is neither Approved nor Approvals Requested.")

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusDRAFT, move.Status)
		suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItem.Status)
		suite.Nil(updatedServiceItem)
	})

	suite.Run("does not approve the move if unacknowledged amended orders exist", func() {

		eTag, serviceItem, move := suite.createServiceItemForMoveWithUnacknowledgedAmendedOrders()
		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		err = suite.DB().Find(&serviceItem, serviceItem.ID)
		suite.NoError(err)

		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)
		suite.Equal(models.MTOServiceItemStatusApproved, serviceItem.Status)
		suite.Nil(serviceItem.RejectionReason)
		suite.Nil(serviceItem.RejectedAt)
		suite.NotNil(serviceItem.ApprovedAt)
		suite.NotNil(updatedServiceItem)
	})

	suite.Run("Returns an error when eTag is stale", func() {
		_, serviceItem, _ := suite.createServiceItem()
		rejectionReason = models.StringPointer("incomplete")

		_, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusRejected, rejectionReason, "")

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), serviceItem.ID.String())
	})

	suite.Run("When TOO rejects a DOFSIT service item and converts it to the customer expense", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITOriginHHGActualAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITOriginHHGOriginalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		var shipment models.MTOShipment
		err = suite.DB().Find(&shipment, serviceItem.MTOShipmentID)
		suite.NoError(err)
		suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
	})

	suite.Run("When TOO rejects a IOFSIT service item and converts it to the customer expense", func() {
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOFSIT,
				},
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITOriginHHGActualAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITOriginHHGOriginalAddress,
			},
		}, nil)

		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)

		updatedServiceItem, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), serviceItem.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag)
		suite.NoError(err)

		// ApproveOrRejectServiceItem doesn't return the service item with the updated move
		// get move from the db to check the updated status
		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVED, move.Status)

		var shipment models.MTOShipment
		err = suite.DB().Find(&shipment, serviceItem.MTOShipmentID)
		suite.NoError(err)
		suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)

		suite.Equal(models.MTOServiceItemStatusApproved, updatedServiceItem.Status)
		suite.NotNil(updatedServiceItem.ApprovedAt)
		suite.Nil(updatedServiceItem.RejectionReason)
		suite.Nil(updatedServiceItem.RejectedAt)
		suite.NotNil(updatedServiceItem)
	})

	suite.Run("Returns a not found error if the updater can't find the ReService code for DOFSIT in the DB.", func() {
		_, err := updater.ConvertItemToCustomerExpense(
			suite.AppContextForTest(), &models.MTOShipment{}, models.StringPointer("test"), true)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns a not found error if the updater can't find the MTO Shipment in the DB.", func() {
		// Create ReService in DB so that ConvertItemToCustomerExpense makes it to the MTO Shipment check.
		testdatagen.FetchReService(suite.DB(), testdatagen.Assertions{ReService: models.ReService{Code: "DOFSIT"}})
		_, err := updater.ConvertItemToCustomerExpense(
			suite.AppContextForTest(), &models.MTOShipment{}, models.StringPointer("test"), true)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Handles TOO unassignment properly", func() {
		serviceItems, officeUsers, move := suite.setupAssignmentTestData()

		officeUser1 := officeUsers[0]
		officeUser2 := officeUsers[1]

		originServiceItem1 := serviceItems[0]
		originServiceItem2 := serviceItems[1]
		eTag1 := etag.GenerateEtag(originServiceItem1.UpdatedAt)
		eTag2 := etag.GenerateEtag(originServiceItem2.UpdatedAt)

		// confirm move has origin and destination assignments
		suite.Equal(officeUser1.ID, *move.TOOTaskOrderAssignedID)
		suite.Equal(officeUser2.ID, *move.TOODestinationAssignedID)

		_, err := updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), originServiceItem1.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag1)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)

		// confirm assignments have not changed
		suite.Equal(officeUser1.ID, *move.TOOTaskOrderAssignedID)
		suite.Equal(officeUser2.ID, *move.TOODestinationAssignedID)

		_, err = updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), originServiceItem2.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag2)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)

		// confirm origin TOO is now unassigned and destination TOO remains assigned
		suite.Nil(move.TOOTaskOrderAssignedID)
		suite.Equal(officeUser2.ID, *move.TOODestinationAssignedID)

		destinationServiceItem1 := serviceItems[2]
		destinationServiceItem2 := serviceItems[3]
		eTag3 := etag.GenerateEtag(destinationServiceItem1.UpdatedAt)
		eTag4 := etag.GenerateEtag(destinationServiceItem2.UpdatedAt)

		_, err = updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), destinationServiceItem1.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag3)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)

		// confirm destination TOO remains assigned
		suite.Equal(officeUser2.ID, *move.TOODestinationAssignedID)

		_, err = updater.ApproveOrRejectServiceItem(
			suite.AppContextForTest(), destinationServiceItem2.ID, models.MTOServiceItemStatusApproved, rejectionReason, eTag4)
		suite.NoError(err)

		err = suite.DB().Find(&move, move.ID)
		suite.NoError(err)

		// confirm destination TOO is now unassigned
		suite.Nil(move.TOODestinationAssignedID)
	})
}

func (suite *MTOServiceItemServiceSuite) setupServiceItemData() {
	startDate := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)
	endDate := time.Date(2020, time.December, 31, 12, 0, 0, 0, time.UTC)

	testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			StartDate: startDate,
			EndDate:   endDate,
		},
	})

	originalDomesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea:      "004",
			ServicesSchedule: 2,
		},
		ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
	})

	testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            originalDomesticServiceArea.Contract,
			ContractID:          originalDomesticServiceArea.ContractID,
			DomesticServiceArea: originalDomesticServiceArea,
			Zip3:                "902",
		},
	})

	testdatagen.FetchOrMakeReDomesticLinehaulPrice(suite.DB(), testdatagen.Assertions{
		ReDomesticLinehaulPrice: models.ReDomesticLinehaulPrice{
			Contract:              originalDomesticServiceArea.Contract,
			ContractID:            originalDomesticServiceArea.ContractID,
			DomesticServiceArea:   originalDomesticServiceArea,
			DomesticServiceAreaID: originalDomesticServiceArea.ID,
			WeightLower:           unit.Pound(500),
			WeightUpper:           unit.Pound(9999),
			MilesLower:            250,
			MilesUpper:            9999,
			PriceMillicents:       unit.Millicents(606800),
			IsPeakPeriod:          false,
		},
	})
}

func (suite *MTOServiceItemServiceSuite) TestUpdateMTOServiceItemPricingEstimate() {
	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	shipmentRouter := mtoshipment.NewShipmentRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	updater := NewMTOServiceItemUpdater(planner, builder, moveRouter, shipmentRouter, shipmentFetcher, addressCreator, portlocation.NewPortLocationFetcher(), ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	setupServiceItem := func() (models.MTOServiceItem, string) {
		serviceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
		return serviceItem, eTag
	}

	setupServiceItems := func() models.MTOServiceItems {
		serviceItems := testdatagen.MakeMTOServiceItems(suite.DB())
		return serviceItems
	}

	suite.Run("Validation Error", func() {
		suite.setupServiceItemData()
		serviceItem, eTag := setupServiceItem()
		invalidServiceItem := serviceItem
		invalidServiceItem.MoveTaskOrderID = serviceItem.ID // invalid Move ID

		updatedServiceItem, err := updater.UpdateMTOServiceItemPricingEstimate(suite.AppContextForTest(), &invalidServiceItem, serviceItem.MTOShipment, eTag)

		suite.Nil(updatedServiceItem)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		invalidInputError := err.(apperror.InvalidInputError)
		suite.True(invalidInputError.ValidationErrors.HasAny())
		suite.Contains(invalidInputError.ValidationErrors.Keys(), "moveTaskOrderID")
	})

	suite.Run("Returns updated service item on success wihtout error", func() {
		suite.setupServiceItemData()
		serviceItems := setupServiceItems()

		for _, serviceItem := range serviceItems {
			eTag := etag.GenerateEtag(serviceItem.UpdatedAt)
			updatedServiceItem, err := updater.UpdateMTOServiceItemPricingEstimate(suite.AppContextForTest(), &serviceItem, serviceItem.MTOShipment, eTag)

			suite.NotNil(updatedServiceItem)
			suite.Nil(err)
		}
	})
}

// Helper function to create a rejected service item
func buildRejectedServiceItem(suite *MTOServiceItemServiceSuite, reServiceCode models.ReServiceCode, reason string, contactDatePlusGracePeriod, aMonthAgo, now, sitRequestedDelivery time.Time, requestApprovalsRequestedStatus bool) models.MTOServiceItem {
	return factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: reServiceCode,
			},
		},
		{
			Model: models.MTOServiceItem{
				SITDepartureDate:                  &contactDatePlusGracePeriod,
				SITEntryDate:                      &aMonthAgo,
				SITCustomerContacted:              &now,
				SITRequestedDelivery:              &sitRequestedDelivery,
				Status:                            models.MTOServiceItemStatusRejected,
				RequestedApprovalsRequestedStatus: &requestApprovalsRequestedStatus,
				Reason:                            &reason,
			},
		},
	}, nil)
}
