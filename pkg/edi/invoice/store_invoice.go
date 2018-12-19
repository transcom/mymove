package ediinvoice

import (
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
	"go.uber.org/zap"
	"io"
)

// StoreInvoice858C stores the EDI/Invoice to S3
func StoreInvoice858C(edi string, invoiceID uuid.UUID, storer *storage.FileStorer, logger *zap.Logger, userID uuid.UUID) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	// Create path for EDI file
	// {application-bucket}/app/invoice/{invoice_id}.edi
	ediFilename := invoiceID.String() + ".edi"
	ediFilePath := "/app/invoice/"
	ediTmpFile := ediFilePath + ediFilename

	fs := afero.NewMemMapFs()

	f, err := fs.Create(ediTmpFile)
	if err != nil {
		return verrs, errors.Wrap(err, "afero.Create Failed in StoreInvoice858C()")
	}
	defer f.Close()

	_, err = io.WriteString(f, edi)
	if err != nil {
		return verrs, errors.Wrap(err, "io.WriteString(edi) Failed in StoreInvoice858C()")
	}

	err = f.Sync()
	if err != nil {
		verrs.Add(validators.GenerateKey("Sync EDI file"), err.Error())
	}

	loader := uploader.NewUploader(nil, logger, *storer)
	_, err = loader.CreateUploadS3Only(userID, &f)

	if verrs.HasAny() {
		logger.Error("Errors encountered for StoreInvoice858C():",
			zap.Any("verrors", verrs.Error()))
	}

	if err != nil {
		logger.Error("Errors encountered for storStoreInvoice858CeEDI():",
			zap.Any("err", err.Error()))
	}

	return verrs, err
}
