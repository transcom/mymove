package models_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicEDIProcessingInstantiation() {
	testCases := map[string]struct {
		ediProcessing models.EDIProcessing
		expectedErrs  map[string][]string
	}{
		"Successful Create": {
			ediProcessing: models.EDIProcessing{
				ID:               uuid.Must(uuid.NewV4()),
				ProcessStartedAt: time.Now(),
				ProcessEndedAt:   time.Now(),
				EDIType:          models.EDIType997,
				NumEDIsProcessed: 6,
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			ediProcessing: models.EDIProcessing{},
			expectedErrs: map[string][]string{
				"process_started_at": {"ProcessStartedAt can not be blank."},
				"process_ended_at":   {"ProcessEndedAt can not be blank."},
				"editype":            {"EDIType is not in the list [810, 824, 858, 997]."},
			},
		},
		"Other Errors": {
			ediProcessing: models.EDIProcessing{
				ID:               uuid.Must(uuid.NewV4()),
				ProcessStartedAt: time.Now(),
				ProcessEndedAt:   time.Now(),
				NumEDIsProcessed: -1,
				EDIType:          "models.EDIType997",
			},
			expectedErrs: map[string][]string{
				"num_edis_processed": {"-1 is not greater than -1."},
				"editype":            {"EDIType is not in the list [810, 824, 858, 997]."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.ediProcessing, test.expectedErrs)
		})
	}

}
