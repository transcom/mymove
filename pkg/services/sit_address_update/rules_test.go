package sitaddressupdate

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *SITAddressUpdateServiceSuite) TestCheckRequiredFields() {
	suite.Run("Success", func() {
		suite.Run("Create SITAddressUpdate", func() {
			oldAddress := factory.BuildAddress(suite.DB(), nil, nil)
			mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusApproved,
					},
				},
				{
					Model: models.Address{},
					Type:  &factory.Addresses.SITDestinationFinalAddress,
				},
			}, nil)
			sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
				{
					Model:    oldAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.SITAddressUpdateOldAddress,
				},
				{
					Model:    mtoServiceItem,
					LinkOnly: true,
				},
			}, nil)

			err := checkAndValidateRequiredFields().Validate(
				suite.AppContextForTest(),
				&sitAddressUpdate,
			)

			suite.NilOrNoVerrs(err)
		})
	})

	suite.Run("Failure", func() {
		suite.Run("Create SITAddressUpdate with missing values", func() {
			sitAddressUpdate := models.SITAddressUpdate{}
			err := checkAndValidateRequiredFields().Validate(
				suite.AppContextForTest(),
				&sitAddressUpdate,
			)

			suite.Error(err)
			suite.IsType(&validate.Errors{}, err)
			suite.Contains(err.Error(), "NewAddress is required")
			suite.Contains(err.Error(), "MTOServiceItem is required")
			suite.Contains(err.Error(), "MTOServiceItem was not found")
			suite.Contains(err.Error(), "SITDestinationFinalAddressID is required")
		})

		suite.Run("Create SITAddressUpdate with rejected service item", func() {
			mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusRejected,
					},
				},
			}, nil)
			sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
				{
					Model:    mtoServiceItem,
					LinkOnly: true,
				},
			}, nil)

			err := checkAndValidateRequiredFields().Validate(
				suite.AppContextForTest(),
				&sitAddressUpdate,
			)

			suite.Error(err)
			suite.IsType(&validate.Errors{}, err)
			suite.Contains(err.Error(), "MTOServiceItem must be approved")
		})

		suite.Run("Create SITAddressUpdate with no service item", func() {
			sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
				{
					Model: models.SITAddressUpdate{
						OfficeRemarks: models.StringPointer("office remarks"),
					},
				},
			}, nil)

			err := checkAndValidateRequiredFields().Validate(
				suite.AppContextForTest(),
				&sitAddressUpdate,
			)

			suite.Error(err)
			suite.IsType(&validate.Errors{}, err)
			suite.Contains(err.Error(), "MTOServiceItem was not found")
		})

		suite.Run("Create SITAddressUpdate with missing SITDestinationFinalAddressID", func() {
			oldAddress := factory.BuildAddress(suite.DB(), nil, nil)
			mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusApproved,
					},
				}}, nil)
			sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
				{
					Model:    oldAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.SITAddressUpdateOldAddress,
				},
				{
					Model:    mtoServiceItem,
					LinkOnly: true,
				},
			}, nil)

			err := checkAndValidateRequiredFields().Validate(
				suite.AppContextForTest(),
				&sitAddressUpdate,
			)

			suite.Error(err)
			suite.IsType(&validate.Errors{}, err)
			suite.Contains(err.Error(), "SITDestinationFinalAddressID is required")
		})
	})
}

func (suite *SITAddressUpdateServiceSuite) TestCheckTOORequiredFields() {
	suite.Run("Success", func() {
		suite.Run("Create SITAddressUpdate", func() {
			oldAddress := factory.BuildAddress(suite.DB(), nil, nil)
			mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
				{
					Model: models.MTOServiceItem{
						Status: models.MTOServiceItemStatusApproved,
					},
				},
			}, nil)
			sitAddressUpdate := factory.BuildSITAddressUpdate(nil, []factory.Customization{
				{
					Model:    oldAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.SITAddressUpdateOldAddress,
				},
				{
					Model:    mtoServiceItem,
					LinkOnly: true,
				},
				{
					Model: models.SITAddressUpdate{
						OfficeRemarks: models.StringPointer("office remarks"),
					},
				},
			}, nil)

			err := checkTOORequiredFields().Validate(
				suite.AppContextForTest(),
				&sitAddressUpdate,
			)

			suite.NilOrNoVerrs(err)
		})
	})

	suite.Run("Failure", func() {
		suite.Run("Create SITAddressUpdate with missing values", func() {
			sitAddressUpdate := models.SITAddressUpdate{}

			err := checkTOORequiredFields().Validate(
				suite.AppContextForTest(),
				&sitAddressUpdate,
			)

			suite.Error(err)
			suite.IsType(&validate.Errors{}, err)
			suite.Contains(err.Error(), "are required")
		})
	})
}
