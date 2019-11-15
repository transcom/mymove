package ghcimport

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) helperImportRERateAreaTC2(action string) {
	var err error
	// Update domestic US6B name "Texas-South" to something else and verify it was changed back when done
	var texas *models.ReRateArea
	texas, err = models.FetchReRateAreaItem(suite.DB(), "US6B")
	suite.NoError(err)
	//suite.NotNil(texas)
	suite.FatalNil(texas)
	fmt.Printf("\nFetch US6B rate area %v\n\n", texas)
	suite.Equal("Texas-South", texas.Name)

	// Update oconus US8101000 name "Alaska (Zone) I" to something else and verify it was changed back when done
	var alaska *models.ReRateArea
	alaska, err = models.FetchReRateAreaItem(suite.DB(), "US8101000")
	suite.NoError(err)
	suite.NotNil(alaska)
	suite.Equal("Alaska (Zone) I", alaska.Name)

	// Update oconus AS11 name "New South Wales/Australian Capital Territory"
	var wales *models.ReRateArea
	wales, err = models.FetchReRateAreaItem(suite.DB(), "AS11")
	suite.NoError(err)
	suite.NotNil(wales)
	suite.Equal("New South Wales/Australian Capital Territory", wales.Name)

	if action == "setup" {
		modifiedName := "New name"
		texas.Name = modifiedName
		suite.MustSave(texas)
		texas, err = models.FetchReRateAreaItem(suite.DB(), "US6B")
		suite.NoError(err)
		suite.Equal(modifiedName, texas.Name)

		modifiedName = "New name 2"
		alaska.Name = modifiedName
		suite.MustSave(alaska)
		alaska, err = models.FetchReRateAreaItem(suite.DB(), "US8101000")
		suite.NoError(err)
		suite.Equal(modifiedName, alaska.Name)

		modifiedName = "New name 3"
		wales.Name = modifiedName
		suite.MustSave(wales)
		wales, err = models.FetchReRateAreaItem(suite.DB(), "AS11")
		suite.NoError(err)
		suite.Equal(modifiedName, wales.Name)
	} else if action == "verify" {
		// nothing to do, verify happens at the top
	}
}

func (suite *GHCRateEngineImportSuite) helperImportRERateAreaTC3(action string) {
	if action == "setup" {
		// drop a staging table that we are depending on to do import
		dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", "stage_conus_to_oconus_prices")
		dropErr := suite.DB().RawQuery(dropQuery).Exec()
		suite.NoError(dropErr)
	}
}

func (suite *GHCRateEngineImportSuite) helperImportRERateAreaVerifyImportComplete() {
	var rateArea models.ReRateArea
	count, countErr := suite.DB().Count(&rateArea)
	suite.NoError(countErr)
	suite.Equal(120, count)
}

func (suite *GHCRateEngineImportSuite) TestGHCRateEngineImporter_importRERateArea() {
	type fields struct {
		Logger Logger
	}
	type args struct {
		dbTx *pop.Connection
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TC0: Successfully run import with staged staging data (empty RE tables)",
			fields: fields{
				Logger: suite.logger,
			},
			args: args{
				dbTx: suite.DB(),
			},
			wantErr: false,
		},
		{
			name: "TC1: Successfully run import, 2nd time, with staged staging data and filled in RE tables",
			fields: fields{
				Logger: suite.logger,
			},
			args: args{
				dbTx: suite.DB(),
			},
			wantErr: false,
		},
		{
			name: "TC2: Successfully run import, prefilled re_rate_areas, update existing rate area from import",
			fields: fields{
				Logger: suite.logger,
			},
			args: args{
				dbTx: suite.DB(),
			},
			wantErr: false,
		},
		{
			name: "TC3: Fail to run import, missing staging table",
			fields: fields{
				Logger: suite.logger,
			},
			args: args{
				dbTx: suite.DB(),
			},
			wantErr: true,
		},
	}
	for tc, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			gre := &GHCRateEngineImporter{
				Logger: tt.fields.Logger,
			}
			// Run any necessary setup functions
			if tc == 2 {
				//suite.NoError(gre.importRERateArea(tt.args.dbTx))
				//suite.helperImportRERateAreaTC2("setup")
			} else if tc == 3 {
				//suite.helperImportRERateAreaTC3("setup")
			}
			// Execute function under test
			if err := gre.importRERateArea(tt.args.dbTx); (err != nil) != tt.wantErr {
				t.Errorf("GHCRateEngineImporter.importRERateArea() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Run any additional verify functions
			suite.helperImportRERateAreaVerifyImportComplete()
			if tc == 2 {
				//suite.helperImportRERateAreaTC2("verify")
			}
		})
	}
}