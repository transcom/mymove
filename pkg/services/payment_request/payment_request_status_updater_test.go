package paymentrequest

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PaymentRequestServiceSuite) TestUpdatePaymentRequestStatus() {
	builder := query.NewQueryBuilder()
	suite.Run("If we get a payment request pointer with a status it should update and return no ", func() {
		setupTestData := func() models.OfficeUser {
			transportationOffice := factory.BuildTransportationOffice(suite.DB(), []factory.Customization{
				{
					Model: models.TransportationOffice{
						ProvidesCloseout: true,
					},
				},
			}, nil)

			officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
				{
					Model:    transportationOffice,
					LinkOnly: true,
					Type:     &factory.TransportationOffices.CloseoutOffice,
				},
			}, []roles.RoleType{roles.RoleTypeTIO})

			return officeUser
		}
		officeUser := setupTestData()
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &factory.OfficeUsers.TIOAssignedUser,
			},
		}, nil)

		paymentRequest.Status = models.PaymentRequestStatusReviewed
		updater := NewPaymentRequestStatusUpdater(builder)
		updatedPr, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))

		suite.NoError(err)
		suite.Nil(updatedPr.MoveTaskOrder.TIOAssignedID)
		suite.Nil(updatedPr.MoveTaskOrder.TIOAssignedUser)
	})

	suite.Run("If we get a payment request pointer with a status it should update and return no error", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.NoError(err)
	})

	suite.Run("Should return a ConflictError if the payment request has any service items that have not been reviewed", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)

		psiCost := unit.Cents(10000)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &psiCost,
					Status:     models.PaymentServiceItemStatusRequested,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &psiCost,
					Status:     models.PaymentServiceItemStatusApproved,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)

		paymentRequest.Status = models.PaymentRequestStatusReviewed
		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
	})

	suite.Run("Should update and return no error if the payment request has service items that have all been reviewed", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)

		psiCost := unit.Cents(10000)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &psiCost,
					Status:     models.PaymentServiceItemStatusApproved,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &psiCost,
					Status:     models.PaymentServiceItemStatusDenied,
				},
			}, {
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil)

		paymentRequest.Status = models.PaymentRequestStatusReviewed
		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(paymentRequest.UpdatedAt))
		suite.NoError(err)
	})

	suite.Run("Should return a PreconditionFailedError with a stale etag", func() {
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		paymentRequest.Status = models.PaymentRequestStatusReviewed

		updater := NewPaymentRequestStatusUpdater(builder)

		_, err := updater.UpdatePaymentRequestStatus(suite.AppContextForTest(), &paymentRequest, etag.GenerateEtag(time.Now()))
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

}
