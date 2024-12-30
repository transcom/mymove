package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/primev2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestMoveTaskOrder() {
	moveTaskOrderID, _ := uuid.NewV4()
	ordersID, _ := uuid.NewV4()
	referenceID := "testID"
	primeTime := time.Now()
	submittedAt := time.Now()
	excessWeightQualifiedAt := time.Now()
	excessUnaccompaniedBaggageWeightQualifiedAt := time.Now()
	excessWeightAcknowledgedAt := time.Now()
	excessUnaccompaniedBaggageWeightAcknowledgedAt := time.Now()
	excessWeightUploadID := uuid.Must(uuid.NewV4())
	ordersType := primev2messages.OrdersTypeRETIREMENT
	originDutyGBLOC := "KKFA"
	shipmentGBLOC := "AGFM"
	packingInstructions := models.InstructionsBeforeContractNumber + factory.DefaultContractNumber + models.InstructionsAfterContractNumber

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
		},
		ReferenceID:             &referenceID,
		PaymentRequests:         models.PaymentRequests{},
		SubmittedAt:             &submittedAt,
		UpdatedAt:               time.Now(),
		Status:                  models.MoveStatusAPPROVED,
		SignedCertifications:    models.SignedCertifications{},
		MTOServiceItems:         models.MTOServiceItems{},
		MTOShipments:            models.MTOShipments{},
		ExcessWeightQualifiedAt: &excessWeightQualifiedAt,
		ExcessUnaccompaniedBaggageWeightQualifiedAt:    &excessUnaccompaniedBaggageWeightQualifiedAt,
		ExcessWeightAcknowledgedAt:                     &excessWeightAcknowledgedAt,
		ExcessUnaccompaniedBaggageWeightAcknowledgedAt: &excessUnaccompaniedBaggageWeightAcknowledgedAt,
		ExcessWeightUploadID:                           &excessWeightUploadID,
		Contractor: &models.Contractor{
			ContractNumber: factory.DefaultContractNumber,
		},
		ShipmentGBLOC: models.MoveToGBLOCs{
			models.MoveToGBLOC{GBLOC: &shipmentGBLOC},
		},
	}

	suite.Run("Success - Returns a basic move payload with no payment requests, service items or shipments", func() {
		returnedModel := MoveTaskOrder(suite.AppContextForTest(), &basicMove)

		suite.IsType(&primev2messages.MoveTaskOrder{}, returnedModel)
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
		suite.True(returnedModel.ExcessWeightQualifiedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightQualifiedAt)))
		suite.True(returnedModel.ExcessUnaccompaniedBaggageWeightQualifiedAt.Equal(strfmt.DateTime(*basicMove.ExcessUnaccompaniedBaggageWeightQualifiedAt)))
		suite.True(returnedModel.ExcessWeightAcknowledgedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightAcknowledgedAt)))
		suite.True(returnedModel.ExcessUnaccompaniedBaggageWeightAcknowledgedAt.Equal(strfmt.DateTime(*basicMove.ExcessUnaccompaniedBaggageWeightAcknowledgedAt)))
		suite.Require().NotNil(returnedModel.ExcessWeightUploadID)
		suite.Equal(strfmt.UUID(basicMove.ExcessWeightUploadID.String()), *returnedModel.ExcessWeightUploadID)
		suite.Equal(factory.DefaultContractNumber, returnedModel.ContractNumber)
		suite.Equal(models.SupplyAndServicesCostEstimate, returnedModel.Order.SupplyAndServicesCostEstimate)
		suite.Equal(models.MethodOfPayment, returnedModel.Order.MethodOfPayment)
		suite.Equal(models.NAICS, returnedModel.Order.Naics)
		suite.Equal(packingInstructions, returnedModel.Order.PackingAndShippingInstructions)
	})
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

		suite.IsType(&primev2messages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(basicMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(strfmt.DateTime(basicMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotEmpty(returnedModel.ETag)
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

		suite.IsType(&primev2messages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(reweigh.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(reweigh.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primev2messages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
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

		suite.IsType(&primev2messages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(reweigh.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(reweigh.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primev2messages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
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

		suite.IsType(&primev2messages.SITExtension{}, returnedPayload)
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

		suite.IsType(&primev2messages.SITExtension{}, returnedPayload)
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
		entitlement.SetWeightAllotment(string(models.ServiceMemberGradeE5), internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

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
		entitlement.SetWeightAllotment(string(models.ServiceMemberGradeE5), internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

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
	mtoShipment := &models.MTOShipment{}

	mtoShipment.MTOServiceItems = nil
	payload := MTOShipment(mtoShipment)
	suite.NotNil(payload)
	suite.Empty(payload.MtoServiceItems())

	mtoShipment.MTOServiceItems = models.MTOServiceItems{
		models.MTOServiceItem{},
	}
	payload = MTOShipment(mtoShipment)
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
	sitOrigin, ok := resultDOFSIT.(*primev2messages.MTOServiceItemOriginSIT)
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
	basicItem, ok := resultDefault.(*primev2messages.MTOServiceItemBasic)
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
	destSIT, ok := resultDestSIT.(*primev2messages.MTOServiceItemDestSIT)
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

	_, ok := resultDCRT.(*primev2messages.MTOServiceItemDomesticCrating)

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

	_, ok := resultICRT.(*primev2messages.MTOServiceItemInternationalCrating)
	suite.True(ok)

	_, ok = resultIUCRT.(*primev2messages.MTOServiceItemInternationalCrating)
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

	_, ok := resultDDSHUT.(*primev2messages.MTOServiceItemShuttle)

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
