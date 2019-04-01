package models

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/validate/validators"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

type GoldenTicket struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	MoveID    *uuid.UUID `json:"move_id" db:"move_id"`
	Code      string     `json:"code" db:"code"`
	MoveType  string     `json:"move_type" db:"move_type"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (g GoldenTicket) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// GoldenTickets is not required by pop and may be deleted
type GoldenTickets []GoldenTicket

// String is not required by pop and may be deleted
func (g GoldenTickets) String() string {
	jg, _ := json.Marshal(g)
	return string(jg)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (g *GoldenTicket) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: g.Code, Name: "Code"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (g *GoldenTicket) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (g *GoldenTicket) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: g.Code, Name: "Code"},
	), nil
}

func MakeGoldenTicket(db *pop.Connection, moveType SelectedMoveType) (*GoldenTicket, *validate.Errors, error) {
	var err error
	var responseError error
	responseVErrors := validate.NewErrors()
	gt := GoldenTicket{MoveType: string(moveType)}
	gt.Code, err = GenerateGoldenTicketCode()
	if err != nil {
		responseError = errors.Wrap(err, "Error creating golden ticket")
		return &GoldenTicket{}, responseVErrors, err
	}
	verrs, err := db.ValidateAndCreate(&gt)
	if err != nil || verrs.HasAny() {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error creating golden ticket")
		return &GoldenTicket{}, responseVErrors, responseError
	}
	if err != nil {
		responseError = errors.Wrap(err, "Error creating golden ticket")
		return &GoldenTicket{}, responseVErrors, err
	}
	return &gt, responseVErrors, err
}

func GenerateGoldenTicketCode() (string, error) {
	// Not committing to uuid, yet, but just use a placeholder for now
	id, err := uuid.NewV4()
	if err == nil {
		return string(id.String()), nil
	}
	return "", err
}

func ValidateGoldenTicket(db *pop.Connection, code string, move Move) (*GoldenTicket, bool) {
	gt := GoldenTicket{}
	err := db.
		Where("code = ?", code).
		Where("move_id IS NULL").
		//TODO where / when does move type get assigned. Assuming that has already been set
		//TODO when sm enters golden ticket code
		Where("move_type = ?", move.SelectedMoveType).
		First(&gt)
	if err != nil {
		return &gt, false
	}
	return &gt, true
}

func UseGoldenTicket(db *pop.Connection, code string, move Move) (*GoldenTicket, *validate.Errors, error) {
	var err error
	var responseError error
	responseVErrors := validate.NewErrors()

	gt, isValid := ValidateGoldenTicket(db, code, move)
	if !isValid {
		err := errors.New("invalid code")
		responseError = errors.Wrap(err, "Error using golden ticket")
		return &GoldenTicket{}, responseVErrors, responseError
	}

	gt.MoveID = &move.ID
	verrs, err := db.ValidateAndUpdate(gt)
	if err != nil || verrs.HasAny() {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error using golden ticket")
		return &GoldenTicket{}, responseVErrors, responseError
	}
	return gt, responseVErrors, responseError
}

type GoldenTicketCounts map[SelectedMoveType]int

func MakeGoldenTickets(db *pop.Connection, moveTypes GoldenTicketCounts) (GoldenTickets, *validate.Errors, error) {
	verrs := validate.NewErrors()
	var goldenTickets []GoldenTicket
	for k, v := range moveTypes {
		for i := 0; i < v; i++ {
			gt, verrs, err := MakeGoldenTicket(db, k)
			if err != nil || verrs.HasAny() {
				return goldenTickets, verrs, err
			}
			goldenTickets = append(goldenTickets, *gt)
		}
	}
	return goldenTickets, verrs, nil
}
