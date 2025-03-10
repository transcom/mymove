package officeuser

import (
	"database/sql"

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

func (o *officeUserDeleter) DeleteOfficeUser(appCtx appcontext.AppContext, id uuid.UUID) error {
	// need to fetch the office user and any downstream associations (roles, privileges)
	var officeUser models.OfficeUser
	err := appCtx.DB().EagerPreload(
		"User",
		"User.Roles",
		"User.Privileges",
	).Where("id = ?", id).Find(&officeUser, id)
	if err == sql.ErrNoRows {
		return apperror.NewNotFoundError(id, "while looking for OfficeUser")
	} else if err != nil {
		return err
	}

	user := officeUser.User
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

		// delete the office user (office_users)
		err = o.builder.DeleteOne(appCtx, &officeUser)
		if err != nil {
			return err
		}

		// finally, delete the user (user)
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

	return nil
}

// NewOfficeUserDeleter returns a new office user deleter builder
func NewOfficeUserDeleter(builder officeUserQueryBuilder) services.OfficeUserDeleter {
	return &officeUserDeleter{builder}
}
