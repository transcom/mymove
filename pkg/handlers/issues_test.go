package handlers

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/markbates/pop"

	issueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/issues"
	"github.com/transcom/mymove/pkg/gen/internalmodel"
)

func TestSubmitIssueHandler(t *testing.T) {

	testDescription := "This is a test issue. The tests are not working. 🍏🍎😍"
	newIssuePayload := internalmodel.CreateIssuePayload{Description: &testDescription}

	newIssueParams := issueop.CreateIssueParams{CreateIssuePayload: &newIssuePayload}

	response := CreateIssueHandler(newIssueParams)

	// assert we got back the 201 response
	createdResponse := response.(*issueop.CreateIssueCreated)
	createdIssuePayload := createdResponse.Payload

	if *createdIssuePayload.Description != testDescription {
		t.Error("Didn't get the same description back")
	}

	if createdIssuePayload.ReporterName != nil {
		t.Error("We should not have sent anything back for the reporter name")
	}

}

func TestSubmitDueDate(t *testing.T) {
	testDescription := "This is a test issue. The tests are not working. 🍏🍎😍"
	testDate := fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC))
	newIssuePayload := internalmodel.CreateIssuePayload{Description: &testDescription, DueDate: testDate}
	newIssueParams := issueop.CreateIssueParams{CreateIssuePayload: &newIssuePayload}

	response := CreateIssueHandler(newIssueParams)

	// assert we got back the 201 response
	createdResponse := response.(*issueop.CreateIssueCreated)
	createdIssuePayload := createdResponse.Payload

	if createdIssuePayload.DueDate != testDate {
		t.Error("Didn't get the same date back")
	}

	if createdIssuePayload.ReporterName != nil {
		t.Error("We should not have sent anything back for the reporter name")
	}
}

func TestIndexIssuesHandler(t *testing.T) {
	// Given: An issue
	testDescription := "This is a test issue for your indexIssueHandler."
	newIssuePayload := internalmodel.CreateIssuePayload{Description: &testDescription}

	// When: New issue is posted
	newIssueParams := issueop.CreateIssueParams{CreateIssuePayload: &newIssuePayload}

	createResponse := CreateIssueHandler(newIssueParams)
	// Assert we got back the 201 response
	_ = createResponse.(*issueop.CreateIssueCreated)

	// And: All issues are queried
	indexIssuesParams := issueop.NewIndexIssuesParams()
	indexResponse := IndexIssuesHandler(indexIssuesParams)

	// Then: Expect a 200 status code
	okResponse := indexResponse.(*issueop.IndexIssuesOK)
	issues := okResponse.Payload

	// And: Returned query to include our posted issue
	issueExists := false
	for _, issue := range issues {
		if *issue.Description == testDescription {
			issueExists = true
			break
		}
	}

	if issueExists == false {
		t.Errorf("Expected an issue to contain '%v'. None do.", testDescription)
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
