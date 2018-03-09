package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_SignedCertificationValidations() {
	signedCert := &SignedCertification{}

	expErrors := map[string][]string{
		"certification_text": {"CertificationText can not be blank."},
		"signature":          {"Signature can not be blank."},
	}

	suite.verifyValidationErrors(signedCert, expErrors)
}
