package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *WeightTicketSuite) TestWeightTicketCreator() {
	suite.Run("Successfully creates a WeightTicket", func() {
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
		suite.Equal(serviceMember.ID, weightTicket.EmptyDocument.ServiceMemberID)
		suite.NotNil(weightTicket.FullDocumentID)
		suite.Equal(serviceMember.ID, weightTicket.FullDocument.ServiceMemberID)
		suite.NotNil(weightTicket.ProofOfTrailerOwnershipDocumentID)
		suite.Equal(serviceMember.ID, weightTicket.ProofOfTrailerOwnershipDocument.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used", func() {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
		}

		weightTicketCreator := NewCustomerWeightTicketCreator()
		weightTicket, err := weightTicketCreator.CreateWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(weightTicket)
		suite.NotNil(err)
	})

	suite.Run("Fails when session has invalid serviceMemberID", func() {
		session := &auth.Session{
			ServiceMemberID: uuid.Nil,
		}
		ppmShipment := testdatagen.MakeMinimalDefaultPPMShipment(suite.DB())

		weightTicketCreator := NewCustomerWeightTicketCreator()
		weightTicket, err := weightTicketCreator.CreateWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(weightTicket)
		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while creating the Document.", err.Error())
	})
}
