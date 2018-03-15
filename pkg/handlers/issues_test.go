package handlers

import (
	"time"

	"github.com/go-openapi/swag"

	issueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/issues"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func (suite *HandlerSuite) TestSubmitIssueHandler() {
	t := suite.T()

	testDescription := "This is a test issue. The tests are not working. 🍏🍎😍"
	newIssuePayload := internalmessages.CreateIssuePayload{Description: &testDescription}

	newIssueParams := issueop.CreateIssueParams{CreateIssuePayload: &newIssuePayload}

	handler := CreateIssueHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(newIssueParams)

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
	newIssuePayload := internalmessages.CreateIssuePayload{Description: &testDescription, DueDate: testDate}
	newIssueParams := issueop.CreateIssueParams{CreateIssuePayload: &newIssuePayload}

	handler := CreateIssueHandler(NewHandlerContext(suite.db, suite.logger))
	response := handler.Handle(newIssueParams)

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
	newIssuePayload := internalmessages.CreateIssuePayload{Description: &testDescription, ReporterName: swag.String("Jackie")}

	// When: New issue is posted
	newIssueParams := issueop.CreateIssueParams{CreateIssuePayload: &newIssuePayload}

	handlerContext := NewHandlerContext(suite.db, suite.logger)

	handler := CreateIssueHandler(handlerContext)
	createResponse := handler.Handle(newIssueParams)
	// Assert we got back the 201 response
	_ = createResponse.(*issueop.CreateIssueCreated)

	// And: All issues are queried
	indexIssuesParams := issueop.NewIndexIssuesParams()
	indexHandler := IndexIssuesHandler(handlerContext)
	indexResponse := indexHandler.Handle(indexIssuesParams)

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

	if !issueExists {
		t.Errorf("Expected an issue to contain '%v'. None do.", testDescription)
	}
}
