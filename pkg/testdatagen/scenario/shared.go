package scenario

import (
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop"
)

// NamedScenario is a data generation scenario that has a name
type NamedScenario struct {
	Name string
}

// May15_2018 is a date in May 2018
var May15_2018 = time.Date(2018, time.May, 15, 0, 0, 0, 0, time.UTC)

// Oct1_2018 is October 1, 2018
var Oct1_2018 = time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)

// Dec31_2018 is December 31, 2018
var Dec31_2018 = time.Date(2018, time.December, 31, 0, 0, 0, 0, time.UTC)

// May14_2019 is May 14, 2019
var May14_2019 = time.Date(2019, time.May, 14, 0, 0, 0, 0, time.UTC)

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
