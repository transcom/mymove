// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used to clean up file created for unit test
// RA: Given the functions causing the lint errors are used to clean up local storage space after a unit test, it does not present a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package mtoserviceitem

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *MTOServiceItemServiceSuite) TestCreateUploadSuccess() {
	var fakeS3 *test.FakeS3Storage
	var contractor models.Contractor
	var testFile *os.File
	var mtoServiceItem models.MTOServiceItem

	setupTestData := func() {
		contractor = factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

		fakeS3 = test.NewFakeS3Storage(true)
		mtoServiceItemID := uuid.Must(uuid.NewV4())

		moveTaskOrder := factory.BuildMove(suite.DB(), nil, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    moveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					ID: mtoServiceItemID,
				},
			},
		}, nil)
		var err error
		testFile, err = os.Open("../../testdatagen/testdata/test.pdf")
		suite.NoError(err)
	}

	suite.Run("ServiceRequestDocumentUpload is created successfully", func() {
		setupTestData()
		uploadCreator := NewServiceRequestDocumentUploadCreator(fakeS3)
		upload, err := uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, mtoServiceItem.ID, contractor.ID, "unit-test-file.pdf")
		suite.NoError(err)

		found := regexp.MustCompile(fmt.Sprintf(`/mto-service-item/%s/unit-test-file-\d{14}.pdf`, mtoServiceItem.ID)).FindString(upload.Filename)
		suite.NotEmpty(found, "Regex must match filename: %s", upload.Filename)

		suite.Equal(int64(10596), upload.Bytes)
		suite.Equal(uploader.FileTypePDF, upload.ContentType)

		var serviceRequestDocument models.ServiceRequestDocument
		serviceRequestDocumentExists, err := suite.DB().Q().
			LeftJoin("mto_service_items si", "si.id = service_request_documents.mto_service_item_id").
			LeftJoin("service_request_document_uploads sr", "service_request_documents.id = sr.service_request_documents_id").
			LeftJoin("uploads u", "sr.upload_id = u.id").
			Where("u.id = $1", upload.ID).Where("si.id = $2", mtoServiceItem.ID).
			Eager("ServiceRequestDocumentUploads.Upload").
			Exists(&serviceRequestDocument)
		suite.NoError(err)
		suite.Equal(true, serviceRequestDocumentExists)
	})

	testFile.Close()
}

func (suite *MTOServiceItemServiceSuite) TestCreateServiceRequestUploadFailure() {
	var contractor models.Contractor

	fakeS3 := test.NewFakeS3Storage(true)

	setupTestData := func() {
		contractor = factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
		factory.BuildMTOServiceItem(suite.DB(), nil, nil)
	}

	suite.Run("invalid mto service item ID", func() {
		setupTestData()
		testFile, err := os.Open("../../testdatagen/testdata/test.pdf")
		suite.NoError(err)
		defer func() {
			if closeErr := testFile.Close(); closeErr != nil {
				suite.T().Error("Failed to close file", zap.Error(closeErr))
			}
		}()

		uploadCreator := NewServiceRequestDocumentUploadCreator(fakeS3)
		_, err = uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, uuid.FromStringOrNil("96b77644-4028-48c2-9ab8-754f33309db9"), contractor.ID, "unit-test-file.pdf")
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

		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)
		uploadCreator := NewServiceRequestDocumentUploadCreator(fakeS3)
		_, err = uploadCreator.CreateUpload(suite.AppContextForTest(), testFile, mtoServiceItem.ID, uuid.FromStringOrNil("806e2f96-f9f9-4cbb-9a3d-d2f488539a1f"), "unit-test-file.pdf")
		suite.Error(err)
	})

	suite.Run("invalid file type", func() {
		setupTestData()
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)
		uploadCreator := NewServiceRequestDocumentUploadCreator(fakeS3)
		wrongTypeFile, err := os.Open("../../testdatagen/testdata/test.txt")
		suite.NoError(err)

		defer func() {
			if closeErr := wrongTypeFile.Close(); closeErr != nil {
				suite.T().Error("Failed to close file", zap.Error(closeErr))
			}
		}()

		_, err = uploadCreator.CreateUpload(suite.AppContextForTest(), wrongTypeFile, mtoServiceItem.ID, contractor.ID, "unit-test-file.pdf")
		suite.Error(err)
	})

}
