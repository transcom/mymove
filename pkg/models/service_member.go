package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"

	"github.com/transcom/mymove/pkg/gen/messages"
)

// ServiceMember holds basic identifying data about a service member
type ServiceMember struct {
	ID            uuid.UUID                  `json:"id" db:"id"`
	CreatedAt     time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at" db:"updated_at"`
	FirstName     string                     `json:"first_name" db:"first_name"`
	MiddleInitial string                     `json:"middle_initial" db:"middle_initial"`
	LastName      string                     `json:"last_name" db:"last_name"`
	Rank          messages.ServiceMemberRank `json:"rank" db:"rank"`
	Ssn           string                     `json:"ssn" db:"ssn"`
	Agency        *string                    `json:"agency" db:"agency"`
}

// String is not required by pop and may be deleted
func (s ServiceMember) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// TODO amitchell - add getservicememberid and fetchaddressbyid

// ServiceMembers is not required by pop and may be deleted
type ServiceMembers []ServiceMember

// String is not required by pop and may be deleted
func (s ServiceMembers) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *ServiceMember) Validate(tx *pop.Connection) (*validate.Errors, error) {
	verrs := validate.NewErrors()

	stringFields := map[string]string{
		"FirstName":     a.FirstName,
		"MiddleInitial": a.MiddleInitial,
		"LastName":      a.LastName,
		"Rank":          string(a.Rank),
		"Ssn":           a.Ssn,
	}

	for key, field := range stringFields {
		if field == "" {
			verrs.Add(key, fmt.Sprintf("%s must not be blank!", key))
		}
	}
	return verrs, nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *ServiceMember) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *ServiceMember) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
