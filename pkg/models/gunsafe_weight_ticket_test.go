package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestGunSafeWeightTicketValidation() {
	blankStatus := models.PPMDocumentStatus("")
	validStatuses := strings.Join(models.AllowedPPMDocumentStatuses, ", ")

	testCases := map[string]struct {
		gunSafeWeightTicket models.GunSafeWeightTicket
		expectedErrs        map[string][]string
	}{
		"Successful create": {
			gunSafeWeightTicket: models.GunSafeWeightTicket{
				PPMShipmentID: uuid.Must(uuid.NewV4()),
				DocumentID:    uuid.Must(uuid.NewV4()),
			},
			expectedErrs: nil,
		},
		"Missing UUIDs": {
			gunSafeWeightTicket: models.GunSafeWeightTicket{},
			expectedErrs: map[string][]string{
				"ppmshipment_id": {"PPMShipmentID can not be blank."},
				"document_id":    {"DocumentID can not be blank."},
			},
		},
		"Optional fields are invalid": {
			gunSafeWeightTicket: models.GunSafeWeightTicket{
				PPMShipmentID:   uuid.Must(uuid.NewV4()),
				DocumentID:      uuid.Must(uuid.NewV4()),
				Description:     models.StringPointer(""),
				Weight:          models.PoundPointer(unit.Pound(-1)),
				SubmittedWeight: models.PoundPointer(unit.Pound(-1)),
				Status:          &blankStatus,
				Reason:          models.StringPointer(""),
				DeletedAt:       models.TimePointer(time.Time{}),
			},
			expectedErrs: map[string][]string{
				"description":      {"Description can not be blank."},
				"weight":           {"-1 is less than zero."},
				"submitted_weight": {"-1 is less than zero."},
				"status":           {fmt.Sprintf("Status is not in the list [%s].", validStatuses)},
				"reason":           {"Reason can not be blank."},
				"deleted_at":       {"DeletedAt can not be blank."},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.gunSafeWeightTicket, tc.expectedErrs)
		})
	}
}
