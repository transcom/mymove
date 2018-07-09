package scenario

import (
	"fmt"
	"log"
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

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered saving %v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("Validation errors encountered saving %v: %v", model, verrs))
	}
}
