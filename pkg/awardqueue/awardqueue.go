package awardqueue

import (
	"fmt"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

var dbConnection *pop.Connection
var numQualBands = 4

type qualityBand []string
type qualityBands []qualityBand

func findAllUnawardedShipments() ([]models.PossiblyAwardedShipment, error) {
	shipments, err := models.FetchAwardedShipments(dbConnection)
	return shipments, err
}

func selectTSPToAwardShipment(shipment models.PossiblyAwardedShipment) error {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := dbConnection.Find(&tdl, shipment.TrafficDistributionListID)

	// Find TSPs in that TDL sorted by shipment_awards[asc] and bvs[desc]
	// tspssba stands for TSPs sorted by award
	tspsba, err := models.FetchTSPsInTDLSortByAward(dbConnection, tdl.ID)

	for _, consideredTSP := range tspsba {
		fmt.Printf("\tConsidering TSP: %v\n", consideredTSP)

		tsp := models.TransportationServiceProvider{}
		err := dbConnection.Find(&tsp, consideredTSP.TransportationServiceProviderID)
		if err == nil {
			// We found a valid TSP to award to!
			err := models.CreateShipmentAward(dbConnection, shipment.ID, tsp.ID, false)
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

	return err
}

// getTSPsPerBand detemines how many TSPs should be assigned to each Quality Band
// If the number of TSPs in the TDL does not divide evenly into 4 bands, the remainder
// is divided from the top band down. Function takes length of TSPs array (tsp count) as arg.
func getTSPsPerBand(tspc int) []int {
	tspPerBandList := make([]int, numQualBands)
	// determine floor of tsps per band
	// tspp is TSP per band
	tsppb := int(math.Floor(float64(tspc) / float64(numQualBands)))
	// assign tsppb to each band in tspPerBandList
	for i := range tspPerBandList {
		tspPerBandList[i] = tsppb
	}
	// add one for bands that will have extra tsps
	for i := 0; i < tspc%numQualBands; i++ {
		tspPerBandList[i] = tsppb + 1
	}
	return tspPerBandList
}

func assignTSPsToBands(tsps []string) qualityBands {
	tspIndex := 0
	qbs := make(qualityBands, 4, 4)
	// Determine how many TSPs should be in each band
	tsppbl := getTSPsPerBand(len(tsps))
	// take the number of tsps that should be assigned each band
	for i, tsppb := range tsppbl {
		// add that number of tsps to the relevant quality band
		for j := tspIndex; j < tspIndex+tsppb; j++ {
			qbs[i] = append(qbs[i], tsps[j])
		}
		tspIndex += tsppb
	}
	return qbs
}

func assignQualityBands() {
	// Find TSPs in that TDL sorted bvs[desc]
	fmt.Printf("Assigning TSPs quality bands")
	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	// tspsbb stands for TSPs sorted by BVS
	tspsbb, err := models.FetchTSPsInTDLSortByBVS(dbConnection, tdl.ID)
	// Assign TSPs to bands and return slice of TSP slices divided by band
	return assignTSPsToBands(tspsbb)
}

// Run will execute the award queue algorithm.
func Run(db *pop.Connection) {
	dbConnection = db

	fmt.Println("TSP Award Queue running.")

	shipments, err := findAllUnawardedShipments()
	if err == nil {
		count := -1
		for i, shipment := range shipments {
			err = selectTSPToAwardShipment(shipment)
			if err != nil {
				fmt.Printf("Failed to award shipment: %s\n", err)
			}
			count = i
		}
		fmt.Printf("Awarded %d shipments.\n", count+1)
	} else {
		fmt.Printf("Failed to query for shipments: %s", err)
	}
}
