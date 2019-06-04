package rateengine

import (
	"fmt"

	"github.com/transcom/mymove/pkg/unit"
)

type pricerTestCase struct {
	pricer   pricer
	rate     unit.Cents
	quantity unit.BaseQuantity
	discount *unit.DiscountRate
	expected unit.Cents
}

func discountPtr(d float64) *unit.DiscountRate {
	rate := unit.DiscountRate(d)
	return &rate
}

var pricersTestCases = []pricerTestCase{
	{newBasicQuantityPricer(), unit.Cents(100), unit.BaseQuantityFromInt(1), nil, unit.Cents(100)},
	{newBasicQuantityPricer(), unit.Cents(100), unit.BaseQuantityFromInt(1), discountPtr(0.5), unit.Cents(50)},

	{newMinimumQuantityPricer(10), unit.Cents(100), unit.BaseQuantityFromInt(100), nil, unit.Cents(10000)},
	{newMinimumQuantityPricer(10), unit.Cents(100), unit.BaseQuantityFromInt(5), nil, unit.Cents(1000)},
	{newMinimumQuantityPricer(10), unit.Cents(100), unit.BaseQuantityFromInt(100), discountPtr(0.5), unit.Cents(5000)},

	{newMinimumQuantityHundredweightPricer(100), unit.Cents(100), unit.BaseQuantityFromInt(1000), nil, unit.Cents(1000)},
	{newMinimumQuantityHundredweightPricer(100), unit.Cents(100), unit.BaseQuantityFromInt(50), nil, unit.Cents(100)},
	{newMinimumQuantityHundredweightPricer(100), unit.Cents(100), unit.BaseQuantityFromInt(1000), discountPtr(0.5), unit.Cents(500)},

	{newFlatRatePricer(), unit.Cents(100), unit.BaseQuantityFromInt(9999), nil, unit.Cents(100)},
	{newFlatRatePricer(), unit.Cents(100), unit.BaseQuantityFromInt(9999), discountPtr(0.5), unit.Cents(50)},
}

func (suite *RateEngineSuite) TestPricersTestCases() {
	for i, testCase := range pricersTestCases {
		result := testCase.pricer.price(testCase.rate, testCase.quantity, testCase.discount)
		if !suite.Equal(result, testCase.expected) {
			fmt.Printf("Failure on test case %d (0 indexed)\n", i)
		}
	}
}
