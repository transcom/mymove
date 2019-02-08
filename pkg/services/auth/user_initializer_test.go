package auth

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func uuidFormatter(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

// Mocks a session object, as far as the UserInitializer is concerned
type appDetectorMock struct {
	isOfficeApp bool
	isTspApp    bool
}

func (o appDetectorMock) IsOfficeApp() bool {
	return o.isOfficeApp
}

func (o appDetectorMock) IsTspApp() bool {
	return o.isTspApp
}

type dataHelperParams struct {
	isBaseUser   bool
	isOfficeUser bool
	isTspUser    bool
}

var emailCounter int

// Builds related user data depends on which users you want
func dataHelper(db *pop.Connection, params dataHelperParams) goth.User {
	var user models.User
	if params.isBaseUser {
		user = testdatagen.MakeDefaultUser(db)
	}

	gothUserID := user.ID
	gothEmail := user.LoginGovEmail
	if gothUserID == uuid.Nil {
		gothUserID = uuid.Must(uuid.NewV4())
		gothEmail = fmt.Sprintf("leo_spaceman_test_%d@example.com", emailCounter)
		emailCounter = emailCounter + 1
	}

	if params.isOfficeUser {
		testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				User:   user,
				UserID: uuidFormatter(user.ID),
				Email:  gothEmail,
			},
		})
	}

	if params.isTspUser {
		testdatagen.MakeTspUser(db, testdatagen.Assertions{
			TspUser: models.TspUser{
				User:   user,
				UserID: uuidFormatter(user.ID),
				Email:  gothEmail,
			},
		})
	}

	return goth.User{
		UserID: gothUserID.String(),
		Email:  gothEmail,
	}
}

// Verifies office user data exists and NOT tsp user data
func (suite *UserInitializerSuite) verifyOfficeResponse(r services.InitializeUserResponse) bool {
	return suite.NotEqual(uuid.Nil, r.UserID) &&
		suite.NotEqual(uuid.Nil, r.OfficeUserID) &&
		suite.Equal(uuid.Nil, r.TspUserID)
}

// Verifies tsp user data exists and NOT office user data
func (suite *UserInitializerSuite) verifyTspResponse(r services.InitializeUserResponse) bool {
	return suite.NotEqual(uuid.Nil, r.UserID) &&
		suite.Equal(uuid.Nil, r.OfficeUserID) &&
		suite.NotEqual(uuid.Nil, r.TspUserID)
}

// Verifies that NEITHER office or TSP user data is provided
func (suite *UserInitializerSuite) verifyMilmoveResponse(r services.InitializeUserResponse) bool {
	return suite.NotEqual(uuid.Nil, r.UserID) &&
		suite.Equal(uuid.Nil, r.OfficeUserID) &&
		suite.Equal(uuid.Nil, r.TspUserID)
}

func (suite *UserInitializerSuite) TestOfficeAppOfficeUser() {
	// On the office app
	detector := appDetectorMock{
		isOfficeApp: true,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// A previously authed Office user should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: true,
		isTspUser:    false,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyOfficeResponse(response)
}

func (suite *UserInitializerSuite) TestOfficeAppNewOfficeUser() {
	// On the office app
	detector := appDetectorMock{
		isOfficeApp: true,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// A brand new office user should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   false,
		isOfficeUser: true,
		isTspUser:    false,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyOfficeResponse(response)
}

func (suite *UserInitializerSuite) TestOfficeAppNotOfficeUser() {
	// On the office app
	detector := appDetectorMock{
		isOfficeApp: true,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// A base-only user should fail
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: false,
		isTspUser:    false,
	})
	_, verrs, err := initializer.InitializeUser(detector, user)
	suite.Error(err)
	suite.False(verrs.HasAny())
}

func (suite *UserInitializerSuite) TestOfficeAppTspUser() {
	// On the office app
	detector := appDetectorMock{
		isOfficeApp: true,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// A base-only user should fail
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: false,
		isTspUser:    true,
	})
	_, verrs, err := initializer.InitializeUser(detector, user)
	suite.Error(err)
	suite.False(verrs.HasAny())
}

func (suite *UserInitializerSuite) TestTspAppTspUser() {
	// On the TSP app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    true,
	}

	initializer := NewUserInitializer(suite.DB())

	// A previously authed TSP user should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: false,
		isTspUser:    true,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyTspResponse(response)
}

func (suite *UserInitializerSuite) TestTspAppNewTspUser() {
	// On the TSP app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    true,
	}

	initializer := NewUserInitializer(suite.DB())

	// A brand new TSP user should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   false,
		isOfficeUser: false,
		isTspUser:    true,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyTspResponse(response)
}

func (suite *UserInitializerSuite) TestTspAppNotTspUser() {
	// On the TSP app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    true,
	}

	initializer := NewUserInitializer(suite.DB())

	// A base-only user should fail
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: false,
		isTspUser:    false,
	})
	_, verrs, err := initializer.InitializeUser(detector, user)
	suite.Error(err)
	suite.False(verrs.HasAny())
}

func (suite *UserInitializerSuite) TestTspAppOfficeUser() {
	// On the TSP app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    true,
	}

	initializer := NewUserInitializer(suite.DB())

	// An office user should fail
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: true,
		isTspUser:    false,
	})
	_, verrs, err := initializer.InitializeUser(detector, user)
	suite.Error(err)
	suite.False(verrs.HasAny())
}

func (suite *UserInitializerSuite) TestMilmoveAppOfficeUser() {
	// On the Milmove app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// An office user should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: true,
		isTspUser:    false,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyMilmoveResponse(response)
}

func (suite *UserInitializerSuite) TestMilmoveAppTspUser() {
	// On the Milmove app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// A TSP user should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: false,
		isTspUser:    true,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyMilmoveResponse(response)
}

func (suite *UserInitializerSuite) TestMilmoveAppAllUsers() {
	// On the Milmove app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// A user with all roles should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: true,
		isTspUser:    true,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyMilmoveResponse(response)
}

func (suite *UserInitializerSuite) TestMilmoveAppBaseUser() {
	// On the Milmove app
	detector := appDetectorMock{
		isOfficeApp: false,
		isTspApp:    false,
	}

	initializer := NewUserInitializer(suite.DB())

	// A base-only user should succeed
	user := dataHelper(suite.DB(), dataHelperParams{
		isBaseUser:   true,
		isOfficeUser: false,
		isTspUser:    false,
	})
	response, verrs, err := initializer.InitializeUser(detector, user)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.verifyMilmoveResponse(response)
}

type UserInitializerSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *UserInitializerSuite) SetupTest() {
	suite.DB().TruncateAll()
}
func TestUserInitializerSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &UserInitializerSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
