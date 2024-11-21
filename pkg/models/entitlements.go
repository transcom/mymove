package models

import (
	"fmt"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// WeightAllotment represents the weights allotted for a rank
type WeightAllotment struct {
	TotalWeightSelf               int
	TotalWeightSelfPlusDependents int
	ProGearWeight                 int
	ProGearWeightSpouse           int
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

// allotment by orders type
var studentTravel = WeightAllotment{
	TotalWeightSelf:               350,
	TotalWeightSelfPlusDependents: 350,
	ProGearWeight:                 0,
	ProGearWeightSpouse:           0,
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

var entitlementsByOrdersType = map[internalmessages.OrdersType]WeightAllotment{
	internalmessages.OrdersTypeSTUDENTTRAVEL: studentTravel,
}

func getEntitlement(grade internalmessages.OrderPayGrade) (WeightAllotment, error) {
	if entitlement, ok := entitlements[grade]; ok {
		return entitlement, nil
	}
	return WeightAllotment{}, fmt.Errorf("no entitlement found for pay grade %s", grade)
}

func getEntitlementByOrdersType(ordersType internalmessages.OrdersType) (WeightAllotment, error) {
	if entitlement, ok := entitlementsByOrdersType[ordersType]; ok {
		return entitlement, nil
	}
	return WeightAllotment{}, fmt.Errorf("no entitlement found for orders type %s", ordersType)
}

// AllWeightAllotments returns all the weight allotments for each rank.
func AllWeightAllotments() map[internalmessages.OrderPayGrade]WeightAllotment {
	return entitlements
}

// GetWeightAllotment returns the weight allotments for a given pay grade or an orders type.
func GetWeightAllotment(grade internalmessages.OrderPayGrade, ordersType internalmessages.OrdersType) WeightAllotment {
	var entitlement WeightAllotment
	var err error

	if ordersType == internalmessages.OrdersTypeSTUDENTTRAVEL { // currently only applies to student travel order that limits overall authorized weight
		entitlement, err = getEntitlementByOrdersType(ordersType)
	} else {
		entitlement, err = getEntitlement(grade)
	}
	if err != nil {
		return WeightAllotment{}
	}
	return entitlement
}
