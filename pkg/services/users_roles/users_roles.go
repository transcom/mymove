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
	rToAdd := UserRoles(allRoles, toAdd, userID)
	err = u.db.Create(rToAdd)
	if err != nil {
		return usersRoles, err
	}
	rToRemove := UserRoles(allRoles, toRemove, userID)
	var rIDSTOdelete []uuid.UUID
	for _, m := range rToRemove {
		rIDSTOdelete = append(rIDSTOdelete, m.RoleID)
	}
	var ur []models.UsersRoles
	err = u.db.Where("role_id in (?)", rIDSTOdelete).Where("user_id = ?", userID).All(&ur)
	if err != nil {
		return usersRoles, err
	}
	log.Println("ToRemove", toRemove)
	err = u.db.Destroy(ur)
	if err != nil {
		return usersRoles, err
	}
	return usersRoles, nil
}

func UserRoles(all roles.Roles, filterOn []roles.RoleType, userID uuid.UUID) []models.UsersRoles {
	var output []models.UsersRoles
	for _, s := range filterOn {
		for _, role := range all {
			if s == role.RoleType {
				record := models.UsersRoles{
					UserID: userID,
					RoleID: role.ID,
				}
				output = append(output, record)
			}
		}
	}
	return output
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
