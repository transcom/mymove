package internalapi

import (
	"github.com/transcom/mymove/pkg/models"
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	documentServices "github.com/transcom/mymove/pkg/services/document"
	userServices "github.com/transcom/mymove/pkg/services/user"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
)

// HandlerSuite is an abstraction of our original suite
type HandlerSuite struct {
	handlers.BaseTestSuite
}

// SetupTest sets up the test suite by preparing the DB
func (suite *HandlerSuite) SetupTest() {
	suite.TestDB().TruncateAll()
}

// AfterTest completes tests by trying to close open files
func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.TestFilesToClose() {
		file.Data.Close()
	}
}

func (suite *HandlerSuite) HandlerContextWithServices() handlers.HandlerContext {
	h := handlers.NewHandlerContext(suite.TestDB(), suite.TestLogger())
	fetchServiceMember := userServices.NewFetchServiceMemberService(models.NewServiceMemberDB(suite.TestDB()))
	h.SetFetchServiceMember(fetchServiceMember)
	documentDB := models.NewDocumentDB(suite.TestDB())
	fetchDocument := documentServices.NewFetchDocumentService(documentDB, fetchServiceMember)
	h.SetFetchDocument(fetchDocument)
	fetchUpload := documentServices.NewFetchUploadService(documentDB, fetchDocument)
	h.SetFetchUpload(fetchUpload)
	return h
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

	hs := &HandlerSuite{}
	hs.SetTestDB(db)
	hs.SetTestLogger(logger)
	hs.SetTestNotificationSender(notifications.NewStubNotificationSender(logger))

	suite.Run(t, hs)
}
