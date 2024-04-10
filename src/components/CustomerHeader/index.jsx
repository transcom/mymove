import React from 'react';
import { string } from 'prop-types';

import styles from './index.module.scss';

import { OrderShape, CustomerShape } from 'types/order';
import { formatCustomerDate, formatLabelReportByDate } from 'utils/formatters';
import {
  CHECK_SPECIAL_ORDERS_TYPES,
  ORDERS_BRANCH_OPTIONS,
  ORDERS_PAY_GRADE_OPTIONS,
  SPECIAL_ORDERS_TYPES,
} from 'constants/orders.js';

const CustomerHeader = ({ customer, order, moveCode }) => {
  // eslint-disable-next-line camelcase
  const { order_type: orderType } = order;

  const isRetireeOrSeparatee = ['RETIREMENT', 'SEPARATION'].includes(orderType);
  const isSpecialMove = CHECK_SPECIAL_ORDERS_TYPES(orderType);

  /**
   * Depending on the order type, this row dt label can be either:
   * Report by date (PERMANENT_CHANGE_OF_STATION)
   * Date of retirement (RETIREMENT)
   * Date of separation (SEPARATION)
   */
  const reportDateLabel = formatLabelReportByDate(orderType);

  return (
    <div className={styles.custHeader}>
      <div>
        <div data-testid="nameBlock" className={styles.nameBlock}>
          <h2>
            {customer.last_name}, {customer.first_name}
          </h2>
          <span className="usa-tag usa-tag--cyan usa-tag--large">#{moveCode}</span>
        </div>
        <div>
          <p>
            <span data-testid="deptPayGrade" className={styles.details}>
              {ORDERS_BRANCH_OPTIONS[`${order.agency}`]} {ORDERS_PAY_GRADE_OPTIONS[`${order.grade}`]}
            </span>
            <span className={styles.verticalBar}>|</span>
            <span data-testid="dodId" className={styles.details}>
              DoD ID {customer.dodID}
            </span>
          </p>
        </div>
      </div>
      {isSpecialMove ? (
        <div data-testid="specialMovesLabel" className={styles.specialMovesLabel}>
          <p>{SPECIAL_ORDERS_TYPES[`${orderType}`]}</p>
        </div>
      ) : null}
      <div data-testid="infoBlock" className={styles.infoBlock}>
        <div>
          <p>Authorized origin</p>
          <h4>{order.originDutyLocation.name}</h4>
        </div>
        {order.destinationDutyLocation.name && (
          <div>
            <p data-testid="destinationLabel">
              {isRetireeOrSeparatee ? 'HOR, HOS or PLEAD' : 'Authorized destination'}
            </p>
            <h4>{order.destinationDutyLocation.name}</h4>
          </div>
        )}
        <div>
          <p data-testid="reportDateLabel">{reportDateLabel}</p>
          <h4>{formatCustomerDate(order.report_by_date)}</h4>
        </div>
      </div>
    </div>
  );
};

CustomerHeader.propTypes = {
  customer: CustomerShape.isRequired,
  order: OrderShape.isRequired,
  moveCode: string.isRequired,
};

export default CustomerHeader;
