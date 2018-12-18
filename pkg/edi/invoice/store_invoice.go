package ediinvoice

import (
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
	ediTmpFile := ediFilePath + ediFilename

	//createLocalFile(ediFilePath, edi)

	//aFile, err := createAFileFromLocalFile(ediFilePath, ediFilename)

	var aFile = afero.NewOsFs()
	f, err := aFile.Create(ediTmpFile)
	f.WriteString(edi)
	err = f.Sync()

	loader := uploader.NewUploader(nil, logger, *fs)
	err = loader.CreateUploadS3OnlyFromString(userID, edi, &f)

	if verrs.HasAny() {
		logger.Error("Errors encountered for storeEDI():",
			zap.Any("verrors", verrs.Error()))
	}

	if err != nil {
		logger.Error("Errors encountered for storeEDI():",
			zap.Any("err", err.Error()))
	}

	return verrs, err
}
