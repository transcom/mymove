package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildOrder() {
	defaultOrdersNumber := "ORDER3"
	defaultTACNumber := "F8E1"
	defaultDepartmentIndicator := "AIR_AND_SPACE_FORCE"
	defaultGrade := "E-1"
	defaultHasDependents := false
	defaultSpouseHasProGear := false
	defaultOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	defaultOrdersTypeDetail := internalmessages.OrdersTypeDetail("HHG_PERMITTED")
	defaultStatus := models.OrderStatusDRAFT
	testYear := 2018
	defaultIssueDate := time.Date(testYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	defaultReportByDate := time.Date(testYear, time.August, 1, 0, 0, 0, 0, time.UTC)
	defaultGBLOC := "KKFA"

	suite.Run("Successful creation of default order", func() {
		// Under test:      BuildOrder
		// Set up:          Create a default order
		// Expected outcome:Create an extended service member,
		// UserUpload, duty location, origin duty location

		// SETUP
		// Create a default order info to compare values

		// Create order
		order := BuildOrder(suite.DB(), nil, nil)

		suite.Equal(defaultOrdersNumber, *order.OrdersNumber)
		suite.Equal(defaultTACNumber, *order.TAC)
		suite.Equal(defaultDepartmentIndicator, *order.DepartmentIndicator)
		suite.Equal(defaultGrade, string(*order.Grade))
		suite.Equal(defaultHasDependents, order.HasDependents)
		suite.Equal(defaultSpouseHasProGear, order.SpouseHasProGear)
		suite.Equal(defaultOrdersType, order.OrdersType)
		suite.Equal(defaultOrdersTypeDetail, *order.OrdersTypeDetail)
		suite.Equal(defaultStatus, order.Status)
		suite.Equal(defaultIssueDate, order.IssueDate)
		suite.Equal(defaultReportByDate, order.ReportByDate)
		suite.Equal(defaultGBLOC, *order.OriginDutyLocationGBLOC)

		// extended service members have backup contacts
		suite.False(order.ServiceMemberID.IsNil())
		suite.False(order.ServiceMember.ID.IsNil())
		suite.False(order.OriginDutyLocationID.IsNil())
		suite.NotEmpty(order.ServiceMember.BackupContacts)
		serviceMemberCountInDB, err := suite.DB().Count(models.ServiceMember{})
		suite.NoError(err)
		suite.Equal(1, serviceMemberCountInDB)

		// uses the default orders NewDutyLocation
		suite.Equal(order.NewDutyLocation.Name, "Fort Eisenhower, GA 30813")

		dutyLocationCountInDB, err := suite.DB().Count(models.DutyLocation{})
		suite.NoError(err)
		// origin and new duty location
		suite.Equal(2, dutyLocationCountInDB)

		// creates uploaded orders as user upload with same service member
		suite.False(order.UploadedOrdersID.IsNil())
		suite.False(order.UploadedOrders.ID.IsNil())
		suite.NotEmpty(order.UploadedOrders.UserUploads)
		suite.Equal(order.UploadedOrders.ServiceMember.UserID,
			order.ServiceMember.UserID)
		modelCountInDB, err := suite.DB().Count(models.Document{})
		suite.NoError(err)
		suite.Equal(1, modelCountInDB)
		userUploadCountInDB, err := suite.DB().Count(models.UserUpload{})
		suite.NoError(err)
		suite.Equal(1, userUploadCountInDB)

		// no amended orders
		suite.Nil(order.UploadedAmendedOrders)
		suite.Nil(order.UploadedAmendedOrdersID)
	})

	suite.Run("Successful creation of customized order", func() {
		originDutyLocation := models.DutyLocation{
			Name: "Custom Origin",
		}
		originDutyLocationTOName := "origin duty location transportation office"
		firstName := "customFirst"
		lastName := "customLast"
		serviceMember := models.ServiceMember{
			FirstName: &firstName,
			LastName:  &lastName,
		}
		uploadedOrders := models.Document{
			ID: uuid.Must(uuid.NewV4()),
		}
		dependents := 7
		entitlement := models.Entitlement{
			TotalDependents: &dependents,
		}
		amendedOrders := models.Document{
			ID: uuid.Must(uuid.NewV4()),
		}
		customs := []Customization{
			{
				Model: originDutyLocation,
				Type:  &DutyLocations.OriginDutyLocation,
			},
			{
				Model: models.TransportationOffice{
					Name: originDutyLocationTOName,
				},
				Type: &TransportationOffices.OriginDutyLocation,
			},
			{
				Model: serviceMember,
			},
			{
				Model: uploadedOrders,
				Type:  &Documents.UploadedOrders,
			},
			{
				Model: entitlement,
			},
			{
				Model: amendedOrders,
				Type:  &Documents.UploadedAmendedOrders,
			},
		}
		// Create order
		order := BuildOrder(suite.DB(), customs, nil)

		suite.Equal(originDutyLocation.Name, order.OriginDutyLocation.Name)
		suite.Equal(originDutyLocationTOName, order.OriginDutyLocation.TransportationOffice.Name)
		suite.Equal(*serviceMember.FirstName, *order.ServiceMember.FirstName)
		suite.Equal(*serviceMember.LastName, *order.ServiceMember.LastName)
		suite.Equal(uploadedOrders.ID, order.UploadedOrdersID)
		suite.Equal(uploadedOrders.ID, order.UploadedOrders.ID)
		suite.Equal(*entitlement.TotalDependents, *order.Entitlement.TotalDependents)
		suite.Equal(amendedOrders.ID, *order.UploadedAmendedOrdersID)
		suite.Equal(amendedOrders.ID, order.UploadedAmendedOrders.ID)
	})
	suite.Run("Successful creation of order with prebuilt uploaded orders", func() {
		userUploadForUploadedOrders := BuildUserUpload(suite.DB(), nil, nil)
		uploadedOrders := userUploadForUploadedOrders.Document
		uploadedOrders.UserUploads = models.UserUploads{userUploadForUploadedOrders}

		userUploadForAmendedOrders := BuildUserUpload(suite.DB(), nil, nil)
		amendedOrders := userUploadForAmendedOrders.Document
		amendedOrders.UserUploads = models.UserUploads{userUploadForAmendedOrders}

		order := BuildOrder(suite.DB(), []Customization{
			{
				Model:    uploadedOrders,
				LinkOnly: true,
				Type:     &Documents.UploadedOrders,
			},
			{
				Model:    amendedOrders,
				LinkOnly: true,
				Type:     &Documents.UploadedAmendedOrders,
			},
		}, nil)
		suite.Equal(order.UploadedOrdersID, uploadedOrders.ID)
		suite.Equal(1, len(order.UploadedOrders.UserUploads))
		suite.Equal(1, len(order.UploadedAmendedOrders.UserUploads))
	})
	suite.Run("Successful creation of stubbed order", func() {
		// Under test:      BuildOrder
		// Set up:          Create an order, but don't pass in a db
		// Expected outcome:Order should be created
		//                  No order should be created in database
		precount, err := suite.DB().Count(&models.Order{})
		suite.NoError(err)

		order := BuildOrder(nil, nil, nil)
		suite.Equal(defaultOrdersNumber, *order.OrdersNumber)
		suite.Equal(defaultTACNumber, *order.TAC)
		suite.Equal(defaultDepartmentIndicator, *order.DepartmentIndicator)
		suite.Equal(defaultGrade, string(*order.Grade))
		suite.Equal(defaultHasDependents, order.HasDependents)
		suite.Equal(defaultSpouseHasProGear, order.SpouseHasProGear)
		suite.Equal(defaultOrdersType, order.OrdersType)
		suite.Equal(defaultOrdersTypeDetail, *order.OrdersTypeDetail)
		suite.Equal(defaultStatus, order.Status)
		suite.Equal(defaultIssueDate, order.IssueDate)
		suite.Equal(defaultReportByDate, order.ReportByDate)

		suite.NotEmpty(order.OriginDutyLocation.Address.PostalCode)

		count, err := suite.DB().Count(&models.Order{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})
	suite.Run("Success creation of order without defaults", func() {
		// Under test:      BuildOrderWithoutDefaults
		// Set up:          Create an order without defaults
		// Expected outcome:Order should be created

		order := BuildOrderWithoutDefaults(suite.DB(), nil, nil)
		suite.Nil(order.OrdersNumber)
		suite.Nil(order.TAC)
		suite.Nil(order.DepartmentIndicator)
		suite.Nil(order.OrdersTypeDetail)
		suite.Equal(defaultGrade, string(*order.Grade))
		suite.Equal(defaultHasDependents, order.HasDependents)
		suite.Equal(defaultSpouseHasProGear, order.SpouseHasProGear)
		suite.Equal(defaultOrdersType, order.OrdersType)
		suite.Equal(defaultStatus, order.Status)
		suite.Equal(defaultIssueDate, order.IssueDate)
		suite.Equal(defaultReportByDate, order.ReportByDate)
	})
	suite.Run("Success creation of stubbed order without defaults", func() {
		// Under test:      BuildOrderWithoutDefaults
		// Set up:          Create an order, but don't pass in a db
		// Expected outcome:Order should be created
		//                  No order should be created in database
		precount, err := suite.DB().Count(&models.Order{})
		suite.NoError(err)

		order := BuildOrderWithoutDefaults(nil, nil, nil)
		suite.Nil(order.OrdersNumber)
		suite.Nil(order.TAC)
		suite.Nil(order.DepartmentIndicator)
		suite.Nil(order.OrdersTypeDetail)
		suite.Equal(defaultGrade, string(*order.Grade))
		suite.Equal(defaultHasDependents, order.HasDependents)
		suite.Equal(defaultSpouseHasProGear, order.SpouseHasProGear)
		suite.Equal(defaultOrdersType, order.OrdersType)
		suite.Equal(defaultStatus, order.Status)
		suite.Equal(defaultIssueDate, order.IssueDate)
		suite.Equal(defaultReportByDate, order.ReportByDate)

		count, err := suite.DB().Count(&models.Order{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}
