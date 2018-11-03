package handlers

import (
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"runtime/debug"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/server"
	documentServices "github.com/transcom/mymove/pkg/services/document"
	userServices "github.com/transcom/mymove/pkg/services/user"
)

// BaseTestSuite abstracts the common methods needed for handler tests
type BaseTestSuite struct {
	suite.Suite
	db                 *pop.Connection
	logger             *zap.Logger
	filesToClose       []*runtime.File
	notificationSender notifications.NotificationSender
}

// TestDB returns a POP db connection for the suite
func (suite *BaseTestSuite) TestDB() *pop.Connection {
	return suite.db
}

// SetTestDB sets a POP db connection for the suite
func (suite *BaseTestSuite) SetTestDB(db *pop.Connection) {
	suite.db = db
}

// TestLogger returns the logger to use in the suite
func (suite *BaseTestSuite) TestLogger() *zap.Logger {
	return suite.logger
}

// SetTestLogger sets the logger to use in the suite
func (suite *BaseTestSuite) SetTestLogger(logger *zap.Logger) {
	suite.logger = logger
}

// TestFilesToClose returns the list of files needed to close at the end of tests
func (suite *BaseTestSuite) TestFilesToClose() []*runtime.File {
	return suite.filesToClose
}

// SetTestFilesToClose sets the list of files needed to close at the end of tests
func (suite *BaseTestSuite) SetTestFilesToClose(filesToClose []*runtime.File) {
	suite.filesToClose = filesToClose
}

// CloseFile adds a single file to close at the end of tests to the list of files
func (suite *BaseTestSuite) CloseFile(file *runtime.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

// TestNotificationSender returns the notification sender to use in the suite
func (suite *BaseTestSuite) TestNotificationSender() notifications.NotificationSender {
	return suite.notificationSender
}

// SetTestNotificationSender sets the notification sender to use in the suite
func (suite *BaseTestSuite) SetTestNotificationSender(notificationSender notifications.NotificationSender) {
	suite.notificationSender = notificationSender
}

// MustSave requires saving without errors
func (suite *BaseTestSuite) MustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

// IsNotErrResponse enforces handler does not return an error response
func (suite *BaseTestSuite) IsNotErrResponse(response middleware.Responder) {
	r, ok := response.(*ErrResponse)
	if ok {
		suite.logger.Error("Received an unexpected error response from handler: ", zap.Error(r.Err))
		// Formally lodge a complaint
		suite.IsType(&ErrResponse{}, response)
	}
}

// CheckErrorResponse verifies error response is what is expected
func (suite *BaseTestSuite) CheckErrorResponse(resp middleware.Responder, code int, name string) {
	errResponse, ok := resp.(*ErrResponse)
	if !ok || errResponse.Code != code {
		suite.T().Fatalf("Expected %s, Response: %v, Code: %v", name, resp, code)
		debug.PrintStack()
	}
}

// CheckResponseBadRequest looks at BadRequest errors
func (suite *BaseTestSuite) CheckResponseBadRequest(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusBadRequest, "BadRequest")
}

// CheckResponseUnauthorized looks at Unauthorized errors
func (suite *BaseTestSuite) CheckResponseUnauthorized(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusUnauthorized, "Unauthorized")
}

// CheckResponseForbidden looks at Forbidden errors
func (suite *BaseTestSuite) CheckResponseForbidden(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusForbidden, "Forbidden")
}

// CheckResponseNotFound looks at NotFound errors
func (suite *BaseTestSuite) CheckResponseNotFound(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusNotFound, "NotFound")
}

// CheckResponseInternalServerError looks at InternalServerError errors
func (suite *BaseTestSuite) CheckResponseInternalServerError(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusInternalServerError, "InternalServerError")
}

// CheckResponseTeapot enforces that response come from a Teapot
func (suite *BaseTestSuite) CheckResponseTeapot(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusTeapot, "Teapot")
}

// AuthenticateRequest Request authenticated with a service member
func (suite *BaseTestSuite) AuthenticateRequest(req *http.Request, serviceMember models.ServiceMember) *http.Request {
	session := server.Session{
		ApplicationName: server.MyApp,
		UserID:          serviceMember.UserID,
		IDToken:         "fake token",
		ServiceMemberID: serviceMember.ID,
	}
	ctx := server.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateUserRequest only authenticated with a bare user - have no idea if they are a service member yet
func (suite *BaseTestSuite) AuthenticateUserRequest(req *http.Request, user models.User) *http.Request {
	session := server.Session{
		ApplicationName: server.MyApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := server.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateOfficeRequest authenticates Office users
func (suite *BaseTestSuite) AuthenticateOfficeRequest(req *http.Request, user models.OfficeUser) *http.Request {
	session := server.Session{
		ApplicationName: server.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	ctx := server.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateTspRequest authenticates TSP users
func (suite *BaseTestSuite) AuthenticateTspRequest(req *http.Request, user models.TspUser) *http.Request {
	session := server.Session{
		ApplicationName: server.TspApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		TspUserID:       user.ID,
	}
	ctx := server.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// Fixture allows us to include a fixture like a PDF in the test
func (suite *BaseTestSuite) Fixture(name string) *runtime.File {
	fixtureDir := "fixtures"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Error(err)
	}

	fixturePath := path.Join(cwd, "..", fixtureDir, name)

	// #nosec never comes from user input
	file, err := os.Open(fixturePath)
	if err != nil {
		suite.logger.Fatal("Error opening fixture file", zap.Error(err))
	}

	info, err := file.Stat()
	if err != nil {
		suite.logger.Fatal("Error accessing fixture stats", zap.Error(err))
	}

	header := multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	returnFile := &runtime.File{
		Header: &header,
		Data:   file,
	}
	suite.CloseFile(returnFile)

	return returnFile
}

// HandlerContextWithServices constructs a handler context with the service layers dependencies populated
func (suite *BaseTestSuite) HandlerContextWithServices() HandlerContext {
	h := NewHandlerContext(suite.TestDB(), suite.TestLogger())
	fetchServiceMember := userServices.NewFetchServiceMemberService(models.NewServiceMemberDB(suite.TestDB()))
	h.SetFetchServiceMember(fetchServiceMember)
	documentDB := models.NewDocumentDB(suite.TestDB())
	fetchDocument := documentServices.NewFetchDocumentService(documentDB, fetchServiceMember)
	h.SetFetchDocument(fetchDocument)
	fetchUpload := documentServices.NewFetchUploadService(documentDB, fetchDocument)
	h.SetFetchUpload(fetchUpload)
	return h
}
