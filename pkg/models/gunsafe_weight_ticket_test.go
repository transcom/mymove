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
			suite.verifyValidationErrors(&tc.gunSafeWeightTicket, tc.expectedErrs, suite.AppContextForTest())
		})
	}
}

func (suite *ModelSuite) TestGunSafeWeightTickets_FilterDeleted() {
	suite.Run("Filters out tickets that have the DeletedAt field present", func() {
		now := time.Now()
		gunSafeTickets := models.GunSafeWeightTickets{
			models.GunSafeWeightTicket{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: &now,
			},
			models.GunSafeWeightTicket{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: &now,
			},
			models.GunSafeWeightTicket{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		filteredTickets := gunSafeTickets.FilterDeleted()
		suite.Equal(1, len(filteredTickets)) // Should only include the one ticket without the DeletedAt field.
	})

	suite.Run("Returns back immediately if an empty object is passed in", func() {
		gunSafeTickets := models.GunSafeWeightTickets{}
		filteredTickets := gunSafeTickets.FilterDeleted()
		suite.Equal(0, len(filteredTickets))
	})
}

func (suite *ModelSuite) TestGunSafeWeightTickets_FilterRejected() {
	suite.Run("returns empty slice when input is empty", func() {
		emptyTickets := models.GunSafeWeightTickets{}
		filtered := emptyTickets.FilterRejected()
		suite.Equal(0, len(filtered))
	})

	suite.Run("returns all tickets when none are rejected", func() {
		status := models.PPMDocumentStatusApproved
		tickets := models.GunSafeWeightTickets{
			{Status: &status},
			{Status: &status},
		}
		filtered := tickets.FilterRejected()
		suite.Equal(2, len(filtered))
	})

	suite.Run("filters out rejected tickets only", func() {
		approved := models.PPMDocumentStatusApproved
		rejected := models.PPMDocumentStatusRejected
		tickets := models.GunSafeWeightTickets{
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
		tickets := models.GunSafeWeightTickets{
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

func (suite *ModelSuite) TestGunSafeWeightTickets_GetApprovedTickets() {
	suite.Run("returns empty slice when input is empty", func() {
		emptyTickets := models.GunSafeWeightTickets{}
		filtered := emptyTickets.GetApprovedTickets()
		suite.Equal(0, len(filtered))
	})

	suite.Run("returns all tickets when all are approved and not deleted", func() {
		status := models.PPMDocumentStatusApproved
		tickets := models.GunSafeWeightTickets{
			{Status: &status, DeletedAt: nil},
			{Status: &status, DeletedAt: nil},
		}
		filtered := tickets.GetApprovedTickets()
		suite.Equal(2, len(filtered))
	})

	suite.Run("filters out non-approved tickets", func() {
		approved := models.PPMDocumentStatusApproved
		rejected := models.PPMDocumentStatusRejected
		tickets := models.GunSafeWeightTickets{
			{Status: &approved},
			{Status: &rejected},
			{Status: nil},
		}
		filtered := tickets.GetApprovedTickets()
		suite.Equal(1, len(filtered))

		for _, ticket := range filtered {
			suite.NotNil(ticket.Status)
			suite.Equal(models.PPMDocumentStatusApproved, *ticket.Status)
		}
	})

	suite.Run("excludes deleted tickets even if approved", func() {
		approved := models.PPMDocumentStatusApproved
		now := time.Now()
		tickets := models.GunSafeWeightTickets{
			{Status: &approved, DeletedAt: &now}, // deleted
			{Status: &approved, DeletedAt: nil},  // not deleted
		}
		filtered := tickets.GetApprovedTickets()
		suite.Equal(1, len(filtered))

		for _, ticket := range filtered {
			suite.Nil(ticket.DeletedAt)
		}
	})

	suite.Run("returns empty slice if no approved tickets remain after filtering", func() {
		rejected := models.PPMDocumentStatusRejected
		now := time.Now()
		tickets := models.GunSafeWeightTickets{
			{Status: &rejected, DeletedAt: nil},
			{Status: nil, DeletedAt: &now},
		}
		filtered := tickets.GetApprovedTickets()
		suite.Equal(0, len(filtered))
	})
}
