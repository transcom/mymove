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
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/services/upload"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestListMovesHandler() {
	suite.Run("Test returns updated with no amendments count", func() {
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
		suite.Equal(0, int(*movesList[0].Amendments.Total))
		suite.Equal(0, int(*movesList[0].Amendments.AvailableSince))
	})

	suite.Run("Test returns updated with amendment count", func() {
		now := time.Now()
		lastFetch := now.Add(-time.Second)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// this move should not be returned
		olderMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		// setup Order and Amendment for move
		primeMoves := make([]models.Move, 0)
		primeMoves = append(primeMoves, move)

		for _, pm := range primeMoves {
			document := factory.BuildDocumentLinkServiceMember(suite.DB(), move.Orders.ServiceMember)

			suite.MustSave(&document)
			suite.Nil(document.DeletedAt)
			pm.Orders.UploadedOrders = document
			pm.Orders.UploadedOrdersID = document.ID

			pm.Orders.UploadedAmendedOrders = &document
			pm.Orders.UploadedAmendedOrdersID = &document.ID
			// nolint:gosec //G601
			suite.MustSave(&pm.Orders)
			upload := models.Upload{
				Filename:    "test.pdf",
				Bytes:       1048576,
				ContentType: uploader.FileTypePDF,
				Checksum:    "ImGQ2Ush0bDHsaQthV5BnQ==",
				UploadType:  models.UploadTypeUSER,
			}
			suite.MustSave(&upload)
			userUpload := models.UserUpload{
				DocumentID: &document.ID,
				UploaderID: document.ServiceMember.UserID,
				UploadID:   upload.ID,
				Upload:     upload,
			}
			suite.MustSave(&userUpload)
		}

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
		suite.Equal(1, int(*movesList[0].Amendments.Total))
		suite.Equal(1, int(*movesList[0].Amendments.AvailableSince))
	})
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
		// Handle the possibility that address.Country is nil
		if address.Country != nil && payload.Country != nil {
			suite.Equal(address.Country.Country, *payload.Country)
		}
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
		suite.NotNil(movePayload.ApprovedAt)
		suite.NotEmpty(movePayload.AvailableToPrimeAt) // checks that the date is not 0001-01-01
		suite.NotEmpty(movePayload.ApprovedAt)         // checks that the date is not 0001-01-01
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
		suite.NotNil(movePayload.ApprovedAt)
		suite.NotEmpty(movePayload.ApprovedAt) // checks that the date is not 0001-01-01
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
		suite.NotNil(movePayload.ApprovedAt)
		suite.NotEmpty(movePayload.ApprovedAt)
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
		suite.NotNil(movePayload.ApprovedAt)
		suite.NotEmpty(movePayload.ApprovedAt) // checks that the date is not 0001-01-01
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
		suite.Equal(successMove.ApprovedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.ApprovedAt).Format(time.RFC3339))
		suite.Equal(successMove.PrimeCounselingCompletedAt.Format(time.RFC3339), handlers.FmtDateTimePtrToPop(movePayload.PrimeCounselingCompletedAt).Format(time.RFC3339))
		suite.Equal(*successMove.PPMType, movePayload.PpmType)
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
		suite.Equal(string(*orders.Grade), string(*ordersPayload.Rank))
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
		now := time.Now()
		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					MTOShipmentID:    &successShipment.ID,
					SITDepartureDate: &now,
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
		payload := primemessages.MTOServiceItemBasic{}
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
		payload := primemessages.MTOServiceItemOriginSIT{}
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
					SITRequestedDelivery: &nowDate,
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
		payload := primemessages.MTOServiceItemDestSIT{}
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
		payload := primemessages.MTOServiceItemShuttle{}
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
		payload := primemessages.MTOServiceItemDomesticCrating{}
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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)

		setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
			mockCreator := &mocks.SignedCertificationCreator{}

			mockCreator.On(
				"CreateSignedCertification",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.SignedCertification"),
			).Return(returnValue...)

			return mockCreator
		}

		setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
			mockUpdater := &mocks.SignedCertificationUpdater{}

			mockUpdater.On(
				"UpdateSignedCertification",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.SignedCertification"),
				mock.AnythingOfType("string"),
			).Return(returnValue...)

			return mockUpdater
		}

		ppmEstimator := &mocks.PPMEstimator{}
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		updater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), ppmEstimator)
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
		moveRouter := moverouter.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
		fetcher := fetch.NewFetcher(queryBuilder)
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)

		setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
			mockCreator := &mocks.SignedCertificationCreator{}

			mockCreator.On(
				"CreateSignedCertification",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.SignedCertification"),
			).Return(returnValue...)

			return mockCreator
		}

		setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
			mockUpdater := &mocks.SignedCertificationUpdater{}

			mockUpdater.On(
				"UpdateSignedCertification",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("models.SignedCertification"),
				mock.AnythingOfType("string"),
			).Return(returnValue...)

			return mockUpdater
		}

		ppmEstimator := &mocks.PPMEstimator{}
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())
		updater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, siCreator, moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), ppmEstimator)
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

func (suite *HandlerSuite) TestDownloadMoveOrderHandler() {
	uri := "/moves/%s/documents"
	paramTypeAll := "ALL"
	suite.Run("Successful DownloadMoveOrder - 200", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		move := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil)

		// Hardcode to true to indicate duty location does not provide GOV counseling
		move.Orders.OriginDutyLocation.ProvidesServicesCounseling = false

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// mock to return nil Error
		mockPrimeDownloadMoveUploadPDFGenerator.On("GenerateDownloadMoveUserUploadPDF",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("services.MoveOrderUploadType"),
			mock.AnythingOfType("models.Move"),
			mock.AnythingOfType("bool")).Return(nil, nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}
		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderOK)

		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderOK{}, downloadMoveOrderResponse)
	})

	suite.Run("Successful DownloadMoveOrder - error generating PDF - 500", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		move := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil)

		// Hardcode to true to indicate duty location does not provide GOV counseling
		move.Orders.OriginDutyLocation.ProvidesServicesCounseling = false

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// mock to return nil Error
		mockPrimeDownloadMoveUploadPDFGenerator.On("GenerateDownloadMoveUserUploadPDF",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("services.MoveOrderUploadType"),
			mock.AnythingOfType("models.Move"),
			mock.AnythingOfType("bool")).Return(nil, errors.New("error"))

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}
		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderInternalServerError)

		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderInternalServerError{}, downloadMoveOrderResponse)
	})

	suite.Run("BadRequest DownloadMoveOrder - missing/empty locator - verify 400", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := ""
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}
		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderBadRequest)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderBadRequest{}, downloadMoveOrderResponse)
	})

	suite.Run("Not Found locator DownloadMoveOrder - 404", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		moves := models.Moves{}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig: handlerConfig,
			MoveSearcher:  &mockMoveSearcher,
			OrderFetcher:  &mockOrderFetcher,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 0, nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}
		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderNotFound)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderNotFound{}, downloadMoveOrderResponse)
	})

	suite.Run("DownloadMoveOrder: move requires counseling but origin duty location does have GOV counseling,  Prime counseling is not needed - 422", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}

		move := factory.BuildMove(suite.DB(), nil, nil)
		// Hardcode to MoveStatusNeedsServiceCounseling status
		//move.Status = models.MoveStatusNeedsServiceCounseling
		// Hardcode to TRUE. TRUE whens GOV counseling available and PRIME counseling NOT needed.
		move.Orders.OriginDutyLocation.ProvidesServicesCounseling = true

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig: handlerConfig,
			MoveSearcher:  &mockMoveSearcher,
			OrderFetcher:  &mockOrderFetcher,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}
		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderUnprocessableEntity)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderUnprocessableEntity{}, downloadMoveOrderResponse)
	})

	suite.Run("DownloadMoveOrder: handles internal errors for search move - 500", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig: handlerConfig,
			MoveSearcher:  &mockMoveSearcher,
			OrderFetcher:  &mockOrderFetcher,
		}

		// mock returning error on move search
		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(nil, 0, apperror.NewInternalServerError("mock"))

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}
		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderInternalServerError)

		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderInternalServerError{}, downloadMoveOrderResponse)
	})

	suite.Run("DownloadMoveOrder: service returns unprocessrntity error - 422", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		move := factory.BuildMove(suite.DB(), nil, nil)

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// mock to return nil Errro
		mockPrimeDownloadMoveUploadPDFGenerator.On("GenerateDownloadMoveUserUploadPDF",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("services.MoveOrderUploadType"),
			mock.AnythingOfType("models.Move"),
			mock.AnythingOfType("bool")).Return(nil, apperror.NewUnprocessableEntityError("test"))

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}
		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderUnprocessableEntity)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderUnprocessableEntity{}, downloadMoveOrderResponse)
	})

	suite.Run("DownloadMoveOrder: service returns internal server error - 500", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		move := factory.BuildMove(suite.DB(), nil, nil)

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// mock to return nil Errro
		mockPrimeDownloadMoveUploadPDFGenerator.On("GenerateDownloadMoveUserUploadPDF",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("services.MoveOrderUploadType"),
			mock.AnythingOfType("models.Move"),
			mock.AnythingOfType("bool")).Return(nil, errors.New("test"))

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}

		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderInternalServerError)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderInternalServerError{}, downloadMoveOrderResponse)
	})

	suite.Run("DownloadMoveOrder: ALL - service returns unprocess entity - 422", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		move := factory.BuildMove(suite.DB(), nil, nil)

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// mock to return nil Errro
		mockPrimeDownloadMoveUploadPDFGenerator.On("GenerateDownloadMoveUserUploadPDF",
			mock.AnythingOfType("*appcontext.appContext"),
			services.MoveOrderUploadAll, //Verify ALL enum is used
			mock.AnythingOfType("models.Move"),
			mock.AnythingOfType("bool")).Return(nil, errors.New("test"))

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAll,
		}

		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderInternalServerError)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderInternalServerError{}, downloadMoveOrderResponse)
	})

	suite.Run("DownloadMoveOrder: Orders Only - service returns unprocess entity - 422", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		move := factory.BuildMove(suite.DB(), nil, nil)

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// mock to return nil Errro
		mockPrimeDownloadMoveUploadPDFGenerator.On("GenerateDownloadMoveUserUploadPDF",
			mock.AnythingOfType("*appcontext.appContext"),
			services.MoveOrderUpload, //Verify Order only enum is used
			mock.AnythingOfType("models.Move"),
			mock.AnythingOfType("bool")).Return(nil, errors.New("test"))

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		paramTypeOrders := "ORDERS"
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeOrders,
		}

		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderInternalServerError)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderInternalServerError{}, downloadMoveOrderResponse)
	})

	suite.Run("DownloadMoveOrder: Orders Only - service returns unprocess entity - 422", func() {
		mockMoveSearcher := mocks.MoveSearcher{}
		mockOrderFetcher := mocks.OrderFetcher{}
		mockPrimeDownloadMoveUploadPDFGenerator := mocks.PrimeDownloadMoveUploadPDFGenerator{}

		move := factory.BuildMove(suite.DB(), nil, nil)

		moves := models.Moves{move}

		handlerConfig := suite.HandlerConfig()
		handler := DownloadMoveOrderHandler{
			HandlerConfig:                       handlerConfig,
			MoveSearcher:                        &mockMoveSearcher,
			OrderFetcher:                        &mockOrderFetcher,
			PrimeDownloadMoveUploadPDFGenerator: &mockPrimeDownloadMoveUploadPDFGenerator,
		}

		mockMoveSearcher.On("SearchMoves",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.MatchedBy(func(_ *services.SearchMovesParams) bool {
				return true
			}),
		).Return(moves, 1, nil)

		// mock to return nil Errro
		mockPrimeDownloadMoveUploadPDFGenerator.On("GenerateDownloadMoveUserUploadPDF",
			mock.AnythingOfType("*appcontext.appContext"),
			services.MoveOrderAmendmentUpload, //Verify Amendment only enum is used
			mock.AnythingOfType("models.Move"),
			mock.AnythingOfType("bool")).Return(nil, errors.New("test"))

		// make the request
		requestUser := factory.BuildUser(nil, nil, nil)
		locator := "test"
		request := httptest.NewRequest("GET", fmt.Sprintf(uri, locator), nil)
		request = suite.AuthenticateUserRequest(request, requestUser)
		paramTypeAmendments := "AMENDMENTS"
		params := movetaskorderops.DownloadMoveOrderParams{
			HTTPRequest: request,
			Locator:     locator,
			Type:        &paramTypeAmendments,
		}

		response := handler.Handle(params)
		downloadMoveOrderResponse := response.(*movetaskorderops.DownloadMoveOrderInternalServerError)
		suite.Assertions.IsType(&movetaskorderops.DownloadMoveOrderInternalServerError{}, downloadMoveOrderResponse)
	})
}
