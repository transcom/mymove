package gunsafeweightticket

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GunSafeWeightTicketSuite) TestDeleteGunSafeWeightTicket() {

	setupForTest := func(appCtx appcontext.AppContext, overrides *models.GunSafeWeightTicket, hasDocumentUploads bool) *models.GunSafeWeightTicket {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)
		gunSafeDocument := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		if hasDocumentUploads {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = models.TimePointer(time.Now())
				}
				factory.BuildUserUpload(suite.DB(), []factory.Customization{
					{
						Model:    gunSafeDocument,
						LinkOnly: true,
					},
					{
						Model: models.UserUpload{
							DeletedAt: deletedAt,
						},
					},
				}, nil)
			}
		}

		originalGunSafeWeightTicket := models.GunSafeWeightTicket{
			PPMShipmentID: ppmShipment.ID,
			Document:      gunSafeDocument,
			DocumentID:    gunSafeDocument.ID,
		}

		if overrides != nil {
			testdatagen.MergeModels(&originalGunSafeWeightTicket, overrides)
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&originalGunSafeWeightTicket)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalGunSafeWeightTicket.ID)

		return &originalGunSafeWeightTicket
	}
	suite.Run("Returns an error if the original doesn't exist", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		notFoundGunSafeWeightTicketID := uuid.Must(uuid.NewV4())
		deleter := NewGunSafeWeightTicketDeleter()

		err := deleter.DeleteGunSafeWeightTicket(session, uuid.Nil, notFoundGunSafeWeightTicketID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for GunSafeWeightTicket", notFoundGunSafeWeightTicketID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Successfully deletes as a customer's gunSafe weight ticket", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		originalGunSafeWeightTicket := setupForTest(session, nil, true)

		deleter := NewGunSafeWeightTicketDeleter()

		suite.Nil(originalGunSafeWeightTicket.DeletedAt)
		err := deleter.DeleteGunSafeWeightTicket(session, originalGunSafeWeightTicket.PPMShipmentID, originalGunSafeWeightTicket.ID)
		suite.NoError(err)

		var gunSafeWeightTicketInDB models.GunSafeWeightTicket
		err = suite.DB().Find(&gunSafeWeightTicketInDB, originalGunSafeWeightTicket.ID)
		suite.NoError(err)
		suite.NotNil(gunSafeWeightTicketInDB.DeletedAt)

		// Should not delete associated PPM shipment
		var dbPPMShipment models.PPMShipment
		suite.NotNil(originalGunSafeWeightTicket.PPMShipmentID)
		err = suite.DB().Find(&dbPPMShipment, originalGunSafeWeightTicket.PPMShipmentID)
		suite.NoError(err)
		suite.Nil(dbPPMShipment.DeletedAt)

		// Should delete associated document
		var dbDocument models.Document
		suite.NotNil(originalGunSafeWeightTicket.DocumentID)
		err = suite.DB().Find(&dbDocument, originalGunSafeWeightTicket.DocumentID)
		suite.NoError(err)
		suite.NotNil(dbDocument.DeletedAt)
	})

	suite.Run("Successfully deletes and totals sum of tickets and updates mto_shipments actual_gunsafe_weight column", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{Model: serviceMember, LinkOnly: true},
		}, nil)
		document := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})

		factoryTickets := []struct {
			weight unit.Pound
		}{
			{weight: 100},
			{weight: 200},
			{weight: 50},
			{weight: 25},
		}

		var tickets []models.GunSafeWeightTicket
		for _, ft := range factoryTickets {
			t := factory.BuildGunSafeWeightTicket(suite.DB(), []factory.Customization{
				{Model: serviceMember, LinkOnly: true},
				{Model: ppmShipment, LinkOnly: true},
				{Model: document, LinkOnly: true},
				{
					Model: models.GunSafeWeightTicket{
						Weight: models.PoundPointer(ft.weight),
					},
				},
			}, nil)
			tickets = append(tickets, t)
		}

		deleter := NewGunSafeWeightTicketDeleter()

		err := deleter.DeleteGunSafeWeightTicket(appCtx, ppmShipment.ID, tickets[3].ID)
		suite.NoError(err)

		var shipment models.MTOShipment
		suite.NoError(suite.DB().
			Q().
			Find(&shipment, ppmShipment.ShipmentID))

		suite.Equal(350, shipment.ActualGunSafeWeight.Int())
	})
}
