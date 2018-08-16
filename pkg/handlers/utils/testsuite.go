package utils

import (
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	suite.Suite
	Db                 *pop.Connection
	Logger             *zap.Logger
	FilesToClose       []*runtime.File
	NotificationSender notifications.NotificationSender
}

// SetupTest is the DB setup
func (suite *HandlerSuite) SetupTest() {
	suite.Db.TruncateAll()
}

// MustSave requires saving without errors
func (suite *HandlerSuite) MustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.Db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

// IsNotErrResponse enforces handler does not return an error response
func (suite *HandlerSuite) IsNotErrResponse(response middleware.Responder) {
	r, ok := response.(*ErrResponse)
	if ok {
		suite.Logger.Error("Received an unexpected error response from handler: ", zap.Error(r.Err))
		// Formally lodge a complaint
		suite.IsType(&ErrResponse{}, response)
	}
}

// CheckErrorResponse verifies error response is what is expected
func (suite *HandlerSuite) CheckErrorResponse(resp middleware.Responder, code int, name string) {
	errResponse, ok := resp.(*ErrResponse)
	if !ok || errResponse.Code != code {
		suite.T().Fatalf("Expected %s Response: %v", name, resp)
		debug.PrintStack()
	}
}

// CheckResponseBadRequest looks at BadRequest errors
func (suite *HandlerSuite) CheckResponseBadRequest(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusBadRequest, "BadRequest")
}

// CheckResponseUnauthorized looks at Unauthorized errors
func (suite *HandlerSuite) CheckResponseUnauthorized(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusUnauthorized, "Unauthorized")
}

// CheckResponseForbidden looks at Forbidden errors
func (suite *HandlerSuite) CheckResponseForbidden(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusForbidden, "Forbidden")
}

// CheckResponseNotFound looks at NotFound errors
func (suite *HandlerSuite) CheckResponseNotFound(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusNotFound, "NotFound")
}

// CheckResponseInternalServerError looks at InternalServerError errors
func (suite *HandlerSuite) CheckResponseInternalServerError(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusInternalServerError, "InternalServerError")
}

// CheckResponseTeapot enforces that response come from a Teapot
func (suite *HandlerSuite) CheckResponseTeapot(resp middleware.Responder) {
	suite.CheckErrorResponse(resp, http.StatusTeapot, "Teapot")
}

// AuthenticateRequest Request authenticated with a service member
func (suite *HandlerSuite) AuthenticateRequest(req *http.Request, serviceMember models.ServiceMember) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          serviceMember.UserID,
		IDToken:         "fake token",
		ServiceMemberID: serviceMember.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateUserRequest only authenticated with a bare user - have no idea if they are a service member yet
func (suite *HandlerSuite) AuthenticateUserRequest(req *http.Request, user models.User) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateOfficeRequest authenticates Office users
func (suite *HandlerSuite) AuthenticateOfficeRequest(req *http.Request, user models.OfficeUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// AuthenticateTspRequest authenticates TSP users
func (suite *HandlerSuite) AuthenticateTspRequest(req *http.Request, user models.TspUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		TspUserID:       user.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// Fixture allows us to include a fixture like a PDF in the test
func (suite *HandlerSuite) Fixture(name string) *runtime.File {
	fixtureDir := "fixtures"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Error(err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)

	// #nosec never comes from user input
	file, err := os.Open(fixturePath)
	if err != nil {
		suite.Logger.Fatal("Error opening fixture file", zap.Error(err))
	}

	info, err := file.Stat()
	if err != nil {
		suite.Logger.Fatal("Error accessing fixture stats", zap.Error(err))
	}

	header := multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	returnFile := &runtime.File{
		Header: &header,
		Data:   file,
	}
	suite.closeFile(returnFile)

	return returnFile
}

// AfterTest tries to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.FilesToClose {
		file.Data.Close()
	}
}

func (suite *HandlerSuite) closeFile(file *runtime.File) {
	suite.FilesToClose = append(suite.FilesToClose, file)
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HandlerSuite{
		Db:                 db,
		Logger:             logger,
		NotificationSender: notifications.NewStubNotificationSender(logger),
	}

	suite.Run(t, hs)
}
