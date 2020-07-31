import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { useQuery, useMutation, queryCache } from 'react-query';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { MatchShape, HistoryShape } from 'types/router';
import samplePDF from 'components/DocumentViewer/sample.pdf';
import styles from 'pages/TIO/PaymentRequestReview.module.scss';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import { updatePaymentRequest as updatePaymentRequestAction } from 'shared/Entities/modules/paymentRequests';
import { PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { getPaymentRequest, getMTOShipments, getMTOServiceItems, patchPaymentServiceItemStatus } from 'services/ghcApi';

const PaymentRequestReview = ({ updatePaymentRequest, history, match }) => {
  const [completeReviewError, setCompleteReviewError] = useState(undefined);
  const { paymentRequestId, moveOrderId } = match.params;

  const { data: { paymentRequests, paymentServiceItems } = {}, ...paymentRequestQuery } = useQuery(
    ['paymentRequest', paymentRequestId],
    getPaymentRequest,
    {
      retry: false,
    },
  );

  const paymentRequest = paymentRequests && paymentRequests[`${paymentRequestId}`];
  const mtoID = paymentRequest?.moveTaskOrderID;

  const { data: { mtoShipments = [] } = {}, ...mtoShipmentQuery } = useQuery(['mtoShipment', mtoID], getMTOShipments, {
    enabled: !!mtoID,
  });

  const { data: { mtoServiceItems = [] } = {}, ...mtoServiceItemQuery } = useQuery(
    ['mtoServiceItem', mtoID],
    getMTOServiceItems,
    {
      enabled: !!mtoID,
    },
  );

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

  const isLoading = paymentRequestQuery.isLoading || mtoShipmentQuery.isLoading || mtoServiceItemQuery.isLoading;
  const isError = paymentRequestQuery.isError || mtoShipmentQuery.isError || mtoServiceItemQuery.isError;
  const error = paymentRequestQuery.error || mtoShipmentQuery.error || mtoServiceItemQuery.error;

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong error={error} />;

  // TODO - util fn
  // TODO - normalize changes?
  // eslint-disable-next-line security/detect-object-injection
  const paymentServiceItemsArr = Object.keys(paymentServiceItems).map((i) => paymentServiceItems[i]);
  // eslint-disable-next-line security/detect-object-injection
  const mtoServiceItemsArr = Object.keys(mtoServiceItems).map((i) => mtoServiceItems[i]);
  // eslint-disable-next-line security/detect-object-injection
  const mtoShipmentsArr = Object.keys(mtoShipments).map((i) => mtoShipments[i]);

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
    // TODO - rewrite with mutation
    // first reset error if there was one
    if (completeReviewError) setCompleteReviewError(undefined);

    const newPaymentRequest = {
      paymentRequestID: paymentRequest.id,
      ifMatchETag: paymentRequest.eTag,
      status: PAYMENT_REQUEST_STATUS.REVIEWED,
    };

    updatePaymentRequest(newPaymentRequest)
      .then(() => {
        // TODO - show flash message?
        history.push(`/`); // Go home
      })
      .catch((e) => {
        const errorMsg = e.response?.response?.body;
        setCompleteReviewError(errorMsg);
      });
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
  updatePaymentRequest: PropTypes.func.isRequired,
};

PaymentRequestReview.defaultProps = {};

const mapStateToProps = () => ({});

const mapDispatchToProps = {
  updatePaymentRequest: updatePaymentRequestAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(PaymentRequestReview));
