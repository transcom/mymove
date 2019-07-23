package migrate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitStatements(t *testing.T) {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyFromStdin.sql"
	f, err := os.Open(fixture)
	require.Nil(t, err)

	in := NewBuffer()
	dropComments := true
	dropBlankLines := true
	dropSearchPath := false

	ReadInSQL(f, in, dropComments, dropBlankLines, dropSearchPath)
	formattedSQL := in.String()

	lines := make(chan string, 1000)
	//read buffer values into the channel
	go func() {
		scanner := bufio.NewScanner(strings.NewReader(formattedSQL))
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	statements := make(chan string, 1000)
	go SplitStatements(lines, statements)
	i := 0
	for stmt := range statements {
		fmt.Println("Statement:", i)
		fmt.Println(stmt)
		i++
	}
	require.Equal(t, i, 11)
}