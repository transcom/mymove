import React from 'react';
import PropTypes from 'prop-types';

import styles from './ReviewDetailsCard.module.scss';

import { toDollarString } from 'shared/formatters';

/** This component represents a Payment Request Review Details Card shown at the end of navigation */
const ReviewDetailsCard = ({ children, completeReviewError, requestedAmount, acceptedAmount, rejectedAmount }) => {
  return (
    <div data-testid="ReviewDetailsCard" className={styles.ReviewDetailsCard}>
      <h4 className={styles.cardHeader}>Review details</h4>
      {completeReviewError && (
        <p className="text-error" data-testid="errorMessage">
          Error: {completeReviewError.detail}
        </p>
      )}
      <dl>
        <dt>Requested</dt>
        <dd data-testid="requested">{toDollarString(requestedAmount)}</dd>

        <dt>Accepted</dt>
        <dd data-testid="accepted">{toDollarString(acceptedAmount)}</dd>

        <dt>Rejected</dt>
        <dd data-testid="rejected">{toDollarString(rejectedAmount)}</dd>
      </dl>

      {children}
    </div>
  );
};

ReviewDetailsCard.propTypes = {
  children: PropTypes.node,
  completeReviewError: PropTypes.shape({
    detail: PropTypes.string,
  }),
  requestedAmount: PropTypes.number,
  acceptedAmount: PropTypes.number,
  rejectedAmount: PropTypes.number,
};

ReviewDetailsCard.defaultProps = {
  children: null,
  completeReviewError: null,
  requestedAmount: 0,
  acceptedAmount: 0,
  rejectedAmount: 0,
};

export default ReviewDetailsCard;
