package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// FetchShipmentSummaryWorksheetFormValues fetches the pages for the Shipment Summary Worksheet for a given Shipment ID
func FetchShipmentSummaryWorksheetFormValues(db *pop.Connection, session *auth.Session, moveID uuid.UUID, preparationDate time.Time) (ShipmentSummaryWorksheetPage1Values, ShipmentSummaryWorksheetPage2Values, error) {
	var err error
	var page1 ShipmentSummaryWorksheetPage1Values
	page2 := ShipmentSummaryWorksheetPage2Values{}

	ssfd, err := FetchDataShipmentSummaryWorksheetFormData(db, session, moveID)
	ssfd.PreparationDate = preparationDate
	if err != nil {
		return page1, page2, err
	}
	page1 = FormatValuesShipmentSummaryWorksheetFormPage1(ssfd)
	return page1, page2, nil
}

// ShipmentSummaryWorksheetPage1Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage1Values struct {
	ServiceMemberName               string
	MaxSITStorageEntitlement        string
	PreferredPhone                  string
	PreferredEmail                  string
	DODId                           string
	ServiceBranch                   string
	Rank                            string
	IssuingBranchOrAgency           string
	OrdersIssueDate                 string
	OrdersTypeAndOrdersNumber       string
	AuthorizedOrigin                string
	AuthorizedDestination           string
	NewDutyAssignment               string
	WeightAllotmentSelf             string
	WeightAllotmentProgear          string
	WeightAllotmentProgearSpouse    string
	TotalWeightAllotment            string
	POVAuthorized                   string
	TAC                             string
	ShipmentNumberAndTypes          []string
	ShipmentPickUpDates             []string
	ShipmentWeights                 []string
	ShipmentCurrentShipmentStatuses []string
	PreparationDate                 string
}

//ShipmentSummaryWorkSheetShipments is and object representing shipment line items on Shipment Summary Worksheet
type ShipmentSummaryWorkSheetShipments struct {
	ShipmentNumberAndType string
	PickUpDate            string
	ShipmentWeight        string
	CurrentShipmentStatus string
}

// ShipmentSummaryWorksheetPage2Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage2Values struct {
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember           ServiceMember
	Order                   Order
	CurrentDutyStation      DutyStation
	NewDutyStation          DutyStation
	WeightAllotment         WeightAllotment
	Shipments               Shipments
	PreparationDate         time.Time
	PersonallyProcuredMoves PersonallyProcuredMoves
}

// FetchDataShipmentSummaryWorksheetFormData fetches the pages for the Shipment Summary Worksheet for a given Move ID
func FetchDataShipmentSummaryWorksheetFormData(db *pop.Connection, session *auth.Session, moveID uuid.UUID) (data ShipmentSummaryFormData, err error) {
	move := Move{}
	err = db.Q().Eager(
		"Orders",
		"Orders.NewDutyStation.Address",
		"Orders.ServiceMember",
		"Orders.ServiceMember.DutyStation.Address",
		"Shipments",
		"PersonallyProcuredMoves",
	).Find(&move, moveID)

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return ShipmentSummaryFormData{}, ErrFetchNotFound
		}
		return ShipmentSummaryFormData{}, err
	}

	_, authErr := FetchOrderForUser(db, session, move.OrdersID)
	if authErr != nil {
		return ShipmentSummaryFormData{}, authErr
	}

	serviceMember := move.Orders.ServiceMember
	var rank ServiceMemberRank
	if serviceMember.Rank != nil {
		rank = ServiceMemberRank(*serviceMember.Rank)
	}
	weightAllotment := GetWeightAllotment(rank)

	ssd := ShipmentSummaryFormData{
		ServiceMember:           serviceMember,
		Order:                   move.Orders,
		CurrentDutyStation:      serviceMember.DutyStation,
		NewDutyStation:          move.Orders.NewDutyStation,
		WeightAllotment:         weightAllotment,
		Shipments:               move.Shipments,
		PersonallyProcuredMoves: move.PersonallyProcuredMoves,
	}
	return ssd, nil
}

// FormatValuesShipmentSummaryWorksheetFormPage1 formats the data for page 1 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage1(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage1Values {
	page1 := ShipmentSummaryWorksheetPage1Values{}
	page1.MaxSITStorageEntitlement = "90 days per each shipment"
	// We don't currently know what allows POV to be authorized, so we are hardcoding it to "No" to start
	page1.POVAuthorized = "NO"
	page1.PreparationDate = FormatDate(data.PreparationDate)

	sm := data.ServiceMember
	page1.ServiceMemberName = FormatServiceMemberFullName(sm)
	page1.PreferredPhone = derefStringTypes(sm.Telephone)
	page1.PreferredEmail = derefStringTypes(sm.PersonalEmail)
	page1.DODId = derefStringTypes(sm.Edipi)

	page1.IssuingBranchOrAgency = FormatServiceMemberAffiliation(sm.Affiliation)
	page1.OrdersIssueDate = FormatDate(data.Order.IssueDate)
	page1.OrdersTypeAndOrdersNumber = FormatOrdersTypeAndOrdersNumber(data.Order)
	page1.TAC = derefStringTypes(data.Order.TAC)

	page1.AuthorizedOrigin = FormatAuthorizedLocation(data.CurrentDutyStation)
	page1.AuthorizedDestination = FormatAuthorizedLocation(data.NewDutyStation)
	page1.NewDutyAssignment = FormatDutyStation(data.NewDutyStation)

	page1.WeightAllotmentSelf = FormatWeights(data.WeightAllotment.TotalWeightSelf)
	page1.WeightAllotmentProgear = FormatWeights(data.WeightAllotment.ProGearWeight)
	page1.WeightAllotmentProgearSpouse = FormatWeights(data.WeightAllotment.ProGearWeightSpouse)
	total := data.WeightAllotment.TotalWeightSelf +
		data.WeightAllotment.ProGearWeight +
		data.WeightAllotment.ProGearWeightSpouse
	page1.TotalWeightAllotment = FormatWeights(total)

	formattedShipments := FormatAllShipments(data.PersonallyProcuredMoves, data.Shipments)
	// This will need to be revised slightly to handle multiple shipments
	if len(formattedShipments) != 0 {
		shipmentNumberAndTypes := make([]string, len(formattedShipments))
		shipmentPickUpDates := make([]string, len(formattedShipments))
		shipmentCurrentShipmentStatuses := make([]string, len(formattedShipments))
		shipmentWeights := make([]string, len(formattedShipments))
		for i, shipment := range formattedShipments {
			shipmentNumberAndTypes[i] = shipment.ShipmentNumberAndType
			shipmentPickUpDates[i] = shipment.PickUpDate
			shipmentCurrentShipmentStatuses[i] = shipment.CurrentShipmentStatus
			shipmentWeights[i] = shipment.ShipmentWeight
		}
		page1.ShipmentNumberAndTypes = shipmentNumberAndTypes
		page1.ShipmentPickUpDates = shipmentPickUpDates
		page1.ShipmentCurrentShipmentStatuses = shipmentCurrentShipmentStatuses
		page1.ShipmentWeights = shipmentWeights
	}
	return page1
}

//FormatAuthorizedLocation formats AuthorizedOrigin and AuthorizedDestination for Shipment Summary Worksheet
func FormatAuthorizedLocation(dutyStation DutyStation) string {
	return fmt.Sprintf("%s, %s %s", dutyStation.Name, dutyStation.Address.State, dutyStation.Address.PostalCode)
}

//FormatServiceMemberFullName formats ServiceMember full name for Shipment Summary Worksheet
func FormatServiceMemberFullName(serviceMember ServiceMember) string {
	lastName := derefStringTypes(serviceMember.LastName)
	suffix := derefStringTypes(serviceMember.Suffix)
	firstName := derefStringTypes(serviceMember.FirstName)
	middleName := derefStringTypes(serviceMember.MiddleName)
	if suffix != "" {
		return fmt.Sprintf("%s %s, %s %s", lastName, suffix, firstName, middleName)
	}
	return strings.TrimSpace(fmt.Sprintf("%s, %s %s", lastName, firstName, middleName))
}

//FormatAllShipments formats Shipment line items for the Shipment Summary Worksheet
func FormatAllShipments(ppms PersonallyProcuredMoves, shipments Shipments) []ShipmentSummaryWorkSheetShipments {
	formattedShipments := make([]ShipmentSummaryWorkSheetShipments, len(shipments)+len(ppms))
	for i, shipment := range shipments {
		formattedShipments[i].ShipmentNumberAndType = FormatShipmentNumberAndType(i)
		formattedShipments[i].PickUpDate = FormatShipmentPickupDate(shipment)
		formattedShipments[i].ShipmentWeight = FormatShipmentWeight(shipment)
		formattedShipments[i].CurrentShipmentStatus = FormatCurrentShipmentStatus(shipment)
	}
	for i, ppm := range ppms {
		j := i + len(shipments)
		formattedShipments[j].ShipmentNumberAndType = FormatPPMNumberAndType(j)
		formattedShipments[j].PickUpDate = FormatPPMPickupDate(ppm)
		// We don't have an actual weight for ppms yet, so we're leaving it blank for now
		formattedShipments[j].ShipmentWeight = ""
		formattedShipments[j].CurrentShipmentStatus = FormatCurrentPPMStatus(ppm)
	}
	return formattedShipments
}

//FormatCurrentShipmentStatus formats FormatCurrentShipmentStatus for the Shipment Summary Worksheet
func FormatCurrentShipmentStatus(shipment Shipment) string {
	return FormatEnum(string(shipment.Status))
}

//FormatCurrentPPMStatus formats FormatCurrentPPMStatus for the Shipment Summary Worksheet
func FormatCurrentPPMStatus(ppm PersonallyProcuredMove) string {
	if ppm.Status == "PAYMENT_REQUESTED" {
		return "At destination"
	}
	return FormatEnum(string(ppm.Status))
}

//FormatShipmentNumberAndType formats FormatShipmentNumberAndType for the Shipment Summary Worksheet
func FormatShipmentNumberAndType(i int) string {
	return fmt.Sprintf("%02d - HHG (GBL)", i+1)
}

//FormatPPMNumberAndType formats FormatShipmentNumberAndType for the Shipment Summary Worksheet
func FormatPPMNumberAndType(i int) string {
	return fmt.Sprintf("%02d - PPM", i+1)
}

//FormatShipmentWeight formats a shipments ShipmentWeight for the Shipment Summary Worksheet
func FormatShipmentWeight(shipment Shipment) string {
	if shipment.NetWeight != nil {
		wtg := FormatWeights(int(*shipment.NetWeight))
		return fmt.Sprintf("%s lbs - FINAL", wtg)
	}
	return ""
}

//FormatShipmentPickupDate formats a shipments ActualPickupDate for the Shipment Summary Worksheet
func FormatShipmentPickupDate(shipment Shipment) string {
	if shipment.ActualPickupDate != nil {
		return FormatDate(*shipment.ActualPickupDate)
	}
	return ""
}

//FormatPPMPickupDate formats a shipments ActualPickupDate for the Shipment Summary Worksheet
func FormatPPMPickupDate(ppm PersonallyProcuredMove) string {
	if ppm.PlannedMoveDate != nil {
		return FormatDate(*ppm.PlannedMoveDate)
	}
	return ""
}

//FormatDutyStation formats DutyStation for Shipment Summary Worksheet
func FormatDutyStation(dutyStation DutyStation) string {
	return fmt.Sprintf("%s, %s", dutyStation.Name, dutyStation.Address.State)
}

//FormatOrdersTypeAndOrdersNumber formats OrdersTypeAndOrdersNumber for Shipment Summary Worksheet
func FormatOrdersTypeAndOrdersNumber(order Order) string {
	issuingBranch := FormatOrdersType(order)
	ordersNumber := derefStringTypes(order.OrdersNumber)
	return fmt.Sprintf("%s/%s", issuingBranch, ordersNumber)
}

//FormatServiceMemberAffiliation formats ServiceMemberAffiliation in human friendly format
func FormatServiceMemberAffiliation(affiliation *ServiceMemberAffiliation) string {
	if affiliation != nil {
		return FormatEnum(string(*affiliation))
	}
	return ""
}

//FormatOrdersType formats OrdersType for Shipment Summary Worksheet
func FormatOrdersType(order Order) string {
	switch order.OrdersType {
	case internalmessages.OrdersTypePERMANENTCHANGEOFSTATION:
		return "PCS"
	default:
		return ""
	}
}

//FormatDate formats Dates for Shipment Summary Worksheet
func FormatDate(date time.Time) string {
	dateLayout := "02-Jan-2006"
	return date.Format(dateLayout)
}

//FormatEnum titlecases string const types (e.g. THIS_CONSTANT -> This Constant)
func FormatEnum(s string) string {
	words := strings.Split(strings.ToLower(s), "_")
	return strings.Title(strings.Join(words, " "))
}

//FormatWeights formats an int using 000s separator
func FormatWeights(wtg int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", wtg)
}

func derefStringTypes(st interface{}) string {
	switch v := st.(type) {
	case *string:
		if v != nil {
			return *v
		}
	case string:
		return v
	}
	return ""
}
