package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// makeSITExtension creates a single SIT Extension and associated set relationships
func makeSITExtension(db *pop.Connection, assertions Assertions) models.SITExtension {

	var MTOShipmentID uuid.UUID
	approvedDays := 90
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
		ApprovedDays:      &approvedDays,
		DecisionDate:      &decisionDate,
		OfficeRemarks:     &officeRemarks,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	// Overwrite values with those from assertions
	mergeModels(&SITExtension, assertions.SITExtensions)

	mustCreate(db, &SITExtension, assertions.Stub)

	return SITExtension
}

// MakeDefaultSITExtension returns a SITExtension with default values
func MakeDefaultSITExtension(db *pop.Connection) models.SITExtension {
	return makeSITExtension(db, Assertions{})
}

// MakeSITExtensions makes an array of SITExtensions
func MakeSITExtensions(db *pop.Connection, assertions Assertions) models.SITExtensions {
	var sitExtensionList models.SITExtensions
	sitExtensionList = append(sitExtensionList, MakeDefaultSITExtension(db))
	return sitExtensionList
}
