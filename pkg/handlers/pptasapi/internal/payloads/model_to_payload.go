package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *pptasmessages.ClientError {
	instanceToUse := strfmt.UUID(traceID.String())
	payload := pptasmessages.ClientError{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: &instanceToUse,
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ListReport payload
func ListReport(move *models.Move) *pptasmessages.ListReport {
	if move == nil {
		return nil
	}

	Orders := move.Orders

	progear := unit.Pound(0)
	sitTotal := unit.Pound(0)
	travelAdvance := unit.Cents(0)

	var moveDate *time.Time
	if move.MTOShipments[0].PPMShipment != nil {
		moveDate = &move.MTOShipments[0].PPMShipment.ExpectedDepartureDate
	} else {
		moveDate = move.MTOShipments[0].ActualPickupDate
	}

	payload := &pptasmessages.ListReport{
		// ID:        *report.ID,
		FirstName:          *Orders.ServiceMember.FirstName,
		LastName:           *Orders.ServiceMember.LastName,
		MiddleInitial:      *Orders.ServiceMember.MiddleName,
		Affiliation:        (*pptasmessages.Affiliation)(Orders.ServiceMember.Affiliation),
		PayGrade:           (*string)(Orders.Grade),
		Edipi:              *Orders.ServiceMember.Edipi,
		PhonePrimary:       *Orders.ServiceMember.Telephone,
		PhoneSecondary:     Orders.ServiceMember.SecondaryTelephone,
		EmailPrimary:       *Orders.ServiceMember.PersonalEmail,
		EmailSecondary:     &Orders.ServiceMember.BackupContacts[0].Email,
		OrdersType:         string(Orders.OrdersType),
		OrdersNumber:       *Orders.OrdersNumber,
		OrdersDate:         strfmt.DateTime(Orders.IssueDate),
		Address:            Address(Orders.ServiceMember.ResidentialAddress),
		OriginAddress:      Address(move.MTOShipments[0].PickupAddress),
		DestinationAddress: Address(move.MTOShipments[0].DestinationAddress),
		OriginGbloc:        Orders.OriginDutyLocationGBLOC,
		DestinationGbloc:   nil, // &move.CloseoutOffice.Gbloc,
		DepCD:              nil,
		TravelAdvance:      models.Float64Pointer(travelAdvance.Float64()), // report.TravelAdvance,
		MoveDate:           (*strfmt.Date)(moveDate),
		Tac:                Orders.TAC,
		FiscalYear:         nil,
		Appro:              nil, // report.Appro,
		Subhead:            nil, // report.Subhead,
		ObjClass:           nil, // report.ObjClass,
		Bcn:                nil, // report.BCN,
		SubAllotCD:         nil, // report.SubAllotCD,
		Aaa:                nil, // report.AAA,
		TypeCD:             nil, // report.TypeCD,
		Paa:                nil, // report.PAA,
		CostCD:             nil, // report.CostCD,
		Ddcd:               nil, // report.DDCD,
		ShipmentNum:        int64(len(move.MTOShipments)),
		WeightEstimate:     calculateTotalWeightEstimate(move.MTOShipments).Float64(),
		TransmitCD:         nil, // report.TransmitCd,
		Dd2278IssueDate:    strfmt.Date(*move.ServiceCounselingCompletedAt),
		Miles:              int64(*move.MTOShipments[0].Distance),
		WeightAuthorized:   0.0, // float64(Orders.Entitlement.WeightAllotted.TotalWeightSelfPlusDependents), // WeightAlloted isn't returning any value
		ShipmentID:         strfmt.UUID(move.ID.String()),
		Scac:               nil, // I don't know what gbloc to use // HSFR
		// Loa:                         nil, // report.LOA,
		// ShipmentType:                "",  // *report.ShipmentType,
		EntitlementWeight: int64(*Orders.Entitlement.DBAuthorizedWeight),
		NetWeight:         int64(models.GetTotalNetWeightForMove(*move)), // this only calculates PPM is that correct?
		PickupDate:        strfmt.Date(*move.MTOShipments[0].ActualPickupDate),
		// Rate:                        nil, // report.Rate,
		// PaidDate:                    (*strfmt.Date)(report.PaidDate),
		// LinehaulTotal:               nil, // report.LinehaulTotal,
		// SitTotal:                    nil, // report.SitTotal,
		// AccessorialTotal:            nil, // report.AccessorialTotal,
		// FuelTotal:                   nil, // report.FuelTotal,
		// OtherTotal:                  nil, // report.OtherTotal,
		// InvoicePaidAmt:              0.0, // report.InvoicePaidAmt.Float64(),
		TravelType:      string(*Orders.OrdersTypeDetail),
		TravelClassCode: string(Orders.OrdersType),
		DeliveryDate:    strfmt.Date(*moveDate),
		// ActualOriginNetWeight:       0, // *report.ActualOriginNetWeight,
		// DestinationReweighNetWeight: 0, // report.DestinationReweighNetWeight.Float64(),
		CounseledDate: strfmt.Date(*move.ServiceCounselingCompletedAt),
	}

	// sharing this for loop for all MTOShipment calculations
	for _, shipment := range move.MTOShipments {
		// calculate total progear for entire move
		if shipment.PPMShipment != nil {
			shipmentTotalProgear := shipment.PPMShipment.ProGearWeight.Float64() + shipment.PPMShipment.SpouseProGearWeight.Float64()
			progear += unit.Pound(shipmentTotalProgear)

			// need to determine which shipment(s) have a ppm and get the travel advances and add them up
			if shipment.PPMShipment.AdvanceAmountReceived != nil {
				travelAdvance += *shipment.PPMShipment.AdvanceAmountReceived
			}

			// add SIT estimated weights
			if *shipment.PPMShipment.SITExpected {
				sitTotal += *shipment.PPMShipment.SITEstimatedWeight

				// SIT Fields
				payload.SitInDate = (*strfmt.Date)(shipment.PPMShipment.SITEstimatedEntryDate)
				payload.SitOutDate = (*strfmt.Date)(shipment.PPMShipment.SITEstimatedDepartureDate)
				// newreport.SitType = // Example data is destination.. ??
			}
		}
	}

	payload.PbpAnde = progear.Float64()

	// SAC is currently optional, is it acceptable to have an empty return here?
	if Orders.SAC != nil {
		payload.OrdersNumber = *Orders.SAC
	} else {
		payload.OrderNumber = ""
	}

	return payload
}

// ListReports payload
func ListReports(moves *models.Moves) []*pptasmessages.ListReport {
	payload := make(pptasmessages.ListReports, len(*moves))

	for i, move := range *moves {
		copyOfMove := move // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = ListReport(&copyOfMove)
	}
	return payload
}

func Address(address *models.Address) *pptasmessages.Address {
	if address == nil {
		return nil
	}
	return &pptasmessages.Address{
		ID:             strfmt.UUID(address.ID.String()),
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           &address.City,
		State:          &address.State,
		PostalCode:     &address.PostalCode,
		Country:        address.Country,
		County:         &address.County,
		ETag:           etag.GenerateEtag(address.UpdatedAt),
	}
}

func calculateTotalWeightEstimate(shipments models.MTOShipments) *unit.Pound {
	var weightEstimate unit.Pound
	for _, shipment := range shipments {
		if shipment.PPMShipment != nil {
			weightEstimate += *shipment.PPMShipment.EstimatedWeight
		}

		if shipment.PrimeEstimatedWeight != nil {
			weightEstimate += *shipment.PrimeEstimatedWeight
		}
	}

	return &weightEstimate
}
