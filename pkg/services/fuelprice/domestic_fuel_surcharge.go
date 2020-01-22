package fuelprice

import "github.com/transcom/mymove/pkg/route"

//DomesticFuelSurchargePricer prices fuel surcharge for domestic GHC moves
type DomesticFuelSurchargePricer interface {
	PriceDomesticFuelSurcharge(planner route.Planner, weight int, source string, destination string) (int, error)
}