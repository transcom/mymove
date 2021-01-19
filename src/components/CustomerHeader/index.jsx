import React from 'react';
import { string } from 'prop-types';

import styles from './index.module.scss';

import { MoveOrderShape, CustomerShape } from 'types/moveOrder';
import { formatCustomerDate } from 'utils/formatters';

const CustomerHeader = ({ customer, moveOrder, moveCode }) => {
  return (
    <div className={styles.custHeader}>
      <div>
        <div data-test="nameBlock" className={styles.nameBlock}>
          <h2>
            {customer.last_name}, {customer.first_name}
          </h2>
          <span className="usa-tag usa-tag--cyan usa-tag--large">#{moveCode}</span>
        </div>
        <div>
          <p>
            <span data-test="deptRank" className={styles.details}>
              {moveOrder.departmentIndicator} {moveOrder.grade}
            </span>
            <span className={styles.verticalBar}>|</span>
            <span data-test="dodId" className={styles.details}>
              DoD ID {customer.dodID}
            </span>
          </p>
        </div>
      </div>
      <div data-test="infoBlock" className={styles.infoBlock}>
        <div>
          <p>Authorized origin</p>
          <h4>{moveOrder.originDutyStation.name}</h4>
        </div>
        <div>
          <p>Authorized destination</p>
          <h4>{moveOrder.destinationDutyStation.name}</h4>
        </div>
        <div>
          <p>Report by</p>
          <h4>{formatCustomerDate(moveOrder.report_by_date)}</h4>
        </div>
      </div>
    </div>
  );
};

CustomerHeader.propTypes = {
  customer: CustomerShape.isRequired,
  moveOrder: MoveOrderShape.isRequired,
  moveCode: string.isRequired,
};

export default CustomerHeader;
