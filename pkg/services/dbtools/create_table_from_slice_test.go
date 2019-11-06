package dbtools

import (
	"testing"
)

func (suite *DBToolsServiceSuite) TestCreateTableFromSlice() {
	tableFromSliceCreator := NewTableFromSliceCreator(suite.DB(), suite.logger, true)

	suite.T().Run("passing in a non-slice", func(t *testing.T) {
		err := tableFromSliceCreator.CreateTableFromSlice(1)
		suite.Error(err)
		suite.Equal("Parameter must be slice or array, but got int", err.Error())
	})

	suite.T().Run("passing in a slice, but not a slice of structs", func(t *testing.T) {
		err := tableFromSliceCreator.CreateTableFromSlice([]int{1, 2, 3})
		suite.Error(err)
		suite.Equal("Elements of slice must be type struct, but got int", err.Error())
	})

	suite.T().Run("passing in a slice of structs, but with a non-string field", func(t *testing.T) {
		var invalidStructSlice []struct {
			field1 string
			field2 int
		}
		err := tableFromSliceCreator.CreateTableFromSlice(invalidStructSlice)
		suite.Error(err)
		suite.Equal("All fields of struct must be string, but field field2 is int", err.Error())
	})

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

	suite.T().Run("valid slice of structs", func(t *testing.T) {
		err := tableFromSliceCreator.CreateTableFromSlice(validSlice)
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
