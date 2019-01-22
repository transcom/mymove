package models

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// FetchShipmentSummaryWorksheetFormValues fetches the pages for the Shipment Summary Worksheet for a given Shipment ID
func FetchShipmentSummaryWorksheetFormValues(db *pop.Connection, moveID uuid.UUID) (ShipmentSummaryWorksheetPage1Values, ShipmentSummaryWorksheetPage2Values, error) {
	var err error
	var ssfd ShipmentSummaryFormData
	var page1 ShipmentSummaryWorksheetPage1Values
	page2 := ShipmentSummaryWorksheetPage2Values{}

	ssfd, err = FetchDataShipmentSummaryWorksFormData(db, moveID)
	if err != nil {
		return page1, page2, err
	}
	page1 = FormatValuesShipmentSummaryWorksheetFormPage1(ssfd)
	return page1, page2, nil
}

// ShipmentSummaryWorksheetPage1Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage1Values struct {
	ServiceMemberName            string
	MaxSITStorageEntitlement     string
	PreferredPhone               string
	PreferredEmail               string
	DODId                        string
	ServiceBranch                string
	Rank                         string
	OrdersNumber                 string
	IssuingAgency                string
	OrderIssueDate               time.Time
	OrdersType                   internalmessages.OrdersType
	DutyStationID                uuid.UUID
	AuthorizedOrigin             DutyStation
	NewDutyStationID             uuid.UUID
	AuthorizedDestination        DutyStation
	WeightAllotmentSelf          string
	WeightAllotmentProgear       string
	WeightAllotmentProgearSpouse string
	TotalWeightAllotment         string
}

// ShipmentSummaryWorksheetPage2Values is an object representing a Shipment Summary Worksheet
type ShipmentSummaryWorksheetPage2Values struct {
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember      ServiceMember
	Order              Order
	CurrentDutyStation DutyStation
	NewDutyStation     DutyStation
	WeightAllotment    WeightAllotment
}

// FetchDataShipmentSummaryWorksFormData fetches the data required for the Shipment Summary Worksheet
func FetchDataShipmentSummaryWorksFormData(db *pop.Connection, moveID uuid.UUID) (data ShipmentSummaryFormData, err error) {
	ssd := ShipmentSummaryFormData{}
	reqFields, err := getRequiredFields(db, moveID)
	if err != nil {
		return ssd, err
	}
	ssd.Order, err = FetchOrder(db, reqFields.OrdersID)
	if err != nil {
		return ssd, err
	}
	ssd.ServiceMember, err = FetchServiceMember(db, reqFields.ServiceMemberID)
	if err != nil {
		return ssd, err
	}
	ssd.CurrentDutyStation, err = FetchDutyStation(context.TODO(), db, reqFields.ServiceMemberDutyStationID)
	if err != nil {
		return ssd, err
	}
	ssd.NewDutyStation, err = FetchDutyStation(context.TODO(), db, ssd.Order.NewDutyStationID)
	if err != nil {
		return ssd, err
	}
	rank := ServiceMemberRank(reqFields.ServiceMemberRank)
	ssd.WeightAllotment = GetWeightAllotment(rank)
	return ssd, nil
}

// FormatValuesShipmentSummaryWorksheetFormPage1 formats the data for page 1 of the Shipment Summary Worksheet
func FormatValuesShipmentSummaryWorksheetFormPage1(data ShipmentSummaryFormData) ShipmentSummaryWorksheetPage1Values {
	page1 := ShipmentSummaryWorksheetPage1Values{}
	page1.MaxSITStorageEntitlement = "90 days per each shipment"

	sm := data.ServiceMember
	lastName := derefStringTypes(sm.LastName)
	suffix := derefStringTypes(sm.Suffix)
	firstName := derefStringTypes(sm.FirstName)
	middleName := derefStringTypes(sm.MiddleName)
	fullName := fmt.Sprintf("%s %s, %s %s", lastName, suffix, firstName, middleName)
	page1.ServiceMemberName = fullName
	page1.PreferredPhone = derefStringTypes(sm.Telephone)
	page1.PreferredEmail = derefStringTypes(sm.PersonalEmail)
	page1.DODId = derefStringTypes(sm.Edipi)
	page1.ServiceBranch = derefStringTypes(sm.Affiliation)
	page1.Rank = derefStringTypes(sm.Rank)

	o := data.Order
	page1.OrdersNumber = derefStringTypes(o.OrdersNumber)
	page1.IssuingAgency = derefStringTypes(o.OrdersIssuingAgency)
	page1.OrderIssueDate = o.IssueDate
	page1.OrdersType = o.OrdersType

	page1.AuthorizedOrigin = data.CurrentDutyStation
	page1.AuthorizedDestination = data.NewDutyStation
	page1.WeightAllotmentSelf = FormatWeights(data.WeightAllotment.TotalWeightSelf)
	page1.WeightAllotmentProgear = FormatWeights(data.WeightAllotment.ProGearWeight)
	page1.WeightAllotmentProgearSpouse = FormatWeights(data.WeightAllotment.ProGearWeightSpouse)
	total := data.WeightAllotment.TotalWeightSelf +
		data.WeightAllotment.ProGearWeight +
		data.WeightAllotment.ProGearWeightSpouse
	page1.TotalWeightAllotment = FormatWeights(total)
	return page1
}

//FormatWeights formats an int using 000s separator
func FormatWeights(wtg int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", wtg)
}

type requiredFields struct {
	OrdersID                   uuid.UUID `db:"orders_id"`
	ServiceMemberID            uuid.UUID `db:"service_member_id"`
	ServiceMemberDutyStationID uuid.UUID `db:"duty_station_id"`
	ServiceMemberRank          string    `db:"rank"`
}

func getRequiredFields(db *pop.Connection, moveID uuid.UUID) (requiredFields, error) {
	var err error
	p := requiredFields{}
	sql := `
		SELECT orders_id,
			   service_member_id,
			   duty_station_id,
			   rank
		FROM moves m
				 INNER JOIN orders o ON m.orders_id = o.id
				 INNER JOIN service_members sm ON o.service_member_id = sm.id
		WHERE m.id = $1`
	err = db.RawQuery(sql, moveID).First(&p)
	return p, err
}

func derefStringTypes(st interface{}) string {
	switch v := st.(type) {
	case *string:
		if v != nil {
			return *v
		}
	case string:
		return v
	case *ServiceMemberRank:
		if v != nil {
			return string(*v)
		}
		return ""
	case *ServiceMemberAffiliation:
		if v != nil {
			return string(*v)
		}
		return ""
	}
	return ""
}
