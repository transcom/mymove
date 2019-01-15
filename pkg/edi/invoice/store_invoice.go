package ediinvoice

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
	invoiceop "github.com/transcom/mymove/pkg/service/invoice"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// StoreInvoice858C stores the EDI/Invoice to S3
func StoreInvoice858C(edi string, invoice *models.Invoice, storer *storage.FileStorer,
	logger *zap.Logger, userID uuid.UUID, db *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	// Create path for EDI file
	// {application-bucket}/app/invoice/{invoice_id}.edi
	invoiceNumber := invoice.ID.String()
	ediFilename := invoiceNumber + ".edi"
	ediFilePath := "/app/invoice/"
	ediTmpFile := path.Join(ediFilePath, ediFilename)

	fs := afero.NewMemMapFs()

	f, err := fs.Create(ediTmpFile)
	if err != nil {
		return verrs, errors.Wrapf(err, "afero.Create Failed in StoreInvoice858C() invoice number: %s", invoiceNumber)
	}
	defer f.Close()

	_, err = io.WriteString(f, edi)
	if err != nil {
		return verrs, errors.Wrapf(err, "io.WriteString(edi) Failed in StoreInvoice858C() invoice number: %s", invoiceNumber)
	}

	err = f.Sync()
	if err != nil {
		verrs.Add(validators.GenerateKey("Sync EDI file Failed for file: "+ediTmpFile), err.Error())
	}

	// Create Upload'r
	loader := uploader.NewUploader(db, logger, *storer)

	// Delete of previous upload, if it exist
	// If Delete of Upload fails, ignoring this error because we still have a new Upload that needs to be saved
	// to the Invoice
	err = invoiceop.UpdateInvoiceUpload{DB: db, Uploader: loader}.DeleteUpload(invoice)
	if err != nil {
		logStr := ""
		if invoice != nil && invoice.UploadID != nil {
			logStr = invoice.UploadID.String()
		}
		logger.Info("Errors encountered for while deleting previous Upload:"+logStr,
			zap.Any("verrors", verrs.Error()))
	}

	// Create and save Upload to s3
	upload, verrs2, err := loader.CreateUploadNoDocument(userID, &f)
	verrs.Append(verrs2)
	if err != nil {
		return verrs, errors.Wrapf(err, "Failed to Create Upload for StoreInvoice858C(), invoice number: %s", invoiceNumber)
	}

	if upload == nil {
		return verrs, errors.New("Failed to Create and Save new Upload object in database, invoice number: " + invoiceNumber)
	}

	// Save Upload to Invoice
	verrs2, err = invoiceop.UpdateInvoiceUpload{DB: db, Uploader: loader}.Call(invoice, upload)
	verrs.Append(verrs2)
	if err != nil {
		return verrs, errors.New("Failed to save Upload to Invoice: " + invoiceNumber)
	}

	if verrs.HasAny() {
		logger.Error("Errors encountered for StoreInvoice858C():",
			zap.Any("verrors", verrs.Error()))
	}

	return verrs, err
}
