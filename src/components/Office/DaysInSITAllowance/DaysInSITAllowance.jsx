import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import { ShipmentPaymentSITBalanceShape } from '../../../types/serviceItems';

import styles from './DaysInSITAllowance.module.scss';

import { formatDate, formatDaysInTransit } from 'utils/formatters';

const DaysInSITAllowance = ({ className, shipmentPaymentSITBalance }) => {
  const {
    previouslyBilledDays,
    previouslyBilledEndDate,
    pendingSITDaysInvoiced,
    pendingBilledEndDate,
    totalSITDaysAuthorized,
    totalSITDaysRemaining,
    totalSITEndDate,
  } = shipmentPaymentSITBalance;
  return (
    <div className={classNames(className, styles.DaysInSITAllowance)} data-testid="DaysInSITAllowance">
      <dl className={styles.daysInSITList}>
        <dt>Prev. billed & accepted</dt>
        <dd data-testid="previouslyBilled">
          {formatDaysInTransit(previouslyBilledDays)}
          {!!previouslyBilledDays && (
            <>
              {', through '}
              {formatDate(previouslyBilledEndDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
            </>
          )}
        </dd>
        <dt>Invoiced & pending</dt>
        <dd data-testid="pendingInvoiced">
          {formatDaysInTransit(pendingSITDaysInvoiced)}, through{' '}
          {formatDate(pendingBilledEndDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
        </dd>
        <dt>Total authorized</dt>
        <dd>{formatDaysInTransit(totalSITDaysAuthorized)}</dd>
        <dt>Authorized remaining</dt>
        <dd data-testid="totalRemaining">
          {formatDaysInTransit(totalSITDaysRemaining)}, ends {formatDate(totalSITEndDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
        </dd>
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
