package gex

import (
	"log"
	"net/http"
	"testing"
	//"github.com/stretchr/testify/mock"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

//type MockedHTTPRequest struct{
//	mock.Mock
//}
//
//func (m *MockedHTTPRequest) MakeRequest(request *http.Client) (resp http.Response) {
//
//	//arg := m.Called(request)
//	//return
//
//}

type GexSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *GexSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestGexSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &GexSuite{db: db, logger: logger}
	suite.Run(t, hs)
}

func (suite *GexSuite) TestGexSend_SendRequest_Actual() {

	//// create an instance of our test object
	//testObj := new(MockedHTTPRequest)
	//
	//// setup expectations
	//testObj.On("client.Do", request).Return(resp, nil)
	//ediString := ""
	//// call the code we are testing
	//SendGex{true}.SendRequest(ediString, "test file")
	//
	//// assert that the expectations were met
	//testObj.AssertExpectations(t)
}

func (suite *GexSuite) TestGexSend_SendRequest_Fake() {
	ediString := ""
	resp, _ := SendGex{false}.SendRequest(ediString, "test file")
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)
}
