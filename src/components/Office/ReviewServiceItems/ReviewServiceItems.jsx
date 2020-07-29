import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import sortServiceItemsByGroup from '../../../utils/serviceItems';

import styles from './ReviewServiceItems.module.scss';
import ServiceItemCard from './ServiceItemCard';

import { ServiceItemCardsShape } from 'types/serviceItemCard';
import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { toDollarString } from 'shared/formatters';

const ReviewServiceItems = ({
  header,
  serviceItemCards,
  handleClose,
  disableScrollIntoView,
  patchPaymentServiceItem,
  onCompleteReview,
  completeReviewError,
}) => {
  const [curCardIndex, setCardIndex] = useState(0);

  const sortedCards = sortServiceItemsByGroup(serviceItemCards);

  const totalCards = serviceItemCards.length;

  const { APPROVED, REVIEWED } = PAYMENT_SERVICE_ITEM_STATUS;

  const handleClick = (index) => {
    setCardIndex(index);
  };

  const handleAuthorizePayment = () => {
    onCompleteReview(REVIEWED);
  };

  const approvedSum = serviceItemCards.filter((s) => s.status === APPROVED).reduce((sum, cur) => sum + cur.amount, 0);
  // const rejectedSum = serviceItemCards.filter((s) => s.status === DENIED).reduce((sum, cur) => sum + cur.amount, 0)
  // const requestedSum = serviceItemCards.reduce((sum, cur) => sum + cur.amount, 0); // TODO - use in Complete review screen

  let firstBasicIndex = null;
  let lastBasicIndex = null;
  sortedCards.forEach((serviceItem, index) => {
    // here we want to set the first and last index
    // of basic service items to know the bounds
    if (!serviceItem.shipmentType) {
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

  const displayCompleteReview = curCardIndex === totalCards;

  const currentCard = !displayCompleteReview && sortedCards[parseInt(curCardIndex, 10)];

  const isBasicServiceItem =
    firstBasicIndex !== null && curCardIndex >= firstBasicIndex && curCardIndex <= lastBasicIndex;

  // Similar to componentDidMount and componentDidUpdate
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
            <XLightIcon />
          </Button>
          <h2 className={styles.header}>Complete request</h2>
        </div>
        <div className={styles.body}>
          <div className={styles.completeReviewCard}>
            <h4>Review details</h4>
            {completeReviewError && (
              <p className="text-error" data-testid="errorMessage">
                Error: {completeReviewError.detail}
              </p>
            )}

            <div className={styles.completeReviewAction}>
              <p>
                <strong>Do you authorize this payment of {toDollarString(approvedSum)}?</strong>
              </p>
              <Button type="button" data-testid="authorizePaymentBtn" onClick={handleAuthorizePayment}>
                Authorize Payment
              </Button>
            </div>
          </div>
        </div>
        <div className={styles.bottom}>
          <Button
            data-testid="prevServiceItem"
            type="button"
            onClick={() => handleClick(curCardIndex - 1)}
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
          <XLightIcon />
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
                // eslint-disable-next-line react/jsx-props-no-spreading
                {...curCard}
              />
            ))
          ) : (
            <ServiceItemCard
              key={`serviceItemCard_${currentCard.id}`}
              patchPaymentServiceItem={patchPaymentServiceItem}
              // eslint-disable-next-line react/jsx-props-no-spreading
              {...currentCard}
            />
          ))}
      </div>
      <div className={styles.bottom}>
        <Button
          data-testid="prevServiceItem"
          type="button"
          onClick={() => handleClick(curCardIndex - 1)}
          secondary
          disabled={curCardIndex === 0}
        >
          Previous
        </Button>
        <Button
          data-testid="nextServiceItem"
          type="button"
          onClick={() => handleClick(curCardIndex + 1)}
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
  serviceItemCards: [],
  disableScrollIntoView: false,
  completeReviewError: undefined,
};

export default ReviewServiceItems;
