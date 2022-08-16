package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestProgearWeightTicketValidation() {
	blankStatus := models.PPMDocumentStatus("")
	validStatuses := strings.Join(models.AllowedPPMDocumentStatuses, ", ")

	testCases := map[string]struct {
		progearWeightTicket models.ProgearWeightTicket
		expectedErrs        map[string][]string
	}{
		"Successful create": {
			progearWeightTicket: models.ProgearWeightTicket{
				PPMShipmentID:               uuid.Must(uuid.NewV4()),
				EmptyDocumentID:             uuid.Must(uuid.NewV4()),
				FullDocumentID:              uuid.Must(uuid.NewV4()),
				ConstructedWeightDocumentID: uuid.Must(uuid.NewV4()),
			},
			expectedErrs: nil,
		},
		"Missing UUIDs": {
			progearWeightTicket: models.ProgearWeightTicket{},
			expectedErrs: map[string][]string{
				"ppmshipment_id":                 {"PPMShipmentID can not be blank."},
				"empty_document_id":              {"EmptyDocumentID can not be blank."},
				"full_document_id":               {"FullDocumentID can not be blank."},
				"constructed_weight_document_id": {"ConstructedWeightDocumentID can not be blank."},
			},
		},
		"Optional fields are invalid": {
			progearWeightTicket: models.ProgearWeightTicket{
				PPMShipmentID:               uuid.Must(uuid.NewV4()),
				EmptyDocumentID:             uuid.Must(uuid.NewV4()),
				FullDocumentID:              uuid.Must(uuid.NewV4()),
				ConstructedWeightDocumentID: uuid.Must(uuid.NewV4()),
				Description:                 models.StringPointer(""),
				EmptyWeight:                 models.PoundPointer(unit.Pound(-1)),
				FullWeight:                  models.PoundPointer(unit.Pound(-1)),
				ConstructedWeight:           models.PoundPointer(unit.Pound(-1)),
				Status:                      &blankStatus,
				Reason:                      models.StringPointer(""),
				DeletedAt:                   models.TimePointer(time.Time{}),
			},
			expectedErrs: map[string][]string{
				"description":        {"Description can not be blank."},
				"empty_weight":       {"-1 is less than zero."},
				"full_weight":        {"-1 is less than zero."},
				"constructed_weight": {"-1 is less than zero."},
				"status":             {fmt.Sprintf("Status is not in the list [%s].", validStatuses)},
				"reason":             {"Reason can not be blank."},
				"deleted_at":         {"DeletedAt can not be blank."},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.progearWeightTicket, tc.expectedErrs)
		})
	}
}
