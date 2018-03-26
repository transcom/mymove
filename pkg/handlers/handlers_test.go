package handlers

import (
	"log"
	"mime/multipart"
	"os"
	"path"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HandlerSuite struct {
	suite.Suite
	db           *pop.Connection
	logger       *zap.Logger
	filesToClose []*os.File
}

func (suite *HandlerSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *HandlerSuite) mustSave(model interface{}) {
	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		suite.T().Fatalf("errors encountered saving %v: %v", model, verrs)
	}
}

func (suite *HandlerSuite) fixture(name string) *runtime.File {
	fixtureDir := "fixtures"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Fatal(err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)

	info, err := os.Stat(fixturePath)
	if err != nil {
		suite.T().Fatal(err)
	}
	header := multipart.FileHeader{
		Filename: name,
		Size:     info.Size(),
	}
	data, err := os.Open(fixturePath)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.closeFile(data)
	return &runtime.File{
		Header: &header,
		Data:   data,
	}
}

func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Close()
	}
}

func (suite *HandlerSuite) closeFile(file *os.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

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
	hs := &HandlerSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
