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
	newIssue := incomingIssue{"This is a test issue. The tests are not working. 🍏🍎😍"}

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
	req := httptest.NewRequest("GET", "/issues", nil)
	w := httptest.NewRecorder()

	indexIssueHandler(w, req)

	resp := w.Result()

	// Then: Expect a 200 status code
	if resp.StatusCode != 200 {
		t.Errorf("Returned status code: %d", resp.StatusCode)
	}
	// And: Returned query to include our posted issue
	expected := `{"body": "This is a test issue. The tests are not working. 🍏🍎😍"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
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
