package models

type StageDomesticServiceArea struct {
	BasePointCity     string `db:"base_point_city"`
	State             string `db:"state"`
	ServiceAreaNumber string `db:"service_area_number"`
	Zip3s             string `db:"zip3s"`
}

type StageDomesticServiceAreas []StageDomesticServiceArea

func (dsa *StageDomesticServiceArea) CSVHeader() []string {
	header := []string{
		"Base Point City",
		"State",
		"Service Area Number",
		"Zip3's",
	}

	return header
}

func (dsa *StageDomesticServiceArea) ToSlice() []string {
	var values []string

	values = append(values, dsa.BasePointCity)
	values = append(values, dsa.State)
	values = append(values, dsa.ServiceAreaNumber)
	values = append(values, dsa.Zip3s)

	return values
}

type StageInternationalServiceArea struct {
	RateArea   string `db:"rate_area"`
	RateAreaID string `db:"rate_area_id"`
}

type StageInternationalServiceAreas []StageInternationalServiceArea

func (sa *StageInternationalServiceArea) CSVHeader() []string {
	header := []string{
		"International Rate Area",
		"Rate Area Id",
	}

	return header
}

func (sa *StageInternationalServiceArea) ToSlice() []string {
	var values []string

	values = append(values, sa.RateArea)
	values = append(values, sa.RateAreaID)

	return values
}
