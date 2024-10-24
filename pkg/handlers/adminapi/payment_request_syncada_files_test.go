// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used set up environment variables
// RA: Given the functions causing the lint errors are used to set environment variables for testing purposes, it does not present a risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package adminapi

import (
	"time"

	"github.com/gofrs/uuid"

	pre "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/payment_request_syncada_files"
	"github.com/transcom/mymove/pkg/models"
	fetch "github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

const (
	edi858cA = "ISA*00*0084182369*00*_   _*ZZ*GOVDPIBS*12*8004171844*20200921*1459*U*00401*100001272*0*T*|"
	edi858cB = "ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *241009*1912*U*00401*404551885*0*T*|"
)

func (suite *HandlerSuite) TestIndexPaymentRequestSyncadaFilesHandler() {
	// test that everything is wired up
	suite.Run("payment request syncada files handler result in ok response", func() {
		prsf := []models.PaymentRequestEdiFile{
			{

				PaymentRequestNumber: "1234-1212-1",
				EdiString:            edi858cA,
				Filename:             "858-2770-1.txt",
			},
			{
				ID:                   uuid.Must(uuid.NewV4()),
				PaymentRequestNumber: "2345-9875-2",
				EdiString:            edi858cB,
				Filename:             "858-0324-1.txt",
			},
		}
		models.CreatePaymentRequestEdiFile(suite.DB(), prsf[0].Filename, prsf[0].EdiString, prsf[0].PaymentRequestNumber)
		models.CreatePaymentRequestEdiFile(suite.DB(), prsf[1].Filename, prsf[1].EdiString, prsf[1].PaymentRequestNumber)

		params := pre.IndexPaymentRequestSyncadaFilesParams{
			HTTPRequest: suite.setupAuthenticatedRequest("GET", "/payment-request-syncada-files"),
		}
		queryBuilder := query.NewQueryBuilder()
		handler := IndexPaymentRequestSyncadaFilesHandler{
			HandlerConfig:  suite.HandlerConfig(),
			NewQueryFilter: query.NewQueryFilter,
			ListFetcher:    fetch.NewListFetcher(queryBuilder),
			NewPagination:  pagination.NewPagination,
		}

		response := handler.Handle(params)
		suite.IsType(&pre.IndexPaymentRequestSyncadaFilesOK{}, response)
		okResponse := response.(*pre.IndexPaymentRequestSyncadaFilesOK)
		suite.Len(okResponse.Payload, 2)
		suite.Equal(prsf[0].PaymentRequestNumber, okResponse.Payload[0].PaymentRequestNumber)
	})

}

func (suite *HandlerSuite) TestPayloadForPaymentRequestEdiFile() {
	testID := uuid.Must(uuid.NewV4())
	testTime := time.Now()

	paymentRequestEdiFile := models.PaymentRequestEdiFile{
		ID:                   testID,
		PaymentRequestNumber: "TEST123",
		Filename:             "test_file.edi",
		EdiString:            "EDI content",
		CreatedAt:            testTime,
	}

	payload := payloadForPaymentRequestEdiFile(paymentRequestEdiFile)

	suite.Equal(testID.String(), payload.ID.String())
	suite.Equal("TEST123", payload.PaymentRequestNumber)
	suite.Equal("test_file.edi", payload.FileName)
	suite.Equal("EDI content", payload.EdiString)
	suite.Equal(testTime.UTC().Format("2006-01-02T15:04:05.000Z"), payload.CreatedAt.String())
}
