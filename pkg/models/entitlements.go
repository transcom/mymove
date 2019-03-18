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

func makeEntitlements() map[ServiceMemberRank]WeightAllotment {
	// the midshipman entitlement is shared with service academy cadet
	midshipman := WeightAllotment{
		TotalWeightSelf:               350,
		TotalWeightSelfPlusDependents: 3000,
		ProGearWeight:                 0,
		ProGearWeightSpouse:           0,
	}

	aviationCadet := WeightAllotment{
		TotalWeightSelf:               7000,
		TotalWeightSelfPlusDependents: 8000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E1 := WeightAllotment{
		TotalWeightSelf:               5000,
		TotalWeightSelfPlusDependents: 8000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E2 := WeightAllotment{
		TotalWeightSelf:               5000,
		TotalWeightSelfPlusDependents: 8000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E3 := WeightAllotment{
		TotalWeightSelf:               5000,
		TotalWeightSelfPlusDependents: 8000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E4 := WeightAllotment{
		TotalWeightSelf:               7000,
		TotalWeightSelfPlusDependents: 8000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E5 := WeightAllotment{
		TotalWeightSelf:               7000,
		TotalWeightSelfPlusDependents: 9000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E6 := WeightAllotment{
		TotalWeightSelf:               8000,
		TotalWeightSelfPlusDependents: 11000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E7 := WeightAllotment{
		TotalWeightSelf:               11000,
		TotalWeightSelfPlusDependents: 13000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E8 := WeightAllotment{
		TotalWeightSelf:               12000,
		TotalWeightSelfPlusDependents: 14000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	E9 := WeightAllotment{
		TotalWeightSelf:               13000,
		TotalWeightSelfPlusDependents: 15000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	// O-1 through O-5 share their entitlements with W-1 through W-5
	O1W1AcademyGraduate := WeightAllotment{
		TotalWeightSelf:               10000,
		TotalWeightSelfPlusDependents: 12000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O2W2 := WeightAllotment{
		TotalWeightSelf:               12500,
		TotalWeightSelfPlusDependents: 13500,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O3W3 := WeightAllotment{
		TotalWeightSelf:               13000,
		TotalWeightSelfPlusDependents: 14500,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O4W4 := WeightAllotment{
		TotalWeightSelf:               14000,
		TotalWeightSelfPlusDependents: 17000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O5W5 := WeightAllotment{
		TotalWeightSelf:               16000,
		TotalWeightSelfPlusDependents: 17500,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O6 := WeightAllotment{
		TotalWeightSelf:               18000,
		TotalWeightSelfPlusDependents: 18000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O7 := WeightAllotment{
		TotalWeightSelf:               18000,
		TotalWeightSelfPlusDependents: 18000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O8 := WeightAllotment{
		TotalWeightSelf:               18000,
		TotalWeightSelfPlusDependents: 18000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O9 := WeightAllotment{
		TotalWeightSelf:               18000,
		TotalWeightSelfPlusDependents: 18000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	O10 := WeightAllotment{
		TotalWeightSelf:               18000,
		TotalWeightSelfPlusDependents: 18000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	civilianEmployee := WeightAllotment{
		TotalWeightSelf:               18000,
		TotalWeightSelfPlusDependents: 18000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	entitlements := map[ServiceMemberRank]WeightAllotment{
		ServiceMemberRankACADEMYCADET:      midshipman,
		ServiceMemberRankAVIATIONCADET:     aviationCadet,
		ServiceMemberRankE1:                E1,
		ServiceMemberRankE2:                E2,
		ServiceMemberRankE3:                E3,
		ServiceMemberRankE4:                E4,
		ServiceMemberRankE5:                E5,
		ServiceMemberRankE6:                E6,
		ServiceMemberRankE7:                E7,
		ServiceMemberRankE8:                E8,
		ServiceMemberRankE9:                E9,
		ServiceMemberRankMIDSHIPMAN:        midshipman,
		ServiceMemberRankO1ACADEMYGRADUATE: O1W1AcademyGraduate,
		ServiceMemberRankO2:                O2W2,
		ServiceMemberRankO3:                O3W3,
		ServiceMemberRankO4:                O4W4,
		ServiceMemberRankO5:                O5W5,
		ServiceMemberRankO6:                O6,
		ServiceMemberRankO7:                O7,
		ServiceMemberRankO8:                O8,
		ServiceMemberRankO9:                O9,
		ServiceMemberRankO10:               O10,
		ServiceMemberRankW1:                O1W1AcademyGraduate,
		ServiceMemberRankW2:                O2W2,
		ServiceMemberRankW3:                O3W3,
		ServiceMemberRankW4:                O4W4,
		ServiceMemberRankW5:                O5W5,
		ServiceMemberRankCIVILIANEMPLOYEE:  civilianEmployee,
	}
	return entitlements
}

// GetWeightAllotment returns the weight allotments for a given rank.
func GetWeightAllotment(rank ServiceMemberRank) WeightAllotment {
	entitlements := makeEntitlements()
	return entitlements[rank]
}

// GetEntitlement calculates the entitlement for a rank, has dependents and has spouseprogear
func GetEntitlement(rank ServiceMemberRank, hasDependents bool, spouseHasProGear bool) (int, error) {

	entitlements := makeEntitlements()
	spouseProGear := 0
	weight := 0

	selfEntitlement, ok := entitlements[rank]
	if !ok {
		return 0, fmt.Errorf("Rank %s not found in entitlement map", rank)
	}

	if hasDependents {
		if spouseHasProGear {
			spouseProGear = selfEntitlement.ProGearWeightSpouse
		}
		weight = selfEntitlement.TotalWeightSelfPlusDependents
	} else {
		weight = selfEntitlement.TotalWeightSelf
	}
	proGear := selfEntitlement.ProGearWeight

	return weight + proGear + spouseProGear, nil
}
