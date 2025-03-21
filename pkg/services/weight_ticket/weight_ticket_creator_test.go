package weightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *WeightTicketSuite) TestWeightTicketCreator() {
	suite.Run("Successfully creates a WeightTicket - Customer", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
			ApplicationName: auth.MilApp,
		}

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

	suite.Run("Fails when an invalid ppmShipmentID is used - Customer", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
			ApplicationName: auth.MilApp,
		}

		weightTicketCreator := NewCustomerWeightTicketCreator()
		weightTicket, err := weightTicketCreator.CreateWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(weightTicket)
		suite.NotNil(err)
	})

	suite.Run("Fails when session has invalid serviceMemberID - Customer", func() {
		session := &auth.Session{
			ServiceMemberID: uuid.Nil,
			ApplicationName: auth.MilApp,
		}
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		weightTicketCreator := NewCustomerWeightTicketCreator()
		weightTicket, err := weightTicketCreator.CreateWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(weightTicket)
		suite.NotNil(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found while looking for PPMShipment")
	})

	suite.Run("Successfully creates a WeightTicket - Office", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		officeId, _ := uuid.NewV4()
		session := &auth.Session{
			OfficeUserID:    officeId,
			ApplicationName: auth.OfficeApp,
		}

		weightTicketCreator := NewOfficeWeightTicketCreator()
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

	suite.Run("Fails when an invalid ppmShipmentID is used - Office", func() {
		officeId, _ := uuid.NewV4()
		session := &auth.Session{
			OfficeUserID:    officeId,
			ApplicationName: auth.OfficeApp,
		}
		weightTicketCreator := NewOfficeWeightTicketCreator()
		weightTicket, err := weightTicketCreator.CreateWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(weightTicket)
		suite.NotNil(err)
	})
}
