package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/gen/messages"
	"github.com/transcom/mymove/pkg/gen/restapi"
	"github.com/transcom/mymove/pkg/gen/restapi/operations"
	issueop "github.com/transcom/mymove/pkg/gen/restapi/operations/issues"
)

var testHandler http.Handler

func TestSubmitIssueHandler(t *testing.T) {
	testDescription := "This is a test issue. The tests are not working. 🍏🍎😍"
	newIssuePayload := messages.CreateIssuePayload{Description: &testDescription}

	createIssueBody, err := json.Marshal(newIssuePayload)
	if err != nil {
		t.Fatalf("Couldn't marshal %e", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/issues", bytes.NewReader(createIssueBody))
	req.Header.Add("Content-Type", "application/json")

	postResp := httptest.NewRecorder()

	testHandler.ServeHTTP(postResp, req)

	if status := postResp.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response contains what we expect.
	var response messages.IssueResponse
	err = json.NewDecoder(postResp.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode submitIssueResponse response - %s", err)
	}

	if *response.Description != testDescription {
		t.Fatal("Didn't get the same description back")
	}

}

func TestIndexIssuesHandler(t *testing.T) {
	// Given: An issue
	testDescription := "This is a test issue for your indexIssueHandler."
	newIssuePayload := messages.CreateIssuePayload{Description: &testDescription}

	newIssueBody, err := json.Marshal(newIssuePayload)
	if err != nil {
		t.Fatalf("Couldn't marshal %e", err)
	}

	postReq := httptest.NewRequest("POST", "/api/v1/issues", bytes.NewReader(newIssueBody))
	postReq.Header.Add("Content-Type", "application/json")

	// When: New issue is posted
	postResp := httptest.NewRecorder()
	testHandler.ServeHTTP(postResp, postReq)

	if status := postResp.Code; status != http.StatusCreated {
		t.Fatalf("create returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// And: All issues are queried
	indexReq := httptest.NewRequest("GET", "/api/v1/issues", nil)
	indexResp := httptest.NewRecorder()
	testHandler.ServeHTTP(indexResp, indexReq)
	resp := indexResp.Result()

	// Then: Expect a 200 status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Returned status code: %d", resp.StatusCode)
	}

	// And: Returned query to include our posted issue
	var response messages.IndexIssuesResponse
	err = json.NewDecoder(indexResp.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode submitIssueResponse response - %s", err)
	}

	issueExists := false
	for _, issue := range response {
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

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewMymoveAPI(swaggerSpec)

	api.IssuesCreateIssueHandler = issueop.CreateIssueHandlerFunc(CreateIssueHandler)
	api.IssuesIndexIssuesHandler = issueop.IndexIssuesHandlerFunc(IndexIssuesHandler)

	testHandler = api.Serve(nil)

	os.Exit(m.Run())
}
