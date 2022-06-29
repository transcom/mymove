package weightticket

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite WeightTicketSuite) TestUpdateWeightTicket() {
	setupForTest := func(appCtx appcontext.AppContext, overrides *models.WeightTicket) *models.WeightTicket {
		serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
		ppmShipment := testdatagen.MakeMinimalDefaultPPMShipment(suite.DB())

		baseDocumentAssertions := testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
			},
		}

		emptyDocument := testdatagen.MakeDocument(appCtx.DB(), baseDocumentAssertions)
		fullDocument := testdatagen.MakeDocument(appCtx.DB(), baseDocumentAssertions)
		proofOfOwnership := testdatagen.MakeDocument(appCtx.DB(), baseDocumentAssertions)

		originalWeightTicket := models.WeightTicket{
			EmptyDocumentID:                   emptyDocument.ID,
			FullDocumentID:                    fullDocument.ID,
			ProofOfTrailerOwnershipDocumentID: proofOfOwnership.ID,
			PPMShipmentID:                     ppmShipment.ID,
		}

		if overrides != nil {
			testdatagen.MergeModels(&originalWeightTicket, overrides)
		}

		verrs, err := appCtx.DB().ValidateAndCreate(&originalWeightTicket)

		suite.NoVerrs(verrs)
		suite.Nil(err)
		suite.NotNil(originalWeightTicket.ID)

		return &originalWeightTicket
	}

	suite.Run("Returns an error if the original doesn't exist", func() {
		badWeightTicket := models.WeightTicket{
			ID: uuid.Must(uuid.NewV4()),
		}

		updater := NewCustomerWeightTicketUpdater()

		updatedWeightTicket, err := updater.UpdateWeightTicket(suite.AppContextForTest(), badWeightTicket, "")

		suite.Nil(updatedWeightTicket)

		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)

			suite.Equal(
				fmt.Sprintf("ID: %s not found while looking for WeightTicket", badWeightTicket.ID.String()),
				err.Error(),
			)
		}
	})

	suite.Run("Returns a PreconditionFailedError if the input eTag is stale/incorrect", func() {
		appCtx := suite.AppContextForTest()

		originalWeightTicket := setupForTest(appCtx, nil)

		updater := NewCustomerWeightTicketUpdater()

		updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *originalWeightTicket, "")

		suite.Nil(updatedWeightTicket)

		if suite.Error(updateErr) {
			suite.IsType(apperror.PreconditionFailedError{}, updateErr)

			suite.Equal(
				fmt.Sprintf("Precondition failed on update to object with ID: '%s'. The If-Match header value did not match the eTag for this record.", originalWeightTicket.ID.String()),
				updateErr.Error(),
			)
		}
	})

	suite.Run("Successfully updates", func() {
		appCtx := suite.AppContextForTest()

		originalWeightTicket := setupForTest(appCtx, nil)

		updater := NewCustomerWeightTicketUpdater()

		desiredWeightTicket := &models.WeightTicket{
			ID:                       originalWeightTicket.ID,
			VehicleDescription:       models.StringPointer("2004 Toyota Prius"),
			EmptyWeight:              models.PoundPointer(3000),
			MissingEmptyWeightTicket: models.BoolPointer(true),
			FullWeight:               models.PoundPointer(4200),
			MissingFullWeightTicket:  models.BoolPointer(true),
			OwnsTrailer:              models.BoolPointer(false),
			TrailerMeetsCriteria:     models.BoolPointer(false),
		}

		updatedWeightTicket, updateErr := updater.UpdateWeightTicket(appCtx, *desiredWeightTicket, etag.GenerateEtag(originalWeightTicket.UpdatedAt))

		suite.Nil(updateErr)
		suite.Equal(originalWeightTicket.ID, updatedWeightTicket.ID)
		suite.Equal(originalWeightTicket.EmptyDocumentID, updatedWeightTicket.EmptyDocumentID)
		suite.Equal(originalWeightTicket.FullDocumentID, updatedWeightTicket.FullDocumentID)
		suite.Equal(originalWeightTicket.ProofOfTrailerOwnershipDocumentID, updatedWeightTicket.ProofOfTrailerOwnershipDocumentID)
		suite.Equal(*desiredWeightTicket.VehicleDescription, *updatedWeightTicket.VehicleDescription)
		suite.Equal(*desiredWeightTicket.EmptyWeight, *updatedWeightTicket.EmptyWeight)
		suite.Equal(*desiredWeightTicket.MissingEmptyWeightTicket, *updatedWeightTicket.MissingEmptyWeightTicket)
		suite.Equal(*desiredWeightTicket.FullWeight, *updatedWeightTicket.FullWeight)
		suite.Equal(*desiredWeightTicket.MissingFullWeightTicket, *updatedWeightTicket.MissingFullWeightTicket)
		suite.Equal(*desiredWeightTicket.OwnsTrailer, *updatedWeightTicket.OwnsTrailer)
		suite.Equal(*desiredWeightTicket.TrailerMeetsCriteria, *updatedWeightTicket.TrailerMeetsCriteria)
	})
}
