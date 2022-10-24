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

// The following data structures are set up for EvaluationReportFormFiller.subsection
// For every subsection in a report, we need:
// - a struct to hold the formatted values to put in the report
// - An array of all the field names in the order that they should be displayed
// - A map that matches up the field name used in the struct with the field name to display in the report
// We use reflection to look the string field names up in the struct.
// This is inspired by a pattern used for ShipmentSummaryWorksheet

type AdditionalKPIData struct {
	ObservedPickupSpreadStartDate string
	ObservedPickupSpreadEndDate   string
	ObservedClaimDate             string
	ObservedPickupDate            string
	ObservedDeliveryDate          string
}

var KPIFieldLabels = map[string]string{
	"ObservedPickupSpreadStartDate": "Observed pickup spread start date",
	"ObservedPickupSpreadEndDate":   "Observed pickup spread end date",
	"ObservedClaimDate":             "Observed claims response date",
	"ObservedPickupDate":            "Observed pickup date",
	"ObservedDeliveryDate":          "Observed delivery date",
}

type InspectionInformationValues struct {
	DateOfInspection                   string
	ReportSubmission                   string
	EvaluationType                     string
	TravelTimeToEvaluation             string
	EvaluationLocation                 string
	ObservedShipmentDeliveryDate       string
	ObservedShipmentPhysicalPickupDate string
	EvaluationLength                   string
	QAERemarks                         string
	ViolationsObserved                 string
	SeriousIncident                    string
	SeriousIncidentDescription         string
}

var InspectionInformationFields = []string{
	"DateOfInspection",
	"ReportSubmission",
	"EvaluationType",
	"TravelTimeToEvaluation",
	"ObservedShipmentDeliveryDate",
	"ObservedShipmentPhysicalPickupDate",
	"EvaluationLocation",
	"EvaluationLength",
}
var InspectionInformationFieldLabels = map[string]string{
	"DateOfInspection":                   "Date of inspection",
	"ReportSubmission":                   "Report submission",
	"EvaluationType":                     "Evaluation type",
	"TravelTimeToEvaluation":             "Travel time to evaluation",
	"EvaluationLocation":                 "Evaluation location",
	"ObservedShipmentPhysicalPickupDate": "Observed pickup date",
	"ObservedShipmentDeliveryDate":       "Observed delivery date",
	"EvaluationLength":                   "Evaluation length",
}

var ViolationsFields = []string{
	"ViolationsObserved",
	"SeriousIncident",
	"SeriousIncidentDescription",
}

var ViolationsFieldLabels = map[string]string{
	"ViolationsObserved":         "Violations observed",
	"SeriousIncident":            "Serious incident",
	"SeriousIncidentDescription": "Serious incident description",
}

var QAERemarksFields = []string{"QAERemarks"}
var QAERemarksFieldLabels = map[string]string{"QAERemarks": "Evaluation remarks"}

// ContactInformationValues holds formatted customer and QAE contact information
type ContactInformationValues struct {
	CustomerFullName    string
	CustomerPhone       string
	CustomerRank        string
	CustomerAffiliation string
	QAEFullName         string
	QAEPhone            string
	QAEEmail            string
}

// ShipmentValues holds formatted values to put in a shipment card
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

// TableRow is used to express the layout of one row of a shipment card
// Each row of the shipment card has two side by side key,value pairs
// The *FieldName properties are used to look up the values in a struct
// The *Label properties contain the label to display in the form
type TableRow struct {
	LeftFieldName  string
	LeftLabel      string
	RightFieldName string
	RightLabel     string
}

// Arrays of TableRow express the layout of a shipment card for a particular kind of shipment.
// They tell us which fields are included, and in what order.

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

func FormatValuesInspectionInformation(report models.EvaluationReport) InspectionInformationValues {
	inspectionInfo := InspectionInformationValues{}
	if report.InspectionDate != nil {
		inspectionInfo.DateOfInspection = report.InspectionDate.Format(dateFormat)
	}
	if report.SubmittedAt != nil {
		inspectionInfo.ReportSubmission = report.SubmittedAt.Format(dateFormat)
	}
	if report.InspectionType != nil {
		inspectionInfo.EvaluationType = formatEnum(string(*report.InspectionType))
	}
	if report.TravelTimeMinutes != nil {
		inspectionInfo.TravelTimeToEvaluation = formatDuration(*report.TravelTimeMinutes)
	}
	if report.Location != nil {
		inspectionInfo.EvaluationLocation = formatEnum(string(*report.Location))
		if *report.Location == models.EvaluationReportLocationTypeOther && report.LocationDescription != nil {
			inspectionInfo.EvaluationLocation += "\n" + *report.LocationDescription
		}
	}

	if report.ObservedShipmentDeliveryDate != nil {
		inspectionInfo.ObservedShipmentDeliveryDate = report.ObservedShipmentDeliveryDate.Format(dateFormat)
	}

	if report.ObservedShipmentPhysicalPickupDate != nil {
		inspectionInfo.ObservedShipmentPhysicalPickupDate = report.ObservedShipmentPhysicalPickupDate.Format(dateFormat)
	}

	if report.EvaluationLengthMinutes != nil {
		inspectionInfo.EvaluationLength = formatDuration(*report.EvaluationLengthMinutes)
	}
	if report.Remarks != nil {
		inspectionInfo.QAERemarks = *report.Remarks
	}
	inspectionInfo.ViolationsObserved = "No"
	if report.ViolationsObserved != nil && *report.ViolationsObserved {
		inspectionInfo.ViolationsObserved = "Yes\nViolations are listed on a subsequent page"
		inspectionInfo.SeriousIncident = "No"
		if report.SeriousIncident != nil && *report.SeriousIncident {
			inspectionInfo.SeriousIncident = "Yes"
			if report.SeriousIncidentDesc != nil {
				inspectionInfo.SeriousIncidentDescription = *report.SeriousIncidentDesc
			}
		}
	}
	return inspectionInfo
}

func formatDuration(minutes int) string {
	hours := minutes / 60
	remainingMinutes := minutes % 60
	return fmt.Sprintf("%d hr %d min", hours, remainingMinutes)
}

func formatEnum(e string) string {
	withSpaces := strings.ReplaceAll(e, "_", " ")
	return strings.ToUpper(withSpaces[:1]) + strings.ToLower(withSpaces[1:])
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
	if shipment.StorageFacility != nil {
		if shipment.StorageFacility.Phone != nil && shipment.StorageFacility.Email != nil {
			vals.StorageFacility = fmt.Sprintf("%s\n%s", *shipment.StorageFacility.Phone, *shipment.StorageFacility.Email)
		} else if shipment.StorageFacility.Phone != nil {
			vals.StorageFacility = *shipment.StorageFacility.Phone
		} else if shipment.StorageFacility.Email != nil {
			vals.StorageFacility = *shipment.StorageFacility.Email
		}

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

func formatSingleLineAddress(address models.Address) string {
	return strings.Join([]string{
		address.StreetAddress1,
		address.City,
		address.State,
		address.PostalCode,
	}, ", ")
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

// TODO we might be able to change the returns here to strings
// TODO the function that uses this doesnt care about any of the individual fields, it just joins them
func FormatContactInformationValues(customer models.ServiceMember, qae models.OfficeUser) ContactInformationValues {
	contactInfo := ContactInformationValues{
		QAEPhone:    qae.Telephone,
		QAEEmail:    qae.Email,
		QAEFullName: fmt.Sprintf("%s, %s", qae.LastName, qae.FirstName),
	}
	if customer.Telephone != nil {
		contactInfo.CustomerPhone = *customer.Telephone
	}
	if customer.Rank != nil {
		contactInfo.CustomerRank = rankDisplayValue[*customer.Rank]
	}
	if customer.Affiliation != nil {
		contactInfo.CustomerAffiliation = serviceMemberAffiliationDisplayValue[*customer.Affiliation]
	}

	contactInfo.CustomerFullName = customer.ReverseNameLineFormat()

	return contactInfo
}

func FormatAdditionalKPIValues(report models.EvaluationReport) AdditionalKPIData {
	additionalKPIData := AdditionalKPIData{}
	if report.ObservedPickupSpreadStartDate != nil {
		additionalKPIData.ObservedPickupSpreadStartDate = report.ObservedPickupSpreadStartDate.Format(dateFormat)
	}
	if report.ObservedPickupSpreadEndDate != nil {
		additionalKPIData.ObservedPickupSpreadEndDate = report.ObservedPickupSpreadEndDate.Format(dateFormat)
	}
	if report.ObservedClaimsResponseDate != nil {
		additionalKPIData.ObservedClaimDate = report.ObservedClaimsResponseDate.Format(dateFormat)
	}
	if report.ObservedPickupDate != nil {
		additionalKPIData.ObservedPickupDate = report.ObservedPickupDate.Format(dateFormat)
	}
	if report.ObservedDeliveryDate != nil {
		additionalKPIData.ObservedDeliveryDate = report.ObservedDeliveryDate.Format(dateFormat)
	}

	return additionalKPIData
}
