package progearweightticket

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite ProgearWeightTicketSuite) TestUpdateProgearWeightTicket() {

	setupForTest := func(appCtx appcontext.AppContext, overrides *models.ProgearWeightTicket, hasWeightDocumentUploads bool) *models.ProgearWeightTicket {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		ppmShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMemberID: serviceMember.ID,
				ServiceMember:   serviceMember,
			},
		})

		document := testdatagen.MakeDocument(appCtx.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		})

		if hasWeightDocumentUploads {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = models.TimePointer(time.Now())
				}
				testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
					UserUpload: models.UserUpload{
						UploaderID: serviceMember.UserID,
						DocumentID: &document.ID,
						Document:   document,
						DeletedAt:  deletedAt,
					},
				})
			}
		}

		// ADD constructed weight?
		oldProgearWeightTicket := models.ProgearWeightTicket{
			PPMShipmentID: ppmShipment.ID,
			Document:      document,
			DocumentID:    document.ID,
		}

		if overrides != nil {
			testdatagen.MergeModels(&oldProgearWeightTicket, overrides)
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&oldProgearWeightTicket)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(oldProgearWeightTicket.ID)

		return &oldProgearWeightTicket
	}

	suite.Run("Returns an error if the old progear weight ticket doesn't exist", func() {
		notFoundProgearWeightTicket := models.ProgearWeightTicket{
			ID: uuid.Must(uuid.NewV4()),
		}

		updater := NewProgearWeightTicketUpdater()

		updatedProgearWeightTicket, err := updater.UpdateProgearWeightTicket(suite.AppContextForTest(), notFoundProgearWeightTicket, "")

		suite.Nil(updatedProgearWeightTicket)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for a ProgearWeightTicket", notFoundProgearWeightTicket.ID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		appCtx := suite.AppContextForTest()

		oldProgearWeightTicket := setupForTest(appCtx, nil, false)

		updater := NewProgearWeightTicketUpdater()

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *oldProgearWeightTicket, "")

		suite.Nil(updatedProgearWeightTicket)

		suite.Error(updateErr)
		suite.IsType(apperror.PreconditionFailedError{}, updateErr)

		suite.Equal(
			fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", oldProgearWeightTicket.ID.String()),
			updateErr.Error(),
		)
	})

	suite.Run("Successfully updates a progear weight ticket", func() {
		appCtx := suite.AppContextForTest()

		oldProgearWeightTicket := setupForTest(appCtx, nil, true)

		updater := NewProgearWeightTicketUpdater()
		rejectedStatus := models.PPMDocumentStatusRejected

		expectedProgearWeightTicket := &models.ProgearWeightTicket{
			ID:               oldProgearWeightTicket.ID,
			BelongsToSelf:    models.BoolPointer(true),
			Description:      models.StringPointer("Self Progear"),
			HasWeightTickets: models.BoolPointer(true),
			Weight:           models.PoundPointer(unit.Pound(100)),
			Status:           &rejectedStatus,
			Reason:           models.StringPointer("Some info missing"),
		}

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *expectedProgearWeightTicket, etag.GenerateEtag(oldProgearWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(oldProgearWeightTicket.ID, updatedProgearWeightTicket.ID)
		suite.Equal(oldProgearWeightTicket.DocumentID, updatedProgearWeightTicket.DocumentID)
		// filters out the deleted upload
		suite.Len(updatedProgearWeightTicket.Document.UserUploads, 1)
		suite.Equal(*expectedProgearWeightTicket.Description, *updatedProgearWeightTicket.Description)
		suite.Equal(*expectedProgearWeightTicket.BelongsToSelf, *updatedProgearWeightTicket.BelongsToSelf)
		suite.Equal(*expectedProgearWeightTicket.HasWeightTickets, *updatedProgearWeightTicket.HasWeightTickets)
		suite.Equal(*expectedProgearWeightTicket.Weight, *updatedProgearWeightTicket.Weight)
		suite.Equal(*expectedProgearWeightTicket.Status, *updatedProgearWeightTicket.Status)
		suite.Equal(*expectedProgearWeightTicket.Reason, *updatedProgearWeightTicket.Reason)
	})

	suite.Run("Successfully clears the reason when status of progear weight ticket is approved", func() {
		appCtx := suite.AppContextForTest()

		rejectedStatus := models.PPMDocumentStatusRejected
		oldProgearWeightTicket := setupForTest(appCtx, &models.ProgearWeightTicket{
			Status: &rejectedStatus,
			Reason: models.StringPointer("Can't add progear for spouse as your own"),
		}, true)

		updater := NewProgearWeightTicketUpdater()

		approvedStatus := models.PPMDocumentStatusApproved
		expectedProgearWeightTicket := &models.ProgearWeightTicket{
			ID:               oldProgearWeightTicket.ID,
			BelongsToSelf:    models.BoolPointer(true),
			Description:      models.StringPointer("Self Progear"),
			HasWeightTickets: models.BoolPointer(true),
			Weight:           models.PoundPointer(unit.Pound(100)),
			Status:           &approvedStatus,
		}

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *expectedProgearWeightTicket, etag.GenerateEtag(oldProgearWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(oldProgearWeightTicket.ID, updatedProgearWeightTicket.ID)
		suite.Equal(oldProgearWeightTicket.DocumentID, updatedProgearWeightTicket.DocumentID)
		suite.Equal(*oldProgearWeightTicket.Description, *updatedProgearWeightTicket.Description)
		suite.Equal(*oldProgearWeightTicket.BelongsToSelf, *updatedProgearWeightTicket.BelongsToSelf)
		suite.Equal(*oldProgearWeightTicket.HasWeightTickets, *updatedProgearWeightTicket.HasWeightTickets)
		suite.Equal(*oldProgearWeightTicket.Weight, *updatedProgearWeightTicket.Weight)
		suite.Equal(*oldProgearWeightTicket.Status, *updatedProgearWeightTicket.Status)
		suite.Nil(updatedProgearWeightTicket.Reason)
	})

	suite.Run("Fails to update when files are missing", func() {
		appCtx := suite.AppContextForTest()

		oldProgearWeightTicket := setupForTest(appCtx, nil, false)

		updater := NewProgearWeightTicketUpdater()

		approvedStatus := models.PPMDocumentStatusApproved
		expectedProgearWeightTicket := &models.ProgearWeightTicket{
			ID:               oldProgearWeightTicket.ID,
			BelongsToSelf:    models.BoolPointer(true),
			Description:      models.StringPointer("Self Progear"),
			HasWeightTickets: models.BoolPointer(false),
			Weight:           models.PoundPointer(unit.Pound(100)),
			Status:           &approvedStatus,
		}

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *expectedProgearWeightTicket, etag.GenerateEtag(oldProgearWeightTicket.UpdatedAt))

		suite.Nil(updatedProgearWeightTicket)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.ErrorContains(updateErr, "Missing Weight Tickets")
	})

	suite.Run("Fails to update when a reason isn't provided for non-approved status", func() {
		appCtx := suite.AppContextForTest()

		oldProgearWeightTicket := setupForTest(appCtx, nil, true)

		updater := NewProgearWeightTicketUpdater()

		rejectedStatus := models.PPMDocumentStatusRejected
		expectedProgearWeightTicket := &models.ProgearWeightTicket{
			ID:               oldProgearWeightTicket.ID,
			BelongsToSelf:    models.BoolPointer(true),
			Description:      models.StringPointer("Self Progear"),
			HasWeightTickets: models.BoolPointer(false),
			Weight:           models.PoundPointer(unit.Pound(100)),
			Status:           &rejectedStatus,
		}

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *expectedProgearWeightTicket, etag.GenerateEtag(oldProgearWeightTicket.UpdatedAt))

		suite.Nil(updatedProgearWeightTicket)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.ErrorContains(updateErr, "A reason must be provided when the status is EXCLUDED or REJECTED")
	})
}
