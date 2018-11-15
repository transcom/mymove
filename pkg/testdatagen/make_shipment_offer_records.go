package testdatagen

import (
	"fmt"
	"github.com/transcom/mymove/pkg/unit"
	"math/rand"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
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

	tspp := assertions.ShipmentOffer.TransportationServiceProviderPerformance
	if isZeroUUID(tspp.ID) || isZeroUUID(assertions.ShipmentOffer.TransportationServiceProviderPerformanceID) {
		tspp = MakeTSPPerformance(db, assertions)
	}

	shipmentOffer := models.ShipmentOffer{
		ShipmentID:                                 shipment.ID,
		Shipment:                                   shipment,
		TransportationServiceProviderID:            tsp.ID,
		TransportationServiceProvider:              tsp,
		TransportationServiceProviderPerformance:   tspp,
		TransportationServiceProviderPerformanceID: tspp.ID,
		AdministrativeShipment:                     false,
		Accepted:                                   nil, // This is a Tri-state and new offers are always nil until accepted
		RejectionReason:                            nil,
	}

	mergeModels(&shipmentOffer, assertions.ShipmentOffer)

	mustCreate(db, &shipmentOffer)

	return shipmentOffer
}

// MakeDefaultShipmentOffer makes a ShipmentOffer with default values
func MakeDefaultShipmentOffer(db *pop.Connection) models.ShipmentOffer {
	return MakeShipmentOffer(db, Assertions{})
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
	sourceGBLOC := "KKFA"
	destinationGBLOC := "HAFC"
	selectedMoveType := models.SelectedMoveTypeHHG
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

	// Make the required Tariff 400 NG Zip3
	MakeDefaultTariff400ngZip3(db)
	MakeTariff400ngZip3(db, Assertions{
		Tariff400ngZip3: models.Tariff400ngZip3{
			Zip3:          "800",
			BasepointCity: "Denver",
			State:         "CO",
			ServiceArea:   "145",
			RateArea:      "US74",
			Region:        "5",
		},
	})

	shouldCreateTariffData := false
	for i := 1; i <= numShipments; i++ {
		// Service Member Details
		smEmail := fmt.Sprintf("leo_spaceman_sm_%d@example.com", i)

		// Shipment Details
		shipmentStatus := statuses[rand.Intn(len(statuses))]

		// New Duty Station
		// Check if Buckley Duty Station exists, if not, create
		newDutyStation, err := models.FetchDutyStationByName(db, "Buckley AFB")
		if err != nil {
			newDutyStationAssertions := Assertions{
				Address: models.Address{
					City:       "Aurora",
					State:      "CO",
					PostalCode: "80011",
				},
				DutyStation: models.DutyStation{
					Name: "Buckley AFB",
				},
			}
			newDutyStation = MakeDutyStation(db, newDutyStationAssertions)
		}

		// Move and Order Details
		moveStatus := models.MoveStatusSUBMITTED
		orderStatus := models.OrderStatusSUBMITTED
		ordTypeDetHHGPermit := internalmessages.OrdersTypeDetailHHGPERMITTED
		if shipmentStatus == models.ShipmentStatusAPPROVED ||
			shipmentStatus == models.ShipmentStatusINTRANSIT ||
			shipmentStatus == models.ShipmentStatusDELIVERED ||
			shipmentStatus == models.ShipmentStatusCOMPLETED {
			moveStatus = models.MoveStatusAPPROVED
			orderStatus = models.OrderStatusAPPROVED
		}

		shipmentAssertions := Assertions{
			User: models.User{
				LoginGovEmail: smEmail,
			},
			Order: models.Order{
				NewDutyStationID: newDutyStation.ID,
				NewDutyStation:   newDutyStation,
				Status:           orderStatus,
				OrdersTypeDetail: &ordTypeDetHHGPermit,
			},
			Move: models.Move{
				SelectedMoveType: &selectedMoveType,
				Status:           moveStatus,
			},
			Shipment: models.Shipment{
				TrafficDistributionList: &tdl,
				SourceGBLOC:             &sourceGBLOC,
				DestinationGBLOC:        &destinationGBLOC,
				Market:                  &market,
				Status:                  shipmentStatus,
				// Let the next method fill in the dates
			},
		}
		shipment := MakeShipment(db, shipmentAssertions)

		durIndex := time.Duration(i + 1)

		// Set dates based on status
		if shipmentStatus == models.ShipmentStatusINTRANSIT || shipmentStatus == models.ShipmentStatusDELIVERED {
			shipment.PmSurveyConductedDate = &Now
			shipment.PmSurveyPlannedPackDate = &NowPlusOneWeek
			shipment.PmSurveyPlannedPickupDate = &NowPlusOneWeek
			shipment.PmSurveyPlannedDeliveryDate = &NowPlusTwoWeeks
			shipment.ActualPackDate = &Now
			// For sortability, we need varying pickup dates
			pickupDate := Now.Add(OneDay * durIndex)
			shipment.ActualPickupDate = &pickupDate

			shipment.NetWeight = shipment.WeightEstimate

			shouldCreateTariffData = true
		}

		if shipmentStatus == models.ShipmentStatusDELIVERED {
			// For sortability, we need varying delivery dates
			deliveryDate := Now.Add(OneWeek * durIndex)
			shipment.ActualDeliveryDate = &deliveryDate
		}

		// Assign a new unique GBL number using source GBLOC
		shipment.AssignGBLNumber(db)
		mustSave(db, &shipment)

		shipmentList = append(shipmentList, shipment)

		// Accepted shipments must have an OSA and DSA
		if shipmentStatus == models.ShipmentStatusACCEPTED || shipmentStatus == models.ShipmentStatusAWARDED {
			originServiceAgentAssertions := Assertions{
				ServiceAgent: models.ServiceAgent{
					ShipmentID: shipment.ID,
					Shipment:   &shipment,
					Role:       models.RoleORIGIN,
				},
			}
			MakeServiceAgent(db, originServiceAgentAssertions)
			destinationServiceAgentAssertions := Assertions{
				ServiceAgent: models.ServiceAgent{
					ShipmentID: shipment.ID,
					Shipment:   &shipment,
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
			var offerState *bool
			if shipment.Status == models.ShipmentStatusACCEPTED || shipment.Status == models.ShipmentStatusAPPROVED {
				offerState = models.BoolPointer(true)
			}
			shipmentOfferAssertions := Assertions{
				ShipmentOffer: models.ShipmentOffer{
					ShipmentID:                      shipment.ID,
					TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
					Accepted:                        offerState,
				},
			}
			shipmentOffer := MakeShipmentOffer(db, shipmentOfferAssertions)
			shipmentOfferList = append(shipmentOfferList, shipmentOffer)
		}
	}

	if shouldCreateTariffData {
		createTariffDataForRateEngine(db, shipmentList[0])
	}

	return tspUserList, shipmentList, shipmentOfferList, nil
}

func createTariffDataForRateEngine(db *pop.Connection, shipment models.Shipment) {
	beforePickupDate := shipment.ActualPickupDate.AddDate(0, -6, 0)
	afterPickupDate := shipment.ActualPickupDate.AddDate(0, 6, 0)

	// $4861 is the cost for a 2000 pound move traveling 1044 miles (90210 to 80011).
	baseLinehaul := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1001,
		DistanceMilesUpper: 1101,
		WeightLbsLower:     2000,
		WeightLbsUpper:     2100,
		RateCents:          386400,
		Type:               "ConusLinehaul",
		EffectiveDateLower: beforePickupDate,
		EffectiveDateUpper: afterPickupDate,
	}
	mustSave(db, &baseLinehaul)

	// Create Service Area entries for Zip3s (which were already created)

	// Create fees for service areas
	sa1 := models.Tariff400ngServiceArea{
		Name:               "Los Angeles, CA",
		ServiceArea:        "56",
		ServicesSchedule:   3,
		LinehaulFactor:     unit.Cents(268),
		ServiceChargeCents: unit.Cents(775),
		EffectiveDateLower: beforePickupDate,
		EffectiveDateUpper: afterPickupDate,
		SIT185ARateCents:   unit.Cents(1626),
		SIT185BRateCents:   unit.Cents(60),
		SITPDSchedule:      3,
	}
	mustSave(db, &sa1)
	sa2 := models.Tariff400ngServiceArea{
		Name:               "Denver, CO Metro",
		ServiceArea:        "145",
		ServicesSchedule:   3,
		LinehaulFactor:     unit.Cents(174),
		ServiceChargeCents: unit.Cents(873),
		EffectiveDateLower: beforePickupDate,
		EffectiveDateUpper: afterPickupDate,
		SIT185ARateCents:   unit.Cents(1532),
		SIT185BRateCents:   unit.Cents(60),
		SITPDSchedule:      3,
	}
	mustSave(db, &sa2)

	fullPackRate := models.Tariff400ngFullPackRate{
		Schedule:           sa1.ServicesSchedule,
		WeightLbsLower:     0,
		WeightLbsUpper:     16001,
		RateCents:          6714,
		EffectiveDateLower: beforePickupDate,
		EffectiveDateUpper: afterPickupDate,
	}
	mustSave(db, &fullPackRate)

	fullUnpackRate := models.Tariff400ngFullUnpackRate{
		Schedule:           sa2.ServicesSchedule,
		RateMillicents:     704970,
		EffectiveDateLower: beforePickupDate,
		EffectiveDateUpper: afterPickupDate,
	}
	mustSave(db, &fullUnpackRate)

	// Set up item codes
	codeLHS := models.Tariff400ngItem{
		Code:                "LHS",
		Item:                "Linehaul Transportation",
		DiscountType:        models.Tariff400ngItemDiscountTypeHHG,
		AllowedLocation:     models.Tariff400ngItemAllowedLocationNEITHER,
		MeasurementUnit1:    models.Tariff400ngItemMeasurementUnitFLATRATE,
		MeasurementUnit2:    models.Tariff400ngItemMeasurementUnitNONE,
		RateRefCode:         models.Tariff400ngItemRateRefCodeTARIFFSECTION,
		RequiresPreApproval: false,
	}
	mustSave(db, &codeLHS)

	code135A := models.Tariff400ngItem{
		Code:                "135A",
		Item:                "Origin Service Charge",
		DiscountType:        models.Tariff400ngItemDiscountTypeHHG,
		AllowedLocation:     models.Tariff400ngItemAllowedLocationORIGIN,
		MeasurementUnit1:    models.Tariff400ngItemMeasurementUnitWEIGHT,
		MeasurementUnit2:    models.Tariff400ngItemMeasurementUnitNONE,
		RateRefCode:         models.Tariff400ngItemRateRefCodePOINTSCHEDULE,
		RequiresPreApproval: false,
	}
	mustSave(db, &code135A)

	code135B := models.Tariff400ngItem{
		Code:                "135B",
		Item:                "Destination Service Charge",
		DiscountType:        models.Tariff400ngItemDiscountTypeHHG,
		AllowedLocation:     models.Tariff400ngItemAllowedLocationDESTINATION,
		MeasurementUnit1:    models.Tariff400ngItemMeasurementUnitWEIGHT,
		MeasurementUnit2:    models.Tariff400ngItemMeasurementUnitNONE,
		RateRefCode:         models.Tariff400ngItemRateRefCodePOINTSCHEDULE,
		RequiresPreApproval: false,
	}
	mustSave(db, &code135B)

	code105A := models.Tariff400ngItem{
		Code:                "105A",
		Item:                "Full Pack",
		DiscountType:        models.Tariff400ngItemDiscountTypeHHG,
		AllowedLocation:     models.Tariff400ngItemAllowedLocationORIGIN,
		MeasurementUnit1:    models.Tariff400ngItemMeasurementUnitWEIGHT,
		MeasurementUnit2:    models.Tariff400ngItemMeasurementUnitNONE,
		RateRefCode:         models.Tariff400ngItemRateRefCodeNONE,
		RequiresPreApproval: false,
	}
	mustSave(db, &code105A)
}
