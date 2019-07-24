package publicapi

import (
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForWeightAllotmentModel(allotment models.WeightAllotment) *apimessages.WeightAllotment {
	return &apimessages.WeightAllotment{
		ProGearWeight:                 handlers.FmtInt64(int64(allotment.ProGearWeight)),
		ProGearWeightSpouse:           handlers.FmtInt64(int64(allotment.ProGearWeightSpouse)),
		TotalWeightSelf:               handlers.FmtInt64(int64(allotment.TotalWeightSelf)),
		TotalWeightSelfPlusDependents: handlers.FmtInt64(int64(allotment.TotalWeightSelfPlusDependents)),
	}
}
