package migrate

import (
	"bufio"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSplitStatements(t *testing.T) {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyMultiFromStdin.sql"
	f, err := os.Open(fixture)
	defer f.Close()
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

	wait := 10 * time.Millisecond
	statements := make(chan string, 1000)
	go SplitStatements(lines, statements, wait)

	expectedStmt := []string{
		"SET statement_timeout = 0;",
		"SET lock_timeout = 0;",
		"SET idle_in_transaction_session_timeout = 0;",
		"SET client_encoding = 'UTF8';",
		"SET standard_conforming_strings = on;",
		"SET check_function_bodies = false;",
		"SET client_min_messages = warning;",
		"SET row_security = off;",
		"COPY public.transportation_service_provider_performances (id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate) FROM stdin;",
		"fbfb095e-6ea3-4c1e-bd3d-7f131d73e295\t2019-01-01\t2019-05-14\t27f1fbeb-090c-4a91-955c-67899de4d6d6\t\\N\t0\t89\t231a7b21-346c-4e94-b6bc-672413733f77\t2018-12-28 18:35:37.147546\t2018-12-28 18:35:37.147546\t2018-10-01\t2019-05-14\t0.55000000000000000000\t0.55000000000000000000",
		"tbqb095e-6ea3-4c1e-bd3d-7f131d343295\t2019-02-01\t2019-05-14\t27f1fbeb-090c-4a91-955c-67899de4d6d6\t\\N\t0\t89\t231a7b21-346c-4e94-b6bc-672413733f77\t2018-12-28 18:35:37.147546\t2018-12-28 18:35:37.147546\t2018-10-01\t2019-05-14\t0.55000000000000000000\t0.55000000000000000000",
	}

	i := 0
	for stmt := range statements {
		require.Equal(t, expectedStmt[i], stmt)
		i++
	}
	require.Equal(t, i, 11)
}
