package testingsuite

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/trace"
)

// AppContextTestHelper keeps track of context per test
type AppContextTestHelper struct {
	testContexts map[string]context.Context
}

// NewAppContextTestHelper creates a test helper
func NewAppContextTestHelper() AppContextTestHelper {
	return AppContextTestHelper{
		testContexts: make(map[string]context.Context),
	}
}

// CurrentTestContext finds or creates a test context for the current test
func (h AppContextTestHelper) CurrentTestContext(testName string) context.Context {
	ctx, ok := h.testContexts[testName]
	if ok {
		return ctx
	}
	// ensure all test contexts have a traceId
	h.testContexts[testName] = trace.NewContext(context.Background(), uuid.Must(uuid.NewV4()))
	return h.testContexts[testName]
}

// UpdateCurrentTestContext overwrites the current test context with
// the provided one
func (h AppContextTestHelper) UpdateCurrentTestContext(ctx context.Context, testName string) {
	h.testContexts[testName] = ctx
}
