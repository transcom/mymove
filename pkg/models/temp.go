package models

import (
	"github.com/gobuffalo/pop"
)

func ServiceAreaForZip3(db *pop.Connection, zip3 string) (int, error) {
	return 3, nil
}

func Rate135A(db *pop.Connection, serviceArea int) (float64, error) {
	return 3.88, nil
}
