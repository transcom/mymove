package ghcimport

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func stringToInteger(rawString string) (int, error) {
	// Get rid of any decimal point
	baseString := strings.Split(rawString, ".")[0]

	// Verify that it's an integer
	asInteger, err := strconv.Atoi(baseString)
	if err != nil {
		return 0, err
	}

	return asInteger, nil
}

func cleanServiceAreaNumber(rawServiceArea string) (string, error) {
	serviceAreaInt, err := stringToInteger(rawServiceArea)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d", serviceAreaInt), nil
}

func cleanZip3(rawZip3 string) (string, error) {
	zip3Int, err := stringToInteger(rawZip3)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d", zip3Int), nil
}

func isPeakPeriod(season string) (bool, error) {
	if strings.EqualFold(season, "Peak") {
		return true, nil
	} else if strings.EqualFold(season, "NonPeak") {
		return false, nil
	}

	return false, fmt.Errorf("invalid season [%s]", season)
}

func priceStringToFloat(rawPrice string) (float64, error) {
	basePrice := strings.Replace(rawPrice, "$", "", -1)

	floatPrice, err := strconv.ParseFloat(basePrice, 64)
	if err != nil {
		return 0, err
	}

	return floatPrice, nil
}

func priceToMillicents(rawPrice string) (int, error) {
	floatPrice, err := priceStringToFloat(rawPrice)
	if err != nil {
		return 0, fmt.Errorf("could not parse price [%s]: %w", rawPrice, err)
	}

	millicents := int(math.Round(floatPrice * 100000))
	return millicents, nil
}
