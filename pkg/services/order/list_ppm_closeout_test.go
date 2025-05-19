package order

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
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

	setupAuthSession := func(officeUserID uuid.UUID) auth.Session {
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			UserID:          officeUserID,
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
					Model: models.Move{
						PPMType:          models.StringPointer(models.MovePPMTypeFULL),
						SubmittedAt:      &latestCloseoutInitiatedDate,
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
					Model: models.Move{
						PPMType:     models.StringPointer(models.MovePPMTypePARTIAL),
						SubmittedAt: &oldestCloseoutInitiatedDate,
						Locator:     "OLDEST",
					},
					Type: &factory.Move,
				},
			},
		)
	}

	suite.Run("default sort is by closeout initiated oldest -> newest", func() {
		setupOrdersToFilterBy(suite.DB())
		servicesCounselor := setupServicesCounselor(suite.DB())
		session := setupAuthSession(servicesCounselor.ID)
		appCtx := suite.AppContextWithSessionForTest(&session)

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
		servicesCounselor := setupServicesCounselor(suite.DB())
		session := setupAuthSession(servicesCounselor.ID)
		appCtx := suite.AppContextWithSessionForTest(&session)

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
		servicesCounselor := setupServicesCounselor(suite.DB())
		session := setupAuthSession(servicesCounselor.ID)
		appCtx := suite.AppContextWithSessionForTest(&session)

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
		servicesCounselor := setupServicesCounselor(suite.DB())
		session := setupAuthSession(servicesCounselor.ID)
		appCtx := suite.AppContextWithSessionForTest(&session)

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
}
