package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type CountryHoliday struct {
	ID              uuid.UUID `json:"id" db:"id" rw:"r"`
	CountryId       uuid.UUID `json:"country_id" db:"country_id" rw:"r"`
	Country         *Country  `belongs_to:"re_countries" fk_id:"country_id" rw:"r"`
	HolidayName     string    `json:"holiday_name" db:"holiday_name" rw:"r"`
	ObservationDate time.Time `json:"observation_date" db:"observation_date" rw:"r"`
	CreatedAt       time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

func (c CountryHoliday) TableName() string {
	return "country_holidays"
}

// CountryHolidays is a slice containing CountryHoliday
type CountryHolidays []CountryHoliday
