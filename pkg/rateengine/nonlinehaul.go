package rateengine

func (re *RateEngine) originServiceFee(weight int, zip string) (float64, error) {
	serviceArea := 3
	rate := 0.45

	return float64(serviceArea) * rate, nil
}
