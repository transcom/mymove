package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

func TestSwaggerYAMLHandler(t *testing.T) {

	req := httptest.NewRequest("GET", "/swagger.yaml", nil)
	w := httptest.NewRecorder()

	swaggerYAMLHandler(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("Returned status code: %d", resp.StatusCode)
	}

}

func TestSubmitIssueHandler(t *testing.T) {
	newIssue := incomingIssue{"This is a test issue. The tests are not working. 🍏🍎😍"}

	body, err := json.Marshal(newIssue)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/issues", bytes.NewReader(body))

	postResp := httptest.NewRecorder()
	submitIssueHandler(postResp, req)
	if status := postResp.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	var response models.Issue
	err = json.NewDecoder(postResp.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode submitIssueResponse response - %s", err)
	}
}

func TestIndexIssuesHandler(t *testing.T) {
	// Given: An issue
	issueBody := "This is a test issue for your indexIssueHandler."
	newIssue := incomingIssue{issueBody}

	body, err := json.Marshal(newIssue)
	if err != nil {
		t.Fatal(err)
	}
	postReq := httptest.NewRequest("POST", "/issues", bytes.NewReader(body))

	// When: New issue is posted
	postResp := httptest.NewRecorder()
	submitIssueHandler(postResp, postReq)

	// And: All issues are queried
	getReq := httptest.NewRequest("GET", "/issues", nil)
	getReqResp := httptest.NewRecorder()
	indexIssueHandler(getReqResp, getReq)
	resp := getReqResp.Result()

	// Then: Expect a 200 status code
	if resp.StatusCode != 200 {
		t.Errorf("Returned status code: %d", resp.StatusCode)
	}

	// And: Returned query to include our posted issue
	var issues []map[string]interface{}
	json.Unmarshal(getReqResp.Body.Bytes(), &issues)

	issueExists := false
	for _, issue := range issues {
		if issue["body"] == issueBody {
			issueExists = true
			break
		}
	}

	if issueExists == false {
		t.Errorf("Expected an issue to contain '%v'. None do.", issueBody)
	}
}

func setupDBConnection() {

	configLocation := "../../config"
	swaggerLocation := "../../../../../swagger.yaml" // ugh.
	pop.AddLookupPaths(configLocation)
	dbConnection, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	Init(dbConnection, swaggerLocation)

}

func TestMain(m *testing.M) {
	setupDBConnection()
	os.Exit(m.Run())
}
