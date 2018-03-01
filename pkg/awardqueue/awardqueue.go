package awardqueue

import (
	"fmt"
	"math"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

const numQualBands = 4

// Minimum Performance Score (MPS) is the lowest BVS a TSP can have and still be assigned shipments.
// TODO: This will eventually need to be configurable; implement as something other than a constant.
const mps = 10

type qualityBand models.TransportationServiceProviderPerformances
type qualityBands []qualityBand

// AwardQueue encapsulates the TSP award queue process
type AwardQueue struct {
	db *pop.Connection
}

func (aq *AwardQueue) findAllUnawardedShipments() ([]models.PossiblyAwardedShipment, error) {
	shipments, err := models.FetchAwardedShipments(aq.db)
	return shipments, err
}

// AttemptShipmentAward will attempt to take the given Shipment and award it to
// a TSP.
func (aq *AwardQueue) attemptShipmentAward(shipment models.PossiblyAwardedShipment) (*models.ShipmentAward, error) {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := aq.db.Find(&tdl, shipment.TrafficDistributionListID)

	if err != nil {
		return nil, fmt.Errorf("Cannot find TDL in database: %s", err)
	}

	tspPerformances, err := models.FetchTSPPerformanceForAwardQueue(aq.db, tdl.ID, mps)

	if err != nil {
		return nil, fmt.Errorf("Cannot award. Database error: %s", err)
	}

	if len(tspPerformances) == 0 {
		return nil, fmt.Errorf("Cannot award. No TSPs found in TDL (%v)", tdl.ID)
	}

	var shipmentAward *models.ShipmentAward

	for _, tspPerformance := range tspPerformances {
		tsp := models.TransportationServiceProvider{}
		if err := aq.db.Find(&tsp, tspPerformance.TransportationServiceProviderID); err == nil {
			fmt.Printf("\tAttempting to award to TSP: %s\n", tsp.Name)
			shipmentAward, err = models.CreateShipmentAward(aq.db, shipment.ID, tsp.ID, false)
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

func (aq *AwardQueue) assignUnawardedShipments() {
	fmt.Println("TSP Award Queue running.")

	shipments, err := aq.findAllUnawardedShipments()
	if err == nil {
		count := 0
		for _, shipment := range shipments {
			_, err = aq.attemptShipmentAward(shipment)
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

// getTSPsPerBand determines how many TSPs should be assigned to each Quality Band
// If the number of TSPs in the TDL does not divide evenly into 4 bands, the remainder
// is divided from the top band down.
//
// count is the number of TSPs to distribute.
func getTSPsPerBand(count int) []int {
	bands := make([]int, numQualBands)
	base := int(math.Floor(float64(count) / float64(numQualBands)))
	for i := range bands {
		bands[i] = base
	}

	for i := 0; i < count%numQualBands; i++ {
		bands[i]++
	}
	return bands
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

// NewAwardQueue creates a new AwardQueue
func NewAwardQueue(db *pop.Connection) *AwardQueue {
	return &AwardQueue{db: db}
}

// Run will execute the award queue algorithm.
func Run(db *pop.Connection) {
	queue := NewAwardQueue(db)
	queue.assignUnawardedShipments()
}
