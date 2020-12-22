package randmcnally

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type randMcNallyZip3DistanceInfo struct {
	db     *pop.Connection
	logger Logger
}

// NewRandMcNallyZip3Distance returns a service that can look up distances from Rand McNally database table
func NewRandMcNallyZip3Distance(db *pop.Connection, logger Logger) services.RandMcNallyPlannerMileage {
	return &randMcNallyZip3DistanceInfo{db: db, logger: logger}
}

func (r *randMcNallyZip3DistanceInfo) RandMcNallyZip3Distance(pickupZip string, destinationZip string) (int, error) {
	var distance models.Zip3Distance
	if pickupZip == destinationZip {
		return 0, services.NewBadDataError("pickupZip cannot be the same as destinationZip")
	} else if pickupZip > destinationZip {
		r.db.Where("from_zip3 = ? and to_zip3 = ?", destinationZip, pickupZip).First(&distance)
	} else {
		r.db.Where("from_zip3 = ? and to_zip3 = ?", pickupZip, destinationZip).First(&distance)
	}
	return distance.DistanceMiles, nil
}
