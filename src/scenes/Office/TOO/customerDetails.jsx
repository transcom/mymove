import React, { Component } from 'react';
import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';
import {
  updateMoveTaskOrderStatus,
  getAllMoveTaskOrders,
  getMoveOrder,
  getCustomer,
  selectCustomer,
} from 'shared/Entities/modules/moveTaskOrders';
import { selectMoveOrder } from 'shared/Entities/modules/moveTaskOrders';

class CustomerDetails extends Component {
  componentDidMount() {
    const { customerId, moveOrderId } = this.props.match.params;
    this.props.getCustomer(customerId);
    this.props.getMoveOrder(moveOrderId).then(({ response: { body: moveOrder } }) => {
      this.props.getAllMoveTaskOrders(moveOrder.id);
    });
  }
  render() {
    const { moveTaskOrder, customer, moveOrder } = this.props;
    const entitlements = get(moveOrder, 'entitlement', {});
    return (
      <>
        <h1>Customer Details Page</h1>
        {customer && (
          <>
            <h2>Customer Info</h2>
            <dl>
              <dt>First Name</dt>
              <dd>{get(customer, 'first_name')}</dd>
              <dt>Last Name</dt>
              <dd>{get(customer, 'last_name')}</dd>
              <dt>ID</dt>
              <dd>{get(customer, 'id')}</dd>
              <dt>DOD ID</dt>
              <dd>{get(customer, 'dodID')}</dd>
            </dl>
          </>
        )}
        {moveOrder && (
          <>
            <h2>Move Orders</h2>
            <dt>Destination Duty Station</dt>
            <dd>{get(moveOrder, 'destinationDutyStation.name', '')}</dd>
            <dt>Destination Duty Station Address</dt>
            <dd>{JSON.stringify(get(moveOrder, 'destinationDutyStation.address', {}))} </dd>
            <dt>Origin Duty Station</dt>
            <dd>{get(moveOrder, 'originDutyStation.name', '')}</dd>
            <dt>Origin Duty Station Address</dt>
            <dd>{JSON.stringify(get(moveOrder, 'originDutyStation.address', {}))} </dd>
            {entitlements && (
              <>
                <h2>Customer Entitlements</h2>
                <dl>
                  <dt>Dependents Authorized</dt>
                  <dd>{get(entitlements, 'dependentsAuthorized', '').toString()}</dd>
                  <dt>Non Temporary Storage</dt>
                  <dd>{get(entitlements, 'nonTemporaryStorage', '').toString()}</dd>
                  <dt>Privately Owned Vehicle</dt>
                  <dd>{get(entitlements, 'privatelyOwnedVehicle', '').toString()}</dd>
                  <dt>ProGear Weight Spouse</dt>
                  <dd>{get(entitlements, 'proGearWeightSpouse')}</dd>
                  <dt>Storage In Transit</dt>
                  <dd>{get(entitlements, 'storageInTransit', '').toString()}</dd>
                  <dt>Total Dependents</dt>
                  <dd>{get(entitlements, 'totalDependents')}</dd>
                </dl>
              </>
            )}
          </>
        )}
        {!isEmpty(moveTaskOrder) && (
          <>
            <h2>Move Task Order</h2>
            <h3>Status: {moveTaskOrder.isAvailableToPrime ? 'Available to Prime' : 'Draft'}</h3>
            <dl>
              <dt>ID</dt>
              <dd>{get(moveTaskOrder, 'id')}</dd>
              <dt>Is Available to Prime</dt>
              <dd>{get(moveTaskOrder, 'isAvailableToPrime').toString()}</dd>
              <dt>Is Canceled</dt>
              <dd>{get(moveTaskOrder, 'isCanceled', false).toString()}</dd>
            </dl>
            <div>
              <button onClick={() => this.props.updateMoveTaskOrderStatus(moveTaskOrder.id)}>Send to Prime</button>
            </div>
          </>
        )}
      </>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const moveOrder = selectMoveOrder(state, ownProps.match.params.moveOrderId);
  const customer = selectCustomer(state, ownProps.match.params.customerId);
  return {
    customer,
    moveOrder,
    // TODO: Change when we start making use of multiple move task orders
    moveTaskOrder: Object.values(get(state, 'entities.moveTaskOrder', {}))[0],
  };
};

const mapDispatchToProps = {
  getMoveOrder,
  getAllMoveTaskOrders,
  updateMoveTaskOrderStatus,
  getCustomer,
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerDetails);
