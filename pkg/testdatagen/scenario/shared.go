package scenario

import (
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/testdatagen"
)

// NamedScenario is a data generation scenario that has a name
type NamedScenario struct {
	Name string
}

// May15TestYear is a May 15 of TestYear
var May15TestYear = time.Date(testdatagen.TestYear, time.May, 15, 0, 0, 0, 0, time.UTC)

// Oct1TestYear is October 1 of TestYear
var Oct1TestYear = time.Date(testdatagen.TestYear, time.October, 1, 0, 0, 0, 0, time.UTC)

// Dec31TestYear is December 31 of TestYear
var Dec31TestYear = time.Date(testdatagen.TestYear, time.December, 31, 0, 0, 0, 0, time.UTC)

// May14FollowingYear is May 14 of the year AFTER TestYear
var May14FollowingYear = time.Date(testdatagen.TestYear+1, time.May, 14, 0, 0, 0, 0, time.UTC)

func save(db *pop.Connection, model interface{}) error {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		return errors.Wrap(err, "Errors encountered saving model")
	}
	if verrs.HasAny() {
		return errors.Errorf("Validation errors encountered saving model: %v", verrs)
	}
	return nil
}
