package progearweightticket

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ProgearWeightTicketSuite) TestUpdateProgearWeightTicket() {
	setupForTest := func(_ *models.ProgearWeightTicket, hasdocFiles bool) *models.ProgearWeightTicket {
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

		originalProgearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
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
				Model: models.ProgearWeightTicket{
					DocumentID:    document.ID,
					PPMShipmentID: ppmShipment.ID,
				},
			},
		}, nil)
		suite.NotNil(originalProgearWeightTicket.ID)

		return &originalProgearWeightTicket
	}

	suite.Run("Returns an error if the original doesn't exist", func() {
		badProgearWeightTicket := models.ProgearWeightTicket{
			ID: uuid.Must(uuid.NewV4()),
		}
		originalProgearWeightTicket := setupForTest(nil, false)

		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalProgearWeightTicket.Document.ServiceMemberID,
		})

		updater := NewCustomerProgearWeightTicketUpdater()

		updatedProgearWeightTicket, err := updater.UpdateProgearWeightTicket(session, badProgearWeightTicket, "")

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
		originalProgearWeightTicket := setupForTest(nil, false)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalProgearWeightTicket.Document.ServiceMemberID,
		})
		suite.Equal(originalProgearWeightTicket.Document.ServiceMemberID, session.Session().ServiceMemberID)
		updater := NewCustomerProgearWeightTicketUpdater()

		updatedProgearWeightTicket, updateErr := updater.UpdateProgearWeightTicket(session, *originalProgearWeightTicket, "")

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
		originalProgearWeightTicket := setupForTest(nil, true)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalProgearWeightTicket.Document.ServiceMemberID,
		})

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
		suite.Equal(*desiredProgearWeightTicket.Weight, *updatedProgearWeightTicket.SubmittedWeight)
		suite.Equal(*desiredProgearWeightTicket.BelongsToSelf, *updatedProgearWeightTicket.SubmittedBelongsToSelf)
		suite.Equal(*desiredProgearWeightTicket.HasWeightTickets, *updatedProgearWeightTicket.SubmittedHasWeightTickets)

		var ppm models.PPMShipment
		err := suite.DB().
			Q().
			EagerPreload("Shipment").
			Find(&ppm, originalProgearWeightTicket.PPMShipmentID)
		suite.NoError(err)
		suite.Equal(int(*desiredProgearWeightTicket.Weight), ppm.Shipment.ActualProGearWeight.Int())
		suite.Nil(ppm.Shipment.ActualSpouseProGearWeight)
	})

	suite.Run("Succesfully updates when files are required", func() {
		originalProgearWeightTicket := setupForTest(nil, true)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalProgearWeightTicket.Document.ServiceMemberID,
		})

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
		suite.Equal(*desiredProgearWeightTicket.Weight, *updatedProgearWeightTicket.SubmittedWeight)
		suite.Equal(*desiredProgearWeightTicket.BelongsToSelf, *updatedProgearWeightTicket.SubmittedBelongsToSelf)
		suite.Equal(*desiredProgearWeightTicket.HasWeightTickets, *updatedProgearWeightTicket.SubmittedHasWeightTickets)
		suite.Len(updatedProgearWeightTicket.Document.UserUploads, 2)
	})

	suite.Run("Fails to update when files are missing", func() {
		originalProgearWeightTicket := setupForTest(nil, false)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: originalProgearWeightTicket.Document.ServiceMemberID,
		})

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
				originalProgearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})
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
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				status := models.PPMDocumentStatusExcluded
				originalProgearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.ProgearWeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

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
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				status := models.PPMDocumentStatusExcluded
				originalProgearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.ProgearWeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

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
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				originalProgearWeightTicket := setupForTest(nil, true)

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
				appCtx := suite.AppContextWithSessionForTest(&auth.Session{
					ApplicationName: auth.OfficeApp,
					OfficeUserID:    uuid.Must(uuid.NewV4()),
				})

				status := models.PPMDocumentStatusExcluded
				originalProgearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
					{
						Model: models.ProgearWeightTicket{
							Status: &status,
							Reason: models.StringPointer("some temporary reason"),
						},
					},
				}, nil)

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
				appCtx := suite.AppContextWithSessionForTest(
					&auth.Session{
						ApplicationName: auth.OfficeApp,
						OfficeUserID:    uuid.Must(uuid.NewV4()),
					})

				originalProgearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)

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

func (suite *ProgearWeightTicketSuite) TestFetchProgearWeightTicketByIDExcludeDeletedUploads() {
	var progearWeightTicket models.ProgearWeightTicket
	var serviceMember models.ServiceMember
	suite.PreloadData(func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		serviceMember = ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		document := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		progearWeightTicket = factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
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
				Model: models.ProgearWeightTicket{
					DocumentID:    document.ID,
					PPMShipmentID: ppmShipment.ID,
				},
			},
		}, nil)
	})

	// Test successful fetch
	suite.Run("Returns a progear weight ticket successfully with correct ID", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})
		fetchedProgearWeightTicket, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(session, progearWeightTicket.ID)
		suite.NoError(err)
		suite.Equal(progearWeightTicket.ID, fetchedProgearWeightTicket.ID)
	})

	// Test 404 fetch
	suite.Run("Returns not found error when progear weight ticket id doesn't exist", func() {
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: serviceMember.ID,
		})
		progearID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(progearID, "while looking for ProgearWeightTicket")

		progear, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(session, progearID)

		suite.Nil(progear)
		suite.Equalf(err, expectedError, "while looking for ProgearWeightTicket")
	})

	suite.Run("404 Not Found Error - progear can only be fetched for service member associated with the current session", func() {
		maliciousSession := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: uuid.Must(uuid.NewV4()),
		})

		progear, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(maliciousSession, progearWeightTicket.ID)
		suite.Error(err)
		suite.Nil(progear)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}

func (suite *ProgearWeightTicketSuite) TestUpdateProgearWeightTicketTotalSumsCorrectly() {
	serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)
	ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{Model: serviceMember, LinkOnly: true},
	}, nil)
	document := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

	factoryTickets := []struct {
		weight        unit.Pound
		belongsToSelf bool
	}{
		{weight: 100, belongsToSelf: true},
		{weight: 200, belongsToSelf: true},
		{weight: 50, belongsToSelf: false},
		{weight: 25, belongsToSelf: false},
	}

	var tickets []models.ProgearWeightTicket
	for _, ft := range factoryTickets {
		t := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
			{Model: serviceMember, LinkOnly: true},
			{Model: ppmShipment, LinkOnly: true},
			{Model: document, LinkOnly: true},
			{
				Model: models.ProgearWeightTicket{
					Weight:        models.PoundPointer(ft.weight),
					BelongsToSelf: models.BoolPointer(ft.belongsToSelf),
				},
			},
		}, nil)
		tickets = append(tickets, t)
	}

	appCtx := suite.AppContextWithSessionForTest(&auth.Session{
		ApplicationName: auth.MilApp,
		ServiceMemberID: serviceMember.ID,
	})

	updater := NewCustomerProgearWeightTicketUpdater()
	for _, t := range tickets {
		et := etag.GenerateEtag(t.UpdatedAt)
		updated, err := updater.UpdateProgearWeightTicket(appCtx, t, et)
		suite.NoError(err, "updating ticket %s", t.ID)
		suite.NotNil(updated)
	}

	var shipment models.MTOShipment
	suite.NoError(suite.DB().
		Q().
		Find(&shipment, ppmShipment.ShipmentID))

	suite.Equal(300, shipment.ActualProGearWeight.Int())
	suite.Equal(75, shipment.ActualSpouseProGearWeight.Int())
}
