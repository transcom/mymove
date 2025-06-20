package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// SITDurationUpdateRequestReason type for SIT Duration Update Request Reason
type SITDurationUpdateRequestReason string

const (
	// SITExtensionRequestReasonSeriousIllnessMember is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonSeriousIllnessMember SITDurationUpdateRequestReason = "SERIOUS_ILLNESS_MEMBER"
	// SITExtensionRequestReasonSeriousIllnessDependent is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonSeriousIllnessDependent SITDurationUpdateRequestReason = "SERIOUS_ILLNESS_DEPENDENT"
	// SITExtensionRequestReasonImpendingAssignment is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonImpendingAssignment SITDurationUpdateRequestReason = "IMPENDING_ASSIGNEMENT"
	// SITExtensionRequestReasonDirectedTemporaryDuty is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonDirectedTemporaryDuty SITDurationUpdateRequestReason = "DIRECTED_TEMPORARY_DUTY"
	// SITExtensionRequestReasonNonavailabilityOfCivilianHousing is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonNonavailabilityOfCivilianHousing SITDurationUpdateRequestReason = "NONAVAILABILITY_OF_CIVILIAN_HOUSING"
	// SITExtensionRequestReasonAwaitingCompletionOfResidence is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonAwaitingCompletionOfResidence SITDurationUpdateRequestReason = "AWAITING_COMPLETION_OF_RESIDENCE"
	// SITExtensionRequestReasonOther is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonOther SITDurationUpdateRequestReason = "OTHER"
)

// SITDurationUpdateStatus type for SIT Duration Update request status
type SITDurationUpdateStatus string

const (
	// SITExtensionStatusPending is a SIT extension status
	SITExtensionStatusPending SITDurationUpdateStatus = "PENDING"
	// SITExtensionStatusApproved is a SIT extension status
	SITExtensionStatusApproved SITDurationUpdateStatus = "APPROVED"
	// SITExtensionStatusDenied is a SIT extension status
	SITExtensionStatusDenied SITDurationUpdateStatus = "DENIED"
	// SITExtensionStatusRemoved is a SIT extension status
	SITExtensionStatusRemoved SITDurationUpdateStatus = "REMOVED"
)

// SITDurationUpdates is a slice containing SITDurationUpdate
// Formerly known as "SITExtensions"
type SITDurationUpdates []SITDurationUpdate

// SITDurationUpdate struct representing one SIT duration update
// Formerly known as "SITExtension"
type SITDurationUpdate struct {
	ID                uuid.UUID                      `db:"id"`
	MTOShipment       MTOShipment                    `belongs_to:"mto_shipments" fk_id:"mto_shipment_id"`
	MTOShipmentID     uuid.UUID                      `db:"mto_shipment_id"`
	RequestReason     SITDurationUpdateRequestReason `db:"request_reason"`
	ContractorRemarks *string                        `db:"contractor_remarks"`
	RequestedDays     int                            `db:"requested_days"`
	Status            SITDurationUpdateStatus        `db:"status"`
	ApprovedDays      *int                           `db:"approved_days"`
	DecisionDate      *time.Time                     `db:"decision_date"`
	OfficeRemarks     *string                        `db:"office_remarks"`
	CustomerExpense   *bool                          `db:"customer_expense"`
	CreatedAt         time.Time                      `db:"created_at"`
	UpdatedAt         time.Time                      `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (m SITDurationUpdate) TableName() string {
	return "sit_extensions"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *SITDurationUpdate) Validate(_ *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.UUIDIsPresent{Field: m.MTOShipmentID, Name: "MTOShipmentID"})
	vs = append(vs, &validators.StringInclusion{Field: string(m.RequestReason), Name: "RequestReason", List: []string{
		string(SITExtensionRequestReasonSeriousIllnessMember),
		string(SITExtensionRequestReasonSeriousIllnessDependent),
		string(SITExtensionRequestReasonImpendingAssignment),
		string(SITExtensionRequestReasonDirectedTemporaryDuty),
		string(SITExtensionRequestReasonNonavailabilityOfCivilianHousing),
		string(SITExtensionRequestReasonAwaitingCompletionOfResidence),
		string(SITExtensionRequestReasonOther),
	}})
	vs = append(vs, &validators.StringInclusion{Field: string(m.Status), Name: "Status", List: []string{
		string(SITExtensionStatusPending),
		string(SITExtensionStatusApproved),
		string(SITExtensionStatusDenied),
		string(SITExtensionStatusRemoved),
	}})

	if m.DecisionDate != nil {
		vs = append(vs, &validators.TimeIsPresent{Field: *m.DecisionDate, Name: "DecisionDate"})
	}

	return validate.Validate(vs...), nil
}
