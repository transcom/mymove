package models_test

import (
	"testing"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/gofrs/uuid"
)

func (suite *ModelSuite) TestEdiErrors() {
	pr := testdatagen.MakeDefaultPaymentRequest(suite.DB())
	prICN := models.PaymentRequestToInterchangeControlNumber{
		ID:                       uuid.Must(uuid.NewV4()),
		PaymentRequestID:         pr.ID,
		InterchangeControlNumber: 5,
	}
	suite.MustCreate(suite.DB(), &prICN)
	testCases := map[string]struct {
		ediError     models.EdiError
		expectedErrs map[string][]string
	}{
		"Successful Create": {
			ediError: models.EdiError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    models.EDI824,
				PaymentRequestID:           pr.ID,
				InterchangeControlNumberID: prICN.ID,
				Code:                       swag.String("B"),
				Description:                swag.String("EDI Error happened to field 99"),
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			ediError: models.EdiError{},
			expectedErrs: map[string][]string{
				/*
					"description":                   {"Code or Description must be present"},
					"code":                          {"Code or Description must be present"},
				*/
				"payment_request_id":            {"PaymentRequestID can not be blank."},
				"interchange_control_number_id": {"InterchangeControlNumberID can not be blank."},
				"editype":                       {"EDIType is not in the list [810, 824, 858, 997]."},
			},
		},
		"Message Type Invalid": {
			ediError: models.EdiError{
				ID:                         uuid.Must(uuid.NewV4()),
				EDIType:                    "EDI956",
				PaymentRequestID:           pr.ID,
				InterchangeControlNumberID: prICN.ID,
				Code:                       swag.String("C"),
				Description:                swag.String("EDI Error happened to field 123"),
			},
			expectedErrs: map[string][]string{
				"editype": {"EDIType is not in the list [810, 824, 858, 997]."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.ediError, test.expectedErrs)
		})
	}
}
