package handlers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/gen/messages"

	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
)

func compareCreateAndResponsePayloads(createdPayload messages.CreateForm1299Payload, responsePayload messages.Form1299Payload) bool {
	v := reflect.ValueOf(createdPayload)
	t := v.Type()
	responseValue := reflect.ValueOf(responsePayload)

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i, t.Field(i).Name, f.Type(), f.Interface())
		responseField := responseValue.FieldByName(t.Field(i).Name)

		createInterface := f.Interface()
		responseInterface := responseField.Interface()

		switch createInterface.(type) {
		case *string:
			if createInterface == responseInterface && createInterface == nil {
				fmt.Println("PASS")
			} else if createInterface == nil || responseInterface == nil {
				fmt.Println("FAIL")
				// } else if *createInterface == *responseInterface {
				// 	fmt.Println("PASS")
			} else {
				fmt.Println("FAIL")
			}
		}

		if f.Type().String() == "*string" {
			fmt.Printf("Equal? %s, %s \n", reflect.Indirect(f).String(), reflect.Indirect(responseField).String())
		}
		if f.Type().String() == "*strfmt.Date" {
			df := time.Time(reflect.Indirect(f).Interface().(strfmt.Date))
			dr := time.Time(reflect.Indirect(responseField).Interface().(strfmt.Date))

			fmt.Printf("EqDate? %s, %s, %v\n", df, dr, df.Equal(dr))
		}
		if reflect.Indirect(f).Interface() != reflect.Indirect(responseField).Interface() {
			fmt.Printf("ERROR: Field %s is not equal. Created: %s Received: %s\n\n", t.Field(i).Name, reflect.Indirect(f).Interface(), reflect.Indirect(responseField).Interface())
			return false
		}
		// fmt.Printf("%d: %s %s = %v\n", i, responseField.Name , responseField.Type(), responseField.Interface())
	}

	return *createdPayload.ShipmentNumber == *responsePayload.ShipmentNumber &&
		(*createdPayload.DatePrepared == *responsePayload.DatePrepared) &&
		(*createdPayload.NameOfPreparingOffice == *responsePayload.NameOfPreparingOffice) &&
		(*createdPayload.DestOfficeName == *responsePayload.DestOfficeName) &&
		(*createdPayload.OriginOfficeAddressName == *responsePayload.OriginOfficeAddressName) &&
		(*createdPayload.OriginOfficeAddress == *responsePayload.OriginOfficeAddress) &&
		(*createdPayload.ServiceMemberFirstName == *responsePayload.ServiceMemberFirstName) &&
		(*createdPayload.ServiceMemberMiddleInitial == *responsePayload.ServiceMemberMiddleInitial) &&
		(*createdPayload.ServiceMemberLastName == *responsePayload.ServiceMemberLastName) &&
		(*createdPayload.ServiceMemberSsn == *responsePayload.ServiceMemberSsn) &&
		(*createdPayload.ServiceMemberAgency == *responsePayload.ServiceMemberAgency) &&
		(*createdPayload.HhgTotalPounds == *responsePayload.HhgTotalPounds) &&
		(*createdPayload.HhgProgearPounds == *responsePayload.HhgProgearPounds) &&
		(*createdPayload.HhgValuableItemsCartons == *responsePayload.HhgValuableItemsCartons) &&
		(*createdPayload.MobileHomeSerialNumber == *responsePayload.MobileHomeSerialNumber) &&
		(*createdPayload.MobileHomeLengthFt == *responsePayload.MobileHomeLengthFt) &&
		(*createdPayload.MobileHomeLengthInches == *responsePayload.MobileHomeLengthInches) &&
		(*createdPayload.MobileHomeWidthFt == *responsePayload.MobileHomeWidthFt) &&
		(*createdPayload.MobileHomeWidthInches == *responsePayload.MobileHomeWidthInches) &&
		(*createdPayload.MobileHomeHeightFt == *responsePayload.MobileHomeHeightFt) &&
		(*createdPayload.MobileHomeHeightInches == *responsePayload.MobileHomeHeightInches) &&
		(*createdPayload.MobileHomeTypeExpando == *responsePayload.MobileHomeTypeExpando) &&
		(*createdPayload.MobileHomeServicesRequested == *responsePayload.MobileHomeServicesRequested) &&
		(*createdPayload.StationOrdersType == *responsePayload.StationOrdersType) &&
		(*createdPayload.StationOrdersIssuedBy == *responsePayload.StationOrdersIssuedBy) &&
		(*createdPayload.StationOrdersNewAssignment == *responsePayload.StationOrdersNewAssignment) &&
		(*createdPayload.StationOrdersDate == *responsePayload.StationOrdersDate) &&
		(*createdPayload.StationOrdersNumber == *responsePayload.StationOrdersNumber) &&
		(*createdPayload.StationOrdersParagraphNumber == *responsePayload.StationOrdersParagraphNumber) &&
		(*createdPayload.StationOrdersInTransitTelephone == *responsePayload.StationOrdersInTransitTelephone) &&
		(*createdPayload.InTransitAddress == *responsePayload.InTransitAddress) &&
		(*createdPayload.PickupAddress == *responsePayload.PickupAddress) &&
		(*createdPayload.PickupAddressMobileCourtName == *responsePayload.PickupAddressMobileCourtName) &&
		(*createdPayload.PickupTelephone == *responsePayload.PickupTelephone) &&
		(*createdPayload.DestAddress == *responsePayload.DestAddress) &&
		(*createdPayload.DestAddressMobileCourtName == *responsePayload.DestAddressMobileCourtName) &&
		(*createdPayload.AgentToReceiveHhg == *responsePayload.AgentToReceiveHhg) &&
		(*createdPayload.ExtraAddress == *responsePayload.ExtraAddress) &&
		(*createdPayload.PackScheduledDate == *responsePayload.PackScheduledDate) &&
		(*createdPayload.PickupScheduledDate == *responsePayload.PickupScheduledDate) &&
		(*createdPayload.DeliveryScheduledDate == *responsePayload.DeliveryScheduledDate) &&
		(*createdPayload.Remarks == *responsePayload.Remarks) &&
		(*createdPayload.OtherMoveFrom == *responsePayload.OtherMoveFrom) &&
		(*createdPayload.OtherMoveTo == *responsePayload.OtherMoveTo) &&
		(*createdPayload.OtherMoveNetPounds == *responsePayload.OtherMoveNetPounds) &&
		(*createdPayload.OtherMoveProgearPounds == *responsePayload.OtherMoveProgearPounds) &&
		(*createdPayload.ServiceMemberSignature == *responsePayload.ServiceMemberSignature) &&
		(*createdPayload.DateSigned == *responsePayload.DateSigned) &&
		(*createdPayload.ContractorAddress == *responsePayload.ContractorAddress) &&
		(*createdPayload.ContractorName == *responsePayload.ContractorName) &&
		(*createdPayload.NonavailabilityOfSignatureReason == *responsePayload.NonavailabilityOfSignatureReason) &&
		(*createdPayload.CertifiedBySignature == *responsePayload.CertifiedBySignature) &&
		(*createdPayload.TitleOfCertifiedBySignature == *responsePayload.TitleOfCertifiedBySignature)
}

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
	if !compareCreateAndResponsePayloads(newForm1299Payload, *createdForm1299Payload) {
		t.Error("The response does not match what was created")
	}

	// Then cofirm the same thing is returned by GET
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: *createdForm1299Payload.ID}

	showResponse := ShowForm1299Handler(showFormParams)
	showOKResponse := showResponse.(*form1299op.ShowForm1299OK)
	showFormPayload := showOKResponse.Payload

	fmt.Println(newForm1299Payload)
	fmt.Println(*showFormPayload)

	b1, _ := json.MarshalIndent(newForm1299Payload, "", "  ")
	fmt.Println(string(b1))

	b2, _ := json.MarshalIndent(*createdForm1299Payload, "", "  ")
	fmt.Println(string(b2))

	b, _ := json.MarshalIndent(*showFormPayload, "", "  ")
	fmt.Println(string(b))

	if !compareCreateAndResponsePayloads(newForm1299Payload, *showFormPayload) {
		t.Error("The GET response does not match what was created")
	}

}

func TestShowUnknownHandler(t *testing.T) {

	badID := strfmt.UUID("2400c3c5-019d-4031-9c27-8a553e022297")
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: badID}

	response := ShowForm1299Handler(showFormParams)

	// assert we got back the 404 response
	_ = response.(*form1299op.ShowForm1299NotFound)
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
