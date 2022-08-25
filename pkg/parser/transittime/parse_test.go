package transittime

import (
	"testing"

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
