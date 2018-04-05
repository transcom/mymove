package rateengine

import (
	"errors"
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

func (re *RateEngine) serviceFee(cwt int, zip3 string) (float64, error) {
	serviceArea, err := models.ServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0.0, err
	}
	// TODO: Fetch 135A or 135B (serviceArea)
	rate, err := models.Rate135A(re.db, serviceArea)
	if err != nil {
		return 0.0, err
	}
	return float64(cwt) * rate, nil
}

func (re *RateEngine) fullPack(cwt int, zip3 string) (float64, error) {
	serviceArea, err := models.ServiceAreaForZip3(re.db, zip3)
	if err != nil {
		return 0.0, err
	}
	fmt.Println(serviceArea)

	// TODO: Fetch service schedule from service area
	serviceSchedule := 1
	fmt.Print(serviceSchedule)
	// TODO: Fetch fullpack rate
	rate := 55.00
	return float64(cwt) * rate, nil
}

func (re *RateEngine) fullUnpack(cwt int, zip string) (float64, error) {
	fmt.Print(zip)
	// TODO: Fetch service area from zip
	serviceArea := 3
	fmt.Print(serviceArea)
	// TODO: Fetch service schedule from service area
	serviceSchedule := 1
	fmt.Print(serviceSchedule)
	// TODO: Fetch full unpack rate
	rate := 5.00
	return float64(cwt) * rate, nil
}

func (re *RateEngine) nonLinehaulChargeTotal(originZip string, destinationZip string, inverseDiscount float64) (float64, error) {
	weight := 4000
	cwt := re.determineCWT(weight)
	originServiceFee, err := re.serviceFee(cwt, originZip)
	destinationServiceFee, err := re.serviceFee(cwt, destinationZip)
	pack, err := re.fullPack(cwt, originZip)
	unpack, err := re.fullUnpack(cwt, destinationZip)
	if err != nil {
		err = errors.New("Oops nonlinehaulChargeTotal")
	}
	return (originServiceFee + destinationServiceFee + pack + unpack) * inverseDiscount, err
}
