package ediinvoice

import (
	"fmt"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
	"go.uber.org/zap"
	"os"
)

func createLocalFile(path string, data string) error {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	// close f on exit and check for its returned error
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	numBytes, err := f.WriteString(data)
	if err != nil {
		panic(err)
	}

	if numBytes == 0 {
		// log/return error
	}

	f.Sync()

	return nil
}

/*
func createAFileFromLocalFile(path string, filename string) (*afero.File, error) {
	newFilename := filepath.Join(path, filepath.Clean(filename))
	file, err := os.Open(newFilename)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	var fs afero.Afero
	outputFile, err := fs.Create(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error creating afero file")
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		return nil, errors.Wrap(err, "error copying over file contents")
	}

	return &outputFile, nil
}
*/

// StoreInvoice858C stores the EDI/Invoice to S3
func StoreInvoice858C(edi string, invoiceID uuid.UUID, fs *storage.FileStorer, logger *zap.Logger, userID uuid.UUID) (*validate.Errors, error) {

	//userID = session.UserID
	verrs := validate.NewErrors()

	// Create path for EDI file
	// {application-bucket}/app/invoice/{invoice_id}.edi
	ediFilename := invoiceID.String() + ".edi"
	ediFilePath := "/tmp/milmoves/app/invoice/"
	// TODO: what is the correct file path
	//ediFilePath := "/app/invoice/"
	ediTmpFile := ediFilePath + ediFilename

	//createLocalFile(ediFilePath, edi)

	//aFile, err := createAFileFromLocalFile(ediFilePath, ediFilename)

	var aFile = afero.NewOsFs()

	pathExist, err := afero.DirExists(aFile, ediFilePath)
	if err == nil {
		if !pathExist {
			modeType := os.ModeDir | os.ModeTemporary
			err = aFile.MkdirAll(ediFilePath, modeType)
			if err != nil {
				verrs.Add(validators.GenerateKey("MkdirAll"), err.Error())
				return verrs, errors.Errorf("ERROR: Could not create dir for StoreInvoice858C dir: %s", ediFilePath)
			}
		}
	}

	f, err := aFile.Create(ediTmpFile)
	f.WriteString(edi)
	err = f.Sync()

	loader := uploader.NewUploader(nil, logger, *fs)
	_, err = loader.CreateUploadS3OnlyFromString(userID, edi, &f)

	// Remove temp EDI/Invoice file from local filesystem after uploading to S3
	err = aFile.Remove(ediTmpFile)
	if err != nil {
		verrs.Add(validators.GenerateKey("Remove EDI File"),
			fmt.Sprintf("Failed to remove file: %s", ediTmpFile))
	}

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
