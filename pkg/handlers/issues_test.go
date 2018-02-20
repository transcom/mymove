package handlers

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/messages"
	issueop "github.com/transcom/mymove/pkg/gen/restapi/operations/issues"
)

func (suite *HandlerSuite) TestSubmitIssueHandler() {
	t := suite.T()

	testDescription := "This is a test issue. The tests are not working. 🍏🍎😍"
	newIssuePayload := messages.CreateIssuePayload{Description: &testDescription}

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

func (suite *HandlerSuite) TestSubmitDueDate() {
	t := suite.T()

	testDescription := "This is a test issue. The tests are not working. 🍏🍎😍"
	testDate := fmtDate(time.Date(2019, 2, 8, 0, 0, 0, 0, time.UTC))
	newIssuePayload := messages.CreateIssuePayload{Description: &testDescription, DueDate: testDate}
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

func (suite *HandlerSuite) TestIndexIssuesHandler() {
	t := suite.T()

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
