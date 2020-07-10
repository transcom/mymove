import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './ReviewServiceItems.module.scss';

import { ServiceItemCardsShape } from 'types/serviceItemCard';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';

const ReviewServiceItems = ({ header, serviceItemCards, handleClose }) => {
  // const [curServiceItemCard] = useState(serviceItemCards[0]);
  const [curCardIndex] = useState(0);
  const totalCards = serviceItemCards.length;

  // debugging
  // console.log(curServiceItemCard);

  return (
    <div data-testid="ReviewServiceItems" className={styles.ReviewServiceItems}>
      <div className={styles.top}>
        <Button data-testid="closeSidebar" type="button" onClick={handleClose} unstyled>
          <XLightIcon />
        </Button>
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
  handleClose: PropTypes.func.isRequired,
};

ReviewServiceItems.defaultProps = {
  header: 'Review service items',
};

export default ReviewServiceItems;
