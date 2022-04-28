package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"

	"github.com/gofrs/uuid"
)

// Entitlement is an object representing entitlements for orders
type Entitlement struct {
	ID                    uuid.UUID `db:"id"`
	DependentsAuthorized  *bool     `db:"dependents_authorized"`
	TotalDependents       *int      `db:"total_dependents"`
	NonTemporaryStorage   *bool     `db:"non_temporary_storage"`
	PrivatelyOwnedVehicle *bool     `db:"privately_owned_vehicle"`
	//DBAuthorizedWeight is AuthorizedWeight when not null
	DBAuthorizedWeight                           *int             `db:"authorized_weight"`
	weightAllotment                              *WeightAllotment `db:"-"`
	StorageInTransit                             *int             `db:"storage_in_transit"`
	RequiredMedicalEquipmentWeight               int              `db:"required_medical_equipment_weight"`
	OrganizationalClothingAndIndividualEquipment bool             `db:"organizational_clothing_and_individual_equipment"`
	ProGearWeight                                int              `db:"pro_gear_weight"`
	ProGearWeightSpouse                          int              `db:"pro_gear_weight_spouse"`
	CreatedAt                                    time.Time        `db:"created_at"`
	UpdatedAt                                    time.Time        `db:"updated_at"`
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *Entitlement) Validate(*pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: e.ProGearWeight, Compared: -1, Name: "ProGearWeight"},
		&validators.IntIsLessThan{Field: e.ProGearWeight, Compared: 2001, Name: "ProGearWeight"},
		&validators.IntIsGreaterThan{Field: e.ProGearWeightSpouse, Compared: -1, Name: "ProGearWeightSpouse"},
		&validators.IntIsLessThan{Field: e.ProGearWeightSpouse, Compared: 501, Name: "ProGearWeightSpouse"},
	), nil
}

// SetWeightAllotment sets the weight allotment
//TODO probably want to reconsider keeping grade a string rather than enum
//TODO and possibly consider creating ghc specific GetWeightAllotment should the two
//TODO diverge in the future
func (e *Entitlement) SetWeightAllotment(grade string) {
	wa := GetWeightAllotment(ServiceMemberRank(grade))
	e.weightAllotment = &wa
}

// WeightAllotment returns the weight allotment
func (e *Entitlement) WeightAllotment() *WeightAllotment {
	return e.weightAllotment
}

// AuthorizedWeight returns authorized weight. If authorized weight has not been
// stored in DBAuthorizedWeight use either TotalWeightSelf with no dependents or TotalWeightSelfPlusDependents
// with dependents.
func (e *Entitlement) AuthorizedWeight() *int {
	switch {
	case e.DBAuthorizedWeight != nil:
		return e.DBAuthorizedWeight
	case e.WeightAllotment() != nil:
		if e.DependentsAuthorized != nil && *e.DependentsAuthorized {
			return &e.WeightAllotment().TotalWeightSelfPlusDependents
		}
		return &e.WeightAllotment().TotalWeightSelf
	default:
		return nil
	}
}

// WeightAllowance will return the service member's weight allotment based on their rank and if dependents are
// authorized
func (e *Entitlement) WeightAllowance() *int {
	if weightAllotment := e.WeightAllotment(); weightAllotment != nil {
		if e.DependentsAuthorized != nil && *e.DependentsAuthorized {
			return &weightAllotment.TotalWeightSelfPlusDependents
		}
		return &weightAllotment.TotalWeightSelf
	}

	return nil
}
