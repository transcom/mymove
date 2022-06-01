package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
)

// SITExtensionRequestReason type for SIT Extension Request Reason
type SITExtensionRequestReason string

const (
	// SITExtensionRequestReasonSeriousIllnessMember is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonSeriousIllnessMember SITExtensionRequestReason = "SERIOUS_ILLNESS_MEMBER"
	// SITExtensionRequestReasonSeriousIllnessDependent is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonSeriousIllnessDependent SITExtensionRequestReason = "SERIOUS_ILLNESS_DEPENDENT"
	// SITExtensionRequestReasonImpendingAssignment is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonImpendingAssignment SITExtensionRequestReason = "IMPENDING_ASSIGNEMENT"
	// SITExtensionRequestReasonDirectedTemporaryDuty is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonDirectedTemporaryDuty SITExtensionRequestReason = "DIRECTED_TEMPORARY_DUTY"
	// SITExtensionRequestReasonNonavailabilityOfCivilianHousing is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonNonavailabilityOfCivilianHousing SITExtensionRequestReason = "NONAVAILABILITY_OF_CIVILIAN_HOUSING"
	// SITExtensionRequestReasonAwaitingCompletionOfResidence is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonAwaitingCompletionOfResidence SITExtensionRequestReason = "AWAITING_COMPLETION_OF_RESIDENCE"
	// SITExtensionRequestReasonOther is the sit extension request reason type for SIT extensions
	SITExtensionRequestReasonOther SITExtensionRequestReason = "OTHER"
)

// SITExtensionStatus type for SIT Extension status
type SITExtensionStatus string

const (
	// SITExtensionStatusPending is a SIT extension status
	SITExtensionStatusPending SITExtensionStatus = "PENDING"
	// SITExtensionStatusApproved is a SIT extension status
	SITExtensionStatusApproved SITExtensionStatus = "APPROVED"
	// SITExtensionStatusDenied is a SIT extension status
	SITExtensionStatusDenied SITExtensionStatus = "DENIED"
)

// SITExtensions is a slice containing SITExtension
type SITExtensions []SITExtension

// SITExtension struct representing one SIT extension request
type SITExtension struct {
	ID                uuid.UUID                 `db:"id"`
	MTOShipment       MTOShipment               `belongs_to:"mto_shipments" fk_id:"mto_shipment_id"`
	MTOShipmentID     uuid.UUID                 `db:"mto_shipment_id"`
	RequestReason     SITExtensionRequestReason `db:"request_reason"`
	ContractorRemarks *string                   `db:"contractor_remarks"`
	RequestedDays     int                       `db:"requested_days"`
	Status            SITExtensionStatus        `db:"status"`
	ApprovedDays      *int                      `db:"approved_days"`
	DecisionDate      *time.Time                `db:"decision_date"`
	OfficeRemarks     *string                   `db:"office_remarks"`
	CreatedAt         time.Time                 `db:"created_at"`
	UpdatedAt         time.Time                 `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (m SITExtension) TableName() string {
	return "sit_extensions"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *SITExtension) Validate(tx *pop.Connection) (*validate.Errors, error) {
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
	vs = append(vs, &validators.IntIsGreaterThan{Field: m.RequestedDays, Compared: 0, Name: "RequestedDays"})
	vs = append(vs, &validators.StringInclusion{Field: string(m.Status), Name: "Status", List: []string{
		string(SITExtensionStatusPending),
		string(SITExtensionStatusApproved),
		string(SITExtensionStatusDenied),
	}})

	if m.ApprovedDays != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: *m.ApprovedDays, Compared: 0, Name: "ApprovedDays"})
	}

	if m.DecisionDate != nil {
		vs = append(vs, &validators.TimeIsPresent{Field: *m.DecisionDate, Name: "DecisionDate"})
	}

	return validate.Validate(vs...), nil
}
