package models

import "time"

type LastTableUpdateResponse struct {
	StatusCode string    `json:"physicalName"`
	DateTime   time.Time `json:"dateTime"`
	LastUpdate time.Time `json:"lastUpdate"`
}
