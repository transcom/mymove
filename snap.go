package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/transcom/mymove/pkg/models"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a move locator as a command-line argument")
	}

	moveLocator := os.Args[1]
	outputFilename := "moveSnapshot.json"
	if len(os.Args) > 2 {
		outputFilename = os.Args[2]
	}

	move, err := getMoveByLocator(moveLocator)
	if err != nil {
		log.Fatalf("Error retrieving move: %v", err)
	}

	moveSnapshot := map[string]interface{}{
		"move": move,
		// Add other associated components here, e.g., shipments, service items, etc.
	}

	jsonData, err := json.MarshalIndent(moveSnapshot, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile(outputFilename, jsonData, 0600)
	if err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}

	fmt.Printf("Move snapshot generated successfully: %s\n", outputFilename)
}

func getMoveByLocator(moveLocator string) (*models.Move, error) {
	// Open a connection to the database
	db, err := sql.Open("postgres", "postgres://postgres@localhost:5432/dev_db?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}
	defer db.Close()

	// Prepare the SQL query to retrieve the move by locator
	query := "SELECT * FROM moves WHERE locator = $1"
	row := db.QueryRow(query, moveLocator)

	// Create a Move struct to hold the retrieved data
	move := &models.Move{}
	err = row.Scan(
		&move.ID,
		&move.CreatedAt,
		&move.UpdatedAt,
		&move.OrdersID,
		&move.LockExpiresAt,
		&move.AdditionalDocumentsID,
		&move.ApprovedAt,
		&move.Show,
		&move.ContractorID,
		&move.AvailableToPrimeAt,
		&move.SubmittedAt,
		&move.ServiceCounselingCompletedAt,
		&move.ExcessWeightQualifiedAt,
		&move.ExcessWeightUploadID,
		&move.ExcessWeightAcknowledgedAt,
		&move.BillableWeightsReviewedAt,
		&move.FinancialReviewFlag,
		&move.FinancialReviewFlagSetAt,
		&move.PrimeCounselingCompletedAt,
		&move.CloseoutOfficeID,
		&move.ApprovalsRequestedAt,
		&move.ShipmentSeqNum,
		&move.Status,
		&move.Locator,
		&move.CancelReason,
		&move.ReferenceID,
		&move.FinancialReviewRemarks,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no move found with locator '%s'", moveLocator)
		}
		return nil, fmt.Errorf("error scanning move data: %v", err)
	}

	return move, nil
}
