import React from 'react';

import styles from './ReviewDetailsCard.module.scss';

import { toDollarString } from 'shared/formatters';

/** This component represents a Payment Request Review Details Card shown at the end of navigation */
const ReviewDetailsCard = () => {
  return (
    <div data-testid="ReviewDetailsCard" className={styles.ReviewDetailsCard}>
      <h4 className={styles.cardHeader}>Review details</h4>
      <dl>
        <dt>Requested</dt>
        <dd data-testid="requested">{toDollarString(1234.12)}</dd>

        <dt>Accepted</dt>
        <dd data-testid="accepted">{toDollarString(1234.12)}</dd>

        <dt>Rejected</dt>
        <dd data-testid="rejected">{toDollarString(1234.12)}</dd>
      </dl>
    </div>
  );
};

ReviewDetailsCard.propTypes = {};

ReviewDetailsCard.defaultProps = {};

export default ReviewDetailsCard;
