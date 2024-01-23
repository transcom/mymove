package progearweightticket

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
)

func (suite *ProgearWeightTicketSuite) TestDeleteProgearWeightTicket() {

	setupForTest := func(appCtx appcontext.AppContext, overrides *models.ProgearWeightTicket, hasDocumentUploads bool) *models.ProgearWeightTicket {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)
		progearDocument := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		if hasDocumentUploads {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = models.TimePointer(time.Now())
				}
				factory.BuildUserUpload(suite.DB(), []factory.Customization{
					{
						Model:    progearDocument,
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

		originalProgearWeightTicket := models.ProgearWeightTicket{
			PPMShipmentID: ppmShipment.ID,
			Document:      progearDocument,
			DocumentID:    progearDocument.ID,
		}

		if overrides != nil {
			testdatagen.MergeModels(&originalProgearWeightTicket, overrides)
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&originalProgearWeightTicket)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalProgearWeightTicket.ID)

		return &originalProgearWeightTicket
	}
	suite.Run("Returns an error if the original doesn't exist", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		notFoundProgearWeightTicketID := uuid.Must(uuid.NewV4())
		deleter := NewProgearWeightTicketDeleter()

		err := deleter.DeleteProgearWeightTicket(session, notFoundProgearWeightTicketID)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for ProgearWeightTicket", notFoundProgearWeightTicketID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Successfully deletes as a customer's progear weight ticket", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		originalProgearWeightTicket := setupForTest(session, nil, true)

		deleter := NewProgearWeightTicketDeleter()

		suite.Nil(originalProgearWeightTicket.DeletedAt)
		err := deleter.DeleteProgearWeightTicket(session, originalProgearWeightTicket.ID)
		suite.NoError(err)

		var progearWeightTicketInDB models.ProgearWeightTicket
		err = suite.DB().Find(&progearWeightTicketInDB, originalProgearWeightTicket.ID)
		suite.NoError(err)
		suite.NotNil(progearWeightTicketInDB.DeletedAt)

		// Should not delete associated PPM shipment
		var dbPPMShipment models.PPMShipment
		suite.NotNil(originalProgearWeightTicket.PPMShipmentID)
		err = suite.DB().Find(&dbPPMShipment, originalProgearWeightTicket.PPMShipmentID)
		suite.NoError(err)
		suite.Nil(dbPPMShipment.DeletedAt)

		// Should delete associated document
		var dbDocument models.Document
		suite.NotNil(originalProgearWeightTicket.DocumentID)
		err = suite.DB().Find(&dbDocument, originalProgearWeightTicket.DocumentID)
		suite.NoError(err)
		suite.NotNil(dbDocument.DeletedAt)
	})
}
