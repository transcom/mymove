package ghcrateengine

import (
"time"

"github.com/gobuffalo/pop"
"github.com/pkg/errors"
"go.uber.org/zap"

"github.com/transcom/mymove/pkg/services"
"github.com/transcom/mymove/pkg/unit"
)

// minDomesticWeight is the minimum weight used in domestic calculations (weights below this are upgraded to the min)
const minDomesticWeight = unit.Pound(500)

// NewDomesticPerWeightServicePricer is the public constructor for a DomesticLinehaulPricer using Pop
func NewDomesticPerWeightservicePricer(db *pop.Connection, logger Logger, contractCode string) services.DomesticPerWeightPricer {
	return &domesticPerWeightServicePricer{
		db:           db,
		logger:       logger,
		contractCode: contractCode,
	}
}

// domesticPerWeightServicePricer is a service object to price domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days
// domestic other prices: pack, unpack, and sit p/d costs for a GHC move
type domesticPerWeightServicePricer struct {
	db           *pop.Connection
	logger       Logger
	contractCode string
}

func (*domesticPerWeightServicePricer) PricePerWeightService (moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string, isDomesticOtherService bool) (unit.Cents, error) {

}