package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ClientLogCollectorSuite struct {
	*testingsuite.PopTestSuite
}

func TestClientLogCollectorSuite(t *testing.T) {
	hs := &ClientLogCollectorSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *ClientLogCollectorSuite) TestClientLogHandler() {
	buf := bytes.NewBuffer(make([]byte, 0))
	// Create logger that writes to the buffer instead of stdout/stderr
	logger := suite.Logger().WithOptions(zap.WrapCore(func(_ zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(buf),
			zapcore.DebugLevel,
		)
	}))
	appCtx := appcontext.NewAppContext(suite.DB(), logger, nil, nil)
	handler := NewClientLogHandler(appCtx)

	logUpload := ClientLogUpload{
		App:         "test",
		LoggerStats: ClientLoggerStats{},
		LogEntries: []ClientLogEntry{
			{
				Level: "info",
				Args:  []interface{}{"one_arg", 2},
			},
		},
	}
	data, err := json.Marshal(&logUpload)
	suite.NoError(err)
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", "/", body)
	suite.NoError(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	out := strings.TrimSpace(buf.String()) // remove trailing new line
	lines := strings.Split(out, "\n")
	// test that there are 2 log lines (stats message and entry)
	suite.Len(lines, 2, lines)
	suite.Contains(lines[0], "logEntryCount")
	suite.Contains(lines[1], "logLevel")
}
