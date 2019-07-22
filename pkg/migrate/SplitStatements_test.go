package migrate

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitStatements(t *testing.T) {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyFromStdin.sql"
	f, err := os.Open(fixture)
	require.Nil(t, err)

	lines := make(chan string, 1000)
	go func() {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	statements := make(chan string, 1000)
	go SplitStatements(lines, statements)
	i := 0
	for stmt := range statements {
		fmt.Println("---------------------------------------------------")
		fmt.Println("Statement:", i)
		fmt.Println("---------------------------------------------------")
		fmt.Println(stmt)
		i++
	}
	require.Equal(t, i, 14)
}