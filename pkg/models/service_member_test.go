package models_test

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicServiceMemberInstantiation() {
	servicemember := &ServiceMember{}

	expErrors := map[string][]string{
		"user_id": {"UserID can not be blank."},
	}

	suite.verifyValidationErrors(servicemember, expErrors)
}

func (suite *ModelSuite) TestIsProfileCompleteWithIncompleteSM() {
	// Given: a user and a service member
	lgu := uuid.Must(uuid.NewV4())
	user1 := User{
		LoginGovUUID:  &lgu,
		LoginGovEmail: "whoever@example.com",
	}
	verrs, err := user1.Validate(nil)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")

	// And: a service member is incompletely initialized with almost all required values
	edipi := "12345567890"
	affiliation := AffiliationARMY
	rank := ServiceMemberRankE5
	firstName := "bob"
	lastName := "sally"
	telephone := "510 555-5555"
	email := "bobsally@gmail.com"
	fakeAddress := testdatagen.MakeStubbedAddress(suite.DB())
	fakeBackupAddress := testdatagen.MakeStubbedAddress(suite.DB())
	location := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		DutyLocation: DutyLocation{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})

	serviceMember := ServiceMember{
		ID:                     uuid.Must(uuid.NewV4()),
		UserID:                 user1.ID,
		Edipi:                  &edipi,
		Affiliation:            &affiliation,
		Rank:                   &rank,
		FirstName:              &firstName,
		LastName:               &lastName,
		Telephone:              &telephone,
		PersonalEmail:          &email,
		ResidentialAddressID:   &fakeAddress.ID,
		BackupMailingAddressID: &fakeBackupAddress.ID,
		DutyStationID:          &location.ID,
	}

	suite.Equal(false, serviceMember.IsProfileComplete())

	// When: all required fields are set
	emailPreferred := true
	serviceMember.EmailIsPreferred = &emailPreferred

	contactAssertions := testdatagen.Assertions{
		BackupContact: BackupContact{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
		Stub: true,
	}
	backupContact := testdatagen.MakeBackupContact(suite.DB(), contactAssertions)
	serviceMember.BackupContacts = append(serviceMember.BackupContacts, backupContact)

	suite.Equal(true, serviceMember.IsProfileComplete())
}

func (suite *ModelSuite) TestFetchServiceMemberForUser() {
	user1 := testdatagen.MakeDefaultUser(suite.DB())
	user2 := testdatagen.MakeDefaultUser(suite.DB())

	firstName := "Oliver"
	resAddress := testdatagen.MakeDefaultAddress(suite.DB())
	sm := ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.MustSave(&sm)

	// User is authorized to fetch service member
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user1.ID,
		ServiceMemberID: sm.ID,
	}
	goodSm, err := FetchServiceMemberForUser(suite.DB(), session, sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddress.ID, goodSm.ResidentialAddress.ID)
	}

	// Wrong ServiceMember
	wrongID, _ := uuid.NewV4()
	_, err = FetchServiceMemberForUser(suite.DB(), session, wrongID)
	if suite.Error(err) {
		suite.Equal(ErrFetchNotFound, err)
	}

	// User is forbidden from fetching order
	session.UserID = user2.ID
	session.ServiceMemberID = uuid.Nil
	_, err = FetchServiceMemberForUser(suite.DB(), session, sm.ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}
}

func (suite *ModelSuite) TestFetchServiceMemberNotForUser() {
	user1 := testdatagen.MakeDefaultUser(suite.DB())

	firstName := "Nino"
	resAddress := testdatagen.MakeDefaultAddress(suite.DB())
	sm := ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.MustSave(&sm)

	goodSm, err := FetchServiceMember(suite.DB(), sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddressID, goodSm.ResidentialAddressID)
	}
}

func (suite *ModelSuite) TestFetchLatestOrders() {
	user := testdatagen.MakeDefaultUser(suite.DB())

	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

	dutyLocation := testdatagen.FetchOrMakeDefaultCurrentDutyLocation(suite.DB())
	dutyLocation2 := testdatagen.FetchOrMakeDefaultNewOrdersDutyLocation(suite.DB())
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hasDependents := true
	spouseHasProGear := true
	uploadedOrder := Document{
		ServiceMember:   serviceMember,
		ServiceMemberID: serviceMember.ID,
	}
	deptIndicator := testdatagen.DefaultDepartmentIndicator
	TAC := testdatagen.DefaultTransportationAccountingCode
	suite.MustSave(&uploadedOrder)

	SAC := "N002214CSW32Y9"
	ordersNumber := "FD4534JFJ"

	order := Order{
		ServiceMemberID:      serviceMember.ID,
		ServiceMember:        serviceMember,
		IssueDate:            issueDate,
		ReportByDate:         reportByDate,
		OrdersType:           ordersType,
		HasDependents:        hasDependents,
		SpouseHasProGear:     spouseHasProGear,
		OriginDutyLocationID: &dutyLocation.ID,
		OriginDutyLocation:   &dutyLocation,
		NewDutyLocationID:    dutyLocation2.ID,
		NewDutyLocation:      dutyLocation2,
		UploadedOrdersID:     uploadedOrder.ID,
		UploadedOrders:       uploadedOrder,
		Status:               OrderStatusSUBMITTED,
		OrdersNumber:         &ordersNumber,
		TAC:                  &TAC,
		SAC:                  &SAC,
		DepartmentIndicator:  &deptIndicator,
		Grade:                swag.String("E-1"),
	}
	suite.MustSave(&order)

	// User is authorized to fetch service member
	session := &auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		ServiceMemberID: serviceMember.ID,
	}

	actualOrder, err := serviceMember.FetchLatestOrder(session, suite.DB())

	if suite.NoError(err) {
		suite.Equal(order.Grade, actualOrder.Grade)
		suite.Equal(order.OriginDutyLocationID, actualOrder.OriginDutyLocationID)
		suite.Equal(order.NewDutyLocationID, actualOrder.NewDutyLocationID)
		suite.True(order.IssueDate.Equal(actualOrder.IssueDate))
		suite.True(order.ReportByDate.Equal(actualOrder.ReportByDate))
		suite.Equal(order.OrdersType, actualOrder.OrdersType)
		suite.Equal(order.HasDependents, actualOrder.HasDependents)
		suite.Equal(order.SpouseHasProGear, actualOrder.SpouseHasProGear)
		suite.Equal(order.UploadedOrdersID, actualOrder.UploadedOrdersID)

	}

	// Wrong ServiceMember
	wrongID, _ := uuid.NewV4()
	_, err = FetchServiceMemberForUser(suite.DB(), session, wrongID)
	if suite.Error(err) {
		suite.Equal(ErrFetchNotFound, err)
	}
}
