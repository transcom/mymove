package handlers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"go.uber.org/zap"

	form1299op "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/form1299s"
)

func compareRequestAndResponsePayloads(t *testing.T, requestPayload interface{}, responsePayload interface{}) {
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
		case internalmessages.Address:
			compareRequestAndResponsePayloads(t, requestInterface, responseInterface)
		default:
			if requestInterface != responseInterface {
				t.Errorf("%s doesn't match, request: %s response: %s", fieldName, requestInterface, responseInterface)
				continue
			}
		}
	}
}

func fakeAddress() *internalmessages.Address {
	return &internalmessages.Address{
		StreetAddress1: swag.String("An address"),
		StreetAddress2: swag.String("Apt. 2"),
		City:           swag.String("Happytown"),
		State:          swag.String("AL"),
		Zip:            swag.String("01234"),
	}
}

// Sets up a basic logger so logs are printed to stdout during tests
func setUpLogger() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.L().Info("replaced zap's global loggers")
}

func (suite *HandlerSuite) TestSubmitForm1299HandlerAllValues() {
	t := suite.T()

	var rankE6 = internalmessages.ServiceMemberRankE6
	setUpLogger()
	// Given: an instance of Form1299 with all valid values
	newForm1299Payload := internalmessages.CreateForm1299Payload{
		ShipmentNumber:                         swag.String("23098eifjsldkjf"),
		DatePrepared:                           fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		NameOfPreparingOffice:                  swag.String("random string bla"),
		DestOfficeName:                         swag.String("random string bla"),
		OriginOfficeAddressName:                swag.String("random string bla"),
		OriginOfficeAddress:                    fakeAddress(),
		ServiceMemberFirstName:                 swag.String("random string bla"),
		ServiceMemberMiddleInitial:             swag.String("random string bla"),
		ServiceMemberLastName:                  swag.String("random string bla"),
		ServiceMemberSsn:                       swag.String("random string bla"),
		ServiceMemberAgency:                    swag.String("random string bla"),
		ServiceMemberRank:                      &rankE6,
		HhgTotalPounds:                         fmtInt64(10500),
		HhgProgearPounds:                       fmtInt64(100),
		HhgValuableItemsCartons:                fmtInt64(100),
		MobileHomeSerialNumber:                 swag.String("random string bla"),
		MobileHomeLengthFt:                     fmtInt64(100),
		MobileHomeLengthInches:                 fmtInt64(100),
		MobileHomeWidthFt:                      fmtInt64(100),
		MobileHomeWidthInches:                  fmtInt64(100),
		MobileHomeHeightFt:                     fmtInt64(100),
		MobileHomeHeightInches:                 fmtInt64(100),
		MobileHomeTypeExpando:                  swag.String("random string bla"),
		MobileHomeContentsPackedRequested:      swag.Bool(true),
		MobileHomeBlockedRequested:             swag.Bool(false),
		MobileHomeUnblockedRequested:           swag.Bool(true),
		MobileHomeStoredAtOriginRequested:      swag.Bool(false),
		MobileHomeStoredAtDestinationRequested: swag.Bool(true),
		StationOrdersType:                      swag.String("random string bla"), // enum validation not happening at server layer
		StationOrdersIssuedBy:                  swag.String("random string bla"),
		StationOrdersNewAssignment:             swag.String("random string bla"),
		StationOrdersDate:                      fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		StationOrdersNumber:                    swag.String("random string bla"),
		StationOrdersParagraphNumber:           swag.String("random string bla"),
		StationOrdersInTransitTelephone:        swag.String("random string bla"),
		InTransitAddress:                       fakeAddress(),
		PickupAddress:                          fakeAddress(),
		PickupTelephone:                        swag.String("random string bla"),
		DestAddress:                            fakeAddress(),
		AgentToReceiveHhg:                      swag.String("random string bla"),
		ExtraAddress:                           fakeAddress(),
		PackScheduledDate:                      fmtDate(time.Date(2019, 2, 6, 0, 0, 0, 0, time.UTC)),
		PickupScheduledDate:                    fmtDate(time.Date(2019, 2, 7, 0, 0, 0, 0, time.UTC)),
		DeliveryScheduledDate:                  fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		Remarks:                                swag.String("random string bla"),
		OtherMove1From:                         swag.String("random string bla"),
		OtherMove1To:                           swag.String("random string bla"),
		OtherMove1NetPounds:                    fmtInt64(100),
		OtherMove1ProgearPounds:                fmtInt64(100),
		ServiceMemberSignature:                 swag.String("random string bla"),
		DateSigned:                             fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		ContractorAddress:                      fakeAddress(),
		ContractorName:                         swag.String("random string bla"),
		NonavailabilityOfSignatureReason:       swag.String("random string bla"),
		CertifiedBySignature:                   swag.String("random string bla"),
		TitleOfCertifiedBySignature:            swag.String("random string bla"),
	}

	// When: New Form1299 is posted
	newForm1299Params := form1299op.CreateForm1299Params{CreateForm1299Payload: &newForm1299Payload}
	response := CreateForm1299Handler(newForm1299Params)

	// Then: Assert we got back the 201 response
	createdResponse := response.(*form1299op.CreateForm1299Created)
	createdForm1299Payload := createdResponse.Payload

	// And: verify the values returned match expected values
	compareRequestAndResponsePayloads(t, newForm1299Payload, *createdForm1299Payload)

	// Then confirm the same thing is returned by GET
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: *createdForm1299Payload.ID}

	showResponse := ShowForm1299Handler(showFormParams)
	showOKResponse := showResponse.(*form1299op.ShowForm1299OK)
	showFormPayload := showOKResponse.Payload

	b1, _ := json.MarshalIndent(newForm1299Payload, "", "  ")
	fmt.Println(string(b1))

	b2, _ := json.MarshalIndent(*createdForm1299Payload, "", "  ")
	fmt.Println(string(b2))

	b, _ := json.MarshalIndent(*showFormPayload, "", "  ")
	fmt.Println(string(b))

	compareRequestAndResponsePayloads(t, newForm1299Payload, *showFormPayload)
}

func (suite *HandlerSuite) TestShowUnknown() {
	unknownID := strfmt.UUID("2400c3c5-019d-4031-9c27-8a553e022297")
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: unknownID}

	response := ShowForm1299Handler(showFormParams)

	// assert we got back the 404 response
	_ = response.(*form1299op.ShowForm1299NotFound)
}

func (suite *HandlerSuite) TestShowBadID() {
	badID := strfmt.UUID("2400c3c5-019d-4031-9c27-8a553e022297xxx")
	showFormParams := form1299op.ShowForm1299Params{Form1299ID: badID}

	response := ShowForm1299Handler(showFormParams)

	// assert we got back the 400 response
	_ = response.(*form1299op.ShowForm1299BadRequest)
}

func (suite *HandlerSuite) TestSubmitForm1299HandlerNoRequiredValues() {
	t := suite.T()

	// Given: an instance of Form1299 with no values
	// When: New Form1299 is posted
	newForm1299Payload := internalmessages.CreateForm1299Payload{
		MobileHomeContentsPackedRequested:      swag.Bool(false),
		MobileHomeBlockedRequested:             swag.Bool(false),
		MobileHomeUnblockedRequested:           swag.Bool(false),
		MobileHomeStoredAtOriginRequested:      swag.Bool(false),
		MobileHomeStoredAtDestinationRequested: swag.Bool(false),
	}
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

func (suite *HandlerSuite) TestSubmitForm1299HandlerSomeValues() {
	t := suite.T()

	// Given: an instance of Form1299 with some values
	newForm1299Payload := internalmessages.CreateForm1299Payload{
		Remarks:                                swag.String("random string bla"),
		OtherMove1From:                         swag.String("random string bla"),
		OtherMove1To:                           swag.String("random string bla"),
		OtherMove1NetPounds:                    fmtInt64(100),
		OtherMove1ProgearPounds:                fmtInt64(100),
		ServiceMemberSignature:                 swag.String("random string bla"),
		DateSigned:                             fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC)),
		MobileHomeContentsPackedRequested:      swag.Bool(false),
		MobileHomeBlockedRequested:             swag.Bool(false),
		MobileHomeUnblockedRequested:           swag.Bool(false),
		MobileHomeStoredAtOriginRequested:      swag.Bool(false),
		MobileHomeStoredAtDestinationRequested: swag.Bool(false),
		ContractorAddress:                      fakeAddress(),
		ContractorName:                         swag.String("random string bla"),
		NonavailabilityOfSignatureReason:       swag.String("random string bla"),
		CertifiedBySignature:                   swag.String("random string bla"),
		TitleOfCertifiedBySignature:            swag.String("random string bla"),
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

func (suite *HandlerSuite) TestIndexForm1299sHandler() {
	t := suite.T()

	// Given: A Form1299
	destOfficeName := "This is a test Form1299 for your indexForm1299Handler."
	newForm1299Payload := internalmessages.CreateForm1299Payload{
		DestOfficeName:                         &destOfficeName,
		MobileHomeContentsPackedRequested:      swag.Bool(false),
		MobileHomeBlockedRequested:             swag.Bool(false),
		MobileHomeUnblockedRequested:           swag.Bool(false),
		MobileHomeStoredAtOriginRequested:      swag.Bool(false),
		MobileHomeStoredAtDestinationRequested: swag.Bool(false),
		ContractorAddress:                      fakeAddress(),
	}

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
		if form1299.DestOfficeName != nil {
			if *form1299.DestOfficeName == destOfficeName {
				form1299Exists = true
			}
		}
	}

	if form1299Exists == false {
		t.Errorf("Expected an form1299 to contain '%v'. None do.", destOfficeName)
	}
}
