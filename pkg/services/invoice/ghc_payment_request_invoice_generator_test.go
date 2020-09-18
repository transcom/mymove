package invoice

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// 1) Create a generated EDI in the test
// 2) Pull in the Golden file EDI
// 3) Create test to compare the golden EDI (expected) to the generated edi (actual)

type GHCInvoiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *GHCInvoiceSuite) SetupTest() {
	errTruncateAll := suite.DB().TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
}

func TestGHCInvoiceSuite(t *testing.T) {
	ts := &GHCInvoiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("ghcinvoice")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

const testDateFormat = "20060102"
const testTimeFormat = "1504"

func (suite *GHCInvoiceSuite) TestAllGenerateEdi() {
	currentTime := time.Now()
	// generator := GHCPaymentRequestInvoiceGenerator{DB: suite.DB()}
	basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(dateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilledActual,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip3,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "242",
		},
	}
	// var paymentServiceItems models.PaymentServiceItems

	dlh := testdatagen.MakePaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeDLH,
		basicPaymentServiceItemParams,
	)
	// fsc := testdatagen.MakePaymentServiceItemWithParams(
	// 	suite.DB(),
	// 	models.ReServiceCodeFSC,
	// 	basicPaymentServiceItemParams,
	// )
	// cs := testdatagen.MakePaymentServiceItemWithParams(
	// 	suite.DB(),
	// 	models.ReServiceCodeCS,
	// 	basicPaymentServiceItemParams,
	// )
	// ms := testdatagen.MakePaymentServiceItemWithParams(
	// 	suite.DB(),
	// 	models.ReServiceCodeMS,
	// 	basicPaymentServiceItemParams,
	// )
	// paymentServiceItems = append(paymentServiceItems, dlh)
	suite.Equal(dlh, "result")

	// serviceMember := testdatagen.MakeExtendedServiceMember(suite.DB(), testdatagen.Assertions{
	// 	ServiceMember: models.ServiceMember{
	// 		UserID: uuid.FromStringOrNil("e038eeb6-f154-4cfe-b395-839b4278205d"),
	// 	},
	// })
	// mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})
	// paymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
	// 	PaymentRequest: models.PaymentRequest{
	// 		ID:                  uuid.FromStringOrNil("d66d9f35-218c-4b85-b9d1-631449b9d984"),
	// 		MoveTaskOrder:       mto,
	// 		IsFinal:             false,
	// 		Status:              models.PaymentRequestStatusPending,
	// 		RejectionReason:     nil,
	// 		PaymentServiceItems: paymentServiceItems,
	// 	},
	// })

	// result, err := generator.Generate(paymentRequest, false)
	// suite.NoError(err, "%f", paymentRequest)
	// actualEDIString, err := result.EDIString()
	// suite.NoError(err, "Failed to get invoice 858C as EDI string")
	//Test Invoice Start and End Segments
	// suite.Equal(paymentRequest, result)

}

func helperLoadExpectedEDI(suite *GHCInvoiceSuite, name string) string {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	suite.NoError(err, "error loading expected EDI fixture")
	return string(bytes)
}
