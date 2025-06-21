package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/storage/test"
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
	ordersType := primemessages.OrdersTypeRETIREMENT
	originDutyGBLOC := "KKFA"
	shipmentGBLOC := "AGFM"

	backupContacts := models.BackupContacts{}
	backupContacts = append(backupContacts, models.BackupContact{
		FirstName: "Backup",
		LastName:  "contact name",
		Phone:     "555-555-5555",
		Email:     "backup@backup.com",
	})
	serviceMember := models.ServiceMember{
		BackupContacts: backupContacts,
	}

	basicMove := models.Move{
		ID:                      moveTaskOrderID,
		Locator:                 "TESTTEST",
		CreatedAt:               time.Now(),
		AvailableToPrimeAt:      &primeTime,
		ApprovedAt:              &primeTime,
		OrdersID:                ordersID,
		Orders:                  models.Order{OrdersType: internalmessages.OrdersType(ordersType), OriginDutyLocationGBLOC: &originDutyGBLOC, ServiceMember: serviceMember},
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
		ShipmentGBLOC: models.MoveToGBLOCs{
			models.MoveToGBLOC{GBLOC: &shipmentGBLOC},
		},
		PrimeAcknowledgedAt: &primeAcknowledgedAt,
	}

	suite.Run("Success - Returns a basic move payload with no payment requests, service items or shipments", func() {
		returnedModel := MoveTaskOrder(suite.AppContextForTest(), &basicMove)

		suite.IsType(&primemessages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.ApprovedAt), returnedModel.ApprovedAt)
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
		suite.Equal(basicMove.Orders.ServiceMember.BackupContacts[0].FirstName, returnedModel.Order.Customer.BackupContact.FirstName)
		suite.Equal(basicMove.Orders.ServiceMember.BackupContacts[0].LastName, returnedModel.Order.Customer.BackupContact.LastName)
		suite.Equal(basicMove.Orders.ServiceMember.BackupContacts[0].Phone, returnedModel.Order.Customer.BackupContact.Phone)
		suite.Equal(basicMove.Orders.ServiceMember.BackupContacts[0].Email, returnedModel.Order.Customer.BackupContact.Email)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.PrimeAcknowledgedAt), returnedModel.PrimeAcknowledgedAt)
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

		suite.IsType(&primemessages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(reweigh.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(reweigh.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primemessages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
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

		suite.IsType(&primemessages.Reweigh{}, returnedPayload)
		suite.Equal(strfmt.UUID(reweigh.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(reweigh.ShipmentID.String()), returnedPayload.ShipmentID)
		suite.Equal(strfmt.DateTime(reweigh.RequestedAt), returnedPayload.RequestedAt)
		suite.Equal(primemessages.ReweighRequester(reweigh.RequestedBy), returnedPayload.RequestedBy)
		suite.Equal(strfmt.DateTime(reweigh.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(reweigh.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Equal(handlers.FmtPoundPtr(reweigh.Weight), returnedPayload.Weight)
		suite.Equal(handlers.FmtStringPtr(reweigh.VerificationReason), returnedPayload.VerificationReason)
		suite.Equal(handlers.FmtDateTimePtr(reweigh.VerificationProvidedAt), returnedPayload.VerificationProvidedAt)
		suite.NotEmpty(returnedPayload.ETag)
	})
}

func (suite *PayloadsSuite) TestExcessWeightRecord() {
	id, err := uuid.NewV4()
	suite.Require().NoError(err, "Unexpected error when generating new UUID")

	now := time.Now()
	fakeFileStorer := test.NewFakeS3Storage(true)

	suite.Run("Success - all data populated", func() {
		// Get stubbed upload with ID and timestamps
		upload := factory.BuildUpload(nil, []factory.Customization{
			{
				Model: models.Upload{ID: uuid.Must(uuid.NewV4())},
			},
		}, []factory.Trait{factory.GetTraitTimestampedUpload})

		move := models.Move{
			ID:                      id,
			ExcessWeightQualifiedAt: &now,
			ExcessUnaccompaniedBaggageWeightQualifiedAt:    &now,
			ExcessWeightAcknowledgedAt:                     &now,
			ExcessUnaccompaniedBaggageWeightAcknowledgedAt: &now,
			ExcessWeightUploadID:                           &upload.ID,
			ExcessWeightUpload:                             &upload,
		}

		excessWeightRecord := ExcessWeightRecord(suite.AppContextForTest(), fakeFileStorer, &move)
		suite.Equal(move.ID.String(), excessWeightRecord.MoveID.String())
		suite.Equal(strfmt.DateTime(*move.ExcessWeightQualifiedAt).String(), excessWeightRecord.MoveExcessWeightQualifiedAt.String())
		suite.Equal(strfmt.DateTime(*move.ExcessUnaccompaniedBaggageWeightQualifiedAt).String(), excessWeightRecord.MoveExcessUnaccompaniedBaggageWeightQualifiedAt.String())
		suite.Equal(strfmt.DateTime(*move.ExcessWeightAcknowledgedAt).String(), excessWeightRecord.MoveExcessWeightAcknowledgedAt.String())
		suite.Equal(strfmt.DateTime(*move.ExcessUnaccompaniedBaggageWeightAcknowledgedAt).String(), excessWeightRecord.MoveExcessUnaccompaniedBaggageWeightAcknowledgedAt.String())

		suite.Equal(move.ExcessWeightUploadID.String(), excessWeightRecord.ID.String())
		suite.Equal(move.ExcessWeightUpload.ID.String(), excessWeightRecord.ID.String())
	})

	suite.Run("Success - some nil data, but no errors", func() {
		move := models.Move{ID: id}

		excessWeightRecord := ExcessWeightRecord(suite.AppContextForTest(), fakeFileStorer, &move)
		suite.Equal(move.ID.String(), excessWeightRecord.MoveID.String())
		suite.Nil(excessWeightRecord.MoveExcessWeightQualifiedAt)
		suite.Nil(excessWeightRecord.MoveExcessWeightAcknowledgedAt)
	})
}

func (suite *PayloadsSuite) TestUpload() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	// Get stubbed upload with ID and timestamps
	upload := factory.BuildUpload(nil, []factory.Customization{
		{
			Model: models.Upload{ID: uuid.Must(uuid.NewV4())},
		},
	}, []factory.Trait{factory.GetTraitTimestampedUpload})

	uploadPayload := Upload(suite.AppContextForTest(), fakeFileStorer, &upload)
	suite.Equal(upload.ID.String(), uploadPayload.ID.String())
	suite.Equal(strfmt.DateTime(upload.CreatedAt), uploadPayload.CreatedAt)
	suite.Equal(strfmt.DateTime(upload.UpdatedAt), uploadPayload.UpdatedAt)

	suite.NotEmpty(uploadPayload.URL)
	suite.NotEmpty(uploadPayload.Status)

	suite.Require().NotNil(uploadPayload.Bytes)
	suite.Require().NotNil(uploadPayload.ContentType)
	suite.Require().NotNil(uploadPayload.Filename)
	suite.Equal(upload.Bytes, *uploadPayload.Bytes)
	suite.Equal(upload.ContentType, *uploadPayload.ContentType)
	suite.Equal(upload.Filename, *uploadPayload.Filename)
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

		suite.IsType(&primemessages.SITExtension{}, returnedPayload)
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

		suite.IsType(&primemessages.SITExtension{}, returnedPayload)
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
			WeightRestriction:   models.IntPointer(1000),
			UBWeightRestriction: models.IntPointer(1200),
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
		suite.Equal(int64(1000), *payload.WeightRestriction)
		suite.Equal(int64(1200), *payload.UbWeightRestriction)
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

func (suite *PayloadsSuite) TestNotImplementedError() {
	traceID, _ := uuid.NewV4()
	detail := "Err"

	noDetailError := NotImplementedError(nil, traceID)
	suite.Equal(handlers.NotImplementedErrMessage, *noDetailError.Title)
	suite.Equal(traceID.String(), noDetailError.Instance.String())

	detailError := NotImplementedError(&detail, traceID)
	suite.Equal(handlers.NotImplementedErrMessage, *detailError.Title)
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

func (suite *PayloadsSuite) TestMTOShipmentWithoutServiceItems() {
	// Create the addresses
	pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationType := models.DestinationTypeHomeOfRecord
	secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
	secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
	dlhTestWeight := unit.Pound(4000)
	now := time.Now()

	// Create the MTOShipment with populated PickupAddress and DestinationAddress
	mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PickupAddressID:            &pickupAddress.ID,
				DestinationAddressID:       &destinationAddress.ID,
				DestinationType:            &destinationType,
				SecondaryPickupAddressID:   &secondaryPickupAddress.ID,
				SecondaryDeliveryAddressID: &secondaryDeliveryAddress.ID,
				PrimeEstimatedWeight:       models.PoundPointer(unit.Pound(980)),
				PrimeActualWeight:          &dlhTestWeight,
				NTSRecordedWeight:          models.PoundPointer(unit.Pound(249)),
				TerminatedAt:               &now,
				TerminationComments:        models.StringPointer("get in the choppuh"),
			},
		},
	}, nil)
	shipmentWithoutServiceItem := MTOShipmentWithoutServiceItems(&mtoShipment)
	suite.NotNil(shipmentWithoutServiceItem)
	suite.NotNil(shipmentWithoutServiceItem.TerminatedAt)
	suite.NotNil(shipmentWithoutServiceItem.TerminationComments)
}

func (suite *PayloadsSuite) TestMTOShipmentsWithoutServiceItems() {
	// Create the addresses
	pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationType := models.DestinationTypeHomeOfRecord
	secondaryPickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
	secondaryDeliveryAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
	dlhTestWeight := unit.Pound(4000)

	// Create the MTOShipment
	mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PickupAddressID:            &pickupAddress.ID,
				DestinationAddressID:       &destinationAddress.ID,
				DestinationType:            &destinationType,
				SecondaryPickupAddressID:   &secondaryPickupAddress.ID,
				SecondaryDeliveryAddressID: &secondaryDeliveryAddress.ID,
				PrimeEstimatedWeight:       models.PoundPointer(unit.Pound(980)),
				PrimeActualWeight:          &dlhTestWeight,
				NTSRecordedWeight:          models.PoundPointer(unit.Pound(249)),
			},
		},
	}, nil)

	// Create the MTOShipment
	mtoShipmentTwo := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				PickupAddressID:            &pickupAddress.ID,
				DestinationAddressID:       &destinationAddress.ID,
				DestinationType:            &destinationType,
				SecondaryPickupAddressID:   &secondaryPickupAddress.ID,
				SecondaryDeliveryAddressID: &secondaryDeliveryAddress.ID,
				PrimeEstimatedWeight:       models.PoundPointer(unit.Pound(980)),
				PrimeActualWeight:          &dlhTestWeight,
				NTSRecordedWeight:          models.PoundPointer(unit.Pound(249)),
			},
		},
	}, nil)
	shipmentList := models.MTOShipments{}
	shipmentList = append(shipmentList, mtoShipmentTwo, mtoShipment)
	value := MTOShipmentsWithoutServiceItems(&shipmentList)
	suite.NotNil(value)
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
	destSIT, ok := resultDestSIT.(*primemessages.MTOServiceItemDestSIT)
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
	destSIT, ok := resultDestSIT.(*primemessages.MTOServiceItemInternationalDestSIT)
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

func (suite *PayloadsSuite) TestMTOServiceItemDCRTandDOFSITandDDFSIT() {
	reServiceCode := models.ReServiceCodeDCRT
	reServiceCodeSIT := models.ReServiceCodeDOFSIT
	reServiceCodeDDFSIT := models.ReServiceCodeDDFSIT

	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)

	mtoServiceItemDCRT := &models.MTOServiceItem{
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
	year, month, day := time.Now().Date()
	aWeekAgo := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
	departureDate := aWeekAgo.Add(time.Hour * 24 * 30)
	actualPickupAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
	requestApprovalRequestedStatus := false
	mtoServiceItemDOFSIT := &models.MTOServiceItem{
		ID:                        uuid.Must(uuid.NewV4()),
		ReService:                 models.ReService{Code: reServiceCodeSIT},
		Reason:                    &reason,
		SITDepartureDate:          &departureDate,
		SITEntryDate:              &aWeekAgo,
		SITPostalCode:             models.StringPointer("90210"),
		SITOriginHHGActualAddress: &actualPickupAddress,
		SITCustomerContacted:      &aWeekAgo,
		SITRequestedDelivery:      &aWeekAgo,
		SITOriginHHGOriginalAddress: &models.Address{
			StreetAddress1: "dummyStreet2",
			City:           "dummyCity2",
			State:          "FL",
			PostalCode:     "55555",
		},
		RequestedApprovalsRequestedStatus: &requestApprovalRequestedStatus,
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
	mtoServiceItemDDFSIT := &models.MTOServiceItem{
		ID:                        uuid.Must(uuid.NewV4()),
		ReService:                 models.ReService{Code: reServiceCodeDDFSIT},
		Reason:                    &reason,
		SITDepartureDate:          &departureDate,
		SITEntryDate:              &aWeekAgo,
		SITPostalCode:             models.StringPointer("90210"),
		SITOriginHHGActualAddress: &actualPickupAddress,
		SITCustomerContacted:      &aWeekAgo,
		SITRequestedDelivery:      &aWeekAgo,
		SITOriginHHGOriginalAddress: &models.Address{
			StreetAddress1: "dummyStreet2",
			City:           "dummyCity2",
			State:          "FL",
			PostalCode:     "55555",
		},
		RequestedApprovalsRequestedStatus: &requestApprovalRequestedStatus,
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
	resultDOFSIT := MTOServiceItem(mtoServiceItemDOFSIT)
	resultDDFSIT := MTOServiceItem(mtoServiceItemDDFSIT)

	suite.NotNil(resultDCRT)
	suite.NotNil(resultDOFSIT)
	suite.NotNil(resultDDFSIT)
	_, ok := resultDCRT.(*primemessages.MTOServiceItemDomesticCrating)

	suite.True(ok)
}

func (suite *PayloadsSuite) TestMTOServiceItemICRTandIOFSITandIDFSIT() {
	reServiceCode := models.ReServiceCodeICRT
	reServiceCodeSIT := models.ReServiceCodeIOFSIT
	reServiceCodeIDFSIT := models.ReServiceCodeIDFSIT

	reason := "reason"
	dateOfContact1 := time.Now()
	timeMilitary1 := "1500Z"
	firstAvailableDeliveryDate1 := dateOfContact1.AddDate(0, 0, 10)
	dateOfContact2 := time.Now().AddDate(0, 0, 5)
	timeMilitary2 := "1300Z"
	firstAvailableDeliveryDate2 := dateOfContact2.AddDate(0, 0, 10)

	mtoServiceItemICRT := &models.MTOServiceItem{
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
	year, month, day := time.Now().Date()
	aWeekAgo := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
	departureDate := aWeekAgo.Add(time.Hour * 24 * 30)
	actualPickupAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress2})
	requestApprovalRequestedStatus := false
	mtoServiceItemIOFSIT := &models.MTOServiceItem{
		ID:                        uuid.Must(uuid.NewV4()),
		ReService:                 models.ReService{Code: reServiceCodeSIT},
		Reason:                    &reason,
		SITDepartureDate:          &departureDate,
		SITEntryDate:              &aWeekAgo,
		SITPostalCode:             models.StringPointer("90210"),
		SITOriginHHGActualAddress: &actualPickupAddress,
		SITCustomerContacted:      &aWeekAgo,
		SITRequestedDelivery:      &aWeekAgo,
		SITOriginHHGOriginalAddress: &models.Address{
			StreetAddress1: "dummyStreet2",
			City:           "dummyCity2",
			State:          "FL",
			PostalCode:     "55555",
		},
		RequestedApprovalsRequestedStatus: &requestApprovalRequestedStatus,
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
	mtoServiceItemIDFSIT := &models.MTOServiceItem{
		ID:                        uuid.Must(uuid.NewV4()),
		ReService:                 models.ReService{Code: reServiceCodeIDFSIT},
		Reason:                    &reason,
		SITDepartureDate:          &departureDate,
		SITEntryDate:              &aWeekAgo,
		SITPostalCode:             models.StringPointer("90210"),
		SITOriginHHGActualAddress: &actualPickupAddress,
		SITCustomerContacted:      &aWeekAgo,
		SITRequestedDelivery:      &aWeekAgo,
		SITOriginHHGOriginalAddress: &models.Address{
			StreetAddress1: "dummyStreet2",
			City:           "dummyCity2",
			State:          "FL",
			PostalCode:     "55555",
		},
		RequestedApprovalsRequestedStatus: &requestApprovalRequestedStatus,
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
	resultIOFSIT := MTOServiceItem(mtoServiceItemIOFSIT)
	resultIDFSIT := MTOServiceItem(mtoServiceItemIDFSIT)

	suite.NotNil(resultICRT)
	suite.NotNil(resultIOFSIT)
	suite.NotNil(resultIDFSIT)
	_, ok := resultICRT.(*primemessages.MTOServiceItemInternationalCrating)

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

	_, ok := resultICRT.(*primemessages.MTOServiceItemInternationalCrating)
	suite.True(ok)

	_, ok = resultIUCRT.(*primemessages.MTOServiceItemInternationalCrating)
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
	standaloneCrate := false

	mtoServiceItemDDSHUT := &models.MTOServiceItem{
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

	resultDDSHUT := MTOServiceItem(mtoServiceItemDDSHUT)

	suite.NotNil(resultDDSHUT)

	_, ok := resultDDSHUT.(*primemessages.MTOServiceItemDomesticShuttle)

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

	_, ok := resultIDSHUT.(*primemessages.MTOServiceItemInternationalShuttle)

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

	_, ok := resultIOSHUT.(*primemessages.MTOServiceItemInternationalShuttle)

	suite.True(ok)
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

		suite.IsType(&primemessages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(basicMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(strfmt.DateTime(basicMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotEmpty(returnedModel.ETag)
	})
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

func (suite *PayloadsSuite) TestMTOServiceItemsPOEFSC() {
	portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
		{
			Model: models.Port{
				PortCode: "SEA",
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
	serviceItemPayload := MTOServiceItem(&poefscServiceItem)
	suite.NotNil(serviceItemPayload)
	internationalFuelSurchargeItem, ok := serviceItemPayload.(*primemessages.MTOServiceItemInternationalFuelSurcharge)
	suite.True(ok)
	suite.Equal(portLocation.Port.PortCode, internationalFuelSurchargeItem.PortCode)
	suite.Equal(poefscServiceItem.ReService.Code.String(), internationalFuelSurchargeItem.ReServiceCode)
}

func (suite *PayloadsSuite) TestMTOServiceItemsPODFSC() {
	portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
		{
			Model: models.Port{
				PortCode: "SEA",
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
	serviceItemPayload := MTOServiceItem(&podfscServiceItem)
	suite.NotNil(serviceItemPayload)
	internationalFuelSurchargeItem, ok := serviceItemPayload.(*primemessages.MTOServiceItemInternationalFuelSurcharge)
	suite.True(ok)
	suite.Equal(portLocation.Port.PortCode, internationalFuelSurchargeItem.PortCode)
	suite.Equal(podfscServiceItem.ReService.Code.String(), internationalFuelSurchargeItem.ReServiceCode)
}

func (suite *PayloadsSuite) TestVLocation() {
	suite.Run("correctly maps VLocation with all fields populated", func() {
		city := "LOS ANGELES"
		state := "CA"
		postalCode := "90210"
		county := "LOS ANGELES"
		usPostRegionCityID := uuid.Must(uuid.NewV4())

		vLocation := &models.VLocation{
			CityName:             city,
			StateName:            state,
			UsprZipID:            postalCode,
			UsprcCountyNm:        county,
			UsPostRegionCitiesID: &usPostRegionCityID,
		}

		payload := VLocation(vLocation)

		suite.IsType(payload, &primemessages.VLocation{})
		suite.Equal(handlers.FmtUUID(usPostRegionCityID), &payload.UsPostRegionCitiesID, "Expected UsPostRegionCitiesID to match")
		suite.Equal(city, payload.City, "Expected City to match")
		suite.Equal(state, payload.State, "Expected State to match")
		suite.Equal(postalCode, payload.PostalCode, "Expected PostalCode to match")
		suite.Equal(county, *(payload.County), "Expected County to match")
	})
}

func (suite *PayloadsSuite) TestListMoves() {
	suite.Run("Correctly maps Move Prime Acknowledge At date", func() {
		yesterday := time.Now().AddDate(0, 1, -1)
		address1 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "90210",
				},
			},
		}, nil)
		dutyLocation1 := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model:    address1,
				LinkOnly: true,
			},
		}, nil)
		orders1 := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    dutyLocation1,
				LinkOnly: true,
				Type:     &factory.DutyLocations.NewDutyLocation,
			},
		}, nil)
		move1 := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &yesterday,
				},
			},
			{
				Model:    orders1,
				LinkOnly: true,
			},
		}, nil)

		address2 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "62225",
				},
			},
		}, nil)
		dutyLocation2 := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model:    address2,
				LinkOnly: true,
			},
		}, nil)
		orders2 := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    dutyLocation2,
				LinkOnly: true,
				Type:     &factory.DutyLocations.NewDutyLocation,
			},
		}, nil)
		move2 := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model:    orders2,
				LinkOnly: true,
			},
		}, nil)

		moves := models.Moves{
			move1,
			move2,
		}
		moveOrderAmendmentAvailableSinceCounts := make(services.MoveOrderAmendmentAvailableSinceCounts, 0)

		payload := ListMoves(&moves, suite.AppContextForTest(), moveOrderAmendmentAvailableSinceCounts)
		suite.NotNil(payload)
		suite.Equal(2, len(payload))

		// Move 1
		suite.Equal(strfmt.UUID(move1.ID.String()), payload[0].ID)
		suite.Equal(move1.Locator, payload[0].MoveCode)
		suite.Equal(strfmt.DateTime(move1.CreatedAt), payload[0].CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(move1.AvailableToPrimeAt), payload[0].AvailableToPrimeAt)
		suite.Equal(handlers.FmtDateTimePtr(move1.ApprovedAt), payload[0].ApprovedAt)
		gbloc, err := move1.GetDestinationGBLOC(suite.DB())
		suite.NotNil(gbloc)
		suite.NoError(err)
		suite.Equal(gbloc, payload[0].DestinationGBLOC)
		suite.Equal(move1.Orders.NewDutyLocation.Address.PostalCode, payload[0].DestinationPostalCode)
		suite.Equal(strfmt.UUID(move1.Orders.ID.String()), payload[0].OrderID)
		suite.Equal(*move1.ReferenceID, payload[0].ReferenceID)
		suite.Equal(strfmt.DateTime(move1.UpdatedAt), payload[0].UpdatedAt)
		suite.Equal(etag.GenerateEtag(move1.UpdatedAt), payload[0].ETag)
		suite.Equal(handlers.FmtDateTimePtr(move1.PrimeAcknowledgedAt), payload[0].PrimeAcknowledgedAt)

		// Move 2
		suite.Equal(move2.Locator, payload[1].MoveCode)
		suite.Equal(strfmt.UUID(move2.ID.String()), payload[1].ID)
		suite.Equal(strfmt.DateTime(move2.CreatedAt), payload[1].CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(move2.AvailableToPrimeAt), payload[1].AvailableToPrimeAt)
		suite.Equal(handlers.FmtDateTimePtr(move2.ApprovedAt), payload[1].ApprovedAt)
		gbloc, err = move2.GetDestinationGBLOC(suite.DB())
		suite.NotNil(gbloc)
		suite.NoError(err)
		suite.Equal(gbloc, payload[1].DestinationGBLOC)
		suite.Equal(move2.Orders.NewDutyLocation.Address.PostalCode, payload[1].DestinationPostalCode)
		suite.Equal(strfmt.UUID(move2.Orders.ID.String()), payload[1].OrderID)
		suite.Equal(*move2.ReferenceID, payload[1].ReferenceID)
		suite.Equal(strfmt.DateTime(move2.UpdatedAt), payload[1].UpdatedAt)
		suite.Equal(etag.GenerateEtag(move2.UpdatedAt), payload[1].ETag)
		suite.Nil(payload[1].PrimeAcknowledgedAt)
	})
}

func (suite *PayloadsSuite) TestCountriesPayload() {
	suite.Run("Correctly transform array of countries into payload", func() {
		countries := make([]models.Country, 0)
		countries = append(countries, models.Country{Country: "US", CountryName: "UNITED STATES"})
		payload := Countries(countries)
		suite.True(len(payload) == 1)
		suite.Equal(payload[0].Code, "US")
		suite.Equal(payload[0].Name, "UNITED STATES")
	})

	suite.Run("empty array of countries into payload", func() {
		countries := make([]models.Country, 0)
		payload := Countries(countries)
		suite.True(len(payload) == 0)
	})

	suite.Run("nil countries into payload", func() {
		payload := Countries(nil)
		suite.True(len(payload) == 0)
	})
}

func (suite *PayloadsSuite) TestVIntlLocation() {
	suite.Run("correctly maps VIntlLocation with all fields populated", func() {
		city := "LONDON"
		principalDivision := "CARDIFF"
		intlCityCountriesId := uuid.Must(uuid.NewV4())

		vIntlLocation := &models.VIntlLocation{
			CityName:            &city,
			CountryPrnDivName:   &principalDivision,
			IntlCityCountriesID: &intlCityCountriesId,
		}

		payload := VIntlLocation(vIntlLocation)

		suite.IsType(payload, &primemessages.VIntlLocation{})
		suite.Equal(handlers.FmtUUID(intlCityCountriesId), &payload.IntlCityCountriesID, "Expected IntlCityCountriesID to match")
		suite.Equal(city, payload.City, "Expected City to match")
		suite.Equal(principalDivision, payload.PrincipalDivision, "Expected State to match")
	})
}
