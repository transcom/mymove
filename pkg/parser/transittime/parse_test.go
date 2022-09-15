package transittime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/dbtools"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TransitTimeParserSuite struct {
	*testingsuite.PopTestSuite
	tableFromSliceCreator services.TableFromSliceCreator
	xlsxFilename          string
	xlsxFile              *xlsx.File
}

func TestTransitTimeParserSuite(t *testing.T) {
	hs := &TransitTimeParserSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		xlsxFilename: "fixtures/Appendix_C(i)_-_Transit_Time_Tables_Fake_Data.xlsx",
	}

	hs.tableFromSliceCreator = dbtools.NewTableFromSliceCreator(true, false)

	var err error
	hs.xlsxFile, err = xlsx.OpenFile(hs.xlsxFilename)
	if err != nil {
		hs.Logger().Panic("could not open XLSX file", zap.Error(err))
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *TransitTimeParserSuite) Test_process() {

	xlsxDataSheets := suite.helperTestSetup()

	type args struct {
		params     ParamConfig
		sheetIndex int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TC 2 run fake process & verify function 1, no error",
			args: args{
				params: ParamConfig{
					RunTime: time.Now(),
				},
				sheetIndex: 0,
			},
			wantErr: false,
		},
		{
			name: "TC 2 run fake process & verify function 2, no error",
			args: args{
				params: ParamConfig{
					RunTime: time.Now(),
				},
				sheetIndex: 1,
			},
			wantErr: false,
		},
		{
			name: "TC 3 run fake process & verify function 3, with error",
			args: args{
				params: ParamConfig{
					RunTime: time.Now(),
				},
				sheetIndex: 2,
			},
			wantErr: true,
		},
		{
			name: "TC 4 run fake process methods & verify function 4, with suffix",
			args: args{
				params: ParamConfig{
					RunTime: time.Now(),
				},
				sheetIndex: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if err := process(suite.AppContextForTest(), xlsxDataSheets, tt.args.params, tt.args.sheetIndex, suite.tableFromSliceCreator); (err != nil) != tt.wantErr {
				suite.T().Errorf("process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
