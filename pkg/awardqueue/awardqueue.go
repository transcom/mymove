package awardqueue

import (
	"fmt"
	"math"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

var db *pop.Connection

const numQualBands = 4

// Minimum Performance Score (MPS) is the lowest BVS a TSP can have and still be assigned shipments.
// TODO: This will eventually need to be configurable; implement as something other than a constant.
const mps = 10

type qualityBand models.TransportationServiceProviderPerformances
type qualityBands []qualityBand

func findAllUnawardedShipments() ([]models.PossiblyAwardedShipment, error) {
	shipments, err := models.FetchAwardedShipments(db)
	return shipments, err
}

// AttemptShipmentAward will attempt to take the given Shipment and award it to
// a TSP.
func AttemptShipmentAward(shipment models.PossiblyAwardedShipment) (*models.ShipmentAward, error) {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := db.Find(&tdl, shipment.TrafficDistributionListID)

	tspPerformances, err := models.FetchTSPPerformanceForAwardQueue(db, tdl.ID, mps)

	if err != nil {
		return nil, fmt.Errorf("Cannot award. Database error: %s", err)
	}

	if len(tspPerformances) == 0 {
		return nil, fmt.Errorf("Cannot award. No TSPs found in TDL (%v)", tdl.ID)
	}

	var shipmentAward *models.ShipmentAward

	for _, tspPerformance := range tspPerformances {
		tsp := models.TransportationServiceProvider{}
		if err := db.Find(&tsp, tspPerformance.TransportationServiceProviderID); err == nil {
			fmt.Printf("\tAttempting to award to TSP: %s\n", tsp.Name)
			shipmentAward, err = models.CreateShipmentAward(db, shipment.ID, tsp.ID, false)
			if err == nil {
				fmt.Print("\tShipment awarded to TSP!\n")
				break
			} else {
				fmt.Printf("\tFailed to award to TSP: %v\n", err)
			}
		} else {
			fmt.Printf("\tFailed to award to TSP: %v\n", err)
		}
	}

	return shipmentAward, err
}

// getTSPsPerBand detemines how many TSPs should be assigned to each Quality Band
// If the number of TSPs in the TDL does not divide evenly into 4 bands, the remainder
// is divided from the top band down. Function takes length of TSPs array as arg.
func getTSPsPerBand(tspc int) []int {
	// tsppb is TSP per band
	tsppbList := make([]int, numQualBands)
	tsppb := int(math.Floor(float64(tspc) / float64(numQualBands)))
	for i := range tsppbList {
		tsppbList[i] = tsppb
	}

	for i := 0; i < tspc%numQualBands; i++ {
		tsppbList[i]++
	}
	return tsppbList
}

// assignTSPsToBands takes slice of tsps and returns
// slice of slices in which they're sorted into 4 bands
func assignTSPsToBands(tspPerfs models.TransportationServiceProviderPerformances) qualityBands {
	tspIndex := 0
	qbs := make(qualityBands, numQualBands)
	tsppbList := getTSPsPerBand(len(tspPerfs))

	for i, tsppb := range tsppbList {
		for j := tspIndex; j < tspIndex+tsppb; j++ {
			qbs[i] = append(qbs[i], tspPerfs[j])
		}
		tspIndex += tsppb
	}
	return qbs
}

// Assign TSPs to bands and return struct slice of band slices
func assignQualityBands() (qualityBands, error) {
	fmt.Printf("Assigning TSPs quality bands")
	tdl := models.TrafficDistributionList{}
	tspPerfs, err := models.FetchTSPPerformanceForQualityBandAssignment(db, tdl.ID, mps)
	return assignTSPsToBands(tspPerfs), err
}

// Run will execute the award queue algorithm.
func Run(dbConnection *pop.Connection) {
	db = dbConnection

	fmt.Println("TSP Award Queue running.")

	shipments, err := findAllUnawardedShipments()
	if err == nil {
		count := 0
		for _, shipment := range shipments {
			_, err = AttemptShipmentAward(shipment)
			if err != nil {
				fmt.Printf("\tFailed to award shipment: %s\n", err)
			} else {
				count++
			}
		}
		fmt.Printf("Awarded %d shipments.\n", count)
	} else {
		fmt.Printf("Failed to query for shipments: %s", err)
	}
}
