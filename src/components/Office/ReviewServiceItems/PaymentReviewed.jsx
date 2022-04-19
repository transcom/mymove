import React from 'react';
import PropTypes from 'prop-types';

import styles from './PaymentReviewed.module.scss';

import { toDollarString, formatDateFromIso } from 'utils/formatters';

/**
 * This component represents a section shown in the ReviewDetailsCard at the end of navigation.
 * Only shows after a payment request has been authorized or rejected and moved to reviewed status.
 * */
const PaymentReviewed = ({ authorizedAmount, dateAuthorized }) => {
  return (
    <div data-testid="PaymentReviewed" className={styles.PaymentReviewed}>
      <div data-testid="content" className={styles.content}>
        {authorizedAmount > 0 ? (
          <>
            <p data-testid="paymentAuthorizedAmt">Payment authorized: {toDollarString(authorizedAmount)}</p>
            <p data-testid="reviewedOn">On: {formatDateFromIso(dateAuthorized, 'DD MMM YYYY')}</p>
          </>
        ) : (
          <p data-testid="paymentAuthorizedAmt">Payment authorized: none</p>
        )}
      </div>
    </div>
  );
};

PaymentReviewed.propTypes = {
  authorizedAmount: PropTypes.number,
  dateAuthorized: PropTypes.string,
};

PaymentReviewed.defaultProps = {
  authorizedAmount: 0,
  dateAuthorized: null,
};

export default PaymentReviewed;
