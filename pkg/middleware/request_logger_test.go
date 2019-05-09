package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (suite *testSuite) TestRequestLogger() {
	buf := bytes.NewBuffer(make([]byte, 0))
	// Create logger that writes to the buffer instead of stdout/stderr
	logger := suite.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(buf),
			zapcore.DebugLevel,
		))
	}))
	requestLogger := RequestLogger(logger)
	rr := httptest.NewRecorder()
	suite.do(requestLogger, suite.ok, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	out := strings.TrimSpace(string(buf.Bytes()))      // remove trailing new line
	suite.NotEmpty(out, "log was empty")
	lines := strings.Split(out, "\n")
	suite.Len(lines, 1) // there is 1 INFO log line
	parts := strings.Split(lines[0], "\t")
	suite.Len(parts, 4)
	//suite.Equal(parts[0], "") // The Date Time
	suite.Equal(parts[1], "INFO", "log level is invalid")
	suite.Equal(parts[2], "Request", "log message is invalid")
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(parts[3]), &m)
	suite.Nil(err, "log fields are not valid json")
	suite.Contains(m, "method")
	suite.Equal("GET", m["method"])
	suite.Contains(m, "resp-status")
	suite.EqualValues(200, m["resp-status"])
}
