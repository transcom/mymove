package invoice

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *InvoiceServiceSuite) openLocalFile(path string) (afero.File, error) {
	var fs = afero.NewMemMapFs()

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		suite.logger.Fatal("Error opening local file", zap.Error(err))
	}

	suite.NotNil(fs)
	outputFile, err := fs.Create(path)
	if err != nil {
		suite.logger.Fatal("Error creating afero file", zap.Error(err))
	}

	_, err = io.Copy(outputFile, file)
	if err != nil {
		suite.logger.Fatal("Error copying to afero file", zap.Error(err))
	}

	return outputFile, nil
}

func (suite *InvoiceServiceSuite) helperCreateUpload(storer *storage.FileStorer) *models.Upload {
	document := testdatagen.MakeDefaultDocument(suite.DB())
	userID := document.ServiceMember.UserID
	up, err := uploader.NewUploader(suite.DB(), suite.logger, *storer, 25*uploader.MB)
	suite.NoError(err)

	// Create file to use for upload
	testFile, err := os.Open("../../testdatagen/testdata/test.pdf")
	suite.NoError(err)

	// Create Upload and save it
	upload, verrs, err := up.CreateUpload(userID, uploader.File{File: testFile}, uploader.AllowedTypesPDF)
	suite.Nil(err, "CreateUpload() failed to create upload")
	suite.Empty(verrs.Error(), "CreateUpload() verrs returned error")
	suite.NotNil(upload, "CreateUpload() failed to create upload structure")
	if upload == nil {
		suite.T().Fatalf("failed to create a upload object: %s", err)
	}
	// Call Close on file after CreateUploadForDocument is complete
	testFile.Close()
	return upload
}

func (suite *InvoiceServiceSuite) helperCreateFileStorer() *storage.FileStorer {
	var storer storage.FileStorer
	fakeS3 := storageTest.NewFakeS3Storage(true)
	storer = fakeS3
	return &storer
}

func (suite *InvoiceServiceSuite) helperFetchInvoice(invoiceID uuid.UUID) (*models.Invoice, error) {
	var invoice models.Invoice
	err := suite.DB().Eager().Find(&invoice, invoiceID)
	if err != nil {
		fmt.Print(err.Error())
		if errors.Cause(err).Error() == "sql: no rows in result set" {
			return nil, errors.New("Record not found")
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	return &invoice, nil
}
