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

	"dp3/pkg/models"
)

func TestSubmitIssueHandler(t *testing.T) {
	newIssue := incomingIssue{"This is a test issue. The tests are not working. 🍏🍎😍"}

	body, err := json.Marshal(newIssue)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/issues", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(submitIssueHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	var response models.Issue
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode submitIssueResponse response - %s", err)
	}
}

func TestIndexIssuesHandler(t *testing.T) {
	// Given: An issue
	issueBody := "This is a test issue. The tests are not working. 🍏🍎😍"
	newIssue := incomingIssue{issueBody}

	body, err := json.Marshal(newIssue)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/issues", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	// When: New issue is posted
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(submitIssueHandler)
	handler.ServeHTTP(rr, req)

	// And: All issues are queried
	req = httptest.NewRequest("GET", "/issues", nil)
	w := httptest.NewRecorder()

	indexIssueHandler(w, req)

	resp := w.Result()

	// Then: Expect a 200 status code
	if resp.StatusCode != 200 {
		t.Errorf("Returned status code: %d", resp.StatusCode)
	}
	// And: Returned query to include our posted issue
	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)

	if m["body"] != issueBody {
		t.Errorf("Expected issue body to be '%v'. Got '%v'", issueBody, m["body"])
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
