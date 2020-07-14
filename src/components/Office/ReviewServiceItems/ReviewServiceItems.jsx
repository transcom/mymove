import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './ReviewServiceItems.module.scss';

import { ServiceItemCardsShape } from 'types/serviceItemCard';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import ServiceItemCard from 'components/Office/ReviewServiceItems/ServiceItemCard';

// Sort ascending by objects with string iso timestamps
const dateCreatedSort = (a, b) => {
  return Date.parse(a.createdAt) - Date.parse(b.createdAt);
};

const sortByGroup = (serviceItemCards) => {
  // Will populate with earliest service item of each shipment id
  const shipmentOrder = [];
  // Contains sorted service items keyed by shipment id or undefined for basic items
  const shipmentServiceItems = {};

  serviceItemCards.sort(dateCreatedSort);

  serviceItemCards.map((serviceItem) => {
    const { shipmentId } = serviceItem;
    // We've already added the earliest service item for this shipment, continue until we get to the next
    if (shipmentServiceItems[`${shipmentId}`]) {
      return false;
    }

    shipmentServiceItems[`${shipmentId}`] = serviceItemCards.filter((item) => item.shipmentId === shipmentId);
    shipmentOrder.push(serviceItem);
    return true;
  });

  shipmentOrder.sort(dateCreatedSort);

  const sortedCards = [];
  shipmentOrder.map((shipment) => {
    sortedCards.push(...shipmentServiceItems[`${shipment.shipmentId}`]);
    return true;
  });

  return sortedCards;
};

const ReviewServiceItems = ({ header, serviceItemCards, handleClose }) => {
  const [curCardIndex, setCardIndex] = useState(0);
  // eslint-disable-next-line no-unused-vars
  const [totalApproved, setTotalApproved] = useState(0);
  const [sortedCards] = useState(sortByGroup(serviceItemCards));

  const totalCards = serviceItemCards.length;

  const handleClick = (index) => {
    setCardIndex(index);
  };

  return (
    <div data-cy="ReviewServiceItems" className={styles.ReviewServiceItems}>
      <div className={styles.top}>
        <Button data-testid="closeSidebar" type="button" onClick={handleClose} unstyled>
          <XLightIcon />
        </Button>
        <div className={styles.eyebrowTitle}>{`${curCardIndex + 1} OF ${totalCards} ITEMS`}</div>
        <h2 className={styles.header}>{header}</h2>
      </div>
      <div className={styles.body}>
        {/* eslint-disable-next-line react/jsx-props-no-spreading */}
        <ServiceItemCard {...sortedCards[parseInt(curCardIndex, 10)]} />
      </div>
      <div className={styles.bottom}>
        <Button
          data-cy="prevServiceItem"
          type="button"
          onClick={() => handleClick(curCardIndex - 1)}
          secondary
          disabled={curCardIndex === 0}
        >
          Previous
        </Button>
        <Button
          data-cy="nextServiceItem"
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
  serviceItemCards: ServiceItemCardsShape.isRequired,
  handleClose: PropTypes.func.isRequired,
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
};

export default ReviewServiceItems;
