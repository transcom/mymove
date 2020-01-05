package usersroles

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type usersRolesCreator struct {
	db *pop.Connection
}

// NewNewUsersRolesCreator creates a new struct with the service dependencies
func NewUsersRolesCreator(db *pop.Connection) services.UserRoleAssociator {
	return usersRolesCreator{db}
}

//AssociateUserRoles associates a given user with a set of roles
func (u usersRolesCreator) AssociateUserRoles(userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	var usersRoles []models.UsersRoles
	allRoles := roles.Roles{}
	err := u.db.All(&allRoles)
	log.Println("allRoles", allRoles)
	if err != nil {
		return usersRoles, err
	}
	dbRoles, err := u.FetchUserRoles(userID)
	if err != nil {
		return usersRoles, err
	}
	log.Println("dbRoles", dbRoles)
	log.Println("rs", rs)
	toAdd := Difference(rs, dbRoles)
	log.Println("ToADD", toAdd)
	toRemove := Difference(dbRoles, rs)
	//SELECT r.id AS role_id, userID::uuid AS user_id,
	//	FROM roles r
	//LEFT JOIN users_roles ur ON ur.user_id = userID
	//AND r.id = ur.role_id`;
	var rolesToAdd []models.UsersRoles
	if len(toAdd) > 0 {
		err = u.db.Select("r.id as role_id, ? as user_id").
			RightJoin("roles r", "r.id=users_roles.role_id", userID).
			Where(`role_type IN (?) and users_roles.id IS NULL`, toAdd).All(&rolesToAdd)
		if err != nil {
			return usersRoles, err

		}
		err = u.db.Create(rolesToAdd)
		if err != nil {
			return usersRoles, err

		}
	}
	var ur []models.UsersRoles
	if len(toRemove) > 0 {
		err = u.db.Select("users_roles.id, r.id as role_id, ? as user_id").
			RightJoin("roles r", "r.id=users_roles.role_id", userID).
			Where(`role_type NOT IN (?) and users_roles.id IS NOT NULL`, toRemove).All(&ur)
		if err != nil {
			return usersRoles, err
		}
		err = u.db.Destroy(ur)
		if err != nil {
			return usersRoles, err
		}
	}
	log.Println("ToRemove", toRemove)

	return usersRoles, nil
}

// Set Difference: A - B
func Difference(a, b []roles.RoleType) (diff []roles.RoleType) {
	m := make(map[roles.RoleType]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func (u usersRolesCreator) FetchUserRoles(userID uuid.UUID) ([]roles.RoleType, error) {
	// select all roles not already associated with this user
	var userRoles []roles.RoleType
	err := u.db.RawQuery(
		`SELECT role_type
		FROM roles
		JOIN users_roles ur ON roles.id = ur.role_id
		WHERE user_id = $1`, userID).All(&userRoles)
	if err != nil {
		return userRoles, err
	}
	return userRoles, nil
}

func (u usersRolesCreator) fetchUnassociatedRoles(userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	// select all roles not already associated with this user
	var userRoles []models.UsersRoles
	rss := make([]interface{}, len(rs))
	for i := 1; i < len(rss); {
		rss[i] = rs[i]
		i++
	}
	err := u.db.RawQuery(
		`SELECT $1::uuid as user_id,
					  roles.id as role_id
	FROM roles
	WHERE role_type NOT IN (
		SELECT role_type
		FROM roles
				 JOIN users_roles ur ON roles.id = ur.role_id
		WHERE user_id = $1)`, userID).All(&userRoles)
	log.Println("fetchUnassociatedRoles.userRoles", userRoles)
	if err != nil {
		return userRoles, err
	}
	return userRoles, nil
}
