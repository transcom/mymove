package pricing

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/suite"
	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services/dbtools"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PricingParserSuite struct {
	*testingsuite.PopTestSuite
	xlsxFilename string
	xlsxFile     *xlsx.File
}

func TestPricingParserSuite(t *testing.T) {
	hs := &PricingParserSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		xlsxFilename: "fixtures/pricing_template_2019-09-19_fake-data.xlsx",
	}

	var err error
	hs.xlsxFile, err = xlsx.OpenFile(hs.xlsxFilename)
	if err != nil {
		hs.Logger().Panic("could not open XLSX file", zap.Error(err))
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *PricingParserSuite) Test_xlsxDataSheetInfo_generateOutputFilename() {

	type fields struct {
		description    *string
		process        *processXlsxSheet
		verify         *verifyXlsxSheet
		outputFilename *string
	}
	type args struct {
		index   int
		runTime time.Time
	}

	currentTime := time.Now()

	tests := []struct {
		name       string
		fields     fields
		args       args
		adtlSuffix *string
		want       string
	}{
		{
			name: "TC 1: generate filename with outputFilename provided",
			fields: fields{
				description:    swag.String("test_case_1"),
				outputFilename: swag.String("test_case_1"),
				// process not needed for this function
				// verify not needed for this function
			},
			args: args{
				index:   0,
				runTime: currentTime,
			},
			want: "0_test_case_1_" + currentTime.Format("20060102150405") + ".csv",
		},
		{
			name: "TC 2: generate filename with outputFilename not provided",
			args: args{
				index:   1,
				runTime: currentTime,
			},
			want: "1_rate_engine_ghc_parse_" + currentTime.Format("20060102150405") + ".csv",
		},
		{
			name: "TC 3: generate filename with suffix",
			args: args{
				index:   2,
				runTime: currentTime,
			},
			adtlSuffix: swag.String("adtlSuffix"),
			want:       "2_rate_engine_ghc_parse_adtlSuffix_" + currentTime.Format("20060102150405") + ".csv",
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			x := &XlsxDataSheetInfo{
				Description: tt.fields.description,
				ProcessMethods: []xlsxProcessInfo{{
					process:    tt.fields.process,
					adtlSuffix: tt.adtlSuffix,
				},
				},
				verify:         tt.fields.verify,
				outputFilename: tt.fields.outputFilename,
			}
			got := x.generateOutputFilename(tt.args.index, tt.args.runTime, tt.adtlSuffix)
			suite.Equal(got, tt.want, "xlsxDataSheetInfo.generateOutputFilename()")
		})
	}
}

type TestStruct1 struct{ Field1 string }
type TestStruct2 struct{ Field1 string }
type TestStruct4 struct{ Field1 string }
type TestStruct5 struct{ Field1 string }
type TestStruct6 struct{ Field1 string }

var testVerifyFunc1 verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	return nil
}

var testVerifyFunc2 verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	return nil
}

var testVerifyFunc3 verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	return fmt.Errorf("forced test error from function testVerifyFunc3 with index %d", sheetIndex)
}

var testVerifyFunc4 verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	return nil
}

var testProcessFunc1 processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	return []TestStruct1{}, nil
}

var testProcessFunc2 processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	return []TestStruct2{}, nil
}

var testProcessFunc3 processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	return nil, fmt.Errorf("forced test error from function testProcessFunc3 with index %d", sheetIndex)
}

var testProcessFunc4 processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	return []TestStruct4{}, nil
}

var testProcessFunc5 processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	return []TestStruct5{}, nil
}

var testProcessFunc6 processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	return []TestStruct6{}, nil
}

func (suite *PricingParserSuite) helperTestSetup() []XlsxDataSheetInfo {
	xlsxDataSheets := make([]XlsxDataSheetInfo, xlsxSheetsCountMax)

	// 0:
	xlsxDataSheets[0] = XlsxDataSheetInfo{
		Description:    swag.String("0) Test Process 1"),
		outputFilename: swag.String("0_test_process_1"),
		ProcessMethods: []xlsxProcessInfo{{
			process: &testProcessFunc1,
		},
		},
		verify: &testVerifyFunc1,
	}

	// 1:
	xlsxDataSheets[1] = XlsxDataSheetInfo{
		Description:    swag.String("1) Test Process 2"),
		outputFilename: swag.String("1_test_process_2"),
		ProcessMethods: []xlsxProcessInfo{{
			process: &testProcessFunc2,
		},
		},
		verify: &testVerifyFunc2,
	}

	// 2:
	xlsxDataSheets[2] = XlsxDataSheetInfo{
		Description:    swag.String("2) Test Process 3"),
		outputFilename: swag.String("2_test_process_3"),
		ProcessMethods: []xlsxProcessInfo{{
			process: &testProcessFunc3,
		},
		},
		verify: &testVerifyFunc3,
	}

	// 3:
	xlsxDataSheets[3] = XlsxDataSheetInfo{
		Description:    swag.String("3) Test Process 4"),
		outputFilename: swag.String("3_test_process_4"),
		ProcessMethods: []xlsxProcessInfo{
			{
				process:     &testProcessFunc4,
				description: swag.String("suffix1 description"),
				adtlSuffix:  swag.String("suffix1"),
			},
			{
				process:     &testProcessFunc5,
				description: swag.String("suffix2 description"),
				adtlSuffix:  swag.String("suffix2"),
			},
			{
				process:     &testProcessFunc6,
				description: swag.String("suffix4 description"),
				adtlSuffix:  swag.String("suffix4"),
			},
		},
		verify: &testVerifyFunc4,
	}

	return xlsxDataSheets
}

func (suite *PricingParserSuite) Test_process() {

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
			tableFromSliceCreator := dbtools.NewTableFromSliceCreator(true, false)

			err := process(suite.AppContextForTest(), xlsxDataSheets, tt.args.params, tt.args.sheetIndex,
				tableFromSliceCreator)
			suite.Equal(err != nil, tt.wantErr, "process() err %v", err)
		})
	}
}

func (suite *PricingParserSuite) Test_getInt() {
	type args struct {
		from string
	}
	tests := []struct {
		name     string
		args     args
		want     int
		hasError bool
	}{
		{
			name: "TC 1: convert string 1",
			args: args{
				from: "1",
			},
			want:     1,
			hasError: false,
		},
		{
			name: "TC 2: convert string 1.0",
			args: args{
				from: "1.0",
			},
			want:     1,
			hasError: false,
		},
		{
			name: "TC 3: convert string 1sldkjf",
			args: args{
				from: "1sldkjf",
			},
			want:     0,
			hasError: true,
		},
		{
			name: "TC 4: convert string 10.sldk",
			args: args{
				from: "10.sldk",
			},
			want:     0,
			hasError: true,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, gotErr := getInt(tt.args.from)
			suite.Equal(tt.want, got)
			if tt.hasError {
				suite.Error(gotErr)
			} else {
				suite.NoError(gotErr)
			}
		})
	}
}

func (suite *PricingParserSuite) Test_removeFirstDollarSign() {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TC 1: $$$",
			args: args{
				s: "$$$",
			},
			want: "$$",
		},
		{
			name: "TC 2: $3.99",
			args: args{
				s: "$3.99",
			},
			want: "3.99",
		},
		{
			name: "TC 2: $3.$99",

			args: args{
				s: "$3.$99",
			},
			want: "3.$99",
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := removeFirstDollarSign(tt.args.s)
			suite.Equal(got, tt.want, "removeFirstDollarSign")
		})
	}
}

func (suite *PricingParserSuite) helperTestExpectedFileOutput(goldenFilename string, currentOutputFilename string) {
	expected := filepath.Join("fixtures", goldenFilename) // relative path
	expectedBytes, err := ioutil.ReadFile(filepath.Clean(expected))
	suite.NoErrorf(err, "error loading expected CSV file output fixture <%s>", expected)

	currentBytes, err := ioutil.ReadFile(filepath.Clean(currentOutputFilename)) // relative path
	suite.NoErrorf(err, "error loading current/new output file <%s>", currentOutputFilename)

	suite.Equal(string(expectedBytes), string(currentBytes))

	// Remove file generated from test after compare is finished
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to clean up file created for unit test
	//RA: Given the functions causing the lint errors are used to clean up local storage space after a unit test, it does not present a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	os.Remove(currentOutputFilename)
}
