import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { sortServiceItemsByGroup } from '../../../utils/serviceItems';

import styles from './NewReviewServiceItems.module.scss';
import ServiceItemCard from './NewServiceItemCard';
import ReviewDetailsCard from './ReviewDetailsCard';
import AuthorizePayment from './AuthorizePayment';
import NeedsReview from './NeedsReview';
import RejectRequest from './RejectRequest';

import Alert from 'shared/Alert';
import { ServiceItemCardsShape } from 'types/serviceItems';
import { MTOServiceItemShape } from 'types/order';
import { PAYMENT_SERVICE_ITEM_STATUS, PAYMENT_REQUEST_STATUS } from 'shared/constants';
import { toDollarString } from 'utils/formatters';
import { PaymentRequestShape } from 'types/index';
import { AccountingCodesShape } from 'types/accountingCodes';

const { APPROVED, DENIED, REQUESTED } = PAYMENT_SERVICE_ITEM_STATUS;

const ReviewServiceItems = ({
  header,
  paymentRequest,
  serviceItemCards,
  handleClose,
  patchPaymentServiceItem,
  onCompleteReview,
  completeReviewError,
  TACs,
  SACs,
}) => {
  const requestReviewed = paymentRequest?.status !== PAYMENT_REQUEST_STATUS.PENDING;

  const sortedCards = sortServiceItemsByGroup(serviceItemCards);

  const totalCards = sortedCards.length;

  const [curCardIndex, setCardIndex] = useState(requestReviewed ? totalCards : 0);

  const handleServiceItemNavBtnClick = (index) => {
    setCardIndex(index);
  };

  const handleAuthorizePaymentClick = (allServiceItemsRejected = false) => {
    onCompleteReview(allServiceItemsRejected);
  };

  const findAdditionalServiceItemData = (mtoServiceItemCode) => {
    const serviceItemCard = serviceItemCards?.find((item) => item.mtoServiceItemCode === mtoServiceItemCode);

    const additionalServiceItems = serviceItemCard
      ? serviceItemCard.mtoServiceItems?.find((mtoItem) => mtoItem.reServiceCode === mtoServiceItemCode)
      : [];

    return additionalServiceItems;
  };

  // calculating the sums
  const approvedSum = sortedCards.filter((s) => s.status === APPROVED).reduce((sum, cur) => sum + cur.amount, 0);
  const rejectedSum = sortedCards.filter((s) => s.status === DENIED).reduce((sum, cur) => sum + cur.amount, 0);
  const requestedSum = sortedCards.reduce((sum, cur) => sum + cur.amount, 0);

  let itemsNeedsReviewLength;
  let showNeedsReview;
  let allServiceItemsRejected;
  let firstItemNeedsReviewIndex;

  if (!requestReviewed) {
    itemsNeedsReviewLength = sortedCards.filter((s) => s.status === REQUESTED)?.length;
    showNeedsReview = sortedCards.some((s) => s.status === REQUESTED);
    allServiceItemsRejected = sortedCards.every((s) => s.status === DENIED);
    firstItemNeedsReviewIndex = showNeedsReview && sortedCards.findIndex((s) => s.status === REQUESTED);
  }

  const displayCompleteReview = curCardIndex === totalCards;

  const currentCard = !displayCompleteReview && sortedCards[parseInt(curCardIndex, 10)];

  // Determines which ReviewDetailsCard will be shown
  let renderCompleteAction = (
    <AuthorizePayment amount={approvedSum} onClick={() => handleAuthorizePaymentClick(requestReviewed)} />
  );
  if (requestReviewed) {
    renderCompleteAction = null;
  } else if (showNeedsReview) {
    renderCompleteAction = (
      <NeedsReview numberOfItems={itemsNeedsReviewLength} onClick={() => setCardIndex(firstItemNeedsReviewIndex)} />
    );
  } else if (allServiceItemsRejected) {
    renderCompleteAction = <RejectRequest onClick={() => handleAuthorizePaymentClick(allServiceItemsRejected)} />;
  }

  if (displayCompleteReview)
    return (
      <div data-testid="ReviewServiceItems" className={styles.ReviewServiceItems}>
        <div className={styles.top}>
          <Button
            data-testid="closeSidebar"
            type="button"
            onClick={handleClose}
            unstyled
            aria-label="Close Service Item review"
          >
            <FontAwesomeIcon icon="times" title="Close Service Item review" alt=" " />
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
            authorized={requestReviewed}
            dateAuthorized={paymentRequest?.reviewedAt}
            TACs={TACs}
            SACs={SACs}
            cards={sortedCards}
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
        <Button
          data-testid="closeSidebar"
          type="button"
          onClick={handleClose}
          unstyled
          aria-label="Close Service Item review"
        >
          <FontAwesomeIcon icon="times" alt=" " />
        </Button>
        <div data-testid="itemCount" className={styles.eyebrowTitle}>
          {curCardIndex + 1} OF {totalCards} ITEMS
        </div>
        <h2 className={styles.header}>{header}</h2>
      </div>
      <div className={styles.body}>
        {currentCard && (
          <ServiceItemCard
            key={`serviceItemCard_${currentCard.id}`}
            patchPaymentServiceItem={patchPaymentServiceItem}
            // eslint-disable-next-line react/jsx-props-no-spreading
            {...currentCard}
            requestComplete={requestReviewed}
            additionalServiceItemData={findAdditionalServiceItemData(currentCard.mtoServiceItemCode)}
          />
        )}
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
  onCompleteReview: PropTypes.func.isRequired,
  completeReviewError: PropTypes.shape({
    detail: PropTypes.string,
    title: PropTypes.string,
  }),
  mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
  paymentRequest: undefined,
  serviceItemCards: [],
  completeReviewError: undefined,
  mtoServiceItems: [],
  TACs: {},
  SACs: {},
};

export default ReviewServiceItems;
