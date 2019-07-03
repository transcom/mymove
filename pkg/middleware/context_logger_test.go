package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"

	//"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"go.uber.org/zap/zaptest"
)

func (suite *testSuite) TestContextLoggerWithoutTrace() {
	buf := bytes.NewBuffer(make([]byte, 0))
	// Create logger that writes to the buffer instead of stdout/stderr
	logger := suite.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(buf),
			zapcore.DebugLevel,
		)
	}))
	mw := ContextLogger("", logger)
	rr := httptest.NewRecorder()
	suite.do(mw, suite.log, rr, httptest.NewRequest("GET", testURL, nil))
	suite.Equal(http.StatusOK, rr.Code, errStatusCode) // check status code
	out := strings.TrimSpace(string(buf.Bytes()))      // remove trailing new line
	suite.NotEmpty(out, "log was empty")
	lines := strings.Split(out, "\n")
	suite.Len(lines, 2) // test that there are 2 log lines (info message and error message)
	parts := strings.Split(lines[0], "\t")
	suite.Len(parts, 3)
	suite.Equal(parts[1], "INFO")
	suite.Equal(parts[2], "Placeholder for info message")
	parts = strings.Split(lines[1], "\t")
	suite.Len(parts, 3)
	suite.Equal(parts[1], "ERROR")
	suite.Equal(parts[2], "Placeholder for error message")
}
