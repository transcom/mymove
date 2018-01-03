package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSwaggerYAMLHandler(t *testing.T) {

	req := httptest.NewRequest("GET", "/swagger.yaml", nil)
	w := httptest.NewRecorder()

	swaggerYAMLHandler(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		fmt.Println("Returned status code: ", resp.StatusCode)
		t.Fail()
	}

}

func TestSubmitIssueHandler(t *testing.T) {
	newIssue := issue{"This is a test issue. The tests are not working. 🍏🍎😍"}

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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var response newIssueResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode submitIssueResponse response - %s", err)
	}
}
