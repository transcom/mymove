package adminapi

import (
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestGenerateQueryFilters() {
	fnZero := func(c string) []services.QueryFilter {
		return []services.QueryFilter{}
	}
	converters := map[string]func(string) []services.QueryFilter{"omega": fnZero}
	suite.Run("error paths", func() {
		suite.Run("params nil", func() {
			qfs := generateQueryFilters(suite.TestLogger(), nil, converters)
			suite.Equal(0, len(qfs))
		})
		suite.Run("params not JSON", func() {
			s := `{"delta": "d",`
			qfs := generateQueryFilters(suite.TestLogger(), &s, converters)
			suite.Equal(0, len(qfs))
		})
	})

	params := `{"alpha": "a", "beta": "b", "gamma": "c"}`
	converters["alpha"] = fnZero
	converters["beta"] = func(c string) []services.QueryFilter {
		return []services.QueryFilter{
			query.NewQueryFilter("double", "=", c+"0"),
			query.NewQueryFilter("double", "=", c+"1"),
		}
	}
	converters["gamma"] = func(c string) []services.QueryFilter {
		return []services.QueryFilter{
			query.NewQueryFilter("triple", "=", c+"0"),
			query.NewQueryFilter("triple", "=", c+"1"),
			query.NewQueryFilter("triple", "=", c+"2"),
		}
	}
	suite.Run("happy path", func() {
		qfs := generateQueryFilters(suite.TestLogger(), &params, converters)
		suite.Equal(5, len(qfs))
	})
}
