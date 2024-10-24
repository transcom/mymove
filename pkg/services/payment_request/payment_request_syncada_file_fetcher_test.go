package paymentrequest

import (
	"reflect"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testingsuite"
)

func (suite *PaymentRequestSyncadaFileFetcherSuite) TestFetchPaymentRequestSyncadaFile() {
	builder := query.NewQueryBuilder()
	fetcher := NewPaymentRequestSyncadaFileFetcher(builder)

	suite.Run("Fetch Syncada files", func() {
		paymentRequestEdiFile := BuildPaymentRequestEdiRecord("858.rec1", "someStringedi", "1234-7654-1")
		err := suite.DB().Create(&paymentRequestEdiFile)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{})
		suite.NoError(err)
		suite.NotNil(result)
	})
	type fields struct {
		builder paymentReqeustSyncadaFileQueryBuilder
	}
	type args struct {
		appCtx  appcontext.AppContext
		filters []services.QueryFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.PaymentRequestEdiFile
		wantErr bool
	}{
		{
			name: "Successful fetch",
			fields: fields{
				builder: query.NewQueryBuilder(),
			},
			args: args{
				appCtx:  suite.AppContextForTest(),
				filters: []services.QueryFilter{},
			},
			want:    BuildPaymentRequestEdiRecord("858.rec1", "someStringedi", "1234-7654-1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			p := &paymentRequestSyncadaFileFetcher{
				builder: tt.fields.builder,
			}
			got, err := p.FetchPaymentRequestSyncadaFile(tt.args.appCtx, tt.args.filters)
			if (err != nil) != tt.wantErr {
				suite.Equal("paymentRequestSyncadaFileFetcher.FetchPaymentRequestSyncadaFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				suite.Equal("paymentRequestSyncadaFileFetcher.FetchPaymentRequestSyncadaFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

type PaymentRequestSyncadaFileFetcherSuite struct {
	testingsuite.PopTestSuite
}

func BuildPaymentRequestEdiRecord(fileName string, ediString string, prNumber string) models.PaymentRequestEdiFile {
	var paymentRequestEdiFile models.PaymentRequestEdiFile
	paymentRequestEdiFile.ID = uuid.Must(uuid.NewV4())
	paymentRequestEdiFile.EdiString = ediString
	paymentRequestEdiFile.Filename = fileName
	paymentRequestEdiFile.PaymentRequestNumber = prNumber

	return paymentRequestEdiFile
}
func (suite *PaymentRequestSyncadaFileFetcherSuite) TestFetchPaymentRequestSyncadaFileEdgeCases() {
	builder := query.NewQueryBuilder()
	fetcher := NewPaymentRequestSyncadaFileFetcher(builder)

	suite.Run("Fetch non-existent Syncada file", func() {
		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("id", "=", uuid.Must(uuid.NewV4())),
		})

		suite.Error(err)
		suite.Equal(models.PaymentRequestEdiFile{}, result)
	})

	suite.Run("Fetch Syncada file with specific filters", func() {
		paymentRequestEdiFile := BuildPaymentRequestEdiRecord("858.rec2", "testEdiString", "9876-5432-1")
		err := suite.DB().Create(&paymentRequestEdiFile)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("filename", "=", "858.rec2"),
			query.NewQueryFilter("payment_request_number", "=", "9876-5432-1"),
		})

		suite.NoError(err)
		suite.NotNil(result)
		suite.Equal(paymentRequestEdiFile.ID, result.ID)
		suite.Equal("858.rec2", result.Filename)
		suite.Equal("testEdiString", result.EdiString)
		suite.Equal("9876-5432-1", result.PaymentRequestNumber)
	})

	suite.Run("Fetch Syncada file with invalid filter", func() {
		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("invalid_column", "=", "some_value"),
		})

		suite.Error(err)
		suite.Equal(models.PaymentRequestEdiFile{}, result)
	})

	suite.Run("Fetch Syncada file with multiple matching records", func() {
		paymentRequestEdiFile1 := BuildPaymentRequestEdiRecord("858.rec3", "testEdiString1", "1111-1111-1")
		paymentRequestEdiFile2 := BuildPaymentRequestEdiRecord("858.rec3", "testEdiString2", "1111-1111-1")
		err := suite.DB().Create(&paymentRequestEdiFile1)
		suite.NoError(err)
		err = suite.DB().Create(&paymentRequestEdiFile2)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("filename", "=", "858.rec3"),
			query.NewQueryFilter("payment_request_number", "=", "1111-1111-1"),
		})

		suite.Error(err)
		suite.Equal(models.PaymentRequestEdiFile{}, result)
	})
}

func (suite *PaymentRequestSyncadaFileFetcherSuite) TestFetchPaymentRequestSyncadaFile_justOne() {
	builder := query.NewQueryBuilder()
	fetcher := NewPaymentRequestSyncadaFileFetcher(builder)

	testCases := []struct {
		name    string
		filters []services.QueryFilter
		want    models.PaymentRequestEdiFile
		wantErr bool
	}{
		{
			name:    "Fetch Syncada files",
			filters: []services.QueryFilter{},
			want:    BuildPaymentRequestEdiRecord("858.rec1", "someStringedi", "1234-7654-1"),
			wantErr: false,
		},
		{
			name: "Successful fetch with specific filters",
			filters: []services.QueryFilter{
				query.NewQueryFilter("filename", "=", "858.rec3"),
				query.NewQueryFilter("payment_request_number", "=", "1111-1111-1"),
			},
			want:    BuildPaymentRequestEdiRecord("858.rec3", "testEdiString1", "1111-1111-1"),
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.name == "Fetch Syncada files" {
				paymentRequestEdiFile := tc.want
				err := suite.DB().Create(&paymentRequestEdiFile)
				suite.NoError(err)
			}

			result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), tc.filters)

			if tc.wantErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
				suite.NotNil(result)
				suite.Equal(tc.want.Filename, result.Filename)
				suite.Equal(tc.want.EdiString, result.EdiString)
				suite.Equal(tc.want.PaymentRequestNumber, result.PaymentRequestNumber)
			}
		})
	}
}

func (suite *PaymentRequestSyncadaFileFetcherSuite) TestFetchPaymentRequestSyncadaFile_NewCases() {
	builder := query.NewQueryBuilder()
	fetcher := NewPaymentRequestSyncadaFileFetcher(builder)

	suite.Run("Fetch Syncada file with partial match filter", func() {
		paymentRequestEdiFile := BuildPaymentRequestEdiRecord("858.rec4", "partialMatchTest", "2222-3333-4")
		err := suite.DB().Create(&paymentRequestEdiFile)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("filename", "LIKE", "%rec4%"),
		})

		suite.NoError(err)
		suite.NotNil(result)
		suite.Equal(paymentRequestEdiFile.ID, result.ID)
		suite.Equal("858.rec4", result.Filename)
		suite.Equal("partialMatchTest", result.EdiString)
		suite.Equal("2222-3333-4", result.PaymentRequestNumber)
	})

	suite.Run("Fetch Syncada file with case-insensitive filter", func() {
		paymentRequestEdiFile := BuildPaymentRequestEdiRecord("UPPERCASE.REC", "caseInsensitiveTest", "5555-6666-7")
		err := suite.DB().Create(&paymentRequestEdiFile)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("filename", "ILIKE", "uppercase.rec"),
		})

		suite.NoError(err)
		suite.NotNil(result)
		suite.Equal(paymentRequestEdiFile.ID, result.ID)
		suite.Equal("UPPERCASE.REC", result.Filename)
		suite.Equal("caseInsensitiveTest", result.EdiString)
		suite.Equal("5555-6666-7", result.PaymentRequestNumber)
	})

	suite.Run("Fetch Syncada file with multiple filters", func() {
		paymentRequestEdiFile := BuildPaymentRequestEdiRecord("multi.filter.rec", "multiFilterTest", "7777-8888-9")
		err := suite.DB().Create(&paymentRequestEdiFile)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("filename", "=", "multi.filter.rec"),
			query.NewQueryFilter("payment_request_number", "=", "7777-8888-9"),
			query.NewQueryFilter("edi_string", "=", "multiFilterTest"),
		})

		suite.NoError(err)
		suite.NotNil(result)
		suite.Equal(paymentRequestEdiFile.ID, result.ID)
		suite.Equal("multi.filter.rec", result.Filename)
		suite.Equal("multiFilterTest", result.EdiString)
		suite.Equal("7777-8888-9", result.PaymentRequestNumber)
	})

	suite.Run("Fetch Syncada file with empty filters", func() {
		paymentRequestEdiFile := BuildPaymentRequestEdiRecord("empty.filter.rec", "emptyFilterTest", "9999-0000-1")
		err := suite.DB().Create(&paymentRequestEdiFile)
		suite.NoError(err)

		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{})

		suite.NoError(err)
		suite.NotNil(result)
		suite.NotEqual(uuid.Nil, result.ID)
	})

	suite.Run("Fetch Syncada file with invalid operator in filter", func() {
		result, err := fetcher.FetchPaymentRequestSyncadaFile(suite.AppContextForTest(), []services.QueryFilter{
			query.NewQueryFilter("filename", "INVALID_OPERATOR", "some_value"),
		})

		suite.Error(err)
		suite.Equal(models.PaymentRequestEdiFile{}, result)
	})
}

func (suite *PaymentRequestSyncadaFileFetcherSuite) Test_paymentRequestSyncadaFileFetcher_FetchPaymentRequestSyncadaFile() {

	type fields struct {
		builder query.Builder
	}
	type args struct {
		appCtx  appcontext.AppContext
		filters []services.QueryFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.PaymentRequestEdiFile
		wantErr bool
	}{
		{
			name: "Fetch Syncada files",
			fields: fields{
				builder: <-suite.AppContextForTest().DB().Context().Done(),
			},
			args: args{
				appCtx: suite.AppContextForTest(),
				filters: []services.QueryFilter{
					query.NewQueryFilter("filename", "=", "858.rec3"),
					query.NewQueryFilter("payment_request_number", "=", "1111-1111-1"),
				},
			},
			want:    BuildPaymentRequestEdiRecord("858.rec3", "testEdiString1", "1111-1111-1"),
			wantErr: false,
		},
		{
			name: "Fetch Syncada files with partial match filter",
			fields: fields{
				builder: <-suite.AppContextForTest().DB().Context().Done(),
			},
			args: args{
				appCtx: suite.AppContextForTest(),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			p := &paymentRequestSyncadaFileFetcher{
				builder: &tt.fields.builder,
			}
			got, err := p.FetchPaymentRequestSyncadaFile(tt.args.appCtx, tt.args.filters)
			if (err != nil) != tt.wantErr {
				suite.Equal("paymentRequestSyncadaFileFetcher.FetchPaymentRequestSyncadaFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			suite.ElementsMatch(got, tt.want)
		})
	}
}
