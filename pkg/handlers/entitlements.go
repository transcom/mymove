package handlers

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// WeightAllotment represents the weights allotted for a rank
type weightAllotment struct {
	totalWeightSelf               int
	totalWeightSelfPlusDependents int
	proGearWeight                 int
	proGearWeightSpouse           int
}

func makeEntitlements() map[internalmessages.ServiceMemberRank]weightAllotment {
	midshipman := weightAllotment{
		totalWeightSelf:               350,
		totalWeightSelfPlusDependents: 3000,
		proGearWeight:                 0,
		proGearWeightSpouse:           0,
	}

	aviationCadet := weightAllotment{
		totalWeightSelf:               7000,
		totalWeightSelfPlusDependents: 8000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E1 := weightAllotment{
		totalWeightSelf:               5000,
		totalWeightSelfPlusDependents: 8000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E2 := weightAllotment{
		totalWeightSelf:               5000,
		totalWeightSelfPlusDependents: 8000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E3 := weightAllotment{
		totalWeightSelf:               5000,
		totalWeightSelfPlusDependents: 8000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E4 := weightAllotment{
		totalWeightSelf:               7000,
		totalWeightSelfPlusDependents: 8000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E5 := weightAllotment{
		totalWeightSelf:               7000,
		totalWeightSelfPlusDependents: 9000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E6 := weightAllotment{
		totalWeightSelf:               8000,
		totalWeightSelfPlusDependents: 11000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E7 := weightAllotment{
		totalWeightSelf:               11000,
		totalWeightSelfPlusDependents: 13000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E8 := weightAllotment{
		totalWeightSelf:               12000,
		totalWeightSelfPlusDependents: 14000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	E9 := weightAllotment{
		totalWeightSelf:               13000,
		totalWeightSelfPlusDependents: 15000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O1W1AcademyGraduate := weightAllotment{
		totalWeightSelf:               10000,
		totalWeightSelfPlusDependents: 12000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O2W2 := weightAllotment{
		totalWeightSelf:               12500,
		totalWeightSelfPlusDependents: 13500,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O3W3 := weightAllotment{
		totalWeightSelf:               13000,
		totalWeightSelfPlusDependents: 14500,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O4W4 := weightAllotment{
		totalWeightSelf:               14000,
		totalWeightSelfPlusDependents: 17000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O5W5 := weightAllotment{
		totalWeightSelf:               16000,
		totalWeightSelfPlusDependents: 17500,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O6 := weightAllotment{
		totalWeightSelf:               18000,
		totalWeightSelfPlusDependents: 18000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O7 := weightAllotment{
		totalWeightSelf:               18000,
		totalWeightSelfPlusDependents: 18000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O8 := weightAllotment{
		totalWeightSelf:               18000,
		totalWeightSelfPlusDependents: 18000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O9 := weightAllotment{
		totalWeightSelf:               18000,
		totalWeightSelfPlusDependents: 18000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	O10 := weightAllotment{
		totalWeightSelf:               18000,
		totalWeightSelfPlusDependents: 18000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	civilianEmployee := weightAllotment{
		totalWeightSelf:               18000,
		totalWeightSelfPlusDependents: 18000,
		proGearWeight:                 2000,
		proGearWeightSpouse:           500,
	}

	entitlements := map[internalmessages.ServiceMemberRank]weightAllotment{
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

// getEntitlement calculates the entitlement for a rank, has dependents and has spouseprogear
func getEntitlement(rank internalmessages.ServiceMemberRank, hasDependents bool, spouseHasProGear bool) int {

	entitlements := makeEntitlements()
	spouseProGear := 0
	weight := 0

	if hasDependents {
		if spouseHasProGear {
			spouseProGear = entitlements[rank].proGearWeightSpouse
		}
		weight = entitlements[rank].totalWeightSelfPlusDependents
	} else {
		weight = entitlements[rank].totalWeightSelf
	}
	proGear := entitlements[rank].proGearWeight

	return weight + proGear + spouseProGear
}

// ValidateEntitlementHandler validates a weight estimate based on entitlement
type ValidateEntitlementHandler HandlerContext

// Handle is the handler
func (h ValidateEntitlementHandler) Handle(params entitlementop.ValidateEntitlementParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Fetch move, orders, serviceMember and PPM
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	orders, err := models.FetchOrder(h.db, session, move.OrdersID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	serviceMember, err := models.FetchServiceMember(h.db, session, orders.ServiceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Return 404 if there's no PPM or Rank
	if len(move.PersonallyProcuredMoves) < 1 || serviceMember.Rank == nil {
		return entitlementop.NewValidateEntitlementNotFound()
	}
	// PPMs are in descending order - this is the last one created
	weightEstimate := *move.PersonallyProcuredMoves[0].WeightEstimate

	smEntitlement := getEntitlement(*serviceMember.Rank, orders.HasDependents, orders.SpouseHasProGear)
	if int(weightEstimate) > smEntitlement {
		return responseForConflictErrors(h.logger, fmt.Errorf("your estimated weight of %s lbs is above your weight entitlement of %s lbs. \n You will only be paid for the weight you move up to your weight entitlement", humanize.Comma(weightEstimate), humanize.Comma(int64(smEntitlement))))
	}

	return entitlementop.NewValidateEntitlementOK()
}
