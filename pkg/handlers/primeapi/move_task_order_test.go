package primeapi

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/upload"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestListMovesHandlerReturnsUpdated() {
	now := time.Now()
	lastFetch := now.Add(-time.Second)

	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

	// this move should not be returned
	olderMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

	// Pop will overwrite UpdatedAt when saving a model, so use SQL to set it in the past
	suite.Require().NoError(suite.DB().RawQuery("UPDATE moves SET updated_at=? WHERE id=?",
		now.Add(-2*time.Second), olderMove.ID).Exec())
	suite.Require().NoError(suite.DB().RawQuery("UPDATE orders SET updated_at=$1 WHERE id=$2;",
		now.Add(-10*time.Second), olderMove.OrdersID).Exec())

	since := handlers.FmtDateTime(lastFetch)
	request := httptest.NewRequest("GET", fmt.Sprintf("/moves?since=%s", since.String()), nil)
	params := movetaskorderops.ListMovesParams{HTTPRequest: request, Since: since}
	handlerConfig := suite.HandlerConfig()

	// Validate incoming payload: no body to validate

	// make the request
	handler := ListMovesHandler{HandlerConfig: handlerConfig, MoveTaskOrderFetcher: movetaskorder.NewMoveTaskOrderFetcher()}
	response := handler.Handle(params)

	suite.IsNotErrResponse(response)
	listMovesResponse := response.(*movetaskorderops.ListMovesOK)
	movesList := listMovesResponse.Payload

	// Validate outgoing payload
	suite.NoError(movesList.Validate(strfmt.Default))

	suite.Equal(1, len(movesList))
	suite.Equal(move.ID.String(), movesList[0].ID.String())
}

func (suite *HandlerSuite) TestGetMoveTaskOrder() {
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)

	verifyAddressFields := func(address *models.Address, payload *primemessages.Address) {
		suite.Equal(address.ID.String(), payload.ID.String())
		suite.Equal(address.StreetAddress1, *payload.StreetAddress1)
		suite.Equal(*address.StreetAddress2, *payload.StreetAddress2)
		suite.Equal(*address.StreetAddress3, *payload.StreetAddress3)
		suite.Equal(address.City, *payload.City)
		suite.Equal(address.State, *payload.State)
		suite.Equal(address.PostalCode, *payload.PostalCode)
		suite.Equal(*address.Country, *payload.Country)
		suite.NotNil(payload.ETag)
	}

	suite.Run("Success with Prime-available move by ID", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}

		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.ID.String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(movePayload.ID.String(), successMove.ID.String())
		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt) // checks that the date is not 0001-01-01
	})

	suite.Run("Success with Prime-available move by Locator", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(movePayload.ID.String(), successMove.ID.String())
		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt) // checks that the date is not 0001-01-01
	})

	suite.Run("Success returns reweighs on shipments if they exist", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		now := time.Now()
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			Move: successMove,
			Reweigh: models.Reweigh{
				VerificationReason:     models.StringPointer("Justification"),
				VerificationProvidedAt: &nowDate,
				Weight:                 models.PoundPointer(4000),
			},
		})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		reweighPayload := movePayload.MtoShipments[0].Reweigh
		suite.Equal(movePayload.ID.String(), successMove.ID.String())
		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), reweighPayload.ID)
		suite.Equal(reweigh.RequestedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(&reweighPayload.RequestedAt).Format(time.RFC3339))
		suite.Equal(string(reweigh.RequestedBy), string(reweighPayload.RequestedBy))
		suite.Equal(*reweigh.VerificationReason, *reweighPayload.VerificationReason)
		suite.Equal(reweigh.VerificationProvidedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(reweighPayload.VerificationProvidedAt).Format(time.RFC3339))
		suite.Equal(*reweigh.Weight, *handlers.PoundPtrFromInt64Ptr(reweighPayload.Weight))
		suite.Equal(reweigh.ShipmentID.String(), reweighPayload.ShipmentID.String())

		suite.NotNil(reweighPayload.ETag)
		suite.NotNil(reweighPayload.CreatedAt)
		suite.NotNil(reweighPayload.UpdatedAt)
	})

	suite.Run("Success - returns sit extensions on shipments if they exist", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		sitUpdate := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model: models.SITDurationUpdate{
					ContractorRemarks: models.StringPointer("customer wasn't able to finalize apartment"),
					OfficeRemarks:     models.StringPointer("customer mentioned they were finalizing an apt"),
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		sitUpdatePayload := movePayload.MtoShipments[0].SitExtensions[0]
		suite.Equal(successMove.ID.String(), movePayload.ID.String())
		suite.Equal(sitUpdate.ID.String(), sitUpdatePayload.ID.String())
		suite.Equal(sitUpdate.MTOShipmentID.String(), sitUpdatePayload.MtoShipmentID.String())
		suite.Equal(string(sitUpdate.RequestReason), string(sitUpdatePayload.RequestReason))
		suite.Equal(*sitUpdate.ContractorRemarks, *sitUpdatePayload.ContractorRemarks)
		suite.Equal(string(sitUpdate.Status), fmt.Sprintf("%v", sitUpdatePayload.Status))
		suite.Equal(int64(sitUpdate.RequestedDays), sitUpdatePayload.RequestedDays)
		suite.Equal(*handlers.FmtIntPtrToInt64(sitUpdate.ApprovedDays), *sitUpdatePayload.ApprovedDays)
		suite.Equal(sitUpdate.DecisionDate.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(sitUpdatePayload.DecisionDate).Format(time.RFC3339))
		suite.Equal(*sitUpdate.OfficeRemarks, *sitUpdatePayload.OfficeRemarks)

		suite.NotNil(sitUpdatePayload.ETag)
		suite.NotNil(sitUpdatePayload.CreatedAt)
		suite.NotNil(sitUpdatePayload.UpdatedAt)
	})

	suite.Run("Success - returns SitDestinationFinalAddress on related MTO service Items if they exist", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		address := factory.BuildAddress(suite.DB(), nil, nil)
		sitEntryDate := time.Now()

		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status:       models.MTOServiceItemStatusApproved,
					SITEntryDate: &sitEntryDate,
				},
			},
			{
				Model:    address,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT, // DDFSIT - Domestic destination 1st day SIT
				},
			},
		}, nil)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(successMove.ID.String(), movePayload.ID.String())
		if suite.Len(movePayload.MtoServiceItems(), 1) {
			serviceItem := movePayload.MtoServiceItems()[0]

			// Take the service item and marshal it into json
			raw, err := json.Marshal(serviceItem)
			suite.NoError(err)

			// Take that raw json and unmarshal it into a MTOServiceItemDestSIT
			ddfsitServiceItem := primemessages.MTOServiceItemDestSIT{}
			err = ddfsitServiceItem.UnmarshalJSON(raw)
			suite.NoError(err)

			suite.Equal(address.StreetAddress1, *ddfsitServiceItem.SitDestinationFinalAddress.StreetAddress1)
			suite.Equal(address.State, *ddfsitServiceItem.SitDestinationFinalAddress.State)
			suite.Equal(address.City, *ddfsitServiceItem.SitDestinationFinalAddress.City)
		}

	})

	suite.Run("Success - filters shipments handled by an external vendor", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// Create two shipments, one prime, one external.  Only prime one should be returned.
		primeShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: false,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
					UsesExternalVendor: true,
				},
			},
		}, nil)

		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(move.ID.String(), movePayload.ID.String())
		if suite.Len(movePayload.MtoShipments, 1) {
			suite.Equal(primeShipment.ID.String(), movePayload.MtoShipments[0].ID.String())
		}
	})

	suite.Run("Success - returns shipment with attached PpmShipment", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      move.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(move.ID.String(), movePayload.ID.String())
		suite.NotNil(movePayload.MtoShipments[0].PpmShipment)
		suite.Equal(ppmShipment.ShipmentID.String(), movePayload.MtoShipments[0].PpmShipment.ShipmentID.String())
		suite.Equal(ppmShipment.ID.String(), movePayload.MtoShipments[0].PpmShipment.ID.String())
	})

	suite.Run("Success - returns all the fields at the mtoShipment level", func() {
		// This tests fields that aren't other structs and Addresses
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
		destinationType := models.DestinationTypeHomeOfRecord
		secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, nil)
		secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		now := time.Now()
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		yesterDate := nowDate.AddDate(0, 0, -1)
		aWeekAgo := nowDate.AddDate(0, 0, -7)
		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ActualDeliveryDate:               &nowDate,
					CounselorRemarks:                 models.StringPointer("LGTM"),
					DestinationAddressID:             &destinationAddress.ID,
					DestinationType:                  &destinationType,
					FirstAvailableDeliveryDate:       &yesterDate,
					Status:                           models.MTOShipmentStatusApproved,
					NTSRecordedWeight:                models.PoundPointer(unit.Pound(249)),
					PrimeEstimatedWeight:             models.PoundPointer(unit.Pound(980)),
					PrimeEstimatedWeightRecordedDate: &aWeekAgo,
					RequiredDeliveryDate:             &nowDate,
					ScheduledDeliveryDate:            &nowDate,
				},
			},
			{
				Model:    secondaryDeliveryAddress,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
				LinkOnly: true,
			},
			{
				Model:    secondaryPickupAddress,
				Type:     &factory.Addresses.SecondaryPickupAddress,
				LinkOnly: true,
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(movePayload.ID.String(), successMove.ID.String())

		shipment := movePayload.MtoShipments[0]
		suite.Equal(successShipment.ID, handlers.FmtUUIDToPop(shipment.ID))
		suite.Equal(successShipment.ActualDeliveryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.ActualDeliveryDate).Format(time.RFC3339))
		suite.Equal(successShipment.ActualPickupDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.ActualPickupDate).Format(time.RFC3339))
		suite.Equal(successShipment.ApprovedDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.ApprovedDate).Format(time.RFC3339))

		suite.Equal(*successShipment.CounselorRemarks, *shipment.CounselorRemarks)
		suite.Equal(*successShipment.CustomerRemarks, *shipment.CustomerRemarks)

		suite.Equal(destinationAddress.ID, handlers.FmtUUIDToPop(shipment.DestinationAddress.ID))
		verifyAddressFields(&destinationAddress, &shipment.DestinationAddress.Address)

		suite.Equal(string(*successShipment.DestinationType), string(*shipment.DestinationType))

		suite.Equal(successShipment.Diversion, shipment.Diversion)
		suite.Equal(successShipment.FirstAvailableDeliveryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.FirstAvailableDeliveryDate).Format(time.RFC3339))

		suite.Equal(successShipment.MoveTaskOrderID, handlers.FmtUUIDToPop(shipment.MoveTaskOrderID))

		suite.Equal(*successShipment.NTSRecordedWeight, *handlers.PoundPtrFromInt64Ptr(shipment.NtsRecordedWeight))
		verifyAddressFields(successShipment.PickupAddress, &shipment.PickupAddress.Address)

		// TODO: test fields on PpmShipment, existing test "Success - returns shipment with attached PpmShipment"

		suite.Equal(*successShipment.PrimeActualWeight, *handlers.PoundPtrFromInt64Ptr(shipment.PrimeActualWeight))
		suite.Equal(*successShipment.PrimeEstimatedWeight, *handlers.PoundPtrFromInt64Ptr(shipment.PrimeEstimatedWeight))

		suite.Equal(successShipment.PrimeEstimatedWeightRecordedDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.PrimeEstimatedWeightRecordedDate).Format(time.RFC3339))
		suite.Equal(successShipment.RequestedPickupDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.RequestedPickupDate).Format(time.RFC3339))
		suite.Equal(successShipment.RequiredDeliveryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.RequiredDeliveryDate).Format(time.RFC3339))
		suite.Equal(successShipment.RequestedDeliveryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.RequestedDeliveryDate).Format(time.RFC3339))

		suite.Equal(successShipment.ScheduledDeliveryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.ScheduledDeliveryDate).Format(time.RFC3339))
		suite.Equal(successShipment.ScheduledPickupDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(shipment.ScheduledPickupDate).Format(time.RFC3339))
		verifyAddressFields(successShipment.SecondaryDeliveryAddress, &shipment.SecondaryDeliveryAddress.Address)

		verifyAddressFields(successShipment.SecondaryPickupAddress, &shipment.SecondaryPickupAddress.Address)

		suite.Equal(string(successShipment.ShipmentType), string(shipment.ShipmentType))
		suite.Equal(string(successShipment.Status), shipment.Status)

		suite.NotNil(shipment.ETag)
		suite.Equal(successShipment.CreatedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(&shipment.CreatedAt).Format(time.RFC3339))
		suite.Equal(successShipment.UpdatedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(&shipment.UpdatedAt).Format(time.RFC3339))

		suite.NotNil(movePayload.AvailableToPrimeAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt) // checks that the date is not 0001-01-01
	})

	suite.Run("Success - returns all the fields associated with StorageFacility within MtoShipments", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)

		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
					Status:       models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		storageFacilityPayload := movePayload.MtoShipments[0].StorageFacility
		suite.Equal(successMove.ID.String(), movePayload.ID.String())
		suite.Equal(successShipment.StorageFacilityID.String(), storageFacilityPayload.ID.String())
		suite.Equal(successShipment.StorageFacility.ID.String(), storageFacilityPayload.ID.String())
		suite.Equal(successShipment.StorageFacility.FacilityName, storageFacilityPayload.FacilityName)
		suite.Equal(*successShipment.StorageFacility.LotNumber, *storageFacilityPayload.LotNumber)
		suite.Equal(*successShipment.StorageFacility.Phone, *storageFacilityPayload.Phone)
		suite.Equal(*successShipment.StorageFacility.Email, *storageFacilityPayload.Email)

		verifyAddressFields(&successShipment.StorageFacility.Address, storageFacilityPayload.Address)

		suite.NotNil(storageFacilityPayload.ETag)
	})

	suite.Run("Success - returns all the fields associated with Agents within MtoShipments", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		agent := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		agentPayload := movePayload.MtoShipments[0].Agents[0]
		suite.Equal(successMove.ID.String(), movePayload.ID.String())
		suite.Equal(agent.MTOShipmentID.String(), agentPayload.MtoShipmentID.String())
		suite.Equal(agent.ID.String(), agentPayload.ID.String())
		suite.Equal(*agent.FirstName, *agentPayload.FirstName)
		suite.Equal(*agent.LastName, *agentPayload.LastName)
		suite.Equal(*agent.Email, *agentPayload.Email)
		suite.Equal(*agent.Phone, *agentPayload.Phone)
		suite.Equal(string(agent.MTOAgentType), string(agentPayload.AgentType))

		suite.NotNil(agentPayload.ETag)
		suite.NotNil(agentPayload.CreatedAt)
		suite.NotNil(agentPayload.UpdatedAt)
	})

	suite.Run("Success - return all base fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		now := time.Now()
		aWeekAgo := now.AddDate(0, 0, -7)
		upload := factory.BuildUpload(suite.DB(), nil, nil)
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeCounselingCompletedAt: &aWeekAgo,
					PPMEstimatedWeight:         models.PoundPointer(1000),
					ExcessWeightQualifiedAt:    &aWeekAgo,
					ExcessWeightAcknowledgedAt: &now,
					ExcessWeightUploadID:       &upload.ID,
				},
			},
		}, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(successMove.ID.String(), movePayload.ID.String())
		suite.Equal(successMove.Locator, movePayload.MoveCode)
		suite.Equal(successMove.OrdersID.String(), movePayload.OrderID.String())
		suite.Equal(*successMove.ReferenceID, movePayload.ReferenceID)
		suite.Equal(successMove.AvailableToPrimeAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.AvailableToPrimeAt).Format(time.RFC3339))
		suite.Equal(successMove.PrimeCounselingCompletedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.PrimeCounselingCompletedAt).Format(time.RFC3339))
		suite.Equal(*successMove.PPMType, movePayload.PpmType)
		suite.Equal(*handlers.FmtPoundPtr(successMove.PPMEstimatedWeight), movePayload.PpmEstimatedWeight)
		suite.Equal(successMove.ExcessWeightQualifiedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.ExcessWeightQualifiedAt).Format(time.RFC3339))
		suite.Equal(successMove.ExcessWeightAcknowledgedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.ExcessWeightAcknowledgedAt).Format(time.RFC3339))
		suite.Equal(successMove.ExcessWeightUploadID.String(), movePayload.ExcessWeightUploadID.String())

		suite.NotNil(movePayload.CreatedAt)
		suite.NotNil(movePayload.UpdatedAt)
		suite.NotNil(movePayload.ETag)
	})

	suite.Run("Success - return all Order fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		currentAddress := factory.BuildAddress(suite.DB(), nil, nil)
		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model:    currentAddress,
				Type:     &factory.Addresses.ResidentialAddress,
				LinkOnly: true,
			},
		}, nil)
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		ordersPayload := movePayload.Order
		orders := successMove.Orders
		suite.Equal(orders.ID.String(), ordersPayload.ID.String())
		suite.Equal(orders.ServiceMemberID.String(), ordersPayload.CustomerID.String())
		suite.Equal(*orders.OriginDutyLocationGBLOC, ordersPayload.OriginDutyLocationGBLOC)
		suite.Equal(*orders.Grade, *ordersPayload.Rank)
		suite.Equal(orders.ReportByDate.Format(time.RFC3339), time.Time(ordersPayload.ReportByDate).Format(time.RFC3339))
		suite.Equal(string(orders.OrdersType), string(ordersPayload.OrdersType))
		suite.Equal(*orders.OrdersNumber, *ordersPayload.OrderNumber)
		suite.Equal(*orders.TAC, *ordersPayload.LinesOfAccounting)

		suite.NotNil(ordersPayload.ETag)

		// verify customer object aka service member
		suite.Equal(orders.ServiceMember.ID.String(), ordersPayload.Customer.ID.String())
		suite.Equal(*orders.ServiceMember.Edipi, ordersPayload.Customer.DodID)
		suite.Equal(orders.ServiceMember.UserID.String(), ordersPayload.Customer.UserID.String())

		verifyAddressFields(orders.ServiceMember.ResidentialAddress, ordersPayload.Customer.CurrentAddress)

		suite.Equal(*orders.ServiceMember.FirstName, ordersPayload.Customer.FirstName)
		suite.Equal(*orders.ServiceMember.LastName, ordersPayload.Customer.LastName)
		suite.Equal(string(*orders.ServiceMember.Affiliation), ordersPayload.Customer.Branch)
		suite.Equal(*orders.ServiceMember.Telephone, ordersPayload.Customer.Phone)
		suite.Equal(*orders.ServiceMember.PersonalEmail, ordersPayload.Customer.Email)
		suite.NotNil(ordersPayload.Customer.ETag)

		// verify entitlement object
		suite.Equal(orders.Entitlement.ID.String(), ordersPayload.Entitlement.ID.String())
		suite.Equal(int64(*orders.Entitlement.DBAuthorizedWeight), *ordersPayload.Entitlement.AuthorizedWeight)
		suite.Equal(*orders.Entitlement.DependentsAuthorized, *ordersPayload.Entitlement.DependentsAuthorized)
		suite.Equal(*orders.Entitlement.NonTemporaryStorage, *ordersPayload.Entitlement.NonTemporaryStorage)
		suite.Equal(*orders.Entitlement.PrivatelyOwnedVehicle, *ordersPayload.Entitlement.PrivatelyOwnedVehicle)
		suite.Equal(int64(orders.Entitlement.ProGearWeight), ordersPayload.Entitlement.ProGearWeight)
		suite.Equal(int64(orders.Entitlement.ProGearWeightSpouse), ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.Equal(int64(orders.Entitlement.RequiredMedicalEquipmentWeight), ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
		suite.Equal(orders.Entitlement.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(int64(*orders.Entitlement.StorageInTransit), ordersPayload.Entitlement.StorageInTransit)
		suite.Equal(int64(*orders.Entitlement.WeightAllowance()), ordersPayload.Entitlement.TotalWeight)
		suite.Equal(int64(*orders.Entitlement.TotalDependents), ordersPayload.Entitlement.TotalDependents)
		suite.NotNil(ordersPayload.Entitlement.ETag)

		// verify destinationDutyLocation object
		suite.Equal(orders.NewDutyLocation.ID.String(), ordersPayload.DestinationDutyLocation.ID.String())
		suite.Equal(orders.NewDutyLocation.Name, ordersPayload.DestinationDutyLocation.Name)
		suite.Equal(orders.NewDutyLocation.AddressID.String(), ordersPayload.DestinationDutyLocation.AddressID.String())

		verifyAddressFields(&orders.NewDutyLocation.Address, ordersPayload.DestinationDutyLocation.Address)

		suite.NotNil(ordersPayload.DestinationDutyLocation.ETag)

		// verify originDutyLocation object
		suite.Equal(orders.OriginDutyLocation.ID.String(), ordersPayload.OriginDutyLocation.ID.String())
		suite.Equal(orders.OriginDutyLocation.Name, ordersPayload.OriginDutyLocation.Name)
		suite.Equal(orders.OriginDutyLocation.AddressID.String(), ordersPayload.OriginDutyLocation.AddressID.String())

		verifyAddressFields(&orders.OriginDutyLocation.Address, ordersPayload.OriginDutyLocation.Address)
		suite.NotNil(ordersPayload.OriginDutyLocation.ETag)
	})

	suite.Run("Failure 'Not Found' for non-available move", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}
		failureMove := factory.BuildMove(suite.DB(), nil, nil) // default is not available to Prime
		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      failureMove.ID.String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderNotFound{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderNotFound)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Contains(*movePayload.Detail, failureMove.ID.String())
	})
}

func (suite *HandlerSuite) TestCreateExcessWeightRecord() {
	request := httptest.NewRequest("POST", "/move-task-orders/{moveTaskOrderID}", nil)
	fakeS3 := storageTest.NewFakeS3Storage(true)

	suite.Run("Success - Created an excess weight record", func() {
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handler := CreateExcessWeightRecordHandler{
			handlerConfig,
			// Must use the Prime service object in particular:
			moverouter.NewPrimeMoveExcessWeightUploader(upload.NewUploadCreator(fakeS3)),
		}

		now := time.Now()
		availableMove := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ExcessWeightQualifiedAt: &now,
				},
			},
		}, nil)

		params := movetaskorderops.CreateExcessWeightRecordParams{
			HTTPRequest:     request,
			File:            suite.Fixture("test.pdf"),
			MoveTaskOrderID: strfmt.UUID(availableMove.ID.String()),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Require().IsType(&movetaskorderops.CreateExcessWeightRecordCreated{}, response)

		okResponse := response.(*movetaskorderops.CreateExcessWeightRecordCreated)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(availableMove.ID.String(), okResponse.Payload.MoveID.String())
		suite.NotNil(okResponse.Payload.MoveExcessWeightQualifiedAt)
		suite.Equal(okResponse.Payload.MoveExcessWeightQualifiedAt.String(), strfmt.DateTime(*availableMove.ExcessWeightQualifiedAt).String())
		suite.NotEmpty(okResponse.Payload.ID)
	})

	suite.Run("Fail - Move not found - 404", func() {
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handler := CreateExcessWeightRecordHandler{
			handlerConfig,
			// Must use the Prime service object in particular:
			moverouter.NewPrimeMoveExcessWeightUploader(upload.NewUploadCreator(fakeS3)),
		}

		params := movetaskorderops.CreateExcessWeightRecordParams{
			HTTPRequest:     request,
			File:            suite.Fixture("test.pdf"),
			MoveTaskOrderID: strfmt.UUID("00000000-0000-0000-0000-000000000123"),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Require().IsType(&movetaskorderops.CreateExcessWeightRecordNotFound{}, response)
		notFoundResponse := response.(*movetaskorderops.CreateExcessWeightRecordNotFound)

		// Validate outgoing payload
		suite.NoError(notFoundResponse.Payload.Validate(strfmt.Default))

		suite.Require().NotNil(notFoundResponse.Payload.Detail)
		suite.Contains(*notFoundResponse.Payload.Detail, params.MoveTaskOrderID.String())
	})

	suite.Run("Fail - Move not Prime-available - 404", func() {
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handler := CreateExcessWeightRecordHandler{
			handlerConfig,
			// Must use the Prime service object in particular:
			moverouter.NewPrimeMoveExcessWeightUploader(upload.NewUploadCreator(fakeS3)),
		}

		unavailableMove := factory.BuildMove(suite.DB(), nil, nil) // default move is not available to Prime
		params := movetaskorderops.CreateExcessWeightRecordParams{
			HTTPRequest:     request,
			File:            suite.Fixture("test.pdf"),
			MoveTaskOrderID: strfmt.UUID(unavailableMove.ID.String()),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.Require().IsType(&movetaskorderops.CreateExcessWeightRecordNotFound{}, response)
		notFoundResponse := response.(*movetaskorderops.CreateExcessWeightRecordNotFound)

		// Validate outgoing payload
		suite.NoError(notFoundResponse.Payload.Validate(strfmt.Default))

		suite.Require().NotNil(notFoundResponse.Payload.Detail)
		suite.Contains(*notFoundResponse.Payload.Detail, unavailableMove.ID.String())
	})
}

func (suite *HandlerSuite) TestUpdateMTOPostCounselingInfo() {

	suite.Run("Successful patch - Integration Test", func() {
		requestUser := factory.BuildUser(nil, nil, nil)
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(mto.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", mto.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := movetaskorderops.UpdateMTOPostCounselingInformationParams{
			HTTPRequest:     req,
			MoveTaskOrderID: mto.ID.String(),
			IfMatch:         eTag,
		}
		// Create two shipments, one prime, one external.  Only prime one should be returned.
		primeShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: false,
				},
			},
		}, nil)
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
					UsesExternalVendor: true,
				},
			},
		}, nil)
		factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS, // CS - Counseling Services
				},
			},
		}, nil)

		queryBuilder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(queryBuilder)
		moveRouter := moverouter.NewMoveRouter()
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
		updater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter)
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			suite.HandlerConfig(),
			fetcher,
			updater,
			mtoChecker,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationOK{}, response)

		okResponse := response.(*movetaskorderops.UpdateMTOPostCounselingInformationOK)
		okPayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(mto.ID.String(), okPayload.ID.String())
		suite.NotNil(okPayload.ETag)

		if suite.Len(okPayload.MtoShipments, 1) {
			suite.Equal(primeShipment.ID.String(), okPayload.MtoShipments[0].PpmShipment.ID.String())
			suite.Equal(primeShipment.ShipmentID.String(), okPayload.MtoShipments[0].ID.String())
		}

		suite.NotNil(okPayload.PrimeCounselingCompletedAt)
		suite.Equal(primemessages.PPMShipmentStatusWAITINGONCUSTOMER, okPayload.MtoShipments[0].PpmShipment.Status)
	})

	suite.Run("Unsuccessful patch - Integration Test - patch fail MTO not available", func() {
		requestUser := factory.BuildUser(nil, nil, nil)
		defaultMTO := factory.BuildMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(defaultMTO.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", defaultMTO.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		defaultMTOParams := movetaskorderops.UpdateMTOPostCounselingInformationParams{
			HTTPRequest:     req,
			MoveTaskOrderID: defaultMTO.ID.String(),
			IfMatch:         eTag,
		}

		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		fetcher := fetch.NewFetcher(queryBuilder)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
		updater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter)
		handler := UpdateMTOPostCounselingInformationHandler{
			suite.HandlerConfig(),
			fetcher,
			updater,
			mtoChecker,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(defaultMTOParams)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationNotFound{}, response)
		payload := response.(*movetaskorderops.UpdateMTOPostCounselingInformationNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Patch failure - 500", func() {
		requestUser := factory.BuildUser(nil, nil, nil)
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(mto.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", mto.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			suite.HandlerConfig(),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		internalServerErr := errors.New("ServerError")
		params := movetaskorderops.UpdateMTOPostCounselingInformationParams{
			HTTPRequest:     req,
			MoveTaskOrderID: mto.ID.String(),
			IfMatch:         eTag,
		}

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.AnythingOfType("*appcontext.appContext"),
			mto.ID,
			eTag,
		).Return(nil, internalServerErr)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationInternalServerError{}, response)
		payload := response.(*movetaskorderops.UpdateMTOPostCounselingInformationInternalServerError).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Patch failure - 404", func() {
		requestUser := factory.BuildUser(nil, nil, nil)
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(mto.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", mto.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			suite.HandlerConfig(),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}
		params := movetaskorderops.UpdateMTOPostCounselingInformationParams{
			HTTPRequest:     req,
			MoveTaskOrderID: mto.ID.String(),
			IfMatch:         eTag,
		}

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.AnythingOfType("*appcontext.appContext"),
			mto.ID,
			eTag,
		).Return(nil, apperror.NotFoundError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationNotFound{}, response)
		payload := response.(*movetaskorderops.UpdateMTOPostCounselingInformationNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Patch failure - 409", func() {
		requestUser := factory.BuildUser(nil, nil, nil)
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(mto.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", mto.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			suite.HandlerConfig(),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}
		params := movetaskorderops.UpdateMTOPostCounselingInformationParams{
			HTTPRequest:     req,
			MoveTaskOrderID: mto.ID.String(),
			IfMatch:         eTag,
		}
		mockUpdater.On("UpdatePostCounselingInfo",
			mock.AnythingOfType("*appcontext.appContext"),
			mto.ID,
			eTag,
		).Return(nil, apperror.ConflictError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationConflict{}, response)
		payload := response.(*movetaskorderops.UpdateMTOPostCounselingInformationConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Patch failure - 422", func() {
		requestUser := factory.BuildUser(nil, nil, nil)
		mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(mto.UpdatedAt)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/post-counseling-info", mto.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MoveTaskOrderUpdater{}
		mtoChecker := movetaskorder.NewMoveTaskOrderChecker()

		handler := UpdateMTOPostCounselingInformationHandler{
			suite.HandlerConfig(),
			&mockFetcher,
			&mockUpdater,
			mtoChecker,
		}

		mockUpdater.On("UpdatePostCounselingInfo",
			mock.AnythingOfType("*appcontext.appContext"),
			mto.ID,
			eTag,
		).Return(nil, apperror.NewInvalidInputError(uuid.Nil, nil, validate.NewErrors(), ""))
		params := movetaskorderops.UpdateMTOPostCounselingInformationParams{
			HTTPRequest:     req,
			MoveTaskOrderID: mto.ID.String(),
			IfMatch:         eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&movetaskorderops.UpdateMTOPostCounselingInformationUnprocessableEntity{}, response)
		payload := response.(*movetaskorderops.UpdateMTOPostCounselingInformationUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
