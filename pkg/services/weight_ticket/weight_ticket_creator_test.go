package weightticket

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WeightTicketSuite) TestWeightTicketCreator() {
	suite.Run("Can successfully create a WeightTicket", func() {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
		}

		ppmShipment := testdatagen.MakeMinimalDefaultPPMShipment(suite.DB())
		weightTicketCreator := NewCustomerWeightTicketCreator()
		weightTicket, err := weightTicketCreator.CreateWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(weightTicket)
		suite.Equal(ppmShipment.ID, weightTicket.PPMShipmentID)
		suite.NotNil(weightTicket.EmptyDocumentID)
		suite.NotNil(weightTicket.FullDocumentID)
		suite.NotNil(weightTicket.ProofOfTrailerOwnershipDocumentID)
	})
}
