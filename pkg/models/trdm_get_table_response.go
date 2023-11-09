package models

import (
	"time"
)

type GetTableResponse struct {
	RowCount   int64     `json:"rowCount"`
	StatusCode string    `json:"statusCode"`
	DateTime   time.Time `json:"dateTime"`
	Attachment []byte    `json:"attachment"`
}
