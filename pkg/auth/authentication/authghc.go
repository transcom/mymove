package authentication

import (
	"strings"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// ErrTOOUnauthorized is too unauthorized error
var ErrTOOUnauthorized = errors.New("too unauthorized user")

// ErrUnauthorized is unauthorized user
var ErrUnauthorized = errors.New("unauthorized user")

// ErrUserDeactivated is user deactivated error
var ErrUserDeactivated = errors.New("user is deactivated")

// UserCreator creates users
//go:generate mockery -name UserCreator
type UserCreator interface {
	CreateUser(id string, email string) (*models.User, error)
}

// RoleAssociator associates roles to users
//go:generate mockery -name RoleAssociator
type RoleAssociator interface {
	AdminUserAssociator
	OfficeUserAssociator
	TOORoleChecker
}

// OfficeUserAssociator interface
//go:generate mockery -name OfficeUserAssociator
type OfficeUserAssociator interface {
	FetchOfficeUser(email string) (*models.OfficeUser, error)
	AssociateOfficeUser(user *models.User) (uuid.UUID, error)
}

// AdminUserAssociator interface
//go:generate mockery -name AdminUserAssociator
type AdminUserAssociator interface {
	FetchAdminUser(email string) (*models.AdminUser, error)
	AssociateAdminUser(user *models.User) (uuid.UUID, error)
}

// TOORoleChecker checks TOO roles
//go:generate mockery -name TOORoleChecker
type TOORoleChecker interface {
	FetchUserIdentity(user *models.User) (*models.UserIdentity, error)
	VerifyHasTOORole(identity *models.UserIdentity) (roles.Role, error)
}

// UnknownUserAuthorizer is an unknown user authorizer
type UnknownUserAuthorizer struct {
	logger Logger
	UserCreator
	RoleAssociator
}

// NewUnknownUserAuthorizer returns a new unknown user authorizer
func NewUnknownUserAuthorizer(db *pop.Connection, logger Logger) *UnknownUserAuthorizer {
	uc := userCreator{db}
	oa := officeUserAssociator{db, logger}
	aa := adminUserAssociator{db, logger}
	ta := tooRoleChecker{db, logger}
	ra := roleAssociator{
		db:                   db,
		logger:               logger,
		OfficeUserAssociator: oa,
		AdminUserAssociator:  aa,
		TOORoleChecker:       ta,
	}
	return &UnknownUserAuthorizer{
		logger:         logger,
		UserCreator:    uc,
		RoleAssociator: ra,
	}
}

// AuthorizeUnknownUser will authorize an unknown user
func (uua UnknownUserAuthorizer) AuthorizeUnknownUser(openIDUser goth.User, session *auth.Session) error {
	user, err := uua.CreateUser(openIDUser.UserID, openIDUser.Email)
	if err != nil {
		uua.logger.Error("Error creating user", zap.Error(err))
		return err
	}
	session.UserID = user.ID
	if session.IsAdminApp() {
		session.AdminUserID, err = uua.AssociateAdminUser(user)
		if err != nil {
			return err
		}
	}
	if session.IsOfficeApp() {
		session.OfficeUserID, err = uua.AssociateOfficeUser(user)
		if err != nil {
			switch err {
			case ErrUnauthorized:
				var tooErr error
				userIdentity, tooErr := uua.FetchUserIdentity(user)
				if tooErr != nil {
					return tooErr
				}
				tooRole, tooErr := uua.VerifyHasTOORole(userIdentity)
				if tooErr == nil && !session.Roles.HasRole(roles.RoleTypeTOO) {
					session.Roles = append(session.Roles, tooRole)
				}
				if tooErr != nil {
					return tooErr
				}
			default:
				return err
			}
		}
	}

	return nil
}

type roleAssociator struct {
	db     *pop.Connection
	logger Logger
	OfficeUserAssociator
	AdminUserAssociator
	TOORoleChecker
}

type officeUserAssociator struct {
	db     *pop.Connection
	logger Logger
}

// AssociatedOfficeUser associates an office user
func (oua officeUserAssociator) AssociateOfficeUser(user *models.User) (uuid.UUID, error) {
	officeUser, err := oua.FetchOfficeUser(user.LoginGovEmail)
	if err == models.ErrFetchNotFound {
		oua.logger.Error("No Office user found", zap.String("email", user.LoginGovEmail))
		return uuid.UUID{}, ErrUnauthorized
	}
	if err != nil {
		oua.logger.Error("Checking for office user", zap.String("email", user.LoginGovEmail), zap.Error(err))
		return uuid.UUID{}, err
	}
	if !officeUser.Active {
		oua.logger.Error("Office user is deactivated", zap.String("email", user.LoginGovEmail))
		return uuid.UUID{}, ErrUserDeactivated
	}
	if officeUser.ID != uuid.Nil && officeUser.UserID == nil {
		officeUser.UserID = &user.ID
		err = oua.db.UpdateColumns(officeUser, "user_id")
		if err != nil {
			oua.logger.Error("Error creating user", zap.Error(err))
			return uuid.UUID{}, err
		}
	}
	return officeUser.ID, nil
}

// FetchOfficeUser fetches an office user
func (oua officeUserAssociator) FetchOfficeUser(email string) (*models.OfficeUser, error) {
	officeUser, err := models.FetchOfficeUserByEmail(oua.db, email)
	return officeUser, err
}

type adminUserAssociator struct {
	db     *pop.Connection
	logger Logger
}

// AssociateAdminuser associates an admin user
func (aua adminUserAssociator) AssociateAdminUser(user *models.User) (uuid.UUID, error) {
	adminUser, err := aua.FetchAdminUser(user.LoginGovEmail)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			aua.logger.Error("no admin user found", zap.String("email", user.LoginGovEmail))
			return uuid.UUID{}, ErrUnauthorized
		}
		aua.logger.Error("checking for admin user", zap.String("email", user.LoginGovEmail), zap.Error(err))
		return uuid.UUID{}, err
	}
	if !adminUser.Active {
		aua.logger.Error("admin user is deactivated", zap.String("email", user.LoginGovEmail))
		return uuid.UUID{}, ErrUserDeactivated
	}
	if adminUser.ID != uuid.Nil && adminUser.UserID == nil {
		adminUser.UserID = &user.ID
		err = aua.db.UpdateColumns(adminUser, "user_id")
		if err != nil {
			aua.logger.Error("error creating user", zap.Error(err))
			return uuid.UUID{}, err
		}
	}
	return adminUser.ID, nil
}

// FetchAdminUser fetches an admin user
func (aua adminUserAssociator) FetchAdminUser(email string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := aua.db.Where("LOWER(email) = $1", strings.ToLower(email)).First(&adminUser)
	return &adminUser, err
}

type userCreator struct {
	db *pop.Connection
}

// CreateUser creates a user
func (uc userCreator) CreateUser(id string, email string) (*models.User, error) {
	return models.CreateUser(uc.db, id, email)
}

type tooRoleChecker struct {
	db     *pop.Connection
	logger Logger
}

// FetchUserIdentity fetches a user identity
func (t tooRoleChecker) FetchUserIdentity(user *models.User) (*models.UserIdentity, error) {
	return models.FetchUserIdentity(t.db, user.LoginGovUUID.String())
}

// VerifyHasTOORole verifies user has TOO Role
// Probably want to update this to return roles to add to session
func (t tooRoleChecker) VerifyHasTOORole(identity *models.UserIdentity) (roles.Role, error) {
	if role, ok := identity.Roles.GetRole(roles.RoleTypeTOO); ok {
		return role, nil
	}
	return roles.Role{}, ErrTOOUnauthorized
}
