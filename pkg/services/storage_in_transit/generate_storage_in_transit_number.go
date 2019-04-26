package storageintransit

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/services"
)

type generateStorageInTransitNumber struct {
	db *pop.Connection
}

func (s *generateStorageInTransitNumber) generateSequenceNumber(year int, dayOfYear int) (int, error) {
	if year <= 0 {
		return 0, errors.Errorf("Year (%d) must be non-negative", year)
	}

	if dayOfYear <= 0 {
		return 0, errors.Errorf("Day of year (%d) must be non-negative", dayOfYear)
	}

	var sequenceNumber int
	sql := `INSERT INTO storage_in_transit_number_trackers as trackers (year, day_of_year, sequence_number)
			VALUES ($1, $2, 1)
		ON CONFLICT (year, day_of_year)
		DO
			UPDATE
				SET sequence_number = trackers.sequence_number + 1
				WHERE trackers.year = $1 AND trackers.day_of_year = $2
		RETURNING sequence_number
	`

	err := s.db.RawQuery(sql, year, dayOfYear).First(&sequenceNumber)
	if err != nil {
		return 0, errors.Wrapf(err, "Error when incrementing storage in transit sequence number for %d/%d", year, dayOfYear)
	}

	return sequenceNumber, nil
}

// GenerateStorageInTransitNumber creates a new storage in transit number
func (s *generateStorageInTransitNumber) GenerateStorageInTransitNumber(placeInSitTime time.Time) (string, error) {
	utcTime := placeInSitTime.UTC()

	fullYear := utcTime.Year()
	dayOfYear := utcTime.YearDay()
	sequenceNumber, err := s.generateSequenceNumber(fullYear, dayOfYear)
	if err != nil {
		return "", errors.Wrap(err, "Could not generate storage in transit number")
	}

	return fmt.Sprintf("%02d%03d%04d", fullYear%100, dayOfYear, sequenceNumber), nil
}

// NewStorageInTransitNumberGenerator instantiates a new storage in transit number generator implementation
func NewStorageInTransitNumberGenerator(db *pop.Connection) services.StorageInTransitNumberGenerator {
	return &generateStorageInTransitNumber{db}
}
