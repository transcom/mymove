import React, { Component } from 'react';
import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';
import {
  getEntitlements,
  updateMoveTaskOrderStatus,
  getMoveTaskOrder,
  getCustomer,
  selectMoveTaskOrder,
} from 'shared/Entities/modules/moveTaskOrders';
import { selectCustomer } from 'shared/Entities/modules/customer';

const fakeMoveTaskOrderID = '5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
class CustomerDetails extends Component {
  componentDidMount() {
    this.props.getEntitlements('fake_move_task_order_id');
    this.props.getMoveTaskOrder(fakeMoveTaskOrderID);
    this.props.getCustomer(this.props.match.params.customerId);
  }

  render() {
    const { entitlements, moveTaskOrder, customer } = this.props;
    return (
      <>
        <h1>Customer Details Page</h1>
        {customer && (
          <>
            <h2>Customer Info</h2>
            <dl>
              <dt>ID</dt>
              <dd>{get(customer, 'id')}</dd>
              <dt>DOD ID</dt>
              <dd>{get(customer, 'dodID')}</dd>
            </dl>
          </>
        )}
        {entitlements && (
          <>
            <h2>Customer Entitlements</h2>
            <dl>
              <dt>Weight Entitlement</dt>
              <dd>{entitlements.totalWeightSelf}</dd>
              <dt>SIT Entitlement</dt>
              <dd>{entitlements.storageInTransit}</dd>
            </dl>
          </>
        )}
        {!isEmpty(moveTaskOrder) && (
          <>
            <h2>Move Task Order</h2>
            <dl>
              <dt>ID</dt>
              <dd>{get(moveTaskOrder, 'id')}</dd>
              <dt>Reference ID</dt>
              <dd>{get(moveTaskOrder, 'referenceId')}</dd>
              <dt>Is Available to Prime</dt>
              <dd>{get(moveTaskOrder, 'isAvailableToPrime').toString()}</dd>
              <dt>Is Canceled</dt>
              <dd>{get(moveTaskOrder, 'isCanceled').toString()}</dd>
            </dl>
          </>
        )}
        <div>
          <button onClick={() => this.props.updateMoveTaskOrderStatus(fakeMoveTaskOrderID, 'DRAFT')}>
            Generate MTO
          </button>
        </div>
      </>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  //TODO hard coding mto for now
  const fakeMoveTaskOrderID = '5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
  const entitlements = get(state, 'entities.entitlements');
  const moveTaskOrder = selectMoveTaskOrder(state, fakeMoveTaskOrderID);
  return {
    entitlements: entitlements && Object.values(entitlements).length > 0 ? Object.values(entitlements)[0] : null,
    moveTaskOrder,
    customer: selectCustomer(state, ownProps.match.params.customerId),
  };
};

const mapDispatchToProps = {
  getEntitlements,
  getMoveTaskOrder,
  updateMoveTaskOrderStatus,
  getCustomer,
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerDetails);
