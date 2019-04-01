package auth

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func uuidFormatter(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

type dataHelperParams struct {
	isBaseUser   bool
	isOfficeUser bool
	isTspUser    bool
}

var emailCounter int

// Builds related user data depends on which users you want
func dataHelper(db *pop.Connection, params dataHelperParams) goth.User {
	gothUserID := uuid.Must(uuid.NewV4())
	gothEmail := fmt.Sprintf("leo_spaceman_test_%d@example.com", emailCounter)
	emailCounter = emailCounter + 1

	if params.isOfficeUser {
		testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
			OfficeUser: models.OfficeUser{
				Email: gothEmail,
			},
		})
	}

	if params.isTspUser {
		testdatagen.MakeTspUser(db, testdatagen.Assertions{
			TspUser: models.TspUser{
				Email: gothEmail,
			},
		})
	}

	return goth.User{
		UserID: gothUserID.String(),
		Email:  gothEmail,
	}
}

func (suite *UserInitializerSuite) TestUserInitializerBaseUser() {
	// Given a service member user
	user := dataHelper(suite.DB(), dataHelperParams{
		isOfficeUser: false,
		isTspUser:    false,
	})

	initializer := NewUserInitializer(suite.DB())
	identity, err := initializer.InitializeUser(user)
	suite.NoError(err)
	suite.NotEqual(uuid.Nil, identity.ID)
	suite.Nil(identity.OfficeUserID)
	suite.Nil(identity.TspUserID)
}

func (suite *UserInitializerSuite) TestUserInitializerOfficeUser() {
	// Given an office user
	user := dataHelper(suite.DB(), dataHelperParams{
		isOfficeUser: true,
		isTspUser:    false,
	})

	initializer := NewUserInitializer(suite.DB())
	identity, err := initializer.InitializeUser(user)
	suite.NoError(err)
	suite.NotEqual(uuid.Nil, identity.ID)
	suite.NotNil(identity.OfficeUserID)
	suite.Nil(identity.TspUserID)
}

func (suite *UserInitializerSuite) TestUserInitializerTSPUser() {
	// Given a TSP user
	user := dataHelper(suite.DB(), dataHelperParams{
		isOfficeUser: false,
		isTspUser:    true,
	})

	initializer := NewUserInitializer(suite.DB())
	identity, err := initializer.InitializeUser(user)
	suite.NoError(err)
	suite.NotEqual(uuid.Nil, identity.ID)
	suite.Nil(identity.OfficeUserID)
	suite.NotNil(identity.TspUserID)
}

func (suite *UserInitializerSuite) TestUserInitializerBothUser() {
	// Given a combo office/TSP user
	user := dataHelper(suite.DB(), dataHelperParams{
		isOfficeUser: true,
		isTspUser:    true,
	})

	initializer := NewUserInitializer(suite.DB())
	identity, err := initializer.InitializeUser(user)
	suite.NoError(err)
	suite.NotEqual(uuid.Nil, identity.ID)
	suite.NotNil(identity.OfficeUserID)
	suite.NotNil(identity.TspUserID)
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
