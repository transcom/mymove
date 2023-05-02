package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *FactorySuite) TestBuildPrimeUpload() {
	suite.Run("Successful creation of default prime upload", func() {
		// Under test:      BuildPrimeUpload
		// Set up:          Create a default prime upload
		// Expected outcome:Create a ProofOfServiceDoc and upload
		// This file doesn't actually exist

		// Create prime upload
		primeUpload := BuildPrimeUpload(suite.DB(), nil, nil)

		suite.NotNil(primeUpload.Contractor)
		suite.False(primeUpload.Contractor.ID.IsNil())
		suite.False(primeUpload.ContractorID.IsNil())

		suite.NotNil(primeUpload.ProofOfServiceDoc)
		suite.False(primeUpload.ProofOfServiceDoc.ID.IsNil())
		suite.False(primeUpload.ProofOfServiceDocID.IsNil())

		suite.NotNil(primeUpload.Upload)
		suite.False(primeUpload.Upload.ID.IsNil())
		suite.False(primeUpload.UploadID.IsNil())
	})

	suite.Run("Successful creation of customized prime upload", func() {
		// Under test:       BuildPrimeUpload
		// Set up:           Create a customized upload (no uploader)
		// Expected outcome: All fields should match

		// Create ProofOfServiceDoc
		customProofOfServiceDoc := models.ProofOfServiceDoc{
			ID: uuid.Must(uuid.NewV4()),
		}

		// Create upload
		customUpload := models.Upload{
			ID: uuid.Must(uuid.NewV4()),
		}

		customPrimeUpload := models.PrimeUpload{
			ID: uuid.Must(uuid.NewV4()),
		}

		// Create prime upload
		primeUpload := BuildPrimeUpload(suite.DB(), []Customization{
			{
				Model: customPrimeUpload,
			},
			{
				Model: customUpload,
				Type:  &Uploads.UploadTypePrime,
			},
			{
				Model: customProofOfServiceDoc,
			},
		}, nil)

		suite.Equal(customPrimeUpload.ID, primeUpload.ID)
		suite.NotNil(primeUpload.ProofOfServiceDoc)
		suite.Equal(customProofOfServiceDoc.ID, primeUpload.ProofOfServiceDoc.ID)
		suite.Equal(customProofOfServiceDoc.ID, primeUpload.ProofOfServiceDocID)
		suite.NotNil(primeUpload.Upload)
		suite.Equal(customUpload.ID, primeUpload.Upload.ID)
		suite.Equal(customUpload.ID, primeUpload.UploadID)
	})

	suite.Run("Successful creation of customized prime upload when contractor already exists", func() {
		// Under test:       BuildPrimeUpload
		// Set up:           Create a customized contractor
		// Expected outcome: Should use the already existing contractor

		contractor := BuildContractor(suite.DB(), []Customization{
			{
				Model: models.Contractor{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)

		customPrimeUpload := models.PrimeUpload{
			ID: uuid.Must(uuid.NewV4()),
		}

		// Create prime upload
		primeUpload := BuildPrimeUpload(suite.DB(), []Customization{
			{
				Model: customPrimeUpload,
			},
		}, nil)
		suite.NotNil(primeUpload.Contractor)
		suite.Equal(contractor.ID, primeUpload.Contractor.ID)
		suite.Equal(contractor.ID, primeUpload.ContractorID)
		suite.Equal(customPrimeUpload.ID, primeUpload.ID)
	})

	suite.Run("Successful creation of customized prime upload with customized orders upload and customized upload for PrimeUpload", func() {
		// Under test:       BuildPrimeUpload
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

		// Create prime upload
		primeUpload := BuildPrimeUpload(suite.DB(), []Customization{
			{
				Model: customUploadForPrime,
				Type:  &Uploads.UploadTypePrime,
			},
			{
				Model: customUploadForUser,
				Type:  &Uploads.UploadTypeUser,
			},
		}, nil)

		suite.Equal(customUploadForPrime.Filename, primeUpload.Upload.Filename)
		suite.Equal(customUploadForPrime.Bytes, primeUpload.Upload.Bytes)
		suite.Equal(customUploadForPrime.ContentType, primeUpload.Upload.ContentType)
		suite.Equal(customUploadForPrime.Checksum, primeUpload.Upload.Checksum)
		suite.Equal(customUploadForPrime.CreatedAt, primeUpload.Upload.CreatedAt)
		suite.False(primeUpload.ProofOfServiceDocID.IsNil())
		suite.False(primeUpload.ProofOfServiceDoc.ID.IsNil())
		suite.False(primeUpload.UploadID.IsNil())
		suite.False(primeUpload.Upload.ID.IsNil())

		ordersUpload := primeUpload.ProofOfServiceDoc.PaymentRequest.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0].Upload
		suite.Equal(customUploadForUser.Filename, ordersUpload.Filename)
		suite.Equal(customUploadForUser.Bytes, ordersUpload.Bytes)
		suite.Equal(customUploadForUser.ContentType, ordersUpload.ContentType)
		suite.Equal(customUploadForUser.Checksum, ordersUpload.Checksum)
		suite.Equal(customUploadForUser.CreatedAt, ordersUpload.CreatedAt)
	})

	suite.Run("Successful creation of prime upload with basic uploader", func() {
		// Under test:      BuildPrimeUpload
		// Mocked:          None
		// Set up:          Create an upload with an uploader and default file
		// Expected outcome:Upload filename should be the default file
		storer := storageTest.NewFakeS3Storage(true)
		primeUploader, err := uploader.NewPrimeUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		defaultFileName := "testdata/test.pdf"
		primeUpload := BuildPrimeUpload(suite.DB(), []Customization{
			{
				Model: models.PrimeUpload{},
				ExtendedParams: &PrimeUploadExtendedParams{
					PrimeUploader: primeUploader,
					AppContext:    suite.AppContextForTest(),
				},
			},
		}, nil)

		upload := primeUpload.Upload

		// no need to test every bit of how the PrimeUploader works
		suite.Contains(upload.Filename, defaultFileName)
		suite.Equal(models.UploadTypePRIME, upload.UploadType)

		// Ensure the associated models are created
		suite.False(primeUpload.ProofOfServiceDocID.IsNil())
		suite.False(primeUpload.ProofOfServiceDoc.ID.IsNil())
		suite.False(primeUpload.ContractorID.IsNil())
		suite.False(primeUpload.UploadID.IsNil())
		suite.False(primeUpload.Upload.ID.IsNil())
	})

	suite.Run("Failed creation of upload - no appcontext", func() {
		// Under test:      BuildPrimeUpload
		// Mocked:          None
		// Set up:          Create a prime upload with a user uploader
		//                  but no appcontext
		// Expected outcome:Should cause a panic
		storer := storageTest.NewFakeS3Storage(true)
		primeUploader, err := uploader.NewPrimeUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		suite.Panics(func() {
			BuildPrimeUpload(suite.DB(), []Customization{
				{
					Model: models.PrimeUpload{},
					ExtendedParams: &PrimeUploadExtendedParams{
						PrimeUploader: primeUploader,
					},
				},
			}, nil)
		})

	})

	suite.Run("Successful creation of customized prime upload with uploader and custom file", func() {
		// Under test:      BuildUserContractor
		// Mocked:          None
		// Set up:          Create a prime upload with a specific file
		// Expected outcome:PrimeUpload should be created with default values
		storer := storageTest.NewFakeS3Storage(true)
		primeUploader, err := uploader.NewPrimeUploader(storer, 100*uploader.MB)
		suite.NoError(err)

		contractor := models.Contractor{
			ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
		}

		// Going down the chain, BuildMove is the first place that a contractor is set.
		// If we want to customize the contractor for Prime Uploads
		// we need to set it here
		paymentRequest := BuildPaymentRequest(suite.DB(), []Customization{
			{
				Model: contractor,
			},
		}, nil)

		uploadFile := "testdata/test.jpg"
		primeUpload := BuildPrimeUpload(suite.DB(), []Customization{
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model: models.PrimeUpload{
					ID: uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-1b557a4e425b"),
				},
				ExtendedParams: &PrimeUploadExtendedParams{
					File:          FixtureOpen("test.jpg"),
					PrimeUploader: primeUploader,
					AppContext:    suite.AppContextForTest(),
				},
			},
		}, nil)

		suite.False(primeUpload.Upload.ID.IsNil())
		suite.Contains(primeUpload.Upload.Filename, uploadFile)
		suite.Equal(models.UploadTypePRIME, primeUpload.Upload.UploadType)
		suite.Equal(paymentRequest.ID, primeUpload.ProofOfServiceDoc.PaymentRequest.ID)
		suite.NotNil(primeUpload.Contractor)

		// Make sure only one contractor is created and used in both the move and prime upload
		suite.Equal(contractor.ID, paymentRequest.MoveTaskOrder.Contractor.ID)
		suite.Equal(contractor.ID, primeUpload.Contractor.ID)
		suite.Equal(contractor.ID, primeUpload.ContractorID)
	})

}
