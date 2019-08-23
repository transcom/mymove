package migrate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomErrors(t *testing.T) {
	r := require.New(t)

	errInvalidDirection := ErrInvalidDirection{Value: "bad"}
	r.Equal("invalid direction \"bad\", expecting up", errInvalidDirection.Error())

	errInvalidFormat := ErrInvalidFormat{Value: "bad"}
	r.Equal("invalid format \"bad\", expecting sql or fizz", errInvalidFormat.Error())

	errInvalidPath := ErrInvalidPath{Value: "bad"}
	r.Equal("invalid path \"bad\", should start with file:// or s3://", errInvalidPath.Error())
}
