import React, { useState } from 'react';
import PropTypes from 'prop-types';

import styles from './ReviewServiceItems.module.scss';

import { ServiceItemCardsShape } from 'types/serviceItemCard';

const ReviewServiceItems = ({ header, serviceItemCards }) => {
  // const [curServiceItemCard] = useState(serviceItemCards[0]);
  const [curCardIndex] = useState(0);
  const totalCards = serviceItemCards.length;

  // debugging
  // console.log(curServiceItemCard);

  return (
    <div data-testid="ReviewServiceItems" className={styles.ReviewServiceItems}>
      <div className={styles.top}>
        <div className={styles.eyebrowTitle}>{`${curCardIndex + 1} OF ${totalCards} ITEMS`}</div>
        <h2 className={styles.header}>{header}</h2>
      </div>
      <div className={styles.body}>BODY</div>
      <div className={styles.bottom}>BOTTOM</div>
    </div>
  );
};

ReviewServiceItems.propTypes = {
  header: PropTypes.string,
  serviceItemCards: ServiceItemCardsShape.isRequired,
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
};

export default ReviewServiceItems;
