package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// OfficeUserStatus represents the status of an office user
type OfficeUserStatus string

const (
	// OfficeUserStatusAPPROVED captures enum value "APPROVED"
	OfficeUserStatusAPPROVED OfficeUserStatus = "APPROVED"
	// OfficeUserStatusREJECTED captures enum value "REJECTED"
	OfficeUserStatusREJECTED OfficeUserStatus = "REJECTED"
	// OfficeUserStatusREQUESTED captures enum value "REQUESTED"
	OfficeUserStatusREQUESTED OfficeUserStatus = "REQUESTED"
)

// OfficeUser is someone who works in one of the TransportationOffices
type OfficeUser struct {
	ID                              uuid.UUID                       `json:"id" db:"id"`
	UserID                          *uuid.UUID                      `json:"user_id" db:"user_id"`
	User                            User                            `belongs_to:"user" fk_id:"user_id"`
	LastName                        string                          `json:"last_name" db:"last_name"`
	FirstName                       string                          `json:"first_name" db:"first_name"`
	MiddleInitials                  *string                         `json:"middle_initials" db:"middle_initials"`
	Email                           string                          `json:"email" db:"email"`
	Telephone                       string                          `json:"telephone" db:"telephone"`
	TransportationOfficeID          uuid.UUID                       `json:"transportation_office_id" db:"transportation_office_id"`
	TransportationOffice            TransportationOffice            `belongs_to:"transportation_office" fk_id:"transportation_office_id"`
	TransportationOfficeAssignments TransportationOfficeAssignments `has_many:"transportation_office_assignments" fk_id:"id" order_by:"created_at asc"`
	CreatedAt                       time.Time                       `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time                       `json:"updated_at" db:"updated_at"`
	Active                          bool                            `json:"active" db:"active"`
	Status                          *OfficeUserStatus               `json:"status" db:"status"`
	EDIPI                           *string                         `json:"edipi" db:"edipi"`
	OtherUniqueID                   *string                         `json:"other_unique_id" db:"other_unique_id"`
	RejectionReason                 *string                         `json:"rejection_reason" db:"rejection_reason"`
	RejectedOn                      *time.Time                      `json:"rejected_on" db:"rejected_on"`
}

type OfficeUserWithWorkload struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Workload  int       `json:"workload" db:"workload"`
}

// TableName overrides the table name used by Pop.
func (o OfficeUser) TableName() string {
	return "office_users"
}

type OfficeUsers []OfficeUser

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (o *OfficeUser) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.EmailIsPresent{Field: o.Email, Name: "Email"},
		&validators.StringIsPresent{Field: o.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: o.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: o.Telephone, Name: "Telephone"},
		&validators.UUIDIsPresent{Field: o.TransportationOfficeID, Name: "TransportationOfficeID"},
	), nil
}

// FetchOfficeUserByEmail looks for an office user with a specific email
func FetchOfficeUserByEmail(tx *pop.Connection, email string) (*OfficeUser, error) {
	var users OfficeUsers
	err := tx.Where("LOWER(email) = $1", strings.ToLower(email)).All(&users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrFetchNotFound
	}
	return &users[0], nil
}

// FetchOfficeUserByID fetches an office user by ID
func FetchOfficeUserByID(tx *pop.Connection, id uuid.UUID) (*OfficeUser, error) {
	var user OfficeUser
	err := tx.Find(&user, id)
	return &user, err
}

// FetchOfficeUserByUserID fetches an office user by UserID
func FetchOfficeUserByUserID(tx *pop.Connection, id uuid.UUID) (*OfficeUser, error) {
	var user OfficeUser
	err := tx.Where("user_id = ?", id).First(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FetchOfficeUserByID fetches an office user by ID
func GetAssignedGBLOCs(o OfficeUser) []string {
	var assignedGblocs []string
	for _, toa := range o.TransportationOfficeAssignments {
		assignedGblocs = append(assignedGblocs, toa.TransportationOffice.Gbloc)
	}
	return assignedGblocs
}
