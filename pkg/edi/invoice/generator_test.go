package ediinvoice_test

import (
	"log"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"go.uber.org/zap"
)

func (suite *InvoiceSuite) TestGenerate858C() {
	shipments := make([]models.Shipment, 1)
	shipments[0] = testdatagen.MakeDefaultShipment(suite.db)
	err := shipments[0].AssignGBLNumber(suite.db)
	suite.mustSave(&shipments[0])
	suite.NoError(err, "could not assign GBLNumber")

	var cost rateengine.CostComputation
	costByShipment := rateengine.CostByShipment{
		Shipment: shipments[0],
		Cost:     cost,
	}
	var costsByShipments []rateengine.CostByShipment
	costsByShipments = append(costsByShipments, costByShipment)

	generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())
	assert.Equal(suite.T(), generatedTransactions, ediinvoice.Invoice858C{})
	//var b bytes.Buffer
	//writer := edi.NewWriter(&b)
	//writer.WriteAll(generatedTransactions)
	//suite.NoError(err, "generates error")
	//suite.NotEmpty(b.String(), "result is empty")

	//re := regexp.MustCompile("\\*" + "T" + "\\*")
	//suite.True(re.MatchString(b.String()), "This fails if the EDI string does not have the environment flag set to T."+
	//	" This is set by the if statement in Generate858C() that checks a boolean variable named sendProductionInvoice")

	//assert.Equal(suite.T(), expectedEDI, b.String())
}

const expectedEDI = `ISA*00*0000000000*00*0000000000*ZZ*MYMOVE         *12*8004171844     *691231*1600*U*00401*000000001*1*T*|
GS*SI*MYMOVE*8004171844*19691231*1600*1*X*004010
ST*858*0001
BX*00*J*PP*KKFA7000001*MCCG**4
N9*DY*SC**
N9*CN*ABCD00001-1**
N9*PQ*ABBV2708**
N9*OQ*ORDER3*ARMY*20180315
N1*SF*Spacemen**
N3*123 Any Street*P.O. Box 12345
N4*Beverly Hills*CA*90210*US**
N1*RG*LKNQ*27*LKNQ
N1*RH*MLNQ*27*MLNQ
FA1*DZ
FA2*TA*F8E1
L10*108.200*B*L
HL*303**SS
L0*1*1.000*FR********
L1*0*0.0000*RC*0********LHS
HL*303**SS
L0*1***108.200*B******L
L1*0*65.7700*RC*0********105A
HL*304**SS
L0*1***108.200*B******L
L1*0*65.7700*RC*0********105C
HL*303**SS
L0*1***108.200*B******L
L1*0*4.0700*RC*0********135A
HL*304**SS
L0*1***108.200*B******L
L1*0*4.0700*RC*0********135B
HL*303**SS
L0*1*1.000*FR********
L1*0*0.0300*RC*22742********16A
SE*33*0001
GE*1*1
IEA*1*000000001
`

type InvoiceSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *InvoiceSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *InvoiceSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

func TestInvoiceSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &InvoiceSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
