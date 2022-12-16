package ghcdieselfuelprice

import "fmt"

type GHCAPIValidationError struct {
	message string
}

// NewGHCAPIValidationError returns a new GHCAPIValidationError
func NewGHCAPIValidationError(message string) GHCAPIValidationError {
	return GHCAPIValidationError{
		message: message,
	}
}

// Error is the string representation of the InternalServerError
func (e GHCAPIValidationError) Error() string {
	return e.message
}

func (e EIAData) validateEIAData() error {
	err := e.validateFuelDataExists()
	if err != nil {
		return err
	}

	err = e.validateDateFormat()
	if err != nil {
		return err
	}

	err = e.validateFrequency()
	if err != nil {
		return err
	}

	err = e.validateDuoArea()
	if err != nil {
		return err
	}

	err = e.validateAreaName()
	if err != nil {
		return err
	}

	err = e.validateProduct()
	if err != nil {
		return err
	}

	err = e.validateProcess()
	if err != nil {
		return err
	}

	err = e.validateSeries()
	if err != nil {
		return err
	}

	err = e.validateUnits()
	if err != nil {
		return err
	}

	return nil
}

func (e EIAData) validateDuoArea() error {
	if e.ResponseData.FuelData[0].DuoArea != "NUS" {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected DuoArea to be NUS, received %s", e.ResponseData.FuelData[0].DuoArea))
	}
	return nil
}

func (e EIAData) validateAreaName() error {
	if e.ResponseData.FuelData[0].AreaName != "U.S." {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected AreaName to be U.S., received %s", e.ResponseData.FuelData[0].AreaName))
	}
	return nil
}

func (e EIAData) validateProduct() error {
	if e.ResponseData.FuelData[0].Product != "EPD2D" {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected Product to be EPD2D, received %s", e.ResponseData.FuelData[0].Product))
	}
	return nil
}

func (e EIAData) validateProcess() error {
	if e.ResponseData.FuelData[0].Process != "PTE" {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected Process to be PTE, received %s", e.ResponseData.FuelData[0].Process))
	}
	return nil
}

func (e EIAData) validateSeries() error {
	if e.ResponseData.FuelData[0].Series != "EMD_EPD2D_PTE_NUS_DPG" {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected Series to be EMD_EPD2D_PTE_NUS_DPG, received %s", e.ResponseData.FuelData[0].Series))
	}
	return nil
}

func (e EIAData) validateUnits() error {
	if e.ResponseData.FuelData[0].Units != "$/GAL" {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected Units to be $/GAL, received %s", e.ResponseData.FuelData[0].Units))
	}
	return nil
}

func (e EIAData) validateDateFormat() error {
	if e.ResponseData.DateFormat != "YYYY-MM-DD" {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected DateFormat to be YYYY-MM-DD, received %s", e.ResponseData.DateFormat))
	}
	return nil
}

func (e EIAData) validateFrequency() error {
	if e.ResponseData.Frequency != "weekly" {
		return NewGHCAPIValidationError(
			fmt.Sprintf("Expected Frequency to be weekly, received %s", e.ResponseData.Frequency))
	}
	return nil
}

func (e EIAData) validateFuelDataExists() error {
	if len(e.ResponseData.FuelData) == 0 {
		return NewGHCAPIValidationError("received empty array of fuel data")
	}
	return nil
}
