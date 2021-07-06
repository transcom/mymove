package handlers

import (
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// BaseHandlerTestSuite abstracts the common methods needed for handler tests
type BaseHandlerTestSuite struct {
	testingsuite.PopTestSuite
	logger             Logger
	filesToClose       []*runtime.File
	notificationSender notifications.NotificationSender
}

// NewBaseHandlerTestSuite returns a new BaseHandlerTestSuite
func NewBaseHandlerTestSuite(logger Logger, sender notifications.NotificationSender, packageName testingsuite.PackageName, opts ...testingsuite.PopTestSuiteOption) BaseHandlerTestSuite {
	return BaseHandlerTestSuite{
		PopTestSuite:       testingsuite.NewPopTestSuite(packageName, opts...),
		logger:             logger,
		notificationSender: sender,
	}
}

// TestLogger returns the logger to use in the suite
func (suite *BaseHandlerTestSuite) TestLogger() Logger {
	return suite.logger
}

// TestFilesToClose returns the list of files needed to close at the end of tests
func (suite *BaseHandlerTestSuite) TestFilesToClose() []*runtime.File {
	return suite.filesToClose
}

// SetTestFilesToClose sets the list of files needed to close at the end of tests
func (suite *BaseHandlerTestSuite) SetTestFilesToClose(filesToClose []*runtime.File) {
	suite.filesToClose = filesToClose
}

// CloseFile adds a single file to close at the end of tests to the list of files
func (suite *BaseHandlerTestSuite) CloseFile(file *runtime.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

// TestNotificationSender returns the notification sender to use in the suite
func (suite *BaseHandlerTestSuite) TestNotificationSender() notifications.NotificationSender {
	return suite.notificationSender
}

// HasWebhookNotification checks that there's a record on the WebhookNotifications table for the object and trace IDs
func (suite *BaseHandlerTestSuite) HasWebhookNotification(objectID uuid.UUID, traceID uuid.UUID) {
	notification := &models.WebhookNotification{}
	err := suite.DB().Where("object_id = $1 AND trace_id = $2", objectID.String(), traceID.String()).First(notification)
	suite.NoError(err)
}

// HasNoWebhookNotification checks that there's no record on the WebhookNotifications table for the object and trace IDs
func (suite *BaseHandlerTestSuite) HasNoWebhookNotification(objectID uuid.UUID, traceID uuid.UUID) {
	notification := &models.WebhookNotification{}
	numRows, err := suite.DB().Where("object_id = $1 AND trace_id = $2", objectID.String(), traceID.String()).Count(notification)
	suite.NoError(err)
	suite.Equal(numRows, 0)
}

// IsNotErrResponse enforces handler does not return an error response
func (suite *BaseHandlerTestSuite) IsNotErrResponse(response middleware.Responder) {
	r, ok := response.(*ErrResponse)
	if ok {
		suite.logger.Error("Received an unexpected error response from handler: ", zap.Error(r.Err))
		// Formally lodge a complaint
		suite.IsType(&ErrResponse{}, response)
	}
}

// CheckErrorResponse verifies error response is what is expected
func (suite *BaseHandlerTestSuite) CheckErrorResponse(resp middleware.Responder, code int, name string) {
	errResponse, ok := resp.(*ErrResponse)
	if !ok || errResponse.Code != code {
		suite.T().Fatalf("Expected %s, Response: %v, Code: %v", name, resp, code)
		debug.PrintStack()
	}
}

// CheckNotErrorResponse verifies there is no error response
func (suite *BaseHandlerTestSuite) CheckNotErrorResponse(resp middleware.Responder) {
	errResponse, ok := resp.(*ErrResponse)
	if ok {
		suite.NoError(errResponse.Err)
		suite.FailNowf("Received error response", "Code: %v", errResponse.Code)
	}
}

// CheckResponseBadRequest looks at BadRequest errors
func (suite *BaseHandlerTestSuite) CheckResponseBadRequest(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusBadRequest, "BadRequest")
}

// CheckResponseUnauthorized looks at Unauthorized errors
func (suite *BaseHandlerTestSuite) CheckResponseUnauthorized(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusUnauthorized, "Unauthorized")
}

// CheckResponseForbidden looks at Forbidden errors
func (suite *BaseHandlerTestSuite) CheckResponseForbidden(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusForbidden, "Forbidden")
}

// CheckResponseNotFound looks at NotFound errors
func (suite *BaseHandlerTestSuite) CheckResponseNotFound(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusNotFound, "NotFound")
}

// CheckResponseInternalServerError looks at InternalServerError errors
func (suite *BaseHandlerTestSuite) CheckResponseInternalServerError(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusInternalServerError, "InternalServerError")
}

// CheckResponseTeapot enforces that response come from a Teapot
func (suite *BaseHandlerTestSuite) CheckResponseTeapot(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusTeapot, "Teapot")
}

// AuthenticateRequest Request authenticated with a service member
func (suite *BaseHandlerTestSuite) AuthenticateRequest(req *http.Request, serviceMember models.ServiceMember) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          serviceMember.UserID,
		IDToken:         "fake token",
		ServiceMemberID: serviceMember.ID,
		Email:           serviceMember.User.LoginGovEmail,
	}
	session.Roles = append(session.Roles, serviceMember.User.Roles...)
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateUserRequest only authenticated with a bare user - have no idea if they are a service member yet
func (suite *BaseHandlerTestSuite) AuthenticateUserRequest(req *http.Request, user models.User) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateOfficeRequest authenticates Office users
func (suite *BaseHandlerTestSuite) AuthenticateOfficeRequest(req *http.Request, user models.OfficeUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	session.Roles = append(session.Roles, user.User.Roles...)
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateDpsRequest authenticates DPS users
func (suite *BaseHandlerTestSuite) AuthenticateDpsRequest(req *http.Request, serviceMember models.ServiceMember, dpsUser models.DpsUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          serviceMember.UserID,
		IDToken:         "fake token",
		DpsUserID:       dpsUser.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateAdminRequest authenticates DPS users
func (suite *BaseHandlerTestSuite) AuthenticateAdminRequest(req *http.Request, user models.User) *http.Request {
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// Fixture allows us to include a fixture like a PDF in the test
func (suite *BaseHandlerTestSuite) Fixture(name string) *runtime.File {
	fixtureDir := "testdatagen/testdata"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Error(err)
	}

	fixturePath := path.Join(cwd, "..", "..", fixtureDir, name)

	file, err := os.Open(filepath.Clean(fixturePath))
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

// EqualDatePtr compares the time.Time from the model with the strfmt.date from the payload
// If one is nil, both should be nil, else they should match in value
// This is to be strictly used for dates as it drops any time parameters in the comparison
func (suite *BaseHandlerTestSuite) EqualDatePtr(expected *time.Time, actual *strfmt.Date) {
	if expected == nil || actual == nil {
		suite.Nil(expected)
		suite.Nil(actual)
	} else {
		isoDate := "2006-01-02" // Create a date format
		suite.Equal(expected.Format(isoDate), time.Time(*actual).Format(isoDate))
	}
}
