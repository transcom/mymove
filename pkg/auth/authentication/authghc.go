package authentication

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

var ErrUnauthorized = errors.New("unauthorized user")
var ErrUserDeactivated = errors.New("user is deactivated")

//go:generate mockery -name UserCreator
type UserCreator interface {
	CreateUser(id string, email string) (*models.User, error)
}

//go:generate mockery -name RoleAssociator
type RoleAssociator interface {
	AdminUserAssociator
	OfficeUserAssociator
	CustomerCreatorAndAssociator
}

//go:generate mockery -name CustomerCreatorAndAssociator
type CustomerCreatorAndAssociator interface {
	CreateAndAssociateCustomer(userID uuid.UUID) error
}

//go:generate mockery -name OfficeUserAssociator
type OfficeUserAssociator interface {
	FetchOfficeUser(email string) (*models.OfficeUser, error)
	AssociateOfficeUser(user *models.User) (uuid.UUID, error)
}

//go:generate mockery -name AdminUserAssociator
type AdminUserAssociator interface {
	FetchAdminUser(email string) (*models.AdminUser, error)
	AssociateAdminUser(user *models.User) (uuid.UUID, error)
}

type UnknownUserAuthorizer struct {
	logger Logger
	UserCreator
	RoleAssociator
}

func NewUnknownUserAuthorizer(db *pop.Connection, logger Logger) *UnknownUserAuthorizer {
	uc := userCreator{db}
	oa := officeUserAssociator{db, logger}
	ca := customerAssociator{db, logger}
	aa := adminUserAssociator{db, logger}
	ra := roleAssociator{
		db:                           db,
		logger:                       logger,
		OfficeUserAssociator:         oa,
		AdminUserAssociator:          aa,
		CustomerCreatorAndAssociator: ca,
	}
	return &UnknownUserAuthorizer{
		logger:         logger,
		UserCreator:    uc,
		RoleAssociator: ra,
	}
}

func (uua UnknownUserAuthorizer) AuthorizeUnknownUser(openIDUser goth.User, session *auth.Session) error {
	user, err := uua.CreateUser(openIDUser.UserID, openIDUser.Email)
	if err != nil {
		uua.logger.Error("Error creating user", zap.Error(err))
		return err
	}
	if session.IsAdminApp() {
		session.AdminUserID, err = uua.AssociateAdminUser(user)
		if err != nil {
			return err
		}
	}
	if session.IsOfficeApp() {
		session.OfficeUserID, err = uua.AssociateOfficeUser(user)
		if err != nil {
			return err
		}
	}
	if session.IsMilApp() {
		err = uua.CreateAndAssociateCustomer(user.ID)
		if err != nil {
			return err
		}
	}
	uua.logger.Info("logged in", zap.Any("session", session))
	return nil
}

type roleAssociator struct {
	db     *pop.Connection
	logger Logger
	OfficeUserAssociator
	AdminUserAssociator
	CustomerCreatorAndAssociator
}

type officeUserAssociator struct {
	db     *pop.Connection
	logger Logger
}

//TODO make idempotent
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
	if officeUser.ID != uuid.Nil {
		officeUser.UserID = &user.ID
		err = oua.db.UpdateColumns(officeUser, "user_id")
		if err != nil {
			oua.logger.Error("Error creating user", zap.Error(err))
			return uuid.UUID{}, err
		}
	}
	return officeUser.ID, nil
}

func (oua officeUserAssociator) FetchOfficeUser(email string) (*models.OfficeUser, error) {
	officeUser, err := models.FetchOfficeUserByEmail(oua.db, email)
	return officeUser, err
}

type adminUserAssociator struct {
	db     *pop.Connection
	logger Logger
}

//TODO make idempotent
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
	if adminUser.ID != uuid.Nil && adminUser.UserID != nil {
		adminUser.UserID = &user.ID
		err = aua.db.UpdateColumns(adminUser, "user_id")
		if err != nil {
			aua.logger.Error("error creating user", zap.Error(err))
			return uuid.UUID{}, err
		}
	}
	return adminUser.ID, nil
}

func (aua adminUserAssociator) FetchAdminUser(email string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	queryBuilder := query.NewQueryBuilder(aua.db)
	filters := []services.QueryFilter{
		query.NewQueryFilter("email", "=", email),
	}
	err := queryBuilder.FetchOne(&adminUser, filters)
	return &adminUser, err
}

type customerAssociator struct {
	db     *pop.Connection
	logger Logger
}

//TODO make idempotent
func (ca customerAssociator) CreateAndAssociateCustomer(userID uuid.UUID) error {
	if userID == uuid.Nil {
		ca.logger.Error("error creating customer, user id cannot be nil")
		return errors.New("user id is nil")
	}
	customer := models.Customer{}
	customer.UserID = userID
	err := ca.db.Create(&customer)
	if err != nil {
		ca.logger.Error("error creating customer", zap.Error(err))
		return err
	}
	return nil
}

type userCreator struct {
	db *pop.Connection
}

func (uc userCreator) CreateUser(id string, email string) (*models.User, error) {
	return models.CreateUser(uc.db, id, email)
}
