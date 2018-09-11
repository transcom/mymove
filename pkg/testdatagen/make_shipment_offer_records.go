package testdatagen

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

// MakeShipmentOffer creates a single shipment offer record
func MakeShipmentOffer(db *pop.Connection, assertions Assertions) models.ShipmentOffer {

	// Test for Shipment first before creating a new Shipment
	shipment := assertions.ShipmentOffer.Shipment
	if isZeroUUID(assertions.ShipmentOffer.ShipmentID) {
		shipment = MakeShipment(db, assertions)
	}

	// Test for TSP ID first before creating a new TSP

	tsp := assertions.ShipmentOffer.TransportationServiceProvider
	if isZeroUUID(tsp.ID) || isZeroUUID(assertions.ShipmentOffer.TransportationServiceProviderID) {
		tsp = MakeTSP(db, assertions)
	}
	shipmentOffer := models.ShipmentOffer{
		ShipmentID:                      shipment.ID,
		Shipment:                        shipment,
		TransportationServiceProviderID: tsp.ID,
		TransportationServiceProvider:   tsp,
		AdministrativeShipment:          false,
		Accepted:                        nil, // This is a Tri-state and new offers are always nil until accepted
		RejectionReason:                 nil,
	}

	mergeModels(&shipmentOffer, assertions.ShipmentOffer)

	mustCreate(db, &shipmentOffer)

	return shipmentOffer
}

// MakeDefaultShipmentOffer makes a ShipmentOffer with default values
func MakeDefaultShipmentOffer(db *pop.Connection) models.ShipmentOffer {
	return MakeShipmentOffer(db, Assertions{})
}

// MakeShipmentOfferData creates one offered shipment record
func MakeShipmentOfferData(db *pop.Connection) {
	// Get a shipment ID
	shipmentList := []models.Shipment{}
	err := db.All(&shipmentList)
	if err != nil {
		fmt.Println("Shipment ID import failed.")
	}

	// Get a TSP ID
	tspList := []models.TransportationServiceProvider{}
	err = db.All(&tspList)
	if err != nil {
		fmt.Println("TSP ID import failed.")
	}

	// Add one offered shipment record for each shipment and a random TSP IDs
	for _, shipment := range shipmentList {
		shipmentOfferAssertions := Assertions{
			ShipmentOffer: models.ShipmentOffer{
				ShipmentID:                      shipment.ID,
				TransportationServiceProviderID: tspList[rand.Intn(len(tspList))].ID,
				AdministrativeShipment:          false,
				Accepted:                        nil, // See note about Tri-state above
				RejectionReason:                 nil,
			},
		}
		MakeShipmentOffer(db, shipmentOfferAssertions)
	}
}

// CreateShipmentOfferData creates a list of TSP Users, Shipments, and Shipment Offers
// Must pass in the number of tsp users to create and number of shipments.
// The split of shipment offers should be the length of TSP users and the sum should equal the number of shipments
func CreateShipmentOfferData(db *pop.Connection, numTspUsers int, numShipments int, numShipmentOfferSplit []int, statuses []models.ShipmentStatus) ([]models.TspUser, []models.Shipment, []models.ShipmentOffer, error) {
	var tspUserList []models.TspUser
	var shipmentList []models.Shipment
	var shipmentOfferList []models.ShipmentOffer

	// Error check some inputs
	if len(numShipmentOfferSplit) != numTspUsers {
		err := errors.New("Length of numShipmentOfferSplit should equal numTspUsers")
		return tspUserList, shipmentList, shipmentOfferList, err
	}

	soSplitSum := 0
	for _, val := range numShipmentOfferSplit {
		soSplitSum += val
	}
	if soSplitSum != numShipments {
		err := errors.New("Number of shipment offers in split should equal numShipments")
		return tspUserList, shipmentList, shipmentOfferList, err
	}

	// Create TSP Users
	for i := 1; i <= numTspUsers; i++ {
		email := fmt.Sprintf("leo_spaceman_tsp_%d@example.com", i)
		tspUserAssertions := Assertions{
			User: models.User{
				LoginGovEmail: email,
			},
			TspUser: models.TspUser{
				Email: email,
			},
		}
		tspUser := MakeTspUser(db, tspUserAssertions)
		tspUserList = append(tspUserList, tspUser)
	}

	// Create shipments
	tdl := MakeTDL(
		db, Assertions{
			TrafficDistributionList: models.TrafficDistributionList{
				SourceRateArea:    DefaultSrcRateArea,
				DestinationRegion: DefaultDstRegion,
				CodeOfService:     DefaultCOS,
			},
		})
	market := "dHHG"
	sourceGBLOC := "OHAI"
	oneWeek, _ := time.ParseDuration("7d")
	selectedMoveType := "HHG"
	if len(statuses) == 0 {
		// Statuses for shipments attached to a shipment offer should not be DRAFT or SUBMITTED
		// because this should be after the award queue has run and SUBMITTED shipments have been awarded
		statuses = []models.ShipmentStatus{
			models.ShipmentStatusAWARDED,
			models.ShipmentStatusACCEPTED,
			models.ShipmentStatusAPPROVED,
			models.ShipmentStatusINTRANSIT,
			models.ShipmentStatusDELIVERED}
	}
	for i := 1; i <= numShipments; i++ {
		now := time.Now()
		nowPlusOne := now.Add(oneWeek)
		nowPlusTwo := now.Add(oneWeek * 2)

		// Service Member Details
		smEmail := fmt.Sprintf("leo_spaceman_sm_%d@example.com", i)

		// Shipment Details
		shipmentStatus := statuses[rand.Intn(len(statuses))]

		// Move Details
		moveStatus := models.MoveStatusSUBMITTED
		if shipmentStatus == models.ShipmentStatusAPPROVED {
			moveStatus = models.MoveStatusAPPROVED
		}

		shipmentAssertions := Assertions{
			User: models.User{
				LoginGovEmail: smEmail,
			},
			Move: models.Move{
				SelectedMoveType: &selectedMoveType,
				Status:           moveStatus,
			},
			Shipment: models.Shipment{
				RequestedPickupDate:     &now,
				ActualPickupDate:        &nowPlusOne,
				DeliveryDate:            &nowPlusTwo,
				TrafficDistributionList: &tdl,
				SourceGBLOC:             &sourceGBLOC,
				Market:                  &market,
				Status:                  shipmentStatus,
			},
		}
		shipment := MakeShipment(db, shipmentAssertions)
		shipmentList = append(shipmentList, shipment)

		// Accepted shipments must have an OSA and DSA
		if shipmentStatus == models.ShipmentStatusACCEPTED || shipmentStatus == models.ShipmentStatusAWARDED {
			originServiceAgentAssertions := Assertions{
				ServiceAgent: models.ServiceAgent{
					ShipmentID: shipment.ID,
					Role:       models.RoleORIGIN,
				},
			}
			MakeServiceAgent(db, originServiceAgentAssertions)
			destinationServiceAgentAssertions := Assertions{
				ServiceAgent: models.ServiceAgent{
					ShipmentID: shipment.ID,
					Role:       models.RoleDESTINATION,
				},
			}
			MakeServiceAgent(db, destinationServiceAgentAssertions)
		}

		// Approved shipments need to collect Weight Ticket documents
		if shipmentStatus == models.ShipmentStatusAPPROVED {
			docAssertions := Assertions{
				Document: models.Document{
					ServiceMemberID: shipment.Move.Orders.ServiceMember.ID,
					ServiceMember:   shipment.Move.Orders.ServiceMember,
				},
				MoveDocument: models.MoveDocument{
					MoveID:           shipment.Move.ID,
					Move:             shipment.Move,
					ShipmentID:       &shipment.ID,
					Shipment:         shipment,
					MoveDocumentType: models.MoveDocumentTypeWEIGHTTICKET,
					Title:            fmt.Sprintf("move_document_%d", i),
				},
			}
			MakeMoveDocument(db, docAssertions)
		}
	}

	// A Shipment Offer is created for each Shipment and split among TSPs
	count := 0
	for index, split := range numShipmentOfferSplit {
		tspUser := tspUserList[index]
		subShipmentList := shipmentList[count : count+split]
		count += split
		for _, shipment := range subShipmentList {
			shipmentOfferAssertions := Assertions{
				ShipmentOffer: models.ShipmentOffer{
					ShipmentID:                      shipment.ID,
					TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
				},
			}
			shipmentOffer := MakeShipmentOffer(db, shipmentOfferAssertions)
			shipmentOfferList = append(shipmentOfferList, shipmentOffer)
		}
	}

	return tspUserList, shipmentList, shipmentOfferList, nil
}
