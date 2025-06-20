package gunsafeweightticket

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *GunSafeWeightTicketSuite) TestUpdateGunSafeWeightTicket() {
	setupForTest := func(_ *models.GunSafeWeightTicket, hasdocFiles bool) *models.GunSafeWeightTicket {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)

		document := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		now := time.Now()
		if hasdocFiles {
			for i := 0; i < 2; i++ {
				var deletedAt *time.Time
				if i == 1 {
					deletedAt = &now
				}
				factory.BuildUserUpload(suite.DB(), []factory.Customization{
					{
						Model:    document,
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

		originalGunSafeWeightTicket := factory.BuildGunSafeWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
			{
				Model:    document,
				LinkOnly: true,
			},
			{
				Model: models.GunSafeWeightTicket{
					DocumentID:    document.ID,
					PPMShipmentID: ppmShipment.ID,
				},
			},
		}, nil)
		suite.NotNil(originalGunSafeWeightTicket.ID)

		return &originalGunSafeWeightTicket
	}

	suite.Run("Returns an error if the original doesn't exist", func() {
		badGunSafeWeightTicket := models.GunSafeWeightTicket{
			ID: uuid.Must(uuid.NewV4()),
		}
		originalGunSafeWeightTicket := setupForTest(nil, false)

		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalGunSafeWeightTicket.Document.ServiceMemberID,
		})

		updater := NewCustomerGunSafeWeightTicketUpdater()

		updatedGunSafeWeightTicket, err := updater.UpdateGunSafeWeightTicket(session, badGunSafeWeightTicket, "")

		suite.Nil(updatedGunSafeWeightTicket)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for GunSafeWeightTicket", badGunSafeWeightTicket.ID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		originalGunSafeWeightTicket := setupForTest(nil, false)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalGunSafeWeightTicket.Document.ServiceMemberID,
		})
		suite.Equal(originalGunSafeWeightTicket.Document.ServiceMemberID, session.Session().ServiceMemberID)
		updater := NewCustomerGunSafeWeightTicketUpdater()

		updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(session, *originalGunSafeWeightTicket, "")

		suite.Nil(updatedGunSafeWeightTicket)

		if suite.Error(updateErr) {
			suite.IsType(apperror.PreconditionFailedError{}, updateErr)

			suite.Equal(
				fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", originalGunSafeWeightTicket.ID.String()),
				updateErr.Error(),
			)
		}
	})

	suite.Run("Successfully updates", func() {
		originalGunSafeWeightTicket := setupForTest(nil, true)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalGunSafeWeightTicket.Document.ServiceMemberID,
		})

		updater := NewCustomerGunSafeWeightTicketUpdater()

		desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
			ID:               originalGunSafeWeightTicket.ID,
			Description:      models.StringPointer("Self gunsafe"),
			Weight:           models.PoundPointer(500),
			HasWeightTickets: models.BoolPointer(true),
		}

		updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalGunSafeWeightTicket.ID, updatedGunSafeWeightTicket.ID)
		suite.Equal(originalGunSafeWeightTicket.DocumentID, updatedGunSafeWeightTicket.DocumentID)
		suite.Equal(*desiredGunSafeWeightTicket.Description, *updatedGunSafeWeightTicket.Description)
		suite.Equal(*desiredGunSafeWeightTicket.Weight, *updatedGunSafeWeightTicket.Weight)
		suite.Equal(*desiredGunSafeWeightTicket.HasWeightTickets, *updatedGunSafeWeightTicket.HasWeightTickets)
		suite.Equal(*desiredGunSafeWeightTicket.Weight, *updatedGunSafeWeightTicket.SubmittedWeight)
		suite.Equal(*desiredGunSafeWeightTicket.HasWeightTickets, *updatedGunSafeWeightTicket.SubmittedHasWeightTickets)
	})

	suite.Run("Succesfully updates when files are required", func() {
		originalGunSafeWeightTicket := setupForTest(nil, true)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalGunSafeWeightTicket.Document.ServiceMemberID,
		})

		updater := NewCustomerGunSafeWeightTicketUpdater()

		desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
			ID:               originalGunSafeWeightTicket.ID,
			Description:      models.StringPointer("Self gunsafe"),
			Weight:           models.PoundPointer(500),
			HasWeightTickets: models.BoolPointer(true),
		}

		updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalGunSafeWeightTicket.ID, updatedGunSafeWeightTicket.ID)
		suite.Equal(originalGunSafeWeightTicket.DocumentID, updatedGunSafeWeightTicket.DocumentID)
		suite.Equal(*desiredGunSafeWeightTicket.Description, *updatedGunSafeWeightTicket.Description)
		suite.Equal(*desiredGunSafeWeightTicket.Weight, *updatedGunSafeWeightTicket.Weight)
		suite.Equal(*desiredGunSafeWeightTicket.HasWeightTickets, *updatedGunSafeWeightTicket.HasWeightTickets)
		suite.Equal(*desiredGunSafeWeightTicket.Weight, *updatedGunSafeWeightTicket.SubmittedWeight)
		suite.Equal(*desiredGunSafeWeightTicket.HasWeightTickets, *updatedGunSafeWeightTicket.SubmittedHasWeightTickets)
		suite.Len(updatedGunSafeWeightTicket.Document.UserUploads, 2)
	})

	suite.Run("Fails to update when files are missing", func() {
		originalGunSafeWeightTicket := setupForTest(nil, false)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalGunSafeWeightTicket.Document.ServiceMemberID,
		})

		updater := NewCustomerGunSafeWeightTicketUpdater()

		desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
			ID:               originalGunSafeWeightTicket.ID,
			Description:      models.StringPointer("Self gunsafe"),
			Weight:           models.PoundPointer(500),
			HasWeightTickets: models.BoolPointer(true),
		}

		updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

		suite.Nil(updatedGunSafeWeightTicket)
		suite.NotNil(updateErr)
		suite.IsType(apperror.InvalidInputError{}, updateErr)
		suite.Equal("Invalid input found while validating the gunSafe weight ticket.", updateErr.Error())
	})

	suite.Run("Status and reason related", func() {
		suite.Run("successfully", func() {

			suite.Run("changes status and reason", func() {
				originalGunSafeWeightTicket := factory.BuildGunSafeWeightTicket(suite.DB(), nil, nil)
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})
				updater := NewOfficeGunSafeWeightTicketUpdater()

				status := models.PPMDocumentStatusExcluded

				desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
					ID:     originalGunSafeWeightTicket.ID,
					Status: &status,
					Reason: models.StringPointer("bad data"),
				}

				updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedGunSafeWeightTicket)
				suite.Equal(*desiredGunSafeWeightTicket.Status, *updatedGunSafeWeightTicket.Status)
				suite.Equal(*desiredGunSafeWeightTicket.Reason, *updatedGunSafeWeightTicket.Reason)
			})

			suite.Run("changes reason", func() {
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				status := models.PPMDocumentStatusExcluded
				originalGunSafeWeightTicket := factory.BuildGunSafeWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.GunSafeWeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

				updater := NewOfficeGunSafeWeightTicketUpdater()

				desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
					ID:     originalGunSafeWeightTicket.ID,
					Reason: models.StringPointer("bad data"),
				}

				updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedGunSafeWeightTicket)
				suite.Equal(status, *updatedGunSafeWeightTicket.Status)
				suite.Equal(*desiredGunSafeWeightTicket.Reason, *updatedGunSafeWeightTicket.Reason)
			})

			suite.Run("changes reason from rejected to approved", func() {
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				status := models.PPMDocumentStatusExcluded
				originalGunSafeWeightTicket := factory.BuildGunSafeWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.GunSafeWeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

				updater := NewOfficeGunSafeWeightTicketUpdater()

				desiredStatus := models.PPMDocumentStatusApproved
				desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
					ID:     originalGunSafeWeightTicket.ID,
					Status: &desiredStatus,
					Reason: models.StringPointer(""),
				}

				updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

				suite.Nil(updateErr)
				suite.NotNil(updatedGunSafeWeightTicket)
				suite.Equal(desiredStatus, *updatedGunSafeWeightTicket.Status)
				suite.Equal((*string)(nil), updatedGunSafeWeightTicket.Reason)
			})
		})

		suite.Run("fails", func() {
			suite.Run("to update when status or reason are changed", func() {
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				originalGunSafeWeightTicket := setupForTest(nil, true)

				updater := NewCustomerGunSafeWeightTicketUpdater()

				status := models.PPMDocumentStatusExcluded

				desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
					ID:               originalGunSafeWeightTicket.ID,
					Description:      models.StringPointer("Self gunsafe"),
					Weight:           models.PoundPointer(500),
					HasWeightTickets: models.BoolPointer(true),
					Status:           &status,
					Reason:           models.StringPointer("bad data"),
				}

				updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

				suite.Nil(updatedGunSafeWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("Invalid input found while validating the gunSafe weight ticket.", updateErr.Error())
			})

			suite.Run("to update status", func() {
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				status := models.PPMDocumentStatusExcluded
				originalGunSafeWeightTicket := factory.BuildGunSafeWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.GunSafeWeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

				updater := NewOfficeGunSafeWeightTicketUpdater()

				desiredStatus := models.PPMDocumentStatusApproved
				desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
					ID:     originalGunSafeWeightTicket.ID,
					Status: &desiredStatus,
					Reason: models.StringPointer("bad data"),
				}

				updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

				suite.Nil(updatedGunSafeWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("Invalid input found while validating the gunSafe weight ticket.", updateErr.Error())
			})

			suite.Run("to update because of invalid status", func() {
				appCtx := suite.AppContextWithSessionForTest(
					&auth.Session{
						ApplicationName: auth.OfficeApp,
						OfficeUserID:    uuid.Must(uuid.NewV4()),
					})

				originalGunSafeWeightTicket := factory.BuildGunSafeWeightTicket(suite.DB(), nil, nil)

				updater := NewOfficeGunSafeWeightTicketUpdater()

				status := models.PPMDocumentStatus("invalid status")
				desiredGunSafeWeightTicket := &models.GunSafeWeightTicket{
					ID:     originalGunSafeWeightTicket.ID,
					Status: &status,
					Reason: models.StringPointer("bad data"),
				}

				updatedGunSafeWeightTicket, updateErr := updater.UpdateGunSafeWeightTicket(appCtx, *desiredGunSafeWeightTicket, etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt))

				suite.Nil(updatedGunSafeWeightTicket)
				suite.NotNil(updateErr)
				suite.IsType(apperror.InvalidInputError{}, updateErr)
				suite.Equal("invalid input found while updating the GunSafeWeightTicket", updateErr.Error())
			})
		})
	})
}

func (suite *GunSafeWeightTicketSuite) TestFetchGunSafeWeightTicketByIDExcludeDeletedUploads() {
	var gunsafeWeightTicket models.GunSafeWeightTicket
	var serviceMember models.ServiceMember
	suite.PreloadData(func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMember = ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		document := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		gunsafeWeightTicket = factory.BuildGunSafeWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    document,
				LinkOnly: true,
			},
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
			{
				Model: models.GunSafeWeightTicket{
					DocumentID:    document.ID,
					PPMShipmentID: ppmShipment.ID,
				},
			},
		}, nil)
	})

	// Test successful fetch
	suite.Run("Returns a gunsafe weight ticket successfully with correct ID", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})
		fetchedGunSafeWeightTicket, err := FetchGunSafeWeightTicketByIDExcludeDeletedUploads(session, gunsafeWeightTicket.ID)
		suite.NoError(err)
		suite.Equal(gunsafeWeightTicket.ID, fetchedGunSafeWeightTicket.ID)
	})

	// Test 404 fetch
	suite.Run("Returns not found error when gunsafe weight ticket id doesn't exist", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})
		gunsafeID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(gunsafeID, "while looking for GunSafeWeightTicket")

		gunsafe, err := FetchGunSafeWeightTicketByIDExcludeDeletedUploads(session, gunsafeID)

		suite.Nil(gunsafe)
		suite.Equalf(err, expectedError, "while looking for GunSafeWeightTicket")
	})

	suite.Run("404 Not Found Error - gunsafe can only be fetched for service member associated with the current session", func() {
		maliciousSession := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})

		gunsafe, err := FetchGunSafeWeightTicketByIDExcludeDeletedUploads(maliciousSession, gunsafeWeightTicket.ID)
		suite.Error(err)
		suite.Nil(gunsafe)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
