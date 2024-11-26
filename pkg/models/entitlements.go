package models

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// WeightAllotment represents the weights allotted for a rank
type WeightAllotment struct {
	TotalWeightSelf               int
	TotalWeightSelfPlusDependents int
	ProGearWeight                 int
	ProGearWeightSpouse           int
	UnaccompaniedBaggageAllowance int
}

// the midshipman entitlement is shared with service academy cadet
var midshipman = WeightAllotment{
	TotalWeightSelf:               350,
	TotalWeightSelfPlusDependents: 350,
	ProGearWeight:                 0,
	ProGearWeightSpouse:           0,
}

var aviationCadet = WeightAllotment{
	TotalWeightSelf:               7000,
	TotalWeightSelfPlusDependents: 8000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e1 = WeightAllotment{
	TotalWeightSelf:               5000,
	TotalWeightSelfPlusDependents: 8000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e2 = WeightAllotment{
	TotalWeightSelf:               5000,
	TotalWeightSelfPlusDependents: 8000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e3 = WeightAllotment{
	TotalWeightSelf:               5000,
	TotalWeightSelfPlusDependents: 8000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e4 = WeightAllotment{
	TotalWeightSelf:               7000,
	TotalWeightSelfPlusDependents: 8000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e5 = WeightAllotment{
	TotalWeightSelf:               7000,
	TotalWeightSelfPlusDependents: 9000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e6 = WeightAllotment{
	TotalWeightSelf:               8000,
	TotalWeightSelfPlusDependents: 11000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e7 = WeightAllotment{
	TotalWeightSelf:               11000,
	TotalWeightSelfPlusDependents: 13000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e8 = WeightAllotment{
	TotalWeightSelf:               12000,
	TotalWeightSelfPlusDependents: 14000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e9 = WeightAllotment{
	TotalWeightSelf:               13000,
	TotalWeightSelfPlusDependents: 15000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var e9SpecialSeniorEnlisted = WeightAllotment{
	TotalWeightSelf:               14000,
	TotalWeightSelfPlusDependents: 17000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

// O-1 through O-5 share their entitlements with W-1 through W-5
var o1W1AcademyGraduate = WeightAllotment{
	TotalWeightSelf:               10000,
	TotalWeightSelfPlusDependents: 12000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o2W2 = WeightAllotment{
	TotalWeightSelf:               12500,
	TotalWeightSelfPlusDependents: 13500,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o3W3 = WeightAllotment{
	TotalWeightSelf:               13000,
	TotalWeightSelfPlusDependents: 14500,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o4W4 = WeightAllotment{
	TotalWeightSelf:               14000,
	TotalWeightSelfPlusDependents: 17000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o5W5 = WeightAllotment{
	TotalWeightSelf:               16000,
	TotalWeightSelfPlusDependents: 17500,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o6 = WeightAllotment{
	TotalWeightSelf:               18000,
	TotalWeightSelfPlusDependents: 18000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o7 = WeightAllotment{
	TotalWeightSelf:               18000,
	TotalWeightSelfPlusDependents: 18000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o8 = WeightAllotment{
	TotalWeightSelf:               18000,
	TotalWeightSelfPlusDependents: 18000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o9 = WeightAllotment{
	TotalWeightSelf:               18000,
	TotalWeightSelfPlusDependents: 18000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var o10 = WeightAllotment{
	TotalWeightSelf:               18000,
	TotalWeightSelfPlusDependents: 18000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var civilianEmployee = WeightAllotment{
	TotalWeightSelf:               18000,
	TotalWeightSelfPlusDependents: 18000,
	ProGearWeight:                 2000,
	ProGearWeightSpouse:           500,
}

var entitlements = map[internalmessages.OrderPayGrade]WeightAllotment{
	ServiceMemberGradeACADEMYCADET:            midshipman,
	ServiceMemberGradeAVIATIONCADET:           aviationCadet,
	ServiceMemberGradeE1:                      e1,
	ServiceMemberGradeE2:                      e2,
	ServiceMemberGradeE3:                      e3,
	ServiceMemberGradeE4:                      e4,
	ServiceMemberGradeE5:                      e5,
	ServiceMemberGradeE6:                      e6,
	ServiceMemberGradeE7:                      e7,
	ServiceMemberGradeE8:                      e8,
	ServiceMemberGradeE9:                      e9,
	ServiceMemberGradeE9SPECIALSENIORENLISTED: e9SpecialSeniorEnlisted,
	ServiceMemberGradeMIDSHIPMAN:              midshipman,
	ServiceMemberGradeO1ACADEMYGRADUATE:       o1W1AcademyGraduate,
	ServiceMemberGradeO2:                      o2W2,
	ServiceMemberGradeO3:                      o3W3,
	ServiceMemberGradeO4:                      o4W4,
	ServiceMemberGradeO5:                      o5W5,
	ServiceMemberGradeO6:                      o6,
	ServiceMemberGradeO7:                      o7,
	ServiceMemberGradeO8:                      o8,
	ServiceMemberGradeO9:                      o9,
	ServiceMemberGradeO10:                     o10,
	ServiceMemberGradeW1:                      o1W1AcademyGraduate,
	ServiceMemberGradeW2:                      o2W2,
	ServiceMemberGradeW3:                      o3W3,
	ServiceMemberGradeW4:                      o4W4,
	ServiceMemberGradeW5:                      o5W5,
	ServiceMemberGradeCIVILIANEMPLOYEE:        civilianEmployee,
}

func getEntitlement(grade internalmessages.OrderPayGrade) (WeightAllotment, error) {
	if entitlement, ok := entitlements[grade]; ok {
		return entitlement, nil
	}
	return WeightAllotment{}, fmt.Errorf("no entitlement found for pay grade %s", grade)
}

// AllWeightAllotments returns all the weight allotments for each rank.
func AllWeightAllotments() map[internalmessages.OrderPayGrade]WeightAllotment {
	return entitlements
}

// GetWeightAllotment returns the weight allotments for a given pay grade.
func GetWeightAllotment(grade internalmessages.OrderPayGrade) WeightAllotment {
	entitlement, err := getEntitlement(grade)
	if err != nil {
		return WeightAllotment{}
	}
	return entitlement
}

// GetUBWeightAllowance returns the UB weight allowance for a UB shipment, part of the overall entitlements for an order
func GetUBWeightAllowance(appCtx appcontext.AppContext, originDutyLocationIsOconus *bool, newDutyLocationIsOconus *bool, branch *ServiceMemberAffiliation, grade *internalmessages.OrderPayGrade, orderType *internalmessages.OrdersType, dependentsAuthorized *bool, isAccompaniedTour *bool, dependentsUnderTwelve *int, dependentsTwelveAndOver *int) (int, error) {
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

	// only calculate UB allowance if either origin or new duty locations are OCONUS
	if originDutyLocationIsOconusValue || newDutyLocationIsOconusValue {

		const civilianBaseUBAllowance = 350
		const dependents12AndOverUBAllowance = 350
		const depedentsUnder12UBAllowance = 175
		const maxWholeFamilyCivilianUBAllowance = 2000
		const studentTravelMaxAllowance = 350
		ubAllowance := 0

		if *orderType == internalmessages.OrdersTypeSTUDENTTRAVEL {
			ubAllowance = studentTravelMaxAllowance
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
