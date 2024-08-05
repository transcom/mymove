package invoice

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ProcessTPPSPaidInvoiceReportSuite struct {
	*testingsuite.PopTestSuite
}

func TestProcessTPPSPaidInvoiceReportSuite(t *testing.T) {
	ts := &ProcessTPPSPaidInvoiceReportSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

type FakeTPPSFile struct {
	contents string
}

func (suite *ProcessTPPSPaidInvoiceReportSuite) TestParsingTPPSPaidInvoiceReport() {
	tppsPaidInvoiceReportProcessor := NewTPPSPaidInvoiceReportProcessor()

	suite.Run("successfully proccesses a valid TPPSPaidInvoiceReport", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		sampleTPPSPaidInvoiceReportString := FakeTPPSFile{
			`Invoice Number From Invoice	Document Create Date	Seller Paid Date	Invoice Total Charges	Line Description	Product Description	Line Billing Units	Line Unit Price	Line Net Charge	PO/TCN	Line Number	First Note Code	First Note Code Description	First Note To	First Note Message	Second Note Code	Second Note Code Description	Second Note To	Second Note Message	Third Note Code	Third Note Code Description	Third Note To	Third Note Message
1841-7267-3	2024-07-29	2024-07-30	1151.55	DDP	DDP	3760	0.0077	28.95	1841-7267-826285fc	1                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066
1841-7267-3	2024-07-29	2024-07-30	1151.55	FSC	FSC	3760	0.0014	5.39	1841-7267-aeb3cfea	4                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066
1841-7267-3	2024-07-29	2024-07-30	1151.55	DLH	DLH	3760	0.2656	998.77	1841-7267-c8ea170b	2                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066
1841-7267-3	2024-07-29	2024-07-30	1151.55	DUPK	DUPK	3760	0.0315	118.44	1841-7267-265c16d7	3                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066
9436-4123-3	2024-07-29	2024-07-30	125.25	DDP	DDP	7500	0.0167	125.25	9436-4123-93761f93	1                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50057
`}
		factory.BuildPaymentRequestToInterchangeControlNumber(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequestToInterchangeControlNumber{
					InterchangeControlNumber: 100001251,
					EDIType:                  models.EDIType858,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		err := tppsPaidInvoiceReportProcessor.ProcessFile(suite.AppContextForTest(), "", sampleTPPSPaidInvoiceReportString.contents)
		suite.NoError(err)
	})
}
