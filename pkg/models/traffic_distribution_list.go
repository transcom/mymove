package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// TrafficDistributionList items are essentially different markets, based on
// source and destination, in which Transportation Service Providers (TSPs)
// bid on shipments.
type TrafficDistributionList struct {
	ID                uuid.UUID `json:"id" db:"id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	SourceRateArea    string    `json:"source_rate_area" db:"source_rate_area"`
	DestinationRegion string    `json:"destination_region" db:"destination_region"`
	CodeOfService     string    `json:"code_of_service" db:"code_of_service"`
}

// String is not required by pop and may be deleted
func (t TrafficDistributionList) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// TrafficDistributionLists is not required by pop and may be deleted
type TrafficDistributionLists []TrafficDistributionList

// String is not required by pop and may be deleted
func (t TrafficDistributionLists) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TrafficDistributionList) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.SourceRateArea, Name: "SourceRateArea"},
		&validators.StringIsPresent{Field: t.DestinationRegion, Name: "DestinationRegion"},
		&validators.StringIsPresent{Field: t.CodeOfService, Name: "CodeOfService"},
	), nil
}

// FetchTDLsAwaitingBandAssignment returns TDLs with at least one TransportationServiceProviderPerformance containing a null QualityBand.
func FetchTDLsAwaitingBandAssignment(db *pop.Connection) (TrafficDistributionLists, error) {
	tdls := TrafficDistributionLists{}

	sql := `SELECT
				tdl.*
			FROM
				traffic_distribution_lists AS tdl
			LEFT JOIN
				transportation_service_provider_performances AS tspp ON
					tspp.traffic_distribution_list_id = tdl.id
			WHERE
				tspp.quality_band IS NULL
			GROUP BY
				tdl.id
			ORDER BY
				tdl.id
			`

	err := db.RawQuery(sql).All(&tdls)

	return tdls, err
}
