package models

import (
	"time"

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
	DBAuthorizedWeight *int             `db:"authorized_weight"`
	weightAllotment    *WeightAllotment `db:"-"`
	StorageInTransit   *int             `db:"storage_in_transit"`
	CreatedAt          time.Time        `db:"created_at"`
	UpdatedAt          time.Time        `db:"updated_at"`
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
		if e.DependentsAuthorized != nil && *e.DependentsAuthorized == true {
			return &e.WeightAllotment().TotalWeightSelfPlusDependents
		}
		return &e.WeightAllotment().TotalWeightSelf
	default:
		return nil
	}
}
