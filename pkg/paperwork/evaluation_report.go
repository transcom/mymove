package paperwork

import (
	"fmt"
	"strings"

	"github.com/transcom/mymove/pkg/models"
)

const (
	controlledUnclassifiedInformationText = "CONTROLLED UNCLASSIFIED INFORMATION"
	dateFormat                            = "02 January 2006"
)

type AdditionalKPIData struct {
	ObservedPickupSpreadStartDate string
	ObservedPickupSpreadEndDate   string
	ObservedClaimDate             string
	ObservedPickupDate            string
	ObservedDeliveryDate          string
}
type InspectionInformationValues struct {
	DateOfInspection           string
	ReportSubmission           string
	EvaluationType             string
	TravelTimeToEvaluation     string
	EvaluationLocation         string
	ObservedPickupDate         string
	ObservedDeliveryDate       string
	EvaluationLength           string
	QAERemarks                 string
	ViolationsObserved         string
	SeriousIncident            string
	SeriousIncidentDescription string
}
type ShipmentValues struct {
	ShipmentID            string
	ShipmentType          string
	PickupAddress         string
	ScheduledPickupDate   string
	RequestedPickupDate   string
	ActualPickupDate      string
	ReleasingAgentName    string
	ReleasingAgentPhone   string
	ReleasingAgentEmail   string
	DeliveryAddress       string
	ScheduledDeliveryDate string
	RequiredDeliveryDate  string
	ActualDeliveryDate    string
	ReceivingAgentName    string
	ReceivingAgentPhone   string
	ReceivingAgentEmail   string
	PPMOriginZIP          string
	PPMDestinationZIP     string
	PPMDepartureDate      string
	ReleasingAgent        string
	ReceivingAgent        string
	StorageFacility       string
	StorageFacilityName   string
	ExternalVendor        bool
}

var PPMShipmentCardLayout = []TableRow{
	{
		LeftFieldName:  "PPMOriginZIP",
		LeftLabel:      "Origin zip",
		RightFieldName: "PPMDestinationZIP",
		RightLabel:     "Destination zip",
	},
	{
		LeftFieldName:  "PPMDepartureDate",
		LeftLabel:      "Departure date",
		RightFieldName: "",
		RightLabel:     "",
	},
}
var HHGShipmentCardLayout = []TableRow{
	{
		LeftFieldName:  "ScheduledPickupDate",
		LeftLabel:      "Scheduled pickup date",
		RightFieldName: "ScheduledDeliveryDate",
		RightLabel:     "Scheduled delivery date",
	},
	{
		LeftFieldName:  "RequestedPickupDate",
		LeftLabel:      "Requested pickup date",
		RightFieldName: "RequiredDeliveryDate",
		RightLabel:     "Required delivery date",
	},
	{
		LeftFieldName:  "ActualPickupDate",
		LeftLabel:      "Actual pickup date",
		RightFieldName: "ActualDeliveryDate",
		RightLabel:     "Actual delivery date",
	},
	{
		LeftFieldName:  "ReleasingAgent",
		LeftLabel:      "Releasing agent",
		RightFieldName: "ReceivingAgent",
		RightLabel:     "Receiving agent",
	},
}

var NTSShipmentCardLayout = []TableRow{
	{
		LeftFieldName:  "ScheduledPickupDate",
		LeftLabel:      "Scheduled pickup date",
		RightFieldName: "ScheduledDeliveryDate",
		RightLabel:     "Scheduled delivery date",
	},
	{
		LeftFieldName:  "RequestedPickupDate",
		LeftLabel:      "Requested pickup date",
		RightFieldName: "RequiredDeliveryDate",
		RightLabel:     "Required delivery date",
	},
	{
		LeftFieldName:  "ActualPickupDate",
		LeftLabel:      "Actual pickup date",
		RightFieldName: "ActualDeliveryDate",
		RightLabel:     "Actual delivery date",
	},
	{
		LeftFieldName:  "ReleasingAgent",
		LeftLabel:      "Releasing agent",
		RightFieldName: "StorageFacility",
		RightLabel:     "Storage information",
	},
}

var NTSRShipmentCardLayout = []TableRow{
	{
		LeftFieldName:  "ScheduledPickupDate",
		LeftLabel:      "Scheduled pickup date",
		RightFieldName: "ScheduledDeliveryDate",
		RightLabel:     "Scheduled delivery date",
	},
	{
		LeftFieldName:  "RequestedPickupDate",
		LeftLabel:      "Requested pickup date",
		RightFieldName: "RequiredDeliveryDate",
		RightLabel:     "Required delivery date",
	},
	{
		LeftFieldName:  "ActualPickupDate",
		LeftLabel:      "Actual pickup date",
		RightFieldName: "ActualDeliveryDate",
		RightLabel:     "Actual delivery date",
	},
	{
		LeftFieldName:  "StorageFacility",
		LeftLabel:      "Storage information",
		RightFieldName: "ReceivingAgent",
		RightLabel:     "Receiving agent",
	},
}

func formatDuration(minutes int) string {
	hours := minutes / 60
	remainingMinutes := minutes % 60
	return fmt.Sprintf("%d hr %d min", hours, remainingMinutes)
}
func FormatValuesInspectionInformation(report models.EvaluationReport, violations models.PWSViolations) InspectionInformationValues {
	inspectionInfo := InspectionInformationValues{}
	if report.InspectionDate != nil {
		inspectionInfo.DateOfInspection = report.InspectionDate.Format(dateFormat)
	}
	if report.SubmittedAt != nil {
		inspectionInfo.ReportSubmission = report.SubmittedAt.Format(dateFormat)
	}
	if report.InspectionType != nil {
		inspectionInfo.EvaluationType = string(*report.InspectionType)
	}
	if report.TravelTimeMinutes != nil {
		inspectionInfo.TravelTimeToEvaluation = formatDuration(*report.TravelTimeMinutes)
	}
	if report.Location != nil {
		inspectionInfo.EvaluationLocation = string(*report.Location)
		if report.ObservedDate != nil {
			if *report.Location == models.EvaluationReportLocationTypeOrigin {
				inspectionInfo.ObservedPickupDate = report.ObservedDate.String()
			} else if *report.Location == models.EvaluationReportLocationTypeDestination {
				inspectionInfo.ObservedDeliveryDate = report.ObservedDate.String()
			}
		}
	}
	if report.EvaluationLengthMinutes != nil {
		inspectionInfo.EvaluationLength = formatDuration(*report.EvaluationLengthMinutes)
	}
	if report.Remarks != nil {
		inspectionInfo.QAERemarks = *report.Remarks
	}
	inspectionInfo.ViolationsObserved = "no"
	if report.ViolationsObserved != nil && *report.ViolationsObserved {
		inspectionInfo.ViolationsObserved = "yes\nAny violations recorded can be found on the following page"
		inspectionInfo.SeriousIncident = "yes"
		inspectionInfo.SeriousIncidentDescription = "a serious incident has happened!!! a serious incident has happened!!! a serious incident has happened!!! a serious incident has happened!!! a serious incident has happened!!! a serious incident has happened!!! a serious incident has happened!!! a serious incident has happened!!! a serious incident has happened!!!"
		// TODO seems like we don't have this in the model yet :(
		//if report.SeriousIncident != nil && *report.SeriousIncident {
		//	inspectionInfo.SeriousIncident = "yes"
		//	if report.SeriousIncidentDescription != nil {
		//		inspectionInfo.SeriousIncidentDescription = report.SeriousIncidentDescription
		//	}
		//}
	}
	return inspectionInfo
}

func formatSingleLineAddress(address models.Address) string {
	return strings.Join([]string{address.StreetAddress1, address.City, address.State, address.PostalCode}, ", ")
}

func formatMTOAgentInfo(agent models.MTOAgent) string {
	var lastName, firstName, phone, email string
	if agent.LastName != nil {
		lastName = *agent.LastName
	}
	if agent.FirstName != nil {
		firstName = *agent.FirstName
	}
	if agent.Phone != nil {
		phone = *agent.Phone
	}
	if agent.Email != nil {
		email = *agent.Email
	}
	contactInfo := fmt.Sprintf("%s, %s\n%s\n%s", lastName, firstName, phone, email)
	return contactInfo
}
func FormatValuesShipment(shipment models.MTOShipment) ShipmentValues {
	vals := ShipmentValues{
		ShipmentID:   strings.ToUpper(shipment.ID.String()[:5]),
		ShipmentType: string(shipment.ShipmentType),
	}
	if shipment.PPMShipment != nil {
		vals.PPMOriginZIP = shipment.PPMShipment.PickupPostalCode
		vals.PPMDestinationZIP = shipment.PPMShipment.DestinationPostalCode
		vals.PPMDepartureDate = shipment.PPMShipment.ExpectedDepartureDate.Format(dateFormat)
	}
	if shipment.StorageFacility != nil || shipment.StorageFacilityID != nil {
		fmt.Println("storage facility", shipment.StorageFacility)
		vals.StorageFacility = fmt.Sprintf("%s\n%s", *shipment.StorageFacility.Phone, *shipment.StorageFacility.Email)
		if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
			vals.PickupAddress = formatSingleLineAddress(shipment.StorageFacility.Address)
		}
		if shipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTSDom {
			vals.DeliveryAddress = formatSingleLineAddress(shipment.StorageFacility.Address)
		}
		vals.StorageFacilityName = strings.ToUpper(shipment.StorageFacility.FacilityName)
	}
	vals.ExternalVendor = shipment.UsesExternalVendor

	for _, agent := range shipment.MTOAgents {
		fmt.Println("found an agent", agent)
		contactInfo := formatMTOAgentInfo(agent)
		if agent.MTOAgentType == models.MTOAgentReleasing {
			vals.ReleasingAgent = contactInfo
		}
		if agent.MTOAgentType == models.MTOAgentReceiving {
			vals.ReceivingAgent = contactInfo
		}
	}
	if shipment.PickupAddress != nil {
		vals.PickupAddress = formatSingleLineAddress(*shipment.PickupAddress)
	}
	if shipment.DestinationAddress != nil {
		vals.DeliveryAddress = formatSingleLineAddress(*shipment.DestinationAddress)
	}
	if shipment.ScheduledPickupDate != nil {
		vals.ScheduledPickupDate = shipment.ScheduledPickupDate.Format(dateFormat)
	}
	if shipment.ActualPickupDate != nil {
		vals.ActualPickupDate = shipment.ActualPickupDate.Format(dateFormat)
	}
	if shipment.RequestedPickupDate != nil {
		vals.RequestedPickupDate = shipment.RequestedPickupDate.Format(dateFormat)
	}
	if shipment.RequiredDeliveryDate != nil {
		vals.RequiredDeliveryDate = shipment.RequiredDeliveryDate.Format(dateFormat)
	}
	return vals
}

func PickShipmentCardLayout(shipmentType models.MTOShipmentType) []TableRow {
	switch shipmentType {
	case models.MTOShipmentTypeHHG, models.MTOShipmentTypeHHGLongHaulDom, models.MTOShipmentTypeHHGShortHaulDom:
		return HHGShipmentCardLayout
	case models.MTOShipmentTypePPM:
		return PPMShipmentCardLayout
	case models.MTOShipmentTypeHHGIntoNTSDom:
		return NTSShipmentCardLayout
	case models.MTOShipmentTypeHHGOutOfNTSDom:
		return NTSRShipmentCardLayout
	case models.MTOShipmentTypeMotorhome:
		return []TableRow{}
	case models.MTOShipmentTypeBoatHaulAway:
		return []TableRow{}
	case models.MTOShipmentTypeBoatTowAway:
		return []TableRow{}
	default:
		return []TableRow{}
	}
}
