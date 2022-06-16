package dbtools

import (
	"github.com/transcom/mymove/pkg/appcontext"
)

type TestStruct struct {
	Name              string `db:"name"`
	ServiceAreaNumber string `db:"service_area_number"`
	Zip3              string `db:"zip3"`
}

var validSlice = []TestStruct{
	// Order matters here for test comparison
	{
		Name:              "Amanda",
		ServiceAreaNumber: "120",
		Zip3:              "292",
	},
	{
		Name:              "James",
		ServiceAreaNumber: "444",
		Zip3:              "361",
	},
	{
		Name:              "John",
		ServiceAreaNumber: "004",
		Zip3:              "309",
	},
}

func (suite *DBToolsServiceSuite) TestCreateTableFromSlice() {
	tableFromSliceCreator := NewTableFromSliceCreator(true, false)

	suite.Run("passing in a non-slice", func() {
		err := tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), 1)
		suite.Error(err)
		suite.Equal("Parameter must be slice or array, but got int", err.Error())
	})

	suite.Run("passing in a slice, but not a slice of structs", func() {
		err := tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), []int{1, 2, 3})
		suite.Error(err)
		suite.Equal("Elements of slice must be type struct, but got int", err.Error())
	})

	suite.Run("passing in a slice of structs, but with a non-string field", func() {
		var invalidStructSlice []struct {
			field1 string
			field2 int
		}
		err := tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), invalidStructSlice)
		suite.Error(err)
		suite.Equal("All fields of struct must be string, but field field2 is int", err.Error())
	})

	suite.Run("valid slice of structs", func() {
		err := tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), validSlice)
		suite.NoError(err)

		var testStructs []TestStruct
		err = suite.DB().Order("name").All(&testStructs)
		suite.NoError(err)
		suite.Len(testStructs, 3)
		for i, testStruct := range testStructs {
			suite.Equal(validSlice[i], testStruct)
		}
	})

	suite.Run("errors out when table exists", func() {
		err := tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), validSlice)
		suite.NoError(err)
		err = tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), validSlice)
		suite.Error(err)
		// TODO: Fix this DB error string literal comparison when we move the COPY-related functionality to jackc/pgx.
		if err != nil {
			suite.Equal("Error creating table: 'test_structs': pq: relation \"test_structs\" already exists", err.Error())
		}
	})
}

func (suite *DBToolsServiceSuite) TestCreateTableFromSlicePermTable() {
	tableFromSliceCreator := NewTableFromSliceCreator(true, true)

	suite.Run("two runs no error when drop flag is true", func() {
		err := tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), validSlice)
		suite.NoError(err)
		err = tableFromSliceCreator.CreateTableFromSlice(suite.AppContextForTest(), validSlice)
		suite.NoError(err)

		var testStructs []TestStruct
		err = suite.DB().Order("name").All(&testStructs)
		suite.NoError(err)
		suite.Len(testStructs, 3)
		for i, testStruct := range testStructs {
			suite.Equal(validSlice[i], testStruct)
		}
	})
}

func (suite *DBToolsServiceSuite) TestCreateTableFromSliceWithinTransaction() {
	var testStructs []TestStruct
	suite.Run("create table from slice in a transaction", func() {
		txnErr := suite.AppContextForTest().NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			tableFromSliceCreator := NewTableFromSliceCreator(true, true)
			err := tableFromSliceCreator.CreateTableFromSlice(txnAppCtx, validSlice)
			suite.NoError(err)

			err = txnAppCtx.DB().Order("name").All(&testStructs)
			suite.NoError(err)
			suite.Len(testStructs, 3)
			for i, testStruct := range testStructs {
				suite.Equal(validSlice[i], testStruct)
			}
			return nil
		})
		suite.NoError(txnErr)

		// Verify data still in table after transaction
		err := suite.DB().Order("name").All(&testStructs)
		suite.NoError(err)
		suite.Len(testStructs, 3)
		for i, testStruct := range testStructs {
			suite.Equal(validSlice[i], testStruct)
		}
	})
}
