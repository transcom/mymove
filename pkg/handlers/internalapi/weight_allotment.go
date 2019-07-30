package internalapi

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForWeightAllotmentModel(allotment models.WeightAllotment) *internalmessages.WeightAllotment {
	return &internalmessages.WeightAllotment{
		ProGearWeight:                 handlers.FmtInt64(int64(allotment.ProGearWeight)),
		ProGearWeightSpouse:           handlers.FmtInt64(int64(allotment.ProGearWeightSpouse)),
		TotalWeightSelf:               handlers.FmtInt64(int64(allotment.TotalWeightSelf)),
		TotalWeightSelfPlusDependents: handlers.FmtInt64(int64(allotment.TotalWeightSelfPlusDependents)),
	}
}
