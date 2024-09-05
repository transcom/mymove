import React, { useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import styles from './PaymentRequestReview.module.scss';

import { sortServiceItemsByGroup } from 'utils/serviceItems';
import { formatCityStateAndPostalCode, getShipmentModificationType } from 'utils/shipmentDisplay';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import DocumentViewer from 'components/DocumentViewer/DocumentViewer';
import ReviewServiceItems from 'components/Office/ReviewServiceItems/ReviewServiceItems';
import { LOA_TYPE, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { bulkDownloadPaymentRequest, patchPaymentRequest, patchPaymentServiceItemStatus } from 'services/ghcApi';
import { usePaymentRequestQueries } from 'hooks/queries';
import { PAYMENT_REQUESTS } from 'constants/queryKeys';
import { OrderShape } from 'types';
import AsyncPacketDownloadLink from 'shared/AsyncPacketDownloadLink/AsyncPacketDownloadLink';

export const PaymentRequestReview = ({ order }) => {
  const navigate = useNavigate();
  const [completeReviewError, setCompleteReviewError] = useState(undefined);
  const { paymentRequestId, moveCode } = useParams();
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
  const queryClient = useQueryClient();

  const { mutate: mutatePaymentRequest } = useMutation(patchPaymentRequest, {
    onSuccess: (data, variables) => {
      const { paymentRequestID } = variables;
      queryClient.setQueryData([PAYMENT_REQUESTS, paymentRequestID], {
        paymentRequests: data.paymentRequests,
        paymentServiceItems,
      });
      // TODO - show flash message?
      navigate(`/moves/${moveCode}/payment-requests`);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      setCompleteReviewError(errorMsg);
    },
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

      const normalizedServiceItems = {};
      selectedShipment?.mtoServiceItems?.forEach((serviceItem) => {
        normalizedServiceItems[serviceItem.id] = serviceItem;
      });
      const selectedServiceItem = normalizedServiceItems[item.mtoServiceItemID];

      const shipmentSITBalance = shipmentsPaymentSITBalance
        ? shipmentsPaymentSITBalance[item.mtoShipmentID]
        : undefined;
      return {
        id: item.id,
        mtoShipmentID: item.mtoShipmentID,
        mtoShipmentType: item.mtoShipmentType,
        mtoShipmentDepartureDate: selectedShipment?.actualPickupDate,
        mtoShipmentPickupAddress: selectedShipment
          ? formatCityStateAndPostalCode(selectedShipment.pickupAddress)
          : undefined,
        mtoShipmentDestinationAddress: selectedShipment
          ? formatCityStateAndPostalCode(selectedShipment.destinationAddress)
          : undefined,
        mtoShipmentTacType: item.mtoShipmentType === LOA_TYPE.HHG ? LOA_TYPE.HHG : selectedShipment?.tacType,
        mtoShipmentSacType: item.mtoShipmentType === LOA_TYPE.HHG ? LOA_TYPE.HHG : selectedShipment?.sacType,
        mtoShipmentModificationType: selectedShipment ? getShipmentModificationType(selectedShipment) : undefined,
        mtoServiceItemCode: item.mtoServiceItemCode,
        mtoServiceItemName: item.mtoServiceItemName,
        mtoServiceItems: selectedShipment?.mtoServiceItems,
        mtoServiceItemStandaloneCrate: selectedServiceItem?.standaloneCrate,
        amount: item.priceCents ? item.priceCents / 100 : 0,
        createdAt: item.createdAt,
        status: item.status,
        rejectionReason: item.rejectionReason,
        paymentServiceItemParams: item.paymentServiceItemParams,
        shipmentSITBalance,
      };
    });
  }, [paymentServiceItems, mtoShipments, shipmentsPaymentSITBalance]);

  const requestReviewed = paymentRequest?.status !== PAYMENT_REQUEST_STATUS.PENDING;

  const sortedCards = sortServiceItemsByGroup(serviceItemCards);

  const totalCards = sortedCards.length;

  const [curCardIndex, setCardIndex] = useState(requestReviewed ? totalCards : 0);

  const [shouldAdvanceOnSubmit, setShouldAdvanceOnSubmit] = useState(false);

  const handlePrevious = () => {
    setCardIndex(curCardIndex - 1);
  };

  const { mutate: mutatePaymentServiceItemStatus } = useMutation(patchPaymentServiceItemStatus, {
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

      queryClient.setQueryData([PAYMENT_REQUESTS, paymentRequestId], {
        paymentRequests,
        paymentServiceItems: {
          ...paymentServiceItems,
          [`${variables.paymentServiceItemID}`]: newPaymentServiceItem,
        },
      });
      if (shouldAdvanceOnSubmit) {
        setCardIndex(curCardIndex + 1);
      }
    },
    throwOnError: true,
  });

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
    navigate(`/moves/${moveCode}/payment-requests`);
  };

  const paymentPacketDownload = (
    <div>
      <dd data-testid="bulkPacketDownload">
        <p className={styles.downloadLink}>
          <AsyncPacketDownloadLink
            id={paymentRequestId}
            label="Download All Files (PDF)"
            asyncRetrieval={bulkDownloadPaymentRequest}
          />
        </p>
      </dd>
    </div>
  );

  return (
    <div data-testid="PaymentRequestReview" className={styles.PaymentRequestReview}>
      <div className={styles.embed}>
        {uploads.length > 0 ? (
          <>
            {paymentPacketDownload}
            <DocumentViewer files={uploads} />
          </>
        ) : (
          <h2>No documents provided</h2>
        )}
      </div>
      <div className={styles.sidebar}>
        <ReviewServiceItems
          handleClose={handleClose}
          paymentRequest={paymentRequest}
          serviceItemCards={sortedCards}
          patchPaymentServiceItem={handleUpdatePaymentServiceItemStatus}
          onCompleteReview={handleCompleteReview}
          completeReviewError={completeReviewError}
          TACs={{ HHG: tac, NTS: ntsTac }}
          SACs={{ HHG: sac, NTS: ntsSac }}
          curCardIndex={curCardIndex}
          setCardIndex={setCardIndex}
          handlePrevious={handlePrevious}
          requestReviewed={requestReviewed}
          setShouldAdvanceOnSubmit={setShouldAdvanceOnSubmit}
        />
      </div>
    </div>
  );
};

PaymentRequestReview.propTypes = {
  order: OrderShape.isRequired,
};

export default PaymentRequestReview;
