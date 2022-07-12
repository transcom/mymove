package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestWeightTicketValidation() {
	testCases := map[string]struct {
		weightTicket models.WeightTicket
		expectedErrs map[string][]string
	}{
		"Successful create": {
			weightTicket: models.WeightTicket{
				EmptyWeight:                       models.PoundPointer(unit.Pound(0)),
				EmptyDocumentID:                   uuid.Must(uuid.NewV4()),
				FullWeight:                        models.PoundPointer(unit.Pound(0)),
				FullDocumentID:                    uuid.Must(uuid.NewV4()),
				ProofOfTrailerOwnershipDocumentID: uuid.Must(uuid.NewV4()),
			},
			expectedErrs: nil,
		},
		"Missing UUIDs": {
			weightTicket: models.WeightTicket{},
			expectedErrs: map[string][]string{
				"empty_document_id":                      {"EmptyDocumentID can not be blank."},
				"full_document_id":                       {"FullDocumentID can not be blank."},
				"proof_of_trailer_ownership_document_id": {"ProofOfTrailerOwnershipDocumentID can not be blank."},
			},
		},
		"Optional fields are valid": {
			weightTicket: models.WeightTicket{
				DeletedAt:                         models.TimePointer(time.Time{}),
				VehicleDescription:                models.StringPointer(""),
				EmptyWeight:                       models.PoundPointer(unit.Pound(-1)),
				EmptyDocumentID:                   uuid.Must(uuid.NewV4()),
				FullWeight:                        models.PoundPointer(unit.Pound(-1)),
				FullDocumentID:                    uuid.Must(uuid.NewV4()),
				ProofOfTrailerOwnershipDocumentID: uuid.Must(uuid.NewV4()),
			},
			expectedErrs: map[string][]string{
				"deleted_at":          {"DeletedAt can not be blank."},
				"vehicle_description": {"VehicleDescription can not be blank."},
				"empty_weight":        {"-1 is less than zero."},
				"full_weight":         {"-1 is less than zero."},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.weightTicket, tc.expectedErrs)
		})
	}
}
