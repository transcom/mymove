package publicapi

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
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	suite.Suite
	db                 *pop.Connection
	logger             *zap.Logger
	filesToClose       []*runtime.File
	notificationSender notifications.NotificationSender
}

// SetupTest is the DB setup
func (suite *HandlerSuite) SetupTest() {
	suite.db.TruncateAll()
}

// mustSave requires saving without errors
func (suite *HandlerSuite) mustSave(model interface{}) {
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

// isNotErrResponse enforces handler does not return an error response
func (suite *HandlerSuite) isNotErrResponse(response middleware.Responder) {
	r, ok := response.(*utils.ErrResponse)
	if ok {
		suite.logger.Error("Received an unexpected error response from handler: ", zap.Error(r.Err))
		// Formally lodge a complaint
		suite.IsType(&utils.ErrResponse{}, response)
	}
}

// checkErrorResponse verifies error response is what is expected
func (suite *HandlerSuite) checkErrorResponse(resp middleware.Responder, code int, name string) {
	errResponse, ok := resp.(*utils.ErrResponse)
	if !ok || errResponse.Code != code {
		suite.T().Fatalf("Expected %s Response: %v", name, resp)
		debug.PrintStack()
	}
}

// checkResponseBadRequest looks at BadRequest errors
func (suite *HandlerSuite) checkResponseBadRequest(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusBadRequest, "BadRequest")
}

// checkResponseUnauthorized looks at Unauthorized errors
func (suite *HandlerSuite) checkResponseUnauthorized(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusUnauthorized, "Unauthorized")
}

// checkResponseForbidden looks at Forbidden errors
func (suite *HandlerSuite) checkResponseForbidden(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusForbidden, "Forbidden")
}

// checkResponseNotFound looks at NotFound errors
func (suite *HandlerSuite) checkResponseNotFound(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusNotFound, "NotFound")
}

// checkResponseInternalServerError looks at InternalServerError errors
func (suite *HandlerSuite) checkResponseInternalServerError(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusInternalServerError, "InternalServerError")
}

// checkResponseTeapot enforces that response come from a Teapot
func (suite *HandlerSuite) checkResponseTeapot(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusTeapot, "Teapot")
}

// authenticateRequest Request authenticated with a service member
func (suite *HandlerSuite) authenticateRequest(req *http.Request, serviceMember models.ServiceMember) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          serviceMember.UserID,
		IDToken:         "fake token",
		ServiceMemberID: serviceMember.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// authenticateUserRequest only authenticated with a bare user - have no idea if they are a service member yet
func (suite *HandlerSuite) authenticateUserRequest(req *http.Request, user models.User) *http.Request {
	session := auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          user.ID,
		IDToken:         "fake token",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// authenticateOfficeRequest authenticates Office users
func (suite *HandlerSuite) authenticateOfficeRequest(req *http.Request, user models.OfficeUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// authenticateTspRequest authenticates TSP users
func (suite *HandlerSuite) authenticateTspRequest(req *http.Request, user models.TspUser) *http.Request {
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		TspUserID:       user.ID,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	return req.WithContext(ctx)
}

// fixture allows us to include a fixture like a PDF in the test
func (suite *HandlerSuite) fixture(name string) *runtime.File {
	fixtureDir := "fixtures"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Error(err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)

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
	suite.closeFile(returnFile)

	return returnFile
}

// AfterTest tries to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Data.Close()
	}
}

func (suite *HandlerSuite) closeFile(file *runtime.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

// TestHandlerSuite creates our test suite
func TestHandlerSuite(t *testing.T) {
	configLocation := "../../../config"
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
		db:                 db,
		logger:             logger,
		notificationSender: notifications.NewStubNotificationSender(logger),
	}

	suite.Run(t, hs)
}
