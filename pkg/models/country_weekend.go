package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type CountryWeekend struct {
	ID                 uuid.UUID `json:"id" db:"id" rw:"r"`
	CountryId          uuid.UUID `json:"country_id" db:"country_id" rw:"r"`
	Country            *Country  `belongs_to:"re_countries" fk_id:"country_id" rw:"r"`
	IsMondayWeekend    bool      `json:"is_monday_weekend" db:"is_monday_weekend" rw:"r"`
	IsTuesdayWeekend   bool      `json:"is_tuesday_weekend" db:"is_tuesday_weekend" rw:"r"`
	IsWednesdayWeekend bool      `json:"is_wednesday_weekend" db:"is_wednesday_weekend" rw:"r"`
	IsThursdayWeekend  bool      `json:"is_thursday_weekend" db:"is_thursday_weekend" rw:"r"`
	IsFridayWeekend    bool      `json:"is_friday_weekend" db:"is_friday_weekend" rw:"r"`
	IsSaturdayWeekend  bool      `json:"is_saturday_weekend" db:"is_saturday_weekend" rw:"r"`
	IsSundayWeekend    bool      `json:"is_sunday_weekend" db:"is_sunday_weekend" rw:"r"`
	CreatedAt          time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

func (c CountryWeekend) TableName() string {
	return "country_weekends"
}
