package privileges

import (
	"slices"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models/roles"
	usersprivileges "github.com/transcom/mymove/pkg/services/users_privileges"
)

func (suite *PrivilegesServiceSuite) TestFetchPrivileges() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	priv1 := roles.Privilege{
		ID:            id1,
		PrivilegeType: "priv1",
	}
	id2, _ := uuid.NewV4()
	priv2 := roles.Privilege{
		ID:            id2,
		PrivilegeType: "priv2",
	}
	// Create privileges
	ps := roles.Privileges{priv1, priv2}
	err := suite.DB().Create(ps)
	suite.NoError(err)
	// Associate privileges
	var privilegeTypes []roles.PrivilegeType
	for _, p := range ps {
		privilegeTypes = append(privilegeTypes, p.PrivilegeType)
	}
	upc := usersprivileges.NewUsersPrivilegesCreator()
	_, err = upc.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, privilegeTypes)
	suite.NoError(err)
	// Fetch privileges
	pf := NewPrivilegesFetcher()
	fps, err := pf.FetchPrivilegesForUser(suite.AppContextForTest(), *officeUser.UserID)
	suite.NoError(err)
	suite.Len(fps, 2)
}

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
