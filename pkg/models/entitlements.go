package models

import (
	"fmt"
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

var entitlements = map[ServiceMemberRank]WeightAllotment{
	ServiceMemberRankACADEMYCADET:            midshipman,
	ServiceMemberRankAVIATIONCADET:           aviationCadet,
	ServiceMemberRankE1:                      e1,
	ServiceMemberRankE2:                      e2,
	ServiceMemberRankE3:                      e3,
	ServiceMemberRankE4:                      e4,
	ServiceMemberRankE5:                      e5,
	ServiceMemberRankE6:                      e6,
	ServiceMemberRankE7:                      e7,
	ServiceMemberRankE8:                      e8,
	ServiceMemberRankE9:                      e9,
	ServiceMemberRankE9SPECIALSENIORENLISTED: e9SpecialSeniorEnlisted,
	ServiceMemberRankMIDSHIPMAN:              midshipman,
	ServiceMemberRankO1ACADEMYGRADUATE:       o1W1AcademyGraduate,
	ServiceMemberRankO2:                      o2W2,
	ServiceMemberRankO3:                      o3W3,
	ServiceMemberRankO4:                      o4W4,
	ServiceMemberRankO5:                      o5W5,
	ServiceMemberRankO6:                      o6,
	ServiceMemberRankO7:                      o7,
	ServiceMemberRankO8:                      o8,
	ServiceMemberRankO9:                      o9,
	ServiceMemberRankO10:                     o10,
	ServiceMemberRankW1:                      o1W1AcademyGraduate,
	ServiceMemberRankW2:                      o2W2,
	ServiceMemberRankW3:                      o3W3,
	ServiceMemberRankW4:                      o4W4,
	ServiceMemberRankW5:                      o5W5,
	ServiceMemberRankCIVILIANEMPLOYEE:        civilianEmployee,
}

func getEntitlement(rank ServiceMemberRank) (WeightAllotment, error) {
	if entitlement, ok := entitlements[rank]; ok {
		return entitlement, nil
	}
	return WeightAllotment{}, fmt.Errorf("no entitlement found for rank %s", rank)
}

// AllWeightAllotments returns all the weight allotments for each rank.
func AllWeightAllotments() map[ServiceMemberRank]WeightAllotment {
	return entitlements
}

// GetWeightAllotment returns the weight allotments for a given rank.
func GetWeightAllotment(rank ServiceMemberRank) WeightAllotment {
	entitlement, err := getEntitlement(rank)
	if err != nil {
		return WeightAllotment{}
	}
	return entitlement
}

// GetEntitlement calculates the entitlement weight based on rank and dependents.
// Only includes either TotalWeightSelf or TotalWeightSelfPlusDependents.
func GetEntitlement(rank ServiceMemberRank, hasDependents bool) (int, error) {
	weight := 0

	selfEntitlement, err := getEntitlement(rank)
	if err != nil {
		return 0, fmt.Errorf("Rank %s not found in entitlement map", rank)
	}

	if hasDependents {
		weight = selfEntitlement.TotalWeightSelfPlusDependents
	} else {
		weight = selfEntitlement.TotalWeightSelf
	}

	return weight, nil
}
