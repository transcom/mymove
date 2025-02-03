package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *FactorySuite) TestBuildServiceRequestDocumentUpload() {
	suite.Run("Successful creation of default service request document upload", func() {
		// Under test:      BuildServiceRequestDocumentUpload
		// Set up:          Create a default service request document upload
		// Expected outcome:Create a ServiceRequestDocument and upload
		// This file doesn't actually exist

		// Create service request document upload
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), nil, nil)

		suite.NotNil(serviceRequestDocumentUpload.Contractor)
		suite.False(serviceRequestDocumentUpload.Contractor.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.ContractorID.IsNil())

		suite.NotNil(serviceRequestDocumentUpload.ServiceRequestDocument)
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocumentID.IsNil())

		suite.NotNil(serviceRequestDocumentUpload.Upload)
		suite.False(serviceRequestDocumentUpload.Upload.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.UploadID.IsNil())
	})

	suite.Run("Successful creation of customized service request document upload", func() {
		// Under test:       BuildServiceRequestDocumentUpload
		// Set up:           Create a customized upload (no uploader)
		// Expected outcome: All fields should match

		// Create ServiceRequestDocument
		customServiceRequestDocument := models.ServiceRequestDocument{
			ID: uuid.Must(uuid.NewV4()),
		}

		// Create upload
		customUpload := models.Upload{
			ID: uuid.Must(uuid.NewV4()),
		}

		customServiceRequestDocumentUpload := models.ServiceRequestDocumentUpload{
			ID: uuid.Must(uuid.NewV4()),
		}

		// Create service request document upload
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{
				Model: customServiceRequestDocumentUpload,
			},
			{
				Model: customUpload,
				Type:  &Uploads.UploadTypePrime,
			},
			{
				Model: customServiceRequestDocument,
			},
		}, nil)

		suite.Equal(customServiceRequestDocumentUpload.ID, serviceRequestDocumentUpload.ID)
		suite.NotNil(serviceRequestDocumentUpload.ServiceRequestDocument)
		suite.Equal(customServiceRequestDocument.ID, serviceRequestDocumentUpload.ServiceRequestDocument.ID)
		suite.Equal(customServiceRequestDocument.ID, serviceRequestDocumentUpload.ServiceRequestDocumentID)
		suite.NotNil(serviceRequestDocumentUpload.Upload)
		suite.Equal(customUpload.ID, serviceRequestDocumentUpload.Upload.ID)
		suite.Equal(customUpload.ID, serviceRequestDocumentUpload.UploadID)
	})

	suite.Run("Successful creation of customized service request document upload when contractor already exists", func() {
		// Under test:       BuildServiceRequestDocumentUpload
		// Set up:           Create a customized contractor
		// Expected outcome: Should use the already existing contractor

		contractor := BuildContractor(suite.DB(), []Customization{
			{
				Model: models.Contractor{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)

		customServiceRequestDocumentUpload := models.ServiceRequestDocumentUpload{
			ID: uuid.Must(uuid.NewV4()),
		}

		// Create service request document upload
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{
				Model: customServiceRequestDocumentUpload,
			},
		}, nil)
		suite.NotNil(serviceRequestDocumentUpload.Contractor)
		suite.Equal(contractor.ID, serviceRequestDocumentUpload.Contractor.ID)
		suite.Equal(contractor.ID, serviceRequestDocumentUpload.ContractorID)
		suite.Equal(customServiceRequestDocumentUpload.ID, serviceRequestDocumentUpload.ID)
	})

	suite.Run("Successful creation of customized service request document upload with customized orders upload and customized upload for ServiceRequestDocumentUpload", func() {
		// Under test:       BuildServiceRequestDocumentUpload
		// Set up:           Create a customized upload with no uploader
		// Expected outcome: All fields should match

		customUploadForPrime := models.Upload{
			Filename:    "BaisWinery.jpg",
			Bytes:       int64(6081979),
			ContentType: "application/jpg",
			Checksum:    "GauMarJosbDHsaQthV5BnQ==",
			CreatedAt:   time.Now(),
			UploadType:  models.UploadTypePRIME,
		}
		customUploadForUser := models.Upload{
			Filename:    "orders.pdf",
			Bytes:       int64(1048576),
			ContentType: "application/pdf",
			Checksum:    "GauMarJosbDHsaQthV5BnQ==",
			CreatedAt:   time.Now(),
		}

		// Create service request document upload
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{
				Model: customUploadForPrime,
				Type:  &Uploads.UploadTypePrime,
			},
			{
				Model: customUploadForUser,
				Type:  &Uploads.UploadTypeUser,
			},
		}, nil)

		suite.Equal(customUploadForPrime.Filename, serviceRequestDocumentUpload.Upload.Filename)
		suite.Equal(customUploadForPrime.Bytes, serviceRequestDocumentUpload.Upload.Bytes)
		suite.Equal(customUploadForPrime.ContentType, serviceRequestDocumentUpload.Upload.ContentType)
		suite.Equal(customUploadForPrime.Checksum, serviceRequestDocumentUpload.Upload.Checksum)
		suite.Equal(customUploadForPrime.CreatedAt, serviceRequestDocumentUpload.Upload.CreatedAt)
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocumentID.IsNil())
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.UploadID.IsNil())
		suite.False(serviceRequestDocumentUpload.Upload.ID.IsNil())

		ordersUpload := serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItem.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0].Upload
		suite.Equal(customUploadForUser.Filename, ordersUpload.Filename)
		suite.Equal(customUploadForUser.Bytes, ordersUpload.Bytes)
		suite.Equal(customUploadForUser.ContentType, ordersUpload.ContentType)
		suite.Equal(customUploadForUser.Checksum, ordersUpload.Checksum)
		suite.Equal(customUploadForUser.CreatedAt, ordersUpload.CreatedAt)
	})

	suite.Run("Successful creation of service request document upload with basic uploader", func() {
		// Under test:      BuildServiceRequestDocumentUpload
		// Mocked:          None
		// Set up:          Create an upload with an uploader and default file
		// Expected outcome:Upload filename should be the default file
		storer := storageTest.NewFakeS3Storage(true)
		serviceRequestDocumentUploader, err := uploader.NewServiceRequestUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		defaultFileName := "testdata/test.pdf"
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{
				Model: models.ServiceRequestDocumentUpload{},
				ExtendedParams: &ServiceRequestDocumentUploadExtendedParams{
					ServiceRequestDocumentUploader: serviceRequestDocumentUploader,
					AppContext:                     suite.AppContextForTest(),
				},
			},
		}, nil)

		upload := serviceRequestDocumentUpload.Upload

		// no need to test every bit of how the ServiceRequestDocumentUploader works
		suite.Contains(upload.Filename, defaultFileName)
		suite.Equal(models.UploadTypePRIME, upload.UploadType)

		// Ensure the associated models are created
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocumentID.IsNil())
		suite.False(serviceRequestDocumentUpload.ServiceRequestDocument.ID.IsNil())
		suite.False(serviceRequestDocumentUpload.ContractorID.IsNil())
		suite.False(serviceRequestDocumentUpload.UploadID.IsNil())
		suite.False(serviceRequestDocumentUpload.Upload.ID.IsNil())
	})

	suite.Run("Successful return of linkOnly ServiceRequestDocumentUpload", func() {
		// Under test:       BuildServiceRequestDocument
		// Set up:           Pass in a linkOnly ServiceRequestDocumentUpload
		// Expected outcome: No new ServiceRequestDocumentUpload should be created
		// Check num ServiceRequestDocument records
		precount, err := suite.DB().Count(&models.ServiceRequestDocumentUpload{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{
				Model: models.ServiceRequestDocumentUpload{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.ServiceRequestDocumentUpload{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, serviceRequestDocumentUpload.ID)
	})

	suite.Run("Failed creation of upload - no appcontext", func() {
		// Under test:      BuildServiceRequestDocumentUpload
		// Mocked:          None
		// Set up:          Create a service request document upload with a user uploader
		//                  but no appcontext
		// Expected outcome:Should cause a panic
		storer := storageTest.NewFakeS3Storage(true)
		serviceRequestDocumentUploader, err := uploader.NewServiceRequestUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		suite.Panics(func() {
			BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
				{
					Model: models.ServiceRequestDocumentUpload{},
					ExtendedParams: &ServiceRequestDocumentUploadExtendedParams{
						ServiceRequestDocumentUploader: serviceRequestDocumentUploader,
					},
				},
			}, nil)
		})

	})

	suite.Run("Successful creation of customized service request document upload with uploader and custom file", func() {
		// Under test:      BuildUserContractor
		// Mocked:          None
		// Set up:          Create a service request document upload with a specific file
		// Expected outcome:ServiceRequestDocumentUpload should be created with default values
		storer := storageTest.NewFakeS3Storage(true)
		serviceRequestDocumentUploader, err := uploader.NewServiceRequestUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		contractor := models.Contractor{
			ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
		}

		// Going down the chain, BuildMove is the first place that a contractor is set.
		// If we want to customize the contractor for Prime Uploads
		// we need to set it here
		paymentRequest := BuildMTOServiceItem(suite.DB(), []Customization{
			{
				Model: contractor,
			},
		}, nil)

		uploadFile := "testdata/test.jpg"
		serviceRequestDocumentUpload := BuildServiceRequestDocumentUpload(suite.DB(), []Customization{
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model: models.ServiceRequestDocumentUpload{
					ID: uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-1b557a4e425b"),
				},
				ExtendedParams: &ServiceRequestDocumentUploadExtendedParams{
					File:                           FixtureOpen("test.jpg"),
					ServiceRequestDocumentUploader: serviceRequestDocumentUploader,
					AppContext:                     suite.AppContextForTest(),
				},
			},
		}, nil)

		suite.False(serviceRequestDocumentUpload.Upload.ID.IsNil())
		suite.Contains(serviceRequestDocumentUpload.Upload.Filename, uploadFile)
		suite.Equal(models.UploadTypePRIME, serviceRequestDocumentUpload.Upload.UploadType)
		suite.Equal(paymentRequest.ID, serviceRequestDocumentUpload.ServiceRequestDocument.MTOServiceItem.ID)
		suite.NotNil(serviceRequestDocumentUpload.Contractor)

		// Make sure only one contractor is created and used in both the move and service request document upload
		suite.Equal(contractor.ID, paymentRequest.MoveTaskOrder.Contractor.ID)
		suite.Equal(contractor.ID, serviceRequestDocumentUpload.Contractor.ID)
		suite.Equal(contractor.ID, serviceRequestDocumentUpload.ContractorID)
	})
}
