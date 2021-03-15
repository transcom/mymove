//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to clean up file created for unit test
//RA: Given the functions causing the lint errors are used to clean up local storage space after a unit test, it does not present a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
// #nosec G307
package migrate

import (
	"bufio"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSplitStatementsCopyFromStdin(t *testing.T) {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyFromStdin.sql"
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
		"5147b246-19c4-487a-b3fd-a503f889daf7\t2019-01-01\t2019-05-14\t27f1fbeb-090c-4a91-955c-67899de4d6d6\t\\N\t0\t92\t231a7b21-346c-4e94-b6bc-672413733f77\t2018-12-28 18:35:37.147546\t2018-12-28 18:35:37.147546\t2018-10-01\t2019-05-14\t0.67000000000000000000\t0.60000000000000000000",
		"\\.",
	}

	i := 0
	for stmt := range statements {
		require.Equal(t, expectedStmt[i], stmt)
		i++
	}
	require.Equal(t, i, 12)
}

func TestSplitStatementsCopyFromStdinMultiple(t *testing.T) {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyFromStdinMultiple.sql"
	f, err := os.Open(fixture)

	defer f.Close() //lint:ignore SA5001
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
		"COPY public.tariff400ng_full_unpack_rates (id, schedule, rate_millicents, effective_date_lower, effective_date_upper, created_at, updated_at) FROM stdin;",
		"759fc749-2c32-4529-ad8a-131df722cb17\t1\t629370\t2020-05-15\t2021-05-15\t2020-03-23 16:31:50.313853\t2020-03-23 16:31:50.313853",
		"8bd34c5c-ef6e-45b1-9d3d-8a89e46af743\t2\t696150\t2020-05-15\t2021-05-15\t2020-03-23 16:31:50.313853\t2020-03-23 16:31:50.313853",
		"\\.",
		"COPY public.tariff400ng_shorthaul_rates (id, cwt_miles_lower, cwt_miles_upper, rate_cents, effective_date_lower, effective_date_upper, created_at, updated_at) FROM stdin;",
		"55c74996-f208-414d-b8de-022938dbfe1e\t0\t16001\t35511\t2020-05-15\t2021-05-15\t2020-03-23 16:31:50.324166\t2020-03-23 16:31:50.324166",
		"b8e87afb-6287-4837-8b49-a6cf3aad0d1a\t16001\t32001\t31567\t2020-05-15\t2021-05-15\t2020-03-23 16:31:50.324166\t2020-03-23 16:31:50.324166",
		"b36e9835-9794-4465-8c43-b63088c5ebe1\t32001\t64001\t27620\t2020-05-15\t2021-05-15\t2020-03-23 16:31:50.324166\t2020-03-23 16:31:50.324166",
		"\\.",
	}

	i := 0
	for stmt := range statements {
		require.Equal(t, expectedStmt[i], stmt)
		i++
	}
	require.Equal(t, i, 9)
}

func TestSplitStatementsLoop(t *testing.T) {

	// Load the fixture with the sql example
	fixture := "./fixtures/loop.sql"
	f, err := os.Open(fixture)

	// TODO
	defer f.Close() // #nosec G307
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
		"CREATE TABLE IF NOT EXISTS shipments (id serial PRIMARY KEY);",
		"ALTER TABLE invoices ADD COLUMN IF NOT EXISTS shipment_id uuid NULL;",
		"LOCK TABLE invoice_number_trackers, invoices IN SHARE MODE;",
		"DELETE\nFROM invoice_number_trackers;",
		`DO $do$

DECLARE
current_shipment_id     UUID;
current_invoice_id      UUID;
scac                    TEXT;
invoice_count           INT;
shipment_year           INT;
shipment_two_digit_year VARCHAR(2);
base_invoice_number     VARCHAR(255);
new_sequence_number     INT;
target_invoice_number   VARCHAR(255);
BEGIN
FOR current_shipment_id IN SELECT DISTINCT shipment_id, MIN(created_at) as min_created_at
FROM invoices
GROUP BY shipment_id
ORDER BY MIN(created_at)
LOOP
scac := NULL;
SELECT tsp.standard_carrier_alpha_code,
EXTRACT(YEAR FROM s.created_at),
to_char(s.created_at, 'YY')
INTO scac, shipment_year, shipment_two_digit_year
FROM shipments s
INNER JOIN shipment_offers so ON s.id = so.shipment_id
INNER JOIN transportation_service_provider_performances tspp
ON so.transportation_service_provider_performance_id = tspp.id
INNER JOIN transportation_service_providers tsp ON tspp.transportation_service_provider_id = tsp.id
WHERE s.id = current_shipment_id
AND so.accepted = TRUE
ORDER BY so.created_at
LIMIT 1;
IF scac IS NULL THEN
RAISE EXCEPTION 'Shipment ID % has no accepted shipment offer, so unable to generate proper invoice number.', current_shipment_id;
END IF;
invoice_count := 0;
base_invoice_number := NULL;
FOR current_invoice_id IN SELECT id FROM invoices WHERE shipment_id = current_shipment_id ORDER BY created_at
LOOP
IF invoice_count = 0 THEN
INSERT INTO invoice_number_trackers as trackers (standard_carrier_alpha_code, year, sequence_number)
VALUES (scac, shipment_year, 1)
ON CONFLICT (
standard_carrier_alpha_code,
year)
DO UPDATE

SET sequence_number = trackers.sequence_number + 1
WHERE trackers.standard_carrier_alpha_code = scac AND trackers.year = shipment_year
RETURNING sequence_number INTO new_sequence_number;
base_invoice_number := scac || shipment_two_digit_year || to_char(new_sequence_number, 'fm0000');
target_invoice_number := base_invoice_number;
ELSE
target_invoice_number := base_invoice_number || '-' || to_char(invoice_count, 'fm00');
END IF;
UPDATE invoices
SET invoice_number = target_invoice_number,
updated_at     = now()
WHERE id = current_invoice_id;
invoice_count := invoice_count + 1;
END LOOP;
END LOOP;
END $do$;`,
		"ALTER TABLE invoices DROP COLUMN IF EXISTS shipment_id;",
		"DROP TABLE IF EXISTS shipments;",
	}

	i := 0
	for stmt := range statements {
		require.Equal(t, expectedStmt[i], stmt)
		i++
	}
	require.Equal(t, i, 7)
}
