package models_test

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestSITDurationUpdateCreation() {

	suite.Run("test valid SITDurationUpdate", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		suite.NotNil(shipment)
		suite.NotEqual(uuid.Nil, shipment.ID)
		approvedDays := 90
		decisionDate := time.Now()
		contractorRemarks := "some remarks here from the contractor"
		officeRemarks := "some remarks here from the office"
		validSITExtension := models.SITDurationUpdate{
			MTOShipment:       shipment,
			MTOShipmentID:     shipment.ID,
			RequestReason:     models.SITExtensionRequestReasonSeriousIllnessMember,
			ContractorRemarks: &contractorRemarks,
			RequestedDays:     42,
			Status:            models.SITExtensionStatusPending,
			ApprovedDays:      &approvedDays,
			DecisionDate:      &decisionDate,
			OfficeRemarks:     &officeRemarks,
			CustomerExpense:   models.BoolPointer(false),
		}

		suite.MustSave(&validSITExtension)
		suite.NotNil(validSITExtension.ID)
		suite.NotEqual(uuid.Nil, validSITExtension.ID)
		suite.Equal(shipment.ID, validSITExtension.MTOShipmentID)
		suite.NotEqual(time.Time{}, validSITExtension.CreatedAt)
		suite.NotEqual(time.Time{}, validSITExtension.UpdatedAt)
	})

	suite.Run("test minimal valid SITDurationUpdate", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		suite.NotNil(shipment)
		suite.NotEqual(uuid.Nil, shipment.ID)
		validSITExtension := models.SITDurationUpdate{
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
	suite.Run("test valid SITExtension", func() {
		approvedDays := 1
		decisionDate := time.Now()
		contractorRemarks := "some remarks here from the contractor"
		officeRemarks := "some remarks here from the office"
		validSITExtension := models.SITDurationUpdate{
			MTOShipmentID:     uuid.Must(uuid.NewV4()),
			RequestReason:     models.SITExtensionRequestReasonSeriousIllnessMember,
			ContractorRemarks: &contractorRemarks,
			RequestedDays:     42,
			Status:            models.SITExtensionStatusPending,
			ApprovedDays:      &approvedDays,
			DecisionDate:      &decisionDate,
			OfficeRemarks:     &officeRemarks,
			CustomerExpense:   models.BoolPointer(false),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validSITExtension, expErrors, nil)
	})

	suite.Run("test valid SITDurationUpdate for a SIT duration decrease", func() {
		approvedDays := -2
		decisionDate := time.Now()
		contractorRemarks := "some remarks here from the contractor"
		officeRemarks := "some remarks here from the office"
		validSITExtension := models.SITDurationUpdate{
			MTOShipmentID:     uuid.Must(uuid.NewV4()),
			RequestReason:     models.SITExtensionRequestReasonSeriousIllnessMember,
			ContractorRemarks: &contractorRemarks,
			RequestedDays:     -2,
			Status:            models.SITExtensionStatusPending,
			ApprovedDays:      &approvedDays,
			DecisionDate:      &decisionDate,
			OfficeRemarks:     &officeRemarks,
			CustomerExpense:   models.BoolPointer(false),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validSITExtension, expErrors, nil)
	})

	reasons := []models.SITDurationUpdateRequestReason{
		models.SITExtensionRequestReasonSeriousIllnessMember,
		models.SITExtensionRequestReasonSeriousIllnessDependent,
		models.SITExtensionRequestReasonImpendingAssignment,
		models.SITExtensionRequestReasonDirectedTemporaryDuty,
		models.SITExtensionRequestReasonNonavailabilityOfCivilianHousing,
		models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		models.SITExtensionRequestReasonOther,
	}

	for _, reason := range reasons {
		suite.Run(fmt.Sprintf("test valid SITDurationUpdate Reasons (%s)", reason), func() {
			validSITExtension := models.SITDurationUpdate{
				MTOShipmentID: uuid.Must(uuid.NewV4()),
				RequestReason: reason,
				RequestedDays: 42,
				Status:        models.SITExtensionStatusPending,
			}
			expErrors := map[string][]string{}
			suite.verifyValidationErrors(&validSITExtension, expErrors, nil)
		})
	}

	statuses := []models.SITDurationUpdateStatus{
		models.SITExtensionStatusPending,
		models.SITExtensionStatusApproved,
		models.SITExtensionStatusDenied,
		models.SITExtensionStatusRemoved,
	}

	for _, status := range statuses {
		suite.Run(fmt.Sprintf("test valid SITDurationUpdate Status (%s)", status), func() {
			validSITExtension := models.SITDurationUpdate{
				MTOShipmentID: uuid.Must(uuid.NewV4()),
				RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
				RequestedDays: 42,
				Status:        status,
			}
			expErrors := map[string][]string{}
			suite.verifyValidationErrors(&validSITExtension, expErrors, nil)
		})
	}

	suite.Run("test invalid SITDurationUpdate", func() {
		const badReason models.SITDurationUpdateRequestReason = "bad reason"
		const badStatus models.SITDurationUpdateStatus = "bad status"
		approvedDays := 0
		badDecisionDate := time.Time{}
		validSITExtension := models.SITDurationUpdate{
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
			"status":         {"Status is not in the list [PENDING, APPROVED, DENIED, REMOVED]."},
			"decision_date":  {"DecisionDate can not be blank."},
		}
		suite.verifyValidationErrors(&validSITExtension, expErrors, nil)
	})
}
