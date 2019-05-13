package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TrafficDistributionList items are essentially different markets, based on
// source and destination, in which Transportation Service Providers (TSPs)
// bid on shipments.
type TrafficDistributionList struct {
	ID                uuid.UUID `json:"id" db:"id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	SourceRateArea    string    `json:"source_rate_area" db:"source_rate_area"`
	DestinationRegion string    `json:"destination_region" db:"destination_region"`
	CodeOfService     string    `json:"code_of_service" db:"code_of_service"`
}

// TrafficDistributionLists is not required by pop and may be deleted
type TrafficDistributionLists []TrafficDistributionList

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TrafficDistributionList) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.SourceRateArea, Name: "SourceRateArea"},
		&validators.RegexMatch{Field: t.SourceRateArea, Name: "SourceRateArea", Expr: "^US.*$"},
		&validators.StringIsPresent{Field: t.DestinationRegion, Name: "DestinationRegion"},
		&validators.RegexMatch{Field: t.DestinationRegion, Name: "DestinationRegion", Expr: "^[0-9]+"},
		&validators.StringIsPresent{Field: t.CodeOfService, Name: "CodeOfService"},
	), nil
}

// MarshalLogObject is required to be able to zap.Object log TDLs
func (t TrafficDistributionList) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("src", t.SourceRateArea)
	encoder.AddString("dest", t.DestinationRegion)
	encoder.AddString("cos", t.CodeOfService)
	return nil
}

// FetchTDL attempts to return a TDL based on SourceRateArea, Region, and CodeOfService (COS).
func FetchTDL(db *pop.Connection, rateArea string, region string, codeOfService string) (TrafficDistributionList, error) {
	var trafficDistributionList TrafficDistributionList
	err := db.Where("source_rate_area = ?", rateArea).
		Where("destination_region = ?", region).
		Where("code_of_service = ?", codeOfService).
		First(&trafficDistributionList)

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return TrafficDistributionList{}, ErrFetchNotFound
		}
		return TrafficDistributionList{}, err
	}

	return trafficDistributionList, nil
}

// FetchOrCreateTDL attempts to return a TDL based on SourceRateArea, Region, and CodeOfService (COS)
// and creates one to return if it doesn't already exist.
func FetchOrCreateTDL(db *pop.Connection, rateArea string, region string, codeOfService string) (TrafficDistributionList, error) {
	// Fetch TDL and return it immediately if found.
	trafficDistributionList, fetchTDLErr := FetchTDL(db, rateArea, region, codeOfService)
	if fetchTDLErr == nil {
		return trafficDistributionList, fetchTDLErr
	}

	// If we didn't find the TDL, create it.
	if fetchTDLErr == ErrFetchNotFound {
		trafficDistributionList := TrafficDistributionList{
			SourceRateArea:    rateArea,
			DestinationRegion: region,
			CodeOfService:     codeOfService,
		}
		verrs, err := db.ValidateAndSave(&trafficDistributionList)
		if err != nil {
			zap.L().Error("DB insertion error", zap.Error(err))
			return TrafficDistributionList{}, err
		} else if verrs.HasAny() {
			zap.L().Error("Validation errors", zap.Error(verrs))
			return TrafficDistributionList{}, errors.New("Validation error on TDL")
		}
		return trafficDistributionList, err
	}

	// If we get here, an unexpected error occurred.
	return TrafficDistributionList{}, fetchTDLErr
}
