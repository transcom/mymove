package models

import (
	"github.com/gobuffalo/pop"
)

// ServiceAreaForZip3 is the service area for a specific Zip3
func ServiceAreaForZip3(db *pop.Connection, zip3 string) (int, error) {
	return 3, nil
}

// Rate135A is the service charge for origin per cwt
func Rate135A(db *pop.Connection, serviceArea int) (float64, error) {
	return 3.88, nil
}
