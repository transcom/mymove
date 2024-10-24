package primeapiv3

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primev3api/primev3operations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestGetMoveTaskOrder() {
	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)

	verifyAddressFields := func(address *models.Address, payload *primev3messages.Address) {
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
		tertiaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, nil)
		tertiaryPickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
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
				Model:    tertiaryDeliveryAddress,
				Type:     &factory.Addresses.TertiaryDeliveryAddress,
				LinkOnly: true,
			},
			{
				Model:    tertiaryPickupAddress,
				Type:     &factory.Addresses.TertiaryPickupAddress,
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

		verifyAddressFields(successShipment.SecondaryDeliveryAddress, shipment.SecondaryDeliveryAddress)
		verifyAddressFields(successShipment.SecondaryPickupAddress, shipment.SecondaryPickupAddress)
		verifyAddressFields(successShipment.TertiaryDeliveryAddress, shipment.TertiaryDeliveryAddress)
		verifyAddressFields(successShipment.TertiaryPickupAddress, shipment.TertiaryPickupAddress)

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
		suite.Equal(successMove.ExcessWeightQualifiedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.ExcessWeightQualifiedAt).Format(time.RFC3339))
		suite.Equal(successMove.ExcessWeightAcknowledgedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.ExcessWeightAcknowledgedAt).Format(time.RFC3339))
		suite.Equal(successMove.ExcessWeightUploadID.String(), movePayload.ExcessWeightUploadID.String())
		suite.Equal(successMove.Contractor.ContractNumber, movePayload.ContractNumber)

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
		suite.Equal(string(*orders.Grade), string(*ordersPayload.Rank))
		suite.Equal(orders.ReportByDate.Format(time.RFC3339), time.Time(ordersPayload.ReportByDate).Format(time.RFC3339))
		suite.Equal(string(orders.OrdersType), string(ordersPayload.OrdersType))
		suite.Equal(*orders.OrdersNumber, *ordersPayload.OrderNumber)
		suite.Equal(*orders.TAC, *ordersPayload.LinesOfAccounting)

		suite.NotNil(ordersPayload.ETag)

		suite.Equal(orders.SupplyAndServicesCostEstimate, ordersPayload.SupplyAndServicesCostEstimate)
		suite.Equal(orders.PackingAndShippingInstructions, ordersPayload.PackingAndShippingInstructions)
		suite.Equal(orders.MethodOfPayment, ordersPayload.MethodOfPayment)
		suite.Equal(orders.NAICS, ordersPayload.Naics)

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

	suite.Run("Success - return all PaymentRequests fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}

		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					MTOShipmentID: &successShipment.ID,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model:    successShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDLH,
				},
			},
		}, nil)
		recalcPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					SequenceNumber: 2,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:                         true,
					Status:                          models.PaymentRequestStatusReviewed,
					RejectionReason:                 models.StringPointer("no good"),
					SequenceNumber:                  1,
					RecalculationOfPaymentRequestID: &recalcPaymentRequest.ID,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)

		paymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   "TEST",
			},
			{
				Key:     models.ServiceItemParamNameMTOAvailableToPrimeAt,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   "2023-05-03T14:38:30Z",
			},
		}
		paymentServiceItem1 := factory.BuildPaymentServiceItemWithParams(suite.DB(), serviceItem.ReService.Code, paymentServiceItemParams, []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					RejectionReason: models.StringPointer("rejection reason"),
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)

		paymentServiceItem2 := factory.BuildPaymentServiceItemWithParams(suite.DB(), models.ReServiceCodeMS, paymentServiceItemParams, []factory.Customization{
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)

		proofOfServiceDoc := factory.BuildProofOfServiceDoc(suite.DB(), []factory.Customization{
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)

		uploads := factory.BuildPrimeUpload(suite.DB(), []factory.Customization{
			{
				Model:    proofOfServiceDoc,
				LinkOnly: true,
			},
		}, nil)

		paymentRequest.PaymentServiceItems = models.PaymentServiceItems{paymentServiceItem1, paymentServiceItem2}
		proofOfServiceDoc.PrimeUploads = models.PrimeUploads{uploads}
		paymentRequest.ProofOfServiceDocs = models.ProofOfServiceDocs{proofOfServiceDoc}

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

		suite.Len(movePayload.PaymentRequests, 2)
		paymentRequestPayload := movePayload.PaymentRequests[0]
		suite.Equal(paymentRequest.ID.String(), paymentRequestPayload.ID.String())
		suite.Equal(successMove.ID.String(), paymentRequestPayload.MoveTaskOrderID.String())
		suite.Equal(paymentRequest.IsFinal, *paymentRequestPayload.IsFinal)
		suite.Equal(*paymentRequest.RejectionReason, *paymentRequestPayload.RejectionReason)
		suite.Equal(paymentRequest.Status.String(), string(paymentRequestPayload.Status))
		suite.Equal(paymentRequest.PaymentRequestNumber, paymentRequestPayload.PaymentRequestNumber)
		suite.Equal(paymentRequest.RecalculationOfPaymentRequestID.String(), paymentRequestPayload.RecalculationOfPaymentRequestID.String())

		// verify paymentServiceItems
		suite.Len(paymentRequestPayload.PaymentServiceItems, 2)
		PSI1 := paymentRequest.PaymentServiceItems[0]
		PSI1Payload := paymentRequestPayload.PaymentServiceItems[0]
		suite.Equal(PSI1.ID.String(), PSI1Payload.ID.String())
		suite.Equal(PSI1.PaymentRequestID.String(), PSI1Payload.PaymentRequestID.String())
		suite.Equal(PSI1.MTOServiceItemID.String(), PSI1Payload.MtoServiceItemID.String())
		suite.Equal(PSI1.Status.String(), string(PSI1Payload.Status))
		suite.Equal(*handlers.FmtCost(PSI1.PriceCents), *PSI1Payload.PriceCents)
		suite.Equal(*PSI1.RejectionReason, *PSI1Payload.RejectionReason)
		suite.Equal(PSI1.ReferenceID, PSI1Payload.ReferenceID)
		suite.NotNil(PSI1Payload.ETag)
		// verify payment service Items
		suite.Len(PSI1Payload.PaymentServiceItemParams, 2)
		PSIP1 := PSI1.PaymentServiceItemParams[0]
		PSIP1Payload := PSI1Payload.PaymentServiceItemParams[0]
		suite.Equal(PSIP1.ID.String(), PSIP1Payload.ID.String())
		suite.Equal(PSIP1.PaymentServiceItemID.String(), PSIP1Payload.PaymentServiceItemID.String())
		suite.Equal(PSIP1.ServiceItemParamKey.Key.String(), string(PSIP1Payload.Key))
		suite.Equal(PSIP1.Value, PSIP1Payload.Value)
		suite.Equal(PSIP1.ServiceItemParamKey.Type.String(), string(PSIP1Payload.Type))
		suite.Equal(PSIP1.ServiceItemParamKey.Origin.String(), string(PSIP1Payload.Origin))
		suite.NotNil(PSIP1Payload.ETag)

		// verify proofOfServiceDocs
		upload := paymentRequest.ProofOfServiceDocs[0].PrimeUploads[0].Upload
		uploadPayload := paymentRequestPayload.ProofOfServiceDocs[0].Uploads[0]
		suite.Equal(upload.ID.String(), uploadPayload.ID.String())
		suite.Equal(upload.Filename, *uploadPayload.Filename)
		suite.Equal(upload.Bytes, *uploadPayload.Bytes)
		suite.Equal(upload.ContentType, *uploadPayload.ContentType)
		suite.Empty(uploadPayload.URL)
		suite.Empty(uploadPayload.Status)
		suite.NotNil(uploadPayload.CreatedAt)
		suite.NotNil(uploadPayload.UpdatedAt)

		suite.NotNil(paymentRequestPayload.ETag)
	})

	suite.Run("Success - return all MTOServiceItemBasic fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}

		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					RejectionReason: models.StringPointer("not applicable"),
					MTOShipmentID:   &successShipment.ID,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model:    successShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeCS,
				},
			},
		}, nil)

		// Validate incoming payload: no body to validate

		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Len(movePayload.MtoServiceItems(), 1)

		serviceItemPayload := movePayload.MtoServiceItems()[0]

		json, err := json.Marshal(serviceItemPayload)
		suite.NoError(err)
		payload := primev3messages.MTOServiceItemBasic{}
		err = payload.UnmarshalJSON(json)
		suite.NoError(err)

		suite.Equal(serviceItem.MoveTaskOrderID.String(), payload.MoveTaskOrderID().String())
		suite.Equal(serviceItem.MTOShipmentID.String(), payload.MtoShipmentID().String())
		suite.Equal(serviceItem.ID.String(), payload.ID().String())
		suite.Equal("MTOServiceItemBasic", string(payload.ModelType()))
		suite.Equal(string(serviceItem.ReService.Code), string(*payload.ReServiceCode))
		suite.Equal(serviceItem.ReService.Name, payload.ReServiceName())
		suite.Equal(string(serviceItem.Status), string(payload.Status()))
		suite.Equal(*serviceItem.RejectionReason, *payload.RejectionReason())

		suite.NotNil(payload.ETag())
	})

	suite.Run("Success - return all MTOServiceItemOriginSIT fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}

		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)

		now := time.Now()
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		later := nowDate.AddDate(0, 0, 3) // this is an arbitrary amount
		originalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		actualAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "177 Q st",
					City:           "Solomons",
					State:          "MD",
					PostalCode:     "20688",
				},
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					RejectionReason:  models.StringPointer("not applicable"),
					MTOShipmentID:    &successShipment.ID,
					Reason:           models.StringPointer("there was a delay in getting the apartment"),
					SITEntryDate:     &nowDate,
					SITDepartureDate: &later,
					SITPostalCode:    models.StringPointer("90210"),
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model:    successShipment,
				LinkOnly: true,
			},
			{
				Model:    actualAddress,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
				LinkOnly: true,
			},
			{
				Model:    originalAddress,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		// Validate incoming payload: no body to validate

		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Len(movePayload.MtoServiceItems(), 1)

		serviceItemPayload := movePayload.MtoServiceItems()[0]

		json, err := json.Marshal(serviceItemPayload)
		suite.NoError(err)
		payload := primev3messages.MTOServiceItemOriginSIT{}
		err = payload.UnmarshalJSON(json)
		suite.NoError(err)

		suite.Equal(serviceItem.MoveTaskOrderID.String(), payload.MoveTaskOrderID().String())
		suite.Equal(serviceItem.MTOShipmentID.String(), payload.MtoShipmentID().String())
		suite.Equal(serviceItem.ID.String(), payload.ID().String())
		suite.Equal("MTOServiceItemOriginSIT", string(payload.ModelType()))
		suite.Equal(string(serviceItem.ReService.Code), string(*payload.ReServiceCode))
		suite.Equal(serviceItem.ReService.Name, payload.ReServiceName())
		suite.Equal(string(serviceItem.Status), string(payload.Status()))
		suite.Equal(*serviceItem.RejectionReason, *payload.RejectionReason())
		suite.Equal(*serviceItem.Reason, *payload.Reason)
		suite.Equal(serviceItem.SITEntryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.SitEntryDate).Format(time.RFC3339))
		suite.Equal(serviceItem.SITDepartureDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.SitDepartureDate).Format(time.RFC3339))
		suite.Equal(*serviceItem.SITPostalCode, *payload.SitPostalCode)
		verifyAddressFields(serviceItem.SITOriginHHGActualAddress, payload.SitHHGActualOrigin)
		verifyAddressFields(serviceItem.SITOriginHHGOriginalAddress, payload.SitHHGOriginalOrigin)

		suite.NotNil(payload.ETag())
	})

	suite.Run("Success - return all MTOServiceItemDestSIT fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}

		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)

		now := time.Now()
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		later := nowDate.AddDate(0, 0, 3) // this is an arbitrary amount
		finalAddress := factory.BuildAddress(suite.DB(), nil, nil)

		contact1 := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
			MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
				DateOfContact:              time.Date(2023, time.December, 04, 0, 0, 0, 0, time.UTC),
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Date(2023, time.December, 02, 0, 0, 0, 0, time.UTC),
				Type:                       models.CustomerContactTypeFirst,
			},
		})

		contact2 := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{
			MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
				DateOfContact:              time.Date(2023, time.December, 8, 0, 0, 0, 0, time.UTC),
				TimeMilitary:               "1600Z",
				FirstAvailableDeliveryDate: time.Date(2023, time.December, 07, 0, 0, 0, 0, time.UTC),
				Type:                       models.CustomerContactTypeSecond,
			},
		})
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					RejectionReason:      models.StringPointer("not applicable"),
					MTOShipmentID:        &successShipment.ID,
					Reason:               models.StringPointer("there was a delay in getting the apartment"),
					SITEntryDate:         &nowDate,
					SITDepartureDate:     &later,
					CustomerContacts:     models.MTOServiceItemCustomerContacts{contact1, contact2},
					SITCustomerContacted: &nowDate,
					SITRequestedDelivery: &later,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model:    successShipment,
				LinkOnly: true,
			},
			{
				Model:    finalAddress,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		// Validate incoming payload: no body to validate

		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Len(movePayload.MtoServiceItems(), 1)

		serviceItemPayload := movePayload.MtoServiceItems()[0]

		json, err := json.Marshal(serviceItemPayload)
		suite.NoError(err)
		payload := primev3messages.MTOServiceItemDestSIT{}
		err = payload.UnmarshalJSON(json)
		suite.NoError(err)

		suite.Equal(serviceItem.MoveTaskOrderID.String(), payload.MoveTaskOrderID().String())
		suite.Equal(serviceItem.MTOShipmentID.String(), payload.MtoShipmentID().String())
		suite.Equal(serviceItem.ID.String(), payload.ID().String())
		suite.Equal("MTOServiceItemDestSIT", string(payload.ModelType()))
		suite.Equal(string(serviceItem.ReService.Code), string(*payload.ReServiceCode))
		suite.Equal(serviceItem.ReService.Name, payload.ReServiceName())
		suite.Equal(string(serviceItem.Status), string(payload.Status()))
		suite.Equal(*serviceItem.RejectionReason, *payload.RejectionReason())
		suite.Equal(*serviceItem.Reason, *payload.Reason)
		suite.Equal(serviceItem.SITEntryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.SitEntryDate).Format(time.RFC3339))
		suite.Equal(serviceItem.SITDepartureDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.SitDepartureDate).Format(time.RFC3339))
		suite.Equal(serviceItem.SITCustomerContacted.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.SitCustomerContacted).Format(time.RFC3339))
		suite.Equal(serviceItem.SITRequestedDelivery.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.SitRequestedDelivery).Format(time.RFC3339))
		suite.Equal(serviceItem.CustomerContacts[0].DateOfContact.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.DateOfContact1).Format(time.RFC3339))
		suite.Equal(serviceItem.CustomerContacts[0].TimeMilitary, *payload.TimeMilitary1)
		suite.Equal(serviceItem.CustomerContacts[0].FirstAvailableDeliveryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.FirstAvailableDeliveryDate1).Format(time.RFC3339))
		suite.Equal(serviceItem.CustomerContacts[1].DateOfContact.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.DateOfContact2).Format(time.RFC3339))
		suite.Equal(serviceItem.CustomerContacts[1].TimeMilitary, *payload.TimeMilitary2)
		suite.Equal(serviceItem.CustomerContacts[1].FirstAvailableDeliveryDate.Format(time.RFC3339), handlers.FmtDatePtrToPop(payload.FirstAvailableDeliveryDate2).Format(time.RFC3339))
		verifyAddressFields(serviceItem.SITDestinationFinalAddress, payload.SitDestinationFinalAddress)

		suite.NotNil(payload.ETag())
	})

	suite.Run("Success - return all MTOServiceItemShuttle fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}

		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)

		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					RejectionReason: models.StringPointer("not applicable"),
					MTOShipmentID:   &successShipment.ID,
					Reason:          models.StringPointer("this is a special item"),
					EstimatedWeight: models.PoundPointer(400),
					ActualWeight:    models.PoundPointer(500),
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model:    successShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDSHUT,
				},
			},
		}, nil)

		// Validate incoming payload: no body to validate

		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Len(movePayload.MtoServiceItems(), 1)

		serviceItemPayload := movePayload.MtoServiceItems()[0]

		json, err := json.Marshal(serviceItemPayload)
		suite.NoError(err)
		payload := primev3messages.MTOServiceItemShuttle{}
		err = payload.UnmarshalJSON(json)
		suite.NoError(err)

		suite.Equal(serviceItem.MoveTaskOrderID.String(), payload.MoveTaskOrderID().String())
		suite.Equal(serviceItem.MTOShipmentID.String(), payload.MtoShipmentID().String())
		suite.Equal(serviceItem.ID.String(), payload.ID().String())
		suite.Equal("MTOServiceItemShuttle", string(payload.ModelType()))
		suite.Equal(string(serviceItem.ReService.Code), string(*payload.ReServiceCode))
		suite.Equal(serviceItem.ReService.Name, payload.ReServiceName())
		suite.Equal(string(serviceItem.Status), string(payload.Status()))
		suite.Equal(*serviceItem.RejectionReason, *payload.RejectionReason())
		suite.Equal(*serviceItem.Reason, *payload.Reason)
		suite.Equal(*handlers.FmtPoundPtr(serviceItem.EstimatedWeight), *payload.EstimatedWeight)
		suite.Equal(*handlers.FmtPoundPtr(serviceItem.ActualWeight), *payload.ActualWeight)

		suite.NotNil(payload.ETag())
	})

	suite.Run("Success - return all MTOServiceItemDomesticCrating fields assoicated with the getMoveTaskOrder", func() {
		handler := GetMoveTaskOrderHandler{
			suite.HandlerConfig(),
			movetaskorder.NewMoveTaskOrderFetcher(),
		}

		successMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		successShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
		}, nil)

		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					RejectionReason: models.StringPointer("not applicable"),
					MTOShipmentID:   &successShipment.ID,
					Reason:          models.StringPointer("needs extra care"),
					Description:     models.StringPointer("ATV"),
				},
			},
			{
				Model:    successMove,
				LinkOnly: true,
			},
			{
				Model:    successShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDCRT,
				},
			},
		}, nil)

		cratingDimension := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeCrate,
					Length:    12000,
					Height:    12000,
					Width:     12000,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)

		itemDimension := factory.BuildMTOServiceItemDimension(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItemDimension{
					Type:      models.DimensionTypeItem,
					Length:    11000,
					Height:    11000,
					Width:     11000,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)

		serviceItem.Dimensions = []models.MTOServiceItemDimension{cratingDimension, itemDimension}

		// Validate incoming payload: no body to validate

		params := movetaskorderops.GetMoveTaskOrderParams{
			HTTPRequest: request,
			MoveID:      successMove.Locator,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&movetaskorderops.GetMoveTaskOrderOK{}, response)

		moveResponse := response.(*movetaskorderops.GetMoveTaskOrderOK)
		movePayload := moveResponse.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Len(movePayload.MtoServiceItems(), 1)

		serviceItemPayload := movePayload.MtoServiceItems()[0]

		json, err := json.Marshal(serviceItemPayload)
		suite.NoError(err)
		payload := primev3messages.MTOServiceItemDomesticCrating{}
		err = payload.UnmarshalJSON(json)
		suite.NoError(err)

		suite.Equal(serviceItem.MoveTaskOrderID.String(), payload.MoveTaskOrderID().String())
		suite.Equal(serviceItem.MTOShipmentID.String(), payload.MtoShipmentID().String())
		suite.Equal(serviceItem.ID.String(), payload.ID().String())
		suite.Equal("MTOServiceItemDomesticCrating", string(payload.ModelType()))
		suite.Equal(string(serviceItem.ReService.Code), string(*payload.ReServiceCode))
		suite.Equal(serviceItem.ReService.Name, payload.ReServiceName())
		suite.Equal(string(serviceItem.Status), string(payload.Status()))
		suite.Equal(*serviceItem.RejectionReason, *payload.RejectionReason())
		suite.Equal(*serviceItem.Reason, *payload.Reason)
		suite.Equal(*serviceItem.Description, *payload.Description)
		suite.Equal(serviceItem.Dimensions[0].ID.String(), payload.Crate.ID.String())
		suite.Equal(*serviceItem.Dimensions[0].Height.Int32Ptr(), *payload.Crate.Height)
		suite.Equal(*serviceItem.Dimensions[0].Width.Int32Ptr(), *payload.Crate.Width)
		suite.Equal(*serviceItem.Dimensions[0].Length.Int32Ptr(), *payload.Crate.Length)
		suite.Equal(serviceItem.Dimensions[1].ID.String(), payload.Item.ID.String())
		suite.Equal(*serviceItem.Dimensions[1].Height.Int32Ptr(), *payload.Item.Height)
		suite.Equal(*serviceItem.Dimensions[1].Width.Int32Ptr(), *payload.Item.Width)
		suite.Equal(*serviceItem.Dimensions[1].Length.Int32Ptr(), *payload.Item.Length)

		suite.NotNil(payload.ETag())
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
