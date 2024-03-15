package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildPrivilege creates a Privilege
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildPrivilege(db *pop.Connection, customs []Customization, traits []Trait) models.Privilege {
	customs = setupCustomizations(customs, traits)

	// Find privilege assertion and convert to privileges Privilege
	var cPrivilege models.Privilege
	if result := findValidCustomization(customs, Privilege); result != nil {
		cPrivilege = result.Model.(models.Privilege)
		if result.LinkOnly {
			return cPrivilege
		}
	}

	// create privilege
	privilegeUUID := uuid.Must(uuid.NewV4())
	privilege := models.Privilege{
		ID:            privilegeUUID,
		PrivilegeType: models.PrivilegeTypeSupervisor,
		PrivilegeName: "Supervisor",
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&privilege, cPrivilege)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &privilege)
	}

	return privilege
}

// lookup a privilege by privilege type, if it doesn't exist make it
func FetchOrBuildPrivilegeByPrivilegeType(db *pop.Connection, privilegeType models.PrivilegeType) models.Privilege {
	privilegeName := models.PrivilegeName(cases.Title(language.Und).String(string(privilegeType)))

	if db == nil {
		return BuildPrivilege(db, []Customization{
			{
				Model: models.Privilege{
					PrivilegeType: privilegeType,
					PrivilegeName: privilegeName,
				},
			},
		}, nil)
	}

	var privilege models.Privilege
	err := db.Where("privilege_type=$1", privilegeType).First(&privilege)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return privilege
	}

	return BuildPrivilege(db, []Customization{
		{
			Model: models.Privilege{
				PrivilegeType: privilegeType,
				PrivilegeName: privilegeName,
			},
		},
	}, nil)
}
