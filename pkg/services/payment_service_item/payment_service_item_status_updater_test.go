package paymentserviceitem

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PaymentServiceItemSuite) TestUpdatePaymentServiceItemStatus() {
	suite.Run("Successfully approves a payment service item", func() {
		paymentServiceItem := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		updater := NewPaymentServiceItemStatusUpdater()

		updatedPaymentServiceItem, verrs, err := updater.UpdatePaymentServiceItemStatus(suite.AppContextForTest(),
			paymentServiceItem.ID, models.PaymentServiceItemStatusApproved, nil, eTag)

		suite.NoError(err)
		suite.NoVerrs(verrs)
		suite.Equal(paymentServiceItem.ID, updatedPaymentServiceItem.ID)
		suite.Equal(models.PaymentServiceItemStatusApproved, updatedPaymentServiceItem.Status)
		suite.NotNil(updatedPaymentServiceItem.ApprovedAt)
		suite.Nil(updatedPaymentServiceItem.RejectionReason)
		suite.Nil(updatedPaymentServiceItem.DeniedAt)

	})

	suite.Run("Successfully rejects a payment service item", func() {
		paymentServiceItem := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		updater := NewPaymentServiceItemStatusUpdater()

		updatedPaymentServiceItem, verrs, err := updater.UpdatePaymentServiceItemStatus(suite.AppContextForTest(),
			paymentServiceItem.ID, models.PaymentServiceItemStatusDenied, models.StringPointer("reasons"), eTag)

		suite.NoError(err)
		suite.NoVerrs(verrs)
		suite.Equal(paymentServiceItem.ID, updatedPaymentServiceItem.ID)
		suite.Equal(models.PaymentServiceItemStatusDenied, updatedPaymentServiceItem.Status)
		suite.NotNil(updatedPaymentServiceItem.DeniedAt)
		suite.Equal("reasons", *updatedPaymentServiceItem.RejectionReason)
		suite.Nil(updatedPaymentServiceItem.ApprovedAt)

	})

	suite.Run("Fails if we can't find an existing paymentServiceItem", func() {
		paymentServiceItem := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		updater := NewPaymentServiceItemStatusUpdater()
		wrongUUID, _ := uuid.NewV4()

		_, _, err := updater.UpdatePaymentServiceItemStatus(suite.AppContextForTest(),
			wrongUUID, models.PaymentServiceItemStatusApproved, nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Fails if we have a stale eTag", func() {
		paymentServiceItem := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		// Arbitrary date time that isn't the record updatedAt used here
		badETag := etag.GenerateEtag(testdatagen.DateInsidePerformancePeriod)
		updater := NewPaymentServiceItemStatusUpdater()

		_, _, err := updater.UpdatePaymentServiceItemStatus(suite.AppContextForTest(),
			paymentServiceItem.ID, models.PaymentServiceItemStatusApproved, nil, badETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Fails if we attempt to reject without a rejection reason", func() {
		paymentServiceItem := factory.BuildPaymentServiceItem(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(paymentServiceItem.UpdatedAt)
		updater := NewPaymentServiceItemStatusUpdater()

		_, _, err := updater.UpdatePaymentServiceItemStatus(suite.AppContextForTest(),
			paymentServiceItem.ID, models.PaymentServiceItemStatusDenied, nil, eTag)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

}
