package migrate

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestReadInSQLLine(t *testing.T) {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyFromStdin.sql"
	f, err := os.Open(fixture)

	defer func() {
		if fixtureCloseErr := f.Close(); fixtureCloseErr != nil {
			t.Error("Failed to close fixture", zap.Error(fixtureCloseErr))
		}
	}()
	require.Nil(t, err)

	lines := make(chan string, 1000)
	dropComments := true
	dropSearchPath := true
	go func() {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines <- ReadInSQLLine(scanner.Text(), dropComments, dropSearchPath)
		}
		close(lines)
	}()

	expectedStmt := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"SET statement_timeout = 0;",
		"SET lock_timeout = 0;",
		"SET idle_in_transaction_session_timeout = 0;",
		"SET client_encoding = 'UTF8';",
		"SET standard_conforming_strings = on;",
		"",
		"SET check_function_bodies = false;",
		"SET client_min_messages = warning;",
		"SET row_security = off;",
		"",
		"",
		"",
		"",
		"",
		"COPY public.re_services (id, code, name, created_at, updated_at, priority) FROM stdin;",
		"10000012-2c32-4529-ad8a-131df722cb17\t12\tTwelve\t2020-03-23 16:31:50.313853\t2020-03-23 16:31:50.313853\t1",
		"10000013-ef6e-45b1-9d3d-8a89e46af743\t13\tThirteen\t2020-03-23 16:31:50.313853\t2020-03-23 16:31:50.313853\t2",
		"\\.",
		"",
		"",
		"",
		"",
		"",
		"",
	}

	i := 0
	for stmt := range lines {
		require.Equal(t, expectedStmt[i], stmt)
		i++
	}
	require.Equal(t, i, 30)
}
