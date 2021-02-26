package invoice

import (
	"io"
	"path"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// StoreInvoice858C is a service object to store an invoice's EDI in S3.
type StoreInvoice858C struct {
	DB     *pop.Connection
	Logger Logger
	Storer *storage.FileStorer
}

// Call stores the EDI/Invoice to S3.
func (s StoreInvoice858C) Call(edi string, invoice *models.Invoice, userID uuid.UUID) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	// Create path for EDI file
	// {application-bucket}/app/invoice/{invoice_id}.edi
	invoiceID := invoice.ID.String()
	ediFilename := invoiceID + ".edi"
	ediFilePath := "/app/invoice/"
	ediTmpFile := path.Join(ediFilePath, ediFilename)

	fs := afero.NewMemMapFs()

	f, err := fs.Create(ediTmpFile)
	if err != nil {
		return verrs, errors.Wrapf(err, "afero.Create Failed in StoreInvoice858C() invoice ID: %s", invoiceID)
	}
	defer f.Close()

	_, err = io.WriteString(f, edi)
	if err != nil {
		return verrs, errors.Wrapf(err, "io.WriteString(edi) Failed in StoreInvoice858C() invoice ID: %s", invoiceID)
	}

	err = f.Sync()
	if err != nil {
		verrs.Add(validators.GenerateKey("Sync EDI file Failed for file: "+ediTmpFile), err.Error())
	}

	// Create UserUpload
	loader, err := uploader.NewUserUploader(s.DB, s.Logger, *s.Storer, 25*uploader.MB)
	if err != nil {
		s.Logger.Fatal("could not instantiate uploader", zap.Error(err))
	}
	// Set Storagekey path for S3
	loader.SetUploadStorageKey(ediTmpFile)

	// Delete of previous upload, if it exist
	// If Delete of UserUpload fails, ignoring this error because we still have a new UserUpload that needs to be saved
	// to the Invoice
	err = UploadUpdater{DB: s.DB, UserUploader: loader}.DeleteUpload(invoice)
	if err != nil {
		logStr := ""
		if invoice != nil && invoice.UserUploadID != nil {
			logStr = invoice.UserUploadID.String()
		}
		s.Logger.Info("Errors encountered for while deleting previous UserUpload:"+logStr,
			zap.Any("verrors", verrs.Error()))
	}

	// Create and save UserUpload to s3
	userUpload, verrs2, err := loader.CreateUserUpload(userID, uploader.File{File: f}, uploader.AllowedTypesText)
	verrs.Append(verrs2)
	if err != nil {
		return verrs, errors.Wrapf(err, "Failed to Create UserUpload for StoreInvoice858C(), invoice ID: %s", invoiceID)
	}

	if userUpload == nil {
		return verrs, errors.New("Failed to Create and Save new UserUpload object in database, invoice ID: " + invoiceID)
	}

	// Save UserUpload to Invoice
	verrs2, err = UploadUpdater{DB: s.DB, UserUploader: loader}.Call(invoice, userUpload)
	verrs.Append(verrs2)
	if err != nil {
		return verrs, errors.New("Failed to save UserUpload to Invoice: " + invoiceID)
	}

	if verrs.HasAny() {
		s.Logger.Error("Errors encountered for StoreInvoice858C():",
			zap.Any("verrors", verrs.Error()))
	}

	return verrs, err
}
