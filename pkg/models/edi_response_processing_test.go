package models_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicEDIResponseProcessingInstantiation() {
	testCases := map[string]struct {
		ediResponseProcessing models.EDIResponseProcessing
		expectedErrs          map[string][]string
	}{
		"Successful Create": {
			ediResponseProcessing: models.EDIResponseProcessing{
				ID:               uuid.Must(uuid.NewV4()),
				ProcessStartedAt: time.Now(),
				ProcessEndedAt:   time.Now(),
				MessageType:      models.EDI997,
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			ediResponseProcessing: models.EDIResponseProcessing{},
			expectedErrs: map[string][]string{
				"process_started_at": {"ProcessStartedAt can not be blank."},
				"process_ended_at":   {"ProcessEndedAt can not be blank."},
				"message_type":       {"MessageType is not in the list [997, 824, 810]."},
			},
		},
		"Message Type Invalid": {
			ediResponseProcessing: models.EDIResponseProcessing{
				ID:               uuid.Must(uuid.NewV4()),
				ProcessStartedAt: time.Now(),
				ProcessEndedAt:   time.Now(),
				MessageType:      "models.EDI997",
			},
			expectedErrs: map[string][]string{
				"message_type": {"MessageType is not in the list [997, 824, 810]."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.ediResponseProcessing, test.expectedErrs)
		})
	}

}
