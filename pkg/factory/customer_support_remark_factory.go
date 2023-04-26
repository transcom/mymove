package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildCustomerSupportRemark creates a single CustomerSupportRemark.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildCustomerSupportRemark(db *pop.Connection, customs []Customization, traits []Trait) models.CustomerSupportRemark {
	customs = setupCustomizations(customs, traits)

	// Find customerSupportRemark customization and extract the custom customerSupportRemark
	var cCustomerSupportRemark models.CustomerSupportRemark
	if result := findValidCustomization(customs, CustomerSupportRemark); result != nil {
		cCustomerSupportRemark = result.Model.(models.CustomerSupportRemark)
		if result.LinkOnly {
			return cCustomerSupportRemark
		}
	}

	move := BuildMove(db, customs, traits)

	officeUser := BuildOfficeUser(db, customs, traits)

	defaultContent := "This is an office remark."

	// Create default CustomerSupportRemark
	customerSupportRemark := models.CustomerSupportRemark{
		Content:      defaultContent,
		OfficeUserID: officeUser.ID,
		OfficeUser:   officeUser,
		MoveID:       move.ID,
		Move:         move,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&customerSupportRemark, cCustomerSupportRemark)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &customerSupportRemark)
	}

	return customerSupportRemark
}
