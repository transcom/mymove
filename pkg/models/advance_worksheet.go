package models

import (
	"time"

	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// AdvanceWorksheet contains the information needed to populate the PPM Advance Worksheet, which is only generated when a PPM advance is requested
// This pulls information from the service_members, orders, personally_procured_moves, and backup_contacts tables
// This represents the fields in the Advance Worksheet currently represented in our data models. Missing fields include storage type and delivery date, POV shipment authorized, and all Reimburseable Expenses and Enclosed Documentation details
type AdvanceWorksheet struct {
	ID                                   uuid.UUID                           `db:"id"`
	CreatedAt                            time.Time                           `db:"created_at"`
	UpdatedAt                            time.Time                           `db:"updated_at"`
	FirstName                            string                              `db:"first_name"`
	MiddleName                           *string                             `db:"middle_name"`
	LastName                             string                              `db:"last_name"`
	PreferredPhoneNumber                 string                              `db:"preferred_phone_number"`
	Edipi                                *string                             `db:"edipi"`
	Branch                               *internalmessages.Affiliation       `db:"branch"`
	Rank                                 *internalmessages.ServiceMemberRank `db:"rank"`
	Email                                *string                             `db:"email"`
	OrderIssueDate                       time.Time                           `db:"order_issue_date"`
	OrdersType                           internalmessages.OrdersType         `db:"orders_type"`
	OrdersNumber                         *string                             `db:"orders_number"`
	IssuingBranch                        *string                             `db:"issuing_branch"`
	NewDutyAssignment                    string                              `db:"new_duty_assignment"`
	AuthorizedOrigin                     string                              `db:"authorized_origin"`
	AuthorizedDestination                string                              `db:"authorized_destination"`
	ShipmentPickupDate                   time.Time                           `db:"shipment_pickup_date"`
	ShipmentWeight                       *int64                              `db:"shipment_weight"`
	CurrentShipmentStatus                PPMStatus                           `db:"current_shipment_status"`
	StorageTotalDays                     int64                               `db:"storage_total_days"`
	CurrentPaymentRequestClaim           uuid.UUID                           `db:"current_payment_request_claim"`
	CurrentPaymentRequestTransactionType MethodOfReceipt                     `db:"current_payment_request_transaction_type"`
	CurrentPaymentAmount                 unit.Cents                          `db:"current_payment_amount"`
	TrustedAgentName                     string                              `db:"trusted_agent_name"`
	TrustedAgentAuthorizationDate        time.Time                           `db:"trusted_agent_authorization_date"`
	TrustedAgentEmail                    string                              `db:"trusted_agent_email"`
	TrustedAgentPhone                    string                              `db:"trusted_agent_phone"`
}

// String is not required by pop and may be deleted
func (a AdvanceWorksheet) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// AdvanceWorksheets is not required by pop and may be deleted
type AdvanceWorksheets []AdvanceWorksheet

// String is not required by pop and may be deleted
func (a AdvanceWorksheets) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AdvanceWorksheet) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.FirstName, Name: "FirstName"},
		&StringIsNilOrNotBlank{Field: a.MiddleName, Name: "MiddleName"},
		&validators.StringIsPresent{Field: a.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: a.PreferredPhoneNumber, Name: "PreferredPhoneNumber"},
		&StringIsNilOrNotBlank{Field: a.Edipi, Name: "Edipi"},
		&StringIsNilOrNotBlank{Field: a.Email, Name: "Email"},
		&validators.TimeIsPresent{Field: a.OrderIssueDate, Name: "OrderIssueDate"},
		&StringIsNilOrNotBlank{Field: a.OrdersNumber, Name: "OrdersNumber"},
		&StringIsNilOrNotBlank{Field: a.IssuingBranch, Name: "IssuingBranch"},
		&validators.StringIsPresent{Field: a.NewDutyAssignment, Name: "NewDutyAssignment"},
		&validators.StringIsPresent{Field: a.AuthorizedOrigin, Name: "AuthorizedOrigin"},
		&validators.StringIsPresent{Field: a.AuthorizedDestination, Name: "AuthorizedDestination"},
		&validators.TimeIsPresent{Field: a.ShipmentPickupDate, Name: "ShipmentPickupDate"},
		&Int64IsPresent{Field: a.StorageTotalDays, Name: "StorageTotalDays"},
		&validators.StringIsPresent{Field: a.TrustedAgentName, Name: "TrustedAgentName"},
		&validators.TimeIsPresent{Field: a.TrustedAgentAuthorizationDate, Name: "TrustedAgentAuthorizationDate"},
		&validators.StringIsPresent{Field: a.TrustedAgentEmail, Name: "TrustedAgentEmail"},
		&validators.StringIsPresent{Field: a.TrustedAgentPhone, Name: "TrustedAgentPhone"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *AdvanceWorksheet) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *AdvanceWorksheet) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
