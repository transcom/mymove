package models

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// WeightAllotment represents the weights allotted for a rank
type WeightAllotment struct {
	TotalWeightSelf               int
	TotalWeightSelfPlusDependents int
	ProGearWeight                 int
	ProGearWeightSpouse           int
}

func makeEntitlements() map[internalmessages.ServiceMemberRank]WeightAllotment {
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

	entitlements := map[internalmessages.ServiceMemberRank]WeightAllotment{
		internalmessages.ServiceMemberRankACADEMYCADETMIDSHIPMAN: midshipman,
		internalmessages.ServiceMemberRankAVIATIONCADET:          aviationCadet,
		internalmessages.ServiceMemberRankE1:                     E1,
		internalmessages.ServiceMemberRankE2:                     E2,
		internalmessages.ServiceMemberRankE3:                     E3,
		internalmessages.ServiceMemberRankE4:                     E4,
		internalmessages.ServiceMemberRankE5:                     E5,
		internalmessages.ServiceMemberRankE6:                     E6,
		internalmessages.ServiceMemberRankE7:                     E7,
		internalmessages.ServiceMemberRankE8:                     E8,
		internalmessages.ServiceMemberRankE9:                     E9,
		internalmessages.ServiceMemberRankO1W1ACADEMYGRADUATE:    O1W1AcademyGraduate,
		internalmessages.ServiceMemberRankO2W2:                   O2W2,
		internalmessages.ServiceMemberRankO3W3:                   O3W3,
		internalmessages.ServiceMemberRankO4W4:                   O4W4,
		internalmessages.ServiceMemberRankO5W5:                   O5W5,
		internalmessages.ServiceMemberRankO6:                     O6,
		internalmessages.ServiceMemberRankO7:                     O7,
		internalmessages.ServiceMemberRankO8:                     O8,
		internalmessages.ServiceMemberRankO9:                     O9,
		internalmessages.ServiceMemberRankO10:                    O10,
		internalmessages.ServiceMemberRankCIVILIANEMPLOYEE:       civilianEmployee,
	}
	return entitlements
}

// GetWeightAllotment returns the weight allotments for a given rank.
func GetWeightAllotment(rank ServiceMemberRank) WeightAllotment {
	entitlements := makeEntitlements()
	return entitlements[internalmessages.ServiceMemberRank(rank)]
}

// GetEntitlement calculates the entitlement for a rank, has dependents and has spouseprogear
func GetEntitlement(rank internalmessages.ServiceMemberRank, hasDependents bool, spouseHasProGear bool) int {

	entitlements := makeEntitlements()
	spouseProGear := 0
	weight := 0

	if hasDependents {
		if spouseHasProGear {
			spouseProGear = entitlements[rank].ProGearWeightSpouse
		}
		weight = entitlements[rank].TotalWeightSelfPlusDependents
	} else {
		weight = entitlements[rank].TotalWeightSelf
	}
	proGear := entitlements[rank].ProGearWeight

	return weight + proGear + spouseProGear
}
