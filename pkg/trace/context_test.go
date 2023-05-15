package trace

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type traceSuite struct {
	testingsuite.BaseTestSuite
}

func TestTraceSuite(t *testing.T) {
	ss := &traceSuite{}
	suite.Run(t, ss)
}

func (suite *traceSuite) TestTraceContextRoundtrip() {
	ctx := context.Background()
	uuid := uuid.Must(uuid.NewV4())
	traceCtx := NewContext(ctx, uuid)
	suite.Equal(uuid, FromContext(traceCtx))
}

func (suite *traceSuite) TestTraceAwsXrayContextRoundtrip() {
	ctx := context.Background()
	xrayID := "1-xray-id"
	traceCtx := AwsXrayNewContext(ctx, xrayID)
	suite.Equal(xrayID, AwsXrayFromContext(traceCtx))
}
