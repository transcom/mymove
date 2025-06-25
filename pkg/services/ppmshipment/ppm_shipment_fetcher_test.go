package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *PPMShipmentSuite) TestPPMShipmentFetcher() {
	fakeS3 := storageTest.NewFakeS3Storage(true)

	userUploader, uploaderErr := uploader.NewUserUploader(fakeS3, uploader.MaxCustomerUserUploadFileSizeLimit)

	suite.FatalNoError(uploaderErr)

	fetcher := NewPPMShipmentFetcher()

	checkWeightTicketUploads := func(expected *models.PPMShipment, actual *models.PPMShipment) {
		suite.Equal(len(expected.WeightTickets), len(actual.WeightTickets))
		suite.Greater(len(actual.WeightTickets), 0)
		for i := range expected.WeightTickets {
			if suite.False(
				actual.WeightTickets[i].EmptyDocument.ID.IsNil(),
				"EmptyDocument ID should not be nil",
			) {
				suite.Equal(expected.WeightTickets[i].EmptyDocument.ID, actual.WeightTickets[i].EmptyDocument.ID)

				suite.Equal(
					len(expected.WeightTickets[i].EmptyDocument.UserUploads),
					len(actual.WeightTickets[i].EmptyDocument.UserUploads),
				)

				if suite.False(
					actual.WeightTickets[i].EmptyDocument.UserUploads[0].Upload.ID.IsNil(),
					"EmptyDocument UserUploads[0] ID should not be nil",
				) {
					suite.Equal(
						expected.WeightTickets[i].EmptyDocument.UserUploads[0].Upload.ID,
						actual.WeightTickets[i].EmptyDocument.UserUploads[0].Upload.ID,
					)
				}
			}

			if suite.False(
				actual.WeightTickets[i].FullDocument.ID.IsNil(),
				"FullDocument ID should not be nil",
			) {
				suite.Equal(expected.WeightTickets[i].FullDocument.ID, actual.WeightTickets[i].FullDocument.ID)

				suite.Equal(
					len(expected.WeightTickets[i].FullDocument.UserUploads),
					len(actual.WeightTickets[i].FullDocument.UserUploads),
				)

				if suite.False(
					actual.WeightTickets[i].FullDocument.UserUploads[0].Upload.ID.IsNil(),
					"FullDocument UserUploads[0] ID should not be nil",
				) {
					suite.Equal(
						expected.WeightTickets[i].FullDocument.UserUploads[0].Upload.ID,
						actual.WeightTickets[i].FullDocument.UserUploads[0].Upload.ID,
					)
				}
			}
		}
	}

	suite.Run("fetches PPM with GCC Multiplier data when applicable", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), factory.GetTraitActiveOfficeUser(), nil)
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          officeUser.User.ID,
			OfficeUserID:    officeUser.ID,
		})

		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: validGccMultiplierDate,
				},
			},
		}, nil)

		ppmShipmentReturned, err := fetcher.GetPPMShipment(
			appCtx,
			ppmShipment.ID,
			nil,
			nil,
		)

		if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
			suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)
			suite.NotNil(ppmShipment.GCCMultiplierID)
			suite.NotNil(ppmShipment.GCCMultiplier)
		}
	})

	suite.Run("GetPPMShipment", func() {
		suite.Run("Can fetch a PPM Shipment if there is no session (e.g. a prime request)", func() {
			appCtx := suite.AppContextWithSessionForTest(nil)

			ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				appCtx,
				ppmShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)
				suite.Nil(ppmShipment.GCCMultiplierID)
				suite.Nil(ppmShipment.GCCMultiplier)
			}
		})

		suite.Run("Can fetch a PPM Shipment if it is an office user making a request from the office app", func() {
			officeUser := factory.BuildOfficeUser(suite.DB(), factory.GetTraitActiveOfficeUser(), nil)

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.OfficeApp,
				UserID:          officeUser.User.ID,
				OfficeUserID:    officeUser.ID,
			})

			ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				appCtx,
				ppmShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)
				suite.Nil(ppmShipment.GCCMultiplierID)
				suite.Nil(ppmShipment.GCCMultiplier)
			}
		})

		suite.Run("Can fetch a PPM Shipment if it is a customer app request by the customer it belongs to", func() {
			ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
			serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.MilApp,
				UserID:          serviceMember.User.ID,
				ServiceMemberID: serviceMember.ID,
			})

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				appCtx,
				ppmShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)
				suite.Nil(ppmShipment.GCCMultiplierID)
				suite.Nil(ppmShipment.GCCMultiplier)
			}
		})

		suite.Run("Returns a not found error if it is a customer app request by a customer that it doesn't belong to", func() {
			maliciousUser := factory.BuildExtendedServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

			appCtx := suite.AppContextWithSessionForTest(&auth.Session{
				ApplicationName: auth.MilApp,
				UserID:          maliciousUser.User.ID,
				ServiceMemberID: maliciousUser.ID,
			})

			ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				appCtx,
				ppmShipment.ID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(ppmShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", ppmShipment.ID), err.Error())
			}
		})

		associationTestCases := map[string]struct {
			eagerPreloadAssociations []string
			successAssertionFunc     func(*models.PPMShipment, *models.PPMShipment)
		}{
			"No associations": {
				eagerPreloadAssociations: nil,
				successAssertionFunc: func(_ *models.PPMShipment, actual *models.PPMShipment) {
					suite.True(actual.Shipment.ID.IsNil())
					suite.Nil(actual.WeightTickets)
					suite.Nil(actual.ProgearWeightTickets)
					suite.Nil(actual.MovingExpenses)
					suite.NotNil(actual.W2Address)
					suite.Nil(actual.AOAPacket)
					suite.Nil(actual.PaymentPacket)
				},
			},
			"Shipment only": {
				eagerPreloadAssociations: []string{EagerPreloadAssociationShipment},
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					if suite.False(
						actual.Shipment.ID.IsNil(),
						"Shipment ID should not be nil",
					) {
						suite.Equal(expected.Shipment.ID, actual.Shipment.ID)
					}

					suite.Nil(actual.WeightTickets)
					suite.Nil(actual.ProgearWeightTickets)
					suite.Nil(actual.MovingExpenses)
					suite.NotNil(actual.W2Address)
					suite.Nil(actual.AOAPacket)
					suite.Nil(actual.PaymentPacket)
				},
			},
			"Weight tickets only": {
				eagerPreloadAssociations: []string{EagerPreloadAssociationWeightTickets},
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.True(
						actual.Shipment.ID.IsNil(),
						"Shipment ID should be nil",
					)

					suite.NotNil(actual.WeightTickets)
					suite.Equal(len(expected.WeightTickets), len(actual.WeightTickets))
					suite.Equal(expected.WeightTickets[0].ID, actual.WeightTickets[0].ID)

					suite.Nil(actual.ProgearWeightTickets)
					suite.Nil(actual.MovingExpenses)
					suite.NotNil(actual.W2Address)
					suite.Nil(actual.AOAPacket)
					suite.Nil(actual.PaymentPacket)
				},
			},
			"All eager preload associations": {
				eagerPreloadAssociations: GetListOfAllPreloadAssociations(),
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					if suite.False(
						actual.Shipment.ID.IsNil(),
						"Shipment ID should not be nil",
					) {
						suite.Equal(expected.Shipment.ID, actual.Shipment.ID)
					}

					if suite.False(actual.Shipment.MoveTaskOrder.ID.IsNil(), "MoveTaskOrder ID should not be nil") &&
						suite.False(actual.Shipment.MoveTaskOrder.Orders.ID.IsNil(), "Orders ID should not be nil") &&
						suite.False(
							actual.Shipment.MoveTaskOrder.Orders.ServiceMember.ID.IsNil(),
							"ServiceMember ID should not be nil",
						) {
						suite.Equal(expected.Shipment.MoveTaskOrder.Orders.ServiceMember.ID, actual.Shipment.MoveTaskOrder.Orders.ServiceMember.ID)
					}

					if suite.NotNil(actual.WeightTickets) {
						suite.Equal(len(expected.WeightTickets), len(actual.WeightTickets))
						suite.Equal(expected.WeightTickets[0].ID, actual.WeightTickets[0].ID)
					}

					if suite.NotNil(actual.ProgearWeightTickets) {
						suite.Equal(len(expected.ProgearWeightTickets), len(actual.ProgearWeightTickets))
						suite.Equal(expected.ProgearWeightTickets[0].ID, actual.ProgearWeightTickets[0].ID)

					}

					if suite.NotNil(actual.MovingExpenses) {
						suite.Equal(len(expected.MovingExpenses), len(actual.MovingExpenses))
						suite.Equal(expected.MovingExpenses[0].ID, actual.MovingExpenses[0].ID)
					}

					if suite.NotNil(actual.W2Address) {
						suite.Equal(expected.W2Address.ID, actual.W2Address.ID)
					}

					if suite.NotNil(actual.PickupAddress) {
						suite.Equal(expected.PickupAddress.ID, actual.PickupAddress.ID)
					}
					if suite.NotNil(actual.DestinationAddress) {
						suite.Equal(expected.DestinationAddress.ID, actual.DestinationAddress.ID)
					}

					if suite.NotNil(actual.AOAPacket) {
						suite.Equal(expected.AOAPacket.ID, actual.AOAPacket.ID)
					}

					if suite.NotNil(actual.PaymentPacket) {
						suite.Equal(expected.PaymentPacket.ID, actual.PaymentPacket.ID)
					}
				},
			},
		}

		for name, testCase := range associationTestCases {
			name, testCase := name, testCase

			suite.Run(fmt.Sprintf("Can fetch a PPM Shipment with associations: %s", name), func() {
				ppmShipment := factory.BuildPPMShipmentWithAllDocTypesApproved(
					suite.DB(),
					userUploader,
				)

				ppmShipmentReturned, err := fetcher.GetPPMShipment(
					suite.AppContextForTest(),
					ppmShipment.ID,
					testCase.eagerPreloadAssociations,
					nil,
				)

				if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
					suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)
					suite.Equal(ppmShipment.ShipmentID, ppmShipmentReturned.ShipmentID)

					testCase.successAssertionFunc(&ppmShipment, ppmShipmentReturned)
				}
			})
		}

		suite.Run("Return a shipment that has secondary and tertiary addresses", func() {
			ppmShipment := factory.BuildFullAddressPPMShipment(suite.DB(), nil, nil)

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				nil,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.NotNil(ppmShipmentReturned.TertiaryPickupAddressID)
				suite.NotNil(ppmShipmentReturned.TertiaryDestinationAddressID)
			}
		})

		suite.Run("Returns a not found error if the PPM Shipment does not exist", func() {
			nonexistentID := uuid.Must(uuid.NewV4())

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				nonexistentID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(ppmShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", nonexistentID), err.Error())
			}
		})

		suite.Run("Returns an error if an invalid association is requested", func() {
			ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

			invalidAssociation := "invalid"
			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				[]string{invalidAssociation},
				nil,
			)

			if suite.Error(err) && suite.Nil(ppmShipmentReturned) {
				suite.IsType(apperror.NotImplementedError{}, err)

				suite.Contains(
					err.Error(),
					fmt.Sprintf("Requested eager preload association %s is not implemented", invalidAssociation),
				)
			}
		})

		suite.Run("Returns an error if the shipment has been deleted", func() {
			ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

			err := utilities.SoftDestroy(suite.DB(), &ppmShipment)
			suite.FatalNoError(err)

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				nil,
				nil,
			)

			if suite.Error(err) && suite.Nil(ppmShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", ppmShipment.ID), err.Error())
			}
		})

		suite.Run("Excludes deleted documents", func() {
			ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(
				suite.DB(),
				userUploader)

			// create new ppm documents that are deleted
			now := time.Now()

			factory.AddWeightTicketToPPMShipment(suite.DB(), &ppmShipment,
				userUploader, &models.WeightTicket{
					DeletedAt: &now,
				})

			factory.AddProgearWeightTicketToPPMShipment(suite.DB(), &ppmShipment,
				userUploader, &models.ProgearWeightTicket{
					DeletedAt: &now,
				})

			factory.AddMovingExpenseToPPMShipment(suite.DB(), &ppmShipment,
				userUploader, &models.MovingExpense{
					DeletedAt: &now,
				})

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				[]string{
					EagerPreloadAssociationWeightTickets,
					EagerPreloadAssociationProgearWeightTickets,
					EagerPreloadAssociationMovingExpenses,
				},
				nil,
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.Equal(len(ppmShipment.WeightTickets)-1, len(ppmShipmentReturned.WeightTickets))
				suite.Equal(ppmShipment.WeightTickets[0].ID, ppmShipmentReturned.WeightTickets[0].ID)

				suite.Equal(len(ppmShipment.ProgearWeightTickets)-1, len(ppmShipmentReturned.ProgearWeightTickets))
				suite.Equal(ppmShipment.ProgearWeightTickets[0].ID, ppmShipmentReturned.ProgearWeightTickets[0].ID)

				suite.Equal(len(ppmShipment.MovingExpenses)-1, len(ppmShipmentReturned.MovingExpenses))
				suite.Equal(ppmShipment.MovingExpenses[0].ID, ppmShipmentReturned.MovingExpenses[0].ID)
			}
		})

		suite.Run("Can fetch a ppm shipment and get both eagerPreloadAssociations and postloadAssociations", func() {
			appCtx := suite.AppContextForTest()

			ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				appCtx,
				ppmShipment.ID,
				[]string{
					EagerPreloadAssociationWeightTickets,
				},
				[]string{
					PostLoadAssociationWeightTicketUploads,
				},
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)

				suite.NotNil(ppmShipmentReturned.WeightTickets)
				suite.Equal(len(ppmShipment.WeightTickets), len(ppmShipmentReturned.WeightTickets))
				suite.Equal(ppmShipment.WeightTickets[0].ID, ppmShipmentReturned.WeightTickets[0].ID)

				checkWeightTicketUploads(&ppmShipment, ppmShipmentReturned)
			}
		})

		suite.Run("Doesn't return postload association if a necessary higher level association isn't eagerly preloaded", func() {
			appCtx := suite.AppContextForTest()

			ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), userUploader, nil)

			suite.FatalTrue(len(ppmShipment.WeightTickets) > 0, "Test data that was set up is invalid, no weight tickets found")

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				appCtx,
				ppmShipment.ID,
				nil,
				[]string{
					PostLoadAssociationWeightTicketUploads,
				},
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)

				suite.Nil(ppmShipmentReturned.WeightTickets)
			}
		})

		suite.Run("Returns an error if an invalid postload association is passed in", func() {
			ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)

			invalidAssociation := "invalid"
			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				nil,
				[]string{invalidAssociation},
			)

			if suite.Error(err) && suite.Nil(ppmShipmentReturned) {
				suite.IsType(apperror.NotImplementedError{}, err)

				suite.Contains(
					err.Error(),
					fmt.Sprintf("Requested post load association %s is not implemented", invalidAssociation),
				)
			}
		})
	})

	suite.Run("PostloadAssociations", func() {
		postloadAssociationTestCases := map[string]struct {
			postloadAssociations []string
			successAssertionFunc func(*models.PPMShipment, *models.PPMShipment)
		}{
			"signed certification": {
				postloadAssociations: []string{PostLoadAssociationSignedCertification},
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.NotNil(actual.SignedCertification)
					suite.Equal(expected.SignedCertification.ID, actual.SignedCertification.ID)

					suite.Greater(len(actual.WeightTickets), 0)
					for i := range actual.WeightTickets {
						suite.True(actual.WeightTickets[i].EmptyDocument.ID.IsNil())
						suite.True(actual.WeightTickets[i].FullDocument.ID.IsNil())
					}

					suite.Greater(len(actual.ProgearWeightTickets), 0)
					for i := range actual.ProgearWeightTickets {
						suite.True(actual.ProgearWeightTickets[i].Document.ID.IsNil())
					}

					suite.Greater(len(actual.MovingExpenses), 0)
					for i := range actual.MovingExpenses {
						suite.True(actual.MovingExpenses[i].Document.ID.IsNil())
					}
				},
			},
			"weight ticket uploads": {
				postloadAssociations: []string{PostLoadAssociationWeightTicketUploads},
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.Nil(actual.SignedCertification)

					checkWeightTicketUploads(expected, actual)

					suite.Greater(len(actual.ProgearWeightTickets), 0)
					for i := range actual.ProgearWeightTickets {
						suite.True(actual.ProgearWeightTickets[i].Document.ID.IsNil())
					}

					suite.Greater(len(actual.MovingExpenses), 0)
					for i := range actual.MovingExpenses {
						suite.True(actual.MovingExpenses[i].Document.ID.IsNil())
					}
				},
			},
			"all post load associations": {
				postloadAssociations: GetListOfAllPostloadAssociations(),
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.NotNil(actual.SignedCertification)
					suite.Equal(expected.SignedCertification.ID, actual.SignedCertification.ID)

					checkWeightTicketUploads(expected, actual)

					suite.Equal(len(expected.ProgearWeightTickets), len(actual.ProgearWeightTickets))
					suite.Greater(len(actual.ProgearWeightTickets), 0)
					for i := range expected.ProgearWeightTickets {
						if suite.False(
							actual.ProgearWeightTickets[i].Document.ID.IsNil(),
							fmt.Sprintf("Expected ProgearWeightTicket %d document ID to not be nil", i),
						) {
							suite.Equal(
								expected.ProgearWeightTickets[i].Document.ID,
								actual.ProgearWeightTickets[i].Document.ID,
							)

							suite.Equal(
								len(expected.ProgearWeightTickets[i].Document.UserUploads),
								len(actual.ProgearWeightTickets[i].Document.UserUploads),
							)

							if suite.False(actual.ProgearWeightTickets[i].Document.UserUploads[0].Upload.ID.IsNil(),
								fmt.Sprintf("Expected ProgearWeightTicket %d document user upload ID to not be nil", i),
							) {
								suite.Equal(
									expected.ProgearWeightTickets[i].Document.UserUploads[0].Upload.ID,
									actual.ProgearWeightTickets[i].Document.UserUploads[0].Upload.ID,
								)
							}
						}
					}

					suite.Equal(len(expected.MovingExpenses), len(actual.MovingExpenses))
					suite.Greater(len(actual.MovingExpenses), 0)
					for i := range expected.MovingExpenses {
						if suite.False(
							actual.MovingExpenses[i].Document.ID.IsNil(),
							fmt.Sprintf("Expected MovingExpense %d document ID to not be nil", i),
						) {
							suite.Equal(expected.MovingExpenses[i].Document.ID, actual.MovingExpenses[i].Document.ID)

							suite.Equal(
								len(expected.MovingExpenses[i].Document.UserUploads),
								len(actual.MovingExpenses[i].Document.UserUploads),
							)

							if suite.False(actual.MovingExpenses[i].Document.UserUploads[0].Upload.ID.IsNil(),
								fmt.Sprintf("Expected MovingExpense %d document user upload ID to not be nil", i),
							) {
								suite.Equal(
									expected.MovingExpenses[i].Document.UserUploads[0].Upload.ID,
									actual.MovingExpenses[i].Document.UserUploads[0].Upload.ID,
								)
							}
						}
					}

					if suite.False(actual.Shipment.MoveTaskOrder.ID.IsNil(), "MoveTaskOrder ID should not be nil") &&
						suite.False(actual.Shipment.MoveTaskOrder.Orders.ID.IsNil(), "Orders ID should not be nil") &&
						suite.False(
							actual.Shipment.MoveTaskOrder.Orders.UploadedOrders.ID.IsNil(),
							"UploadedOrders ID should not be nil",
						) &&
						suite.False(
							actual.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0].ID.IsNil(),
							"Expected uploaded orders user uploads to be loaded",
						) &&
						suite.False(
							actual.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0].Upload.ID.IsNil(),
							"Expected uploaded orders to be loaded",
						) {
						suite.Equal(
							expected.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0].Upload.ID,
							actual.Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads[0].Upload.ID,
						)
					}
				},
			},
		}

		for name, testCase := range postloadAssociationTestCases {
			name, testCase := name, testCase

			suite.Run(fmt.Sprintf("Can load %s", name), func() {
				ppmShipment := factory.BuildPPMShipmentWithAllDocTypesApproved(
					suite.DB(),
					userUploader,
				)

				// Fetch the shipment fresh from the DB because the ppmShipment var already has all the associations
				// loaded
				ppmShipmentReturned, err := fetcher.GetPPMShipment(
					suite.AppContextForTest(),
					ppmShipment.ID,
					GetListOfAllPreloadAssociations(),
					nil,
				)

				suite.FatalNoError(err, "failed to fetch PPM Shipment")

				err = fetcher.PostloadAssociations(
					suite.AppContextForTest(),
					ppmShipmentReturned,
					testCase.postloadAssociations,
				)

				if suite.NoError(err) {
					testCase.successAssertionFunc(&ppmShipment, ppmShipmentReturned)
				}
			})
		}

		suite.Run("Excludes deleted uploads", func() {
			ppmShipment := factory.BuildPPMShipmentThatNeedsCloseoutWithAllDocTypes(
				suite.DB(),
				userUploader,
			)

			// Create a deleted upload for a weight ticket
			originalWeightTicket := ppmShipment.WeightTickets[0]
			numValidEmptyUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
			suite.FatalTrue(numValidEmptyUploads > 0)

			now := time.Now()
			deletedWeightTicketUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    originalWeightTicket.EmptyDocument,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{
						DeletedAt: &now,
					},
					ExtendedParams: &factory.UserUploadExtendedParams{
						AppContext: suite.AppContextForTest(),
					},
				},
				{
					Model: models.Upload{
						DeletedAt: &now,
					},
				},
			}, nil)

			suite.FatalNotNil(deletedWeightTicketUpload.Upload.DeletedAt)
			suite.FatalNotNil(deletedWeightTicketUpload.DeletedAt)

			// Create a deleted upload for a progear weight ticket
			originalProgearWeightTicket := ppmShipment.ProgearWeightTickets[0]
			numValidProgearWeightTicketUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
			suite.FatalTrue(numValidProgearWeightTicketUploads > 0)

			deletedProgearWeightTicketUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    originalProgearWeightTicket.Document,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{
						DeletedAt: &now,
					},
					ExtendedParams: &factory.UserUploadExtendedParams{
						AppContext: suite.AppContextForTest(),
					},
				},
				{
					Model: models.Upload{
						DeletedAt: &now,
					},
				},
			}, nil)

			// Create a deleted upload for a moving expense
			originalMovingExpense := ppmShipment.MovingExpenses[0]
			numValidMovingExpenseUploads := len(originalWeightTicket.EmptyDocument.UserUploads)
			suite.FatalTrue(numValidMovingExpenseUploads > 0)

			deletedMovingExpenseUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
				{
					Model:    originalMovingExpense.Document,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{
						DeletedAt: &now,
					},
					ExtendedParams: &factory.UserUploadExtendedParams{
						AppContext: suite.AppContextForTest(),
					},
				},
				{
					Model: models.Upload{
						DeletedAt: &now,
					},
				},
			}, nil)

			// Fetch the shipment fresh from the DB because the ppmShipment var already has all the associations
			// loaded
			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				GetListOfAllPreloadAssociations(),
				nil,
			)

			suite.FatalNoError(err, "failed to fetch PPM Shipment")

			err = fetcher.PostloadAssociations(
				suite.AppContextForTest(),
				ppmShipmentReturned,
				GetListOfAllPostloadAssociations(),
			)

			if suite.NoError(err) {
				suite.Equal(len(ppmShipment.WeightTickets), len(ppmShipmentReturned.WeightTickets))

				suite.Equal(originalWeightTicket.ID, ppmShipmentReturned.WeightTickets[0].ID)
				retrievedWeightTicket := ppmShipmentReturned.WeightTickets[0]

				if suite.Equal(numValidEmptyUploads, len(retrievedWeightTicket.EmptyDocument.UserUploads)) {
					for _, upload := range retrievedWeightTicket.EmptyDocument.UserUploads {
						suite.NotEqual(deletedWeightTicketUpload.ID, upload.ID)
						suite.Nil(upload.DeletedAt)
					}
				}

				suite.Equal(len(ppmShipment.MovingExpenses), len(ppmShipmentReturned.MovingExpenses))

				suite.Equal(originalMovingExpense.ID, ppmShipmentReturned.MovingExpenses[0].ID)
				retrievedMovingExpense := ppmShipmentReturned.MovingExpenses[0]

				if suite.Equal(numValidMovingExpenseUploads, len(retrievedMovingExpense.Document.UserUploads)) {
					for _, upload := range retrievedMovingExpense.Document.UserUploads {
						suite.NotEqual(deletedMovingExpenseUpload.ID, upload.ID)
						suite.Nil(upload.DeletedAt)
					}
				}

				suite.Equal(len(ppmShipment.ProgearWeightTickets), len(ppmShipmentReturned.ProgearWeightTickets))

				suite.Equal(originalProgearWeightTicket.ID, ppmShipmentReturned.ProgearWeightTickets[0].ID)
				retrievedProgearWeightTicket := ppmShipmentReturned.ProgearWeightTickets[0]

				if suite.Equal(
					numValidProgearWeightTicketUploads,
					len(retrievedProgearWeightTicket.Document.UserUploads),
				) {
					for _, upload := range retrievedProgearWeightTicket.Document.UserUploads {
						suite.NotEqual(deletedProgearWeightTicketUpload.ID, upload.ID)
						suite.Nil(upload.DeletedAt)
					}
				}
			}
		})

		suite.Run("Returns an error if an invalid association is passed in", func() {
			appCtx := suite.AppContextForTest()

			ppmShipment := factory.BuildPPMShipment(appCtx.DB(), nil, nil)

			invalidAssociation := "invalid_association"

			// Fetch the shipment fresh from the DB because the ppmShipment var already has all the associations
			err := fetcher.PostloadAssociations(appCtx, &ppmShipment, []string{invalidAssociation})

			if suite.Error(err) {
				suite.IsType(apperror.NotImplementedError{}, err)

				suite.Contains(
					err.Error(),
					fmt.Sprintf("Requested post load association %s is not implemented", invalidAssociation),
				)
			}
		})
	})
}

func (suite *PPMShipmentSuite) TestFetchPPMShipment() {
	fetcher := NewPPMShipmentFetcher()
	suite.Run("FindPPMShipmentWithDocument - document belongs to weight ticket", func() {
		weightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.EmptyDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for empty weight document")

		err = FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.FullDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for full weight document")

		err = FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.FullDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for trailer ownership document")
	})

	suite.Run("FindPPMShipmentWithDocument - document belongs to pro gear", func() {
		proGear := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), proGear.PPMShipmentID, proGear.DocumentID)
		suite.NoError(err, "expected to find PPM Shipment for weight document")
	})

	suite.Run("FindPPMShipmentWithDocument - document belongs to moving expenses", func() {
		movingExpense := factory.BuildMovingExpense(suite.DB(), nil, nil)

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), movingExpense.PPMShipmentID, movingExpense.DocumentID)
		suite.NoError(err, "expected to find PPM Shipment for moving expense document")
	})

	suite.Run("FindPPMShipmentWithDocument - document not found", func() {
		weightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)
		factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    weightTicket.PPMShipment,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    weightTicket.PPMShipment,
				LinkOnly: true,
			},
		}, nil)

		documentID := uuid.Must(uuid.NewV4())
		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, documentID)
		suite.Error(err, "expected to return not found error for unknown document id")
		suite.Equal(fmt.Sprintf("ID: %s not found document does not exist for the given shipment", documentID), err.Error())
	})

	suite.Run("FindPPMShipment - loads weight tickets association", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		// No uploads are added by default for the ProofOfTrailerOwnershipDocument to the WeightTicket model
		testdatagen.GetOrCreateDocumentWithUploads(suite.DB(),
			ppmShipment.WeightTickets[0].ProofOfTrailerOwnershipDocument,
			testdatagen.Assertions{ServiceMember: ppmShipment.WeightTickets[0].EmptyDocument.ServiceMember})

		actualShipment, err := FindPPMShipment(suite.AppContextForTest(), ppmShipment.ID)
		suite.NoError(err)

		suite.Len(actualShipment.WeightTickets, 1)
		suite.NotEmpty(actualShipment.WeightTickets[0].EmptyDocument.UserUploads[0].Upload)
		suite.NotEmpty(actualShipment.WeightTickets[0].FullDocument.UserUploads[0].Upload)
		suite.NotEmpty(actualShipment.WeightTickets[0].ProofOfTrailerOwnershipDocument.UserUploads[0].Upload)
	})

	suite.Run("FindPPMShipment - loads ProgearWeightTicket and MovingExpense associations", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		actualShipment, err := FindPPMShipment(suite.AppContextForTest(), ppmShipment.ID)
		suite.NoError(err)

		suite.Len(actualShipment.ProgearWeightTickets, 1)
		suite.NotEmpty(actualShipment.ProgearWeightTickets[0].Document.UserUploads[0].Upload)

		suite.Len(actualShipment.MovingExpenses, 1)
		suite.NotEmpty(actualShipment.MovingExpenses[0].Document.UserUploads[0].Upload)
	})

	suite.Run("FindPPMShipment - loads signed certification", func() {
		signedCertification := factory.BuildSignedCertification(suite.DB(), nil, nil)
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)
		signedCertification.PpmID = &ppmShipment.ID
		suite.NoError(suite.DB().Save(&signedCertification))

		actualShipment, err := FindPPMShipment(suite.AppContextForTest(), *signedCertification.PpmID)
		suite.NoError(err)

		if actualCertification := actualShipment.SignedCertification; suite.NotNil(actualCertification.ID) {
			suite.Equal(signedCertification.ID, actualCertification.ID)
			suite.Equal(signedCertification.CertificationText, actualCertification.CertificationText)
			suite.Equal(signedCertification.CertificationType, actualCertification.CertificationType)
			suite.True(signedCertification.Date.UTC().Truncate(time.Millisecond).
				Equal(actualCertification.Date.UTC().Truncate(time.Millisecond)))
			suite.Equal(signedCertification.MoveID, actualCertification.MoveID)
			suite.Equal(signedCertification.PpmID, actualCertification.PpmID)
			suite.Equal(signedCertification.Signature, actualCertification.Signature)
			suite.Equal(signedCertification.SubmittingUserID, actualCertification.SubmittingUserID)
		}
	})

	suite.Run("FindPPMShipment - returns not found for unknown id", func() {
		badID := uuid.Must(uuid.NewV4())
		_, err := FindPPMShipment(suite.AppContextForTest(), badID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", badID), err.Error())
	})

	suite.Run("FindPPMShipment - returns not found for deleted shipment", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		err := utilities.SoftDestroy(suite.DB(), &ppmShipment)
		suite.NoError(err)

		_, err = FindPPMShipment(suite.AppContextForTest(), ppmShipment.ID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", ppmShipment.ID), err.Error())
	})

	suite.Run("FindPPMShipment - deleted uploads are removed", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		testdatagen.GetOrCreateDocumentWithUploads(suite.DB(),
			ppmShipment.WeightTickets[0].ProofOfTrailerOwnershipDocument,
			testdatagen.Assertions{ServiceMember: ppmShipment.WeightTickets[0].EmptyDocument.ServiceMember})

		factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		actualShipment, err := FindPPMShipment(suite.AppContextForTest(), ppmShipment.ID)
		suite.NoError(err)

		suite.Len(actualShipment.WeightTickets[0].EmptyDocument.UserUploads, 1)
		suite.Len(actualShipment.WeightTickets[0].FullDocument.UserUploads, 1)
		suite.Len(actualShipment.WeightTickets[0].ProofOfTrailerOwnershipDocument.UserUploads, 1)
		suite.Len(actualShipment.ProgearWeightTickets[0].Document.UserUploads, 1)
		suite.Len(actualShipment.MovingExpenses[0].Document.UserUploads, 1)

		err = utilities.SoftDestroy(suite.DB(), &actualShipment.WeightTickets[0].EmptyDocument.UserUploads[0])
		suite.NoError(err)

		err = utilities.SoftDestroy(suite.DB(), &actualShipment.WeightTickets[0].FullDocument.UserUploads[0])
		suite.NoError(err)

		err = utilities.SoftDestroy(suite.DB(), &actualShipment.WeightTickets[0].ProofOfTrailerOwnershipDocument.UserUploads[0])
		suite.NoError(err)

		err = utilities.SoftDestroy(suite.DB(), &actualShipment.ProgearWeightTickets[0].Document.UserUploads[0])
		suite.NoError(err)

		err = utilities.SoftDestroy(suite.DB(), &actualShipment.MovingExpenses[0].Document.UserUploads[0])
		suite.NoError(err)

		actualShipment, err = FindPPMShipment(suite.AppContextForTest(), ppmShipment.ID)
		suite.NoError(err)

		suite.Len(actualShipment.WeightTickets[0].EmptyDocument.UserUploads, 0)
		suite.Len(actualShipment.WeightTickets[0].FullDocument.UserUploads, 0)
		suite.Len(actualShipment.WeightTickets[0].ProofOfTrailerOwnershipDocument.UserUploads, 0)
		suite.Len(actualShipment.ProgearWeightTickets[0].Document.UserUploads, 0)
		suite.Len(actualShipment.MovingExpenses[0].Document.UserUploads, 0)
	})

	suite.Run("FetchPPMShipmentFromMTOShipmentID - finds records", func() {
		ppm := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		retrievedPPM, _ := FetchPPMShipmentFromMTOShipmentID(suite.AppContextForTest(), ppm.ShipmentID)

		suite.Equal(retrievedPPM.ID, ppm.ID)
		suite.Equal(retrievedPPM.ShipmentID, ppm.ShipmentID)

	})

	suite.Run("FetchPPMShipmentFromMTOShipmentID  - returns not found for unknown id", func() {
		badID := uuid.Must(uuid.NewV4())
		_, err := FetchPPMShipmentFromMTOShipmentID(suite.AppContextForTest(), badID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", badID), err.Error())
	})

	suite.Run("FetchPPMShipmentFromMTOShipmentID  - returns not found for deleted shipment", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)

		err := utilities.SoftDestroy(suite.DB(), &ppmShipment)
		suite.NoError(err)

		_, err = FetchPPMShipmentFromMTOShipmentID(suite.AppContextForTest(), ppmShipment.ShipmentID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", ppmShipment.ShipmentID), err.Error())
	})

	suite.Run("FindPPMShipmentAndWeightTickets - Success", func() {
		weightTicket := factory.BuildWeightTicket(suite.DB(), nil, nil)
		foundPPMShipment, err := FindPPMShipmentAndWeightTickets(suite.AppContextForTest(), weightTicket.PPMShipmentID)

		suite.Nil(err)
		suite.Equal(weightTicket.PPMShipmentID, foundPPMShipment.ID)
		suite.Equal(weightTicket.PPMShipment.Status, foundPPMShipment.Status)
		suite.Len(foundPPMShipment.WeightTickets, 1)
		suite.Equal(*weightTicket.EmptyWeight, *foundPPMShipment.WeightTickets[0].EmptyWeight)
		suite.Equal(*weightTicket.FullWeight, *foundPPMShipment.WeightTickets[0].FullWeight)
	})

	suite.Run("FindPPMShipmentAndWeightTickets - still returns if weightTicket does not exist", func() {
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), nil, nil)
		foundPPMShipment, err := FindPPMShipmentAndWeightTickets(suite.AppContextForTest(), ppmShipment.ID)

		suite.Nil(err)
		suite.Equal(ppmShipment.ID, foundPPMShipment.ID)
		suite.Equal(ppmShipment.ShipmentID, foundPPMShipment.ShipmentID)
	})

	suite.Run("FindPPMShipmentAndWeightTickets - errors if ID isn't found", func() {
		id := uuid.Must(uuid.NewV4())
		foundPPMShipment, err := FindPPMShipmentAndWeightTickets(suite.AppContextForTest(), id)

		suite.Nil(foundPPMShipment)
		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for PPMShipmentAndWeightTickets", id.String()),
				err.Error(),
			)
		}
	})

	suite.Run("FindPPMShipmentByMTOID - Success deleted line items are excluded", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, nil)

		weightTicketToDelete := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		err := utilities.SoftDestroy(suite.DB(), &weightTicketToDelete)
		suite.NoError(err)

		factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		proGearToDelete := factory.BuildProgearWeightTicket(suite.DB(),
			[]factory.Customization{
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
			}, nil)

		err = utilities.SoftDestroy(suite.DB(), &proGearToDelete)
		suite.NoError(err)

		factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		movingExpenseToDelete := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		err = utilities.SoftDestroy(suite.DB(), &movingExpenseToDelete)
		suite.NoError(err)

		actualShipment, err := FindPPMShipmentByMTOID(suite.AppContextForTest(), ppmShipment.ShipmentID)
		suite.NoError(err)

		suite.Len(actualShipment.WeightTickets, 1)
		suite.Len(actualShipment.ProgearWeightTickets, 1)
		suite.Len(actualShipment.MovingExpenses, 1)
	})

	suite.Run("GetPPMShipment filters rejected weight tickets", func() {
		ppmShipment := factory.BuildPPMShipmentWithAllDocTypesApproved(suite.DB(), nil)
		rejectedStatus := models.PPMDocumentStatusRejected
		rejectedWeightTicket := factory.BuildWeightTicket(suite.DB(), []factory.Customization{
			{
				Model: models.WeightTicket{
					Status: &rejectedStatus,
				},
			},
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		result, err := fetcher.GetPPMShipment(
			suite.AppContextForTest(),
			ppmShipment.ID,
			nil,
			nil,
		)

		suite.NoError(err)
		suite.NotNil(result)
		for _, wt := range result.WeightTickets {
			suite.NotEqual(rejectedWeightTicket.ID, wt.ID)
			suite.NotEqual(models.PPMDocumentStatusRejected, *wt.Status)
		}
	})
	suite.Run("GetPPMShipment filters rejected moving expenses", func() {
		ppmShipment := factory.BuildPPMShipmentWithAllDocTypesApproved(suite.DB(), nil)

		rejectedStatus := models.PPMDocumentStatusRejected
		rejectedMovingExpense := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model: models.MovingExpense{
					Status: &rejectedStatus,
				},
			},
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		result, err := fetcher.GetPPMShipment(
			suite.AppContextForTest(),
			ppmShipment.ID,
			nil,
			nil,
		)

		suite.NoError(err)
		suite.NotNil(result)
		for _, exp := range result.MovingExpenses {
			suite.NotEqual(rejectedMovingExpense.ID, exp.ID)
			suite.NotEqual(models.PPMDocumentStatusRejected, *exp.Status)
		}
	})
}
