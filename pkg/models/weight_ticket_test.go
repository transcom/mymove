package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestWeightTicketValidation() {
	blankStatusType := models.PPMDocumentStatus("")
	validStatuses := strings.Join(models.AllowedPPMDocumentStatuses, ", ")
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
				Status:                            &blankStatusType,
				Reason:                            models.StringPointer(""),
				AdjustedNetWeight:                 models.PoundPointer(unit.Pound(-1)),
				NetWeightRemarks:                  models.StringPointer(""),
			},
			expectedErrs: map[string][]string{
				"deleted_at":          {"DeletedAt can not be blank."},
				"vehicle_description": {"VehicleDescription can not be blank."},
				"empty_weight":        {"-1 is less than zero."},
				"full_weight":         {"-1 is less than zero."},
				"status":              {fmt.Sprintf("Status is not in the list [%s].", validStatuses)},
				"reason":              {"Reason can not be blank."},
				"adjusted_net_weight": {"-1 is less than zero."},
				"net_weight_remarks":  {"NetWeightRemarks can not be blank."},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.weightTicket, tc.expectedErrs, nil)
		})
	}
}
func (suite *ModelSuite) TestWeightTickets_FilterRejected() {
	suite.Run("returns empty slice when input is empty", func() {
		emptyTickets := models.WeightTickets{}
		filtered := emptyTickets.FilterRejected()
		suite.Equal(0, len(filtered))
	})

	suite.Run("returns all tickets when none are rejected", func() {
		status := models.PPMDocumentStatusApproved
		tickets := models.WeightTickets{
			{Status: &status},
			{Status: &status},
		}
		filtered := tickets.FilterRejected()
		suite.Equal(2, len(filtered))
	})

	suite.Run("filters out rejected tickets only", func() {
		approved := models.PPMDocumentStatusApproved
		rejected := models.PPMDocumentStatusRejected
		tickets := models.WeightTickets{
			{Status: &approved},
			{Status: &rejected},
			{Status: &rejected},
			{Status: nil},
		}
		filtered := tickets.FilterRejected()
		suite.Equal(2, len(filtered))

		for _, ticket := range filtered {
			if ticket.Status != nil {
				suite.NotEqual(models.PPMDocumentStatusRejected, *ticket.Status)
			}
		}
	})

	suite.Run("handles nil status values correctly", func() {
		rejected := models.PPMDocumentStatusRejected
		tickets := models.WeightTickets{
			{Status: nil},
			{Status: &rejected},
			{Status: nil},
		}
		filtered := tickets.FilterRejected()
		suite.Equal(2, len(filtered))

		for _, ticket := range filtered {
			suite.Nil(ticket.Status)
		}
	})
}
