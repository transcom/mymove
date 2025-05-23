package payloads

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestMoveTaskOrder() {
	moveTaskOrderID, _ := uuid.NewV4()
	ordersID, _ := uuid.NewV4()
	referenceID := "testID"
	primeTime := time.Now()
	primeAcknowledgedAt := time.Now().AddDate(0, 0, -3)
	submittedAt := time.Now()
	excessWeightQualifiedAt := time.Now()
	excessUnaccompaniedBaggageWeightQualifiedAt := time.Now()
	excessWeightAcknowledgedAt := time.Now()
	excessUnaccompaniedBaggageWeightAcknowledgedAt := time.Now()
	excessWeightUploadID := uuid.Must(uuid.NewV4())
	ordersType := primev3messages.OrdersTypeRETIREMENT
	originDutyGBLOC := "KKFA"
	shipmentGBLOC := "AGFM"
	packingInstructions := models.InstructionsBeforeContractNumber + factory.DefaultContractNumber + models.InstructionsAfterContractNumber

	streetAddress2 := "Apt 1"
	streetAddress3 := "Apt 1"

	backupContacts := models.BackupContacts{}
	backupContacts = append(backupContacts, models.BackupContact{
		Name:  "Backup contact name",
		Phone: "555-555-5555",
		Email: "backup@backup.com",
	})
	serviceMember := models.ServiceMember{
		BackupContacts: backupContacts,
	}

	basicMove := models.Move{
		ID:                 moveTaskOrderID,
		Locator:            "TESTTEST",
		CreatedAt:          time.Now(),
		AvailableToPrimeAt: &primeTime,
		OrdersID:           ordersID,
		Orders: models.Order{
			OrdersType:                     internalmessages.OrdersType(ordersType),
			OriginDutyLocationGBLOC:        &originDutyGBLOC,
			SupplyAndServicesCostEstimate:  models.SupplyAndServicesCostEstimate,
			MethodOfPayment:                models.MethodOfPayment,
			NAICS:                          models.NAICS,
			PackingAndShippingInstructions: packingInstructions,
			ServiceMember:                  serviceMember,
		},
		ReferenceID:          &referenceID,
		PaymentRequests:      models.PaymentRequests{},
		SubmittedAt:          &submittedAt,
		UpdatedAt:            time.Now(),
		Status:               models.MoveStatusAPPROVED,
		SignedCertifications: models.SignedCertifications{},
		MTOServiceItems:      models.MTOServiceItems{},
		MTOShipments: models.MTOShipments{
			models.MTOShipment{
				PickupAddress: &models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: &streetAddress2,
					StreetAddress3: &streetAddress3,
					City:           "Washington",
					State:          "DC",
					PostalCode:     "20001",
					County:         models.StringPointer("my county"),
				},
			},
		},
		ExcessWeightQualifiedAt:                        &excessWeightQualifiedAt,
		ExcessWeightAcknowledgedAt:                     &excessWeightAcknowledgedAt,
		ExcessUnaccompaniedBaggageWeightQualifiedAt:    &excessUnaccompaniedBaggageWeightQualifiedAt,
		ExcessUnaccompaniedBaggageWeightAcknowledgedAt: &excessUnaccompaniedBaggageWeightAcknowledgedAt,
		ExcessWeightUploadID:                           &excessWeightUploadID,
		Contractor: &models.Contractor{
			ContractNumber: factory.DefaultContractNumber,
		},
		ShipmentGBLOC: models.MoveToGBLOCs{
			models.MoveToGBLOC{GBLOC: &shipmentGBLOC},
		},
		PrimeAcknowledgedAt: &primeAcknowledgedAt,
	}

	suite.Run("Success - Returns a basic move payload with no payment requests, service items or shipments", func() {
		returnedModel := MoveTaskOrder(suite.AppContextForTest(), &basicMove)

		suite.IsType(&primev3messages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(basicMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(ordersType, returnedModel.Order.OrdersType)
		suite.Equal(shipmentGBLOC, returnedModel.Order.OriginDutyLocationGBLOC)
		suite.Equal(referenceID, returnedModel.ReferenceID)
		suite.Equal(strfmt.DateTime(basicMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotEmpty(returnedModel.ETag)
		suite.True(returnedModel.ExcessUnaccompaniedBaggageWeightQualifiedAt.Equal(strfmt.DateTime(*basicMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)))
		suite.True(returnedModel.ExcessWeightQualifiedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightQualifiedAt)))
		suite.True(returnedModel.ExcessUnaccompaniedBaggageWeightAcknowledgedAt.Equal(strfmt.DateTime(*basicMove.ExcessUnaccompaniedBaggageWeightAcknowledgedAt)))
		suite.True(returnedModel.ExcessWeightAcknowledgedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightAcknowledgedAt)))
		suite.Require().NotNil(returnedModel.ExcessWeightUploadID)
		suite.Equal(strfmt.UUID(basicMove.ExcessWeightUploadID.String()), *returnedModel.ExcessWeightUploadID)
		suite.Equal(factory.DefaultContractNumber, returnedModel.ContractNumber)
		suite.Equal(models.SupplyAndServicesCostEstimate, returnedModel.Order.SupplyAndServicesCostEstimate)
		suite.Equal(models.MethodOfPayment, returnedModel.Order.MethodOfPayment)
		suite.Equal(models.NAICS, returnedModel.Order.Naics)
		suite.Equal(packingInstructions, returnedModel.Order.PackingAndShippingInstructions)
		suite.Require().NotEmpty(returnedModel.MtoShipments)
		suite.Equal(basicMove.MTOShipments[0].PickupAddress.County, returnedModel.MtoShipments[0].PickupAddress.County)
		suite.Equal(basicMove.Orders.ServiceMember.BackupContacts[0].Name, returnedModel.Order.Customer.BackupContact.Name)
		suite.Equal(basicMove.Orders.ServiceMember.BackupContacts[0].Phone, returnedModel.Order.Customer.BackupContact.Phone)
		suite.Equal(basicMove.Orders.ServiceMember.BackupContacts[0].Email, returnedModel.Order.Customer.BackupContact.Email)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.PrimeAcknowledgedAt), returnedModel.PrimeAcknowledgedAt)
	})

	suite.Run("Success - payload with RateArea", func() {
		cloneMove := func(orig *models.Move) (*models.Move, error) {
			origJSON, err := json.Marshal(orig)
			if err != nil {
				return nil, err
			}

			clone := models.Move{}
			if err = json.Unmarshal(origJSON, &clone); err != nil {
				return nil, err
			}

			return &clone, nil
		}

		newMove, err := cloneMove(&basicMove)
		suite.NotNil(newMove)
		suite.Nil(err)

		const fairbanksAlaskaPostalCode = "99716"
		const anchorageAlaskaPostalCode = "99521"
		const wasillaAlaskaPostalCode = "99652"
		const beverlyHillsCAPostalCode = "90210"

		//clear MTOShipment and rebuild with specifics for test
		newMove.MTOShipments = newMove.MTOShipments[:0]

		newMove.MTOShipments = append(newMove.MTOShipments, models.MTOShipment{
			MarketCode: models.MarketCodeInternational,
			PickupAddress: &models.Address{
				StreetAddress1: "123 Main St",
				StreetAddress2: &streetAddress2,
				StreetAddress3: &streetAddress3,
				City:           "Fairbanks",
				State:          "AK",
				PostalCode:     fairbanksAlaskaPostalCode,
			},
			DestinationAddress: &models.Address{
				StreetAddress1:   "123 Main St",
				StreetAddress2:   &streetAddress2,
				StreetAddress3:   &streetAddress3,
				City:             "Anchorage",
				State:            "AK",
				PostalCode:       anchorageAlaskaPostalCode,
				DestinationGbloc: models.StringPointer("JEAT"),
			},
		})
		newMove.MTOShipments = append(newMove.MTOShipments, models.MTOShipment{
			MarketCode: models.MarketCodeInternational,
			PickupAddress: &models.Address{
				StreetAddress1: "123 Main St",
				StreetAddress2: &streetAddress2,
				StreetAddress3: &streetAddress3,
				City:           "Wasilla",
				State:          "AK",
				PostalCode:     wasillaAlaskaPostalCode,
			},
			DestinationAddress: &models.Address{
				StreetAddress1:   "123 Main St",
				StreetAddress2:   &streetAddress2,
				StreetAddress3:   &streetAddress3,
				City:             "Wasilla",
				State:            "AK",
				PostalCode:       wasillaAlaskaPostalCode,
				DestinationGbloc: models.StringPointer("JEAT"),
			},
		})
		newMove.MTOShipments = append(newMove.MTOShipments, models.MTOShipment{
			ShipmentType: models.MTOShipmentTypePPM,
			MarketCode:   models.MarketCodeInternational,
			PPMShipment: &models.PPMShipment{
				ID:                    uuid.Must(uuid.NewV4()),
				ApprovedAt:            models.TimePointer(time.Now()),
				Status:                models.PPMShipmentStatusNeedsAdvanceApproval,
				ActualMoveDate:        models.TimePointer(time.Now()),
				HasReceivedAdvance:    models.BoolPointer(true),
				AdvanceAmountReceived: models.CentPointer(unit.Cents(340000)),
				FinalIncentive:        models.CentPointer(50000000),
				PickupAddress: &models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: &streetAddress2,
					StreetAddress3: &streetAddress3,
					City:           "Wasilla",
					State:          "AK",
					PostalCode:     wasillaAlaskaPostalCode,
				},
				DestinationAddress: &models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: &streetAddress2,
					StreetAddress3: &streetAddress3,
					City:           "Wasilla",
					State:          "AK",
					PostalCode:     wasillaAlaskaPostalCode,
				},
			},
		})
		newMove.MTOShipments = append(newMove.MTOShipments, models.MTOShipment{
			ShipmentType: models.MTOShipmentTypePPM,
			MarketCode:   models.MarketCodeDomestic,
			PPMShipment: &models.PPMShipment{
				ID:                    uuid.Must(uuid.NewV4()),
				ApprovedAt:            models.TimePointer(time.Now()),
				Status:                models.PPMShipmentStatusNeedsAdvanceApproval,
				ActualMoveDate:        models.TimePointer(time.Now()),
				HasReceivedAdvance:    models.BoolPointer(true),
				AdvanceAmountReceived: models.CentPointer(unit.Cents(340000)),
				FinalIncentive:        models.CentPointer(50000000),
				PickupAddress: &models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: &streetAddress2,
					StreetAddress3: &streetAddress3,
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     beverlyHillsCAPostalCode,
				},
				DestinationAddress: &models.Address{
					StreetAddress1: "123 Main St",
					StreetAddress2: &streetAddress2,
					StreetAddress3: &streetAddress3,
					City:           "Beverly Hills",
					State:          "CA",
					PostalCode:     beverlyHillsCAPostalCode,
				},
			},
		})
		newMove.MTOShipments = append(newMove.MTOShipments, models.MTOShipment{
			MarketCode: models.MarketCodeDomestic,
			PickupAddress: &models.Address{
				StreetAddress1:   "123 Main St",
				StreetAddress2:   &streetAddress2,
				StreetAddress3:   &streetAddress3,
				City:             "Beverly Hills",
				State:            "CA",
				PostalCode:       beverlyHillsCAPostalCode,
				DestinationGbloc: models.StringPointer("JEAT"),
			},
			DestinationAddress: &models.Address{
				StreetAddress1:   "123 Main St",
				StreetAddress2:   &streetAddress2,
				StreetAddress3:   &streetAddress3,
				City:             "Beverly Hills",
				State:            "CA",
				PostalCode:       beverlyHillsCAPostalCode,
				DestinationGbloc: models.StringPointer("JEAT"),
			},
		})

		// no ShipmentPostalCodeRateArea passed in
		returnedModel := MoveTaskOrderWithShipmentRateAreas(suite.AppContextForTest(), newMove, nil)

		suite.IsType(&primev3messages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(newMove.ID.String()), returnedModel.ID)
		suite.Equal(newMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(newMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(newMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(newMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(ordersType, returnedModel.Order.OrdersType)
		suite.Equal(shipmentGBLOC, returnedModel.Order.OriginDutyLocationGBLOC)
		suite.Equal(referenceID, returnedModel.ReferenceID)
		suite.Equal(strfmt.DateTime(newMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotEmpty(returnedModel.ETag)
		suite.True(returnedModel.ExcessWeightQualifiedAt.Equal(strfmt.DateTime(*newMove.ExcessWeightQualifiedAt)))
		suite.True(returnedModel.ExcessWeightAcknowledgedAt.Equal(strfmt.DateTime(*newMove.ExcessWeightAcknowledgedAt)))
		suite.Require().NotNil(returnedModel.ExcessWeightUploadID)
		suite.Equal(strfmt.UUID(newMove.ExcessWeightUploadID.String()), *returnedModel.ExcessWeightUploadID)
		suite.Equal(factory.DefaultContractNumber, returnedModel.ContractNumber)
		suite.Equal(models.SupplyAndServicesCostEstimate, returnedModel.Order.SupplyAndServicesCostEstimate)
		suite.Equal(models.MethodOfPayment, returnedModel.Order.MethodOfPayment)
		suite.Equal(models.NAICS, returnedModel.Order.Naics)
		suite.Equal(packingInstructions, returnedModel.Order.PackingAndShippingInstructions)
		suite.Require().NotEmpty(returnedModel.MtoShipments)
		suite.Equal(newMove.MTOShipments[0].PickupAddress.County, returnedModel.MtoShipments[0].PickupAddress.County)

		// verify there are no RateArea set because no ShipmentPostalCodeRateArea passed in.
		for _, shipment := range returnedModel.MtoShipments {
			suite.Nil(shipment.OriginRateArea)
			suite.Nil(shipment.DestinationRateArea)
			if shipment.PpmShipment != nil {
				suite.Nil(shipment.PpmShipment.OriginRateArea)
				suite.Nil(shipment.PpmShipment.DestinationRateArea)
			}
		}

		// mock up ShipmentPostalCodeRateArea
		shipmentPostalCodeRateArea := []services.ShipmentPostalCodeRateArea{
			{
				PostalCode: fairbanksAlaskaPostalCode,
				RateArea: &models.ReRateArea{
					ID:   uuid.Must(uuid.NewV4()),
					Code: fairbanksAlaskaPostalCode,
					Name: fairbanksAlaskaPostalCode,
				},
			},
			{
				PostalCode: anchorageAlaskaPostalCode,
				RateArea: &models.ReRateArea{
					ID:   uuid.Must(uuid.NewV4()),
					Code: anchorageAlaskaPostalCode,
					Name: anchorageAlaskaPostalCode,
				},
			},
			{
				PostalCode: wasillaAlaskaPostalCode,
				RateArea: &models.ReRateArea{
					ID:   uuid.Must(uuid.NewV4()),
					Code: wasillaAlaskaPostalCode,
					Name: wasillaAlaskaPostalCode,
				},
			},
			{
				PostalCode: beverlyHillsCAPostalCode,
				RateArea: &models.ReRateArea{
					ID:   uuid.Must(uuid.NewV4()),
					Code: beverlyHillsCAPostalCode,
					Name: beverlyHillsCAPostalCode,
				},
			},
		}

		returnedModel = MoveTaskOrderWithShipmentRateAreas(suite.AppContextForTest(), newMove, &shipmentPostalCodeRateArea)

		var shipmentPostalCodeRateAreaLookupMap = make(map[string]services.ShipmentPostalCodeRateArea)
		for _, i := range shipmentPostalCodeRateArea {
			shipmentPostalCodeRateAreaLookupMap[i.PostalCode] = i
		}

		// test Alaska/Oconus PostCodes have associative RateArea for respective shipment
		expectedAlaskaPostalCodes := []string{fairbanksAlaskaPostalCode, anchorageAlaskaPostalCode, wasillaAlaskaPostalCode}
		for _, shipment := range returnedModel.MtoShipments {
			if shipment.PpmShipment != nil {
				suite.NotNil(shipment.PpmShipment.PickupAddress)
				suite.NotNil(shipment.PpmShipment.DestinationAddress)
				if slices.Contains(expectedAlaskaPostalCodes, *shipment.PpmShipment.PickupAddress.PostalCode) {
					// verify mapping of RateArea is correct
					ra, contains := shipmentPostalCodeRateAreaLookupMap[*shipment.PpmShipment.PickupAddress.PostalCode]
					suite.True(contains)
					suite.NotNil(shipment.PpmShipment.OriginRateArea)
					// for testing purposes RateArea code/names are using postalCodes as value
					suite.Equal(ra.PostalCode, *shipment.PpmShipment.PickupAddress.PostalCode)
					suite.Equal(ra.PostalCode, *shipment.PpmShipment.OriginRateArea.RateAreaName)
				} else {
					suite.Nil(shipment.PpmShipment.OriginRateArea)
				}
				if slices.Contains(expectedAlaskaPostalCodes, *shipment.PpmShipment.DestinationAddress.PostalCode) {
					ra, contains := shipmentPostalCodeRateAreaLookupMap[*shipment.PpmShipment.DestinationAddress.PostalCode]
					suite.True(contains)
					suite.NotNil(shipment.PpmShipment.DestinationRateArea)
					suite.Equal(ra.PostalCode, *shipment.PpmShipment.DestinationAddress.PostalCode)
					suite.Equal(ra.PostalCode, *shipment.PpmShipment.DestinationRateArea.RateAreaName)
				} else {
					suite.Nil(shipment.PpmShipment.DestinationRateArea)
				}
				// because it's PPM verify root doesnt have rateArea for org/dest
				suite.Nil(shipment.OriginRateArea)
				suite.Nil(shipment.DestinationRateArea)
			} else {
				suite.NotNil(shipment.PickupAddress)
				suite.NotNil(shipment.DestinationAddress)
				suite.NotNil(shipment.DestinationAddress.DestinationGbloc)
				if slices.Contains(expectedAlaskaPostalCodes, *shipment.PickupAddress.PostalCode) {
					ra, contains := shipmentPostalCodeRateAreaLookupMap[*shipment.PickupAddress.PostalCode]
					suite.True(contains)
					suite.NotNil(shipment.OriginRateArea)
					suite.Equal(ra.PostalCode, *shipment.PickupAddress.PostalCode)
					suite.Equal(ra.PostalCode, *shipment.OriginRateArea.RateAreaName)
					suite.NotNil(shipment.OriginRateArea)
				} else {
					suite.Nil(shipment.OriginRateArea)
				}
				if slices.Contains(expectedAlaskaPostalCodes, *shipment.DestinationAddress.PostalCode) {
					ra, contains := shipmentPostalCodeRateAreaLookupMap[*shipment.DestinationAddress.PostalCode]
					suite.True(contains)
					suite.NotNil(shipment.DestinationRateArea)
					suite.Equal(ra.PostalCode, *shipment.DestinationAddress.PostalCode)
					suite.Equal(ra.PostalCode, *shipment.DestinationRateArea.RateAreaName)
					suite.NotNil(shipment.OriginRateArea)
				} else {
					suite.Nil(shipment.OriginRateArea)
				}
			}
		}
	})
}

func (suite *PayloadsSuite) TestReweigh() {
	id, _ := uuid.NewV4()
	shipmentID, _ := uuid.NewV4()
	requestedAt := time.Now()
	createdAt := time.Now()
	updatedAt := time.Now()

	reweigh := models.Reweigh{
		ID:          id,
		ShipmentID:  shipmentID,
		RequestedAt: requestedAt,
		RequestedBy: models.ReweighRequesterTOO,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	suite.Run("Success - Returns a reweigh payload without optional fields", func() {
		returnedPayload := Reweigh(&reweigh)

		suite.IsType(&primev3messages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(reweigh.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(reweigh.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primev3messages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
		suite.Equal(strfmt.DateTime(reweigh.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(reweigh.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Nil(returnedPayload.Weight)
		suite.Nil(returnedPayload.VerificationReason)
		suite.Nil(returnedPayload.VerificationProvidedAt)
		suite.NotEmpty(returnedPayload.ETag)

	})

	suite.Run("Success - Returns a reweigh payload with optional fields", func() {
		// Set optional fields
		weight := int64(2000)
		reweigh.Weight = handlers.PoundPtrFromInt64Ptr(&weight)

		verificationProvidedAt := time.Now()
		reweigh.VerificationProvidedAt = &verificationProvidedAt

		verificationReason := "Because I said so"
		reweigh.VerificationReason = &verificationReason

		// Send model through func
		returnedPayload := Reweigh(&reweigh)

		suite.IsType(&primev3messages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(reweigh.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(reweigh.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primev3messages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
		suite.Equal(strfmt.DateTime(reweigh.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(reweigh.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Equal(handlers.FmtPoundPtr(reweigh.Weight), returnedPayload.Weight)
		suite.Equal(handlers.FmtStringPtr(reweigh.VerificationReason), returnedPayload.VerificationReason)
		suite.Equal(handlers.FmtDateTimePtr(reweigh.VerificationProvidedAt), returnedPayload.VerificationProvidedAt)
		suite.NotEmpty(returnedPayload.ETag)
	})
}

func (suite *PayloadsSuite) TestSitExtension() {

	id, _ := uuid.NewV4()
	shipmentID, _ := uuid.NewV4()
	createdAt := time.Now()
	updatedAt := time.Now()

	sitDurationUpdate := models.SITDurationUpdate{
		ID:            id,
		MTOShipmentID: shipmentID,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		RequestedDays: int(30),
		RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		Status:        models.SITExtensionStatusPending,
	}

	suite.Run("Success - Returns a sitextension payload without optional fields", func() {
		returnedPayload := SITDurationUpdate(&sitDurationUpdate)

		suite.IsType(&primev3messages.SITExtension{}, returnedPayload)
		suite.Equal(strfmt.UUID(sitDurationUpdate.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(sitDurationUpdate.MTOShipmentID.String()), returnedPayload.MtoShipmentID)
		suite.Equal(strfmt.DateTime(sitDurationUpdate.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(sitDurationUpdate.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Nil(returnedPayload.ApprovedDays)
		suite.Nil(returnedPayload.ContractorRemarks)
		suite.Nil(returnedPayload.OfficeRemarks)
		suite.Nil(returnedPayload.DecisionDate)
		suite.NotNil(returnedPayload.ETag)

	})

	suite.Run("Success - Returns a sit duration update payload with optional fields", func() {
		// Set optional fields
		approvedDays := int(30)
		sitDurationUpdate.ApprovedDays = &approvedDays

		contractorRemarks := "some reason"
		sitDurationUpdate.ContractorRemarks = &contractorRemarks

		officeRemarks := "some other reason"
		sitDurationUpdate.OfficeRemarks = &officeRemarks

		decisionDate := time.Now()
		sitDurationUpdate.DecisionDate = &decisionDate

		// Send model through func
		returnedPayload := SITDurationUpdate(&sitDurationUpdate)

		suite.IsType(&primev3messages.SITExtension{}, returnedPayload)
		suite.Equal(strfmt.UUID(sitDurationUpdate.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(sitDurationUpdate.MTOShipmentID.String()), returnedPayload.MtoShipmentID)
		suite.Equal(strfmt.DateTime(sitDurationUpdate.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(sitDurationUpdate.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Equal(handlers.FmtIntPtrToInt64(sitDurationUpdate.ApprovedDays), returnedPayload.ApprovedDays)
		suite.Equal(sitDurationUpdate.ContractorRemarks, returnedPayload.ContractorRemarks)
		suite.Equal(sitDurationUpdate.OfficeRemarks, returnedPayload.OfficeRemarks)
		suite.Equal((*strfmt.DateTime)(sitDurationUpdate.DecisionDate), returnedPayload.DecisionDate)
		suite.NotNil(returnedPayload.ETag)

	})
}

func (suite *PayloadsSuite) TestEntitlement() {
	waf := entitlements.NewWeightAllotmentFetcher()
	suite.Run("Success - Returns the entitlement payload with only required fields", func() {
		entitlement := models.Entitlement{
			ID:                             uuid.Must(uuid.NewV4()),
			DependentsAuthorized:           nil,
			TotalDependents:                nil,
			NonTemporaryStorage:            nil,
			PrivatelyOwnedVehicle:          nil,
			DBAuthorizedWeight:             nil,
			UBAllowance:                    nil,
			StorageInTransit:               nil,
			RequiredMedicalEquipmentWeight: 0,
			OrganizationalClothingAndIndividualEquipment: false,
			ProGearWeight:       0,
			ProGearWeightSpouse: 0,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			WeightRestriction:   models.IntPointer(1000),
			UBWeightRestriction: models.IntPointer(1200),
		}

		payload := Entitlement(&entitlement)

		suite.Equal(strfmt.UUID(entitlement.ID.String()), payload.ID)
		suite.Equal(int64(0), payload.RequiredMedicalEquipmentWeight)
		suite.Equal(false, payload.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(int64(0), payload.ProGearWeight)
		suite.Equal(int64(0), payload.ProGearWeightSpouse)
		suite.NotEmpty(payload.ETag)
		suite.Equal(etag.GenerateEtag(entitlement.UpdatedAt), payload.ETag)

		suite.Nil(payload.AuthorizedWeight)
		suite.Nil(payload.DependentsAuthorized)
		suite.Nil(payload.NonTemporaryStorage)
		suite.Nil(payload.PrivatelyOwnedVehicle)

		/* These fields are defaulting to zero if they are nil in the model */
		suite.Equal(int64(0), payload.StorageInTransit)
		suite.Equal(int64(0), payload.TotalDependents)
		suite.Equal(int64(0), payload.TotalWeight)
		suite.Equal(int64(0), *payload.UnaccompaniedBaggageAllowance)
		suite.Equal(int64(1000), *payload.WeightRestriction)
		suite.Equal(int64(1200), *payload.UbWeightRestriction)
	})

	suite.Run("Success - Returns the entitlement payload with all optional fields populated", func() {
		entitlement := models.Entitlement{
			ID:                             uuid.Must(uuid.NewV4()),
			DependentsAuthorized:           handlers.FmtBool(true),
			TotalDependents:                handlers.FmtInt(2),
			NonTemporaryStorage:            handlers.FmtBool(true),
			PrivatelyOwnedVehicle:          handlers.FmtBool(true),
			DBAuthorizedWeight:             handlers.FmtInt(10000),
			UBAllowance:                    handlers.FmtInt(400),
			StorageInTransit:               handlers.FmtInt(45),
			RequiredMedicalEquipmentWeight: 500,
			OrganizationalClothingAndIndividualEquipment: true,
			ProGearWeight:       1000,
			ProGearWeightSpouse: 750,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		// TotalWeight needs to read from the internal weightAllotment, in this case 7000 lbs w/o dependents and
		// 9000 lbs with dependents
		allotment, err := waf.GetWeightAllotment(suite.AppContextForTest(), string(models.ServiceMemberGradeE5), internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
		suite.NoError(err)
		entitlement.WeightAllotted = &allotment

		payload := Entitlement(&entitlement)

		suite.Equal(strfmt.UUID(entitlement.ID.String()), payload.ID)
		suite.True(*payload.DependentsAuthorized)
		suite.Equal(int64(2), payload.TotalDependents)
		suite.True(*payload.NonTemporaryStorage)
		suite.True(*payload.PrivatelyOwnedVehicle)
		suite.Equal(int64(10000), *payload.AuthorizedWeight)
		suite.Equal(int64(400), *payload.UnaccompaniedBaggageAllowance)

		suite.Equal(int64(9000), payload.TotalWeight)
		suite.Equal(int64(45), payload.StorageInTransit)
		suite.Equal(int64(500), payload.RequiredMedicalEquipmentWeight)
		suite.Equal(true, payload.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(int64(1000), payload.ProGearWeight)
		suite.Equal(int64(750), payload.ProGearWeightSpouse)
		suite.NotEmpty(payload.ETag)
		suite.Equal(etag.GenerateEtag(entitlement.UpdatedAt), payload.ETag)
	})

	suite.Run("Success - Returns the entitlement payload with total weight self when dependents are not authorized", func() {
		entitlement := models.Entitlement{
			ID:                             uuid.Must(uuid.NewV4()),
			DependentsAuthorized:           handlers.FmtBool(false),
			TotalDependents:                handlers.FmtInt(2),
			NonTemporaryStorage:            handlers.FmtBool(true),
			PrivatelyOwnedVehicle:          handlers.FmtBool(true),
			DBAuthorizedWeight:             handlers.FmtInt(10000),
			UBAllowance:                    handlers.FmtInt(400),
			StorageInTransit:               handlers.FmtInt(45),
			RequiredMedicalEquipmentWeight: 500,
			OrganizationalClothingAndIndividualEquipment: true,
			ProGearWeight:       1000,
			ProGearWeightSpouse: 750,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		// TotalWeight needs to read from the internal weightAllotment, in this case 7000 lbs w/o dependents and
		// 9000 lbs with dependents
		allotment, err := waf.GetWeightAllotment(suite.AppContextForTest(), string(models.ServiceMemberGradeE5), internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
		suite.NoError(err)
		entitlement.WeightAllotted = &allotment

		payload := Entitlement(&entitlement)

		suite.Equal(strfmt.UUID(entitlement.ID.String()), payload.ID)
		suite.False(*payload.DependentsAuthorized)
		suite.Equal(int64(2), payload.TotalDependents)
		suite.True(*payload.NonTemporaryStorage)
		suite.True(*payload.PrivatelyOwnedVehicle)
		suite.Equal(int64(10000), *payload.AuthorizedWeight)
		suite.Equal(int64(400), *payload.UnaccompaniedBaggageAllowance)
		suite.Equal(int64(7000), payload.TotalWeight)
		suite.Equal(int64(45), payload.StorageInTransit)
		suite.Equal(int64(500), payload.RequiredMedicalEquipmentWeight)
		suite.Equal(true, payload.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(int64(1000), payload.ProGearWeight)
		suite.Equal(int64(750), payload.ProGearWeightSpouse)
		suite.NotEmpty(payload.ETag)
		suite.Equal(etag.GenerateEtag(entitlement.UpdatedAt), payload.ETag)
	})
}

func (suite *PayloadsSuite) TestValidationError() {
	instanceID, _ := uuid.NewV4()
	detail := "Err"

	noValidationErrors := ValidationError(detail, instanceID, nil)
	suite.Equal(handlers.ValidationErrMessage, *noValidationErrors.ClientError.Title)
	suite.Equal(detail, *noValidationErrors.ClientError.Detail)
	suite.Equal(instanceID.String(), noValidationErrors.ClientError.Instance.String())
	suite.Nil(noValidationErrors.InvalidFields)

	valErrors := validate.NewErrors()
	valErrors.Add("Field1", "dummy")
	valErrors.Add("Field2", "dummy")

	withValidationErrors := ValidationError(detail, instanceID, valErrors)
	suite.Equal(handlers.ValidationErrMessage, *withValidationErrors.ClientError.Title)
	suite.Equal(detail, *withValidationErrors.ClientError.Detail)
	suite.Equal(instanceID.String(), withValidationErrors.ClientError.Instance.String())
	suite.NotNil(withValidationErrors.InvalidFields)
	suite.Equal(2, len(withValidationErrors.InvalidFields))
}

func (suite *PayloadsSuite) TestMTOShipment() {
	primeAcknowledgeAt := time.Now().AddDate(0, 0, -5)
	mtoShipment := factory.BuildMTOShipment(nil, []factory.Customization{
		{
			Model: models.MTOShipment{
				PrimeAcknowledgedAt: &primeAcknowledgeAt,
				Status:              models.MTOShipmentStatusApproved,
			},
		},
	}, nil)
	payload := MTOShipment(&mtoShipment)
	suite.NotNil(payload)
	suite.Empty(payload.MtoServiceItems())
	suite.Equal(strfmt.UUID(mtoShipment.ID.String()), payload.ID)
	suite.Equal(handlers.FmtDatePtr(mtoShipment.ActualPickupDate), payload.ActualPickupDate)
	suite.Equal(handlers.FmtDatePtr(mtoShipment.RequestedDeliveryDate), payload.RequestedDeliveryDate)
	suite.Equal(handlers.FmtDatePtr(mtoShipment.RequestedPickupDate), payload.RequestedPickupDate)
	suite.Equal(string(mtoShipment.Status), payload.Status)
	suite.Equal(strfmt.DateTime(mtoShipment.UpdatedAt), payload.UpdatedAt)
	suite.Equal(strfmt.DateTime(mtoShipment.CreatedAt), payload.CreatedAt)
	suite.Equal(etag.GenerateEtag(mtoShipment.UpdatedAt), payload.ETag)
	suite.Equal(handlers.FmtDateTimePtr(mtoShipment.PrimeAcknowledgedAt), payload.PrimeAcknowledgedAt)

	mtoShipment = models.MTOShipment{}
	mtoShipment.MTOServiceItems = models.MTOServiceItems{
		models.MTOServiceItem{},
	}
	payload = MTOShipment(&mtoShipment)
	suite.NotNil(payload)
	suite.NotEmpty(payload.MtoServiceItems())
}

func (suite *PayloadsSuite) TestInternalServerError() {
	traceID, _ := uuid.NewV4()
	detail := "Err"

	noDetailError := InternalServerError(nil, traceID)
	suite.Equal(handlers.InternalServerErrMessage, *noDetailError.Title)
	suite.Equal(handlers.InternalServerErrDetail, *noDetailError.Detail)
	suite.Equal(traceID.String(), noDetailError.Instance.String())

	detailError := InternalServerError(&detail, traceID)
	suite.Equal(handlers.InternalServerErrMessage, *detailError.Title)
	suite.Equal(detail, *detailError.Detail)
	suite.Equal(traceID.String(), detailError.Instance.String())
}

func (suite *PayloadsSuite) TestGetDimension() {
	dimensionType := models.DimensionTypeItem
	dimensions := models.MTOServiceItemDimensions{
		models.MTOServiceItemDimension{
			Type:   dimensionType,
			Length: unit.ThousandthInches(100),
		},
		models.MTOServiceItemDimension{
			Type:   models.DimensionTypeCrate,
			Length: unit.ThousandthInches(200),
		},
	}

	resultDimension := GetDimension(dimensions, dimensionType)
	suite.Equal(dimensionType, resultDimension.Type)
	suite.Equal(unit.ThousandthInches(100), resultDimension.Length)

	emptyResultDimension := GetDimension(models.MTOServiceItemDimensions{}, dimensionType)
	suite.Equal(models.MTOServiceItemDimension{}, emptyResultDimension)
}

func (suite *PayloadsSuite) TestProofOfServiceDoc() {
	proofOfServiceDoc := models.ProofOfServiceDoc{
		PrimeUploads: []models.PrimeUpload{
			{Upload: models.Upload{ID: uuid.Must(uuid.NewV4())}},
		},
	}

	result := ProofOfServiceDoc(proofOfServiceDoc)

	suite.NotNil(result)
	suite.Equal(len(proofOfServiceDoc.PrimeUploads), len(result.Uploads))
}

func (suite *PayloadsSuite) TestPaymentRequest() {
	paymentRequest := models.PaymentRequest{
		ID: uuid.Must(uuid.NewV4()),
	}

	result := PaymentRequest(&paymentRequest)

	suite.NotNil(result)
	suite.Equal(strfmt.UUID(paymentRequest.ID.String()), result.ID)
}

func (suite *PayloadsSuite) TestPaymentRequests() {
	paymentRequests := models.PaymentRequests{
		models.PaymentRequest{ID: uuid.Must(uuid.NewV4())},
	}

	result := PaymentRequests(&paymentRequests)

	suite.NotNil(result)
	suite.Equal(len(paymentRequests), len(*result))
}

func (suite *PayloadsSuite) TestPaymentServiceItem() {
	paymentServiceItem := models.PaymentServiceItem{
		ID: uuid.Must(uuid.NewV4()),
	}

	result := PaymentServiceItem(&paymentServiceItem)

	suite.NotNil(result)
	suite.Equal(strfmt.UUID(paymentServiceItem.ID.String()), result.ID)
}

func (suite *PayloadsSuite) TestPaymentServiceItems() {
	paymentServiceItems := models.PaymentServiceItems{
		models.PaymentServiceItem{ID: uuid.Must(uuid.NewV4())},
	}

	result := PaymentServiceItems(&paymentServiceItems)

	suite.NotNil(result)
	suite.Equal(len(paymentServiceItems), len(*result))
}

func (suite *PayloadsSuite) TestPaymentServiceItemParam() {
	paymentServiceItemParam := models.PaymentServiceItemParam{
		ID: uuid.Must(uuid.NewV4()),
	}

	result := PaymentServiceItemParam(&paymentServiceItemParam)

	suite.NotNil(result)
	suite.Equal(strfmt.UUID(paymentServiceItemParam.ID.String()), result.ID)
}

func (suite *PayloadsSuite) TestPaymentServiceItemParams() {
	paymentServiceItemParams := models.PaymentServiceItemParams{
		models.PaymentServiceItemParam{ID: uuid.Must(uuid.NewV4())},
	}

	result := PaymentServiceItemParams(&paymentServiceItemParams)

	suite.NotNil(result)
	suite.Equal(len(paymentServiceItemParams), len(*result))
}

func (suite *PayloadsSuite) TestServiceRequestDocument() {
	serviceRequestDocument := models.ServiceRequestDocument{
		ServiceRequestDocumentUploads: []models.ServiceRequestDocumentUpload{
			{Upload: models.Upload{ID: uuid.Must(uuid.NewV4())}},
		},
	}

	result := ServiceRequestDocument(serviceRequestDocument)

	suite.NotNil(result)
	suite.Equal(len(serviceRequestDocument.ServiceRequestDocumentUploads), len(result.Uploads))
}

func (suite *PayloadsSuite) TestPPMShipment() {
	isActualExpenseReimbursemnt := true
	ppmShipment := &models.PPMShipment{
		ID:                           uuid.Must(uuid.NewV4()),
		IsActualExpenseReimbursement: &isActualExpenseReimbursemnt,
	}

	result := PPMShipment(ppmShipment)

	suite.NotNil(result)
	suite.Equal(strfmt.UUID(ppmShipment.ID.String()), result.ID)
	suite.True(*ppmShipment.IsActualExpenseReimbursement)
}

func (suite *PayloadsSuite) TestPPMShipmentContainingOptionalDestinationStreet1() {
	now := time.Now()
	ppmShipment := &models.PPMShipment{
		ID: uuid.Must(uuid.NewV4()),
		DestinationAddress: &models.Address{
			ID:             uuid.Must(uuid.NewV4()),
			StreetAddress1: models.STREET_ADDRESS_1_NOT_PROVIDED,
			StreetAddress2: models.StringPointer("1"),
			StreetAddress3: models.StringPointer("2"),
			City:           "SomeCity",
			State:          "CA",
			PostalCode:     "90210",
			County:         models.StringPointer("SomeCounty"),
			UpdatedAt:      now,
		},
	}

	result := PPMShipment(ppmShipment)

	eTag := etag.GenerateEtag(now)

	suite.NotNil(result)
	// expecting empty string on the response side to simulate nothing was provided.
	suite.Equal(result.DestinationAddress.StreetAddress1, models.StringPointer(""))
	suite.Equal(result.DestinationAddress.StreetAddress2, ppmShipment.DestinationAddress.StreetAddress2)
	suite.Equal(result.DestinationAddress.StreetAddress3, ppmShipment.DestinationAddress.StreetAddress3)
	suite.Equal(*result.DestinationAddress.City, ppmShipment.DestinationAddress.City)
	suite.Equal(*result.DestinationAddress.State, ppmShipment.DestinationAddress.State)
	suite.Equal(*result.DestinationAddress.PostalCode, ppmShipment.DestinationAddress.PostalCode)
	suite.Equal(result.DestinationAddress.County, ppmShipment.DestinationAddress.County)
	suite.Equal(result.DestinationAddress.ETag, eTag)
}

func (suite *PayloadsSuite) TestMTOServiceItem() {
	sitPostalCode := "55555"
	mtoServiceItemDOFSIT := &models.MTOServiceItem{
		ID:               uuid.Must(uuid.NewV4()),
		ReService:        models.ReService{Code: models.ReServiceCodeDOFSIT},
		SITDepartureDate: nil,
		SITEntryDate:     nil,
		SITPostalCode:    &sitPostalCode,
		SITOriginHHGActualAddress: &models.Address{
			StreetAddress1: "dummyStreet",
			City:           "dummyCity",
			State:          "FL",
			PostalCode:     "55555",
		},
		SITOriginHHGOriginalAddress: &models.Address{
			StreetAddress1: "dummyStreet2",
			City:           "dummyCity2",
			State:          "FL",
			PostalCode:     "55555",
		},
	}

	resultDOFSIT := MTOServiceItem(mtoServiceItemDOFSIT)
	suite.NotNil(resultDOFSIT)
	sitOrigin, ok := resultDOFSIT.(*primev3messages.MTOServiceItemOriginSIT)
	suite.True(ok)
	suite.Equal("55555", *sitOrigin.SitPostalCode)
	suite.Equal("dummyStreet", *sitOrigin.SitHHGActualOrigin.StreetAddress1)
	suite.Equal("dummyStreet2", *sitOrigin.SitHHGOriginalOrigin.StreetAddress1)

	mtoServiceItemDefault := &models.MTOServiceItem{
		ID:              uuid.Must(uuid.NewV4()),
		ReService:       models.ReService{Code: "SOME_OTHER_SERVICE_CODE"},
		MoveTaskOrderID: uuid.Must(uuid.NewV4()),
	}

	resultDefault := MTOServiceItem(mtoServiceItemDefault)
	suite.NotNil(resultDefault)
	basicItem, ok := resultDefault.(*primev3messages.MTOServiceItemBasic)
	suite.True(ok)
	suite.Equal("SOME_OTHER_SERVICE_CODE", string(*basicItem.ReServiceCode))
	suite.Equal(mtoServiceItemDefault.MoveTaskOrderID.String(), basicItem.MoveTaskOrderID().String())
}

func (suite *PayloadsSuite) TestGetCustomerContact() {
	customerContacts := models.MTOServiceItemCustomerContacts{
		models.MTOServiceItemCustomerContact{Type: models.CustomerContactTypeFirst},
	}
	contactType := models.CustomerContactTypeFirst

	result := GetCustomerContact(customerContacts, contactType)

	suite.Equal(models.CustomerContactTypeFirst, result.Type)
}

func (suite *PayloadsSuite) TestShipmentAddressUpdate() {
	shipmentAddressUpdate := &models.ShipmentAddressUpdate{
		ID: uuid.Must(uuid.NewV4()),
	}

	result := ShipmentAddressUpdate(shipmentAddressUpdate)

	suite.NotNil(result)
	suite.Equal(strfmt.UUID(shipmentAddressUpdate.ID.String()), result.ID)
}

func (suite *PayloadsSuite) TestMTOServiceItemDestSIT() {
	reServiceCode := models.ReServiceCodeDDFSIT
	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)
	sitDepartureDate := time.Now().AddDate(0, 1, 0)
	sitEntryDate := time.Now().AddDate(0, 0, -30)
	finalAddress := models.Address{
		StreetAddress1: "dummyStreet",
		City:           "dummyCity",
		State:          "FL",
		PostalCode:     "55555",
	}
	mtoShipmentID := uuid.Must(uuid.NewV4())

	mtoServiceItemDestSIT := &models.MTOServiceItem{
		ID:                         uuid.Must(uuid.NewV4()),
		ReService:                  models.ReService{Code: reServiceCode},
		Reason:                     &reason,
		SITDepartureDate:           &sitDepartureDate,
		SITEntryDate:               &sitEntryDate,
		SITDestinationFinalAddress: &finalAddress,
		MTOShipmentID:              &mtoShipmentID,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	resultDestSIT := MTOServiceItem(mtoServiceItemDestSIT)
	suite.NotNil(resultDestSIT)
	destSIT, ok := resultDestSIT.(*primev3messages.MTOServiceItemDestSIT)
	suite.True(ok)

	suite.Equal(string(reServiceCode), string(*destSIT.ReServiceCode))
	suite.Equal(reason, *destSIT.Reason)
	suite.Equal(strfmt.Date(sitDepartureDate).String(), destSIT.SitDepartureDate.String())
	suite.Equal(strfmt.Date(sitEntryDate).String(), destSIT.SitEntryDate.String())
	suite.Equal(strfmt.Date(dateOfContact1).String(), destSIT.DateOfContact1.String())
	suite.Equal(timeMilitary1, *destSIT.TimeMilitary1)
	suite.Equal(strfmt.Date(firstAvailableDeliveryDate1).String(), destSIT.FirstAvailableDeliveryDate1.String())
	suite.Equal(strfmt.Date(dateOfContact2).String(), destSIT.DateOfContact2.String())
	suite.Equal(timeMilitary2, *destSIT.TimeMilitary2)
	suite.Equal(strfmt.Date(firstAvailableDeliveryDate2).String(), destSIT.FirstAvailableDeliveryDate2.String())
	suite.Equal(finalAddress.StreetAddress1, *destSIT.SitDestinationFinalAddress.StreetAddress1)
	suite.Equal(finalAddress.City, *destSIT.SitDestinationFinalAddress.City)
	suite.Equal(finalAddress.State, *destSIT.SitDestinationFinalAddress.State)
	suite.Equal(finalAddress.PostalCode, *destSIT.SitDestinationFinalAddress.PostalCode)
	suite.Equal(mtoShipmentID.String(), destSIT.MtoShipmentID().String())
}

func (suite *PayloadsSuite) TestMTOServiceItemInternationalDestSIT() {
	reServiceCode := models.ReServiceCodeIDFSIT
	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)
	sitDepartureDate := time.Now().AddDate(0, 1, 0)
	sitEntryDate := time.Now().AddDate(0, 0, -30)
	finalAddress := models.Address{
		StreetAddress1: "dummyStreet",
		City:           "dummyCity",
		State:          "FL",
		PostalCode:     "55555",
	}
	mtoShipmentID := uuid.Must(uuid.NewV4())

	mtoServiceItemDestSIT := &models.MTOServiceItem{
		ID:                         uuid.Must(uuid.NewV4()),
		ReService:                  models.ReService{Code: reServiceCode},
		Reason:                     &reason,
		SITDepartureDate:           &sitDepartureDate,
		SITEntryDate:               &sitEntryDate,
		SITDestinationFinalAddress: &finalAddress,
		MTOShipmentID:              &mtoShipmentID,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	resultDestSIT := MTOServiceItem(mtoServiceItemDestSIT)
	suite.NotNil(resultDestSIT)
	destSIT, ok := resultDestSIT.(*primev3messages.MTOServiceItemInternationalDestSIT)
	suite.True(ok)

	suite.Equal(string(reServiceCode), string(*destSIT.ReServiceCode))
	suite.Equal(reason, *destSIT.Reason)
	suite.Equal(strfmt.Date(sitDepartureDate).String(), destSIT.SitDepartureDate.String())
	suite.Equal(strfmt.Date(sitEntryDate).String(), destSIT.SitEntryDate.String())
	suite.Equal(strfmt.Date(dateOfContact1).String(), destSIT.DateOfContact1.String())
	suite.Equal(timeMilitary1, *destSIT.TimeMilitary1)
	suite.Equal(strfmt.Date(firstAvailableDeliveryDate1).String(), destSIT.FirstAvailableDeliveryDate1.String())
	suite.Equal(strfmt.Date(dateOfContact2).String(), destSIT.DateOfContact2.String())
	suite.Equal(timeMilitary2, *destSIT.TimeMilitary2)
	suite.Equal(strfmt.Date(firstAvailableDeliveryDate2).String(), destSIT.FirstAvailableDeliveryDate2.String())
	suite.Equal(finalAddress.StreetAddress1, *destSIT.SitDestinationFinalAddress.StreetAddress1)
	suite.Equal(finalAddress.City, *destSIT.SitDestinationFinalAddress.City)
	suite.Equal(finalAddress.State, *destSIT.SitDestinationFinalAddress.State)
	suite.Equal(finalAddress.PostalCode, *destSIT.SitDestinationFinalAddress.PostalCode)
	suite.Equal(mtoShipmentID.String(), destSIT.MtoShipmentID().String())
}

func (suite *PayloadsSuite) TestMTOServiceItemDCRT() {
	reServiceCode := models.ReServiceCodeDCRT
	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)
	standaloneCrate := false

	mtoServiceItemDCRT := &models.MTOServiceItem{
		ID:              uuid.Must(uuid.NewV4()),
		ReService:       models.ReService{Code: reServiceCode},
		Reason:          &reason,
		StandaloneCrate: &standaloneCrate,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	resultDCRT := MTOServiceItem(mtoServiceItemDCRT)

	suite.NotNil(resultDCRT)

	_, ok := resultDCRT.(*primev3messages.MTOServiceItemDomesticCrating)

	suite.True(ok)
}

func (suite *PayloadsSuite) TestMTOServiceItemICRTandIUCRT() {
	icrtReServiceCode := models.ReServiceCodeICRT
	iucrtReServiceCode := models.ReServiceCodeIUCRT
	reason := "reason"
	standaloneCrate := false
	externalCrate := false
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)

	mtoServiceItemICRT := &models.MTOServiceItem{
		ID:              uuid.Must(uuid.NewV4()),
		ReService:       models.ReService{Code: icrtReServiceCode},
		Reason:          &reason,
		StandaloneCrate: &standaloneCrate,
		ExternalCrate:   &externalCrate,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	mtoServiceItemIUCRT := &models.MTOServiceItem{
		ID:              uuid.Must(uuid.NewV4()),
		ReService:       models.ReService{Code: iucrtReServiceCode},
		Reason:          &reason,
		StandaloneCrate: &standaloneCrate,
		ExternalCrate:   &externalCrate,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	resultICRT := MTOServiceItem(mtoServiceItemICRT)
	resultIUCRT := MTOServiceItem(mtoServiceItemIUCRT)

	suite.NotNil(resultICRT)
	suite.NotNil(resultIUCRT)

	_, ok := resultICRT.(*primev3messages.MTOServiceItemInternationalCrating)
	suite.True(ok)

	_, ok = resultIUCRT.(*primev3messages.MTOServiceItemInternationalCrating)
	suite.True(ok)
}

func (suite *PayloadsSuite) TestMTOServiceItemDDSHUT() {
	reServiceCode := models.ReServiceCodeDDSHUT
	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)

	mtoServiceItemDDSHUT := &models.MTOServiceItem{
		ID:        uuid.Must(uuid.NewV4()),
		ReService: models.ReService{Code: reServiceCode},
		Reason:    &reason,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	resultDDSHUT := MTOServiceItem(mtoServiceItemDDSHUT)

	suite.NotNil(resultDDSHUT)

	_, ok := resultDDSHUT.(*primev3messages.MTOServiceItemDomesticShuttle)

	suite.True(ok)
}

func (suite *PayloadsSuite) TestMTOServiceItemIDSHUT() {
	reServiceCode := models.ReServiceCodeIDSHUT
	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)
	standaloneCrate := false

	mtoServiceItemIDSHUT := &models.MTOServiceItem{
		ID:              uuid.Must(uuid.NewV4()),
		ReService:       models.ReService{Code: reServiceCode},
		Reason:          &reason,
		StandaloneCrate: &standaloneCrate,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	resultIDSHUT := MTOServiceItem(mtoServiceItemIDSHUT)

	suite.NotNil(resultIDSHUT)

	_, ok := resultIDSHUT.(*primev3messages.MTOServiceItemInternationalShuttle)

	suite.True(ok)
}

func (suite *PayloadsSuite) TestMTOServiceItemIOSHUT() {
	reServiceCode := models.ReServiceCodeIOSHUT
	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)
	standaloneCrate := false

	mtoServiceItemIOSHUT := &models.MTOServiceItem{
		ID:              uuid.Must(uuid.NewV4()),
		ReService:       models.ReService{Code: reServiceCode},
		Reason:          &reason,
		StandaloneCrate: &standaloneCrate,
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact1,
				TimeMilitary:               timeMilitary1,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate1,
				Type:                       models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				DateOfContact:              dateOfContact2,
				TimeMilitary:               timeMilitary2,
				FirstAvailableDeliveryDate: firstAvailableDeliveryDate2,
				Type:                       models.CustomerContactTypeSecond,
			},
		},
	}

	resultIOSHUT := MTOServiceItem(mtoServiceItemIOSHUT)

	suite.NotNil(resultIOSHUT)

	_, ok := resultIOSHUT.(*primev3messages.MTOServiceItemInternationalShuttle)

	suite.True(ok)
}

func (suite *PayloadsSuite) TestStorageFacilityPayload() {
	phone := "555"
	email := "email"
	facility := "facility"
	lot := "lot"

	storage := &models.StorageFacility{
		ID:           uuid.Must(uuid.NewV4()),
		Address:      models.Address{},
		UpdatedAt:    time.Now(),
		Email:        &email,
		FacilityName: facility,
		LotNumber:    &lot,
		Phone:        &phone,
	}

	suite.NotNil(storage)
}

func (suite *PayloadsSuite) TestMTOAgentPayload() {
	firstName := "John"
	lastName := "Doe"
	phone := "555"
	email := "email"
	mtoAgent := &models.MTOAgent{
		ID:            uuid.Must(uuid.NewV4()),
		MTOAgentType:  models.MTOAgentReceiving,
		FirstName:     &firstName,
		LastName:      &lastName,
		Phone:         &phone,
		Email:         &email,
		MTOShipmentID: uuid.Must(uuid.NewV4()),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	payload := MTOAgent(mtoAgent)
	suite.NotNil(payload)
}

func (suite *PayloadsSuite) TestStorageFacility() {
	storageFacilityID := uuid.Must(uuid.NewV4())
	updatedAt := time.Now()
	dummy := "dummy"
	email := "dummy@example.com"
	facilityName := "dummy"
	lotNumber := "dummy"
	phone := "dummy"
	storage := &models.StorageFacility{
		ID: storageFacilityID,
		Address: models.Address{
			StreetAddress1: dummy,
			City:           dummy,
			State:          dummy,
			PostalCode:     dummy,
		},
		Email:        &email,
		FacilityName: facilityName,
		LotNumber:    &lotNumber,
		Phone:        &phone,
		UpdatedAt:    updatedAt,
	}

	result := StorageFacility(storage)
	suite.NotNil(result)
}

func (suite *PayloadsSuite) TestBoatShipment() {
	id, _ := uuid.NewV4()
	year := 2000
	make := "Test Make"
	model := "Test Model"
	lengthInInches := 400
	widthInInches := 320
	heightInInches := 300
	hasTrailer := true
	IsRoadworthy := false
	boatShipment := &models.BoatShipment{
		ID:             id,
		Type:           models.BoatShipmentTypeHaulAway,
		Year:           &year,
		Make:           &make,
		Model:          &model,
		LengthInInches: &lengthInInches,
		WidthInInches:  &widthInInches,
		HeightInInches: &heightInInches,
		HasTrailer:     &hasTrailer,
		IsRoadworthy:   &IsRoadworthy,
	}

	result := BoatShipment(boatShipment)
	suite.NotNil(result)
}

func (suite *PayloadsSuite) TestDestinationPostalCodeAndGBLOC() {
	moveID := uuid.Must(uuid.NewV4())
	moveLocator := "TESTTEST"
	primeTime := time.Now()
	ordersID := uuid.Must(uuid.NewV4())
	refID := "123456"
	contractNum := "HTC-123-456"
	address := models.Address{PostalCode: "35023"}
	shipment := models.MTOShipment{
		ID:                 uuid.Must(uuid.NewV4()),
		DestinationAddress: &address,
	}
	shipments := models.MTOShipments{shipment}
	contractor := models.Contractor{
		ContractNumber: contractNum,
	}

	basicMove := models.Move{
		ID:                   moveID,
		Locator:              moveLocator,
		CreatedAt:            primeTime,
		ReferenceID:          &refID,
		AvailableToPrimeAt:   &primeTime,
		ApprovedAt:           &primeTime,
		OrdersID:             ordersID,
		Contractor:           &contractor,
		PaymentRequests:      models.PaymentRequests{},
		SubmittedAt:          &primeTime,
		UpdatedAt:            primeTime,
		Status:               models.MoveStatusAPPROVED,
		SignedCertifications: models.SignedCertifications{},
		MTOServiceItems:      models.MTOServiceItems{},
		MTOShipments:         shipments,
	}

	suite.Run("Returns values needed to get the destination postal code and GBLOC", func() {
		returnedModel := MoveTaskOrder(suite.AppContextForTest(), &basicMove)

		suite.IsType(&primev3messages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(basicMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(strfmt.DateTime(basicMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotEmpty(returnedModel.ETag)
	})
}

func (suite *PayloadsSuite) TestMarketCode() {
	suite.Run("returns nil when marketCode is nil", func() {
		var marketCode *models.MarketCode = nil
		result := MarketCode(marketCode)
		suite.Equal(result, "")
	})

	suite.Run("returns string when marketCode is not nil", func() {
		marketCodeDomestic := models.MarketCodeDomestic
		result := MarketCode(&marketCodeDomestic)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("d", result, "Expected result to be 'd' for domestic market code")
	})

	suite.Run("returns string when marketCode is international", func() {
		marketCodeInternational := models.MarketCodeInternational
		result := MarketCode(&marketCodeInternational)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("i", result, "Expected result to be 'i' for international market code")
	})
}

func (suite *PayloadsSuite) TestMTOServiceItemPOEFSC() {

	portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
		{
			Model: models.Port{
				PortCode: "PDX",
			},
		},
	}, nil)

	poefscServiceItem := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model: models.ReService{
				Code:     models.ReServiceCodePOEFSC,
				Priority: 1,
			},
		},
		{
			Model:    portLocation,
			LinkOnly: true,
			Type:     &factory.PortLocations.PortOfEmbarkation,
		},
	}, nil)
	mtoServiceItems := [...]models.MTOServiceItem{poefscServiceItem}

	mtoShipment := models.MTOShipment{
		ID: *poefscServiceItem.MTOShipmentID,
	}
	mtoShipments := [...]models.MTOShipment{mtoShipment}

	move := models.Move{
		MTOShipments:    mtoShipments[:],
		MTOServiceItems: mtoServiceItems[:],
		ReferenceID:     poefscServiceItem.MoveTaskOrder.ReferenceID,
		Contractor: &models.Contractor{
			ContractNumber: factory.DefaultContractNumber,
		},
	}

	mtoPayload := MoveTaskOrder(suite.AppContextForTest(), &move)
	suite.NotNil(mtoPayload)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.PortType, portLocation.Port.PortType.String())
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.PortCode, portLocation.Port.PortCode)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.PortName, portLocation.Port.PortName)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.City, portLocation.City.CityName)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.County, portLocation.UsPostRegionCity.UsprcCountyNm)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.State, portLocation.UsPostRegionCity.UsPostRegion.State.StateName)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.Zip, portLocation.UsPostRegionCity.UsprZipID)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfEmbarkation.Country, portLocation.Country.CountryName)
}

func (suite *PayloadsSuite) TestMTOServiceItemPODFSC() {

	portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
		{
			Model: models.Port{
				PortCode: "PDX",
			},
		},
	}, nil)

	podfscServiceItem := factory.BuildMTOServiceItem(nil, []factory.Customization{
		{
			Model: models.ReService{
				Code:     models.ReServiceCodePODFSC,
				Priority: 1,
			},
		},
		{
			Model:    portLocation,
			LinkOnly: true,
			Type:     &factory.PortLocations.PortOfDebarkation,
		},
	}, nil)
	mtoServiceItems := [...]models.MTOServiceItem{podfscServiceItem}

	mtoShipment := models.MTOShipment{
		ID: *podfscServiceItem.MTOShipmentID,
	}
	mtoShipments := [...]models.MTOShipment{mtoShipment}

	move := models.Move{
		MTOShipments:    mtoShipments[:],
		MTOServiceItems: mtoServiceItems[:],
		ReferenceID:     podfscServiceItem.MoveTaskOrder.ReferenceID,
		Contractor: &models.Contractor{
			ContractNumber: factory.DefaultContractNumber,
		},
	}

	mtoPayload := MoveTaskOrder(suite.AppContextForTest(), &move)
	suite.NotNil(mtoPayload)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.PortType, portLocation.Port.PortType.String())
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.PortCode, portLocation.Port.PortCode)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.PortName, portLocation.Port.PortName)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.City, portLocation.City.CityName)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.County, portLocation.UsPostRegionCity.UsprcCountyNm)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.State, portLocation.UsPostRegionCity.UsPostRegion.State.StateName)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.Zip, portLocation.UsPostRegionCity.UsprZipID)
	suite.Equal(mtoPayload.MtoShipments[0].PortOfDebarkation.Country, portLocation.Country.CountryName)
}
func (suite *PayloadsSuite) TestAddress() {
	usprcId := uuid.Must(uuid.NewV4())
	shipmentAddress := &models.Address{
		ID:                 uuid.Must(uuid.NewV4()),
		StreetAddress1:     "400 Drive St",
		City:               "Charleston",
		County:             models.StringPointer("Charleston"),
		State:              "SC",
		PostalCode:         "29404",
		UsPostRegionCityID: &usprcId,
	}

	result := Address(shipmentAddress)
	suite.NotNil(result)
	suite.Equal(strfmt.UUID(shipmentAddress.ID.String()), result.ID)
	suite.Equal(strfmt.UUID(usprcId.String()), result.UsPostRegionCitiesID)

	result = Address(nil)
	suite.Nil(result)

	usprcId = uuid.Nil
	shipmentAddress.UsPostRegionCityID = &uuid.Nil
	result = Address(shipmentAddress)
	suite.NotNil(result)
	suite.Equal(strfmt.UUID(""), result.UsPostRegionCitiesID)
}
