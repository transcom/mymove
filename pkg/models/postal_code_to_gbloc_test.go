package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_FetchGBLOCForPostalCode() {
	t := suite.T()
	postalCodeToGBLOC := factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "77777", "UUUU")

	gbloc, err := models.FetchGBLOCForPostalCode(suite.DB(), postalCodeToGBLOC.PostalCode)
	if err != nil {
		t.Errorf("Find GBLOC for Postal Code error: %v", err)
	}

	if gbloc.GBLOC != "UUUU" {
		t.Errorf("GBLOC should be loaded for Postal Code %v", postalCodeToGBLOC.PostalCode)
	}
}
