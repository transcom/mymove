package handlers

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/gen/messages"
	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299"
)

func TestSubmitForm1299Handler(t *testing.T) {

	testDescription := "This is a test Form1299. The tests are not working. 🍏🍎😍"
	newForm1299Payload := messages.CreateForm1299Payload{Description: &testDescription}

	newForm1299Params := form1299op.CreateForm1299Params{CreateForm1299Payload: &newForm1299Payload}

	response := CreateForm1299Handler(newForm1299Params)

	// assert we got back the 201 response
	createdResponse := response.(*form1299op.CreateForm1299Created)
	createdForm1299Payload := createdResponse.Payload

	if *createdForm1299Payload.Description != testDescription {
		t.Error("Didn't get the same description back")
	}

	if createdForm1299Payload.ReporterName != nil {
		t.Error("We should not have sent anything back for the reporter name")
	}

}

func TestSubmitDueDate(t *testing.T) {
	testDescription := "This is a test Form1299. The tests are not working. 🍏🍎😍"
	testDate := fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC))
	newForm1299Payload := messages.CreateForm1299Payload{Description: &testDescription, DueDate: testDate}
	newForm1299Params := form1299op.CreateForm1299Params{CreateForm1299Payload: &newForm1299Payload}

	response := CreateForm1299Handler(newForm1299Params)

	// assert we got back the 201 response
	createdResponse := response.(*form1299op.CreateForm1299Created)
	createdForm1299Payload := createdResponse.Payload

	if createdForm1299Payload.DueDate != testDate {
		t.Error("Didn't get the same date back")
	}

	if createdForm1299Payload.ReporterName != nil {
		t.Error("We should not have sent anything back for the reporter name")
	}
}

func TestIndexForm1299sHandler(t *testing.T) {
	// Given: An Form1299
	testDescription := "This is a test Form1299 for your indexForm1299Handler."
	newForm1299Payload := messages.CreateForm1299Payload{Description: &testDescription}

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
	Form1299s := okResponse.Payload

	// And: Returned query to include our posted Form1299
	form1299Exists := false
	for _, form1299 := range form1299s {
		if *form1299.Description == testDescription {
			form1299Exists = true
			break
		}
	}

	if form1299Exists == false {
		t.Errorf("Expected an form1299 to contain '%v'. None do.", testDescription)
	}
}

func setupDBConnection() {

	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	dbConnection, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	Init(dbConnection)

}

func TestMain(m *testing.M) {
	setupDBConnection()

	os.Exit(m.Run())
}
