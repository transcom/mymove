package testingsuite

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
)

// PopTestSuite is a suite for testing
type PopTestSuite struct {
	BaseTestSuite
	db *pop.Connection
}

// NewPopTestSuite returns a new PopTestSuite
func NewPopTestSuite() PopTestSuite {
	// Find root path of testingsuite package
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	configLocation := filepath.Join(basepath, "../../config")
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

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

// MustCreate requires creating without errors
func (suite *PopTestSuite) MustCreate(db *pop.Connection, model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := db.ValidateAndCreate(model)
	if err != nil {
		suite.T().Errorf("Errors encountered creating %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered creating %v: %v", model, verrs)
	}
}

// NoVerrs prints any errors it receives
func (suite *PopTestSuite) NoVerrs(verrs *validate.Errors) bool {
	if !suite.False(verrs.HasAny()) {
		fmt.Println(verrs.String())
		return false
	}
	return true
}
