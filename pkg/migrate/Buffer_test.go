//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package migrate

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {

	in := "hello world"

	buf := NewBuffer()

	go func() {
		time.Sleep(time.Second * 1)
		buf.WriteString(in)
		buf.Close()
	}()

	c := byte(0)
	var err error

	c, err = buf.Index(0)
	require.Equal(t, err, ErrWait)
	require.Equal(t, byte(0), c)

	time.Sleep(time.Second * 2)

	i := 0
	for ; i < len(in); i++ {
		c, err = buf.Index(i)
		require.Nil(t, err)
		require.Equal(t, in[i], c)
	}

	c, err = buf.Index(i)
	require.Equal(t, err, io.EOF)
	require.Equal(t, byte(0), c)

	c, err = buf.Index(i)
	require.Equal(t, err, io.EOF)
	require.Equal(t, byte(0), c)

	cstr, err := buf.Range(1, 0)
	require.NotNil(t, err)
	require.Equal(t, "", cstr)

	cstr, err = buf.Range(0, 1)
	require.Nil(t, err)
	require.Equal(t, "h", cstr)

}

func TestBufferString(t *testing.T) {

	in := "hello world"

	buf := NewBuffer()

	cstr, err := buf.Range(0, 10)
	require.NotNil(t, err)
	require.Equal(t, "", cstr)

	time.Sleep(time.Second * 1)
	buf.WriteString(in)
	buf.Close()

	out := buf.String()
	require.Equal(t, out, in)

	index, err := buf.WriteString(in)
	require.Equal(t, -1, index)
	require.Equal(t, err, ErrClosed)
}

func TestBufferByte(t *testing.T) {

	in := []byte("hello world")

	buf := NewBuffer()

	time.Sleep(time.Second * 1)
	buf.WriteByte(in[0])
	buf.Close()

	out := buf.String()
	require.Equal(t, out, string(in[0]))

	err := buf.WriteByte(in[0])
	require.Equal(t, err, ErrClosed)
}

func TestBufferRune(t *testing.T) {

	in := 'h'

	buf := NewBuffer()

	time.Sleep(time.Second * 1)
	buf.WriteRune(in)
	buf.Close()

	out := buf.String()
	require.Equal(t, out, string(in))

	index, err := buf.WriteRune(in)
	require.Equal(t, -1, index)
	require.Equal(t, err, ErrClosed)
}
