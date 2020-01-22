package fuelprice

import (
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

//DomesticFuelSurchargePricer prices fuel surcharge for domestic GHC moves
type DomesticFuelSurchargePricer interface {
	PriceDomesticFuelSurcharge(planner route.Planner, weight unit.Pound, source string, destination string) (unit.Cents, error)
}