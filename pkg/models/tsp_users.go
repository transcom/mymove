package models

import (
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// TspUser is someone who works for a Transportation Service Provider
type TspUser struct {
	ID                              uuid.UUID                     `json:"id" db:"id"`
	UserID                          *uuid.UUID                    `json:"user_id" db:"user_id"`
	User                            User                          `belongs_to:"user"`
	LastName                        string                        `json:"last_name" db:"last_name"`
	FirstName                       string                        `json:"first_name" db:"first_name"`
	MiddleInitials                  *string                       `json:"middle_initials" db:"middle_initials"`
	Email                           string                        `json:"email" db:"email"`
	Telephone                       string                        `json:"telephone" db:"telephone"`
	TransportationServiceProviderID uuid.UUID                     `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	TransportationServiceProvider   TransportationServiceProvider `belongs_to:"transportation_service_provider"`
	CreatedAt                       time.Time                     `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time                     `json:"updated_at" db:"updated_at"`
}

// TspUsers is not required by pop and may be deleted
type TspUsers []TspUser

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *TspUser) Validate(tx *pop.Connection) (*validate.Errors, error) {

	return validate.Validate(
		&validators.StringIsPresent{Field: o.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: o.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: o.Email, Name: "Email"},
		&validators.StringIsPresent{Field: o.Telephone, Name: "Telephone"},
		&validators.UUIDIsPresent{Field: o.TransportationServiceProviderID, Name: "TransportationServiceProviderID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *TspUser) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *TspUser) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchTspUserByID looks for an tsp user with a specific id
func FetchTspUserByID(tx *pop.Connection, id uuid.UUID) (*TspUser, error) {
	var user TspUser
	err := tx.Find(&user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FetchTspUserByEmail looks for an tsp user with a specific email
func FetchTspUserByEmail(tx *pop.Connection, email string) (*TspUser, error) {
	var users TspUsers
	err := tx.Where("email = $1", strings.ToLower(email)).All(&users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrFetchNotFound
	}
	return &users[0], nil
}
