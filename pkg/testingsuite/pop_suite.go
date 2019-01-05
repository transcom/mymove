package testingsuite

import (
	"github.com/gobuffalo/pop"
)

// PopTestSuite is a suite for testing
type PopTestSuite struct {
	BaseTestSuite
	db *pop.Connection
}

// NewPopTestSuite returns a new PopTestSuite
func NewPopTestSuite(db *pop.Connection) PopTestSuite {
	return PopTestSuite{db: db}
}

// DB returns a db connection
func (suite *PopTestSuite) DB() *pop.Connection {
	return suite.db
}

// MustSave requires saving without errors
func (suite *PopTestSuite) MustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}
