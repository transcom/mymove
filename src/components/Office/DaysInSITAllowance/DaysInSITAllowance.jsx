import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import { ShipmentPaymentSITBalanceShape } from '../../../types/serviceItems';

import styles from './DaysInSITAllowance.module.scss';

import { formatDate } from 'utils/formatters';

const DaysInSITAllowance = ({ className, shipmentPaymentSITBalance }) => {
  const {
    previouslyBilledDays,
    pendingBilledEndDate,
    pendingBilledStartDate,
    totalSITDaysAuthorized,
    totalSITDaysRemaining,
  } = shipmentPaymentSITBalance;
  return (
    <div className={classNames(className, styles.DaysInSITAllowance)} data-testid="DaysInSITAllowance">
      <dl className={styles.daysInSITList}>
        <dt>Prev. billed & accepted days</dt>
        <dd data-testid="previouslyBilled">{previouslyBilledDays || '0'}</dd>
        <dt>Payment start - end date</dt>
        <dd data-testid="pendingBilledStartEndDate">
          {formatDate(pendingBilledStartDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
          {` - `}
          {formatDate(pendingBilledEndDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
        </dd>
        <dt>Total days of SIT approved</dt>
        <dd>{totalSITDaysAuthorized}</dd>
        <dt>Total approved days remaining</dt>
        <dd data-testid="totalRemaining">{totalSITDaysRemaining > 0 ? totalSITDaysRemaining : '0'}</dd>
      </dl>
    </div>
  );
};

DaysInSITAllowance.propTypes = {
  className: PropTypes.string,
  shipmentPaymentSITBalance: ShipmentPaymentSITBalanceShape.isRequired,
};

DaysInSITAllowance.defaultProps = {
  className: undefined,
};

export default DaysInSITAllowance;
