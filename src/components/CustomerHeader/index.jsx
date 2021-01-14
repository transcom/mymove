import React from 'react';
import { string } from 'prop-types';
import classNames from 'classnames/bind';
import { connect } from 'react-redux';

import styles from './index.module.scss';

import { MoveOrderShape, CustomerShape } from 'types/moveOrder';
import { getMoveByLocator as getMoveByLocatorAction, selectMoveByLocator } from 'shared/Entities/modules/moves';
import {
  getCustomer as getCustomerAction,
  selectCustomer,
  selectMoveOrder,
} from 'shared/Entities/modules/moveTaskOrders';

const cx = classNames.bind(styles);

const CustomerHeader = ({ customer, moveOrder, moveCode }) => (
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
        <h4>27 Mar 2020</h4>
      </div>
    </div>
  </div>
);

CustomerHeader.propTypes = {
  customer: CustomerShape.isRequired,
  moveOrder: MoveOrderShape.isRequired,
  moveCode: string.isRequired,
};

function mapStateToProps(state) {
  const moveCode = 'FKLCTR'; // temp until retrieval is worked out
  const move = selectMoveByLocator(state, moveCode);
  const moveOrderId = move?.ordersId;
  const moveOrder = selectMoveOrder(state, moveOrderId);
  const customerId = move?.customerID;
  const customer = selectCustomer(state, customerId);
  return {
    moveOrder,
    customer,
    moveCode,
  };
}

const mapDispatchToProps = {
  getMoveByLocator: getMoveByLocatorAction,
  getCustomer: getCustomerAction,
};

// in order to avoid setting up proxy server only for storybook, pass in stub function so API requests don't fail
const mergeProps = (stateProps, dispatchProps, ownProps) => ({
  ...stateProps,
  ...dispatchProps,
  ...ownProps,
});
export default connect(mapStateToProps, mapDispatchToProps, mergeProps)(CustomerHeader);
