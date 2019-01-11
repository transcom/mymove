package internalapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	shipmentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/shipments"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) verifyAddressFields(expected, actual *internalmessages.Address) {
	suite.T().Helper()
	suite.Equal(expected.StreetAddress1, actual.StreetAddress1, "Street1 did not match")
	suite.Equal(expected.StreetAddress2, actual.StreetAddress2, "Street2 did not match")
	suite.Equal(expected.StreetAddress3, actual.StreetAddress3, "Street3 did not match")
	suite.Equal(expected.City, actual.City, "City did not match")
	suite.Equal(expected.State, actual.State, "State did not match")
	suite.Equal(expected.PostalCode, actual.PostalCode, "PostalCode did not match")
	suite.Equal(expected.Country, actual.Country, "Country did not match")
}

func (suite *HandlerSuite) TestCreateShipmentHandlerAllValues() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
	})
	sm := move.Orders.ServiceMember

	// Make associated lookup table records.
	testdatagen.FetchOrMakeTariff400ngZip3(suite.DB(), testdatagen.Assertions{
		Tariff400ngZip3: models.Tariff400ngZip3{
			Zip3:          "012",
			BasepointCity: "Pittsfield",
			State:         "MA",
			ServiceArea:   "388",
			RateArea:      "US14",
			Region:        "9",
		},
	})

	testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US14",
			DestinationRegion: "9",
			CodeOfService:     "D",
		},
	})

	addressPayload := fakeAddressPayload()
	requestedPickupDate := strfmt.Date(testdatagen.DateInsideNonPeakRateCycle)

	newShipment := internalmessages.Shipment{
		PickupAddress:                addressPayload,
		HasSecondaryPickupAddress:    handlers.FmtBool(true),
		SecondaryPickupAddress:       addressPayload,
		HasDeliveryAddress:           handlers.FmtBool(true),
		DeliveryAddress:              addressPayload,
		HasPartialSitDeliveryAddress: handlers.FmtBool(true),
		PartialSitDeliveryAddress:    addressPayload,
		WeightEstimate:               swag.Int64(4500),
		ProgearWeightEstimate:        swag.Int64(325),
		SpouseProgearWeightEstimate:  swag.Int64(120),
		RequestedPickupDate:          &requestedPickupDate,
	}

	req := httptest.NewRequest("POST", "/moves/move_id/shipment", nil)
	req = suite.AuthenticateRequest(req, sm)

	params := shipmentop.CreateShipmentParams{
		Shipment:    &newShipment,
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: req,
	}

	handler := CreateShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	planner := route.NewTestingPlanner(2000)
	handler.SetPlanner(planner)

	response := handler.Handle(params)

	suite.Assertions.IsType(&shipmentop.CreateShipmentCreated{}, response)
	unwrapped := response.(*shipmentop.CreateShipmentCreated)
	createShipmentPayload := unwrapped.Payload

	suite.Equal(strfmt.UUID(move.ID.String()), createShipmentPayload.MoveID)
	suite.Equal(strfmt.UUID(sm.ID.String()), createShipmentPayload.ServiceMemberID)
	suite.Equal(internalmessages.ShipmentStatusDRAFT, createShipmentPayload.Status)
	suite.Equal(swag.String("D"), createShipmentPayload.CodeOfService)
	suite.Equal(swag.String("dHHG"), createShipmentPayload.Market)
	suite.EqualValues(3, *createShipmentPayload.EstimatedPackDays)
	suite.EqualValues(12, *createShipmentPayload.EstimatedTransitDays)
	suite.verifyAddressFields(addressPayload, createShipmentPayload.PickupAddress)
	suite.Equal(true, *createShipmentPayload.HasSecondaryPickupAddress)
	suite.verifyAddressFields(addressPayload, createShipmentPayload.SecondaryPickupAddress)
	suite.Equal(true, *createShipmentPayload.HasDeliveryAddress)
	suite.verifyAddressFields(addressPayload, createShipmentPayload.DeliveryAddress)
	suite.Equal(true, *createShipmentPayload.HasPartialSitDeliveryAddress)
	suite.verifyAddressFields(addressPayload, createShipmentPayload.PartialSitDeliveryAddress)
	suite.Equal(swag.Int64(4500), createShipmentPayload.WeightEstimate)
	suite.Equal(swag.Int64(325), createShipmentPayload.ProgearWeightEstimate)
	suite.Equal(swag.Int64(120), createShipmentPayload.SpouseProgearWeightEstimate)

	suite.EqualValues(*createShipmentPayload.EstimatedPackDays, 3, "EstimatedPackDays was not updated")
	suite.EqualValues(*createShipmentPayload.EstimatedTransitDays, 12, "EstimatedTransitDays was not updated")
	suite.EqualValues(*createShipmentPayload.RequestedPickupDate, requestedPickupDate, "RequestedPickupDate was not updated")

	expectedOriginalDeliveryDate := time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)
	suite.EqualValues(time.Time(*createShipmentPayload.OriginalDeliveryDate), expectedOriginalDeliveryDate, "OriginalDeliveryDate was not updated")

	expectedOriginalPackDate := time.Date(testdatagen.TestYear, time.September, 27, 0, 0, 0, 0, time.UTC)
	suite.EqualValues(time.Time(*createShipmentPayload.OriginalPackDate), expectedOriginalPackDate, "OriginalPackDate was not updated")

	count, err := suite.DB().Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.Nil(err, "could not count shipments")
	suite.Equal(1, count)
}

func (suite *HandlerSuite) TestCreateShipmentHandlerEmpty() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	sm := move.Orders.ServiceMember

	req := httptest.NewRequest("POST", "/moves/move_id/shipment", nil)
	req = suite.AuthenticateRequest(req, sm)

	newShipment := internalmessages.Shipment{}
	params := shipmentop.CreateShipmentParams{
		Shipment:    &newShipment,
		MoveID:      strfmt.UUID(move.ID.String()),
		HTTPRequest: req,
	}

	handler := CreateShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&shipmentop.CreateShipmentCreated{}, response)
	unwrapped := response.(*shipmentop.CreateShipmentCreated)

	count, err := suite.DB().Where("move_id=$1", move.ID).Count(&models.Shipment{})
	suite.Nil(err, "could not count shipments")
	suite.Equal(1, count)

	suite.Equal(strfmt.UUID(move.ID.String()), unwrapped.Payload.MoveID)
	suite.Equal(strfmt.UUID(sm.ID.String()), unwrapped.Payload.ServiceMemberID)
	suite.Equal(internalmessages.ShipmentStatusDRAFT, unwrapped.Payload.Status)
	suite.Equal(swag.String("dHHG"), unwrapped.Payload.Market)
	suite.Nil(unwrapped.Payload.CodeOfService) // Won't be able to assign a TDL since we do not have a pickup address.
	suite.Nil(unwrapped.Payload.EstimatedPackDays)
	suite.Nil(unwrapped.Payload.EstimatedTransitDays)
	suite.Nil(unwrapped.Payload.ActualPackDate)
	suite.Nil(unwrapped.Payload.ActualPickupDate)
	suite.Nil(unwrapped.Payload.ActualDeliveryDate)
	suite.Equal(internalmessages.ShipmentMoveDatesSummary{}, *unwrapped.Payload.MoveDatesSummary)
	suite.Nil(unwrapped.Payload.PickupAddress)
	suite.Equal(false, *unwrapped.Payload.HasSecondaryPickupAddress)
	suite.Nil(unwrapped.Payload.SecondaryPickupAddress)
	suite.Equal(false, *unwrapped.Payload.HasDeliveryAddress)
	suite.Nil(unwrapped.Payload.DeliveryAddress)
	suite.Equal(false, *unwrapped.Payload.HasPartialSitDeliveryAddress)
	suite.Nil(unwrapped.Payload.PartialSitDeliveryAddress)
	suite.Nil(unwrapped.Payload.WeightEstimate)
	suite.Nil(unwrapped.Payload.ProgearWeightEstimate)
	suite.Nil(unwrapped.Payload.SpouseProgearWeightEstimate)
}

func (suite *HandlerSuite) TestPatchShipmentsHandlerHappyPath() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	sm := move.Orders.ServiceMember

	addressPayload := testdatagen.MakeDefaultAddress(suite.DB())

	shipment1 := models.Shipment{
		MoveID:                       move.ID,
		Status:                       "DRAFT",
		EstimatedPackDays:            swag.Int64(2),
		EstimatedTransitDays:         swag.Int64(5),
		PickupAddress:                &addressPayload,
		HasSecondaryPickupAddress:    true,
		SecondaryPickupAddress:       &addressPayload,
		HasDeliveryAddress:           false,
		HasPartialSITDeliveryAddress: true,
		PartialSITDeliveryAddress:    &addressPayload,
		WeightEstimate:               handlers.PoundPtrFromInt64Ptr(swag.Int64(4500)),
		ProgearWeightEstimate:        handlers.PoundPtrFromInt64Ptr(swag.Int64(325)),
		SpouseProgearWeightEstimate:  handlers.PoundPtrFromInt64Ptr(swag.Int64(120)),
		ServiceMemberID:              sm.ID,
	}
	suite.MustSave(&shipment1)

	req := httptest.NewRequest("POST", "/moves/move_id/shipment/shipment_id", nil)
	req = suite.AuthenticateRequest(req, sm)

	// Make associated lookup table records.
	testdatagen.FetchOrMakeTariff400ngZip3(suite.DB(), testdatagen.Assertions{
		Tariff400ngZip3: models.Tariff400ngZip3{
			Zip3:          "321",
			BasepointCity: "Crescent City",
			State:         "FL",
			ServiceArea:   "184",
			RateArea:      "ZIP",
			Region:        "13",
		},
	})

	newAddress := otherFakeAddressPayload()

	payload := internalmessages.Shipment{
		HasSecondaryPickupAddress:   handlers.FmtBool(false),
		HasDeliveryAddress:          handlers.FmtBool(true),
		DeliveryAddress:             newAddress,
		SpouseProgearWeightEstimate: swag.Int64(100),
	}

	patchShipmentParams := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment1.ID.String()),
		Shipment:    &payload,
	}

	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(patchShipmentParams)

	// assert we got back the 200 response
	okResponse := response.(*shipmentop.PatchShipmentOK)
	patchShipmentPayload := okResponse.Payload

	suite.Equal(*patchShipmentPayload.HasDeliveryAddress, true, "HasDeliveryAddress should have been updated.")
	suite.verifyAddressFields(newAddress, patchShipmentPayload.DeliveryAddress)

	suite.Equal(*patchShipmentPayload.HasSecondaryPickupAddress, false, "HasSecondaryPickupAddress should have been updated.")
	suite.Nil(patchShipmentPayload.SecondaryPickupAddress, "SecondaryPickupAddress should have been updated to nil.")

	suite.Equal(*patchShipmentPayload.SpouseProgearWeightEstimate, int64(100), "SpouseProgearWeightEstimate should have been set to 100")
}

func (suite *HandlerSuite) TestSetShipmentDates() {
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
	})
	sm := move.Orders.ServiceMember
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			Move:   move,
			MoveID: move.ID,
		},
	})

	req := httptest.NewRequest("POST", "/moves/move_id/shipment/shipment_id", nil)
	req = suite.AuthenticateRequest(req, sm)

	requestedPickupDate := strfmt.Date(testdatagen.DateInsideNonPeakRateCycle)
	payload := internalmessages.Shipment{
		RequestedPickupDate: &requestedPickupDate,
	}

	patchShipmentParams := shipmentop.PatchShipmentParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
		Shipment:    &payload,
	}

	planner := route.NewTestingPlanner(2000)
	handler := PatchShipmentHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	handler.SetPlanner(planner)

	response := handler.Handle(patchShipmentParams)
	okResponse := response.(*shipmentop.PatchShipmentOK)
	patchShipmentPayload := okResponse.Payload

	suite.EqualValues(*patchShipmentPayload.EstimatedPackDays, 3, "EstimatedPackDays was not updated")
	suite.EqualValues(*patchShipmentPayload.EstimatedTransitDays, 12, "EstimatedTransitDays was not updated")
	suite.EqualValues(*patchShipmentPayload.RequestedPickupDate, requestedPickupDate, "RequestedPickupDate was not updated")

	expectedOriginalDeliveryDate := time.Date(testdatagen.TestYear, time.October, 15, 0, 0, 0, 0, time.UTC)
	suite.EqualValues(time.Time(*patchShipmentPayload.OriginalDeliveryDate), expectedOriginalDeliveryDate, "OriginalDeliveryDate was not updated")

	expectedOriginalPackDate := time.Date(testdatagen.TestYear, time.September, 27, 0, 0, 0, 0, time.UTC)
	suite.EqualValues(time.Time(*patchShipmentPayload.OriginalPackDate), expectedOriginalPackDate, "OriginalPackDate was not updated")
}

func (suite *HandlerSuite) TestApproveHHGHandler() {
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	shipmentAssertions := testdatagen.Assertions{
		Shipment: models.Shipment{
			Status: "ACCEPTED",
		},
	}
	shipment := testdatagen.MakeShipment(suite.DB(), shipmentAssertions)
	suite.MustSave(&shipment)

	handler := ApproveHHGHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	path := "/shipments/shipment_id/approve"
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := shipmentop.ApproveHHGParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	// assert we got back the 200 response
	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.ApproveHHGOK{}, response)
	okResponse := response.(*shipmentop.ApproveHHGOK)
	suite.Equal("APPROVED", string(okResponse.Payload.Status))
}

func (suite *HandlerSuite) TestCompleteHHGHandler() {
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	shipmentAssertions := testdatagen.Assertions{
		Shipment: models.Shipment{
			Status: "DELIVERED",
		},
	}
	shipment := testdatagen.MakeShipment(suite.DB(), shipmentAssertions)
	suite.MustSave(&shipment)

	handler := CompleteHHGHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	path := "/shipments/shipment_id/complete"
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := shipmentop.CompleteHHGParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	// assert we got back the 200 response
	response := handler.Handle(params)
	suite.Assertions.IsType(&shipmentop.CompleteHHGOK{}, response)
	okResponse := response.(*shipmentop.CompleteHHGOK)
	suite.Equal("COMPLETED", string(okResponse.Payload.Status))
}

/*
func (suite *HandlerSuite) TestShipmentInvoiceHandler() {
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{})
	suite.MustSave(&shipment)

	handler := ShipmentInvoiceHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	path := "/shipments/shipment_id/invoice"
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := shipmentop.SendHHGInvoiceParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	// assert we got back the OK response
	response := handler.Handle(params)
	suite.Equal(shipmentop.NewSendHHGInvoiceOK(), response)

    // check that invoices were saved and have submitted status
	var invoices []models.Invoice
		suite.NoError(suite.DB().All(&invoices)) // needs to filter on ID that was saved
		suite.NotEmpty(invoices)
		for _, invoice := range invoices {
		suite.Equal(models.InvoiceStatusSUBMITTED, invoice.Status)
	}
}
*/

func (suite *HandlerSuite) TestShipmentInvoiceHandlerShipmentWrongState() {
	// Given: an office User
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			Status:                       "DRAFT",
			EstimatedPackDays:            swag.Int64(2),
			EstimatedTransitDays:         swag.Int64(5),
			HasSecondaryPickupAddress:    true,
			HasDeliveryAddress:           false,
			HasPartialSITDeliveryAddress: true,
			WeightEstimate:               handlers.PoundPtrFromInt64Ptr(swag.Int64(4500)),
			ProgearWeightEstimate:        handlers.PoundPtrFromInt64Ptr(swag.Int64(325)),
			SpouseProgearWeightEstimate:  handlers.PoundPtrFromInt64Ptr(swag.Int64(120)),
		},
	})
	shipmentOffer := testdatagen.MakeShipmentOffer(suite.DB(), testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			ShipmentID: shipment.ID,
		},
	})
	suite.MustSave(&shipment)
	suite.MustSave(&shipmentOffer)

	handler := ShipmentInvoiceHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}

	path := "/shipments/shipment_id/invoice"
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := shipmentop.CreateAndSendHHGInvoiceParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	// assert we got back the conflict response
	response := handler.Handle(params)
	suite.Equal(shipmentop.NewCreateAndSendHHGInvoiceConflict(), response)
}
