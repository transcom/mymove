package models

import "time"

type GetTableRequest struct {
	PhysicalName                string    `json:"physicalName"`
	ContentUpdatedSinceDateTime time.Time `json:"contentUpdatedSinceDateTime"`
	ReturnContent               bool      `json:"returnContent"`
}
