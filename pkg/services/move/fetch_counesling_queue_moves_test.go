package move

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	mocks "github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *MoveServiceSuite) TestGetCounselingQueueDBFuncProcess() {
	counselingQueueFetcher := NewCounselingQueueFetcher()

	suite.Run("returns all moves sorted based on its default submitted_at", func() {

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "KKFA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		fetcher := &mocks.OfficeUserGblocFetcher{}
		fetcher.On("FetchGblocForOfficeUser",
			mock.AnythingOfType("*appcontext.appContext"),
			officeUser.ID,
		).Return("KKFA", nil)

		fortGordon := factory.FetchOrBuildOrdersDutyLocation(suite.DB())
		yuma := factory.FetchOrBuildCurrentDutyLocation(suite.DB())

		army := models.AffiliationARMY
		navy := models.AffiliationNAVY

		first1 := "Adam"
		last1 := "Smith"
		first2 := "Will"
		last2 := "Smilmer"
		first3 := "Jason"
		last3 := "Smighler"

		emplid1 := "111188879"
		emplid2 := "111188878"
		emplid3 := "111133333"

		edipi1 := "111145678"
		edipi2 := "111145998"
		edipi3 := "111133333"

		submittedAt1 := time.Date(2022, 05, 01, 0, 0, 0, 0, time.UTC)
		submittedAt2 := time.Date(2022, 05, 01, 0, 0, 0, 0, time.UTC)
		submittedAt3 := time.Date(2022, 05, 01, 0, 0, 0, 0, time.UTC)

		requestedPickupDate1 := time.Date(2022, 07, 01, 0, 0, 0, 0, time.UTC)
		requestedPickupDate2 := time.Date(2022, 07, 01, 0, 0, 0, 0, time.UTC)
		requestedPickupDate3 := time.Date(2022, 07, 01, 0, 0, 0, 0, time.UTC)

		office1 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "JPPSO Testy McTest",
				},
			},
		}, nil)
		office2 := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Name: "PPO Rome Test Office",
				},
			},
		}, nil)

		status := models.OfficeUserStatusAPPROVED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Cam",
					LastName:  "Newton",
					Email:     "camNewton@mail.mil",
					Status:    &status,
					Telephone: "555-555-5555",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		move1 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName:   &first1,
					LastName:    &last1,
					Emplid:      &emplid1,
					Edipi:       &edipi1,
					Affiliation: &army,
				},
			},
			{
				Model: models.Move{
					Locator:            "AAA3T6",
					SubmittedAt:        &submittedAt1,
					SCAssignedID:       &officeUser1.ID,
					CounselingOfficeID: &office1.ID,
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate1,
				},
			},
			{
				Model:    fortGordon,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)
		move2 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName:   &first2,
					LastName:    &last2,
					Affiliation: &army,
					Emplid:      &emplid2,
					Edipi:       &edipi2,
				},
			},
			{
				Model: models.Move{
					Locator:            "AAA3T1",
					SubmittedAt:        &submittedAt2,
					SCAssignedID:       &officeUser1.ID,
					CounselingOfficeID: &office1.ID,
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate2,
				},
			},
			{
				Model:    fortGordon,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)
		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName:   &first3,
					LastName:    &last3,
					Affiliation: &army,
					Emplid:      &emplid3,
					Edipi:       &edipi3,
				},
			},
			{
				Model: models.Move{
					Locator:            "AAA3T0",
					SubmittedAt:        &submittedAt3,
					SCAssignedID:       &officeUser1.ID,
					CounselingOfficeID: &office2.ID,
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate3,
				},
			},
			{
				Model:    fortGordon,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Affiliation: &navy,
				},
			},
			{
				Model: models.Move{
					Locator: "ZZZZZZ",
				},
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate3,
				},
			},
			{
				Model:    yuma,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		// Sort by locator in descending order
		dutyLocationFilter := "fort"
		branchFilter := "Army"
		emplidFilter := "1111"
		edipiFilter := "1111"
		locatorFilter := "AA"
		nameFilter := "Smi"
		submittedAtFilter := time.Date(2022, 05, 01, 0, 0, 0, 0, time.UTC)
		requestedMoveDateFilter := "2022-07-01"
		serviceCounselorFilterName := "New"
		counselingOfficeNameFilter := "JPP"

		sortBy := "Locator"
		desc := "desc"

		ListOrderParams := services.CounselingQueueParams{
			Branch:                 &branchFilter,
			Locator:                &locatorFilter,
			Edipi:                  &edipiFilter,
			Emplid:                 &emplidFilter,
			CustomerName:           &nameFilter,
			OriginDutyLocationName: &dutyLocationFilter,
			SubmittedAt:            &submittedAtFilter,
			RequestedMoveDate:      &requestedMoveDateFilter,
			Sort:                   &sortBy,
			Order:                  &desc,
			CounselingOffice:       &counselingOfficeNameFilter,
			SCAssignedUser:         &serviceCounselorFilterName,
		}

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: officeUser.ID,
		})

		moves, count, err := counselingQueueFetcher.FetchCounselingQueue(appCtx, ListOrderParams)
		suite.NoError(err)
		suite.Equal(int64(2), count)
		suite.Equal(2, len(moves))
		suite.Equal(move1.Locator, moves[0].Locator)
		suite.Equal(move2.Locator, moves[1].Locator)
	})
}
