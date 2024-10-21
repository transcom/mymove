package paymentrequest

import (
	"sort"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListbyMove() {
	suite.Run("Only returns visible (where Move.Show is not false) payment requests", func() {
		paymentRequestListFetcher := NewPaymentRequestListFetcher()

		expectedMove := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "ABC123",
				},
			},
		}, nil)
		// We need a payment request with a move that has a shipment
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove,
				LinkOnly: true,
			},
		}, nil)

		// Hidden move should not be returned
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: models.BoolPointer(false),
				},
			},
		}, nil)

		expectedPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestListByMove(suite.AppContextForTest(), "ABC123")

		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequestsToCheck := *expectedPaymentRequests
		suite.Equal(paymentRequest.ID, paymentRequestsToCheck[0].ID)
	})

}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestList() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	var officeUser models.OfficeUser
	var expectedMove models.Move
	var paymentRequest models.PaymentRequest

	var session auth.Session

	suite.PreloadData(func() {
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		expectedMove = factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		session = auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		// We need a payment request with a move that has a shipment that's within the GBLOC
		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "ABCD",
				},
				Type: &factory.TransportationOffices.OriginDutyLocation,
			},
			{
				Model: models.DutyLocation{
					Name: "KJKJKJKJKJK",
				},
				Type: &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		// Hidden move should not be returned
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Show: models.BoolPointer(false),
				},
			},
		}, nil)
		// Marine Corps payment requests should be excluded even if in the same GBLOC
		marines := models.AffiliationMARINES
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
				},
			},
			{
				Model: models.TransportationOffice{
					Gbloc: "LKNQ",
					ID:    uuid.Must(uuid.NewV4()),
				},
				Type: &factory.TransportationOffices.OriginDutyLocation,
			},
			{
				Model: models.ServiceMember{Affiliation: &marines},
			},
		}, nil)
	})

	suite.Run("Only returns visible (where Move.Show is not false) payment requests matching office user GBLOC", func() {
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2)})

		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))

		paymentRequestsForComparison := *expectedPaymentRequests
		suite.Equal(paymentRequest.ID, paymentRequestsForComparison[0].ID)
	})

	suite.Run("Returns payment request matching an arbitrary filter", func() {
		// Locator
		locator := paymentRequest.MoveTaskOrder.Locator
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2), Locator: &locator})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests := *expectedPaymentRequests
		suite.Equal(paymentRequest.MoveTaskOrder.Locator, paymentRequests[0].MoveTaskOrder.Locator)

		// Branch
		serviceMember := paymentRequest.MoveTaskOrder.Orders.ServiceMember
		affiliation := models.AffiliationAIRFORCE
		serviceMember.Affiliation = &affiliation
		err = suite.DB().Save(&serviceMember)
		suite.NoError(err)

		branch := serviceMember.Affiliation.String()
		expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2), Branch: &branch})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests = *expectedPaymentRequests
		suite.Equal(models.AffiliationAIRFORCE, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
	})
	locationName := paymentRequest.MoveTaskOrder.Orders.OriginDutyLocation.Name
	suite.Run("Returns payment request matching a full originDutyLocation name filter", func() {

		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2), OriginDutyLocation: &locationName})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests := *expectedPaymentRequests
		suite.Equal(locationName, paymentRequests[0].MoveTaskOrder.Orders.OriginDutyLocation.Name)

	})
	suite.Run("Returns payment request matching a partial originDutyLocation filter", func() {
		//Split the location name and retrieve a substring (first string) for the search param
		partialParamSearch := strings.Split(locationName, " ")[0]
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2), OriginDutyLocation: &partialParamSearch})
		suite.NoError(err)
		suite.Equal(1, len(*expectedPaymentRequests))
		paymentRequests := *expectedPaymentRequests
		suite.Equal(locationName, paymentRequests[0].MoveTaskOrder.Orders.OriginDutyLocation.Name)

	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListStatusFilter() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	var officeUser models.OfficeUser
	var allPaymentRequests models.PaymentRequests
	var pendingPaymentRequest, reviewedPaymentRequest, sentToGexPaymentRequest, recByGexPaymentRequest, rejectedPaymentRequest, paidPaymentRequest, deprecatedPaymentRequest, errorPaymentRequest models.PaymentRequest

	var session auth.Session

	suite.PreloadData(func() {
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		session = auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		expectedMove1 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		expectedMove2 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		expectedMove3 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		expectedMove4 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		expectedMove5 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		expectedMove6 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		expectedMove7 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		reviewedPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove1,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		rejectedPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove2,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewedAllRejected,
				},
			},
		}, nil)

		sentToGexPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove3,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusSentToGex,
				},
			},
		}, nil)
		recByGexPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove4,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusTppsReceived,
				},
			},
		}, nil)
		paidPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove5,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPaid,
				},
			},
		}, nil)

		pendingPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove6,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPending,
				},
			},
		}, nil)

		deprecatedPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove6,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status:         models.PaymentRequestStatusDeprecated,
					IsFinal:        false,
					SequenceNumber: 2,
				},
			},
		}, nil)

		errorPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMove7,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusEDIError,
				},
			},
		}, nil)
		allPaymentRequests = []models.PaymentRequest{pendingPaymentRequest, reviewedPaymentRequest, rejectedPaymentRequest, sentToGexPaymentRequest, recByGexPaymentRequest, paidPaymentRequest, deprecatedPaymentRequest, errorPaymentRequest}
	})

	suite.Run("Returns all payment requests when no status filter is specified", func() {
		_, actualCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{})
		suite.NoError(err)
		suite.Equal(len(allPaymentRequests), actualCount)
	})

	suite.Run("Returns all payment requests when all status filters are selected", func() {
		_, actualCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Payment requested", "Reviewed", "Rejected", "Paid", "Deprecated", "Error"}})
		suite.NoError(err)
		suite.Equal(len(allPaymentRequests), actualCount)
	})

	suite.Run("Returns only those payment requests with the exact status", func() {
		pendingPaymentRequests, pendingCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Payment requested"}})
		pending := *pendingPaymentRequests
		suite.NoError(err)
		suite.Equal(1, pendingCount)
		suite.Equal(pendingPaymentRequest.ID, pending[0].ID)

		reviewedPaymentRequests, reviewedCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Reviewed"}})
		reviewed := *reviewedPaymentRequests
		suite.NoError(err)
		suite.Equal(3, reviewedCount)

		reviewedIDs := []uuid.UUID{reviewedPaymentRequest.ID, sentToGexPaymentRequest.ID, recByGexPaymentRequest.ID}
		for _, pr := range reviewed {
			suite.Contains(reviewedIDs, pr.ID)
		}

		rejectedPaymentRequests, rejectedCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Rejected"}})
		rejected := *rejectedPaymentRequests
		suite.NoError(err)
		suite.Equal(1, rejectedCount)
		suite.Equal(rejectedPaymentRequest.ID, rejected[0].ID)

		paidPaymentRequests, paidCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Paid"}})
		paid := *paidPaymentRequests
		suite.NoError(err)
		suite.Equal(1, paidCount)
		suite.Equal(paidPaymentRequest.ID, paid[0].ID)

		deprecatedPaymentRequests, deprecatedCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Deprecated"}})

		deprecated := *deprecatedPaymentRequests
		suite.NoError(err)
		suite.Equal(1, deprecatedCount)
		suite.Equal(deprecatedPaymentRequest.ID, deprecated[0].ID)

		errorPaymentRequests, errorCount, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Status: []string{"Error"}})

		errorPR := *errorPaymentRequests
		suite.NoError(err)
		suite.Equal(1, errorCount)
		suite.Equal(errorPaymentRequest.ID, errorPR[0].ID)
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListUSMCGBLOC() {
	var officeUser, officeUserUSMC models.OfficeUser
	var paymentRequestUSMC, paymentRequestUSMC2 models.PaymentRequest

	var session auth.Session

	suite.PreloadData(func() {
		officeUUID, _ := uuid.NewV4()
		marines := models.AffiliationMARINES
		army := models.AffiliationARMY
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		session = auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		expectedMoveNotUSMC := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		paymentRequestUSMC = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "LKNQ",
					ID:    officeUUID,
				},
				Type: &factory.TransportationOffices.OriginDutyLocation,
			},
			{
				Model: models.Move{
					Status: models.MoveStatusSUBMITTED,
				},
			},
			{
				Model: models.ServiceMember{Affiliation: &marines},
			},
		}, nil)

		paymentRequestUSMC2 = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					SequenceNumber: 2,
				},
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model: models.TransportationOffice{
					Gbloc: "LKNQ",
					ID:    officeUUID,
				},
				Type: &factory.TransportationOffices.OriginDutyLocation,
			},
			{
				Model:    paymentRequestUSMC.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.ServiceMember{Affiliation: &marines},
			},
		}, nil)

		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    expectedMoveNotUSMC,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusPending,
				},
			},
			{
				Model: models.ServiceMember{Affiliation: &army},
			},
		}, nil)

		tioRole := factory.FetchOrBuildRoleByRoleType(suite.DB(), roles.RoleTypeTIO)
		tooRole := factory.FetchOrBuildRoleByRoleType(suite.DB(), roles.RoleTypeTOO)
		officeUserUSMC = factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "USMC",
				},
			},
			{
				Model: models.User{
					Roles: []roles.Role{tioRole, tooRole},
				},
			},
		}, nil)
	})

	suite.Run("returns USMC payment requests", func() {
		paymentRequestListFetcher := NewPaymentRequestListFetcher()
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUserUSMC.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2)})
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(models.AffiliationMARINES, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
		suite.Equal(models.AffiliationMARINES, *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Affiliation)
		expectedPaymentRequests, _, err = paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2)})
		paymentRequests = *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(1, len(paymentRequests))
		suite.Equal(models.AffiliationARMY, *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation)
	})

	suite.Run("returns USMC payment requests for move", func() {
		paymentRequestListFetcher := NewPaymentRequestListFetcher()
		expectedPaymentRequests, err := paymentRequestListFetcher.FetchPaymentRequestListByMove(suite.AppContextWithSessionForTest(&session), paymentRequestUSMC.MoveTaskOrder.Locator)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(paymentRequestUSMC.ID, paymentRequests[0].ID)
		suite.Equal(paymentRequestUSMC2.ID, paymentRequests[1].ID)
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListNoGBLOCMatch() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()

	suite.Run("No results when GBLOC does not match", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "EFGH",
				},
				Type: &factory.TransportationOffices.OriginDutyLocation,
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{
					Gbloc: "ABCD",
				},
				Type: &factory.TransportationOffices.OriginDutyLocation,
			},
		}, nil)

		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2)})

		suite.NoError(err)
		suite.Equal(0, len(*expectedPaymentRequests))
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListFailure() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()

	suite.Run("Error when office user ID does not exist", func() {
		nonexistentOfficeUserID := uuid.Must(uuid.NewV4())
		_, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextForTest(), nonexistentOfficeUserID,
			&services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(2)})

		suite.Error(err)
		suite.Contains(err.Error(), "error fetching transportationOffice for officeUserID")
		suite.Contains(err.Error(), nonexistentOfficeUserID.String())
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestListWithPagination() {
	paymentRequestListFetcher := NewPaymentRequestListFetcher()
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Roles:           officeUser.User.Roles,
		OfficeUserID:    officeUser.ID,
		IDToken:         "fake_token",
		AccessToken:     "fakeAccessToken",
	}

	expectedMove1 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
	expectedMove2 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

	factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				Status: models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    expectedMove1,
			LinkOnly: true,
		},
	}, nil)
	factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				Status: models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    expectedMove2,
			LinkOnly: true,
		},
	}, nil)

	expectedPaymentRequests, count, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &services.FetchPaymentRequestListParams{Page: models.Int64Pointer(1), PerPage: models.Int64Pointer(1)})

	suite.NoError(err)
	suite.Equal(1, len(*expectedPaymentRequests))
	suite.Equal(2, count)

}

func (suite *PaymentRequestServiceSuite) TestListPaymentRequestWithSortOrder() {

	var expectedNameOrder []string
	var expectedDodIDOrder []string
	var expectedEmplidOrder []string
	var expectedStatusOrder []string
	var expectedCreatedAtOrder []time.Time
	var expectedLocatorOrder []string
	var expectedBranchOrder []string
	var expectedOriginDutyLocation []string
	var officeUser models.OfficeUser
	var session auth.Session

	branchNavy := models.AffiliationNAVY
	paymentRequestListFetcher := NewPaymentRequestListFetcher()

	//
	suite.PreloadData(func() {
		officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})

		session = auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		originDutyLocation1 := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "Applewood, CA 99999",
				},
			},
		}, nil)

		originDutyLocation2 := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "Scott AFB",
				},
			},
		}, nil)

		expectedMove1 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					Edipi:       models.StringPointer("EZFG"),
					LastName:    models.StringPointer("Spacemen"),
					FirstName:   models.StringPointer("Lena"),
					Affiliation: &branchNavy,
					Emplid:      models.StringPointer(""),
				},
			},
			{
				Model: models.Move{
					Locator: "AAAA",
				},
			},
			{
				Model:    originDutyLocation1,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		expectedMove2 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Leo"),
					LastName:  models.StringPointer("Spacemen"),
					Edipi:     models.StringPointer("AZFG"),
					Emplid:    models.StringPointer("1111111"),
				},
			},
			{
				Model: models.Move{
					Locator: "ZZZZ",
				},
			},
			{
				Model:    originDutyLocation2,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		// Fake this as a day and a half in the past so floating point age values can be tested
		prevCreatedAt := time.Now().Add(time.Duration(time.Hour * -36))
		paymentRequest1 := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:    models.PaymentRequestStatusPending,
					CreatedAt: prevCreatedAt,
				},
			},
			{
				Model:    expectedMove1,
				LinkOnly: true,
			},
		}, nil)
		paymentRequest2 := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status: models.PaymentRequestStatusReviewed,
				},
			},
			{
				Model:    expectedMove2,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    paymentRequest2.MoveTaskOrder,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		expectedNameOrder = append(expectedNameOrder, *paymentRequest1.MoveTaskOrder.Orders.ServiceMember.FirstName, *paymentRequest2.MoveTaskOrder.Orders.ServiceMember.FirstName)
		expectedDodIDOrder = append(expectedDodIDOrder, *paymentRequest1.MoveTaskOrder.Orders.ServiceMember.Edipi, *paymentRequest2.MoveTaskOrder.Orders.ServiceMember.Edipi)
		expectedEmplidOrder = append(expectedEmplidOrder, *paymentRequest1.MoveTaskOrder.Orders.ServiceMember.Emplid, *paymentRequest2.MoveTaskOrder.Orders.ServiceMember.Emplid)
		expectedStatusOrder = append(expectedStatusOrder, string(paymentRequest1.Status), string(paymentRequest2.Status))
		expectedCreatedAtOrder = append(expectedCreatedAtOrder, paymentRequest1.CreatedAt, paymentRequest2.CreatedAt)
		expectedLocatorOrder = append(expectedLocatorOrder, paymentRequest1.MoveTaskOrder.Locator, paymentRequest2.MoveTaskOrder.Locator)
		expectedBranchOrder = append(expectedBranchOrder, string(*paymentRequest1.MoveTaskOrder.Orders.ServiceMember.Affiliation), string(*paymentRequest2.MoveTaskOrder.Orders.ServiceMember.Affiliation))
		expectedOriginDutyLocation = append(expectedOriginDutyLocation, string(paymentRequest1.MoveTaskOrder.Orders.OriginDutyLocation.Name), string(paymentRequest2.MoveTaskOrder.Orders.OriginDutyLocation.Name))
	})

	suite.Run("Sort by service member name ASC", func() {
		sort.Strings(expectedNameOrder)

		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("lastName"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedNameOrder[0], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName)
		suite.Equal(expectedNameOrder[1], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.FirstName)
	})

	suite.Run("Sort by service member name DESC", func() {
		sort.Strings(expectedNameOrder)

		// Sort by service member name
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("lastName"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedNameOrder[0], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.FirstName)
		suite.Equal(expectedNameOrder[1], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.FirstName)
	})

	suite.Run("Sort by dodID ASC", func() {
		sort.Strings(expectedDodIDOrder)

		// Sort by dodID
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("dodID"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedDodIDOrder[0], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Edipi)
		suite.Equal(expectedDodIDOrder[1], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Edipi)
	})

	suite.Run("Sort by dodID DESC", func() {
		sort.Strings(expectedDodIDOrder)

		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("dodID"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedDodIDOrder[0], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Edipi)
		suite.Equal(expectedDodIDOrder[1], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Edipi)
	})

	suite.Run("Sort by emplid ASC", func() {
		sort.Strings(expectedEmplidOrder)

		// Sort by emplid
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("emplid"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedEmplidOrder[0], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Emplid)
		suite.Equal(expectedEmplidOrder[1], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Emplid)
	})

	suite.Run("Sort by emplid DESC", func() {
		sort.Strings(expectedEmplidOrder)

		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("emplid"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedEmplidOrder[0], *paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Emplid)
		suite.Equal(expectedEmplidOrder[1], *paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Emplid)
	})

	suite.Run("Sort by status ASC", func() {
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("status"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedStatusOrder[0], string(paymentRequests[0].Status))
		suite.Equal(expectedStatusOrder[1], string(paymentRequests[1].Status))
	})

	suite.Run("Sort by status DESC", func() {
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("status"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedStatusOrder[0], string(paymentRequests[1].Status))
		suite.Equal(expectedStatusOrder[1], string(paymentRequests[0].Status))
	})

	suite.Run("Sort by age ASC", func() {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("age"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.Run("Sort by age DESC", func() {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("age"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.Run("Sort by submittedAt ASC", func() {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("submittedAt"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.Run("Sort by submittedAt DESC", func() {
		sort.Slice(expectedCreatedAtOrder, func(i, j int) bool { return expectedCreatedAtOrder[i].Before(expectedCreatedAtOrder[j]) })
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("submittedAt"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedCreatedAtOrder[0].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[1].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
		suite.Equal(expectedCreatedAtOrder[1].Format("2006-01-02T15:04:05.000Z07:00"), paymentRequests[0].CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"))
	})

	suite.Run("Sort by locator ASC", func() {
		sort.Strings(expectedLocatorOrder)
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedLocatorOrder[0], strings.TrimSpace(paymentRequests[0].MoveTaskOrder.Locator))
		suite.Equal(expectedLocatorOrder[1], strings.TrimSpace(paymentRequests[1].MoveTaskOrder.Locator))
	})

	suite.Run("Sort by locator DESC", func() {
		sort.Strings(expectedLocatorOrder)

		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("locator"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedLocatorOrder[0], strings.TrimSpace(paymentRequests[1].MoveTaskOrder.Locator))
		suite.Equal(expectedLocatorOrder[1], strings.TrimSpace(paymentRequests[0].MoveTaskOrder.Locator))
	})

	suite.Run("Sort by branch ASC", func() {
		sort.Strings(expectedBranchOrder)
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("branch"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedBranchOrder[0], string(*paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation))
		suite.Equal(expectedBranchOrder[1], string(*paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Affiliation))
	})

	suite.Run("Sort by branch DESC", func() {
		sort.Strings(expectedBranchOrder)
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("branch"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedBranchOrder[0], string(*paymentRequests[1].MoveTaskOrder.Orders.ServiceMember.Affiliation))
		suite.Equal(expectedBranchOrder[1], string(*paymentRequests[0].MoveTaskOrder.Orders.ServiceMember.Affiliation))
	})

	suite.Run("Sort by originDutyLocation ASC", func() {
		sort.Strings(expectedOriginDutyLocation)
		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("originDutyLocation"), Order: models.StringPointer("asc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		suite.NoError(err)

		paymentRequests := *expectedPaymentRequests
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedOriginDutyLocation[0], string(paymentRequests[0].MoveTaskOrder.Orders.OriginDutyLocation.Name))
		suite.Equal(expectedOriginDutyLocation[1], string(paymentRequests[1].MoveTaskOrder.Orders.OriginDutyLocation.Name))
	})

	suite.Run("Sort by originDutyLocation DESC", func() {
		sort.Strings(expectedOriginDutyLocation)

		params := services.FetchPaymentRequestListParams{Sort: models.StringPointer("originDutyLocation"), Order: models.StringPointer("desc")}
		expectedPaymentRequests, _, err := paymentRequestListFetcher.FetchPaymentRequestList(suite.AppContextWithSessionForTest(&session), officeUser.ID, &params)
		paymentRequests := *expectedPaymentRequests

		suite.NoError(err)
		suite.Equal(2, len(paymentRequests))
		suite.Equal(expectedOriginDutyLocation[0], string(paymentRequests[1].MoveTaskOrder.Orders.OriginDutyLocation.Name))
		suite.Equal(expectedOriginDutyLocation[1], string(paymentRequests[0].MoveTaskOrder.Orders.OriginDutyLocation.Name))
	})
}
