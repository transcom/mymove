package handlers

import (
	"log"
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HandlerSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *HandlerSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestHandlerSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}
	Init(db)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &HandlerSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
