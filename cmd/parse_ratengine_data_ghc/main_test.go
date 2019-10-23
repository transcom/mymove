package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ParseRateEngineGHCXLSXSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *ParseRateEngineGHCXLSXSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestParseRateEngineGHCXLSXSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &ParseRateEngineGHCXLSXSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *ParseRateEngineGHCXLSXSuite) Test_xlsxDataSheetInfo_generateOutputFilename() {

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
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "TC 1: generate filename with outputFilename provided",
			fields: fields{
				description:    stringPointer("test_case_1"),
				outputFilename: stringPointer("test_case_1"),
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
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			x := &xlsxDataSheetInfo{
				description:    tt.fields.description,
				process:        tt.fields.process,
				verify:         tt.fields.verify,
				outputFilename: tt.fields.outputFilename,
			}
			if got := x.generateOutputFilename(tt.args.index, tt.args.runTime); got != tt.want {
				t.Errorf("xlsxDataSheetInfo.generateOutputFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testVerifyFunc1 verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {
	return nil
}

var testVerifyFunc2 verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {
	return nil
}

var testVerifyFunc3 verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {
	return fmt.Errorf("forced test error from function testVerifyFunc3 with index %d", sheetIndex)
}

var testProcessFunc1 processXlsxSheet = func(params paramConfig, sheetIndex int) error {
	return nil
}

var testProcessFunc2 processXlsxSheet = func(params paramConfig, sheetIndex int) error {
	return nil
}

var testProcessFunc3 processXlsxSheet = func(params paramConfig, sheetIndex int) error {
	return fmt.Errorf("forced test error from function testProcessFunc3 with index %d", sheetIndex)
}

func (suite *ParseRateEngineGHCXLSXSuite) helperTestSetup() {
	xlsxDataSheets = make([]xlsxDataSheetInfo, xlsxSheetsCountMax, xlsxSheetsCountMax)

	// 0:
	xlsxDataSheets[0] = xlsxDataSheetInfo{
		description:    stringPointer("0) Test Process 1"),
		outputFilename: stringPointer("0_test_process_1"),
		process:        &testProcessFunc1,
		verify:         &testVerifyFunc1,
	}

	// 1:
	xlsxDataSheets[1] = xlsxDataSheetInfo{
		description:    stringPointer("1) Test Process 2"),
		outputFilename: stringPointer("1_test_process_2"),
		process:        &testProcessFunc2,
		verify:         &testVerifyFunc2,
	}

	// 2:
	xlsxDataSheets[2] = xlsxDataSheetInfo{
		description:    stringPointer("2) Test Process 3"),
		outputFilename: stringPointer("2_test_process_3"),
		process:        &testProcessFunc3,
		verify:         &testVerifyFunc3,
	}
}

func (suite *ParseRateEngineGHCXLSXSuite) Test_process() {

	suite.helperTestSetup()

	type args struct {
		params     paramConfig
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
				params: paramConfig{
					runTime: time.Now(),
				},
				sheetIndex: 0,
			},
			wantErr: false,
		},
		{
			name: "TC 2 run fake process & verify function 2, no error",
			args: args{
				params: paramConfig{
					runTime: time.Now(),
				},
				sheetIndex: 1,
			},
			wantErr: false,
		},
		{
			name: "TC 3 run fake process & verify function 3, with error",
			args: args{
				params: paramConfig{
					runTime: time.Now(),
				},
				sheetIndex: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if err := process(tt.args.params, tt.args.sheetIndex); (err != nil) != tt.wantErr {
				t.Errorf("process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (suite *ParseRateEngineGHCXLSXSuite) Test_getInt() {
	type args struct {
		from string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TC 1: convert string 1",
			args: args{
				from: "1",
			},
			want: 1,
		},
		{
			name: "TC 2: convert string 1.0",
			args: args{
				from: "1.0",
			},
			want: 1,
		},
		{
			name: "TC 3: convert string 1sldkjf",
			args: args{
				from: "1sldkjf",
			},
			want: 0,
		},
		{
			name: "TC 4: convert string 10.sldk",
			args: args{
				from: "10.sldk",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if got := getInt(tt.args.from); got != tt.want {
				t.Errorf("getInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *ParseRateEngineGHCXLSXSuite) Test_removeFirstDollarSign() {
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
		suite.T().Run(tt.name, func(t *testing.T) {
			if got := removeFirstDollarSign(tt.args.s); got != tt.want {
				t.Errorf("removeFirstDollarSign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *ParseRateEngineGHCXLSXSuite) helperTestExpectedFileOutput(goldenFilename string, currentOutputFilename string) {
	expected := filepath.Join("fixtures", goldenFilename) // relative path
	expectedBytes, err := ioutil.ReadFile(expected)
	suite.NoErrorf(err, "error loading expected CSV file output fixture <%s>", expected)

	currentBytes, err := ioutil.ReadFile(currentOutputFilename) // relative path
	suite.NoErrorf(err, "error loading current/new output file <%s>", currentOutputFilename)

	suite.Equal(string(expectedBytes), string(currentBytes))

	// Remove file generated from test after compare is finished
	os.Remove(currentOutputFilename)
}
