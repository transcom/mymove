package gunsafeweightticket

// import (
// 	"github.com/gofrs/uuid"

// 	"github.com/transcom/mymove/pkg/apperror"
// 	"github.com/transcom/mymove/pkg/auth"
// 	"github.com/transcom/mymove/pkg/factory"
// )

// func (suite *ProgearWeightTicketSuite) TestProgearWeightTicketCreator() {
// 	suite.Run("Successfully creates a ProgearWeightTicket - Customer", func() {
// 		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
// 		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID
// 		session := &auth.Session{
// 			ServiceMemberID: serviceMemberID,
// 			ApplicationName: auth.MilApp,
// 		}

// 		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
// 		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

// 		suite.Nil(err)
// 		suite.NotNil(progearWeightTicket)
// 		suite.Equal(ppmShipment.ID, progearWeightTicket.PPMShipmentID)
// 		suite.NotNil(progearWeightTicket.DocumentID)
// 		suite.Equal(serviceMemberID, progearWeightTicket.Document.ServiceMemberID)
// 	})

// 	suite.Run("Fails when an invalid ppmShipmentID is used - Customer", func() {
// 		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
// 		session := &auth.Session{
// 			ServiceMemberID: serviceMember.ID,
// 			ApplicationName: auth.MilApp,
// 		}

// 		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
// 		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

// 		suite.Nil(progearWeightTicket)

// 		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

// 		suite.ErrorIs(err, expectedErr)
// 	})

// 	suite.Run("Fails when session has invalid serviceMemberID - Customer", func() {
// 		session := &auth.Session{
// 			ServiceMemberID: uuid.Must(uuid.NewV4()),
// 			ApplicationName: auth.MilApp,
// 		}
// 		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

// 		progearWeightTicketCreator := NewCustomerProgearWeightTicketCreator()
// 		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

// 		suite.Nil(progearWeightTicket)

// 		expectedErr := apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")

// 		suite.ErrorIs(err, expectedErr)
// 	})

// 	suite.Run("Successfully creates a ProgearWeightTicket - Office", func() {
// 		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
// 		serviceMemberID := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID
// 		officeId, _ := uuid.NewV4()
// 		session := &auth.Session{
// 			OfficeUserID:    officeId,
// 			ApplicationName: auth.OfficeApp,
// 		}

// 		progearWeightTicketCreator := NewOfficeProgearWeightTicketCreator()
// 		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), ppmShipment.ID)

// 		suite.Nil(err)
// 		suite.NotNil(progearWeightTicket)
// 		suite.Equal(ppmShipment.ID, progearWeightTicket.PPMShipmentID)
// 		suite.NotNil(progearWeightTicket.DocumentID)
// 		suite.Equal(serviceMemberID, progearWeightTicket.Document.ServiceMemberID)
// 	})

// 	suite.Run("Fails when an invalid ppmShipmentID is used - Office", func() {
// 		officeId, _ := uuid.NewV4()
// 		session := &auth.Session{
// 			OfficeUserID:    officeId,
// 			ApplicationName: auth.OfficeApp,
// 		}

// 		progearWeightTicketCreator := NewOfficeProgearWeightTicketCreator()
// 		progearWeightTicket, err := progearWeightTicketCreator.CreateProgearWeightTicket(suite.AppContextWithSessionForTest(session), uuid.Nil)

// 		suite.Nil(progearWeightTicket)

// 		expectedErr := apperror.NewNotFoundError(uuid.Nil, "while looking for PPMShipment")

// 		suite.ErrorIs(err, expectedErr)
// 	})
// }
