package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReServiceValidation() {
	suite.Run("test valid ReService", func() {
		validReService := models.ReService{
			Code: "123abc",
			Name: "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReService, expErrors, nil)
	})

	suite.Run("test empty ReService", func() {
		emptyReService := models.ReService{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReService, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchReServiceBycode() {
	suite.Run("success - receive ReService when code is provided", func() {
		reService, err := models.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIHPK)
		suite.NoError(err)
		suite.NotNil(reService)
	})

	suite.Run("failure - receive error when code is not provided", func() {
		var blankReServiceCode models.ReServiceCode
		reService, err := models.FetchReServiceByCode(suite.DB(), blankReServiceCode)
		suite.Error(err)
		suite.Nil(reService)
		suite.Contains(err.Error(), "error fetching from re_services - required code not provided")
	})
}

func (suite *ModelSuite) TestContainsReServiceCode() {
	tests := []struct {
		name      string
		valid     []models.ReServiceCode
		code      models.ReServiceCode
		expectsOK bool
	}{
		{
			name:      "code is in the slice",
			valid:     []models.ReServiceCode{models.ReServiceCodeDOASIT, models.ReServiceCodeDDASIT},
			code:      models.ReServiceCodeDOASIT,
			expectsOK: true,
		},
		{
			name:      "code is not in the slice",
			valid:     []models.ReServiceCode{models.ReServiceCodeDOASIT, models.ReServiceCodeDDASIT},
			code:      models.ReServiceCodeIDFSIT,
			expectsOK: false,
		},
		{
			name:      "empty slice never contains anything",
			valid:     []models.ReServiceCode{},
			code:      models.ReServiceCodeIOFSIT,
			expectsOK: false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ok := models.ContainsReServiceCode(tt.valid, tt.code)
			if tt.expectsOK {
				suite.True(ok, "expected ContainsReServiceCode(%v, %q) to be true", tt.valid, tt.code)
			} else {
				suite.False(ok, "expected ContainsReServiceCode(%v, %q) to be false", tt.valid, tt.code)
			}
		})
	}
}
