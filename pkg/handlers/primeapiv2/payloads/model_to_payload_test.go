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
	excessWeightAcknowledgedAt := time.Now()
	excessWeightUploadID := uuid.Must(uuid.NewV4())
	ordersType := primev2messages.OrdersTypeRETIREMENT
	originDutyGBLOC := "KKFA"
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
		ReferenceID:                &referenceID,
		PaymentRequests:            models.PaymentRequests{},
		SubmittedAt:                &submittedAt,
		UpdatedAt:                  time.Now(),
		PersonallyProcuredMoves:    models.PersonallyProcuredMoves{},
		Status:                     models.MoveStatusAPPROVED,
		SignedCertifications:       models.SignedCertifications{},
		MTOServiceItems:            models.MTOServiceItems{},
		MTOShipments:               models.MTOShipments{},
		ExcessWeightQualifiedAt:    &excessWeightQualifiedAt,
		ExcessWeightAcknowledgedAt: &excessWeightAcknowledgedAt,
		ExcessWeightUploadID:       &excessWeightUploadID,
		Contractor: &models.Contractor{
			ContractNumber: factory.DefaultContractNumber,
		},
	}

	suite.Run("Success - Returns a basic move payload with no payment requests, service items or shipments", func() {
		returnedModel := MoveTaskOrder(&basicMove)

		suite.IsType(&primev2messages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(basicMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(ordersType, returnedModel.Order.OrdersType)
		suite.Equal(originDutyGBLOC, returnedModel.Order.OriginDutyLocationGBLOC)
		suite.Equal(referenceID, returnedModel.ReferenceID)
		suite.Equal(strfmt.DateTime(basicMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotEmpty(returnedModel.ETag)
		suite.True(returnedModel.ExcessWeightQualifiedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightQualifiedAt)))
		suite.True(returnedModel.ExcessWeightAcknowledgedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightAcknowledgedAt)))
		suite.Require().NotNil(returnedModel.ExcessWeightUploadID)
		suite.Equal(strfmt.UUID(basicMove.ExcessWeightUploadID.String()), *returnedModel.ExcessWeightUploadID)
		suite.Equal(factory.DefaultContractNumber, returnedModel.ContractNumber)
		suite.Equal(models.SupplyAndServicesCostEstimate, returnedModel.Order.SupplyAndServicesCostEstimate)
		suite.Equal(models.MethodOfPayment, returnedModel.Order.MethodOfPayment)
		suite.Equal(models.NAICS, returnedModel.Order.Naics)
		suite.Equal(packingInstructions, returnedModel.Order.PackingAndShippingInstructions)
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
	})

	suite.Run("Success - Returns the entitlement payload with all optional fields populated", func() {
		entitlement := models.Entitlement{
			ID:                             uuid.Must(uuid.NewV4()),
			DependentsAuthorized:           handlers.FmtBool(true),
			TotalDependents:                handlers.FmtInt(2),
			NonTemporaryStorage:            handlers.FmtBool(true),
			PrivatelyOwnedVehicle:          handlers.FmtBool(true),
			DBAuthorizedWeight:             handlers.FmtInt(10000),
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
		entitlement.SetWeightAllotment(string(models.ServiceMemberRankE5))

		payload := Entitlement(&entitlement)

		suite.Equal(strfmt.UUID(entitlement.ID.String()), payload.ID)
		suite.True(*payload.DependentsAuthorized)
		suite.Equal(int64(2), payload.TotalDependents)
		suite.True(*payload.NonTemporaryStorage)
		suite.True(*payload.PrivatelyOwnedVehicle)
		suite.Equal(int64(10000), *payload.AuthorizedWeight)
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
		entitlement.SetWeightAllotment(string(models.ServiceMemberRankE5))

		payload := Entitlement(&entitlement)

		suite.Equal(strfmt.UUID(entitlement.ID.String()), payload.ID)
		suite.False(*payload.DependentsAuthorized)
		suite.Equal(int64(2), payload.TotalDependents)
		suite.True(*payload.NonTemporaryStorage)
		suite.True(*payload.PrivatelyOwnedVehicle)
		suite.Equal(int64(10000), *payload.AuthorizedWeight)
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

func (suite *PayloadsSuite) TestSITAddressUpdate() {
	newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
	contractorRemark := "I must update the final address please"
	officeRemark := ""

	suite.Run("Success - Returns a SITAddressUpdate payload as expected", func() {
		sitAddressUpdate := models.SITAddressUpdate{
			ID:                uuid.Must(uuid.NewV4()),
			MTOServiceItemID:  uuid.Must(uuid.NewV4()),
			NewAddressID:      newAddress.ID,
			NewAddress:        newAddress,
			ContractorRemarks: &contractorRemark,
			OfficeRemarks:     &officeRemark,
			Status:            models.SITAddressUpdateStatusRequested,
			UpdatedAt:         time.Now(),
			CreatedAt:         time.Now(),
		}

		payload := SITAddressUpdate(&sitAddressUpdate)

		suite.Equal(payload.ID.String(), sitAddressUpdate.ID.String())
		suite.Equal(payload.MtoServiceItemID.String(), sitAddressUpdate.MTOServiceItemID.String())
		suite.Equal(payload.NewAddressID.String(), sitAddressUpdate.NewAddressID.String())
		suite.Equal(payload.NewAddress.ID.String(), sitAddressUpdate.NewAddress.ID.String())
		suite.Equal(*payload.NewAddress.City, sitAddressUpdate.NewAddress.City)
		suite.Equal(*payload.NewAddress.State, sitAddressUpdate.NewAddress.State)
		suite.Equal(*payload.NewAddress.PostalCode, sitAddressUpdate.NewAddress.PostalCode)
		suite.Equal(*payload.NewAddress.Country, *sitAddressUpdate.NewAddress.Country)
		suite.Equal(*payload.NewAddress.StreetAddress1, sitAddressUpdate.NewAddress.StreetAddress1)
		suite.Equal(payload.ContractorRemarks, sitAddressUpdate.ContractorRemarks)
		suite.Equal(payload.OfficeRemarks, sitAddressUpdate.OfficeRemarks)
		suite.Equal(payload.Status, sitAddressUpdate.Status)
		suite.Equal(strfmt.DateTime(payload.UpdatedAt).String(), strfmt.DateTime(sitAddressUpdate.UpdatedAt).String())
		suite.Equal(strfmt.DateTime(payload.CreatedAt).String(), strfmt.DateTime(sitAddressUpdate.CreatedAt).String())
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
