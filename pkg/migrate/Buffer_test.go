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

}
