package handlers

import (
	// "fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// WeightAllotment represents the weights alloted for a rank
type WeightAllotment struct {
	Weight map[string]int
}

func makeEntitlements() map[internalmessages.ServiceMemberRank]WeightAllotment {
	midshipman := make(map[string]int)
	midshipman["total_weight_self"] = 350
	midshipman["total_weight_self_plus_dependents"] = 350

	aviationCadet := make(map[string]int)
	aviationCadet["total_weight_self"] = 7000
	aviationCadet["total_weight_self_plus_dependents"] = 8000
	aviationCadet["pro_gear_weight"] = 2000
	aviationCadet["pro_gear_weight_spouse"] = 500

	E1 := make(map[string]int)
	E1["total_weight_self"] = 5000
	E1["total_weight_self_plus_dependents"] = 8000
	E1["pro_gear_weight"] = 2000
	E1["pro_gear_weight_spouse"] = 500

	E2 := make(map[string]int)
	E2["total_weight_self"] = 5000
	E2["total_weight_self_plus_dependents"] = 8000
	E2["pro_gear_weight"] = 2000
	E2["pro_gear_weight_spouse"] = 500

	E3 := make(map[string]int)
	E3["total_weight_self"] = 5000
	E3["total_weight_self_plus_dependents"] = 8000
	E3["pro_gear_weight"] = 2000
	E3["pro_gear_weight_spouse"] = 500

	E4 := make(map[string]int)
	E4["total_weight_self"] = 7000
	E4["total_weight_self_plus_dependents"] = 8000
	E4["pro_gear_weight"] = 2000
	E4["pro_gear_weight_spouse"] = 500

	E5 := make(map[string]int)
	E5["total_weight_self"] = 7000
	E5["total_weight_self_plus_dependents"] = 9000
	E5["pro_gear_weight"] = 2000
	E5["pro_gear_weight_spouse"] = 500

	E6 := make(map[string]int)
	E6["total_weight_self"] = 8000
	E6["total_weight_self_plus_dependents"] = 11000
	E6["pro_gear_weight"] = 2000
	E6["pro_gear_weight_spouse"] = 500

	E7 := make(map[string]int)
	E7["total_weight_self"] = 11000
	E7["total_weight_self_plus_dependents"] = 13000
	E7["pro_gear_weight"] = 2000
	E7["pro_gear_weight_spouse"] = 500

	E8 := make(map[string]int)
	E8["total_weight_self"] = 12000
	E8["total_weight_self_plus_dependents"] = 14000
	E8["pro_gear_weight"] = 2000
	E8["pro_gear_weight_spouse"] = 500

	E9 := make(map[string]int)
	E9["total_weight_self"] = 13000
	E9["total_weight_self_plus_dependents"] = 15000
	E9["pro_gear_weight"] = 2000
	E9["pro_gear_weight_spouse"] = 500

	O1W1AcademyGraduate := make(map[string]int)
	O1W1AcademyGraduate["total_weight_self"] = 10000
	O1W1AcademyGraduate["total_weight_self_plus_dependents"] = 12000
	O1W1AcademyGraduate["pro_gear_weight"] = 2000
	O1W1AcademyGraduate["pro_gear_weight_spouse"] = 500

	O2W2 := make(map[string]int)
	O2W2["total_weight_self"] = 12500
	O2W2["total_weight_self_plus_dependents"] = 13500
	O2W2["pro_gear_weight"] = 2000
	O2W2["pro_gear_weight_spouse"] = 500

	O3W3 := make(map[string]int)
	O3W3["total_weight_self"] = 13000
	O3W3["total_weight_self_plus_dependents"] = 14500
	O3W3["pro_gear_weight"] = 2000
	O3W3["pro_gear_weight_spouse"] = 500

	O4W4 := make(map[string]int)
	O4W4["total_weight_self"] = 14000
	O4W4["total_weight_self_plus_dependents"] = 17000
	O4W4["pro_gear_weight"] = 2000
	O4W4["pro_gear_weight_spouse"] = 500

	O5W5 := make(map[string]int)
	O5W5["total_weight_self"] = 16000
	O5W5["total_weight_self_plus_dependents"] = 17500
	O5W5["pro_gear_weight"] = 2000
	O5W5["pro_gear_weight_spouse"] = 500

	O6 := make(map[string]int)
	O6["total_weight_self"] = 18000
	O6["total_weight_self_plus_dependents"] = 18000
	O6["pro_gear_weight"] = 2000
	O6["pro_gear_weight_spouse"] = 500

	O7 := make(map[string]int)
	O7["total_weight_self"] = 18000
	O7["total_weight_self_plus_dependents"] = 18000
	O7["pro_gear_weight"] = 2000
	O7["pro_gear_weight_spouse"] = 500

	O8 := make(map[string]int)
	O8["total_weight_self"] = 18000
	O8["total_weight_self_plus_dependents"] = 18000
	O8["pro_gear_weight"] = 2000
	O8["pro_gear_weight_spouse"] = 500

	O9 := make(map[string]int)
	O9["total_weight_self"] = 18000
	O9["total_weight_self_plus_dependents"] = 18000
	O9["pro_gear_weight"] = 2000
	O9["pro_gear_weight_spouse"] = 500

	O10 := make(map[string]int)
	O10["total_weight_self"] = 18000
	O10["total_weight_self_plus_dependents"] = 18000
	O10["pro_gear_weight"] = 2000
	O10["pro_gear_weight_spouse"] = 500

	civilianEmployee := make(map[string]int)
	civilianEmployee["total_weight_self"] = 18000
	civilianEmployee["total_weight_self_plus_dependents"] = 18000
	civilianEmployee["pro_gear_weight"] = 2000
	civilianEmployee["pro_gear_weight_spouse"] = 500

	entitlements := make(map[internalmessages.ServiceMemberRank]WeightAllotment)
	entitlements["ACADEMY_CADET_MIDSHIPMAN"] = WeightAllotment{midshipman}
	entitlements["AVIATION_CADET"] = WeightAllotment{aviationCadet}
	entitlements["E_1"] = WeightAllotment{E1}
	entitlements["E_2"] = WeightAllotment{E2}
	entitlements["E_3"] = WeightAllotment{E3}
	entitlements["E_4"] = WeightAllotment{E4}
	entitlements["E_5"] = WeightAllotment{E5}
	entitlements["E_6"] = WeightAllotment{E6}
	entitlements["E_7"] = WeightAllotment{E7}
	entitlements["E_8"] = WeightAllotment{E8}
	entitlements["E_9"] = WeightAllotment{E9}
	entitlements["O_1_W_1_ACADEMY_GRADUATE"] = WeightAllotment{O1W1AcademyGraduate}
	entitlements["O_2_W_2"] = WeightAllotment{O2W2}
	entitlements["O_3_W_3"] = WeightAllotment{O3W3}
	entitlements["O_4_W_4"] = WeightAllotment{O4W4}
	entitlements["O_5_W_5"] = WeightAllotment{O5W5}
	entitlements["O_6"] = WeightAllotment{O6}
	entitlements["O_7"] = WeightAllotment{O7}
	entitlements["O_8"] = WeightAllotment{O8}
	entitlements["O_9"] = WeightAllotment{O9}
	entitlements["O_10"] = WeightAllotment{O10}
	entitlements["CIVILIAN_EMPLOYEE"] = WeightAllotment{civilianEmployee}

	return entitlements
}

func getEntitlement(rank *internalmessages.ServiceMemberRank, hasDependents bool, spouseHasProGear bool) int {

	entitlements := makeEntitlements()
	var totalKey string
	var spouseProGear int

	if hasDependents {
		if spouseHasProGear {
			spouseProGear = entitlements[*rank].Weight["pro_gear_weight_spouse"]
		} else {
			spouseProGear = 0
		}
		totalKey = "total_weight_self_plus_dependents"
	} else {
		totalKey = "total_weight_self"
	}
	weight := entitlements[*rank].Weight[totalKey]
	proGear := entitlements[*rank].Weight["pro_gear_weight"]

	return weight + proGear + spouseProGear
}

// ValidateEntitlementHandler validates a weight estimate based on entitlement
type ValidateEntitlementHandler HandlerContext

// Handle is the handler
func (h ValidateEntitlementHandler) Handle(params entitlementop.ValidateEntitlementParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
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
	// TODO: account for more than 1 PPM, or get ppm by moveid
	weightEstimate := *move.PersonallyProcuredMoves[0].WeightEstimate
	smEntitlement := getEntitlement(serviceMember.Rank, orders.HasDependents, orders.SpouseHasProGear)

	if int(weightEstimate) > smEntitlement {
		// TODO: decide which response to use below - one has a logging message and stack trace
		// return responseForConflictErrors(h.logger, fmt.Errorf("Stored weight estimate is above Entitlement %d", smEntitlement))
		return entitlementop.NewValidateEntitlementConflict()
	}

	return entitlementop.NewValidateEntitlementOK()
}
