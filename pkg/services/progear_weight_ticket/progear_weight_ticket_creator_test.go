package progearweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *ProgearWeightTicketSuite) TestProgearWeightTicketCreator() {
	suite.Run("Successfully creates a ProgearWeightTicket", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
		}

		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(progearWeightTicket)
		suite.Equal(ppmShipment.ID, progearWeightTicket.PPMShipmentID)
		suite.NotNil(progearWeightTicket.DocumentID)
		suite.Equal(serviceMember.ID, progearWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
		}

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(progearWeightTicket)
		suite.NotNil(err)
	})

	suite.Run("Fails when session has invalid serviceMemberID", func() {
		session := &auth.Session{
			ServiceMemberID: uuid.Nil,
		}
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(progearWeightTicket)
		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while creating the Document.", err.Error())
	})
}
