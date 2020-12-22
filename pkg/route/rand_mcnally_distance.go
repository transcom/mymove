package route

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

func randMcNallyZip3Distance(db *pop.Connection, pickupZip string, destinationZip string) (int, error) {
	var distance models.Zip3Distance
	if pickupZip == destinationZip {
		return 0, fmt.Errorf("pickupZip (%s) cannot be the same as destinationZip (%s)", pickupZip, destinationZip)
	} else if pickupZip > destinationZip {
		db.Where("from_zip3 = ? and to_zip3 = ?", destinationZip, pickupZip).First(&distance)
	} else {
		db.Where("from_zip3 = ? and to_zip3 = ?", pickupZip, destinationZip).First(&distance)
	}
	return distance.DistanceMiles, nil
}
