package logging

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func parent() error {
	return next()
}

func next() error {
	return inner()
}

func inner() error {
	return errors.New("error doing child work")
}

func TestFilterErrorFields(t *testing.T) {
	// Create a stacktrace that is over the stacktrace_length of 6
	err := parent()

	fields := []zapcore.Field{zap.String("id", "1234"), zap.Error(err)}
	filteredFields := filterErrorFields(fields, 6)
	// filterErrorFields should split the zap.Error field into 2 fields
	// zap.String fields for error and errorVerbose keys
	require.Len(t, filteredFields, 3)
	require.Equal(t, "error", filteredFields[1].Key)
	require.Equal(t, "error doing child work", filteredFields[1].String)

	// then newlines aren't seperating frames groups but between the function
	// name and filepath with line number so we'll get a length of 3 groups + 1
	// github.com/transcom/mymove/pkg/logging.inner
	// 	/Users/duncan/workspace/mymove/pkg/logging/log_test.go:24 github.com/transcom/mymove/pkg/logging.next
	// 	/Users/duncan/workspace/mymove/pkg/logging/log_test.go:20 github.com/transcom/mymove/pkg/logging.parent
	//  /Users/duncan/workspace/mymove/pkg/logging/log_test.go:16
	require.Len(t, strings.Split(filteredFields[2].String, "\n"), 4)
	require.Equal(t, "errorVerbose", filteredFields[2].Key)
}

func TestFilterErrorFieldsWithNamedError(t *testing.T) {
	// Create a stacktrace that is over the stacktrace_length of 6
	err := parent()
	fields := []zapcore.Field{zap.String("id", "1234"), zap.NamedError("custom", err)}
	filteredFields := filterErrorFields(fields, 6)
	// filterErrorFields should split the zap.Error field into 2 fields
	// zap.String fields for error and errorVerbose keys
	require.Len(t, filteredFields, 3)
	require.Equal(t, "custom", filteredFields[1].Key)
	require.Equal(t, "error doing child work", filteredFields[1].String)

	// then newlines aren't seperating frames groups but between the function
	// name and filepath with line number so we'll get a length of 3 groups + 1
	// github.com/transcom/mymove/pkg/logging.inner
	// 	/Users/duncan/workspace/mymove/pkg/logging/log_test.go:24 github.com/transcom/mymove/pkg/logging.next
	// 	/Users/duncan/workspace/mymove/pkg/logging/log_test.go:20 github.com/transcom/mymove/pkg/logging.parent
	//  /Users/duncan/workspace/mymove/pkg/logging/log_test.go:16
	require.Len(t, strings.Split(filteredFields[2].String, "\n"), 4)
	require.Equal(t, "customVerbose", filteredFields[2].Key)
}