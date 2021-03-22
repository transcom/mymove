import React, { useState } from 'react';
import { withRouter } from 'react-router-dom';
import { useMutation, queryCache } from 'react-query';

import styles from './PaymentRequestReview.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MatchShape, HistoryShape } from 'types/router';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import { PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { patchPaymentRequest, patchPaymentServiceItemStatus } from 'services/ghcApi';
import { usePaymentRequestQueries } from 'hooks/queries';
import { PAYMENT_REQUESTS } from 'constants/queryKeys';

export const PaymentRequestReview = ({ history, match }) => {
  const [completeReviewError, setCompleteReviewError] = useState(undefined);
  const { paymentRequestId, moveCode } = match.params;
  const { paymentRequest, paymentRequests, paymentServiceItems, isLoading, isError } = usePaymentRequestQueries(
    paymentRequestId,
  );

  const [mutatePaymentRequest] = useMutation(patchPaymentRequest, {
    onSuccess: (data, variables) => {
      const { paymentRequestID } = variables;
      queryCache.setQueryData([PAYMENT_REQUESTS, paymentRequestID], {
        paymentRequests: data.paymentRequests,
        paymentServiceItems,
      });
      // TODO - show flash message?
      history.push(`/moves/${moveCode}/payment-requests`);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      setCompleteReviewError(errorMsg);
    },
  });

  const [mutatePaymentServiceItemStatus] = useMutation(patchPaymentServiceItemStatus, {
    onSuccess: (data, variables) => {
      const newPaymentServiceItem = data.paymentServiceItems[variables.paymentServiceItemID];
      const oldPaymentServiceItem = paymentServiceItems[variables.paymentServiceItemID];

      // We already have this associated data and it won't change on status update and
      // would be overwritten with null values as they aren't in the payload response
      newPaymentServiceItem.mtoServiceItemName = oldPaymentServiceItem.mtoServiceItemName;
      newPaymentServiceItem.mtoShipmentType = oldPaymentServiceItem.mtoShipmentType;
      newPaymentServiceItem.mtoShipmentID = oldPaymentServiceItem.mtoShipmentID;

      queryCache.setQueryData([PAYMENT_REQUESTS, paymentRequestId], {
        paymentRequests,
        paymentServiceItems: {
          ...paymentServiceItems,
          [`${variables.paymentServiceItemID}`]: newPaymentServiceItem,
        },
      });
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const uploads = paymentRequest.proofOfServiceDocs
    ? paymentRequest.proofOfServiceDocs.flatMap((docs) => docs.uploads.flatMap((primeUploads) => primeUploads))
    : [];
  const paymentServiceItemsArr = Object.values(paymentServiceItems);

  const handleUpdatePaymentServiceItemStatus = (paymentServiceItemID, values) => {
    const paymentServiceItemForRequest = paymentServiceItemsArr.find((s) => s.id === paymentServiceItemID);

    mutatePaymentServiceItemStatus({
      moveTaskOrderID: paymentRequest.moveTaskOrderID,
      paymentServiceItemID,
      status: values.status,
      ifMatchEtag: paymentServiceItemForRequest.eTag,
      rejectionReason: values.rejectionReason,
    });
  };

  const handleCompleteReview = (requestRejected = false) => {
    // first reset error if there was one
    if (completeReviewError) setCompleteReviewError(undefined);

    const updatedPaymentRequest = {
      paymentRequestID: paymentRequest.id,
      ifMatchETag: paymentRequest.eTag,
      status: requestRejected
        ? PAYMENT_REQUEST_STATUS.REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED
        : PAYMENT_REQUEST_STATUS.REVIEWED,
    };

    mutatePaymentRequest(updatedPaymentRequest);
  };

  const handleClose = () => {
    history.push(`/moves/${moveCode}/payment-requests`);
  };

  const serviceItemCards = paymentServiceItemsArr.map((item) => {
    return {
      id: item.id,
      mtoShipmentID: item.mtoShipmentID,
      mtoShipmentType: item.mtoShipmentType,
      mtoServiceItemCode: item.mtoServiceItemCode,
      mtoServiceItemName: item.mtoServiceItemName,
      amount: item.priceCents ? item.priceCents / 100 : 0,
      createdAt: item.createdAt,
      status: item.status,
      rejectionReason: item.rejectionReason,
      paymentServiceItemParams: item.paymentServiceItemParams,
    };
  });

  return (
    <div data-testid="PaymentRequestReview" className={styles.PaymentRequestReview}>
      <div className={styles.embed}>
        <DocumentViewer files={uploads} />
      </div>
      <div className={styles.sidebar}>
        <ReviewServiceItems
          handleClose={handleClose}
          paymentRequest={paymentRequest}
          serviceItemCards={serviceItemCards}
          patchPaymentServiceItem={handleUpdatePaymentServiceItemStatus}
          onCompleteReview={handleCompleteReview}
          completeReviewError={completeReviewError}
        />
      </div>
    </div>
  );
};

PaymentRequestReview.propTypes = {
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
};

export default withRouter(PaymentRequestReview);
