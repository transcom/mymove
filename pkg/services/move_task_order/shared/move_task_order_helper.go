package movetaskordershared

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// GenerateReferenceID generates a reference ID for the MTO
func GenerateReferenceID(db *pop.Connection) (string, error) {
	const maxAttempts = 10
	var referenceID string
	var err error
	for i := 0; i < maxAttempts; i++ {
		referenceID, err = generateReferenceIDHelper(db)
		if err == nil {
			return referenceID, nil
		}
	}
	return "", fmt.Errorf("move_task_order: failed to generate reference id; %w", err)
}

// GenerateReferenceID creates a random ID for an MTO. Format (xxxx-xxxx) with X being a number 0-9 (ex. 0009-1234. 4321-4444)
func generateReferenceIDHelper(db *pop.Connection) (string, error) {
	min := 0
	max := 9999
	firstNum := rand.Intn(max - min + 1)
	secondNum := rand.Intn(max - min + 1)
	newReferenceID := fmt.Sprintf("%04d-%04d", firstNum, secondNum)
	// TODO: change this to use SELECT 1 AS one with a LIMIT of 1 instead of
	// count. Count on a large table is a slow operation. We don't need the
	// actual count here. All we're looking for is whether or not there is
	// already a match, so we can have the DB query return the match if there
	// is one, or none otherwise. Then we can check the size of the result.
	// ActiveRecord has the `.any?` method for this. I'm not sure if Pop has
	// something similar.
	count, err := db.Where(`reference_id= $1`, newReferenceID).Count(&models.Move{})
	if err != nil {
		return "", err
	} else if count > 0 {
		return "", errors.New("move_task_order: reference_id already exists")
	}

	return newReferenceID, nil
}
