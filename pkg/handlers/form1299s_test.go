package handlers

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/gen/messages"

	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
)

func TestSubmitForm1299HandlerAllValues(t *testing.T) {

	// Given: an instance of Form1299 with all valid values
	newForm1299Payload := messages.CreateForm1299Payload{
		ShipmentNumber:                   stringPointer("23098eifjsldkjf"),
		DatePrepared:                     fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		NameOfPreparingOffice:            stringPointer("random string bla"),
		DestOfficeName:                   stringPointer("random string bla"),
		OriginOfficeAddressName:          stringPointer("random string bla"),
		OriginOfficeAddress:              stringPointer("random string bla"),
		ServiceMemberFirstName:           stringPointer("random string bla"),
		ServiceMemberMiddleInitial:       stringPointer("random string bla"),
		ServiceMemberLastName:            stringPointer("random string bla"),
		ServiceMemberSsn:                 stringPointer("random string bla"),
		ServiceMemberAgency:              stringPointer("random string bla"),
		HhgTotalPounds:                   fmtInt64(10500),
		HhgProgearPounds:                 fmtInt64(100),
		HhgValuableItemsCartons:          fmtInt64(100),
		MobileHomeSerialNumber:           stringPointer("random string bla"),
		MobileHomeLengthFt:               fmtInt64(100),
		MobileHomeLengthInches:           fmtInt64(100),
		MobileHomeWidthFt:                fmtInt64(100),
		MobileHomeWidthInches:            fmtInt64(100),
		MobileHomeHeightFt:               fmtInt64(100),
		MobileHomeHeightInches:           fmtInt64(100),
		MobileHomeTypeExpando:            stringPointer("random string bla"),
		MobileHomeServicesRequested:      stringPointer("random string bla"), // enum validation not happening at server layer
		StationOrdersType:                stringPointer("random string bla"), // enum validation not happening at server layer
		StationOrdersIssuedBy:            stringPointer("random string bla"),
		StationOrdersNewAssignment:       stringPointer("random string bla"),
		StationOrdersDate:                fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		StationOrdersNumber:              stringPointer("random string bla"),
		StationOrdersParagraphNumber:     stringPointer("random string bla"),
		StationOrdersInTransitTelephone:  stringPointer("random string bla"),
		InTransitAddress:                 stringPointer("random string bla"),
		PickupAddress:                    stringPointer("random string bla"),
		PickupAddressMobileCourtName:     stringPointer("random string bla"),
		PickupTelephone:                  stringPointer("random string bla"),
		DestAddress:                      stringPointer("random string bla"),
		DestAddressMobileCourtName:       stringPointer("random string bla"),
		AgentToReceiveHhg:                stringPointer("random string bla"),
		ExtraAddress:                     stringPointer("random string bla"),
		PackScheduledDate:                fmtDate(time.Date(2019, 2, 6, 0, 0, 0, 0, time.UTC)),
		PickupScheduledDate:              fmtDate(time.Date(2019, 2, 7, 0, 0, 0, 0, time.UTC)),
		DeliveryScheduledDate:            fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		Remarks:                          stringPointer("random string bla"),
		OtherMoveFrom:                    stringPointer("random string bla"),
		OtherMoveTo:                      stringPointer("random string bla"),
		OtherMoveNetPounds:               fmtInt64(100),
		OtherMoveProgearPounds:           fmtInt64(100),
		ServiceMemberSignature:           stringPointer("random string bla"),
		DateSigned:                       fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		ContractorAddress:                stringPointer("random string bla"),
		ContractorName:                   stringPointer("random string bla"),
		NonavailabilityOfSignatureReason: stringPointer("random string bla"),
		CertifiedBySignature:             stringPointer("random string bla"),
		TitleOfCertifiedBySignature:      stringPointer("random string bla"),
	}

	// When: New Form1299 is posted
	newForm1299Params := form1299op.CreateForm1299Params{CreateForm1299Payload: &newForm1299Payload}
	response := CreateForm1299Handler(newForm1299Params)

	// Then: Assert we got back the 201 response
	createdResponse := response.(*form1299op.CreateForm1299Created)
	createdForm1299Payload := createdResponse.Payload

	// And: verify the values returned match expected values
	if *createdForm1299Payload.ShipmentNumber != *newForm1299Payload.ShipmentNumber ||
		(*createdForm1299Payload.DatePrepared != *newForm1299Payload.DatePrepared) ||
		(*createdForm1299Payload.NameOfPreparingOffice != *newForm1299Payload.NameOfPreparingOffice) ||
		(*createdForm1299Payload.DestOfficeName != *newForm1299Payload.DestOfficeName) ||
		(*createdForm1299Payload.OriginOfficeAddressName != *newForm1299Payload.OriginOfficeAddressName) ||
		(*createdForm1299Payload.OriginOfficeAddress != *newForm1299Payload.OriginOfficeAddress) ||
		(*createdForm1299Payload.ServiceMemberFirstName != *newForm1299Payload.ServiceMemberFirstName) ||
		(*createdForm1299Payload.ServiceMemberMiddleInitial != *newForm1299Payload.ServiceMemberMiddleInitial) ||
		(*createdForm1299Payload.ServiceMemberLastName != *newForm1299Payload.ServiceMemberLastName) ||
		(*createdForm1299Payload.ServiceMemberSsn != *newForm1299Payload.ServiceMemberSsn) ||
		(*createdForm1299Payload.ServiceMemberAgency != *newForm1299Payload.ServiceMemberAgency) ||
		(*createdForm1299Payload.HhgTotalPounds != *newForm1299Payload.HhgTotalPounds) ||
		(*createdForm1299Payload.HhgProgearPounds != *newForm1299Payload.HhgProgearPounds) ||
		(*createdForm1299Payload.HhgValuableItemsCartons != *newForm1299Payload.HhgValuableItemsCartons) ||
		(*createdForm1299Payload.MobileHomeSerialNumber != *newForm1299Payload.MobileHomeSerialNumber) ||
		(*createdForm1299Payload.MobileHomeLengthFt != *newForm1299Payload.MobileHomeLengthFt) ||
		(*createdForm1299Payload.MobileHomeLengthInches != *newForm1299Payload.MobileHomeLengthInches) ||
		(*createdForm1299Payload.MobileHomeWidthFt != *newForm1299Payload.MobileHomeWidthFt) ||
		(*createdForm1299Payload.MobileHomeWidthInches != *newForm1299Payload.MobileHomeWidthInches) ||
		(*createdForm1299Payload.MobileHomeHeightFt != *newForm1299Payload.MobileHomeHeightFt) ||
		(*createdForm1299Payload.MobileHomeHeightInches != *newForm1299Payload.MobileHomeHeightInches) ||
		(*createdForm1299Payload.MobileHomeTypeExpando != *newForm1299Payload.MobileHomeTypeExpando) ||
		(*createdForm1299Payload.MobileHomeServicesRequested != *newForm1299Payload.MobileHomeServicesRequested) ||
		(*createdForm1299Payload.StationOrdersType != *newForm1299Payload.StationOrdersType) ||
		(*createdForm1299Payload.StationOrdersIssuedBy != *newForm1299Payload.StationOrdersIssuedBy) ||
		(*createdForm1299Payload.StationOrdersNewAssignment != *newForm1299Payload.StationOrdersNewAssignment) ||
		(*createdForm1299Payload.StationOrdersDate != *newForm1299Payload.StationOrdersDate) ||
		(*createdForm1299Payload.StationOrdersNumber != *newForm1299Payload.StationOrdersNumber) ||
		(*createdForm1299Payload.StationOrdersParagraphNumber != *newForm1299Payload.StationOrdersParagraphNumber) ||
		(*createdForm1299Payload.StationOrdersInTransitTelephone != *newForm1299Payload.StationOrdersInTransitTelephone) ||
		(*createdForm1299Payload.InTransitAddress != *newForm1299Payload.InTransitAddress) ||
		(*createdForm1299Payload.PickupAddress != *newForm1299Payload.PickupAddress) ||
		(*createdForm1299Payload.PickupAddressMobileCourtName != *newForm1299Payload.PickupAddressMobileCourtName) ||
		(*createdForm1299Payload.PickupTelephone != *newForm1299Payload.PickupTelephone) ||
		(*createdForm1299Payload.DestAddress != *newForm1299Payload.DestAddress) ||
		(*createdForm1299Payload.DestAddressMobileCourtName != *newForm1299Payload.DestAddressMobileCourtName) ||
		(*createdForm1299Payload.AgentToReceiveHhg != *newForm1299Payload.AgentToReceiveHhg) ||
		(*createdForm1299Payload.ExtraAddress != *newForm1299Payload.ExtraAddress) ||
		(*createdForm1299Payload.PackScheduledDate != *newForm1299Payload.PackScheduledDate) ||
		(*createdForm1299Payload.PickupScheduledDate != *newForm1299Payload.PickupScheduledDate) ||
		(*createdForm1299Payload.DeliveryScheduledDate != *newForm1299Payload.DeliveryScheduledDate) ||
		(*createdForm1299Payload.Remarks != *newForm1299Payload.Remarks) ||
		(*createdForm1299Payload.OtherMoveFrom != *newForm1299Payload.OtherMoveFrom) ||
		(*createdForm1299Payload.OtherMoveTo != *newForm1299Payload.OtherMoveTo) ||
		(*createdForm1299Payload.OtherMoveNetPounds != *newForm1299Payload.OtherMoveNetPounds) ||
		(*createdForm1299Payload.OtherMoveProgearPounds != *newForm1299Payload.OtherMoveProgearPounds) ||
		(*createdForm1299Payload.ServiceMemberSignature != *newForm1299Payload.ServiceMemberSignature) ||
		(*createdForm1299Payload.DateSigned != *newForm1299Payload.DateSigned) ||
		(*createdForm1299Payload.ContractorAddress != *newForm1299Payload.ContractorAddress) ||
		(*createdForm1299Payload.ContractorName != *newForm1299Payload.ContractorName) ||
		(*createdForm1299Payload.NonavailabilityOfSignatureReason != *newForm1299Payload.NonavailabilityOfSignatureReason) ||
		(*createdForm1299Payload.CertifiedBySignature != *newForm1299Payload.CertifiedBySignature) ||
		(*createdForm1299Payload.TitleOfCertifiedBySignature != *newForm1299Payload.TitleOfCertifiedBySignature) {
		t.Error("Not all response values match expected values.")
	}

	// Then cofirm the same thing is returned by GET
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: createdForm1299Payload.ID}

	showResponse := ShowForm1299Handler(showFormParams)
	showOKResponse := showResponse.(*form1299op.ShowForm1299OK)
	showFormPayload := showOKResponse.Payload

}

func TestShowUnknownHandler(t *testing.T) {

	badID := strfmt.UUID("2400c3c5-019d-4031-9c27-8a553e022297")
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: badID}

	response := ShowForm1299Handler(showFormParams)

	// assert we got back the 404 response
	_ := response.(*form1299op.ShowForm1299NotFound)
}

func TestSubmitForm1299HandlerNoRequiredValues(t *testing.T) {

	// Given: an instance of Form1299 with no values
	// When: New Form1299 is posted
	newForm1299Payload := messages.CreateForm1299Payload{}
	newForm1299Params := form1299op.CreateForm1299Params{CreateForm1299Payload: &newForm1299Payload}
	response := CreateForm1299Handler(newForm1299Params)

	// Then: Assert we got back the 201 response
	createdResponse := response.(*form1299op.CreateForm1299Created)
	createdForm1299Payload := createdResponse.Payload

	// And: expected fields have values
	if (createdForm1299Payload.CreatedAt == nil) ||
		(createdForm1299Payload.ID == nil) {
		t.Error("CreatedAt time and ID should have values.")
	}

	// And: unset fields should be nil
	if createdForm1299Payload.HhgTotalPounds != nil {
		t.Error("There should not be anything sent back for HhgTotalPounds.")
	}
}

func TestSubmitForm1299HandlerSomeValues(t *testing.T) {
	// Given: an instance of Form1299 with some values
	newForm1299Payload := messages.CreateForm1299Payload{
		Remarks:                          stringPointer("random string bla"),
		OtherMoveFrom:                    stringPointer("random string bla"),
		OtherMoveTo:                      stringPointer("random string bla"),
		OtherMoveNetPounds:               fmtInt64(100),
		OtherMoveProgearPounds:           fmtInt64(100),
		ServiceMemberSignature:           stringPointer("random string bla"),
		DateSigned:                       fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		ContractorAddress:                stringPointer("random string bla"),
		ContractorName:                   stringPointer("random string bla"),
		NonavailabilityOfSignatureReason: stringPointer("random string bla"),
		CertifiedBySignature:             stringPointer("random string bla"),
		TitleOfCertifiedBySignature:      stringPointer("random string bla"),
	}

	// When: a new Form1299 is posted
	newForm1299Params := form1299op.CreateForm1299Params{CreateForm1299Payload: &newForm1299Payload}
	response := CreateForm1299Handler(newForm1299Params)

	// Then: Assert we got back the 201 response
	createdResponse := response.(*form1299op.CreateForm1299Created)
	createdForm1299Payload := createdResponse.Payload

	// And: expected fields have the right values
	if (createdForm1299Payload.CreatedAt == nil) ||
		(createdForm1299Payload.ID == nil) ||
		(*createdForm1299Payload.Remarks != *newForm1299Payload.Remarks) ||
		(*createdForm1299Payload.OtherMoveFrom != *newForm1299Payload.OtherMoveFrom) ||
		(*createdForm1299Payload.OtherMoveTo != *newForm1299Payload.OtherMoveTo) ||
		(*createdForm1299Payload.OtherMoveNetPounds != *newForm1299Payload.OtherMoveNetPounds) ||
		(*createdForm1299Payload.OtherMoveProgearPounds != *newForm1299Payload.OtherMoveProgearPounds) ||
		(*createdForm1299Payload.ServiceMemberSignature != *newForm1299Payload.ServiceMemberSignature) ||
		(*createdForm1299Payload.DateSigned != *newForm1299Payload.DateSigned) ||
		(*createdForm1299Payload.ContractorAddress != *newForm1299Payload.ContractorAddress) ||
		(*createdForm1299Payload.ContractorName != *newForm1299Payload.ContractorName) ||
		(*createdForm1299Payload.NonavailabilityOfSignatureReason != *newForm1299Payload.NonavailabilityOfSignatureReason) ||
		(*createdForm1299Payload.CertifiedBySignature != *newForm1299Payload.CertifiedBySignature) ||
		(*createdForm1299Payload.TitleOfCertifiedBySignature != *newForm1299Payload.TitleOfCertifiedBySignature) {
		t.Error("Not all response values match expected values.")
	}

	// And: unset fields should have nil values
	if (createdForm1299Payload.HhgTotalPounds != nil) ||
		(createdForm1299Payload.ServiceMemberFirstName != nil) {
		t.Error("Unset values should be nil.")
	}

}

func TestIndexForm1299sHandler(t *testing.T) {
	// Given: A Form1299
	shipmentNumber := "This is a test Form1299 for your indexForm1299Handler."
	newForm1299Payload := messages.CreateForm1299Payload{ShipmentNumber: &shipmentNumber}

	// When: New Form1299 is posted
	newForm1299Params := form1299op.CreateForm1299Params{CreateForm1299Payload: &newForm1299Payload}
	createResponse := CreateForm1299Handler(newForm1299Params)

	// Assert we got back the 201 response
	_ = createResponse.(*form1299op.CreateForm1299Created)

	// And: All Form1299s are queried
	indexForm1299sParams := form1299op.NewIndexForm1299sParams()
	indexResponse := IndexForm1299sHandler(indexForm1299sParams)

	// Then: Expect a 200 status code
	okResponse := indexResponse.(*form1299op.IndexForm1299sOK)
	form1299s := okResponse.Payload

	// And: Returned query to include our posted Form1299
	form1299Exists := false
	for _, form1299 := range form1299s {
		if form1299.ShipmentNumber != nil {
			if *form1299.ShipmentNumber == shipmentNumber {
				form1299Exists = true
			}
		}
	}

	if form1299Exists == false {
		t.Errorf("Expected an form1299 to contain '%v'. None do.", shipmentNumber)
	}
}
