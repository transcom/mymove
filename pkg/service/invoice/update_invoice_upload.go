package invoice

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// UpdateInvoiceUpload is a service object to invoices adding an Upload
type UpdateInvoiceUpload struct {
	DB       *pop.Connection
	Uploader *uploader.Uploader
}

// saveInvoice using DB Transaction
func (u UpdateInvoiceUpload) saveInvoice(invoice *models.Invoice) error {
	if invoice == nil {
		return errors.New("Invoice is nil")
	}

	verrs, err := u.DB.ValidateAndSave(invoice)
	if err != nil || verrs.HasAny() {
		var dbError string
		if err != nil {
			dbError = err.Error()
		}
		if verrs.HasAny() {
			dbError = dbError + verrs.Error()
		}
		return errors.Wrapf(err, "error saving invoice with ID: "+invoice.ID.String())
	}
	return nil
}

// DeleteUpload deletes an existing Upload
// This function should be called before adding an Upload to an Invoice so that the
// Upload is removed from the database and from S3 storage before adding a new Upload to Invoice
func (u UpdateInvoiceUpload) DeleteUpload(invoice *models.Invoice) error {

	// Check that there is an upload object
	if invoice.Upload != nil {
		if invoice.Upload.StorageKey != "" {

			deleteUpload := invoice.Upload

			// Remove association to Upload that is to be deleted
			invoice.UploadID = nil
			invoice.Upload = nil
			err := u.saveInvoice(invoice)

			// Delete Upload
			err = u.Uploader.DeleteUpload(deleteUpload)
			var logString string
			if err != nil {
				logString = fmt.Sprintf("Failed to DeleteUpload for Upload.ID [%s] and StorageKey [%s]", deleteUpload.ID, deleteUpload.StorageKey)
				return errors.Wrap(err, logString)
			}
		}
	}
	return nil
}

// Call updates the Invoice Upload and removes an old Upload if present
func (u UpdateInvoiceUpload) Call(invoice *models.Invoice, upload *models.Upload) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	if upload == nil {
		return verrs, errors.New("upload is nil")
	}
	if invoice == nil {
		return verrs, errors.New("invoice is nil")
	}

	var err error
	// Save new Upload to Invoice
	invoice.Upload = upload
	invoice.UploadID = &upload.ID
	err = u.saveInvoice(invoice)
	if err != nil {
		return verrs, errors.Wrap(err, "Could not save Invoice for UpdateInvoiceUpload -- save new upload")
	}

	return verrs, nil
}
