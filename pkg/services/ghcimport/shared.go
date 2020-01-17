package ghcimport

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/transcom/mymove/pkg/models"
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

func getPriceParts(rawPrice string, expectedDecimalPlaces int) (int, int, error) {
	// Get rid of a dollar sign if there is one.
	basePrice := strings.Replace(rawPrice, "$", "", -1)

	// Split the string on the decimal point.
	priceParts := strings.Split(basePrice, ".")
	if len(priceParts) != 2 {
		return 0, 0, fmt.Errorf("expected 2 price parts but found %d for price [%s]", len(priceParts), rawPrice)
	}

	integerPart, err := strconv.Atoi(priceParts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("could not convert integer part of price [%s]", rawPrice)
	}

	if len(priceParts[1]) != expectedDecimalPlaces {
		return 0, 0, fmt.Errorf("expected %d decimal places but found %d for price [%s]", expectedDecimalPlaces,
			len(priceParts[1]), rawPrice)
	}

	fractionalPart, err := strconv.Atoi(priceParts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("could not convert fractional part of price [%s]", rawPrice)
	}

	return integerPart, fractionalPart, nil
}

func priceToMillicents(rawPrice string) (int, error) {
	integerPart, fractionalPart, err := getPriceParts(rawPrice, 3)
	if err != nil {
		return 0, fmt.Errorf("could not parse price [%s]: %w", rawPrice, err)
	}

	millicents := (integerPart * 100000) + (fractionalPart * 100)
	return millicents, nil
}

func priceToCents(rawPrice string) (int, error) {
	integerPart, fractionalPart, err := getPriceParts(rawPrice, 2)
	if err != nil {
		return 0, fmt.Errorf("could not parse price [%s]: %w", rawPrice, err)
	}

	cents := (integerPart * 100) + fractionalPart
	return cents, nil
}

func getMarket(market string) (models.Market, error) {
	if strings.EqualFold(market, "CONUS") {
		return models.MarketConus, nil
	} else if strings.EqualFold(market, "OCONUS") {
		return models.MarketOconus, nil
	}
	return "invalid market", fmt.Errorf("invalid market [%s]", market)
}
