import React from 'react';
import { string } from 'prop-types';

import styles from './index.module.scss';

import { OrderShape, CustomerShape } from 'types/order';
import { formatCustomerDate } from 'utils/formatters';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders.js';

const CustomerHeader = ({ customer, order, moveCode }) => {
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
            <span data-testid="deptRank" className={styles.details}>
              {ORDERS_BRANCH_OPTIONS[`${order.agency}`]} {ORDERS_RANK_OPTIONS[`${order.grade}`]}
            </span>
            <span className={styles.verticalBar}>|</span>
            <span data-testid="dodId" className={styles.details}>
              DoD ID {customer.dodID}
            </span>
          </p>
        </div>
      </div>
      <div data-testid="infoBlock" className={styles.infoBlock}>
        <div>
          <p>Authorized origin</p>
          <h4>{order.originDutyStation.name}</h4>
        </div>
        <div>
          <p>Authorized destination</p>
          <h4>{order.destinationDutyStation.name}</h4>
        </div>
        <div>
          <p>Report by</p>
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
