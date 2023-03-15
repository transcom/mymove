package factory

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildOrder() {
	suite.Run("Successful creation of default order", func() {
		// Under test:      BuildOrder
		// Set up:          Create a default order
		// Expected outcome:Create an extended service member,
		// UserUpload, duty location, origin duty location

		// SETUP
		// Create a default order infor to compare values
		defaultOrdersNumber := "ORDER3"
		defaultTACNumber := "F8E1"
		defaultDepartmentIndicator := "AIR_FORCE"
		defaultGrade := "E_1"
		defaultHasDependents := false
		defaultSpouseHasProGear := false
		defaultOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		defaultOrdersTypeDetail := internalmessages.OrdersTypeDetail("HHG_PERMITTED")
		defaultStatus := models.OrderStatusDRAFT
		testYear := 2018
		defaultIssueDate := time.Date(testYear, time.March, 15, 0, 0, 0, 0, time.UTC)
		defaultReportByDate := time.Date(testYear, time.August, 1, 0, 0, 0, 0, time.UTC)

		// Create order
		order := BuildOrder(suite.DB(), nil, nil)

		suite.Equal(defaultOrdersNumber, *order.OrdersNumber)
		suite.Equal(defaultTACNumber, *order.TAC)
		suite.Equal(defaultDepartmentIndicator, *order.DepartmentIndicator)
		suite.Equal(defaultGrade, *order.Grade)
		suite.Equal(defaultHasDependents, order.HasDependents)
		suite.Equal(defaultSpouseHasProGear, order.SpouseHasProGear)
		suite.Equal(defaultOrdersType, order.OrdersType)
		suite.Equal(defaultOrdersTypeDetail, *order.OrdersTypeDetail)
		suite.Equal(defaultStatus, order.Status)
		suite.Equal(defaultIssueDate, order.IssueDate)
		suite.Equal(defaultReportByDate, order.ReportByDate)

		// extended service members have backup contacts
		suite.False(order.ServiceMemberID.IsNil())
		suite.False(order.ServiceMember.ID.IsNil())
		suite.False(order.ServiceMember.DutyLocationID.IsNil())
		suite.NotEmpty(order.ServiceMember.BackupContacts)
		serviceMemberCountInDB, err := suite.DB().Count(models.ServiceMember{})
		suite.NoError(err)
		suite.Equal(1, serviceMemberCountInDB)

		// uses the same duty location name for service member and
		// orders OriginDutyLocation
		suite.Equal(order.ServiceMember.DutyLocation.Name, order.OriginDutyLocation.Name)

		// uses the default orders NewDutyLocation
		suite.Equal(order.NewDutyLocation.Name, "Fort Gordon")

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
}
