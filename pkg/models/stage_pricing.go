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

type StageDomesticLinehaulPrice struct {
	ServiceAreaNumber string `db:"service_area_number"`
	OriginServiceArea string `db:"origin_service_area"`
	ServicesSchedule  string `db:"services_schedule"`
	Season            string `db:"season"`
	WeightLower       string `db:"weight_lower"`
	WeightUpper       string `db:"weight_upper"`
	MilesLower        string `db:"miles_lower"`
	MilesUpper        string `db:"miles_upper"`
	EscalationNumber  string `db:"escalation_number"`
	Rate              string `db:"rate"`
}

type StageDomesticLinehaulPrices []StageDomesticLinehaulPrice

func (dLh *StageDomesticLinehaulPrice) CSVHeader() []string {
	header := []string{
		"Service Area Number",
		"Origin Service Area",
		"Services Schedule",
		"Season",
		"Weight Lower",
		"Weight Upper",
		"Miles Lower",
		"Miles Upper",
		"Escalation Number",
		"Rate",
	}

	return header
}

func (dLh *StageDomesticLinehaulPrice) ToSlice() []string {
	var values []string

	values = append(values, dLh.ServiceAreaNumber)
	values = append(values, dLh.OriginServiceArea)
	values = append(values, dLh.ServicesSchedule)
	values = append(values, dLh.Season)
	values = append(values, dLh.WeightLower)
	values = append(values, dLh.WeightUpper)
	values = append(values, dLh.MilesLower)
	values = append(values, dLh.MilesUpper)
	values = append(values, dLh.EscalationNumber)
	values = append(values, dLh.Rate)

	return values
}

type StageDomesticServiceAreaPrice struct {
	ServiceAreaNumber                     string `db:"service_area_number"`
	ServiceAreaName                       string `db:"service_area_name"`
	ServicesSchedule                      string `db:"services_schedule"`
	SITPickupDeliverySchedule             string `db:"sit_pickup_delivery_schedule"`
	Season                                string `db:"season"`
	ShorthaulPrice                        string `db:"shorthaul_price"`
	OriginDestinationPrice                string `db:"origin_destination_price"`
	OriginDestinationSITFirstDayWarehouse string `db:"origin_destination_sit_first_day_warehouse"`
	OriginDestinationSITAddlDays          string `db:"origin_destination_sit_addl_days"`
}

type StageDomesticServiceAreaPrices []StageDomesticServiceAreaPrice

func (dSA *StageDomesticServiceAreaPrice) CSVHeader() []string {
	header := []string{
		"Service Area Number",
		"Service Area Name",
		"Services Schedule",
		"SIT Pickup Delivery Schedule",
		"Season",
		"Shorthaul Price",
		"Origin/Destination Price",
		"Origin/Destination SIT First Day & Warehouse",
		"Origin/Destination SIT Add'l Days",
	}

	return header
}

func (dSA *StageDomesticServiceAreaPrice) ToSlice() []string {
	var values []string

	values = append(values, dSA.ServiceAreaNumber)
	values = append(values, dSA.ServiceAreaName)
	values = append(values, dSA.ServicesSchedule)
	values = append(values, dSA.SITPickupDeliverySchedule)
	values = append(values, dSA.Season)
	values = append(values, dSA.ShorthaulPrice)
	values = append(values, dSA.OriginDestinationPrice)
	values = append(values, dSA.OriginDestinationSITFirstDayWarehouse)
	values = append(values, dSA.OriginDestinationSITAddlDays)

	return values
}
