package factory

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildOfficeUser creates an OfficeUser, and a transportation office and transportation office assignment if either doesn't exist
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
// Notes:
//   - To build an office user with one or more roles use BuildOfficeUserWithRoles
//   - There's a uniqueness constraint on office user emails so use the GetTraitOfficeUserEmail trait
//     when creating a test with multiple office users
//   - The OfficeUser returned won't have an ID if the db is nil. If an ID is needed for a stubbed user,
//     use trait GetTraitOfficeUserWithID
func BuildOfficeUser(db *pop.Connection, customs []Customization, traits []Trait) models.OfficeUser {
	customs = setupCustomizations(customs, traits)

	// Find officeuser assertion and convert to models officeuser
	var cOfficeUser models.OfficeUser
	if result := findValidCustomization(customs, OfficeUser); result != nil {
		cOfficeUser = result.Model.(models.OfficeUser)
		if result.LinkOnly {
			return cOfficeUser
		}
	}

	// Find/create the user model
	user := BuildUserAndUsersRoles(db, customs, nil)

	var transportationOffice models.TransportationOffice
	tempCloseoutOfficeCustoms := customs
	tempCounselingOfficeCustoms := customs
	closeoutOfficeResult := findValidCustomization(customs, TransportationOffices.CloseoutOffice)
	counselingOfficeResult := findValidCustomization(customs, TransportationOffices.CounselingOffice)
	if closeoutOfficeResult != nil {
		tempCloseoutOfficeCustoms = convertCustomizationInList(tempCloseoutOfficeCustoms, TransportationOffices.CloseoutOffice, TransportationOffice)
		transportationOffice = BuildTransportationOffice(db, tempCloseoutOfficeCustoms, nil)
	} else if counselingOfficeResult != nil {
		tempCounselingOfficeCustoms = convertCustomizationInList(tempCounselingOfficeCustoms, TransportationOffices.CounselingOffice, TransportationOffice)
		transportationOffice = BuildTransportationOffice(db, tempCounselingOfficeCustoms, nil)
	} else {
		transportationOffice = BuildTransportationOffice(db, customs, nil)
	}

	// create officeuser
	officeUser := models.OfficeUser{
		UserID:                 &user.ID,
		User:                   user,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  "leo_spaceman_office@example.com",
		Telephone:              "415-555-1212",
		TransportationOffice:   transportationOffice,
		TransportationOfficeID: transportationOffice.ID,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&officeUser, cOfficeUser)

	// If db is false, it's a stub. No need to create in database
	// If OfficeUser or Transportation office is a stub there is nothing to link with
	// a transportation office assignment and the link will fail due to nil IDs
	if db != nil {
		mustCreate(db, &officeUser)

		BuildPrimaryTransportationOfficeAssignment(db, []Customization{
			{
				Model: models.OfficeUser{
					ID: officeUser.ID,
				},
				LinkOnly: true,
			},
			{
				Model: models.TransportationOffice{
					ID: transportationOffice.ID,
				},
				LinkOnly: true,
			},
		}, nil)
	}

	return officeUser
}

// BuildOfficeUserWithoutTransportationAssignment creates an OfficeUser
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
// Notes:
//   - To build an office user with one or more roles use BuildOfficeUserWithRoles
//   - There's a uniqueness constraint on office user emails so use the GetTraitOfficeUserEmail trait
//     when creating a test with multiple office users
//   - The OfficeUser returned won't have an ID if the db is nil. If an ID is needed for a stubbed user,
//     use trait GetTraitOfficeUserWithID
func BuildOfficeUserWithoutTransportationOfficeAssignment(db *pop.Connection, customs []Customization, traits []Trait) models.OfficeUser {
	customs = setupCustomizations(customs, traits)

	// Find officeuser assertion and convert to models officeuser
	var cOfficeUser models.OfficeUser
	if result := findValidCustomization(customs, OfficeUser); result != nil {
		cOfficeUser = result.Model.(models.OfficeUser)
		if result.LinkOnly {
			return cOfficeUser
		}
	}

	// Find/create the user model
	user := BuildUserAndUsersRoles(db, customs, nil)

	// Find/create the TransportationOffice model
	transportationOffice := BuildTransportationOffice(db, customs, nil)

	// create officeuser
	officeUser := models.OfficeUser{
		UserID:                 &user.ID,
		User:                   user,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  "leo_spaceman_office@example.com",
		Telephone:              "415-555-1212",
		TransportationOffice:   transportationOffice,
		TransportationOfficeID: transportationOffice.ID,
	}
	// Overwrite values with those from assertions
	testdatagen.MergeModels(&officeUser, cOfficeUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &officeUser)
	}

	return officeUser
}

// BuildOfficeUserWithRoles returns an office user with an ID, unique email, and roles
// Also creates
//   - User
//   - Role
//   - UsersRoles
//
// Notes:
//   - roleTypes passed into the function will overwrite over any roles in a User customization
//   - a unique email for the user will be created
//   - a UUID will be added to the OfficeUser record when it's stubbed
func BuildOfficeUserWithRoles(db *pop.Connection, customs []Customization, roleTypes []roles.RoleType) models.OfficeUser {
	customs = setupCustomizations(customs, nil)

	var rolesList []roles.Role
	for _, roleType := range roleTypes {
		role := roles.Role{
			RoleType: roleType,
		}
		rolesList = append(rolesList, role)
	}

	traits := []Trait{GetTraitOfficeUserEmail}
	if db == nil {
		// UUIDs are only set when saving to a DB, but they're necessary when checking session auths
		traits = append(traits, GetTraitOfficeUserWithID)
	}

	// Find/create the user model
	// If there is a user customization, add the roles to it, otherwise add a new user customization
	var user models.User
	idx, result := findCustomWithIdx(customs, User)
	if result != nil {
		// add roles to the existing user customization
		user = result.Model.(models.User)
		user.Roles = rolesList
		customs[idx].Model = user
	} else {
		// create a new user customization with the correct roles
		user.Roles = rolesList
		customs = append(customs, Customization{Model: user})
	}

	return BuildOfficeUser(db, customs, traits)
}

// BuildOfficeUserWithPrivileges returns an office user with an ID, unique email, and privileges
// Also creates
//   - User
//   - Privilege
//   - UsersPrivileges
//
// Notes:
//   - privilegeTypes passed into the function will overwrite over any privileges in a User customization
//   - a unique email for the user will be created
//   - a UUID will be added to the OfficeUser record when it's stubbed
func BuildOfficeUserWithPrivileges(db *pop.Connection, customs []Customization, traits []Trait) models.OfficeUser {
	customs = setupCustomizations(customs, traits)

	// Find officeuser assertion and convert to models officeuser
	var cOfficeUser models.OfficeUser
	if result := findValidCustomization(customs, OfficeUser); result != nil {
		cOfficeUser = result.Model.(models.OfficeUser)
		if result.LinkOnly {
			return cOfficeUser
		}
	}

	// Find/create the user model
	user := BuildUserAndUsersRolesAndUsersPrivileges(db, customs, nil)

	// Find/create the TransportationOffice model
	basicTransportationOffice := BuildTransportationOffice(db, customs, nil)

	var closeoutOffice models.TransportationOffice
	tempCloseoutOfficeCustoms := customs
	closeoutOfficeResult := findValidCustomization(customs, TransportationOffices.CloseoutOffice)
	if closeoutOfficeResult != nil {
		tempCloseoutOfficeCustoms = convertCustomizationInList(tempCloseoutOfficeCustoms, TransportationOffices.CloseoutOffice, TransportationOffice)
		closeoutOffice = BuildTransportationOffice(db, tempCloseoutOfficeCustoms, nil)
	}

	// create officeuser
	officeUser := models.OfficeUser{
		UserID:    &user.ID,
		User:      user,
		FirstName: "Leo",
		LastName:  "Spaceman",
		Email:     "leo_spaceman_office@example.com",
		Telephone: "415-555-1212",
	}

	if closeoutOfficeResult != nil {
		officeUser.TransportationOffice = closeoutOffice
		officeUser.TransportationOfficeID = closeoutOffice.ID
	} else {
		officeUser.TransportationOffice = basicTransportationOffice
		officeUser.TransportationOfficeID = basicTransportationOffice.ID
	}
	// Overwrite values with those from assertions
	testdatagen.MergeModels(&officeUser, cOfficeUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &officeUser)
	}

	return officeUser
}

// ------------------------
//        TRAITS
// ------------------------

// GetTraitOfficeUserEmail helps comply with the uniqueness constraint on emails
func GetTraitOfficeUserEmail() []Customization {
	// There's a uniqueness constraint on office user emails so add some randomness
	email := strings.ToLower(fmt.Sprintf("leo_spaceman_office_%s@example.com", MakeRandomString(5)))
	return []Customization{
		{
			Model: models.User{
				OktaEmail: email,
			},
			Type: &User,
		},
		{
			Model: models.OfficeUser{
				Email: email,
			},
			Type: &OfficeUser,
		},
	}
}

// GetTraitOfficeUserWithID adds a UUID to the record regardless of whether it's stubbed or not
func GetTraitOfficeUserWithID() []Customization {
	return []Customization{
		{
			Model: models.OfficeUser{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.User{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
	}
}

// GetTraitOfficeUserActive sets the User and OfficeUser as Active
func GetTraitActiveOfficeUser() []Customization {
	return []Customization{
		{
			Model: models.OfficeUser{
				Active: true,
			},
		},
		{
			Model: models.User{
				Active: true,
			},
		},
	}
}

// GetTraitApprovedOfficeUser sets the OfficeUser in an APPROVED status
func GetTraitApprovedOfficeUser() []Customization {
	approvedStatus := models.OfficeUserStatusAPPROVED
	return []Customization{
		{
			Model: models.OfficeUser{
				Status: &approvedStatus,
			},
		},
	}
}

// GetTraitRequestedOfficeUser sets the OfficeUser in an REQUESTED status
func GetTraitRequestedOfficeUser() []Customization {
	requestedStatus := models.OfficeUserStatusREQUESTED
	return []Customization{
		{
			Model: models.OfficeUser{
				Status: &requestedStatus,
			},
		},
	}
}
