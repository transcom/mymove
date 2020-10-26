package edisegment

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SegmentSuite struct {
	testingsuite.BaseTestSuite
	validator *validator.Validate
}

func TestSegmentSuite(t *testing.T) {
	ss := &SegmentSuite{
		validator: validator.New(),
	}

	// Custom validators
	err := ss.validator.RegisterValidation("timeformat", HasTimeFormat)
	if err != nil {
		log.Fatal(err)
	}

	suite.Run(t, ss)
}

func (suite *SegmentSuite) ValidateError(err error, structField, expectedTag string) {
	errs := err.(validator.ValidationErrors)

	found := false
	var fe validator.FieldError

	for i := 0; i < len(errs); i++ {
		if errs[i].StructField() == structField {
			found = true
			fe = errs[i]
			break
		}
	}

	suite.True(found)
	suite.NotNil(fe)
	if fe != nil {
		suite.Equal(expectedTag, fe.Tag())
	}
}

func (suite *SegmentSuite) ValidateErrorLen(err error, length int) {
	errs := err.(validator.ValidationErrors)
	suite.Len(errs, length)
}
