package order

import (
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

	setupServicesCounselor := func(db *pop.Connection) models.OfficeUser {
		return factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
	}

	setupOrdersToFilterBy := func(db *pop.Connection) {
		oldestCloseoutInitiatedDate := time.Date(2022, 04, 02, 0, 0, 0, 0, time.UTC)
		latestCloseoutInitiatedDate := time.Date(2023, 04, 02, 0, 0, 0, 0, time.UTC)

		macDillTransportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "MacDill AFB",
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
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						Locator:          "LATEST",
						CloseoutOfficeID: &macDillTransportationOffice.ID,
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
						PPMType: models.StringPointer(models.MovePPMTypePARTIAL),
						Locator: "OLDEST",
					},
					Type: &factory.Move,
				},
			},
		)
	}

	createUserAndCtx := func(db *pop.Connection) (models.OfficeUser, appcontext.AppContext) {
		servicesCounselor := setupServicesCounselor(db)
		session := setupAuthSession(servicesCounselor.ID)
		appCtx := suite.AppContextWithSessionForTest(&session)
		return servicesCounselor, appCtx
	}

	suite.Run("default sort is by closeout initiated oldest -> newest", func() {
		setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

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

	suite.Run("can sort descending showing newest first and oldest last", func() {
		setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

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
		setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

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
		setupOrdersToFilterBy(suite.DB())
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		suite.Equal(filteredMoves[0].Locator, ppmShipment.Shipment.MoveTaskOrder.Locator)
		suite.Equal(count, 1)
		suite.FatalNoError(err)
	})

	suite.Run("can filter by destination", func() {
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

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
		servicesCounselor, appCtx := createUserAndCtx(suite.DB())

		now := time.Now()
		army := models.AffiliationARMY
		navy := models.AffiliationNAVY
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
		serviceMemberArmy := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		serviceMemberNavy := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &navy,
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

		ppmShipmentNavy := factory.BuildPPMShipmentThatNeedsCloseout(
			suite.DB(),
			nil,
			[]factory.Customization{
				{
					Model: models.Order{
						ServiceMemberID: serviceMemberNavy.ID,
					},
				},
				{
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &now,
						Locator:          "NAVYNA",
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
		// Army should be the first move
		suite.Equal(filteredMoves[0].Locator, ppmShipmentArmy.Shipment.MoveTaskOrder.Locator, "Army move not sorted properly it should be the first in the slice")
		suite.Equal(filteredMoves[1].Locator, ppmShipmentNavy.Shipment.MoveTaskOrder.Locator, "The navy move should be the second in the slice")
	})

}
