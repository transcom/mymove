package handlers

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type ConvertersSuite struct {
	*testingsuite.PopTestSuite
}

func TestConvertersSuite(t *testing.T) {
	cs := &ConvertersSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, cs)
	cs.PopTestSuite.TearDown()
}

func (suite *ConvertersSuite) TestFmtUUIDValue() {
	u := uuid.Must(uuid.NewV4())
	expected := strfmt.UUID(u.String())
	result := FmtUUIDValue(u)
	suite.Equal(expected, result, "FmtUUIDValue should match expected value")
}

func (suite *ConvertersSuite) TestFmtUUID() {
	u := uuid.Must(uuid.NewV4())
	result := FmtUUID(u)
	suite.NotNil(result, "FmtUUID result should not be nil")
	suite.Equal(strfmt.UUID(u.String()), *result, "FmtUUID should match expected value")
}

func (suite *ConvertersSuite) TestFmtUUIDPtr() {
	// nil input
	suite.Nil(FmtUUIDPtr(nil), "FmtUUIDPtr should return nil for nil input")
	// non-nil input
	u := uuid.Must(uuid.NewV4())
	result := FmtUUIDPtr(&u)
	suite.NotNil(result, "FmtUUIDPtr result should not be nil for non-nil input")
	suite.Equal(strfmt.UUID(u.String()), *result, "FmtUUIDPtr should match expected value")
}

func (suite *ConvertersSuite) TestFmtDateTime() {
	// zero time returns nil.
	zeroTime := time.Time{}
	suite.Nil(FmtDateTime(zeroTime), "FmtDateTime should return nil for zero time")

	// non-zero time.
	now := time.Now()
	result := FmtDateTime(now)
	suite.NotNil(result, "FmtDateTime should not return nil for non-zero time")
	suite.True(time.Time(*result).Equal(now), "FmtDateTime should match expected time")
}

func (suite *ConvertersSuite) TestFmtDateTimePtr() {
	// nil input
	suite.Nil(FmtDateTimePtr(nil), "FmtDateTimePtr should return nil for nil input")
	// zero time pointer
	zeroTime := time.Time{}
	suite.Nil(FmtDateTimePtr(&zeroTime), "FmtDateTimePtr should return nil for zero time pointer")
	// non-zero time pointer
	now := time.Now()
	result := FmtDateTimePtr(&now)
	suite.NotNil(result, "FmtDateTimePtr should not return nil for non-zero time pointer")
	suite.True(time.Time(*result).Equal(now), "FmtDateTimePtr should match expected time")
}

func (suite *ConvertersSuite) TestFmtDate() {
	// zero date returns nil
	zeroDate := time.Time{}
	suite.Nil(FmtDate(zeroDate), "FmtDate should return nil for zero date")
	// non-zero date
	now := time.Now()
	result := FmtDate(now)
	suite.NotNil(result, "FmtDate should not return nil for non-zero date")
	suite.True(time.Time(*result).Equal(now), "FmtDate should match expected date")
}

func (suite *ConvertersSuite) TestFmtDateSlice() {
	dates := []time.Time{
		time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.February, 2, 0, 0, 0, 0, time.UTC),
	}
	result := FmtDateSlice(dates)
	suite.Len(result, len(dates), "FmtDateSlice should return a slice with the same length as input")
	for i, d := range dates {
		expected := strfmt.Date(d)
		suite.Equal(expected, result[i], "FmtDateSlice element at index %d should match", i)
	}
}

func (suite *ConvertersSuite) TestFmtDatePtr() {
	// nil input
	suite.Nil(FmtDatePtr(nil), "FmtDatePtr should return nil for nil input")
	// zero time pointer
	zeroTime := time.Time{}
	suite.Nil(FmtDatePtr(&zeroTime), "FmtDatePtr should return nil for zero time pointer")
	// non-zero time pointer
	now := time.Now()
	result := FmtDatePtr(&now)
	suite.NotNil(result, "FmtDatePtr should not return nil for non-zero time pointer")
	suite.True(time.Time(*result).Equal(now), "FmtDatePtr should match expected time")
}

func (suite *ConvertersSuite) TestFmtPoundPtr() {
	// nil input
	suite.Nil(FmtPoundPtr(nil), "FmtPoundPtr should return nil for nil input")
	// non-nil input
	weight := unit.Pound(150)
	result := FmtPoundPtr(&weight)
	expected := weight.Int64()
	suite.NotNil(result, "FmtPoundPtr should not return nil for non-nil input")
	suite.Equal(expected, *result, "FmtPoundPtr should match expected value")
}

func (suite *ConvertersSuite) TestFmtURI() {
	uriStr := "https://example.com"
	result := FmtURI(uriStr)
	suite.NotNil(result, "FmtURI should not return nil")
	suite.Equal(uriStr, string(*result), "FmtURI should match expected value")
}

func (suite *ConvertersSuite) TestFmtIntPtrToInt64() {
	// nil input
	suite.Nil(FmtIntPtrToInt64(nil), "FmtIntPtrToInt64 should return nil for nil input")
	// non-nil input
	i := 42
	result := FmtIntPtrToInt64(&i)
	suite.NotNil(result, "FmtIntPtrToInt64 should not return nil for non-nil input")
	suite.Equal(int64(i), *result, "FmtIntPtrToInt64 should match expected value")
}

func (suite *ConvertersSuite) TestFmtInt64() {
	i := int64(42)
	result := FmtInt64(i)
	suite.NotNil(result, "FmtInt64 should not return nil")
	suite.Equal(i, *result, "FmtInt64 should match expected value")
}

func (suite *ConvertersSuite) TestFmtInt() {
	i := 42
	result := FmtInt(i)
	suite.NotNil(result, "FmtInt should not return nil")
	suite.Equal(i, *result, "FmtInt should match expected value")
}

func (suite *ConvertersSuite) TestFmtBool() {
	b := true
	result := FmtBool(b)
	suite.NotNil(result, "FmtBool should not return nil")
	suite.Equal(b, *result, "FmtBool should match expected value")
}

func (suite *ConvertersSuite) TestFmtBoolPtr() {
	// nil input
	suite.Nil(FmtBoolPtr(nil), "FmtBoolPtr should return nil for nil input")
	// non-nil input
	b := true
	result := FmtBoolPtr(&b)
	suite.NotNil(result, "FmtBoolPtr should not return nil for non-nil input")
	suite.Equal(b, *result, "FmtBoolPtr should match expected value")
}

func (suite *ConvertersSuite) TestFmtEmail() {
	email := "test@example.com"
	result := FmtEmail(email)
	suite.NotNil(result, "FmtEmail should not return nil")
	suite.Equal(email, string(*result), "FmtEmail should match expected email")
}

func (suite *ConvertersSuite) TestFmtEmailPtr() {
	// nil input
	suite.Nil(FmtEmailPtr(nil), "FmtEmailPtr should return nil for nil input")
	// non-nil input
	email := "test@example.com"
	result := FmtEmailPtr(&email)
	suite.NotNil(result, "FmtEmailPtr should not return nil for non-nil input")
	suite.Equal(email, string(*result), "FmtEmailPtr should match expected email")
}

func (suite *ConvertersSuite) TestStringFromEmail() {
	// nil input
	suite.Nil(StringFromEmail(nil), "StringFromEmail should return nil for nil input")
	// non-nil input
	email := FmtEmail("test@example.com")
	result := StringFromEmail(email)
	suite.NotNil(result, "StringFromEmail should not return nil for non-nil input")
	suite.Equal(email.String(), *result, "StringFromEmail should match expected value")
}

func (suite *ConvertersSuite) TestFmtString() {
	s := "hello"
	result := FmtString(s)
	suite.NotNil(result, "FmtString should not return nil")
	suite.Equal(s, *result, "FmtString should match expected value")
}

func (suite *ConvertersSuite) TestFmtStringPtr() {
	// nil input
	suite.Nil(FmtStringPtr(nil), "FmtStringPtr should return nil for nil input")
	// non-nil input
	s := "hello"
	result := FmtStringPtr(&s)
	suite.NotNil(result, "FmtStringPtr should not return nil for non-nil input")
	suite.Equal(s, *result, "FmtStringPtr should match expected value")
}

func (suite *ConvertersSuite) TestFmtStringPtrNonEmpty() {
	// nil input
	suite.Nil(FmtStringPtrNonEmpty(nil), "FmtStringPtrNonEmpty should return nil for nil input")
	// empty string (after trimming)
	empty := "    "
	suite.Nil(FmtStringPtrNonEmpty(&empty), "FmtStringPtrNonEmpty should return nil for empty input")
	// non-empty string
	s := "hello"
	result := FmtStringPtrNonEmpty(&s)
	suite.NotNil(result, "FmtStringPtrNonEmpty should not return nil for non-empty input")
	suite.Equal(s, *result, "FmtStringPtrNonEmpty should match expected value")
}

func (suite *ConvertersSuite) TestFmtSSN() {
	ssnStr := "123-45-6789"
	result := FmtSSN(ssnStr)
	suite.NotNil(result, "FmtSSN should not return nil")
	suite.Equal(ssnStr, string(*result), "FmtSSN should match expected SSN")
}

func (suite *ConvertersSuite) TestStringFromSSN() {
	// nil input
	suite.Nil(StringFromSSN(nil), "StringFromSSN should return nil for nil input")
	// non-nil input
	ssn := FmtSSN("123-45-6789")
	result := StringFromSSN(ssn)
	suite.NotNil(result, "StringFromSSN should not return nil for non-nil input")
	suite.Equal(ssn.String(), *result, "StringFromSSN should match expected value")
}

func (suite *ConvertersSuite) TestFmtCost() {
	// nil input
	suite.Nil(FmtCost(nil), "FmtCost should return nil for nil input")
	// non-nil input
	cents := unit.Cents(1000)
	result := FmtCost(&cents)
	expected := cents.Int64()
	suite.NotNil(result, "FmtCost should not return nil for non-nil input")
	suite.Equal(expected, *result, "FmtCost should match expected cost")
}

func (suite *ConvertersSuite) TestFmtMilliCentsPtr() {
	// nil input
	suite.Nil(FmtMilliCentsPtr(nil), "FmtMilliCentsPtr should return nil for nil input")
	// non-nil input
	mc := unit.Millicents(5000)
	result := FmtMilliCentsPtr(&mc)
	expected := mc.Int64()
	suite.NotNil(result, "FmtMilliCentsPtr should not return nil for non-nil input")
	suite.Equal(expected, *result, "FmtMilliCentsPtr should match expected value")
}
