package models

import "github.com/gofrs/uuid"

// GHCDomesticTransitTime Tab Domestic Transit Times
type GHCDomesticTransitTime struct {
	ID                 uuid.UUID `db:"id" csv:"id"`
	MaxDaysTransitTime int       `db:"max_days_transit_time" csv:"max_days_transit_time"`
	WeightLbsLower     int       `db:"weight_lbs_lower" csv:"weight_lbs_lower"`
	WeightLbsUpper     int       `db:"weight_lbs_upper" csv:"weight_lbs_upper"`
	DistanceMilesLower int       `db:"distance_miles_lower" csv:"distance_miles_lower"`
	DistanceMilesUpper int       `db:"distance_miles_upper" csv:"distance_miles_upper"`
}

// TableName overrides the table name used by Pop.
func (g GHCDomesticTransitTime) TableName() string {
	return "ghc_domestic_transit_times"
}
