package models

import (
	"encoding/json"
	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"time"
)

// BestValueScore (or BVS) is a combination of quality and value score for
// Transportation Service Providers (TSPs). The higher a TSP's BVS, the higher
// the chance they will be awarded more shipments.
type BestValueScore struct {
	ID                              uuid.UUID `json:"id" db:"id"`
	CreatedAt                       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at" db:"updated_at"`
	TransportationServiceProviderID uuid.UUID `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	Score                           int       `json:"score" db:"score"`
	TrafficDistributionListID       uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
}

// String is not required by pop and may be deleted
func (b BestValueScore) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BestValueScores is not required by pop and may be deleted
type BestValueScores []BestValueScore

// String is not required by pop and may be deleted
func (b BestValueScores) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (b *BestValueScore) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		// Best Value Scores can range from 0 - 100, as defined in DTR403. See page 7
		// of https://www.ustranscom.mil/dtr/part-iv/dtr-part-4-403.pdf
		&validators.IntIsGreaterThan{Field: b.Score, Name: "Score", Compared: -1},
		&validators.IntIsLessThan{Field: b.Score, Name: "Score", Compared: 101},
	), nil
}
