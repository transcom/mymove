package user

import (
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type smTestSuite struct {
	suite.Suite
}

func TestModelSuite(t *testing.T) {
	s := &smTestSuite{}
	suite.Run(t, s)
}

type mockServiceMemberDB struct {
	sm *models.ServiceMember
}

func (msmDB mockServiceMemberDB) Save(sm *models.ServiceMember) (models.ValidationErrors, error) {
	return nil, nil
}

func (msmDB mockServiceMemberDB) Fetch(id uuid.UUID, loadEagerly bool) (*models.ServiceMember, error) {
	if msmDB.sm.ID != id {
		return nil, models.ErrFetchNotFound
	}
	return msmDB.sm, nil
}

func (msmDB mockServiceMemberDB) IsTspManagingShipment(tspUserID uuid.UUID, smUserID uuid.UUID) (bool, error) {
	return false, nil
}

func (s *smTestSuite) TestFetchServiceMemberForUser() {
	user1 := testdatagen.GetDefaultUser()
	user2 := testdatagen.GetDefaultUser()

	firstName := "Oliver"
	resAddress := testdatagen.GetDefaultAddress()
	sm := models.ServiceMember{
		User:                 *user1,
		UserID:               user1.ID,
		ID:                   uuid.Must(uuid.NewV4()),
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   resAddress,
	}

	smDB := mockServiceMemberDB{&sm}

	// User is authorized to fetch service member
	session := &server.Session{
		ApplicationName: server.MyApp,
		UserID:          user1.ID,
		ServiceMemberID: sm.ID,
	}
	fetchServiceMemberService := NewFetchServiceMemberService(smDB)

	goodSm, err := fetchServiceMemberService.Execute(sm.ID, session)
	if s.NoError(err) {
		s.Equal(sm.FirstName, goodSm.FirstName)
		s.Equal(sm.ResidentialAddress.ID, goodSm.ResidentialAddress.ID)
	}

	// Wrong ServiceMember
	wrongID, _ := uuid.NewV4()
	_, err = fetchServiceMemberService.Execute(wrongID, session)
	if s.Error(err) {
		s.Equal(services.ErrFetchForbidden, err)
	}

	// User is forbidden from fetching ServiceMember
	session.UserID = user2.ID
	session.ServiceMemberID = uuid.Nil
	_, err = fetchServiceMemberService.Execute(sm.ID, session)
	if s.Error(err) {
		s.Equal(services.ErrFetchForbidden, err)
	}

	// No session, no checks
	noSessionSm, err := fetchServiceMemberService.Execute(sm.ID, nil)
	if s.NoError(err) {
		s.Equal(noSessionSm.ID, sm.ID)
	}
}
