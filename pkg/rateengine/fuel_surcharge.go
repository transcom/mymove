package rateengine

// UpdateFuelDieselPrices retrieves data for the months we do not have prices for, calculates them, and adds them to the database
func UpdateFuelDieselPrices() (err error) {
	findMissingDataMonths()
	// Use the start and end dates of the month(s) we want data for to make api call
	// Get the first Mondays (or non-holiday) values (pub_date, price per gallon)
	var pricePerGallon int
	calculateFuelSurchargeBaselineRate(pricePerGallon)
	// Insert values into fuel_eia_diesel_prices
}

func findMissingDataMonths() (months []int, err error) {
	// Determine month(s) we are pulling data for
	// query db table and pull the pub_dates of last 12 months
	// between now and the last 12 months, return what month's data is missing
	return months, err
}

func calculateFuelSurchargeBaselineRate(pricePerGallon int) (baselineRate int, err error) {
	// Calculate fuel surcharge based on price per gallon
	return baselineRate, err
}
