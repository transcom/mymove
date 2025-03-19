package progearweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *ProgearWeightTicketSuite) TestProgearWeightTicketCreator() {
	suite.Run("Successfully creates a ProgearWeightTicket - Customer", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMemberID,
		}

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(progearWeightTicket)
		suite.Equal(ppmShipment.ID, progearWeightTicket.PPMShipmentID)
		suite.NotNil(progearWeightTicket.DocumentID)
		suite.Equal(serviceMemberID, progearWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		}

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(progearWeightTicket)

		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})

	suite.Run("Fails when session has invalid serviceMemberID", func() {
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		}
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(progearWeightTicket)

		expectedErr := apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})

	suite.Run("Successfully creates a ProgearWeightTicket - Office", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID
		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			ServiceMemberID: serviceMemberID,
		}

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(progearWeightTicket)
		suite.Equal(ppmShipment.ID, progearWeightTicket.PPMShipmentID)
		suite.NotNil(progearWeightTicket.DocumentID)
		suite.Equal(serviceMemberID, progearWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		}

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(progearWeightTicket)

		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})

	suite.Run("Fails when session has invalid serviceMemberID", func() {
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		}
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(progearWeightTicket)

		expectedErr := apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})
}
