package ghcrateengine

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// NewDomesticServiceAreaPricer is the public constructor for a DomesticRateAreaPricer using Pop
func NewDomesticShorthaulPricer(db *pop.Connection, logger Logger, contractCode string) services.DomesticShorthaulPricer {
	return &domesticShorthaulPricer{
		db:           db,
		logger:       logger,
		contractCode: contractCode,
	}
}

// DomesticShorthaulPricer is a service object to price domestic prices: origin and destination service area, SIT day 1, SIT Addt'l days
// domestic other prices: pack, unpack, and sit p/d costs for a GHC move
type domesticShorthaulPricer struct {
	db           *pop.Connection
	logger       Logger
	contractCode string
}

func (*domesticShorthaulPricer) PriceDomesticShorthaul(moveDate time.Time, distance unit.Miles, weight unit.Pound, serviceArea string) (cost unit.Cents, err error) {
	return cost, err
}
