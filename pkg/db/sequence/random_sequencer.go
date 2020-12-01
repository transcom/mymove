package sequence

import (
	"errors"

	"github.com/transcom/mymove/pkg/random"
)

// randomSequencer represents a sequencer that produces a random number between min and max, inclusive
type randomSequencer struct {
	min int64
	max int64
}

// NextVal returns the next random value within the range
func (rs randomSequencer) NextVal() (int64, error) {
	randMax := int(rs.max) - int(rs.min) + 1
	randInt, err := random.GetRandomInt(randMax)
	if err != nil {
		return -1, err
	}
	return int64(randInt) + rs.min, nil
}

// SetVal is a no-op for the random sequence generator.
func (rs randomSequencer) SetVal(val int64) error {
	return nil
}

// NewRandomSequencer is a factory for creating a new random number sequencer
func NewRandomSequencer(min int64, max int64) (Sequencer, error) {
	if min < 0 {
		return nil, errors.New("min (%d) cannot be negative")
	}

	if min > max {
		return nil, errors.New("min (%d) cannot be greater than max (%d)")
	}

	return &randomSequencer{min, max}, nil
}
