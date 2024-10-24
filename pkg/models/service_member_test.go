package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	m "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicServiceMemberInstantiation() {
	servicemember := &m.ServiceMember{}

	expErrors := map[string][]string{
		"user_id": {"UserID can not be blank."},
	}

	suite.verifyValidationErrors(servicemember, expErrors)
}

func (suite *ModelSuite) TestIsProfileCompleteWithIncompleteSM() {
	// Given: a user and a service member
	lgu := uuid.Must(uuid.NewV4())
	user1 := m.User{
		OktaID:    lgu.String(),
		OktaEmail: "whoever@example.com",
	}
	verrs, err := user1.Validate(nil)
	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")

	// And: a service member is incompletely initialized with almost all required values
	edipi := "12345567890"
	affiliation := m.AffiliationARMY
	firstName := "bob"
	lastName := "sally"
	telephone := "510 555-5555"
	email := "bobsally@gmail.com"
	fakeAddress := factory.BuildAddress(nil, []factory.Customization{
		{
			Model: m.Address{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
	}, nil)
	fakeBackupAddress := factory.BuildAddress(nil, []factory.Customization{
		{
			Model: m.Address{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
	}, nil)

	serviceMember := m.ServiceMember{
		ID:                     uuid.Must(uuid.NewV4()),
		UserID:                 user1.ID,
		Edipi:                  &edipi,
		Affiliation:            &affiliation,
		FirstName:              &firstName,
		LastName:               &lastName,
		Telephone:              &telephone,
		PersonalEmail:          &email,
		ResidentialAddressID:   &fakeAddress.ID,
		BackupMailingAddressID: &fakeBackupAddress.ID,
	}

	suite.Equal(false, serviceMember.IsProfileComplete())

	// When: all required fields are set
	emailPreferred := true
	serviceMember.EmailIsPreferred = &emailPreferred

	backupContact := factory.BuildBackupContact(nil, []factory.Customization{
		{
			Model:    serviceMember,
			LinkOnly: true,
		},
	}, nil)
	serviceMember.BackupContacts = append(serviceMember.BackupContacts, backupContact)

	suite.Equal(true, serviceMember.IsProfileComplete())
}

func (suite *ModelSuite) TestFetchServiceMemberForUser() {
	user1 := factory.BuildDefaultUser(suite.DB())
	user2 := factory.BuildDefaultUser(suite.DB())

	firstName := "Oliver"
	resAddress := factory.BuildAddress(suite.DB(), nil, nil)
	sm := m.ServiceMember{
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
	goodSm, err := m.FetchServiceMemberForUser(suite.DB(), session, sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddress.ID, goodSm.ResidentialAddress.ID)
	}

	// Wrong ServiceMember
	wrongID, _ := uuid.NewV4()
	_, err = m.FetchServiceMemberForUser(suite.DB(), session, wrongID)
	if suite.Error(err) {
		suite.Equal(m.ErrFetchNotFound, err)
	}

	// User is forbidden from fetching order
	session.UserID = user2.ID
	session.ServiceMemberID = uuid.Nil
	_, err = m.FetchServiceMemberForUser(suite.DB(), session, sm.ID)
	if suite.Error(err) {
		suite.Equal(m.ErrFetchForbidden, err)
	}
}

func (suite *ModelSuite) TestFetchServiceMemberNotForUser() {
	user1 := factory.BuildDefaultUser(suite.DB())

	firstName := "Nino"
	resAddress := factory.BuildAddress(suite.DB(), nil, nil)
	sm := m.ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.MustSave(&sm)

	goodSm, err := m.FetchServiceMember(suite.DB(), sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddressID, goodSm.ResidentialAddressID)
	}
}

func (suite *ModelSuite) TestFetchLatestOrders() {
	setupTestData := func() (m.Order, *auth.Session) {

		user := factory.BuildDefaultUser(suite.DB())

		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
		dutyLocation2 := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
		issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
		reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
		ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		hasDependents := true
		spouseHasProGear := true
		uploadedOrder := m.Document{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		}
		deptIndicator := testdatagen.DefaultDepartmentIndicator
		TAC := testdatagen.DefaultTransportationAccountingCode
		suite.MustSave(&uploadedOrder)
		SAC := "N002214CSW32Y9"
		ordersNumber := "FD4534JFJ"
		contractor := factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
		packingAndShippingInstructions := m.InstructionsBeforeContractNumber + " " + contractor.ContractNumber + " " + m.InstructionsAfterContractNumber

		newGBLOC, gblocErr := models.FetchGBLOCForPostalCode(suite.DB(), dutyLocation2.Address.PostalCode)
		suite.NoError(gblocErr)
		grade := m.ServiceMemberGradeE1
		order := m.Order{
			ServiceMemberID:                serviceMember.ID,
			ServiceMember:                  serviceMember,
			IssueDate:                      issueDate,
			ReportByDate:                   reportByDate,
			OrdersType:                     ordersType,
			HasDependents:                  hasDependents,
			SpouseHasProGear:               spouseHasProGear,
			OriginDutyLocationID:           &dutyLocation.ID,
			OriginDutyLocation:             &dutyLocation,
			NewDutyLocationID:              dutyLocation2.ID,
			NewDutyLocation:                dutyLocation2,
			DestinationGBLOC:               &newGBLOC.GBLOC,
			UploadedOrdersID:               uploadedOrder.ID,
			UploadedOrders:                 uploadedOrder,
			Status:                         m.OrderStatusSUBMITTED,
			OrdersNumber:                   &ordersNumber,
			TAC:                            &TAC,
			SAC:                            &SAC,
			DepartmentIndicator:            &deptIndicator,
			Grade:                          &grade,
			SupplyAndServicesCostEstimate:  m.SupplyAndServicesCostEstimate,
			MethodOfPayment:                m.MethodOfPayment,
			NAICS:                          m.NAICS,
			PackingAndShippingInstructions: packingAndShippingInstructions,
		}
		suite.MustSave(&order)

		// User is authorized to fetch service member
		session := &auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          user.ID,
			ServiceMemberID: serviceMember.ID,
		}
		return order, session
	}

	suite.Run("successfully returns orders with uploads", func() {
		order, session := setupTestData()
		actualOrder, err := m.FetchLatestOrder(session, suite.DB())

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
		_, err = m.FetchServiceMemberForUser(suite.DB(), session, wrongID)
		if suite.Error(err) {
			suite.Equal(m.ErrFetchNotFound, err)
		}
	})

	suite.Run("successfully returns orders without any existing uploads", func() {
		document := factory.BuildDocument(suite.DB(), nil, nil)
		expectedOrder := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    document,
				LinkOnly: true, // if LinkOnly is true, order factory won't build UserUploads
				Type:     &factory.Documents.UploadedOrders,
			},
		}, nil)
		userSession := auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          expectedOrder.ServiceMember.ID,
			ServiceMemberID: expectedOrder.ServiceMemberID,
		}

		actualOrder, err := m.FetchLatestOrder(&userSession, suite.DB())

		suite.NoError(err)
		suite.Equal(expectedOrder.ID, actualOrder.ID)
		suite.Len(actualOrder.UploadedOrders.UserUploads, 0)
	})

	suite.Run("successfully returns non deleted orders and amended orders uploads", func() {
		nonDeletedOrdersUpload := factory.BuildUserUpload(suite.DB(), nil, nil)
		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    nonDeletedOrdersUpload.Document,
				LinkOnly: true,
			},
			{
				Model: m.UserUpload{
					DeletedAt: m.TimePointer(time.Now()),
				},
			},
		}, nil)

		nonDeletedAmendedUpload := factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model: m.UserUpload{
					UploaderID: nonDeletedOrdersUpload.Document.ServiceMember.UserID,
				},
			},
		}, nil)
		factory.BuildUserUpload(suite.DB(), []factory.Customization{
			{
				Model:    nonDeletedAmendedUpload.Document,
				LinkOnly: true,
			},
			{
				Model: m.UserUpload{
					DeletedAt: m.TimePointer(time.Now()),
				},
			},
		}, nil)

		expectedOrder := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    nonDeletedOrdersUpload.Document.ServiceMember,
				LinkOnly: true,
			},
			{
				Model:    nonDeletedOrdersUpload.Document,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedOrders,
			},
			{
				Model:    nonDeletedAmendedUpload.Document,
				LinkOnly: true,
				Type:     &factory.Documents.UploadedAmendedOrders,
			},
		}, nil)

		userSession := auth.Session{
			ApplicationName: auth.MilApp,
			UserID:          expectedOrder.ServiceMember.ID,
			ServiceMemberID: expectedOrder.ServiceMemberID,
		}

		actualOrder, err := m.FetchLatestOrder(&userSession, suite.DB())

		suite.NoError(err)
		suite.Len(actualOrder.UploadedOrders.UserUploads, 1)
		suite.Equal(actualOrder.UploadedOrders.UserUploads[0].ID, nonDeletedOrdersUpload.ID)
		suite.Len(actualOrder.UploadedAmendedOrders.UserUploads, 1)
		suite.Equal(actualOrder.UploadedAmendedOrders.UserUploads[0].ID, nonDeletedAmendedUpload.ID)
	})
}

func (suite *ModelSuite) TestSaveServiceMember() {
	user1 := factory.BuildDefaultUser(suite.DB())

	firstName := "Billy"
	lastName := "Bob"
	sm := m.ServiceMember{
		User:      user1,
		UserID:    user1.ID,
		FirstName: &firstName,
		LastName:  &lastName,
	}
	suite.MustSave(&sm)
	appCtx := suite.AppContextForTest()

	// updating residential address
	resAddress := m.Address{
		StreetAddress1: "987 Other Avenue",
		City:           "Tulsa",
		State:          "OK",
		PostalCode:     "74133",
	}
	sm.ResidentialAddress = &resAddress
	verrs, err := m.SaveServiceMember(appCtx, &sm)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	// updating backup address
	backupAddress := m.Address{
		StreetAddress1: "987 Backup Avenue",
		City:           "Tulsa",
		State:          "OK",
		PostalCode:     "74133",
	}
	sm.BackupMailingAddress = &backupAddress
	verrs, err = m.SaveServiceMember(appCtx, &sm)
	suite.NoError(err)
	suite.False(verrs.HasAny())
}
