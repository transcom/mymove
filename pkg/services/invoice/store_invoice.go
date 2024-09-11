package invoice

import (
	"io"
	"path"

	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// StoreInvoice858C is a service object to store an invoice's EDI in S3.
type StoreInvoice858C struct {
}

// Call stores the EDI/Invoice to S3.
func (s StoreInvoice858C) Call(appCtx appcontext.AppContext, edi string, paymentRequest *models.PaymentRequest, fileName string) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	// Create path for EDI file
	// {application-bucket}/app/payment-request-{payment_request_id}/mto-{move_id}/{payment_request_id}.txt
	paymentRequestId := paymentRequest.ID.String()
	ediFilePath := "app/payment-request-uploads/" + "payment-request-" + paymentRequestId + "/mto-" + paymentRequest.MoveTaskOrderID.String() + "/edi"
	ediTmpFile := path.Join(ediFilePath, fileName)
	appCtx.Logger().Warn("attemping to store it to:" + ediFilePath)
	fs := afero.NewMemMapFs()

	f, err := fs.Create(ediTmpFile)
	if err != nil {
		return verrs, errors.Wrapf(err, "afero.Create Failed in StoreInvoice858C() payment request ID: %s", paymentRequestId)
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			appCtx.Logger().Info("Errors encountered while closing file", zap.Error(closeErr))
		}
	}()

	_, err = io.WriteString(f, edi)
	if err != nil {
		return verrs, errors.Wrapf(err, "io.WriteString(edi) Failed in StoreInvoice858C() invoice ID: %s", paymentRequestId)
	}

	err = f.Sync()
	if err != nil {
		verrs.Add(validators.GenerateKey("Sync EDI file Failed for file: "+ediTmpFile), err.Error())
	}

	// Create UserUpload
	loader, err := uploader.NewPaymentRequestEDIUploader(storage.NewFilesystem(storage.FilesystemParams{}), uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
	}
	// Set Storagekey path for S3
	loader.SetUploadStorageKey(ediTmpFile)

	// Create and save UserUpload to s3
	filePointer := &uploader.File{
		File: f,
	}

	// Then pass the pointer to the function
	appUpload, verrs2, err := loader.CreatePaymentRequestEDIUploadForDocument(appCtx, filePointer, uploader.AllowedTypesText)
	verrs.Append(verrs2)
	if err != nil {
		return verrs, errors.Wrapf(err, "Failed to Create AppUpload for StoreInvoice858C(), invoice ID: %s", paymentRequestId)
	}

	if appUpload == nil {
		return verrs, errors.New("Failed to Create and Save new ApUpload object in database, invoice ID: " + paymentRequestId)
	}

	if verrs.HasAny() {
		appCtx.Logger().Error("Errors encountered for StoreInvoice858C():",
			zap.Any("verrors", verrs.Error()))
	}

	return verrs, err
}
