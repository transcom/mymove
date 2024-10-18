package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCreatePaymentRequestEdiFile() {
	suite.Run("test create PaymentRequest Edi file db save", func() {
		fileRecord := models.PaymentRequestEdiFile{
			ID:                   uuid.Must(uuid.NewV4()),
			Filename:             "858.file.txt",
			EdiString:            "Sample EDI content",
			PaymentRequestNumber: "1111-2222-1",
		}
		suite.NoError(models.CreatePaymentRequestEdiFile(suite.DB(), fileRecord.Filename, fileRecord.EdiString, fileRecord.PaymentRequestNumber))
	})
	suite.Run("test create PaymentRequest Edi file db save", func() {
		fileRecord := models.PaymentRequestEdiFile{
			ID:                   uuid.Must(uuid.NewV4()),
			EdiString:            "Sample EDI content",
			PaymentRequestNumber: "1111-2222-1",
		}
		suite.Nil(models.CreatePaymentRequestEdiFile(suite.DB(), fileRecord.Filename, fileRecord.EdiString, fileRecord.PaymentRequestNumber))
	})
	suite.Run("TestFetchAllPaymentRequestEdiFiles", func() {
		suite.T().Run("successfully fetches all PaymentRequestEdiFiles", func(t *testing.T) {
			fileRecord := models.PaymentRequestEdiFile{
				ID:                   uuid.Must(uuid.NewV4()),
				Filename:             "858.file.txt",
				EdiString:            "Sample EDI content",
				PaymentRequestNumber: "1111-2222-1",
			}
			err := models.CreatePaymentRequestEdiFile(suite.DB(), fileRecord.Filename, fileRecord.EdiString, fileRecord.PaymentRequestNumber)
			suite.NoError(err)

			files, err := models.FetchAllPaymentRequestEdiFiles(suite.DB())
			suite.NoError(err)
			suite.NotEmpty(files)
		})
	})
	suite.Run("TestFetchPaymentRequestEdiByPaymentRequestNumber", func() {
		fileRecord := models.PaymentRequestEdiFile{
			ID:                   uuid.Must(uuid.NewV4()),
			Filename:             "858.file.txt",
			EdiString:            "Sample EDI content",
			PaymentRequestNumber: "1111-2222-1",
		}
		err := models.CreatePaymentRequestEdiFile(suite.DB(), fileRecord.Filename, fileRecord.EdiString, fileRecord.PaymentRequestNumber)
		suite.NoError(err)

		got, err := models.FetchPaymentRequestEdiByPaymentRequestNumber(suite.DB(), fileRecord.PaymentRequestNumber)
		suite.NoError(err)
		suite.Equal(fileRecord.Filename, got.Filename)
		suite.Equal(fileRecord.EdiString, got.EdiString)
		suite.Equal(fileRecord.PaymentRequestNumber, got.PaymentRequestNumber)
	})

}
