package transittime

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *TransitTimeParserSuite) Test_xlsxDataSheetInfo_generateOutputFilename() {

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
				description:    models.StringPointer("test_case_1"),
				outputFilename: models.StringPointer("test_case_1"),
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
			want: "1_transit_time_ghc_parse_" + currentTime.Format("20060102150405") + ".csv",
		},
		{
			name: "TC 3: generate filename with suffix",
			args: args{
				index:   2,
				runTime: currentTime,
			},
			adtlSuffix: models.StringPointer("adtlSuffix"),
			want:       "2_transit_time_ghc_parse_adtlSuffix_" + currentTime.Format("20060102150405") + ".csv",
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
			if got := x.generateOutputFilename(tt.args.index, tt.args.runTime, tt.adtlSuffix); got != tt.want {
				suite.T().Errorf("xlsxDataSheetInfo.generateOutputFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

type TestStruct1 struct{ Field1 string }
type TestStruct2 struct{ Field1 string }
type TestStruct4 struct{ Field1 string }
type TestStruct5 struct{ Field1 string }
type TestStruct6 struct{ Field1 string }

var testVerifyFunc1 verifyXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) error {
	return nil
}

var testVerifyFunc2 verifyXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) error {
	return nil
}

var testVerifyFunc3 verifyXlsxSheet = func(_ ParamConfig, sheetIndex int, _ *zap.Logger) error {
	return fmt.Errorf("forced test error from function testVerifyFunc3 with index %d", sheetIndex)
}

var testVerifyFunc4 verifyXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) error {
	return nil
}

var testProcessFunc1 processXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) (interface{}, error) {
	return []TestStruct1{}, nil
}

var testProcessFunc2 processXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) (interface{}, error) {
	return []TestStruct2{}, nil
}

var testProcessFunc3 processXlsxSheet = func(_ ParamConfig, sheetIndex int, _ *zap.Logger) (interface{}, error) {
	return nil, fmt.Errorf("forced test error from function testProcessFunc3 with index %d", sheetIndex)
}

var testProcessFunc4 processXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) (interface{}, error) {
	return []TestStruct4{}, nil
}

var testProcessFunc5 processXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) (interface{}, error) {
	return []TestStruct5{}, nil
}

var testProcessFunc6 processXlsxSheet = func(_ ParamConfig, _ int, _ *zap.Logger) (interface{}, error) {
	return []TestStruct6{}, nil
}

func (suite *TransitTimeParserSuite) helperTestSetup() []XlsxDataSheetInfo {
	xlsxDataSheets := make([]XlsxDataSheetInfo, xlsxSheetsCountMax)

	// 0:
	xlsxDataSheets[0] = XlsxDataSheetInfo{
		Description:    models.StringPointer("0) Test Process 1"),
		outputFilename: models.StringPointer("0_test_process_1"),
		ProcessMethods: []xlsxProcessInfo{{
			process: &testProcessFunc1,
		},
		},
		verify: &testVerifyFunc1,
	}

	// 1:
	xlsxDataSheets[1] = XlsxDataSheetInfo{
		Description:    models.StringPointer("1) Test Process 2"),
		outputFilename: models.StringPointer("1_test_process_2"),
		ProcessMethods: []xlsxProcessInfo{{
			process: &testProcessFunc2,
		},
		},
		verify: &testVerifyFunc2,
	}

	// 2:
	xlsxDataSheets[2] = XlsxDataSheetInfo{
		Description:    models.StringPointer("2) Test Process 3"),
		outputFilename: models.StringPointer("2_test_process_3"),
		ProcessMethods: []xlsxProcessInfo{{
			process: &testProcessFunc3,
		},
		},
		verify: &testVerifyFunc3,
	}

	// 3:
	xlsxDataSheets[3] = XlsxDataSheetInfo{
		Description:    models.StringPointer("3) Test Process 4"),
		outputFilename: models.StringPointer("3_test_process_4"),
		ProcessMethods: []xlsxProcessInfo{
			{
				process:    &testProcessFunc4,
				adtlSuffix: models.StringPointer("suffix1"),
			},
			{
				process:    &testProcessFunc5,
				adtlSuffix: models.StringPointer("suffix2"),
			},
			{
				process:    &testProcessFunc6,
				adtlSuffix: models.StringPointer("suffix4"),
			},
		},
		verify: &testVerifyFunc4,
	}

	return xlsxDataSheets
}

func (suite *TransitTimeParserSuite) Test_getInt() {
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
		suite.Run(tt.name, func() {
			if got := getInt(tt.args.from); got != tt.want {
				suite.T().Errorf("getInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (suite *TransitTimeParserSuite) Test_removeFirstDollarSign() {
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
			if got := removeFirstDollarSign(tt.args.s); got != tt.want {
				suite.T().Errorf("removeFirstDollarSign() = %v, want %v", got, tt.want)
			}
		})
	}
}
