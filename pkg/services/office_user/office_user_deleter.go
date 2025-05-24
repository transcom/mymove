package officeuser

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

type officeUserDeleter struct {
	builder officeUserQueryBuilder
}

var foreignKeyPattern = regexp.MustCompile("violates foreign key constraint")

func (o *officeUserDeleter) DeleteOfficeUser(appCtx appcontext.AppContext, id uuid.UUID) error {
	var officeUser models.OfficeUser
	err := appCtx.DB().EagerPreload(
		"User",
	).Where("id = ?", id).Find(&officeUser, id)
	if err == sql.ErrNoRows {
		return apperror.NewNotFoundError(id, "while looking for OfficeUser")
	} else if err != nil {
		return err
	}

	user := officeUser.User
	oktaID := user.OktaID
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		userIdFilter := []services.QueryFilter{query.NewQueryFilter("user_id", "=", user.ID.String())}

		// Delete associated roles (users_roles)
		err = o.builder.DeleteMany(txnAppCtx, &[]models.UsersRoles{}, userIdFilter)
		if err != nil {
			return err
		}

		// Delete associated privileges (users_privileges)
		err = o.builder.DeleteMany(txnAppCtx, &[]models.UsersPrivileges{}, userIdFilter)
		if err != nil {
			return err
		}

		// delete the office user (office_users)
		err = o.builder.DeleteOne(txnAppCtx, &officeUser)
		if err != nil {
			return handleError(id, err)
		}

		// finally, delete the user (user)
		err = o.builder.DeleteOne(txnAppCtx, &user)
		if err != nil {
			return handleError(id, err)
		}

		return nil
	})

	if transactionError != nil {
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

// NewOfficeUserDeleter returns a new office user deleter builder
func NewOfficeUserDeleter(builder officeUserQueryBuilder) services.OfficeUserDeleter {
	return &officeUserDeleter{builder}
}

func handleError(id uuid.UUID, rawError error) error {
	if foreignKeyPattern.MatchString(rawError.Error()) {
		return apperror.NewConflictError(id, rawError.Error())
	}
	return rawError
}
