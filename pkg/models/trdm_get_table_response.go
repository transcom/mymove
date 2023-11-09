package models

import (
	"math/big"
	"time"
)

type GetTableResponse struct {
	RowCount   big.Int   `json:"rowCount"`
	StatusCode string    `json:"statusCode"`
	DateTime   time.Time `json:"dateTime"`
	Attachment []byte    `json:"attachment"`
}
