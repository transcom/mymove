package user

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userDeleter struct {
	builder userQueryBuilder
}

func (o *userDeleter) DeleteUser(appCtx appcontext.AppContext, id uuid.UUID) error {
	// need to fetch the user and any downstream associations (roles, privileges)
	var user models.User
	err := appCtx.DB().EagerPreload(
		"Roles",
		"Privileges",
	).Where("id = ?", id).Find(&user, id)
	if err == sql.ErrNoRows {
		return apperror.NewNotFoundError(id, "while looking for User")
	} else if err != nil {
		return err
	}
	oktaID := user.OktaID

	// find any associated ServiceMember, OfficeUser, or AdminUser

	var serviceMember models.ServiceMember
	serviceMemberCount, err := appCtx.DB().Where("user_id = ?", id).Count(&serviceMember)
	if err == sql.ErrNoRows {
		appCtx.Logger().Debug("Not a ServiceMember")
	} else if err != nil {
		return err
	}

	var officeUser models.OfficeUser
	officeUserCount, err := appCtx.DB().Where("user_id = ?", id).Count(&officeUser)
	if err == sql.ErrNoRows {
		appCtx.Logger().Debug("Not an OfficeUser")
	} else if err != nil {
		return err
	}

	var adminUser models.AdminUser
	adminUserCount, err := appCtx.DB().Where("user_id = ?", id).Count(&adminUser)
	if err == sql.ErrNoRows {
		appCtx.Logger().Debug("Not an AdminUser")
	} else if err != nil {
		return err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		userIdFilter := []services.QueryFilter{query.NewQueryFilter("user_id", "=", user.ID.String())}
		if len(user.Roles) > 0 {
			// Delete associated roles (users_roles)
			err = o.builder.DeleteMany(appCtx, &[]models.UsersRoles{}, userIdFilter)
			if err != nil {
				return err
			}
		}

		if len(user.Privileges) > 0 {
			// Delete associated privileges (users_privileges)
			err = o.builder.DeleteMany(appCtx, &[]models.UsersPrivileges{}, userIdFilter)
			if err != nil {
				return err
			}
		}

		if serviceMemberCount > 0 {
			err = o.builder.DeleteMany(appCtx, &[]models.ServiceMember{}, userIdFilter)
			if err != nil {
				return err
			}
		}

		if officeUserCount > 0 {
			err = o.builder.DeleteMany(appCtx, &[]models.OfficeUser{}, userIdFilter)
			if err != nil {
				return err
			}
		}

		if adminUserCount > 0 {
			err = o.builder.DeleteMany(appCtx, &[]models.AdminUser{}, userIdFilter)
			if err != nil {
				return err
			}
		}

		// delete the user
		err = o.builder.DeleteOne(appCtx, &user)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		appCtx.Logger().Error(transactionError.Error())
		return transactionError
	}

	if oktaID != "" {
		/*
			Now that we have deleted the user from the milmove db, we will remove their okta account.
			We are intentionally keeping this process outside the milmove db delete transaction as it should not impact the ability to process a deletion from milmove db.
			This is considered more of a convenience to clean up the okta account.
		*/
		req := appCtx.HTTPRequest()
		if req == nil {
			appCtx.Logger().Error("failed to retrieve HTTP request from session")
			return nil
		}
		provider, err := okta.GetOktaProviderForRequest(req)
		if err != nil {
			appCtx.Logger().Error("error retrieving Okta provider: %w")
			return nil
		}
		apiKey := models.GetOktaAPIKey()
		err = models.DeleteOktaUser(appCtx, provider, oktaID, apiKey)
		if err != nil {
			appCtx.Logger().Error("error deleting user from okta: %w", zap.Error(err))
			return nil
		}
	}

	return nil
}

// NewUserDeleter returns a new user deleter builder
func NewUserDeleter(builder userQueryBuilder) services.UserDeleter {
	return &userDeleter{builder}
}
