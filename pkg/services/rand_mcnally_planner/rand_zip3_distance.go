package randmcnally

import (
	"github.com/transcom/mymove/pkg/services"
)

type randMcNallyZip3DistanceInfo struct {
	logger Logger
}

// NewRandMcNallyZip3Distance returns a service that can look up distances from Rand McNally database table
func NewRandMcNallyZip3Distance(logger Logger) services.RandMcNallyPlannerMileage {
	return &randMcNallyZip3DistanceInfo{logger: logger}
}

func (r *randMcNallyZip3DistanceInfo) RandMcNallyZip3Distance(pickupZip string, destinationZip string) (int, error) {
	// TODO: Implement this
	return -1, nil
}
