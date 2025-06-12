package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// Entitlement is an object representing entitlements for orders
type Entitlement struct {
	ID                    uuid.UUID `db:"id"`
	DependentsAuthorized  *bool     `db:"dependents_authorized"`
	TotalDependents       *int      `db:"total_dependents" rw:"r"` // DB generated column
	NonTemporaryStorage   *bool     `db:"non_temporary_storage"`
	PrivatelyOwnedVehicle *bool     `db:"privately_owned_vehicle"`
	//DBAuthorizedWeight is AuthorizedWeight when not null
	DBAuthorizedWeight                           *int             `db:"authorized_weight"`
	WeightAllotted                               *WeightAllotment `db:"-"`
	StorageInTransit                             *int             `db:"storage_in_transit"`
	AccompaniedTour                              *bool            `db:"accompanied_tour"`
	DependentsUnderTwelve                        *int             `db:"dependents_under_twelve"`
	DependentsTwelveAndOver                      *int             `db:"dependents_twelve_and_over"`
	UBAllowance                                  *int             `db:"ub_allowance"`
	GunSafe                                      bool             `db:"gun_safe"`
	GunSafeWeight                                int              `db:"gun_safe_weight"`
	RequiredMedicalEquipmentWeight               int              `db:"required_medical_equipment_weight"`
	OrganizationalClothingAndIndividualEquipment bool             `db:"organizational_clothing_and_individual_equipment"`
	ProGearWeight                                int              `db:"pro_gear_weight"`
	ProGearWeightSpouse                          int              `db:"pro_gear_weight_spouse"`
	WeightRestriction                            *int             `db:"weight_restriction"`
	UBWeightRestriction                          *int             `db:"ub_weight_restriction"`
	CreatedAt                                    time.Time        `db:"created_at"`
	UpdatedAt                                    time.Time        `db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (e Entitlement) TableName() string {
	return "entitlements"
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (e *Entitlement) Validate(*pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator

	vs = append(vs,
		&validators.IntIsGreaterThan{Field: e.ProGearWeight, Compared: -1, Name: "ProGearWeight"},
		&validators.IntIsLessThan{Field: e.ProGearWeight, Compared: 2001, Name: "ProGearWeight"},
		&validators.IntIsGreaterThan{Field: e.ProGearWeightSpouse, Compared: -1, Name: "ProGearWeightSpouse"},
		&validators.IntIsLessThan{Field: e.ProGearWeightSpouse, Compared: 501, Name: "ProGearWeightSpouse"},
	)

	if e.DependentsUnderTwelve != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: *e.DependentsUnderTwelve, Compared: -1, Name: "DependentsUnderTwelve"})
	}

	if e.DependentsTwelveAndOver != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: *e.DependentsTwelveAndOver, Compared: -1, Name: "DependentsTwelveAndOver"})
	}

	if e.UBAllowance != nil {
		vs = append(vs, &validators.IntIsGreaterThan{Field: *e.UBAllowance, Compared: -1, Name: "UBAllowance"})
	}

	return validate.Validate(vs...), nil
}

// WeightAllotment returns the weight allotment
func (e *Entitlement) WeightAllotment() *WeightAllotment {
	return e.WeightAllotted
}

// UBWeightAllotment returns the UB weight allotment
func (e *Entitlement) UBWeightAllotment() *int {
	if e.WeightAllotment() != nil {
		if e.WeightAllotment().UnaccompaniedBaggageAllowance >= 0 {
			return &e.WeightAllotment().UnaccompaniedBaggageAllowance
		}
	}
	return nil
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

// WeightAllowance will return the service member's weight allotment based on their grade and if dependents are
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

// UBWeightAllowance returns authorized weight for UB shipments
func (e *Entitlement) UBWeightAllowance() *int {
	switch {
	case e.UBWeightAllotment() != nil:
		return e.UBAllowance
	default:
		return nil
	}
}

// GetUBWeightAllowance returns the UB weight allowance for a UB shipment, part of the overall entitlements for an order
func GetUBWeightAllowance(appCtx appcontext.AppContext, originDutyLocationIsOconus *bool, newDutyLocationIsOconus *bool, branch *ServiceMemberAffiliation, grade *internalmessages.OrderPayGrade, orderType *internalmessages.OrdersType, dependentsAuthorized *bool, isAccompaniedTour *bool, dependentsUnderTwelve *int, dependentsTwelveAndOver *int, civilianTDYUBAllowance *int) (int, error) {
	originDutyLocationIsOconusValue := false
	if originDutyLocationIsOconus != nil {
		originDutyLocationIsOconusValue = *originDutyLocationIsOconus
	}
	newDutyLocationIsOconusValue := false
	if newDutyLocationIsOconus != nil {
		newDutyLocationIsOconusValue = *newDutyLocationIsOconus
	}
	branchOfService := ""
	if branch != nil {
		branchOfService = string(*branch)
	}
	orderPayGrade := ""
	if grade != nil {
		orderPayGrade = string(*grade)
	}
	typeOfOrder := ""
	if orderType != nil {
		typeOfOrder = string(*orderType)
	}
	dependentsAreAuthorized := false
	if dependentsAuthorized != nil {
		dependentsAreAuthorized = *dependentsAuthorized
	}
	isAnAccompaniedTour := false
	if isAccompaniedTour != nil {
		isAnAccompaniedTour = *isAccompaniedTour
	}
	underTwelveDependents := 0
	if dependentsUnderTwelve != nil {
		underTwelveDependents = *dependentsUnderTwelve
	}
	twelveAndOverDependents := 0
	if dependentsTwelveAndOver != nil {
		twelveAndOverDependents = *dependentsTwelveAndOver
	}
	civilianTDYProvidedUBAllowance := 0
	if civilianTDYUBAllowance != nil {
		civilianTDYProvidedUBAllowance = *civilianTDYUBAllowance
	}

	// only calculate UB allowance if either origin or new duty locations are OCONUS
	if originDutyLocationIsOconusValue || newDutyLocationIsOconusValue {

		const civilianBaseUBAllowance = 350
		const dependents12AndOverUBAllowance = 350
		const depedentsUnder12UBAllowance = 175
		const maxWholeFamilyCivilianUBAllowance = 2000
		const studentTravelMaxAllowance = 350
		ubAllowance := 0

		if typeOfOrder == string(internalmessages.OrdersTypeSTUDENTTRAVEL) {
			ubAllowance = studentTravelMaxAllowance
		} else if orderPayGrade == string(internalmessages.OrderPayGradeCIVILIANEMPLOYEE) && typeOfOrder == string(internalmessages.OrdersTypeTEMPORARYDUTY) {
			ubAllowance = civilianTDYProvidedUBAllowance
			return ubAllowance, nil
		} else if orderPayGrade == string(internalmessages.OrderPayGradeCIVILIANEMPLOYEE) && dependentsAreAuthorized && underTwelveDependents == 0 && twelveAndOverDependents == 0 {
			ubAllowance = civilianBaseUBAllowance
		} else if orderPayGrade == string(internalmessages.OrderPayGradeCIVILIANEMPLOYEE) && dependentsAreAuthorized && (underTwelveDependents > 0 || twelveAndOverDependents > 0) {
			ubAllowance = civilianBaseUBAllowance
			// for each dependent 12 and older, add an additional 350 lbs to the civilian's baggage allowance
			ubAllowance += twelveAndOverDependents * dependents12AndOverUBAllowance
			// for each dependent under 12, add an additional 175 lbs to the civilian's baggage allowance
			ubAllowance += underTwelveDependents * depedentsUnder12UBAllowance
			// max allowance of 2,000 lbs for entire family
			if ubAllowance > maxWholeFamilyCivilianUBAllowance {
				ubAllowance = maxWholeFamilyCivilianUBAllowance
			}
		} else {
			if typeOfOrder == string(internalmessages.OrdersTypeLOCALMOVE) {
				// no UB allowance for local moves
				return 0, nil
			} else if typeOfOrder != string(internalmessages.OrdersTypeTEMPORARYDUTY) {
				// all order types other than temporary duty are treated as permanent change of station types for the lookup
				typeOfOrder = string(internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
			}
			// space force members entitled to the same allowance as air force members
			if branchOfService == AffiliationSPACEFORCE.String() {
				branchOfService = AffiliationAIRFORCE.String()
			}
			// e9 special senior enlisted members entitled to the same allowance as e9 members
			if orderPayGrade == string(ServiceMemberGradeE9SPECIALSENIORENLISTED) {
				orderPayGrade = string(ServiceMemberGradeE9)
			}

			var baseUBAllowance UBAllowances
			err := appCtx.DB().Where("branch = ? AND grade = ? AND orders_type = ? AND dependents_authorized = ? AND accompanied_tour = ?", branchOfService, orderPayGrade, typeOfOrder, dependentsAreAuthorized, isAnAccompaniedTour).First(&baseUBAllowance)
			if err != nil {
				if errors.Cause(err).Error() == RecordNotFoundErrorString {
					message := fmt.Sprintf("No UB allowance entry found in ub_allowances table for branch: %s, grade: %s, orders_type: %s, dependents_authorized: %t, accompanied_tour: %t.", branchOfService, orderPayGrade, typeOfOrder, dependentsAreAuthorized, isAnAccompaniedTour)
					appCtx.Logger().Info(message)
					return 0, nil
				}
				return 0, err
			}
			if baseUBAllowance.UBAllowance != nil {
				ubAllowance = *baseUBAllowance.UBAllowance
				return ubAllowance, nil
			} else {
				return 0, nil
			}
		}
		return ubAllowance, nil
	} else {
		appCtx.Logger().Info("No OCONUS duty location found for orders, no UB allowance calculated as part of order entitlement.")
		return 0, nil
	}
}

func GetMaxGunSafeAllowance(appCtx appcontext.AppContext) (int, error) {
	var maxGunSafeAllowance int
	err := appCtx.DB().
		RawQuery(`SELECT parameter_value::int FROM application_parameters WHERE parameter_name = 'maxGunSafeAllowance' LIMIT 1`).
		First(&maxGunSafeAllowance)
	if err != nil {
		return maxGunSafeAllowance, apperror.NewQueryError("ApplicationParameters", err, "error fetching max gun safe allowance")
	}
	return maxGunSafeAllowance, nil
}

// WeightAllotment represents the weights allotted for a rank
type WeightAllotment struct {
	TotalWeightSelf               int
	TotalWeightSelfPlusDependents int
	ProGearWeight                 int
	ProGearWeightSpouse           int
	UnaccompaniedBaggageAllowance int
	GunSafeWeight                 int
}
