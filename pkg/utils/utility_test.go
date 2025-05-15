package utils_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/utils"
)

type UtilitySuite struct {
	testingsuite.BaseTestSuite
}

func TestUtilitySuite(t *testing.T) {
	suite.Run(t, &UtilitySuite{})
}
func (suite *UtilitySuite) TestStringIsNilEmptyOrWhitespace() {
	suite.Run("nil string", func() {
		actual := utils.IsNullOrWhiteSpace(nil)
		suite.True(actual)
	})

	suite.Run("empty string", func() {
		testString := ""
		actual := utils.IsNullOrWhiteSpace(&testString)
		suite.True(actual)
	})

	suite.Run("whitespace string", func() {
		testString := " "
		actual := utils.IsNullOrWhiteSpace(&testString)
		suite.True(actual)
	})
	suite.Run("valid string", func() {
		testString := "hello"
		actual := utils.IsNullOrWhiteSpace(&testString)
		suite.False(actual)
	})
}

func (suite *UtilitySuite) TestAppendTimestampToFilename() {
	originalFilename := "example.txt"
	result := utils.AppendTimestampToFilename(originalFilename)

	suite.Run("Produces correct formatting", func() {
		expectedPattern := `^example-\d{14}\.txt$`
		matched, err := regexp.MatchString(expectedPattern, result)
		suite.NoError(err, "Error in regex matching")
		suite.True(matched, "Format must match expected pattern")
	})

	suite.Run("Current timestamp", func() {
		parts := regexp.MustCompile(`-(\d{14})\.`).FindStringSubmatch(result)

		suite.Len(parts, 2, "Could not extract timestamp from result")

		timestamp, err := time.Parse(utils.VersionTimeFormat, parts[1])
		suite.NoError(err, "Error parsing timestamp")

		timeWithin2Seconds := time.Since(timestamp) <= 2*time.Second
		suite.True(timeWithin2Seconds, "Timestamp should be now()")
	})

	suite.Run("Preserve original name and extension", func() {
		suite.True(strings.HasPrefix(result, "example-"), "Prefix does not match original filename")
		suite.True(strings.HasSuffix(result, ".txt"), "Suffix does not match original filename extension")
	})

	suite.Run("Handle filename without extension", func() {
		result := utils.AppendTimestampToFilename("noextension")
		expectedPattern := `^noextension-\d{14}$`
		matched, err := regexp.MatchString(expectedPattern, result)
		suite.NoError(err, "Error matching regex")
		suite.True(matched, "Result does not match expected pattern")
	})
}
