package privileges

import (
	"slices"

	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *PrivilegesServiceSuite) TestFetchPrivilegeTypes() {
	// Initialize the privileges fetcher
	pf := NewPrivilegesFetcher()

	// Fetch privilege types
	privTypes, err := pf.FetchPrivilegeTypes(suite.AppContextForTest())

	// Check for errors or empty tables
	suite.NoError(err, "Fetching privilege types should not return an error")
	suite.NotEmpty(privTypes, "Expected privileges to be pre-populated in the database with own privilege_type")

	// Example privilege types to match (adjust as needed)
	privilegesToMatch := []roles.PrivilegeType{
		roles.PrivilegeTypeSupervisor,
		roles.PrivilegeTypeSafety,
	}

	// Assert that only expected privilegeTypes are included in list
	for _, privType := range privTypes {
		index := slices.Index(privilegesToMatch, privType)
		suite.NotEqual(-1, index, "PrivilegeType %s not found in privilegesToMatch.", privType)
		privilegesToMatch = slices.Delete(privilegesToMatch, index, index+1)
	}

	suite.Empty(privilegesToMatch, "privTypes should be 1->1 with privilegesToMatch")
}
