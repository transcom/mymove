package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

// BlackoutDate indicates the range of unavailable times for a TSP and includes its TDL as well.
type BlackoutDate struct {
	ID                              uuid.UUID  `json:"id" db:"id"`
	CreatedAt                       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time  `json:"updated_at" db:"updated_at"`
	TransportationServiceProviderID uuid.UUID  `json:"transportation_service_provider_id" db:"transportation_service_provider_id"`
	StartBlackoutDate               time.Time  `json:"start_blackout_date" db:"start_blackout_date"`
	EndBlackoutDate                 time.Time  `json:"end_blackout_date" db:"end_blackout_date"`
	TrafficDistributionListID       *uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	Market                          *string    `json:"market" db:"market"`
	CodeOfService                   *string    `json:"code_of_service" db:"code_of_service"`
	Channel                         *string    `json:"channel" db:"channel"`
	GBLOC                           *string    `json:"gbloc" db:"gbloc"`
	Zip3                            *int       `json:"zip3" db:"zip3"`
	VolumeMove                      *bool      `json:"volume_move" db:"volume_move"`
}

// FetchTSPBlackoutDates runs a SQL query to find all blackout_date records connected to a TSP ID.
func FetchTSPBlackoutDates(tx *pop.Connection, tspID uuid.UUID, pickupDate time.Time, codeOfService *string, channel *string, gbloc *string, market *string) ([]BlackoutDate, error) {
	fmt.Println("Hello! It is the blackout dates query.")
	// If/then it up to build up appropriate query string
	// A fine place to accommodate the logic that already exists - two attributes that cannot be together
	blackoutDates := []BlackoutDate{}
	moreThanOneClause := false
	sqlString := []string{`SELECT
			*
		FROM
			blackout_dates
		WHERE
			transportation_service_provider_id = $1
		AND
			$2 BETWEEN start_blackout_date AND end_blackout_date`}
	// fmt.Println(sqlString)

	// fmt.Printf("%T\n", sqlString)
	if codeOfService != nil {
		sqlString = append(sqlString, `AND (code_of_service = `, *codeOfService, " ")
		moreThanOneClause = true
	}
	// fmt.Println("This is line 54: ", sqlString)
	// fmt.Printf("This is line 55: %T\n", codeOfService)
	// fmt.Printf("This is line 66: %T\n", *codeOfService)
	if channel != nil {
		if !moreThanOneClause {
			sqlString = append(sqlString, `AND (channel = `, *channel, " ")
			moreThanOneClause = true
		} else {
			sqlString = append(sqlString, ` OR channel = `, *channel, " ")
		}
	}
	// fmt.Println("Line 64! ", sqlString)
	if gbloc != nil {
		if !moreThanOneClause {
			sqlString = append(sqlString, `AND (gbloc = `, *gbloc, " ")
			moreThanOneClause = true
		} else {
			sqlString = append(sqlString, ` OR gbloc = `, *gbloc, " ")
		}
	}
	// fmt.Println("Line 74! ", sqlString)
	// A good place to do that if market, no channel thing - need to change order of things.
	if market != nil {
		if !moreThanOneClause {
			sqlString = append(sqlString, `AND (market = `, *market, " ")
			moreThanOneClause = true
		} else {
			sqlString = append(sqlString, ` OR market = `, *market, " ")
		}
	}

	// fmt.Println(sqlString)
	sqlString = append(sqlString, ";") // Address the paren problem later

	sql := strings.Join(sqlString, "")
	// fmt.Println("Line 89! ", sql)

	err := tx.RawQuery(sql, tspID, pickupDate).All(&blackoutDates)
	// fmt.Println("These are the blackout dates: ", blackoutDates)

	return blackoutDates, err
}

// String is not required by pop and may be deleted
func (b BlackoutDate) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BlackoutDates is not required by pop and may be deleted
type BlackoutDates []BlackoutDate

// String is not required by pop and may be deleted
func (b BlackoutDates) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BlackoutDate) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: b.StartBlackoutDate, Name: "StartBlackoutDate"},
		&validators.TimeIsPresent{Field: b.EndBlackoutDate, Name: "EndBlackoutDate"},
		// &validators.StringIsPresent{Field: b.CodeOfService, Name: "CodeOfService"},
		// TODO: write our own validator that can validate pointers; Pop lacks that
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *BlackoutDate) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *BlackoutDate) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
