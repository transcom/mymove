package user

import (
	"database/sql"
	"regexp"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userDeleter struct {
	builder userQueryBuilder
}

var foreignKeyPattern = regexp.MustCompile("violates foreign key constraint")

func (o *userDeleter) DeleteUser(appCtx appcontext.AppContext, id uuid.UUID) error {
	var user models.User
	err := appCtx.DB().Where("id = ?", id).Find(&user, id)
	if err == sql.ErrNoRows {
		return apperror.NewNotFoundError(id, "while looking for User")
	} else if err != nil {
		return err
	}
	oktaID := user.OktaID

	var adminUser models.AdminUser
	adminUserCount, err := appCtx.DB().Where("user_id = ?", id).Count(&adminUser)
	if err == sql.ErrNoRows {
		appCtx.Logger().Debug("Not an AdminUser")
	} else if err != nil {
		return err
	} else if adminUserCount > 0 {
		return apperror.NewForbiddenError("This is an Admin User and cannot be deleted.")
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		userIdFilter := []services.QueryFilter{query.NewQueryFilter("user_id", "=", user.ID.String())}

		err = o.builder.DeleteMany(txnAppCtx, &[]models.UsersRoles{}, userIdFilter)
		if err != nil {
			return err
		}

		err = o.builder.DeleteMany(txnAppCtx, &[]models.UsersPrivileges{}, userIdFilter)
		if err != nil {
			return err
		}

		var serviceMember models.ServiceMember
		err = o.builder.FetchOne(txnAppCtx, &serviceMember, userIdFilter)
		if err == nil {
			serviceMemberIdFilter := []services.QueryFilter{query.NewQueryFilter("service_member_id", "=", serviceMember.ID.String())}
			err = o.builder.DeleteMany(txnAppCtx, &[]models.BackupContact{}, serviceMemberIdFilter)
		}
		err = o.builder.DeleteMany(txnAppCtx, &[]models.ServiceMember{}, userIdFilter)
		if err != nil {
			return handleError(id, err)
		}

		err = o.builder.DeleteMany(txnAppCtx, &[]models.OfficeUser{}, userIdFilter)
		if err != nil {
			return handleError(id, err)
		}

		// delete the user
		err = o.builder.DeleteOne(txnAppCtx, &user)
		if err != nil {
			return handleError(id, err)
		}

		return nil
	})

	if transactionError != nil {
		appCtx.Logger().Error("transaction error...")
		appCtx.Logger().Error(transactionError.Error())
		return transactionError
	}

	/*
		Now that we have deleted the user from the milmove db, we will remove their okta account.
		We are intentionally keeping this process outside the milmove db delete transaction as it should not impact the ability to process a deletion from milmove db.
		This is considered more of a convenience to clean up the okta account.
	*/
	models.DeleteOktaUserHandled(appCtx, oktaID)

	return nil
}

// NewUserDeleter returns a new user deleter builder
func NewUserDeleter(builder userQueryBuilder) services.UserDeleter {
	return &userDeleter{builder}
}

func handleError(id uuid.UUID, rawError error) error {
	if foreignKeyPattern.MatchString(rawError.Error()) {
		return apperror.NewConflictError(id, rawError.Error())
	}
	return rawError
}
