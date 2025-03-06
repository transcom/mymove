package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPPMCloseoutSummaryValidation() {
	suite.Run("Test Valid PPM Closeout Summary", func() {
		validPPMCloseoutSummary := models.PPMCloseoutSummary{
			ID:            uuid.Must(uuid.NewV4()),
			PPMShipmentID: uuid.Must(uuid.NewV4()),
			MaxAdvance:    models.CentPointer(100000),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPPMCloseoutSummary, expErrors)
	})

	suite.Run("Test missing PPMShipmentID", func() {
		invalidPPMCloseoutSummary := models.PPMCloseoutSummary{
			ID:         uuid.Must(uuid.NewV4()),
			MaxAdvance: models.CentPointer(100000),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		expErrors := map[string][]string{
			"ppmshipment_id": {"PPMShipmentID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPPMCloseoutSummary, expErrors)
	})
}
