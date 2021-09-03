package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeSITExtension creates a single SIT Extension and associated set relationships
func MakeSITExtension(db *pop.Connection, assertions Assertions) models.SITExtension {

	var MTOShipmentID uuid.UUID
	approvedDays := 90
	requestedDays := 100
	decisionDate := time.Now()
	createdAt := time.Now()
	updatedAt := time.Now()
	contractorRemarks := "some remarks here from the contractor"
	officeRemarks := "some remarks here from the office"

	SITExtension := models.SITExtension{
		MTOShipmentID:     MTOShipmentID,
		RequestReason:     models.SITExtensionRequestReasonSeriousIllnessMember,
		ContractorRemarks: &contractorRemarks,
		Status:            models.SITExtensionStatusPending,
		RequestedDays:     requestedDays,
		ApprovedDays:      &approvedDays,
		DecisionDate:      &decisionDate,
		OfficeRemarks:     &officeRemarks,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	// Overwrite values with those from assertions
	mergeModels(&SITExtension, assertions.SITExtension)

	mustCreate(db, &SITExtension, assertions.Stub)

	return SITExtension
}
