package reweigh

import (
	"testing"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ReweighSuite) TestReweighUpdater() {
	reweighUpdater := NewReweighUpdater()
	oldReweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{})
	eTag := etag.GenerateEtag(oldReweigh.UpdatedAt)
	newReweigh := oldReweigh

	appCtx := appcontext.NewAppContext(suite.DB(), suite.logger)

	// Test Success - Reweigh updated
	suite.T().Run("Updated reweigh - Success", func(t *testing.T) {
		newWeight := unit.Pound(200)
		newReweigh.Weight = &newWeight
		updatedReweigh, err := reweighUpdater.UpdateReweigh(appCtx, &newReweigh, eTag)

		suite.NoError(err)
		suite.NotNil(updatedReweigh)
		suite.Equal(newWeight, *updatedReweigh.Weight)
		eTag = etag.GenerateEtag(updatedReweigh.UpdatedAt)
	})
	// Test NotFoundError
	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundReweigh := newReweigh
		notFoundReweigh.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedReweigh, err := reweighUpdater.UpdateReweigh(appCtx, &notFoundReweigh, eTag)

		suite.Nil(updatedReweigh)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})
	// PreconditionFailedError
	suite.T().Run("Precondition Failed", func(t *testing.T) {
		updatedReweigh, err := reweighUpdater.UpdateReweigh(appCtx, &newReweigh, "nada") // base validation

		suite.Nil(updatedReweigh)
		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})
	// InvalidInputError
	suite.T().Run("Reweigh with validation errors returns an InvalidInputError", func(t *testing.T) {
		badRequestedby := models.ReweighRequester("not requested by anyone")
		newReweigh.RequestedBy = badRequestedby
		updatedReweigh, err := reweighUpdater.UpdateReweigh(appCtx, &newReweigh, eTag)

		suite.Error(err)
		suite.Nil(updatedReweigh)
		suite.IsType(services.InvalidInputError{}, err)
	})
}
