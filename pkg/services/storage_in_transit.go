package services

import (
	"time"
)

// StorageInTransitNumberGenerator is an interface for generating a storage in transit number
type StorageInTransitNumberGenerator interface {
	GenerateStorageInTransitNumber(placeInSitTime time.Time) (string, error)
}
