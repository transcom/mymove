package invoice

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// UpdateInvoiceUpload is a service object to invoices adding an Upload
type UpdateInvoiceUpload struct {
	DB       *pop.Connection
	Uploader *uploader.Uploader
}

// deleteUpload deletes an existing Upload
// This function should be called before adding and Upload to an Invoice so that the
// Upload is removed from the database and from S3 before adding a new Upload to Invoice
func (u UpdateInvoiceUpload) deleteUpload(upload *models.Upload) error {
	if upload != nil {
		if upload.StorageKey != "" {
			err := u.Uploader.DeleteUpload(upload)
			var logString string
			if err != nil {
				logString = fmt.Sprintf("Failed to DeleteUpload for Upload.ID [%s] and StorageKey [%s]", upload.ID, upload.StorageKey)
				return errors.New(logString)
			}
		}
	}
	return nil
}

// Call updates the Invoice Upload and removes an old Upload if present
func (u UpdateInvoiceUpload) Call(invoice *models.Invoice, upload *models.Upload) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	var err error
	transactionErr := u.DB.Transaction(func(connection *pop.Connection) error {
		err = u.deleteUpload(invoice.Upload)
		invoice.UploadID = nil
		invoice.Upload = nil
		if err != nil {
			verrs.Add(validators.GenerateKey("DeleteUpload"), err.Error())
		}
		invoice.Upload = upload
		invoice.UploadID = &upload.ID
		// Sample code of what eager creation should like
		// https://gobuffalo.io/en/docs/db/relations#eager-creation
		// verrs, err := c.db.Eager().ValidateAndSave(&invoice)
		verrs2, err := u.DB.ValidateAndSave(invoice)
		if err != nil || verrs2.HasAny() {
			return errors.New("error saving invoice")
		}
		return nil
	})
	if transactionErr != nil {
		return verrs, err
	}
	return verrs, nil
}
