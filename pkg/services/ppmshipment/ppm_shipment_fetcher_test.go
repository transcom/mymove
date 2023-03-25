package ppmshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
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

	suite.Run("GetPPMShipment", func() {
		suite.Run("Can fetch a PPM Shipment", func() {
			ppmShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				nil,
			)

			if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
				suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)
			}
		})

		associationTestCases := map[string]struct {
			eagerPreloadAssociations []string
			successAssertionFunc     func(*models.PPMShipment, *models.PPMShipment)
		}{
			"No associations": {
				eagerPreloadAssociations: nil,
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.True(actual.Shipment.ID.IsNil())
					suite.Nil(actual.WeightTickets)
					suite.Nil(actual.ProgearWeightTickets)
					suite.Nil(actual.MovingExpenses)
					suite.Nil(actual.W2Address)
					suite.Nil(actual.AOAPacket)
					suite.Nil(actual.PaymentPacket)
				},
			},
			"Shipment only": {
				eagerPreloadAssociations: []string{EagerPreloadAssociationShipment},
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.False(actual.Shipment.ID.IsNil())
					suite.Equal(expected.Shipment.ID, actual.Shipment.ID)

					suite.Nil(actual.WeightTickets)
					suite.Nil(actual.ProgearWeightTickets)
					suite.Nil(actual.MovingExpenses)
					suite.Nil(actual.W2Address)
					suite.Nil(actual.AOAPacket)
					suite.Nil(actual.PaymentPacket)
				},
			},
			"Weight tickets only": {
				eagerPreloadAssociations: []string{EagerPreloadAssociationWeightTickets},
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.True(actual.Shipment.ID.IsNil())

					suite.NotNil(actual.WeightTickets)
					suite.Equal(len(expected.WeightTickets), len(actual.WeightTickets))
					suite.Equal(expected.WeightTickets[0].ID, actual.WeightTickets[0].ID)

					suite.Nil(actual.ProgearWeightTickets)
					suite.Nil(actual.MovingExpenses)
					suite.Nil(actual.W2Address)
					suite.Nil(actual.AOAPacket)
					suite.Nil(actual.PaymentPacket)
				},
			},
			"All eager preload associations": {
				eagerPreloadAssociations: GetListOfAllPreloadAssociations(),
				successAssertionFunc: func(expected *models.PPMShipment, actual *models.PPMShipment) {
					suite.False(actual.Shipment.ID.IsNil())
					suite.Equal(expected.Shipment.ID, actual.Shipment.ID)

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
				ppmShipment := testdatagen.MakePPMShipmentWithAllDocTypesApproved(
					suite.DB(),
					testdatagen.Assertions{
						UserUploader: userUploader,
					},
				)

				ppmShipmentReturned, err := fetcher.GetPPMShipment(
					suite.AppContextForTest(),
					ppmShipment.ID,
					testCase.eagerPreloadAssociations,
				)

				if suite.NoError(err) && suite.NotNil(ppmShipmentReturned) {
					suite.Equal(ppmShipment.ID, ppmShipmentReturned.ID)
					suite.Equal(ppmShipment.ShipmentID, ppmShipmentReturned.ShipmentID)

					testCase.successAssertionFunc(&ppmShipment, ppmShipmentReturned)
				}
			})
		}

		suite.Run("Returns a not found error if the PPM Shipment does not exist", func() {
			nonexistentID := uuid.Must(uuid.NewV4())

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				nonexistentID,
				nil,
			)

			if suite.Error(err) && suite.Nil(ppmShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", nonexistentID), err.Error())
			}
		})

		suite.Run("Returns an error if an invalid association is requested", func() {
			ppmShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())

			invalidAssociation := "invalid"
			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				[]string{invalidAssociation},
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
			ppmShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())

			err := utilities.SoftDestroy(suite.DB(), &ppmShipment)
			suite.FatalNoError(err)

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				nil,
			)

			if suite.Error(err) && suite.Nil(ppmShipmentReturned) {
				suite.IsType(apperror.NotFoundError{}, err)

				suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", ppmShipment.ID), err.Error())
			}
		})

		suite.Run("Excludes deleted documents", func() {
			ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(
				suite.DB(),
				testdatagen.Assertions{
					UserUploader: userUploader,
				},
			)

			// create new ppm documents that are deleted
			now := time.Now()

			testdatagen.AddWeightTicketToPPMShipment(suite.DB(), &ppmShipment, testdatagen.Assertions{
				UserUploader: userUploader,
				WeightTicket: models.WeightTicket{
					DeletedAt: &now,
				},
			})

			testdatagen.AddProgearWeightTicketToPPMShipment(suite.DB(), &ppmShipment, testdatagen.Assertions{
				UserUploader: userUploader,
				ProgearWeightTicket: models.ProgearWeightTicket{
					DeletedAt: &now,
				},
			})

			testdatagen.AddMovingExpenseToPPMShipment(suite.DB(), &ppmShipment, testdatagen.Assertions{
				UserUploader: userUploader,
				MovingExpense: models.MovingExpense{
					DeletedAt: &now,
				},
			})

			ppmShipmentReturned, err := fetcher.GetPPMShipment(
				suite.AppContextForTest(),
				ppmShipment.ID,
				[]string{
					EagerPreloadAssociationWeightTickets,
					EagerPreloadAssociationProgearWeightTickets,
					EagerPreloadAssociationMovingExpenses,
				},
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
	})
}

func (suite *PPMShipmentSuite) TestFetchPPMShipment() {
	suite.Run("FindPPMShipmentWithDocument - document belongs to weight ticket", func() {
		weightTicket := testdatagen.MakeDefaultWeightTicket(suite.DB())

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.EmptyDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for empty weight document")

		err = FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.FullDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for full weight document")

		err = FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.FullDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for trailer ownership document")
	})

	suite.Run("FindPPMShipmentWithDocument - document belongs to pro gear", func() {
		proGear := testdatagen.MakeDefaultProgearWeightTicket(suite.DB())

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), proGear.PPMShipmentID, proGear.DocumentID)
		suite.NoError(err, "expected to find PPM Shipment for weight document")
	})

	suite.Run("FindPPMShipmentWithDocument - document belongs to moving expenses", func() {
		movingExpense := testdatagen.MakeDefaultMovingExpense(suite.DB())

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), movingExpense.PPMShipmentID, movingExpense.DocumentID)
		suite.NoError(err, "expected to find PPM Shipment for moving expense document")
	})

	suite.Run("FindPPMShipmentWithDocument - document not found", func() {
		weightTicket := testdatagen.MakeDefaultWeightTicket(suite.DB())
		testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{PPMShipment: weightTicket.PPMShipment})
		testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{PPMShipment: weightTicket.PPMShipment})

		documentID := uuid.Must(uuid.NewV4())
		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, documentID)
		suite.Error(err, "expected to return not found error for unknown document id")
		suite.Equal(fmt.Sprintf("ID: %s not found document does not exist for the given shipment", documentID), err.Error())
	})

	suite.Run("FindPPMShipment - loads weight tickets association", func() {
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), testdatagen.Assertions{})

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
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), testdatagen.Assertions{})

		testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		actualShipment, err := FindPPMShipment(suite.AppContextForTest(), ppmShipment.ID)
		suite.NoError(err)

		suite.Len(actualShipment.ProgearWeightTickets, 1)
		suite.NotEmpty(actualShipment.ProgearWeightTickets[0].Document.UserUploads[0].Upload)

		suite.Len(actualShipment.MovingExpenses, 1)
		suite.NotEmpty(actualShipment.MovingExpenses[0].Document.UserUploads[0].Upload)
	})

	suite.Run("FindPPMShipment - loads signed certification", func() {
		signedCertification := testdatagen.MakeSignedCertificationForPPM(suite.DB(), testdatagen.Assertions{})

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
			suite.Nil(actualCertification.PersonallyProcuredMoveID)
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
		ppmShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

		err := utilities.SoftDestroy(suite.DB(), &ppmShipment)
		suite.NoError(err)

		_, err = FindPPMShipment(suite.AppContextForTest(), ppmShipment.ID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", ppmShipment.ID), err.Error())
	})

	suite.Run("FindPPMShipment - deleted uploads are removed", func() {
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), testdatagen.Assertions{})

		testdatagen.GetOrCreateDocumentWithUploads(suite.DB(),
			ppmShipment.WeightTickets[0].ProofOfTrailerOwnershipDocument,
			testdatagen.Assertions{ServiceMember: ppmShipment.WeightTickets[0].EmptyDocument.ServiceMember})

		testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

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
		ppm := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

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
		ppmShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

		err := utilities.SoftDestroy(suite.DB(), &ppmShipment)
		suite.NoError(err)

		_, err = FetchPPMShipmentFromMTOShipmentID(suite.AppContextForTest(), ppmShipment.ShipmentID)
		suite.Error(err)

		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for PPMShipment", ppmShipment.ShipmentID), err.Error())
	})

	suite.Run("FindPPMShipmentAndWeightTickets - Success", func() {
		weightTicket := testdatagen.MakeDefaultWeightTicket(suite.DB())
		foundPPMShipment, err := FindPPMShipmentAndWeightTickets(suite.AppContextForTest(), weightTicket.PPMShipmentID)

		suite.Nil(err)
		suite.Equal(weightTicket.PPMShipmentID, foundPPMShipment.ID)
		suite.Equal(weightTicket.PPMShipment.Status, foundPPMShipment.Status)
		suite.Len(foundPPMShipment.WeightTickets, 1)
		suite.Equal(*weightTicket.EmptyWeight, *foundPPMShipment.WeightTickets[0].EmptyWeight)
		suite.Equal(*weightTicket.FullWeight, *foundPPMShipment.WeightTickets[0].FullWeight)
	})

	suite.Run("FindPPMShipmentAndWeightTickets - still returns if weightTicket does not exist", func() {
		ppmShipment := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})
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
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), testdatagen.Assertions{})

		weightTicketToDelete := testdatagen.MakeWeightTicket(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		err := utilities.SoftDestroy(suite.DB(), &weightTicketToDelete)
		suite.NoError(err)

		testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		proGearToDelete := testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		err = utilities.SoftDestroy(suite.DB(), &proGearToDelete)
		suite.NoError(err)

		testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		movingExpenseToDelete := testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
			PPMShipment: ppmShipment,
		})

		err = utilities.SoftDestroy(suite.DB(), &movingExpenseToDelete)
		suite.NoError(err)

		actualShipment, err := FindPPMShipmentByMTOID(suite.AppContextForTest(), ppmShipment.ShipmentID)
		suite.NoError(err)

		suite.Len(actualShipment.WeightTickets, 1)
		suite.Len(actualShipment.ProgearWeightTickets, 1)
		suite.Len(actualShipment.MovingExpenses, 1)
	})
}
