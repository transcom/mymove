package invoice

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/uploader"
)

// UploadUpdater is a service object to invoices adding an UserUpload
type UploadUpdater struct {
	UserUploader *uploader.UserUploader
}

// saveInvoice using DB Transaction
func (u UploadUpdater) saveInvoice(appCtx appcontext.AppContext, invoice *models.Invoice) error {
	if invoice == nil {
		return errors.New("Invoice is nil")
	}

	verrs, err := appCtx.DB().ValidateAndSave(invoice)
	if err != nil || verrs.HasAny() {
		var dbError string
		if err != nil {
			dbError = err.Error()
		}
		if verrs.HasAny() {
			dbError = dbError + verrs.Error()
		}
		return errors.Wrapf(errors.New(dbError), "error saving invoice with ID: %s", invoice.ID.String())
	}
	return nil
}

// DeleteUpload deletes an existing UserUpload
// This function should be called before adding an UserUpload to an Invoice so that the
// UserUpload is removed from the database and from S3 storage before adding a new UserUpload to Invoice
func (u UploadUpdater) DeleteUpload(appCtx appcontext.AppContext, invoice *models.Invoice) error {

	// Check that there is an upload object
	if invoice.UserUpload != nil && invoice.UserUpload.Upload.ID != uuid.Nil {
		if invoice.UserUpload.Upload.StorageKey != "" {

			deleteUploadForUser := invoice.UserUpload

			// Remove association to UserUpload that is to be deleted
			invoice.UserUploadID = nil
			invoice.UserUpload = nil
			err := u.saveInvoice(appCtx, invoice)
			var logString string
			if err != nil {
				logString = fmt.Sprintf("Failed to saveInvoice with UserUploadID: %s", invoice.UserUploadID)
				return errors.Wrap(err, logString)
			}

			// Delete UserUpload
			err = u.UserUploader.DeleteUserUpload(appCtx, deleteUploadForUser)
			if err != nil {
				var storageKey string
				if deleteUploadForUser.Upload.ID != uuid.Nil {
					storageKey = deleteUploadForUser.Upload.StorageKey
				}

				logString = fmt.Sprintf("Failed to DeleteUpload for UserUpload.ID [%s] and StorageKey [%s]", deleteUploadForUser.ID, storageKey)
				return errors.Wrap(err, logString)
			}
		}
	}
	return nil
}

// Call updates the Invoice UserUpload and removes an old UserUpload if present
func (u UploadUpdater) Call(appCtx appcontext.AppContext, invoice *models.Invoice, userUpload *models.UserUpload) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	if userUpload == nil {
		return verrs, errors.New("userUpload is nil")
	}
	if invoice == nil {
		return verrs, errors.New("invoice is nil")
	}

	var err error
	// Save new UserUpload to Invoice
	invoice.UserUpload = userUpload
	invoice.UserUploadID = &userUpload.ID
	err = u.saveInvoice(appCtx, invoice)
	if err != nil {
		return verrs, errors.Wrap(err, "Could not save Invoice for UploadUpdater -- save new userUpload")
	}

	return verrs, nil
}
