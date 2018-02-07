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

func compareRequestAndResponsePayloads(t *testing.T, requestPayload messages.CreateForm1299Payload, responsePayload messages.Form1299Payload) {
	requestValue := reflect.ValueOf(requestPayload)
	requestType := requestValue.Type()
	responseValue := reflect.ValueOf(responsePayload)

	// iterate through all fields in request payload
	for i := 0; i < requestValue.NumField(); i++ {
		requestField := requestValue.Field(i)
		fieldName := requestType.Field(i).Name
		responseField := responseValue.FieldByName(fieldName)
		if responseField == (reflect.Value{}) {
			t.Errorf("Response has no field named %s", fieldName)
			continue
		}

		// First check that they are the same type. If not, error.
		if requestField.Type() != responseField.Type() {
			t.Errorf("Response field %s is of mismatched type. Request type: %s Response type: %s", fieldName, requestField.Type(), responseField.Type())
			continue
		}

		// Then, check if they are pointers, we have to check on them being nil
		if requestField.Kind() == reflect.Ptr {
			// if xor on the pointer being nil is true, it's a failure.
			// either both must be nil or, neither
			if requestField.IsNil() != responseField.IsNil() {
				t.Errorf("Response and Request field %s are not matching pointers: Request: %v, Response: %v", fieldName, requestField, responseField)
				continue
			} else if requestField.IsNil() {
				// If they are both nil, then they match. Nothing more to check
				continue
			}
			// If we arrive here, we know that they are both pointers and are both not nil.
		}

		// Indirect() turns a pointer into a type and does nothing to a type
		requestInterface := reflect.Indirect(requestField).Interface()
		responseInterface := reflect.Indirect(responseField).Interface()

		switch requestInterface.(type) {
		// Dates need to be compared with .Equal, not ==
		case strfmt.Date:
			df := time.Time(requestInterface.(strfmt.Date))
			dr := time.Time(responseInterface.(strfmt.Date))

			if !df.Equal(dr) {
				t.Errorf("%s doesn't match: request: %s response: %s", fieldName, requestInterface, responseInterface)
				continue
			}
		// Everything else can use == for now. Other cases may develop
		default:
			if requestInterface != responseInterface {
				t.Errorf("%s doesn't match, request: %s response: %s", fieldName, requestInterface, responseInterface)
				continue
			}
		}

	}
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
	compareRequestAndResponsePayloads(t, newForm1299Payload, *createdForm1299Payload)

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

	compareRequestAndResponsePayloads(t, newForm1299Payload, *showFormPayload)

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

	compareRequestAndResponsePayloads(t, newForm1299Payload, *createdForm1299Payload)

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
