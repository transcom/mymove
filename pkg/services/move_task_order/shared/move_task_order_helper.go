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
	return "", errors.New("move_task_order: failed to generate reference id")
}

// GenerateReferenceID creates a random ID for an MTO. Format (xxxx-xxxx) with X being a number 0-9 (ex. 0009-1234. 4321-4444)
func generateReferenceIDHelper(db *pop.Connection) (string, error) {
	min := 0
	max := 9999
	firstNum := rand.Intn(max - min + 1)
	secondNum := rand.Intn(max - min + 1)
	newReferenceID := fmt.Sprintf("%04d-%04d", firstNum, secondNum)
	count, err := db.Where(`reference_id= $1`, newReferenceID).Count(&models.MoveTaskOrder{})
	if err != nil {
		return "", err
	} else if count > 0 {
		return "", errors.New("move_task_order: reference_id already exists")
	}

	return newReferenceID, nil
}
