package user

import (
	"database/sql"

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

	return nil
}

// NewUserDeleter returns a new user deleter builder
func NewUserDeleter(builder userQueryBuilder) services.UserDeleter {
	return &userDeleter{builder}
}
