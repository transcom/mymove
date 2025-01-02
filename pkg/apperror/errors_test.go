package apperror

import (
	"errors"
	"fmt"
	"testing"

	validate "github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
)

type errorsSuite struct {
	suite.Suite
}

func TestErrorSuite(t *testing.T) {
	hs := &errorsSuite{}
	suite.Run(t, hs)
}

func (suite *errorsSuite) TestContextError() {
	suite.Run("ContextError shows error message", func() {
		contextError := NewContextError("This is a context error message")
		suite.Equal("ContextError: This is a context error message", contextError.Error())

	})
}

func (suite *errorsSuite) TestPreconditionFailedError() {
	suite.Run("PreconditionFailedError error function message", func() {
		id := uuid.Must(uuid.NewV4())
		err := errors.New("Precondition Failed Error")
		preconditionFailedError := NewPreconditionFailedError(id, err)
		suite.Equal("Precondition failed on update to object with ID: '"+id.String()+"'. The If-Match header value did not match the eTag for this record.", preconditionFailedError.Error())

	})
}

func (suite *errorsSuite) TestNotFoundError() {
	suite.Run("NotFoundError error function message", func() {
		id := uuid.Must(uuid.NewV4())
		errorMessage := "This is a Not Found Error"
		notFoundError := NewNotFoundError(id, errorMessage)
		suite.Equal("ID: "+id.String()+" not found "+errorMessage, notFoundError.Error())

	})
}

func (suite *errorsSuite) TestUpdateError() {
	suite.Run("UpdateError error function message", func() {
		id := uuid.Must(uuid.NewV4())
		errorMessage := "This is an Update Error"
		updateError := NewUpdateError(id, errorMessage)
		suite.Equal("Update Error "+errorMessage, updateError.Error())

	})
}

func (suite *errorsSuite) TestPPMNotReadyForCloseoutError() {
	suite.Run("PPMNotReadyForCloseoutError error function message", func() {
		id := uuid.Must(uuid.NewV4())
		errorMessage := "This is a PPM Not Ready For Closeout Error"
		ppmNotReadyForCloseoutError := NewPPMNotReadyForCloseoutError(id, errorMessage)
		suite.Equal("ID: "+id.String()+" - PPM Shipment is not ready for closeout. Customer must upload PPM documents. "+errorMessage, ppmNotReadyForCloseoutError.Error())

	})
}

func (suite *errorsSuite) TestPPMNoWeightTicketsError() {
	suite.Run("PPMNoWeightTicketsError error function message", func() {
		id := uuid.Must(uuid.NewV4())
		errorMessage := "This is a PPM No Weight Tickets Error"
		pPMNoWeightTicketsError := NewPPMNoWeightTicketsError(id, errorMessage)
		suite.Equal("ID: "+id.String()+" - PPM Shipment has no weight tickets assigned to it, can't calculate any weights. "+errorMessage, pPMNoWeightTicketsError.Error())

	})
}

func (suite *errorsSuite) TestBadDataError() {
	suite.Run("BadDataError error function message", func() {
		errorMessage := "This is a Bad Data Error"
		badDataError := NewBadDataError(errorMessage)
		suite.Equal(fmt.Sprintf("Data received from requester is bad: %s: %s", badDataError.baseError.code, errorMessage), badDataError.Error())
	})
}

func (suite *errorsSuite) TestUnsupportedPostalCodeError() {
	suite.Run("UnsupportedPostalCodeError error function message", func() {
		postalCode := "36022"
		reason := "This postal code is not supported"
		unsupportedPostalCodeError := NewUnsupportedPostalCodeError(postalCode, reason)
		suite.Equal(fmt.Sprintf("Unsupported postal code (%s): %s", postalCode, reason), unsupportedPostalCodeError.Error())
	})
}

func (suite *errorsSuite) TestUnsupportedPortCodeError() {
	suite.Run("UnsupportedPortCodeError error function message", func() {
		portCode := "ABC"
		reason := "This port code is not legit"
		unsupportedPortCode := NewUnsupportedPortCodeError(portCode, reason)
		suite.Equal(fmt.Sprintf("Unsupported port code (%s): %s", portCode, reason), unsupportedPortCode.Error())
	})
}

func (suite *errorsSuite) TestInvalidInputError() {
	suite.Run("InvalidInputError error function returns the message when it's not empty", func() {
		id := uuid.Must(uuid.NewV4())
		err := errors.New("Invalid Input Error")
		verrs := validate.NewErrors()
		fieldWithErr := "field"
		fieldErrorMsg := "Field error"
		verrs.Add(fieldWithErr, fieldErrorMsg)
		message := "Error messagefor an invalid input"
		invalidInputError := NewInvalidInputError(id, err, verrs, message)
		suite.Equal(message, invalidInputError.Error())
	})

	suite.Run("InvalidInputError error function with no message but validation errors returns correct message", func() {
		id := uuid.Must(uuid.NewV4())
		err := errors.New("Invalid Input Error")
		verrs := validate.NewErrors()
		fieldWithErr := "field"
		fieldErrorMsg := "Field error"
		verrs.Add(fieldWithErr, fieldErrorMsg)
		message := ""
		invalidInputError := NewInvalidInputError(id, err, verrs, message)
		suite.Equal(fmt.Sprintf("Invalid input for ID: %s. %s", id.String(), verrs), invalidInputError.Error())
	})
}
