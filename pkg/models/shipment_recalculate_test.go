package models_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestShipmentRecalculate() {
	testCases := map[string]struct {
		recalculateDates models.ShipmentRecalculate
		expectedErrs     map[string][]string
	}{
		"Successful Create": {
			recalculateDates: models.ShipmentRecalculate{
				ID:                    uuid.Must(uuid.NewV4()),
				ShipmentUpdatedAfter:  time.Date(1970, time.January, 01, 0, 0, 0, 0, time.UTC),
				ShipmentUpdatedBefore: time.Now(),
				Active:                true,
			},
			expectedErrs: nil,
		},

		"Empty Fields": {
			recalculateDates: models.ShipmentRecalculate{},
			expectedErrs: map[string][]string{
				"shipment_updated_before": {"ShipmentUpdatedBefore can not be blank."},
				"shipment_updated_after":  {"ShipmentUpdatedAfter can not be blank."},
			},
		},

		"Bad recalculate date range": {
			recalculateDates: models.ShipmentRecalculate{
				ID:                    uuid.Must(uuid.NewV4()),
				ShipmentUpdatedBefore: time.Date(1970, time.January, 01, 0, 0, 0, 0, time.UTC),
				ShipmentUpdatedAfter:  time.Now(),
				Active:                true,
			},
			expectedErrs: map[string][]string{
				"shipment_updated_before": {"ShipmentUpdatedBefore must be after ShipmentUpdatedAfter."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.recalculateDates, test.expectedErrs)
		})
	}
}

func (suite *ModelSuite) TestShipmentRecalculateTooManyActiveRecord() {

	id := uuid.Must(uuid.NewV4())

	suite.T().Run("Too many active records", func(t *testing.T) {

		// Test for too many active records
		newRecalculateDates := models.ShipmentRecalculate{
			ID:                    id,
			ShipmentUpdatedAfter:  time.Date(1970, time.January, 01, 0, 0, 0, 0, time.UTC),
			ShipmentUpdatedBefore: time.Date(1970, time.February, 01, 0, 0, 0, 0, time.UTC),
			Active:                true,
		}
		verrs, err := suite.DB().ValidateAndCreate(&newRecalculateDates)
		suite.Empty(verrs.Error())
		suite.Nil(err)

		id = uuid.Must(uuid.NewV4())
		newRecalculateDates = models.ShipmentRecalculate{
			ShipmentUpdatedAfter:  time.Date(testdatagen.TestYear, time.January, 01, 0, 0, 0, 0, time.UTC),
			ShipmentUpdatedBefore: time.Date(testdatagen.TestYear, time.March, 01, 0, 0, 0, 0, time.UTC),
			Active:                true,
		}
		verrs, err = suite.DB().ValidateAndCreate(&newRecalculateDates)
		suite.Empty(verrs.Error())
		suite.Nil(err)

		fetchDates, err := models.FetchShipmentRecalculateDates(suite.DB())
		suite.EqualError(err, "Too many active re-calculate date records")
		suite.Nil(fetchDates)

	})

}

func (suite *ModelSuite) TestShipmentRecalculateFetchActiveRecord() {

	suite.T().Run("Fetch active records", func(t *testing.T) {

		id := uuid.Must(uuid.NewV4())

		// Fetch active record
		newRecalculateDatesActive := models.ShipmentRecalculate{
			ID:                    id,
			ShipmentUpdatedAfter:  time.Date(testdatagen.TestYear, time.January, 01, 0, 0, 0, 0, time.UTC),
			ShipmentUpdatedBefore: time.Date(testdatagen.TestYear, time.February, 01, 0, 0, 0, 0, time.UTC),
			Active:                true,
		}
		verrs, err := suite.DB().ValidateAndCreate(&newRecalculateDatesActive)
		suite.Empty(verrs.Error())
		suite.Nil(err)

		id = uuid.Must(uuid.NewV4())

		newRecalculateDatesNotActive := models.ShipmentRecalculate{
			ID:                    id,
			ShipmentUpdatedAfter:  time.Date(1970, time.January, 01, 0, 0, 0, 0, time.UTC),
			ShipmentUpdatedBefore: time.Date(1970, time.February, 01, 0, 0, 0, 0, time.UTC),
			Active:                false,
		}
		verrs, err = suite.DB().ValidateAndCreate(&newRecalculateDatesNotActive)
		suite.Empty(verrs.Error())
		suite.Nil(err)

		fetchDates, err := models.FetchShipmentRecalculateDates(suite.DB())
		suite.Nil(err)
		suite.Empty(verrs.Error())
		suite.NotNil(fetchDates)

	})

}
