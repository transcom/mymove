package move

import (
	"errors"
	"slices"
	"time"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type TestMoves struct {
	// Moves
	defaultMoveWithShipments                  models.Move
	latestMoveWithShipments                   models.Move
	defaultLatestMoveWithShipmentsNeedsSC     models.Move
	safetyMove                                models.Move
	defaultMoveWithShipmentsApproved          models.Move
	defaultLatestMoveWithShipmentsSCCompleted models.Move
	counselingOffice1                         models.TransportationOffice
	assignedOfficeUser                        models.OfficeUser
}

func (suite *MoveServiceSuite) makeCounselingSubtestData() (subtestData *TestMoves) {
	testData := &TestMoves{}

	specifiedTimestampEarliest := time.Date(2022, 04, 01, 0, 0, 0, 0, time.UTC)
	specifiedTimestamp1 := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
	specifiedTimestamp2 := time.Date(2022, 04, 03, 0, 0, 0, 0, time.UTC)
	specifiedTimestamp3 := time.Date(2022, 04, 04, 0, 0, 0, 0, time.UTC)
	specifiedTimestampLatest := time.Date(2022, 04, 05, 0, 0, 0, 0, time.UTC)

	navy := models.AffiliationNAVY

	// Duty Locations
	fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
	yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())

	// Offices
	office1 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name: "JPPSO Random Office 1234",
			},
		},
	}, nil)
	testData.counselingOffice1 = office1

	officeUserApprovedStatus := models.OfficeUserStatusAPPROVED
	officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				FirstName: "Cam",
				LastName:  "Newton",
				Email:     "camNewton@mail.mil",
				Status:    &officeUserApprovedStatus,
				Telephone: "555-555-5555",
			},
		},
	}, []roles.RoleType{roles.RoleTypeServicesCounselor})
	testData.assignedOfficeUser = officeUser1

	edipi1 := "1122334455"
	emplid1 := "1122334455"
	firstName1 := "Grant"

	// Test Moves
	testData.defaultMoveWithShipments = factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{

		{
			Model:    yuma,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Move{
				Locator:                "AAA3T0",
				SubmittedAt:            &specifiedTimestampEarliest,
				Status:                 models.MoveStatusNeedsServiceCounseling,
				SCCounselingAssignedID: &officeUser1.ID,
				CounselingOfficeID:     &office1.ID,
			},
			Type: &factory.Move,
		},
		{
			Model: models.ServiceMember{
				FirstName:   &firstName1,
				Edipi:       &edipi1,
				Emplid:      &emplid1,
				Affiliation: &navy,
			},
		},
		{
			Model: models.MTOShipment{
				RequestedPickupDate: &specifiedTimestampEarliest,
				Status:              models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	testData.safetyMove = factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{

		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Move{
				Locator: "AAA3T3",
				Status:  models.MoveStatusNeedsServiceCounseling,
			},
		},
		{
			Model: models.Order{
				OrdersType: internalmessages.OrdersTypeSAFETY,
			},
		},
	}, nil)

	testData.latestMoveWithShipments = factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{

		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Move{
				Locator:     "AAA3T1",
				SubmittedAt: &specifiedTimestampLatest,
				Status:      models.MoveStatusNeedsServiceCounseling,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	testData.defaultLatestMoveWithShipmentsNeedsSC = factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{

		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Move{
				Locator:     "AAA3T2",
				SubmittedAt: &specifiedTimestamp1,
				Status:      models.MoveStatusNeedsServiceCounseling,
			},
		},
		{
			Model: models.MTOShipment{
				RequestedPickupDate: &specifiedTimestampLatest,
				Status:              models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	testData.defaultLatestMoveWithShipmentsSCCompleted = factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{

		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Move{
				Locator:     "AAA3T4",
				SubmittedAt: &specifiedTimestamp2,
				Status:      models.MoveStatusServiceCounselingCompleted,
			},
		},
	}, nil)

	testData.defaultMoveWithShipmentsApproved = factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{

		{
			Model:    fortGordon,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Move{
				Locator:     "AAA3T5",
				SubmittedAt: &specifiedTimestamp3,
				Status:      models.MoveStatusAPPROVED,
			},
		},
	}, nil)

	return testData
}

func (suite *MoveServiceSuite) TestGetCounselingQueueDBFuncProcess() {

	suite.Run("sorting by requested pickup date returns moves in correct order", func() {

		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		suite.makeCounselingSubtestData()
		sortBy := "requestedPickupDates"
		orderBy := "desc"

		counselingQueueParams := services.CounselingQueueParams{
			Sort:  &sortBy,
			Order: &orderBy,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(3))
		suite.FatalTrue(suite.NotEmpty(returnedMoves, "No moves were found when there should have been 3"))
		suite.Equal("AAA3T1", string(returnedMoves[0].Locator))
		suite.Equal("AAA3T2", string(returnedMoves[1].Locator))
		suite.Equal("AAA3T0", string(returnedMoves[2].Locator))

		// Test that updated_at is returned for move
		suite.NotNil("AAA3T0", returnedMoves[2].UpdatedAt)

		sortBy = "requestedPickupDates"
		orderBy = "asc"
		counselingQueueParams = services.CounselingQueueParams{
			Sort:  &sortBy,
			Order: &orderBy,
		}
		returnedMoves, count, err = counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(3))
		suite.FatalTrue(suite.NotEmpty(returnedMoves, "No moves were found when there should have been 3"))
		suite.Equal("AAA3T0", string(returnedMoves[0].Locator))
		suite.Equal("AAA3T2", string(returnedMoves[1].Locator))
		suite.Equal("AAA3T1", string(returnedMoves[2].Locator))
	})

	suite.Run("sorting by submittedAt returns moves in correct order", func() {

		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		suite.makeCounselingSubtestData()
		sortBy := "submittedAt"
		orderBy := "asc"

		counselingQueueParams := services.CounselingQueueParams{
			Sort:  &sortBy,
			Order: &orderBy,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(3))
		suite.FatalTrue(suite.NotEmpty(returnedMoves, "No moves were found when there should have been 3"))
		suite.Equal("AAA3T0", string(returnedMoves[0].Locator))
		suite.Equal("AAA3T2", string(returnedMoves[1].Locator))
		suite.Equal("AAA3T1", string(returnedMoves[2].Locator))

		sortBy = "submittedAt"
		orderBy = "desc"
		counselingQueueParams = services.CounselingQueueParams{
			Sort:  &sortBy,
			Order: &orderBy,
		}
		returnedMoves, count, err = counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(3))
		suite.FatalTrue(suite.NotEmpty(returnedMoves, "No moves were found when there should have been 3"))
		suite.Equal("AAA3T1", string(returnedMoves[0].Locator))
		suite.Equal("AAA3T2", string(returnedMoves[1].Locator))
		suite.Equal("AAA3T0", string(returnedMoves[2].Locator))
	})

	suite.Run("default sort is by submitted at oldest -> newest", func() {

		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		suite.makeCounselingSubtestData()
		counselingQueueParams := services.CounselingQueueParams{}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(3))
		suite.FatalTrue(suite.NotEmpty(returnedMoves, "No moves were found when there should have been 3"))
		suite.Equal("AAA3T0", string(returnedMoves[0].Locator))
		suite.Equal("AAA3T2", string(returnedMoves[1].Locator))
		suite.Equal("AAA3T1", string(returnedMoves[2].Locator))
	})

	suite.Run("can sort descending showing newest first and oldest last", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		suite.makeCounselingSubtestData()

		desc := "desc"
		counselingQueueParams := services.CounselingQueueParams{
			Order: &desc,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Greater(count, int64(0))
		suite.FatalTrue(suite.NotEmpty(returnedMoves, "No moves were found when there should have been"))
		suite.Equal("AAA3T0", string(returnedMoves[2].Locator))
		suite.Equal("AAA3T1", string(returnedMoves[0].Locator))
	})

	suite.Run("returns moves based on status", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		suite.makeCounselingSubtestData()

		statuses := []string{"NEEDS SERVICE COUNSELING", "SERVICE COUNSELING COMPLETED"}
		counselingQueueParams := services.CounselingQueueParams{
			Status: statuses,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Greater(count, int64(0))
		suite.FatalTrue(suite.NotEmpty(returnedMoves, "No moves were found when there should have been"))
		suite.Equal("AAA3T0", returnedMoves[0].Locator)
		suite.Equal("AAA3T2", returnedMoves[1].Locator)
		suite.Equal("AAA3T4", returnedMoves[2].Locator)
		suite.Equal("AAA3T1", returnedMoves[3].Locator)

		for _, m := range returnedMoves {
			var err error
			if !slices.Contains(statuses, string(m.Status)) {
				err = errors.New("Moves of incorrect status has been returned")
			}
			suite.NoError(err, "Moves of incorrect status has been returned")
		}

	})

	suite.Run("returns moves filtered by customer name", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			CustomerName: testData.defaultMoveWithShipments.Orders.ServiceMember.FirstName,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by edipi", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			Edipi: testData.defaultMoveWithShipments.Orders.ServiceMember.Edipi,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by emplid", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			Emplid: testData.defaultMoveWithShipments.Orders.ServiceMember.Emplid,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by locator", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			Locator: &testData.defaultMoveWithShipments.Locator,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by submitted_at", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			SubmittedAt: testData.defaultMoveWithShipments.SubmittedAt,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by requested pickup date", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		timestamp := "2022-04-01"
		counselingQueueParams := services.CounselingQueueParams{
			RequestedMoveDate: &timestamp,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by branch", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			Branch: (*string)(testData.defaultMoveWithShipments.Orders.ServiceMember.Affiliation),
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by duty location", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			OriginDutyLocationName: &testData.defaultMoveWithShipments.Orders.OriginDutyLocation.Name,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("returns moves filtered by counseling office", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			CounselingOffice: &testData.counselingOffice1.Name,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
	})

	suite.Run("see safety moves if permitted", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
			{
				Model: models.User{
					Privileges: []roles.Privilege{
						{
							PrivilegeType: roles.PrivilegeTypeSupervisor,
						},
						{
							PrivilegeType: roles.PrivilegeTypeSafety,
						},
					},
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypeServicesCounselor,
						},
					},
				},
			},
		}, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		hasSafetyPrivilege := scOfficeUser.User.Privileges.HasPrivilege(roles.PrivilegeTypeSafety)

		counselingQueueParams := services.CounselingQueueParams{
			HasSafetyPrivilege: &hasSafetyPrivilege,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(4))
		suite.Equal(testData.safetyMove.Locator, returnedMoves[3].Locator)
		suite.Equal(internalmessages.OrdersTypeSAFETY, returnedMoves[3].Orders.OrdersType)
	})

	suite.Run("returns moves filtered by assigned office user", func() {
		// Office users
		scOfficeUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: scOfficeUser.ID,
		})

		testData := suite.makeCounselingSubtestData()

		counselingQueueParams := services.CounselingQueueParams{
			SCAssignedUser: &testData.assignedOfficeUser.LastName,
		}

		counselingQueueFetcher := NewCounselingQueueFetcher()
		returnedMoves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, counselingQueueParams)

		suite.FatalNoError(err)
		suite.Equal(count, int64(1))
		suite.Equal(testData.defaultMoveWithShipments.Locator, returnedMoves[0].Locator)
		suite.Equal(testData.defaultMoveWithShipments.SCCounselingAssignedID, returnedMoves[0].SCCounselingAssignedID)
		suite.Equal(testData.defaultMoveWithShipments.SCCounselingAssignedID.String(), returnedMoves[0].SCCounselingAssignedUser.ID.String())
	})
}
