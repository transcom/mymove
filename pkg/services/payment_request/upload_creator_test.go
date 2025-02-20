// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to clean up file created for unit test
// RA: Given the functions causing the lint errors are used to clean up local storage space after a unit test, it does not present a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package paymentrequest

import (
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PaymentRequestServiceSuite) TestCreateUploadSuccess() {
	var fakeS3 *test.FakeS3Storage
	var contractor models.Contractor
	var testFile *os.File
	var paymentRequest models.PaymentRequest

	moveTaskOrderID, err := uuid.FromString("cc4523e2-e418-48cc-804e-57a507fff093")
	suite.NoError(err)

	setupTestData := func() {
		contractor = factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

		fakeS3 = test.NewFakeS3Storage(true)
		paymentRequestID, err := uuid.FromString("9b873071-149f-43c2-8971-e93348ebc5e3")
		suite.NoError(err)

		moveTaskOrder := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ID: moveTaskOrderID,
				},
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					ID: paymentRequestID,
				},
			},
		}, nil)

		testFile, err = os.Open("../../testdatagen/testdata/test.pdf")
		suite.NoError(err)
	}

	suite.Run("PrimeUpload is created successfully", func() {
		setupTestData()
		uploadCreator := NewPaymentRequestUploadCreator(fakeS3)
		upload, err := uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, paymentRequest.ID, contractor.ID, "unit-test-file.pdf", false)

		expectedFilename := fmt.Sprintf("/payment-request-uploads/mto-%s/payment-request-%s", moveTaskOrderID, paymentRequest.ID)
		suite.NoError(err)
		suite.Contains(upload.Filename, expectedFilename)
		suite.Equal(int64(10596), upload.Bytes)
		suite.Equal(uploader.FileTypePDF, upload.ContentType)

		var proofOfServiceDoc models.ProofOfServiceDoc
		proofOfServiceDocExists, err := suite.DB().Q().
			LeftJoin("payment_requests pr", "pr.id = proof_of_service_docs.payment_request_id").
			LeftJoin("prime_uploads pu", "proof_of_service_docs.id = pu.proof_of_service_docs_id").
			LeftJoin("uploads u", "pu.upload_id = u.id").
			Where("u.id = $1", upload.ID).Where("pr.id = $2", paymentRequest.ID).
			Eager("PrimeUploads.Upload").
			Exists(&proofOfServiceDoc)
		suite.NoError(err)
		suite.Equal(true, proofOfServiceDocExists)
	})

	testFile.Close()
}

func (suite *PaymentRequestServiceSuite) TestCreateUploadFailure() {
	var contractor models.Contractor

	fakeS3 := test.NewFakeS3Storage(true)

	setupTestData := func() {
		contractor = factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
		factory.BuildPaymentRequest(suite.DB(), nil, nil)
	}

	suite.Run("invalid payment request ID", func() {
		setupTestData()
		testFile, err := os.Open("../../testdatagen/testdata/test.pdf")
		suite.NoError(err)
		defer func() {
			if closeErr := testFile.Close(); closeErr != nil {
				suite.T().Error("Failed to close file", zap.Error(closeErr))
			}
		}()

		uploadCreator := NewPaymentRequestUploadCreator(fakeS3)
		_, err = uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, uuid.FromStringOrNil("96b77644-4028-48c2-9ab8-754f33309db9"), contractor.ID, "unit-test-file.pdf", false)
		suite.Error(err)
	})

	suite.Run("invalid user ID", func() {
		setupTestData()
		testFile, err := os.Open("../../testdatagen/testdata/test.pdf")
		suite.NoError(err)
		defer func() {
			if closeErr := testFile.Close(); closeErr != nil {
				suite.T().Error("Failed to close file", zap.Error(closeErr))
			}
		}()

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		uploadCreator := NewPaymentRequestUploadCreator(fakeS3)
		_, err = uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, paymentRequest.ID, uuid.FromStringOrNil("806e2f96-f9f9-4cbb-9a3d-d2f488539a1f"), "unit-test-file.pdf", false)
		suite.Error(err)
	})

	suite.Run("invalid file type", func() {
		setupTestData()
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		uploadCreator := NewPaymentRequestUploadCreator(fakeS3)
		wrongTypeFile, err := os.Open("../../testdatagen/testdata/test.txt")
		suite.NoError(err)

		defer func() {
			if closeErr := wrongTypeFile.Close(); closeErr != nil {
				suite.T().Error("Failed to close file", zap.Error(closeErr))
			}
		}()

		_, err = uploadCreator.CreateUpload(suite.AppContextForTest(), wrongTypeFile, paymentRequest.ID, contractor.ID, "unit-test-file.pdf", false)
		suite.Error(err)
	})

}
