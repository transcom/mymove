package hnyzap

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type zapFieldSuite struct {
	testingsuite.BaseTestSuite
	bool    bool
	float32 float32
	float64 float64
	int32   int32
	int64   int64
	string  string
	uint32  uint32
	uint64  uint64
	error   error
}

func TestZapFieldSuite(t *testing.T) {
	fs := &zapFieldSuite{
		bool:    false,
		float32: rand.Float32(),
		float64: rand.Float64(),
		int32:   rand.Int31(),
		int64:   rand.Int63(),
		uint32:  rand.Uint32(),
		uint64:  rand.Uint64(),
		string:  "zap me",
		error:   fmt.Errorf("fail me"),
	}
	suite.Run(t, fs)
}

func (suite *zapFieldSuite) TestZapBoolean() {
	zapField := zap.Bool("bool", suite.bool)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.bool, honeyField)
}

func (suite *zapFieldSuite) TestZapFloat32() {
	zapField := zap.Float32("float32", suite.float32)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.float32, honeyField)
}

func (suite *zapFieldSuite) TestZapFloat64() {
	zapField := zap.Float64("float64", suite.float64)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.float64, honeyField)
}

func (suite *zapFieldSuite) TestZapInt32() {
	zapField := zap.Int32("int32", suite.int32)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.int32, honeyField)
}

func (suite *zapFieldSuite) TestZapInt64() {
	zapField := zap.Int64("int64", suite.int64)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.int64, honeyField)
}

func (suite *zapFieldSuite) TestZapUint32() {
	zapField := zap.Uint32("unint32", suite.uint32)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.uint32, honeyField)
}

func (suite *zapFieldSuite) TestZapUnint64() {
	zapField := zap.Uint64("unint64", suite.uint64)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.uint64, honeyField)
}

func (suite *zapFieldSuite) TestZapString() {
	zapField := zap.String("string", suite.string)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.string, honeyField)
}

func (suite *zapFieldSuite) TestZapError() {
	zapField := zap.Error(suite.error)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal(suite.error.Error(), honeyField)
}

func (suite *zapFieldSuite) TestZapUnsupported() {
	zapField := zap.Duration("time", time.Second)
	honeyField := ZapFieldToHoneycombField(zapField)
	suite.Equal("unsupported field type", honeyField)
}
