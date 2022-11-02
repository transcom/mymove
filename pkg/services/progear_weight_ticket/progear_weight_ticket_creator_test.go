package progearweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ProgearWeightTicketSuite) TestProgearWeightTicketCreator() {
	suite.Run("Successfully creates a ProgearWeightTicket", func() {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
		}

		ppmShipment := testdatagen.MakeMinimalDefaultPPMShipment(suite.DB())
		progearWeightTicketCreator := NewProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(progearWeightTicket)
		suite.Equal(ppmShipment.ID, progearWeightTicket.PPMShipmentID)
		suite.NotNil(progearWeightTicket.DocumentID)
		suite.Equal(serviceMember.ID, progearWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used", func() {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
		}

		progearWeightTicketCreator := NewProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(progearWeightTicket)
		suite.ErrorContains(err, "Valid PPMShipmentID is missing")
	})

	suite.Run("Invalid serviceMemberID", func() {
		session := &auth.Session{
			ServiceMemberID: uuid.Nil,
		}
		ppmShipment := testdatagen.MakeMinimalDefaultPPMShipment(suite.DB())

		progearWeightTicketCreator := NewProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(progearWeightTicket)
		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input due to invalid ServiceMemberID", err.Error())
	})
}
