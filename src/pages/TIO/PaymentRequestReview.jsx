import React, { useState } from 'react';
import { withRouter } from 'react-router-dom';
import { useQuery, useMutation, queryCache } from 'react-query';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MatchShape, HistoryShape } from 'types/router';
import samplePDF from 'components/DocumentViewer/sample.pdf';
import styles from 'pages/TIO/PaymentRequestReview.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import { PAYMENT_REQUEST_STATUS } from 'shared/constants';
import {
  getPaymentRequest,
  getMTOShipments,
  getMTOServiceItems,
  patchPaymentRequest,
  patchPaymentServiceItemStatus,
} from 'services/ghcApi';
import { mapObjectToArray, getQueriesStatus } from 'utils/api';

const usePaymentRequestQueries = (paymentRequestId) => {
  // get payment request by ID
  const { data: { paymentRequests, paymentServiceItems } = {}, ...paymentRequestQuery } = useQuery(
    ['paymentRequest', paymentRequestId],
    getPaymentRequest,
  );

  const paymentRequest = paymentRequests && paymentRequests[`${paymentRequestId}`];
  const mtoID = paymentRequest?.moveTaskOrderID;

  // get MTO shipments
  const { data: { mtoShipments = [] } = {}, ...mtoShipmentQuery } = useQuery(['mtoShipment', mtoID], getMTOShipments, {
    enabled: !!mtoID,
  });

  // get MTO service items
  const { data: { mtoServiceItems = [] } = {}, ...mtoServiceItemQuery } = useQuery(
    ['mtoServiceItem', mtoID],
    getMTOServiceItems,
    {
      enabled: !!mtoID,
    },
  );

  const { isLoading, isError, isSuccess } = getQueriesStatus([
    paymentRequestQuery,
    mtoShipmentQuery,
    mtoServiceItemQuery,
  ]);

  return {
    paymentRequest,
    paymentRequests,
    paymentServiceItems,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
    isSuccess,
  };
};

const PaymentRequestReview = ({ history, match }) => {
  const [completeReviewError, setCompleteReviewError] = useState(undefined);
  const { paymentRequestId, moveOrderId } = match.params;

  const {
    paymentRequest,
    paymentRequests,
    paymentServiceItems,
    mtoShipments,
    mtoServiceItems,
    isLoading,
    isError,
  } = usePaymentRequestQueries(paymentRequestId);

  const [mutatePaymentRequest] = useMutation(patchPaymentRequest, {
    onSuccess: () => {
      // TODO - update cache?
      // TODO - show flash message?
      history.push(`/`); // Go home
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      setCompleteReviewError(errorMsg);
    },
  });

  const [mutatePaymentServiceItemStatus] = useMutation(patchPaymentServiceItemStatus, {
    onSuccess: (data, variables) => {
      const newPaymentServiceItem = data.paymentServiceItems[variables.paymentServiceItemID];
      queryCache.setQueryData(['paymentRequest', paymentRequestId], {
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

  const paymentServiceItemsArr = mapObjectToArray(paymentServiceItems);
  const mtoServiceItemsArr = mapObjectToArray(mtoServiceItems);
  const mtoShipmentsArr = mapObjectToArray(mtoShipments);

  const handleUpdatePaymentServiceItemStatus = (paymentServiceItemID, values) => {
    const paymentServiceItemForRequest = paymentServiceItemsArr.find((s) => s.id === paymentServiceItemID);

    mutatePaymentServiceItemStatus({
      moveTaskOrderID: mtoServiceItemsArr[0].moveTaskOrderID,
      paymentServiceItemID,
      status: values.status,
      ifMatchEtag: paymentServiceItemForRequest.eTag,
      rejectionReason: values.rejectionReason,
    });
  };

  const handleCompleteReview = () => {
    // first reset error if there was one
    if (completeReviewError) setCompleteReviewError(undefined);

    const newPaymentRequest = {
      paymentRequestID: paymentRequest.id,
      ifMatchETag: paymentRequest.eTag,
      status: PAYMENT_REQUEST_STATUS.REVIEWED,
    };

    mutatePaymentRequest(newPaymentRequest);
  };

  const handleClose = () => {
    history.push(`/moves/${moveOrderId}/payment-requests/`);
  };

  const testFiles = [
    {
      filename: 'Test File.pdf',
      fileType: 'pdf',
      filePath: samplePDF,
    },
  ];

  const serviceItemCards = paymentServiceItemsArr.map((item) => {
    const mtoServiceItem = mtoServiceItemsArr.find((s) => s.id === item.mtoServiceItemID);
    const itemShipment = mtoServiceItem && mtoShipmentsArr.find((s) => s.id === mtoServiceItem.mtoShipmentID);

    return {
      id: item.id,
      shipmentId: mtoServiceItem?.mtoShipmentID,
      shipmentType: itemShipment?.shipmentType,
      serviceItemName: mtoServiceItem?.reServiceName,
      amount: item.priceCents ? item.priceCents / 100 : 0,
      createdAt: item.createdAt,
      status: item.status,
      rejectionReason: item.rejectionReason,
    };
  });

  return (
    <div data-testid="PaymentRequestReview" className={styles.PaymentRequestReview}>
      <div className={styles.embed}>
        <DocumentViewer files={testFiles} />
      </div>
      <div className={styles.sidebar}>
        <ReviewServiceItems
          handleClose={handleClose}
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
