package scenario

import (
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop"
)

var may15_2018 = time.Date(2018, time.May, 15, 0, 0, 0, 0, time.UTC)
var oct15_2018 = time.Date(2018, time.October, 15, 0, 0, 0, 0, time.UTC)
var may15_2019 = time.Date(2019, time.May, 15, 0, 0, 0, 0, time.UTC)

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
