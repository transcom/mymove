package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestSITExtensionCreation() {
	shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
	suite.NotNil(shipment)
	suite.NotEqual(uuid.Nil, shipment.ID)

	suite.T().Run("test valid SITExtension", func(t *testing.T) {
		approvedDays := 90
		decisionDate := time.Now()
		contractorRemarks := "some remarks here from the contractor"
		officeRemarks := "some remarks here from the office"
		validSITExtension := models.SITExtension{
			MTOShipment:       shipment,
			MTOShipmentID:     shipment.ID,
			RequestReason:     models.SITExtensionRequestReasonSeriousIllnessMember,
			ContractorRemarks: &contractorRemarks,
			RequestedDays:     42,
			Status:            models.SITExtensionStatusPending,
			ApprovedDays:      &approvedDays,
			DecisionDate:      &decisionDate,
			OfficeRemarks:     &officeRemarks,
		}

		suite.MustSave(&validSITExtension)
		suite.NotNil(validSITExtension.ID)
		suite.NotEqual(uuid.Nil, validSITExtension.ID)
		suite.Equal(shipment.ID, validSITExtension.MTOShipmentID)
		suite.NotEqual(time.Time{}, validSITExtension.CreatedAt)
		suite.NotEqual(time.Time{}, validSITExtension.UpdatedAt)
	})

	suite.T().Run("test minimal valid SITExtension", func(t *testing.T) {
		validSITExtension := models.SITExtension{
			MTOShipment:   shipment,
			MTOShipmentID: shipment.ID,
			RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
			RequestedDays: 42,
			Status:        models.SITExtensionStatusPending,
		}

		suite.MustSave(&validSITExtension)
		suite.NotNil(validSITExtension.ID)
		suite.NotEqual(uuid.Nil, validSITExtension.ID)
		suite.Equal(shipment.ID, validSITExtension.MTOShipmentID)
		suite.NotEqual(time.Time{}, validSITExtension.CreatedAt)
		suite.NotEqual(time.Time{}, validSITExtension.UpdatedAt)
	})
}

func (suite *ModelSuite) TestSITExtensionValidation() {
	suite.T().Run("test valid SITExtension", func(t *testing.T) {
		approvedDays := 1
		decisionDate := time.Now()
		contractorRemarks := "some remarks here from the contractor"
		officeRemarks := "some remarks here from the office"
		validSITExtension := models.SITExtension{
			MTOShipmentID:     uuid.Must(uuid.NewV4()),
			RequestReason:     models.SITExtensionRequestReasonSeriousIllnessMember,
			ContractorRemarks: &contractorRemarks,
			RequestedDays:     42,
			Status:            models.SITExtensionStatusPending,
			ApprovedDays:      &approvedDays,
			DecisionDate:      &decisionDate,
			OfficeRemarks:     &officeRemarks,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validSITExtension, expErrors)
	})

	reasons := []models.SITExtensionRequestReason{
		models.SITExtensionRequestReasonSeriousIllnessMember,
		models.SITExtensionRequestReasonSeriousIllnessDependent,
		models.SITExtensionRequestReasonImpendingAssignment,
		models.SITExtensionRequestReasonDirectedTemporaryDuty,
		models.SITExtensionRequestReasonNonavailabilityOfCivilianHousing,
		models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		models.SITExtensionRequestReasonOther,
	}

	for _, reason := range reasons {
		suite.T().Run(fmt.Sprintf("test valid SITExtension Reasons (%s)", reason), func(t *testing.T) {
			validSITExtension := models.SITExtension{
				MTOShipmentID: uuid.Must(uuid.NewV4()),
				RequestReason: reason,
				RequestedDays: 42,
				Status:        models.SITExtensionStatusPending,
			}
			expErrors := map[string][]string{}
			suite.verifyValidationErrors(&validSITExtension, expErrors)
		})
	}

	statuses := []models.SITExtensionStatus{
		models.SITExtensionStatusPending,
		models.SITExtensionStatusApproved,
		models.SITExtensionStatusDenied,
	}

	for _, status := range statuses {
		suite.T().Run(fmt.Sprintf("test valid SITExtension Status (%s)", status), func(t *testing.T) {
			validSITExtension := models.SITExtension{
				MTOShipmentID: uuid.Must(uuid.NewV4()),
				RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
				RequestedDays: 42,
				Status:        status,
			}
			expErrors := map[string][]string{}
			suite.verifyValidationErrors(&validSITExtension, expErrors)
		})
	}

	suite.T().Run("test invalid sit extension", func(t *testing.T) {
		const badReason models.SITExtensionRequestReason = "bad reason"
		const badStatus models.SITExtensionStatus = "bad status"
		approvedDays := 0
		badDecisionDate := time.Time{}
		validSITExtension := models.SITExtension{
			MTOShipmentID: uuid.Nil,
			RequestReason: badReason,
			RequestedDays: 0,
			ApprovedDays:  &approvedDays,
			Status:        badStatus,
			DecisionDate:  &badDecisionDate,
		}
		expErrors := map[string][]string{
			"mtoshipment_id": {"MTOShipmentID can not be blank."},
			"request_reason": {"RequestReason is not in the list [SERIOUS_ILLNESS_MEMBER, SERIOUS_ILLNESS_DEPENDENT, IMPENDING_ASSIGNEMENT, DIRECTED_TEMPORARY_DUTY, NONAVAILABILITY_OF_CIVILIAN_HOUSING, AWAITING_COMPLETION_OF_RESIDENCE, OTHER]."},
			"requested_days": {"0 is not greater than 0."},
			"status":         {"Status is not in the list [PENDING, APPROVED, DENIED]."},
			"approved_days":  {"0 is not greater than 0."},
			"decision_date":  {"DecisionDate can not be blank."},
		}
		suite.verifyValidationErrors(&validSITExtension, expErrors)
	})
}
