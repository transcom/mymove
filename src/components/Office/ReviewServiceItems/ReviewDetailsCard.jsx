import React from 'react';
import PropTypes from 'prop-types';

import styles from './ReviewDetailsCard.module.scss';
import PaymentReviewed from './PaymentReviewed';
import ReviewAccountingCodes from './ReviewAccountingCodes';

import { toDollarString } from 'shared/formatters';
import { AccountingCodesShape } from 'types/accountingCodes';
import { ServiceItemCardsShape } from 'types/serviceItems';

/** This component represents a Payment Request Review Details Card shown at the end of navigation */
const ReviewDetailsCard = ({
  children,
  completeReviewError,
  requestedAmount,
  acceptedAmount,
  rejectedAmount,
  authorized,
  dateAuthorized,
  TACs,
  SACs,
  cards,
}) => {
  return (
    <div data-testid="ReviewDetailsCard" className={styles.ReviewDetailsCard}>
      <h3 className={styles.cardHeader}>Review details</h3>

      {authorized && <PaymentReviewed authorizedAmount={acceptedAmount} dateAuthorized={dateAuthorized} />}

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

      <ReviewAccountingCodes TACs={TACs} SACs={SACs} cards={cards} />

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
  authorized: PropTypes.bool,
  dateAuthorized: PropTypes.string,
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  cards: ServiceItemCardsShape,
};

ReviewDetailsCard.defaultProps = {
  children: null,
  completeReviewError: null,
  requestedAmount: 0,
  acceptedAmount: 0,
  rejectedAmount: 0,
  authorized: false,
  dateAuthorized: null,
  TACs: {},
  SACs: {},
  cards: [],
};

export default ReviewDetailsCard;
