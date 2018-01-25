package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/gen/messages"
	issueop "github.com/transcom/mymove/pkg/gen/restapi/operations/issues"
)

func TestSubmitIssueHandler(t *testing.T) {

	testDescription := "This is a test issue. The tests are not working. 🍏🍎😍"
	newIssuePayload := messages.CreateIssuePayload{Description: &testDescription}

	newIssueParams := issueop.CreateIssueParams{CreateIssuePayload: &newIssuePayload}

	response := CreateIssueHandler(newIssueParams)

	// assert we got back the 201 response
	createdResponse := response.(*issueop.CreateIssueCreated)
	createdIssuePayload := createdResponse.Payload

	if *createdIssuePayload.Description != testDescription {
		t.Fatal("Didn't get the same description back")
	}

}

func TestIndexIssuesHandler(t *testing.T) {
	// Given: An issue
	testDescription := "This is a test issue for your indexIssueHandler."
	newIssuePayload := messages.CreateIssuePayload{Description: &testDescription}

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
