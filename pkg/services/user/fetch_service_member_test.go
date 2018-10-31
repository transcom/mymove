package user

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchServiceMemberForUser() {
	user1 := testdatagen.MakeDefaultUser(suite.db)
	user2 := testdatagen.MakeDefaultUser(suite.db)

	firstName := "Oliver"
	resAddress := testdatagen.MakeDefaultAddress(suite.db)
	sm := ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.mustSave(&sm)

	// User is authorized to fetch service member
	session := &server.Session{
		ApplicationName: auth.MyApp,
		UserID:          user1.ID,
		ServiceMemberID: sm.ID,
	}
	goodSm, err := FetchServiceMemberForUser(suite.db, session, sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddress.ID, goodSm.ResidentialAddress.ID)
	}

	// Wrong ServiceMember
	wrongID, _ := uuid.NewV4()
	_, err = FetchServiceMemberForUser(suite.db, session, wrongID)
	if suite.Error(err) {
		suite.Equal(ErrFetchNotFound, err)
	}

	// User is forbidden from fetching order
	session.UserID = user2.ID
	session.ServiceMemberID = uuid.Nil
	_, err = FetchServiceMemberForUser(suite.db, session, sm.ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}
}

func (suite *ModelSuite) TestFetchServiceMemberNotForUser() {
	user1 := testdatagen.MakeDefaultUser(suite.db)

	firstName := "Nino"
	resAddress := testdatagen.MakeDefaultAddress(suite.db)
	sm := ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.mustSave(&sm)

	goodSm, err := FetchServiceMember(suite.db, sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddressID, goodSm.ResidentialAddressID)
	}
}
