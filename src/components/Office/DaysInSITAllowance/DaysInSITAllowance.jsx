import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './DaysInSITAllowance.module.scss';

import { formatDaysInTransit, formatDate } from 'shared/formatters';

const DaysInSITAllowance = ({
  className,
  previouslyBilledDays,
  previouslyBilledEndDate,
  pendingSITDaysInvoiced,
  pendingBilledEndDate,
  totalSITDaysAuthorized,
  totalSITDaysRemaining,
  totalSITEndDate,
}) => {
  return (
    <dl className={classNames(className, styles.DaysInSITAllowance)}>
      <dt>Prev. billed & accepted</dt>
      <dd>
        {formatDaysInTransit(previouslyBilledDays)}
        {!!previouslyBilledDays && (
          <>
            {', through '}
            {formatDate(previouslyBilledEndDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
          </>
        )}
      </dd>
      <dt>Invoiced & pending</dt>
      <dd>
        {formatDaysInTransit(pendingSITDaysInvoiced)}, through{' '}
        {formatDate(pendingBilledEndDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
      </dd>
      <dt>Total authorized</dt>
      <dd>{formatDaysInTransit(totalSITDaysAuthorized)}</dd>
      <dt>Authorized remaining</dt>
      <dd>
        {formatDaysInTransit(totalSITDaysRemaining)}, ends {formatDate(totalSITEndDate, 'YYYY-MM-DD', 'DD MMM YYYY')}
      </dd>
    </dl>
  );
};

DaysInSITAllowance.propTypes = {
  className: PropTypes.string,
  previouslyBilledDays: PropTypes.number.isRequired,
  previouslyBilledEndDate: PropTypes.string,
  pendingSITDaysInvoiced: PropTypes.number.isRequired,
  pendingBilledEndDate: PropTypes.string.isRequired,
  totalSITDaysAuthorized: PropTypes.number.isRequired,
  totalSITDaysRemaining: PropTypes.number.isRequired,
  totalSITEndDate: PropTypes.string.isRequired,
};

DaysInSITAllowance.defaultProps = {
  className: undefined,
  previouslyBilledEndDate: undefined,
};

export default DaysInSITAllowance;
