import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './AuthorizePayment.module.scss';

import { toDollarString } from 'shared/formatters';

/**
 * This component represents a section shown in the ReviewDetailsCard at the end of navigation.
 * Only shows if all service items have been reviewed yet and at least 1 is approved.
 * */
const AuthorizePayment = ({ amount, onClick }) => {
  return (
    <div data-testid="AuthorizePayment" className={styles.AuthorizePayment}>
      <div className={styles.content}>{`Do you authorize this payment of ${toDollarString(amount)}?`}</div>
      <Button data-testid="authorizePaymentBtn" type="button" onClick={onClick}>
        Authorize payment
      </Button>
    </div>
  );
};

AuthorizePayment.propTypes = {
  amount: PropTypes.number,
  onClick: PropTypes.func,
};

AuthorizePayment.defaultProps = {
  amount: 0,
  onClick: null,
};

export default AuthorizePayment;
