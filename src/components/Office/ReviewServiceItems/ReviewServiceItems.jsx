import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { sortServiceItemsByGroup } from '../../../utils/serviceItems';

import styles from './ReviewServiceItems.module.scss';
import ServiceItemCard from './ServiceItemCard';
import ReviewDetailsCard from './ReviewDetailsCard';
import AuthorizePayment from './AuthorizePayment';
import NeedsReview from './NeedsReview';
import RejectRequest from './RejectRequest';
import PaymentReviewed from './PaymentReviewed';

import Alert from 'shared/Alert';
import { ServiceItemCardsShape } from 'types/serviceItemCard';
import { PAYMENT_SERVICE_ITEM_STATUS, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { toDollarString } from 'shared/formatters';
import { PaymentRequestShape } from 'types/index';

const ReviewServiceItems = ({
  header,
  paymentRequest,
  serviceItemCards,
  handleClose,
  disableScrollIntoView,
  patchPaymentServiceItem,
  onCompleteReview,
  completeReviewError,
}) => {
  const requestReviewed = paymentRequest?.status !== PAYMENT_REQUEST_STATUS.PENDING;

  const sortedCards = sortServiceItemsByGroup(serviceItemCards);

  const totalCards = sortedCards.length;

  const [curCardIndex, setCardIndex] = useState(requestReviewed ? totalCards : 0);

  const { APPROVED, DENIED, REQUESTED } = PAYMENT_SERVICE_ITEM_STATUS;

  const handleServiceItemNavBtnClick = (index) => {
    setCardIndex(index);
  };

  const handleAuthorizePaymentClick = (allServiceItemsRejected = false) => {
    onCompleteReview(allServiceItemsRejected);
  };

  // calculating the sums
  const approvedSum = sortedCards.filter((s) => s.status === APPROVED).reduce((sum, cur) => sum + cur.amount, 0);
  const rejectedSum = sortedCards.filter((s) => s.status === DENIED).reduce((sum, cur) => sum + cur.amount, 0);
  const requestedSum = sortedCards.reduce((sum, cur) => sum + cur.amount, 0);

  let itemsNeedsReviewLength;
  let showNeedsReview;
  let allServiceItemsRejected;
  let firstItemNeedsReviewIndex;
  let firstBasicIndex = null;
  let lastBasicIndex = null;

  if (!requestReviewed) {
    itemsNeedsReviewLength = sortedCards.filter((s) => s.status === REQUESTED)?.length;
    showNeedsReview = sortedCards.some((s) => s.status === REQUESTED);
    allServiceItemsRejected = sortedCards.every((s) => s.status === DENIED);
    firstItemNeedsReviewIndex = showNeedsReview && sortedCards.findIndex((s) => s.status === REQUESTED);

    sortedCards.forEach((serviceItem, index) => {
      // here we want to set the first and last index
      // of basic service items to know the bounds
      if (!serviceItem.mtoShipmentType) {
        // no shipemntId, then it is a basic service items
        if (firstBasicIndex === null) {
          // if not set yet, set it the first time we see a basic
          // service item
          firstBasicIndex = index;
        }
        // keep setting the last basic index until the last one
        lastBasicIndex = index;
      }
    });
  }

  const displayCompleteReview = curCardIndex === totalCards;

  const currentCard = !displayCompleteReview && sortedCards[parseInt(curCardIndex, 10)];

  const isBasicServiceItem =
    firstBasicIndex !== null && curCardIndex >= firstBasicIndex && curCardIndex <= lastBasicIndex;

  // Determines which ReviewDetailsCard will be shown
  let renderCompleteAction = (
    <AuthorizePayment amount={approvedSum} onClick={() => handleAuthorizePaymentClick(requestReviewed)} />
  );
  if (requestReviewed) {
    renderCompleteAction = (
      <PaymentReviewed authorizedAmount={approvedSum} dateAuthorized={paymentRequest?.reviewedAt} />
    );
  } else if (showNeedsReview) {
    renderCompleteAction = (
      <NeedsReview numberOfItems={itemsNeedsReviewLength} onClick={() => setCardIndex(firstItemNeedsReviewIndex)} />
    );
  } else if (allServiceItemsRejected) {
    renderCompleteAction = <RejectRequest onClick={() => handleAuthorizePaymentClick(allServiceItemsRejected)} />;
  }

  // Similar to component lifecycle methods
  useEffect(() => {
    if (currentCard) {
      const { id } = sortedCards[parseInt(curCardIndex, 10)];
      const element = document.querySelector(`#card-${id}`);
      // scroll into element view
      if (element && !disableScrollIntoView) {
        element.scrollIntoView();
      }
    }
  });

  if (displayCompleteReview)
    return (
      <div data-testid="ReviewServiceItems" className={styles.ReviewServiceItems}>
        <div className={styles.top}>
          <Button data-testid="closeSidebar" type="button" onClick={handleClose} unstyled>
            <FontAwesomeIcon icon="times" title="Close Service Item review" aria-label="Close Service Item review" />
          </Button>
          <h2 className={styles.header}>Complete request</h2>
        </div>
        <div className={styles.body}>
          {requestReviewed && (
            <Alert heading={null} type="success">
              The payment request was successfully submitted.
            </Alert>
          )}
          <ReviewDetailsCard
            completeReviewError={completeReviewError}
            acceptedAmount={approvedSum}
            rejectedAmount={rejectedSum}
            requestedAmount={requestedSum}
          >
            {renderCompleteAction}
          </ReviewDetailsCard>
        </div>
        <div className={styles.bottom}>
          <Button
            data-testid="prevServiceItem"
            type="button"
            onClick={() => handleServiceItemNavBtnClick(curCardIndex - 1)}
            secondary
            disabled={curCardIndex === 0}
          >
            Back
          </Button>
        </div>
      </div>
    );

  return (
    <div data-testid="ReviewServiceItems" className={styles.ReviewServiceItems}>
      <div className={styles.top}>
        <Button data-testid="closeSidebar" type="button" onClick={handleClose} unstyled>
          <FontAwesomeIcon icon="times" aria-label="Close Service Item review" />
        </Button>
        <div data-testid="itemCount" className={styles.eyebrowTitle}>
          {curCardIndex + 1} OF {totalCards} ITEMS
        </div>
        <h2 className={styles.header}>{header}</h2>
      </div>
      <div className={styles.body}>
        {currentCard && // render multiple basic service item cards
          // otherwise, render only one card for shipment
          (isBasicServiceItem ? (
            sortedCards.slice(firstBasicIndex, lastBasicIndex + 1).map((curCard) => (
              <ServiceItemCard
                key={`serviceItemCard_${curCard.id}`}
                patchPaymentServiceItem={patchPaymentServiceItem}
                //  react/jsx-props-no-spreading
                {...curCard}
                requestComplete={requestReviewed}
              />
            ))
          ) : (
            <ServiceItemCard
              key={`serviceItemCard_${currentCard.id}`}
              patchPaymentServiceItem={patchPaymentServiceItem}
              //  react/jsx-props-no-spreading
              {...currentCard}
              requestComplete={requestReviewed}
            />
          ))}
      </div>
      <div className={styles.bottom}>
        <Button
          data-testid="prevServiceItem"
          aria-label="Previous Service Item"
          type="button"
          onClick={() => handleServiceItemNavBtnClick(curCardIndex - 1)}
          secondary
          disabled={curCardIndex === 0}
        >
          Previous
        </Button>
        <Button
          data-testid="nextServiceItem"
          aria-label="Next Service Item"
          type="button"
          onClick={() => handleServiceItemNavBtnClick(curCardIndex + 1)}
          disabled={curCardIndex === totalCards}
        >
          Next
        </Button>
        <div className={styles.totalApproved}>
          <div className={styles.totalLabel}>Total approved</div>
          <div className={styles.totalAmount} data-testid="approvedAmount">
            {toDollarString(approvedSum)}
          </div>
        </div>
      </div>
    </div>
  );
};

ReviewServiceItems.propTypes = {
  header: PropTypes.string,
  paymentRequest: PaymentRequestShape,
  serviceItemCards: ServiceItemCardsShape,
  handleClose: PropTypes.func.isRequired,
  patchPaymentServiceItem: PropTypes.func.isRequired,
  disableScrollIntoView: PropTypes.bool,
  onCompleteReview: PropTypes.func.isRequired,
  completeReviewError: PropTypes.shape({
    detail: PropTypes.string,
    title: PropTypes.string,
  }),
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
  paymentRequest: undefined,
  serviceItemCards: [],
  disableScrollIntoView: false,
  completeReviewError: undefined,
};

export default ReviewServiceItems;
