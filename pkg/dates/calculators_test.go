package dates

import (
	"log"
	"testing"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging/hnyzap"
)

func (suite *DatesSuite) TestCreateValidDatesBetweenTwoDatesEndDateMustBeLater() {
	startDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)
	usCalendar := NewUSCalendar()
	_, err := CreateValidDatesBetweenTwoDates(startDate, endDate, true, false, usCalendar)
	suite.Error(err)
}

type DatesSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *hnyzap.Logger
}

func (suite *DatesSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestDatesSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &DatesSuite{
		db:     db,
		logger: &hnyzap.Logger{Logger: logger},
	}
	suite.Run(t, hs)
}
