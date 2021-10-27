package payloads

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestMoveTaskOrder() {
	moveTaskOrderID, _ := uuid.NewV4()
	ordersID, _ := uuid.NewV4()
	referenceID := "testID"
	primeTime := time.Now()
	submittedAt := time.Now()
	hhgMoveType := models.SelectedMoveTypeHHG
	excessWeightQualifiedAt := time.Now()
	excessWeightAcknowledgedAt := time.Now()
	excessWeightUploadID := uuid.Must(uuid.NewV4())

	basicMove := models.Move{
		ID:                         moveTaskOrderID,
		Locator:                    "TESTTEST",
		CreatedAt:                  time.Now(),
		AvailableToPrimeAt:         &primeTime,
		OrdersID:                   ordersID,
		Orders:                     models.Order{},
		ReferenceID:                &referenceID,
		PaymentRequests:            models.PaymentRequests{},
		SubmittedAt:                &submittedAt,
		UpdatedAt:                  time.Now(),
		SelectedMoveType:           &hhgMoveType,
		PersonallyProcuredMoves:    models.PersonallyProcuredMoves{},
		MoveDocuments:              models.MoveDocuments{},
		Status:                     models.MoveStatusAPPROVED,
		SignedCertifications:       models.SignedCertifications{},
		MTOServiceItems:            models.MTOServiceItems{},
		MTOShipments:               models.MTOShipments{},
		ExcessWeightQualifiedAt:    &excessWeightQualifiedAt,
		ExcessWeightAcknowledgedAt: &excessWeightAcknowledgedAt,
		ExcessWeightUploadID:       &excessWeightUploadID,
	}

	suite.T().Run("Success - Returns a basic move payload with no payment requests, service items or shipments", func(t *testing.T) {
		returnedModel := MoveTaskOrder(&basicMove)

		suite.IsType(&primemessages.MoveTaskOrder{}, returnedModel)
		suite.Equal(strfmt.UUID(basicMove.ID.String()), returnedModel.ID)
		suite.Equal(basicMove.Locator, returnedModel.MoveCode)
		suite.Equal(strfmt.DateTime(basicMove.CreatedAt), returnedModel.CreatedAt)
		suite.Equal(handlers.FmtDateTimePtr(basicMove.AvailableToPrimeAt), returnedModel.AvailableToPrimeAt)
		suite.Equal(strfmt.UUID(basicMove.OrdersID.String()), returnedModel.OrderID)
		suite.Equal(referenceID, returnedModel.ReferenceID)
		suite.Equal(strfmt.DateTime(basicMove.UpdatedAt), returnedModel.UpdatedAt)
		suite.NotEmpty(returnedModel.ETag)
		suite.True(returnedModel.ExcessWeightQualifiedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightQualifiedAt)))
		suite.True(returnedModel.ExcessWeightAcknowledgedAt.Equal(strfmt.DateTime(*basicMove.ExcessWeightAcknowledgedAt)))
		suite.Require().NotNil(returnedModel.ExcessWeightUploadID)
		suite.Equal(strfmt.UUID(basicMove.ExcessWeightUploadID.String()), *returnedModel.ExcessWeightUploadID)
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

	suite.T().Run("Success - Returns a reweigh payload without optional fields", func(t *testing.T) {
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

	suite.T().Run("Success - Returns a reweigh payload with optional fields", func(t *testing.T) {
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
	upload := testdatagen.MakeStubbedUpload(suite.DB(), testdatagen.Assertions{})

	suite.T().Run("Success - all data populated", func(t *testing.T) {
		move := models.Move{
			ID:                         id,
			ExcessWeightQualifiedAt:    &now,
			ExcessWeightAcknowledgedAt: &now,
			ExcessWeightUploadID:       &upload.ID,
			ExcessWeightUpload:         &upload,
		}

		excessWeightRecord := ExcessWeightRecord(suite.TestAppContext(), fakeFileStorer, &move)
		suite.Equal(move.ID.String(), excessWeightRecord.MoveID.String())
		suite.Equal(strfmt.DateTime(*move.ExcessWeightQualifiedAt).String(), excessWeightRecord.MoveExcessWeightQualifiedAt.String())
		suite.Equal(strfmt.DateTime(*move.ExcessWeightAcknowledgedAt).String(), excessWeightRecord.MoveExcessWeightAcknowledgedAt.String())

		suite.Equal(move.ExcessWeightUploadID.String(), excessWeightRecord.ID.String())
		suite.Equal(move.ExcessWeightUpload.ID.String(), excessWeightRecord.ID.String())
	})

	suite.T().Run("Success - some nil data, but no errors", func(t *testing.T) {
		move := models.Move{ID: id}

		excessWeightRecord := ExcessWeightRecord(suite.TestAppContext(), fakeFileStorer, &move)
		suite.Equal(move.ID.String(), excessWeightRecord.MoveID.String())
		suite.Nil(excessWeightRecord.MoveExcessWeightQualifiedAt)
		suite.Nil(excessWeightRecord.MoveExcessWeightAcknowledgedAt)
	})
}

func (suite *PayloadsSuite) TestUpload() {
	fakeFileStorer := test.NewFakeS3Storage(true)
	upload := testdatagen.MakeStubbedUpload(suite.DB(), testdatagen.Assertions{})

	uploadPayload := Upload(suite.TestAppContext(), fakeFileStorer, &upload)
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

	sitExtension := models.SITExtension{
		ID:            id,
		MTOShipmentID: shipmentID,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		RequestedDays: int(30),
		RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		Status:        models.SITExtensionStatusPending,
	}

	suite.T().Run("Success - Returns a sitextension payload without optional fields", func(t *testing.T) {
		returnedPayload := SITExtension(&sitExtension)

		suite.IsType(&primemessages.SITExtension{}, returnedPayload)
		suite.Equal(strfmt.UUID(sitExtension.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(sitExtension.MTOShipmentID.String()), returnedPayload.MtoShipmentID)
		suite.Equal(strfmt.DateTime(sitExtension.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(sitExtension.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Nil(returnedPayload.ApprovedDays)
		suite.Nil(returnedPayload.ContractorRemarks)
		suite.Nil(returnedPayload.OfficeRemarks)
		suite.Nil(returnedPayload.DecisionDate)
		suite.NotNil(returnedPayload.ETag)

	})

	suite.T().Run("Success - Returns a sit extension payload with optional fields", func(t *testing.T) {
		// Set optional fields
		approvedDays := int(30)
		sitExtension.ApprovedDays = &approvedDays

		contractorRemarks := "some reason"
		sitExtension.ContractorRemarks = &contractorRemarks

		officeRemarks := "some other reason"
		sitExtension.OfficeRemarks = &officeRemarks

		decisionDate := time.Now()
		sitExtension.DecisionDate = &decisionDate

		// Send model through func
		returnedPayload := SITExtension(&sitExtension)

		suite.IsType(&primemessages.SITExtension{}, returnedPayload)
		suite.Equal(strfmt.UUID(sitExtension.ID.String()), returnedPayload.ID)
		suite.Equal(strfmt.UUID(sitExtension.MTOShipmentID.String()), returnedPayload.MtoShipmentID)
		suite.Equal(strfmt.DateTime(sitExtension.CreatedAt), returnedPayload.CreatedAt)
		suite.Equal(strfmt.DateTime(sitExtension.UpdatedAt), returnedPayload.UpdatedAt)
		suite.Equal(handlers.FmtIntPtrToInt64(sitExtension.ApprovedDays), returnedPayload.ApprovedDays)
		suite.Equal(sitExtension.ContractorRemarks, returnedPayload.ContractorRemarks)
		suite.Equal(sitExtension.OfficeRemarks, returnedPayload.OfficeRemarks)
		suite.Equal((*strfmt.DateTime)(sitExtension.DecisionDate), returnedPayload.DecisionDate)
		suite.NotNil(returnedPayload.ETag)

	})
}
