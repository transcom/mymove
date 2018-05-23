package scenario

import (
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop"
)

// May15_2018 is a date in May 2018
var May15_2018 = time.Date(2018, time.May, 15, 0, 0, 0, 0, time.UTC)

// Oct15_2018 is a date in October 2018
var Oct15_2018 = time.Date(2018, time.October, 15, 0, 0, 0, 0, time.UTC)

// May15_2019 is a date in May 2019
var May15_2019 = time.Date(2019, time.May, 15, 0, 0, 0, 0, time.UTC)

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
