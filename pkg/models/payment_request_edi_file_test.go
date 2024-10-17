package models_test

import (
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
}
