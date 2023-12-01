package models

import "time"

// The first and second date time filters are used to specifically request a range of data. If not provided,
// it will pull data with no filter. This allows us to specifically gather identified date ranges of data.
type GetTableRequest struct {
	PhysicalName                     string     `json:"physicalName"`
	ContentUpdatedSinceDateTime      time.Time  `json:"contentUpdatedSinceDateTime"`
	ReturnContent                    bool       `json:"returnContent"`
	ContentUpdatedOnOrBeforeDateTime *time.Time `json:"contentUpdatedOnOrBeforeDateTime"` // Optional
}
