package ghcrateengine

import (
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// NewDomesticFuelSurchargePricer is the public constructor for a DomesticFuelSurchargePricer using Pop
func NewDomesticFuelSurchargePricer(db *pop.Connection, logger Logger, contractCode string) services.DomesticFuelSurchargePricer {
	return &domesticFuelSurchargePricer{
		db:           db,
		logger:       logger,
		contractCode: contractCode,
	}
}

// domesticFuelSurchargePricer is a service object to price domestic fuel surcharge
type domesticFuelSurchargePricer struct {
	db           *pop.Connection
	logger       Logger
	contractCode string
}

//PriceDomesticFuelSurcharge is a placeholder to calculate fuel surcharge, which will be done in the Fuel Surcharge epic.
//Zip3TransitDistance in route.Planner can be used to retrieve the distance needed for the calculation
func (d domesticFuelSurchargePricer) PriceDomesticFuelSurcharge(planner route.Planner, weight unit.Pound, source string, destination string) (unit.Cents, error) {
	return 0, errors.New("Error calculating fuel surcharge")
}
