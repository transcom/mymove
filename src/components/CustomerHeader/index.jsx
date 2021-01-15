import React from 'react';
import { string } from 'prop-types';
import classNames from 'classnames/bind';

import styles from './index.module.scss';

import { MoveOrderShape, CustomerShape } from 'types/moveOrder';
import { formatCustomerDate } from 'utils/formatters';

const cx = classNames.bind(styles);

const CustomerHeader = ({ customer, moveOrder, moveCode }) => {
  return (
    <div className={cx('cust-header')}>
      <div>
        <div className={cx('name-block')}>
          <h2>
            {customer.last_name}, {customer.first_name}
          </h2>
          <span className="usa-tag usa-tag--cyan usa-tag--large">#{moveCode}</span>
        </div>
        <div>
          <p>
            {moveOrder.departmentIndicator} {moveOrder.grade}
            <span className={cx('vertical-bar')}>|</span>
            DoD ID {customer.dodID}
          </p>
        </div>
      </div>
      <div className={cx('info-block')}>
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
