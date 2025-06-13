import React, { useRef } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ReviewServiceItems.module.scss';
import ServiceItemCard from './ServiceItemCard';
import ReviewDetailsCard from './ReviewDetailsCard';
import AuthorizePayment from './AuthorizePayment';
import NeedsReview from './NeedsReview';
import RejectRequest from './RejectRequest';

import Alert from 'shared/Alert';
import { ServiceItemCardsShape } from 'types/serviceItems';
import { MTOServiceItemShape } from 'types/order';
import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
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
  curCardIndex,
  setCardIndex,
  handlePrevious,
  requestReviewed,
  setShouldAdvanceOnSubmit,
}) => {
  const formRef = useRef();

  const totalCards = serviceItemCards.length;

  const displayCompleteReview = curCardIndex === totalCards;
  const handleNext = () => {
    setShouldAdvanceOnSubmit(true);
    if (formRef.current && !requestReviewed) {
      formRef.current.handleSubmit();
    } else {
      setCardIndex(curCardIndex + 1);
    }
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
  const approvedSum = serviceItemCards.filter((s) => s.status === APPROVED).reduce((sum, cur) => sum + cur.amount, 0);
  const rejectedSum = serviceItemCards.filter((s) => s.status === DENIED).reduce((sum, cur) => sum + cur.amount, 0);
  const requestedSum = serviceItemCards.reduce((sum, cur) => sum + cur.amount, 0);

  let itemsNeedsReviewLength;
  let showNeedsReview;
  let allServiceItemsRejected;
  let firstItemNeedsReviewIndex;

  if (!requestReviewed) {
    itemsNeedsReviewLength = serviceItemCards.filter((s) => s.status === REQUESTED)?.length;
    showNeedsReview = serviceItemCards.some((s) => s.status === REQUESTED);
    allServiceItemsRejected = serviceItemCards.every((s) => s.status === DENIED);
    firstItemNeedsReviewIndex = showNeedsReview && serviceItemCards.findIndex((s) => s.status === REQUESTED);
  }

  const currentCard = !displayCompleteReview && serviceItemCards[parseInt(curCardIndex, 10)];

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
            cards={serviceItemCards}
          >
            {renderCompleteAction}
          </ReviewDetailsCard>
        </div>
        <div className={styles.bottom}>
          <Button
            data-testid="prevServiceItem"
            type="button"
            onClick={handlePrevious}
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
            formRef={formRef}
            setShouldAdvanceOnSubmit={setShouldAdvanceOnSubmit}
          />
        )}
      </div>
      <div className={styles.bottom}>
        <div className={styles.totalApproved}>
          <div className={styles.totalLabel}>Total approved</div>
          <div className={styles.totalAmount} data-testid="approvedAmount">
            {toDollarString(approvedSum)}
          </div>
        </div>
        <div className={styles.navBtns}>
          <Button
            data-testid="prevServiceItem"
            aria-label="Previous Service Item"
            type="button"
            onClick={() => setCardIndex(curCardIndex - 1)}
            secondary
            disabled={curCardIndex === 0}
          >
            Previous
          </Button>
          <Button
            data-testid="nextServiceItem"
            aria-label="Next Service Item"
            type="button"
            onClick={handleNext}
            disabled={curCardIndex === totalCards}
          >
            Next
          </Button>
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
  curCardIndex: PropTypes.number,
  setCardIndex: PropTypes.func.isRequired,
  handlePrevious: PropTypes.func.isRequired,
  requestReviewed: PropTypes.bool.isRequired,
  setShouldAdvanceOnSubmit: PropTypes.func.isRequired,
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
  paymentRequest: undefined,
  serviceItemCards: [],
  completeReviewError: undefined,
  mtoServiceItems: [],
  TACs: {},
  SACs: {},
  curCardIndex: 0,
};

export default ReviewServiceItems;
