package primeapi

import (
	"log"
	"testing"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseHandlerTestSuite
}

// SetupTest sets up the test suite by preparing the DB
func (suite *HandlerSuite) SetupTest() {
	suite.DB().TruncateAll()
}

// AfterTest completes tests by trying to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.TestFilesToClose() {
		file.Data.Close()
	}
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HandlerSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(logger, notifications.NewStubNotificationSender("milmovelocal", logger), testingsuite.CurrentPackage()),
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// MakeAvailableMoveTaskOrder returns a default MTO with AvailableToPrimeAt set. Shared throughout this test package.
func MakeAvailableMoveTaskOrder(db *pop.Connection) models.MoveTaskOrder {
	now := time.Now()
	mto := testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			AvailableToPrimeAt: &now,
		},
	})
	return mto
}
