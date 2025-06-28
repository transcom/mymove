package order

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
)

// We want to find ppm_shipments tied to mto_shipments given a status of
// models.PPMShipmentStatusNeedsCloseout/ ""
func (suite *OrderServiceSuite) TestListPPMCloseoutOrders() {
	waf := entitlements.NewWeightAllotmentFetcher()
	orderFetcher := NewOrderFetcher(waf)

	setupAuthSession := func(userID uuid.UUID) auth.Session {
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          userID,
			IDToken:         "token",
			Hostname:        handlers.OfficeTestHost,
			Email:           "deactivated@example.com",
		}
		return session
	}

	setupServicesCounselor := func(db *pop.Connection, transportationOffice models.TransportationOffice) models.OfficeUser {
		// Create the SC with role and assigned transportation office
		return factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					TransportationOfficeID: transportationOffice.ID,
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})
	}

	// Return MacDill AFB as a "SHARED" default GBLOC office
	setupOrdersToFilterBy := func(db *pop.Connection) models.TransportationOffice {
		oldestCloseoutInitiatedDate := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		latestCloseoutInitiatedDate := time.Date(2023, 04, 02, 0, 0, 0, 0, time.UTC)

		servicesCounselorAssignedForCloseoutForBothMoves := factory.BuildOfficeUserWithRoles(db, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

		macDillTransportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "MacDill AFB",
					Gbloc: "SHARED",
				},
			},
		}, nil)

		patrickTransportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "Patrick SFB",
					Gbloc: "SHARED",
				},
			},
		}, nil)

		// Full PPM for closeout
		// Latest
		// MacDill AFB
		factory.BuildPPMShipmentThatNeedsCloseout(
			db,
			nil,
			[]factory.Customization{
				{
					Model: models.PPMShipment{
						SubmittedAt: &latestCloseoutInitiatedDate,
					},
				},
				{
					Model: models.Move{
						PPMType:              models.StringPointer(models.MovePPMTypeFULL),
						Locator:              "LATEST",
						CloseoutOfficeID:     &macDillTransportationOffice.ID,
						CounselingOfficeID:   &macDillTransportationOffice.ID,
						SCCloseoutAssignedID: &servicesCounselorAssignedForCloseoutForBothMoves.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// Partial PPM for closeout
		// Oldest
		factory.BuildPPMShipmentThatNeedsCloseout(
			db,
			nil,
			[]factory.Customization{
				{
					Model: models.PPMShipment{
						SubmittedAt: &oldestCloseoutInitiatedDate,
					},
				},
				{
					Model: models.Move{
						PPMType:              models.StringPointer(models.MovePPMTypePARTIAL),
						Locator:              "OLDEST",
						CloseoutOfficeID:     &patrickTransportationOffice.ID,
						CounselingOfficeID:   &patrickTransportationOffice.ID,
						SCCloseoutAssignedID: &servicesCounselorAssignedForCloseoutForBothMoves.ID,
					},
					Type: &factory.Move,
				},
			},
		)
		return macDillTransportationOffice
	}

	createUserAndCtx := func(db *pop.Connection, transportationOffice models.TransportationOffice) (models.OfficeUser, appcontext.AppContext) {
		servicesCounselor := setupServicesCounselor(db, transportationOffice)
		session := setupAuthSession(servicesCounselor.ID)
		appCtx := suite.AppContextWithSessionForTest(&session)
		return servicesCounselor, appCtx
	}

	suite.Run("default sort is by closeout initiated oldest -> newest", func() {
		defaultOffice := setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), defaultOffice)

		defaultMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{NeedsPPMCloseout: models.BoolPointer(true)},
		)
		suite.FatalNoError(err)
		suite.Equal(count, 2)
		suite.FatalTrue(suite.NotEmpty(defaultMoves, "No moves were found when there should be 2"))
		suite.Equal("OLDEST", defaultMoves[0].Locator)
		suite.Equal("LATEST", defaultMoves[1].Locator)
	})

	suite.Run("jsonb mapping from the db func holds the ID values", func() {
		defaultOffice := setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), defaultOffice)

		defaultMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{NeedsPPMCloseout: models.BoolPointer(true)},
		)
		suite.FatalNoError(err)
		suite.Equal(count, 2)
		suite.FatalTrue(suite.NotEmpty(defaultMoves, "No moves were found when there should be 2"))
		for _, move := range defaultMoves {
			suite.True(move.CloseoutOfficeID != nil && *move.CloseoutOfficeID != uuid.Nil)
			suite.True(move.OrdersID != uuid.Nil)
			suite.True(move.CounselingOfficeID != nil && *move.CounselingOfficeID != uuid.Nil)
			suite.True(move.SCCloseoutAssignedID != nil && *move.CloseoutOfficeID != uuid.Nil)
		}
	})

	suite.Run("locked by and updated_at information comes back with the move", func() {
		macDillTransportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:  "MacDill AFB",
					Gbloc: "SHARED",
				},
			},
		}, nil)
		sm := factory.BuildServiceMember(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), macDillTransportationOffice)

		now := time.Now()
		factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(), nil,
			[]factory.Customization{
				{
					Model: models.Order{ServiceMemberID: sm.ID},
				},
				{
					Model: models.Move{
						Locator:              "LOCKED",
						SubmittedAt:          &now,
						LockedByOfficeUserID: models.UUIDPointer(servicesCounselor.ID),
						LockExpiresAt:        &now,
						CloseoutOfficeID:     models.UUIDPointer(macDillTransportationOffice.ID),
					},
					Type: &factory.Move,
				},
			},
		)

		defaultMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{NeedsPPMCloseout: models.BoolPointer(true), Locator: models.StringPointer("LOCKED")},
		)
		suite.FatalNoError(err)
		suite.Equal(1, count)
		suite.NotNil(defaultMoves[0].LockExpiresAt)
		suite.NotNil(defaultMoves[0].LockedByOfficeUserID)
		suite.NotEmpty(defaultMoves[0].UpdatedAt)
	})

	suite.Run("can sort descending showing newest first and oldest last", func() {
		defaultOffice := setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), defaultOffice)

		sortedDescendingMoves, _, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{Sort: models.StringPointer("closeoutInitiated"), Order: models.StringPointer("desc"), NeedsPPMCloseout: models.BoolPointer(true)},
		)

		suite.FatalNoError(err)
		suite.Equal("LATEST", sortedDescendingMoves[0].Locator)
		suite.Equal("OLDEST", sortedDescendingMoves[1].Locator)
	})

	suite.Run("returns moves based on ppm status", func() {
		defaultOffice := setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), defaultOffice)

		partialMoves, count, err := orderFetcher.ListPPMCloseoutOrders(appCtx, servicesCounselor.ID, &services.ListOrderParams{
			PPMType: models.StringPointer(models.MovePPMTypePARTIAL),
		})
		suite.Equal(count, 1)
		suite.NoError(err)
		suite.Len(partialMoves, 1, "Test setup should have create a single partial PPM and it should have been found here")

		fullMoves, count, err := orderFetcher.ListPPMCloseoutOrders(appCtx, servicesCounselor.ID, &services.ListOrderParams{
			PPMType: models.StringPointer(models.MovePPMTypeFULL),
		})
		suite.Equal(count, 1)
		suite.FatalNoError(err)
		suite.Len(fullMoves, 1, "Test setup should have create a single full PPM and it should have been found here")
	})

	suite.Run("can filter by closeout location inclusively", func() {
		defaultOffice := setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), defaultOffice)

		movesWithMacDillFilter, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				NeedsPPMCloseout: models.BoolPointer(true),
				CloseoutLocation: models.StringPointer("dill"),
			},
		)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
		suite.Equal("LATEST", movesWithMacDillFilter[0].Locator)
		suite.Equal("MacDill AFB", movesWithMacDillFilter[0].CloseoutOffice.Name)
	})

	suite.Run("can filter by customer name", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				CustomerName: models.StringPointer(*ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.FirstName),
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by edipi", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()

		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				Edipi: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.Edipi,
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by emplid", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				Edipi: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.Emplid,
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by locator", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				Locator: &ppmShipment.Shipment.MoveTaskOrder.Locator,
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by submitted_at", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				CloseoutInitiated: ppmShipment.SubmittedAt,
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by branch", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				Branch: (*string)(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.Affiliation),
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by move's ppm type", func() {
		// Don't confuse this one with ppm_shipments.ppm_type!
		// They are different! Filtering by ppm shipment ppm type is not supported
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				PPMType: models.StringPointer(string(models.MovePPMTypeFULL)),
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by origin duty location name", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		dutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)

		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						OriginDutyLocationID: &dutyLocation.ID,
					},
				},
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				OriginDutyLocation: []string{dutyLocation.Name},
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by counseling office", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:            models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:        &now,
						Locator:            "LATEST",
						CounselingOfficeID: &transportationOffice.ID,
						CloseoutOfficeID:   &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				CounselingOffice: &transportationOffice.Name,
			},
		)
		suite.FatalNoError(err)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
	})

	suite.Run("can filter by destination", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		dutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						NewDutyLocationID: dutyLocation.ID,
					},
				},
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				DestinationDutyLocation: &dutyLocation.Name,
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by closeout office", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				CloseoutLocation: &transportationOffice.Name,
			},
		)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("see safety moves if permitted", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		privilegedUser := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email:  "officeuser1@example.com",
					Active: true,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
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
							RoleType: roles.RoleTypeTOO,
						},
					},
				},
			},
		}, nil)

		session := setupAuthSession(*privilegedUser.UserID)
		appCtx := suite.AppContextWithSessionForTest(&session)

		now := time.Now()

		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						OrdersType: internalmessages.OrdersTypeSAFETY,
					},
				},
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			privilegedUser.ID,
			&services.ListOrderParams{
				Locator: &ppmShipment.Shipment.MoveTaskOrder.Locator,
			},
		)

		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can not see safety moves by default", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						OrdersType: internalmessages.OrdersTypeSAFETY,
					},
				},
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				Locator: &ppmShipment.Shipment.MoveTaskOrder.Locator,
			},
		)

		suite.Len(filteredMoves, 0)
		suite.Equal(count, 0)
		suite.FatalNoError(err)
	})

	suite.Run("can not see a move entry that has a ppm but the ppm has not entered closeout", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		// By not setting a submittedAt value for the PPM shipment, we are declaring
		// it has not entered the closeout phase
		ppmShipment := factory.BuildPPMShipment(
			suite.DB(),
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "LATEST",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
			nil,
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				Locator: &ppmShipment.Shipment.MoveTaskOrder.Locator,
			},
		)
		suite.NoError(err)
		suite.Equal(0, count)
		suite.Equal(0, len(filteredMoves))
	})

	suite.Run("throws an err if the wrong submitted at param is provided", func() {
		now := time.Now()
		_, _, err := orderFetcher.ListPPMCloseoutOrders(
			suite.AppContextForTest(),
			uuid.Must(uuid.NewV4()),
			&services.ListOrderParams{
				SubmittedAt: &now,
			},
		)
		suite.ErrorContains(err, "submitted at parameter should not be used for PPM closeout queue. Please use closeout initiated instead")
	})

	suite.Run("can filter by sc assigned closeout", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						PPMType:              models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:          &now,
						Locator:              "LATEST",
						CloseoutOfficeID:     &transportationOffice.ID,
						SCCloseoutAssignedID: models.UUIDPointer(servicesCounselor.ID),
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				AssignedTo: &servicesCounselor.FirstName,
			},
		)
		suite.Equal(count, 1)
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.FatalNoError(err)
	})

	suite.Run("can sort by branch affiliation", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		servicesCounselor, appCtx := createUserAndCtx(suite.DB(), transportationOffice)

		now := time.Now()
		army := models.AffiliationARMY
		airForce := models.AffiliationAIRFORCE
		serviceMemberArmy := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		serviceMemberAirforce := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &airForce,
				},
			},
		}, nil)

		ppmShipmentArmy := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						ServiceMemberID: serviceMemberArmy.ID,
					},
				},
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "ARMYAR",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		ppmShipmentAirforce := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						ServiceMemberID: serviceMemberAirforce.ID,
					},
				},
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "AIRFOR",
						CloseoutOfficeID: &transportationOffice.ID,
					},
					Type: &factory.Move,
				},
			},
		)

		// The factory should always return this information
		suite.NotEmpty(ppmShipmentArmy.Shipment.MoveTaskOrder.Orders.ServiceMember)

		filteredMoves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{
				Sort: models.StringPointer("branch"),
			},
		)
		suite.Equal(count, 2)
		suite.FatalNoError(err)
		// Air force should be the first move
		suite.Equal(filteredMoves[0].Locator, ppmShipmentAirforce.Shipment.MoveTaskOrder.Locator, "Air force move not sorted properly it should be the first in the slice")
		suite.Equal(filteredMoves[1].Locator, ppmShipmentArmy.Shipment.MoveTaskOrder.Locator, "The ARmy move should be the second in the slice")
	})

	suite.Run("GBLOC results should be driven by closeout location, not origin location", func() {
		closeoutOfficeAssignedToUser := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:             "Will appear",
					Gbloc:            "FIND",
					ProvidesCloseout: true,
				},
			},
		}, nil)

		closeoutOfficeForAnotherGBLOC := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name:             "Will not appear",
					Gbloc:            "HIDE",
					ProvidesCloseout: true,
				},
			},
		}, nil)

		servicesCounselor := factory.BuildOfficeUserWithRoles(
			suite.DB(),
			[]factory.Customization{
				{
					Model: models.OfficeUser{
						TransportationOfficeID: closeoutOfficeAssignedToUser.ID,
					},
				},
				{
					Model:    closeoutOfficeAssignedToUser,
					LinkOnly: true,
					Type:     &factory.TransportationOffices.CounselingOffice,
				},
			},
			[]roles.RoleType{roles.RoleTypeServicesCounselor},
		)

		session := setupAuthSession(servicesCounselor.ID)
		appCtx := suite.AppContextWithSessionForTest(&session)
		now := time.Now()

		// Move we will find
		factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Move{
						Locator:          "MATCHY",
						CloseoutOfficeID: &closeoutOfficeAssignedToUser.ID,
						SubmittedAt:      &now,
					},
					Type: &factory.Move,
				},
			},
		)

		// Move that will exist but be filtered out
		originDL := factory.BuildDutyLocation(suite.DB(), nil, nil)
		factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						OriginDutyLocationID: &originDL.ID,
					},
				},
				{
					Model: models.Move{
						Locator:          "NOMATC",
						CloseoutOfficeID: &closeoutOfficeForAnotherGBLOC.ID,
						SubmittedAt:      &now,
					},
					Type: &factory.Move,
				},
			},
		)

		moves, count, err := orderFetcher.ListPPMCloseoutOrders(
			appCtx,
			servicesCounselor.ID,
			&services.ListOrderParams{NeedsPPMCloseout: models.BoolPointer(true)},
		)
		suite.Equal(count, 1)
		suite.Equal("MATCHY", moves[0].Locator)
		suite.FatalNoError(err)
	})

	suite.Run("Special filtering (NAVY, MARINES (TVCB), USCG)", func() {
		testCases := []struct {
			name                string
			userGBLOC           string
			closeoutGBLOC       string
			affiliationShow     models.ServiceMemberAffiliation
			affiliationDontShow models.ServiceMemberAffiliation
		}{
			{
				name:                "NAVY counselor sees only NAVY moves",
				userGBLOC:           "NAVY",
				affiliationShow:     models.AffiliationNAVY,
				affiliationDontShow: models.AffiliationARMY,
				closeoutGBLOC:       "NAVY",
			},
			{
				name:                "TVCB counselor sees only moves for MARINES",
				userGBLOC:           "TVCB",
				affiliationShow:     models.AffiliationMARINES,
				affiliationDontShow: models.AffiliationARMY,
				closeoutGBLOC:       "TVCB",
			},
			{
				name:                "USCG counselor sees only COAST GUARD moves",
				userGBLOC:           "USCG",
				affiliationShow:     models.AffiliationCOASTGUARD,
				affiliationDontShow: models.AffiliationARMY,
				closeoutGBLOC:       "USCG",
			},
		}

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				// Closeout office the SC will belong to
				closeoutTO := factory.BuildTransportationOffice(
					suite.DB(),
					[]factory.Customization{
						{
							Model: models.TransportationOffice{
								Name:  fmt.Sprintf("%s closeout", tc.userGBLOC),
								Gbloc: tc.userGBLOC, ProvidesCloseout: true,
							},
						},
					}, nil)

				// Create counselor tied to office
				sc := factory.BuildOfficeUserWithRoles(
					suite.DB(),
					[]factory.Customization{
						{
							Model: models.OfficeUser{
								TransportationOfficeID: closeoutTO.ID,
							},
						},
						{
							Model: closeoutTO, LinkOnly: true,
							Type: &factory.TransportationOffices.CounselingOffice,
						},
					},
					[]roles.RoleType{roles.RoleTypeServicesCounselor},
				)

				session := setupAuthSession(sc.ID)
				appCtx := suite.AppContextWithSessionForTest(&session)

				// Helper func to build PPMs needing closeout
				buildMove := func(aff models.ServiceMemberAffiliation, locator, gbloc string) {
					// Build service member with test case affiliation
					sm := factory.BuildServiceMember(suite.DB(),
						[]factory.Customization{{Model: models.ServiceMember{Affiliation: &aff}}}, nil)

					closeout := factory.BuildTransportationOffice(
						suite.DB(),
						[]factory.Customization{
							{
								Model: models.TransportationOffice{
									Name: locator + " closeout", Gbloc: gbloc, ProvidesCloseout: true,
								},
							},
						}, nil)

					now := time.Now()
					factory.BuildPPMShipmentThatNeedsCloseout(
						suite.DB(), nil,
						[]factory.Customization{
							{
								Model: models.Order{ServiceMemberID: sm.ID},
							},
							{
								Model: models.Move{
									Locator: locator, CloseoutOfficeID: &closeout.ID,
									SubmittedAt: &now,
								},
								Type: &factory.Move,
							},
						},
					)
				}

				// Move to be filtered in
				buildMove(tc.affiliationShow, tc.userGBLOC+"YE", tc.closeoutGBLOC)
				// Move to be exist but be filtered out
				buildMove(tc.affiliationDontShow, tc.userGBLOC+"NO", tc.closeoutGBLOC)

				// Fetch
				moves, count, err := orderFetcher.ListPPMCloseoutOrders(
					appCtx, sc.ID,
					&services.ListOrderParams{NeedsPPMCloseout: models.BoolPointer(true)},
				)
				suite.FatalNoError(err)
				suite.Equal(1, count, "unexpected number of moves visible")
				suite.Equal(tc.userGBLOC+"YE", moves[0].Locator)
			})
		}
	})

	suite.Run("Branch filter prevents leaks", func() {
		// This test is used because service members affiliated with:
		// USMC/NAVY/USCG should only ever appear in their designated, special
		// PPM closeout transportation offices. These branches can not appear
		// for any other GBLOC besides TVCB/NAVY/USCG
		buildMove := func() *models.Move {
			navy := models.AffiliationNAVY

			sm := factory.BuildServiceMember(
				suite.DB(),
				[]factory.Customization{{Model: models.ServiceMember{Affiliation: &navy}}},
				nil,
			)

			closeoutTO := factory.BuildTransportationOffice(
				suite.DB(),
				[]factory.Customization{{
					Model: models.TransportationOffice{
						Name: "HAFC closeout", Gbloc: "HAFC", ProvidesCloseout: true,
					}}},
				nil,
			)

			now := time.Now()
			ppm := factory.BuildPPMShipmentThatNeedsCloseout(
				suite.DB(), nil,
				[]factory.Customization{
					{Model: models.Order{ServiceMemberID: sm.ID}},
					{Model: models.Move{
						Locator:          "LOCATO",
						CloseoutOfficeID: &closeoutTO.ID,
						SubmittedAt:      &now,
					}},
				},
			)
			return &ppm.Shipment.MoveTaskOrder
		}

		testCases := []struct {
			name          string
			userGBLOC     string
			expectVisible bool
		}{
			{
				name:          "HAFC counselor should NOT see NAVY move even though GBLOC matches", // GBLOC should NEVER be assigned in the first place, but this covers an edge case
				userGBLOC:     "HAFC",
				expectVisible: false,
			},
			{
				name:          "NAVY counselor still sees NAVY move even with HAFC GBLOC", // Should also never happen, but just covering the edge case
				userGBLOC:     "NAVY",
				expectVisible: true,
			},
		}

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				move := buildMove()

				to := factory.BuildTransportationOffice(
					suite.DB(),
					[]factory.Customization{{
						Model: models.TransportationOffice{
							Name:  fmt.Sprintf("%s closeout", tc.userGBLOC),
							Gbloc: tc.userGBLOC, ProvidesCloseout: true,
						}}},
					nil)

				user := factory.BuildOfficeUserWithRoles(
					suite.DB(),
					[]factory.Customization{{
						Model: models.OfficeUser{TransportationOfficeID: to.ID},
					}},
					[]roles.RoleType{roles.RoleTypeServicesCounselor},
				)

				session := setupAuthSession(user.ID)
				appCtx := suite.AppContextWithSessionForTest(&session)

				moves, count, err := orderFetcher.ListPPMCloseoutOrders(
					appCtx, user.ID,
					&services.ListOrderParams{NeedsPPMCloseout: models.BoolPointer(true)},
				)
				suite.NoError(err)

				if tc.expectVisible {
					suite.Equal(1, count, "move should be visible")
					suite.Equal(move.Locator, moves[0].Locator)
				} else {
					suite.Equal(0, count, "move should NOT be visible")
				}
			})
		}
	})

}
