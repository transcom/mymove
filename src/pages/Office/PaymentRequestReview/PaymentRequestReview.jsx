import React, { useState } from 'react';
import { withRouter } from 'react-router-dom';
import { queryCache, useMutation } from 'react-query';

import styles from './PaymentRequestReview.module.scss';

import { formatPaymentRequestReviewAddressString, getShipmentModificationType } from 'utils/shipmentDisplay';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { HistoryShape, MatchShape } from 'types/router';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import { LOA_TYPE, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { patchPaymentRequest, patchPaymentServiceItemStatus } from 'services/ghcApi';
import { usePaymentRequestQueries } from 'hooks/queries';
import { PAYMENT_REQUESTS } from 'constants/queryKeys';
import { OrderShape } from 'types';

export const PaymentRequestReview = ({ history, match, order }) => {
  const [completeReviewError, setCompleteReviewError] = useState(undefined);
  const { paymentRequestId, moveCode } = match.params;
  const { tac, sac, ntsTac, ntsSac } = order;
  const {
    paymentRequest,
    paymentRequests,
    paymentServiceItems,
    mtoShipments,
    shipmentsPaymentSITBalance,
    isLoading,
    isError,
  } = usePaymentRequestQueries(paymentRequestId);

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
      newPaymentServiceItem.mtoServiceItemCode = oldPaymentServiceItem.mtoServiceItemCode;
      newPaymentServiceItem.paymentServiceItemParams = oldPaymentServiceItem.paymentServiceItemParams;

      queryCache.setQueryData([PAYMENT_REQUESTS, paymentRequestId], {
        paymentRequests,
        paymentServiceItems: {
          ...paymentServiceItems,
          [`${variables.paymentServiceItemID}`]: newPaymentServiceItem,
        },
      });
    },
    throwOnError: true,
  });

  const serviceItemCards = React.useMemo(() => {
    // avoids needing to search through the same arary to match up shipment ids for payment service items
    // eslint has a rule about function param reassignment so not using array.reduce here
    const normalizedShipments = {};
    mtoShipments?.forEach((shipment) => {
      normalizedShipments[shipment.id] = shipment;
    });

    return Object.values(paymentServiceItems || {}).map((item) => {
      const selectedShipment = normalizedShipments[item.mtoShipmentID];
      const shipmentSITBalance = shipmentsPaymentSITBalance
        ? shipmentsPaymentSITBalance[item.mtoShipmentID]
        : undefined;
      return {
        id: item.id,
        mtoShipmentID: item.mtoShipmentID,
        mtoShipmentType: item.mtoShipmentType,
        mtoShipmentDepartureDate: selectedShipment?.actualPickupDate,
        mtoShipmentPickupAddress: selectedShipment
          ? formatPaymentRequestReviewAddressString(selectedShipment.pickupAddress)
          : undefined,
        mtoShipmentDestinationAddress: selectedShipment
          ? formatPaymentRequestReviewAddressString(selectedShipment.destinationAddress)
          : undefined,
        mtoShipmentTacType: item.mtoShipmentType === LOA_TYPE.HHG ? LOA_TYPE.HHG : selectedShipment?.tacType,
        mtoShipmentSacType: item.mtoShipmentType === LOA_TYPE.HHG ? LOA_TYPE.HHG : selectedShipment?.sacType,
        mtoShipmentModificationType: selectedShipment ? getShipmentModificationType(selectedShipment) : undefined,
        mtoServiceItemCode: item.mtoServiceItemCode,
        mtoServiceItemName: item.mtoServiceItemName,
        mtoServiceItems: selectedShipment?.mtoServiceItems,
        amount: item.priceCents ? item.priceCents / 100 : 0,
        createdAt: item.createdAt,
        status: item.status,
        rejectionReason: item.rejectionReason,
        paymentServiceItemParams: item.paymentServiceItemParams,
        shipmentSITBalance,
      };
    });
  }, [paymentServiceItems, mtoShipments, shipmentsPaymentSITBalance]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const uploads = paymentRequest.proofOfServiceDocs
    ? paymentRequest.proofOfServiceDocs.flatMap((docs) => docs.uploads.flatMap((primeUploads) => primeUploads))
    : [];

  const handleUpdatePaymentServiceItemStatus = (paymentServiceItemID, values) => {
    return mutatePaymentServiceItemStatus({
      moveTaskOrderID: paymentRequest.moveTaskOrderID,
      paymentServiceItemID,
      status: values.status,
      ifMatchEtag: paymentServiceItems[paymentServiceItemID].eTag,
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
          TACs={{ HHG: tac, NTS: ntsTac }}
          SACs={{ HHG: sac, NTS: ntsSac }}
        />
      </div>
    </div>
  );
};

PaymentRequestReview.propTypes = {
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
  order: OrderShape.isRequired,
};

export default withRouter(PaymentRequestReview);
