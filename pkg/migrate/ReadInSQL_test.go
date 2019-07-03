package migrate

import (
	"io/ioutil"
	"os"
)

func (suite *MigrateSuite) TestReadInSQLDefault() {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyFromStdin.sql"
	f, err := os.Open(fixture)
	suite.Nil(err)

	in := NewBuffer()
	dropComments := false
	dropBlankLines := false
	dropSearchPath := false

	ReadInSQL(f, in, dropComments, dropBlankLines, dropSearchPath)

	orig, err := ioutil.ReadFile(fixture)
	suite.Nil(err)
	suite.Equal(in.String(), string(orig))
}

func (suite *MigrateSuite) TestReadInSQLStripAll() {

	// Load the fixture with the sql example
	fixture := "./fixtures/copyFromStdin.sql"
	f, err := os.Open(fixture)
	suite.Nil(err)

	in := NewBuffer()
	dropComments := true
	dropBlankLines := true
	dropSearchPath := true

	ReadInSQL(f, in, dropComments, dropBlankLines, dropSearchPath)

	expected := `SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;
COPY public.transportation_service_provider_performances (id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate) FROM stdin;
fbfb095e-6ea3-4c1e-bd3d-7f131d73e295	2019-01-01	2019-05-14	27f1fbeb-090c-4a91-955c-67899de4d6d6	\N	0	89	231a7b21-346c-4e94-b6bc-672413733f77	2018-12-28 18:35:37.147546	2018-12-28 18:35:37.147546	2018-10-01	2019-05-14	0.55000000000000000000	0.55000000000000000000
\.
`
	suite.Equal(in.String(), expected)
}
