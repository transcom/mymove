package rateengine

import (
	"errors"
	"fmt"
)

func (re *RateEngine) serviceFee(weight int, zip string) (float64, error) {
	fmt.Print(zip)
	// TODO: Fetch service area from zip
	serviceArea := 3
	fmt.Print(serviceArea)
	// TODO: Fetch 135A or 135B Rate
	rate := 3.88
	return float64(weight/100) * rate, nil
}

func (re *RateEngine) fullPack(weight int, zip string) (float64, error) {
	fmt.Print(zip)
	// TODO: Fetch service area from zip
	serviceArea := 3
	fmt.Print(serviceArea)
	// TODO: Fetch service schedule from service area
	serviceSchedule := 1
	fmt.Print(serviceSchedule)
	// TODO: Fetch fullpack rate
	rate := 55.00
	return float64(weight/100) * rate, nil
}

func (re *RateEngine) fullUnpack(weight int, zip string) (float64, error) {
	fmt.Print(zip)
	// TODO: Fetch service area from zip
	serviceArea := 3
	fmt.Print(serviceArea)
	// TODO: Fetch service schedule from service area
	serviceSchedule := 1
	fmt.Print(serviceSchedule)
	// TODO: Fetch full unpack rate
	rate := 5.00
	return float64(weight/100) * rate, nil
}

func (re *RateEngine) nonLinehaulChargeTotal(originZip string, destinationZip string, inverseDiscount float64) (float64, error) {
	weight := 4000
	originServiceFee, err := re.serviceFee(weight, originZip)
	destinationServiceFee, err := re.serviceFee(weight, destinationZip)
	pack, err := re.fullPack(weight, originZip)
	unpack, err := re.fullUnpack(weight, destinationZip)
	if err != nil {
		err = errors.New("Oops nonlinehaulChargeTotal")
	}
	return (originServiceFee + destinationServiceFee + pack + unpack) * inverseDiscount, err
}
