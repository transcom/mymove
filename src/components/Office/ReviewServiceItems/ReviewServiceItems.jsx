import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import sortServiceItemsByGroup from '../../../utils/serviceItems';

import styles from './ReviewServiceItems.module.scss';

import { ServiceItemCardsShape } from 'types/serviceItemCard';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import ServiceItemCard from 'components/Office/ReviewServiceItems/ServiceItemCard';

const ReviewServiceItems = ({ header, serviceItemCards, handleClose, patchPaymentServiceItem }) => {
  const [curCardIndex, setCardIndex] = useState(0);
  // eslint-disable-next-line no-unused-vars
  const [totalApproved, setTotalApproved] = useState(0);

  const sortedCards = sortServiceItemsByGroup(serviceItemCards);

  const totalCards = serviceItemCards.length;

  const handleClick = (index) => {
    setCardIndex(index);
  };

  const formValues = {};
  // TODO - preset these based on existing values
  serviceItemCards.forEach((serviceItem) => {
    formValues[serviceItem.id] = {
      status: serviceItem.status,
      rejectionReason: undefined,
    };
  });

  const currentCard = sortedCards[parseInt(curCardIndex, 10)];

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
        {currentCard && (
          <ServiceItemCard
            key={`serviceItemCard_${currentCard.id}`}
            patchPaymentServiceItem={patchPaymentServiceItem}
            // eslint-disable-next-line react/jsx-props-no-spreading
            {...currentCard}
          />
        )}
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
          disabled={curCardIndex + 1 === totalCards}
        >
          Next
        </Button>
        <div className={styles.totalApproved}>
          <div className={styles.totalLabel}>Total approved</div>
          <div className={styles.totalAmount}>${totalApproved.toFixed(2)}</div>
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
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
  serviceItemCards: [],
};

export default ReviewServiceItems;
