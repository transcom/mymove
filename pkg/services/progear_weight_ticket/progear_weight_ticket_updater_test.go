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
)

func (suite ProgearWeightTicketSuite) TestUpdateProgearWeightTicket() {
	setupForTest := func(appCtx appcontext.AppContext, overrides *models.ProgearWeightTicket, hasdocFiles bool) *models.ProgearWeightTicket {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		ppmShipment := testdatagen.MakeMinimalDefaultPPMShipment(suite.DB())

		baseDocumentAssertions := testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		}

		document := testdatagen.MakeDocument(appCtx.DB(), baseDocumentAssertions)

		now := time.Now()
		if hasdocFiles {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = &now
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

		originalProgearWeightTicket := models.ProgearWeightTicket{
			DocumentID:    document.ID,
			PPMShipmentID: ppmShipment.ID,
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
		badProgearWeightTicket := models.ProgearWeightTicket{
			ID: uuid.Must(uuid.NewV4()),
		}

		updater := NewCustomerProgearWeightTicketUpdater()

		updatedProgearWeightTicket, err := updater.UpdateProgearWeightTicket(suite.AppContextForTest(), badProgearWeightTicket, "")

		suite.Nil(updatedProgearWeightTicket)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for ProgearWeightTicket", badProgearWeightTicket.ID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		appCtx := suite.AppContextForTest()

		originalProgearWeightTicket := setupForTest(appCtx, nil, false)

		updater := NewCustomerProgearWeightTicketUpdater()

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *originalProgearWeightTicket, "")

		suite.Nil(updatedProgearWeightTicket)

		if suite.Error(updateErr) {
			suite.IsType(apperror.PreconditionFailedError{}, updateErr)

			suite.Equal(
				fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", originalProgearWeightTicket.ID.String()),
				updateErr.Error(),
			)
		}
	})

	suite.Run("Successfully updates", func() {
		appCtx := suite.AppContextForTest()

		originalProgearWeightTicket := setupForTest(appCtx, nil, true)

		updater := NewCustomerProgearWeightTicketUpdater()

		desiredProgearWeightTicket := &models.ProgearWeightTicket{
			ID:               originalProgearWeightTicket.ID,
			Description:      models.StringPointer("Self progear"),
			Weight:           models.PoundPointer(3000),
			HasWeightTickets: models.BoolPointer(true),
			BelongsToSelf:    models.BoolPointer(true),
		}

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalProgearWeightTicket.ID, updatedProgearWeightTicket.ID)
		suite.Equal(originalProgearWeightTicket.DocumentID, updatedProgearWeightTicket.DocumentID)
		suite.Equal(*desiredProgearWeightTicket.Description, *updatedProgearWeightTicket.Description)
		suite.Equal(*desiredProgearWeightTicket.Weight, *updatedProgearWeightTicket.Weight)
		suite.Equal(*desiredProgearWeightTicket.HasWeightTickets, *updatedProgearWeightTicket.HasWeightTickets)
		suite.Equal(*desiredProgearWeightTicket.BelongsToSelf, *updatedProgearWeightTicket.BelongsToSelf)
	})

	suite.Run("Succesfully updates when files are required", func() {
		appCtx := suite.AppContextForTest()

		originalProgearWeightTicket := setupForTest(appCtx, nil, true)

		updater := NewCustomerProgearWeightTicketUpdater()

		desiredProgearWeightTicket := &models.ProgearWeightTicket{
			ID:               originalProgearWeightTicket.ID,
			Description:      models.StringPointer("Self progear"),
			Weight:           models.PoundPointer(3000),
			HasWeightTickets: models.BoolPointer(true),
			BelongsToSelf:    models.BoolPointer(true),
		}

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalProgearWeightTicket.ID, updatedProgearWeightTicket.ID)
		suite.Equal(originalProgearWeightTicket.DocumentID, updatedProgearWeightTicket.DocumentID)
		suite.Equal(*desiredProgearWeightTicket.Description, *updatedProgearWeightTicket.Description)
		suite.Equal(*desiredProgearWeightTicket.Weight, *updatedProgearWeightTicket.Weight)
		suite.Equal(*desiredProgearWeightTicket.HasWeightTickets, *updatedProgearWeightTicket.HasWeightTickets)
		suite.Equal(*desiredProgearWeightTicket.BelongsToSelf, *updatedProgearWeightTicket.BelongsToSelf)
		suite.Equal(1, len(updatedProgearWeightTicket.Document.UserUploads))
	})

	suite.Run("Fails to update when files are missing", func() {
		appCtx := suite.AppContextForTest()

		originalProgearWeightTicket := setupForTest(appCtx, nil, false)

		updater := NewCustomerProgearWeightTicketUpdater()

		desiredProgearWeightTicket := &models.ProgearWeightTicket{
			ID:               originalProgearWeightTicket.ID,
			Description:      models.StringPointer("Self progear"),
			Weight:           models.PoundPointer(3000),
			HasWeightTickets: models.BoolPointer(true),
			BelongsToSelf:    models.BoolPointer(true),
		}

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

		suite.Nil(updatedProgearWeightTicket)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.Equal("Invalid input found while validating the progear weight ticket.", updateErr.Error())
	})

	suite.Run("Status and reason related", func() {
		suite.Run("successfully", func() {

			suite.Run("changes status and reason", func() {
				appCtx := suite.AppContextForTest()

				originalProgearWeightTicket := testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{})

				updater := NewOfficeProgearWeightTicketUpdater()

				status := models.PPMDocumentStatusExcluded

				desiredProgearWeightTicket := &models.ProgearWeightTicket{
					ID:     originalProgearWeightTicket.ID,
					Status: &status,
					Reason: models.StringPointer("bad data"),
				}

				updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedProgearWeightTicket)
				suite.Equal(*desiredProgearWeightTicket.Status, *updatedProgearWeightTicket.Status)
				suite.Equal(*desiredProgearWeightTicket.Reason, *updatedProgearWeightTicket.Reason)
			})

			suite.Run("changes reason", func() {
				appCtx := suite.AppContextForTest()

				status := models.PPMDocumentStatusExcluded
				originalProgearWeightTicket := testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
					ProgearWeightTicket: models.ProgearWeightTicket{
						Status: &status,
						Reason: models.StringPointer("some temporary reason"),
					},
				})

				updater := NewOfficeProgearWeightTicketUpdater()

				desiredProgearWeightTicket := &models.ProgearWeightTicket{
					ID:     originalProgearWeightTicket.ID,
					Reason: models.StringPointer("bad data"),
				}

				updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedProgearWeightTicket)
				suite.Equal(status, *updatedProgearWeightTicket.Status)
				suite.Equal(*desiredProgearWeightTicket.Reason, *updatedProgearWeightTicket.Reason)
			})

			suite.Run("changes reason from rejected to approved", func() {
				appCtx := suite.AppContextForTest()

				status := models.PPMDocumentStatusExcluded
				originalProgearWeightTicket := testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
					ProgearWeightTicket: models.ProgearWeightTicket{
						Status: &status,
						Reason: models.StringPointer("some temporary reason"),
					},
				})

				updater := NewOfficeProgearWeightTicketUpdater()

				desiredStatus := models.PPMDocumentStatusApproved
				desiredProgearWeightTicket := &models.ProgearWeightTicket{
					ID:     originalProgearWeightTicket.ID,
					Status: &desiredStatus,
					Reason: models.StringPointer(""),
				}

				updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedProgearWeightTicket)
				suite.Equal(desiredStatus, *updatedProgearWeightTicket.Status)
				suite.Equal((*string)(nil), updatedProgearWeightTicket.Reason)
			})
		})

		suite.Run("fails", func() {
			suite.Run("to update when status or reason are changed", func() {
				appCtx := suite.AppContextForTest()

				originalProgearWeightTicket := setupForTest(appCtx, nil, true)

				updater := NewCustomerProgearWeightTicketUpdater()

				status := models.PPMDocumentStatusExcluded

				desiredProgearWeightTicket := &models.ProgearWeightTicket{
					ID:               originalProgearWeightTicket.ID,
					Description:      models.StringPointer("Self progear"),
					Weight:           models.PoundPointer(3000),
					HasWeightTickets: models.BoolPointer(true),
					BelongsToSelf:    models.BoolPointer(true),
					Status:           &status,
					Reason:           models.StringPointer("bad data"),
				}

				updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

				suite.Nil(updatedProgearWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("Invalid input found while validating the progear weight ticket.", updateErr.Error())
			})

			suite.Run("to update status", func() {
				appCtx := suite.AppContextForTest()

				status := models.PPMDocumentStatusExcluded
				originalProgearWeightTicket := testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
					ProgearWeightTicket: models.ProgearWeightTicket{
						Status: &status,
						Reason: models.StringPointer("some temporary reason"),
					},
				})

				updater := NewOfficeProgearWeightTicketUpdater()

				desiredStatus := models.PPMDocumentStatusApproved
				desiredProgearWeightTicket := &models.ProgearWeightTicket{
					ID:     originalProgearWeightTicket.ID,
					Status: &desiredStatus,
					Reason: models.StringPointer("bad data"),
				}

				updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

				suite.Nil(updatedProgearWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("Invalid input found while validating the progear weight ticket.", updateErr.Error())
			})

			suite.Run("to update because of invalid status", func() {
				appCtx := suite.AppContextForTest()

				originalProgearWeightTicket := testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{})

				updater := NewOfficeProgearWeightTicketUpdater()

				status := models.PPMDocumentStatus("invalid status")
				desiredProgearWeightTicket := &models.ProgearWeightTicket{
					ID:     originalProgearWeightTicket.ID,
					Status: &status,
					Reason: models.StringPointer("bad data"),
				}

				updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(appCtx, *desiredProgearWeightTicket, etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt))

				suite.Nil(updatedProgearWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("invalid input found while updating the ProgearWeightTicket", updateErr.Error())
			})
		})
	})
}
