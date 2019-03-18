package invoice

import (
	"io"
	"path"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
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

	// Create Upload'r
	loader := uploader.NewUploader(s.DB, s.Logger, *s.Storer)
	// Set Storagekey path for S3
	loader.SetUploadStorageKey(ediTmpFile)

	// Delete of previous upload, if it exist
	// If Delete of Upload fails, ignoring this error because we still have a new Upload that needs to be saved
	// to the Invoice
	err = UpdateInvoiceUpload{DB: s.DB, Uploader: loader}.DeleteUpload(invoice)
	if err != nil {
		logStr := ""
		if invoice != nil && invoice.UploadID != nil {
			logStr = invoice.UploadID.String()
		}
		s.Logger.Info("Errors encountered for while deleting previous Upload:"+logStr,
			zap.Any("verrors", verrs.Error()))
	}

	// Create and save Upload to s3
	upload, verrs2, err := loader.CreateUpload(userID, &f, uploader.AllowedTypesText)
	verrs.Append(verrs2)
	if err != nil {
		return verrs, errors.Wrapf(err, "Failed to Create Upload for StoreInvoice858C(), invoice ID: %s", invoiceID)
	}

	if upload == nil {
		return verrs, errors.New("Failed to Create and Save new Upload object in database, invoice ID: " + invoiceID)
	}

	// Save Upload to Invoice
	verrs2, err = UpdateInvoiceUpload{DB: s.DB, Uploader: loader}.Call(invoice, upload)
	verrs.Append(verrs2)
	if err != nil {
		return verrs, errors.New("Failed to save Upload to Invoice: " + invoiceID)
	}

	if verrs.HasAny() {
		s.Logger.Error("Errors encountered for StoreInvoice858C():",
			zap.Any("verrors", verrs.Error()))
	}

	return verrs, err
}
