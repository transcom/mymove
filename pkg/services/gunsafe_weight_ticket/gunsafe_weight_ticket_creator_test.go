package gunsafeweightticket

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *GunSafeWeightTicketSuite) TestGunSafeWeightTicketCreator() {
	suite.Run("Successfully creates a GunSafeWeightTicket - Customer", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID
		session := &auth.Session{
			ServiceMemberID: serviceMemberID,
			ApplicationName: auth.MilApp,
		}

		gunSafeWeightTicketCreator := NewCustomerGunSafeWeightTicketCreator()
		gunSafeWeightTicket, err := gunSafeWeightTicketCreator.CreateGunSafeWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(gunSafeWeightTicket)
		suite.Equal(ppmShipment.ID, gunSafeWeightTicket.PPMShipmentID)
		suite.NotNil(gunSafeWeightTicket.DocumentID)
		suite.Equal(serviceMemberID, gunSafeWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used - Customer", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		session := &auth.Session{
			ServiceMemberID: serviceMember.ID,
			ApplicationName: auth.MilApp,
		}

		gunSafeWeightTicketCreator := NewCustomerGunSafeWeightTicketCreator()
		gunSafeWeightTicket, err := gunSafeWeightTicketCreator.CreateGunSafeWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(gunSafeWeightTicket)

		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})

	suite.Run("Fails when session has invalid serviceMemberID - Customer", func() {
		session := &auth.Session{
			ServiceMemberID: uuid.Must(uuid.NewV4()),
			ApplicationName: auth.MilApp,
		}
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		gunSafeWeightTicketCreator := NewCustomerGunSafeWeightTicketCreator()
		gunSafeWeightTicket, err := gunSafeWeightTicketCreator.CreateGunSafeWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(gunSafeWeightTicket)

		expectedErr := apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})

	suite.Run("Successfully creates a GunSafeWeightTicket - Office", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID
		officeId, _ := uuid.NewV4()
		session := &auth.Session{
			OfficeUserID:    officeId,
			ApplicationName: auth.OfficeApp,
		}

		gunSafeWeightTicketCreator := NewOfficeGunSafeWeightTicketCreator()
		gunSafeWeightTicket, err := gunSafeWeightTicketCreator.CreateGunSafeWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

		suite.Nil(err)
		suite.NotNil(gunSafeWeightTicket)
		suite.Equal(ppmShipment.ID, gunSafeWeightTicket.PPMShipmentID)
		suite.NotNil(gunSafeWeightTicket.DocumentID)
		suite.Equal(serviceMemberID, gunSafeWeightTicket.Document.ServiceMemberID)
	})

	suite.Run("Fails when an invalid ppmShipmentID is used - Office", func() {
		officeId, _ := uuid.NewV4()
		session := &auth.Session{
			OfficeUserID:    officeId,
			ApplicationName: auth.OfficeApp,
		}

		gunSafeWeightTicketCreator := NewOfficeGunSafeWeightTicketCreator()
		gunSafeWeightTicket, err := gunSafeWeightTicketCreator.CreateGunSafeWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

		suite.Nil(gunSafeWeightTicket)

		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

		suite.ErrorIs(err, expectedErr)
	})
}
