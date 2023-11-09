package models

type GetTableRequest struct {
	PhysicalName                string `json:"physicalName"`
	ContentUpdatedSinceDateTime string `json:"contentUpdatedSinceDateTime"`
	ReturnContent               bool   `json:"returnContent"`
}
