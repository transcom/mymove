package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequest() {
	suite.Run("If a payment request is fetched, it should be returned", func() {

		fetcher := NewPaymentRequestFetcher()

		pr := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), pr.ID)

		suite.NoError(err)
		suite.Equal(pr.ID, paymentRequest.ID)
	})

	suite.Run("returns payment request with proof of service docs", func() {

		fetcher := NewPaymentRequestFetcher()

		primeUpload := factory.BuildPrimeUpload(suite.DB(), nil, nil)
		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), primeUpload.ProofOfServiceDoc.PaymentRequestID)

		suite.NoError(err)
		suite.Equal(primeUpload.ProofOfServiceDoc.PaymentRequest.ID, paymentRequest.ID)

		suite.Len(paymentRequest.ProofOfServiceDocs, 1)
		suite.Equal(primeUpload.ProofOfServiceDoc.ID, paymentRequest.ProofOfServiceDocs[0].ID)

		suite.Len(paymentRequest.ProofOfServiceDocs[0].PrimeUploads, 1)
		suite.Equal(primeUpload.ID, paymentRequest.ProofOfServiceDocs[0].PrimeUploads[0].ID)

		suite.Equal(primeUpload.UploadID, paymentRequest.ProofOfServiceDocs[0].PrimeUploads[0].UploadID)
	})

	suite.Run("returns payment request without soft deleted proof of service docs", func() {

		fetcher := NewPaymentRequestFetcher()
		primeUpload := factory.BuildPrimeUpload(suite.DB(), nil, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), primeUpload.ProofOfServiceDoc.PaymentRequest.ID)

		suite.NoError(err)
		suite.Equal(primeUpload.ProofOfServiceDoc.PaymentRequest.ID, paymentRequest.ID)

		suite.Len(paymentRequest.ProofOfServiceDocs, 1)
		suite.Equal(primeUpload.ProofOfServiceDoc.ID, paymentRequest.ProofOfServiceDocs[0].ID)

		suite.Len(paymentRequest.ProofOfServiceDocs[0].PrimeUploads, 0)
	})

	suite.Run("if there is an error, we get it with zero payment request", func() {
		fetcher := NewPaymentRequestFetcher()

		paymentRequest, err := fetcher.FetchPaymentRequest(suite.AppContextForTest(), uuid.Nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(models.PaymentRequest{}, paymentRequest)
	})
}

func (suite *PaymentRequestServiceSuite) TestFetchPaymentRequestsForBulkAssignment() {
	setupTestData := func() (services.PaymentRequestFetcherBulkAssignment, models.TransportationOffice) {
		paymentRequestFetcher := NewPaymentRequestFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		// this move has a transportation office associated with it that matches
		// the TIO's transportation office and should be found
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		move2 := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    move2,
				LinkOnly: true,
			},
		}, nil)

		return paymentRequestFetcher, transportationOffice
	}

	suite.Run("TIO: Returns payment requests that fulfill the query criteria", func() {
		paymentRequestFetcher, _ := setupTestData()
		paymentRequests, err := paymentRequestFetcher.FetchPaymentRequestsForBulkAssignment(suite.AppContextForTest(), "KKFA")
		suite.FatalNoError(err)
		suite.Equal(2, len(paymentRequests))
	})

	suite.Run("Does not return moves that are already assigned", func() {
		paymentRequestFetcher := NewPaymentRequestFetcherBulkAssignment()
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, []roles.RoleType{roles.RoleTypeTIO})

		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &factory.OfficeUsers.TIOAssignedUser,
			},
		}, nil)
		assignedPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentRequests, err := paymentRequestFetcher.FetchPaymentRequestsForBulkAssignment(suite.AppContextForTest(), "KKFA")
		suite.FatalNoError(err)

		// confirm that the assigned move isn't returned
		for _, paymentRequest := range paymentRequests {
			suite.NotEqual(paymentRequest.ID, assignedPaymentRequest.ID)
		}

		// confirm that the rest of the details are correct
		// move is APPROVALS REQUESTED STATUS
		suite.Equal(assignedPaymentRequest.Status, models.PaymentRequestStatusPending)
		// GBLOC is the same
		suite.Equal(*move.Orders.OriginDutyLocationGBLOC, officeUser.TransportationOffice.Gbloc)
		// Show is true
		suite.Equal(move.Show, models.BoolPointer(true))
		// Orders type isn't WW, BB, or Safety
		suite.Equal(move.Orders.OrdersType, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
	})

	suite.Run("TIO: Does not return payment requests with safety, bluebark, or wounded warrior order types", func() {
		paymentRequestFetcher, transportationOffice := setupTestData()
		moveSafety := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					OrdersType: internalmessages.OrdersTypeSAFETY,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    moveSafety,
				LinkOnly: true,
			},
		}, nil)

		moveBB := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					OrdersType: internalmessages.OrdersTypeBLUEBARK,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    moveBB,
				LinkOnly: true,
			},
		}, nil)

		moveWW := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					OrdersType: internalmessages.OrdersTypeWOUNDEDWARRIOR,
				},
			},
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    moveWW,
				LinkOnly: true,
			},
		}, nil)

		paymentRequests, err := paymentRequestFetcher.FetchPaymentRequestsForBulkAssignment(suite.AppContextForTest(), "KKFA")
		suite.FatalNoError(err)
		suite.Equal(2, len(paymentRequests))
	})

	suite.Run("TIO: Does not return payment requests with Marines if GBLOC not USMC", func() {
		paymentRequestFetcher, transportationOffice := setupTestData()

		marine := models.AffiliationMARINES
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusAPPROVALSREQUESTED,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.ServiceMember{
					Affiliation: &marine,
				},
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentRequests, err := paymentRequestFetcher.FetchPaymentRequestsForBulkAssignment(suite.AppContextForTest(), "KKFA")
		suite.FatalNoError(err)
		suite.Equal(2, len(paymentRequests))
	})

	suite.Run("TIO: Only return payment requests with Marines if GBLOC is USMC", func() {
		paymentRequestFetcher, transportationOffice := setupTestData()

		marine := models.AffiliationMARINES
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusServiceCounselingCompleted,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.ServiceMember{
					Affiliation: &marine,
				},
			},
		}, nil)
		factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					ID:              uuid.Must(uuid.NewV4()),
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentRequests, err := paymentRequestFetcher.FetchPaymentRequestsForBulkAssignment(suite.AppContextForTest(), "USMC")
		suite.FatalNoError(err)
		suite.Equal(1, len(paymentRequests))
	})
}
