import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';
import {
  updateMoveTaskOrderStatus,
  getAllMoveTaskOrders,
  getMoveOrder,
  getCustomer,
  selectCustomer,
  selectMoveOrder,
  selectMoveTaskOrders,
} from 'shared/Entities/modules/moveTaskOrders';
import { getMTOAgentList } from 'shared/Entities/modules/mtoAgents';
import { getMTOServiceItems, selectMTOServiceItems } from 'shared/Entities/modules/mtoServiceItems';
import { getMTOShipments, patchMTOShipmentStatus, selectMTOShipments } from 'shared/Entities/modules/mtoShipments';
import { ApproveRejectModal } from 'shared/ApproveRejectModal';
import { selectMTOAgents } from 'shared/Entities/modules/mtoAgents';

class CustomerDetails extends Component {
  componentDidMount() {
    const { customerId, moveOrderId } = this.props.match.params;
    this.props.getCustomer(customerId);
    this.props.getMoveOrder(moveOrderId).then(({ response: { body: moveOrder } }) => {
      this.props.getAllMoveTaskOrders(moveOrder.id).then(({ response: { body: moveTaskOrder } }) => {
        // TODO: would like to do batch fetching later
        moveTaskOrder.forEach((item) => this.props.getMTOServiceItems(item.id));
        moveTaskOrder.forEach((item) =>
          this.props.getMTOShipments(item.id).then(({ response: { body: mtoShipments } }) => {
            mtoShipments.forEach((shipment) => this.props.getMTOAgentList(item.id, shipment.id));
          }),
        );
      });
    });
  }
  render() {
    const { moveTaskOrder, customer, moveOrder, mtoServiceItems, mtoShipments, mtoAgents } = this.props;
    const entitlements = get(moveOrder, 'entitlement', {});
    return (
      <>
        <h1>Customer Details Page</h1>
        {customer && (
          <>
            <h2>Customer Info</h2>
            <dl>
              <dt>ID</dt>
              <dd>{get(customer, 'id')}</dd>
              <dt>First Name</dt>
              <dd>{get(customer, 'first_name')}</dd>
              <dt>Last Name</dt>
              <dd>{get(customer, 'last_name')}</dd>
              <dt>Email</dt>
              <dd>{get(customer, 'email')}</dd>
              <dt>Phone</dt>
              <dd>{get(customer, 'phone')}</dd>
              <dt>Current Address</dt>
              <dd>{JSON.stringify(get(customer, 'current_address'))}</dd>
              <dt>DOD ID</dt>
              <dd>{get(customer, 'dodID')}</dd>
              <dt>Agency</dt>
              <dd>{get(customer, 'agency')}</dd>
            </dl>
          </>
        )}
        {moveOrder && (
          <>
            <h2>Move Orders</h2>
            <dt>Order Number</dt>
            <dd>{get(moveOrder, 'order_number', '')}</dd>
            <dt>Order Type</dt>
            <dd>{get(moveOrder, 'order_type', '')}</dd>
            <dt>Order Type Detail</dt>
            <dd>{get(moveOrder, 'order_type_detail', '')}</dd>
            <dt>Date Issued</dt>
            <dd>{get(moveOrder, 'date_issued', '')}</dd>
            <dt>Report By Date</dt>
            <dd>{get(moveOrder, 'report_by_date', '')}</dd>
            <dt>Destination Duty Station</dt>
            <dd>{get(moveOrder, 'destinationDutyStation.name', '')}</dd>
            <dt>Destination Duty Station Address</dt>
            <dd>{JSON.stringify(get(moveOrder, 'destinationDutyStation.address', {}))} </dd>
            <dt>Origin Duty Station</dt>
            <dd>{get(moveOrder, 'originDutyStation.name', '')}</dd>
            <dt>Origin Duty Station Address</dt>
            <dd>{JSON.stringify(get(moveOrder, 'originDutyStation.address', {}))} </dd>

            <dt>Department Indicator</dt>
            <dd></dd>
            <dt>TAC / MDC</dt>
            <dd></dd>
            <dt>SAC / SDN</dt>
            <dd></dd>

            {entitlements && (
              <>
                <h2>Customer Entitlements</h2>
                <dl>
                  <dt>Rank</dt>
                  <dd>{get(moveOrder, 'grade', '')}</dd>
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
                  <dt>Total Weight</dt>
                  <dd>{get(entitlements, 'totalWeight')}</dd>
                  <dt>Authorized Weight</dt>
                  <dd>{get(entitlements, 'authorizedWeight')}</dd>
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
            <h3>Status: {moveTaskOrder.availableToPrimeAt ? 'Available to Prime' : 'Draft'}</h3>
            <dl>
              <dt>ID</dt>
              <dd>{get(moveTaskOrder, 'id')}</dd>
              <dt>Available to Prime At</dt>
              <dd>{get(moveTaskOrder, 'availableToPrimeAt', 'N/A')}</dd>
              <dt>Is Canceled</dt>
              <dd>{get(moveTaskOrder, 'isCanceled', false).toString()}</dd>
              <dt>Reference ID</dt>
              <dd>{get(moveTaskOrder, 'referenceId', '').toString()}</dd>
              <dt>Updated At</dt>
              <dd>{get(moveTaskOrder, 'updatedAt', '').toString()}</dd>
            </dl>

            <h2>Requested Shipments</h2>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Shipment Type</th>
                  <th>Requested Pick-up Date</th>
                  <th>Scheduled Pick-up Date</th>
                  <th>Pick up Address</th>
                  <th>Secondary Pickup Address</th>
                  <th>Delivery Address</th>
                  <th>Secondary Delivery Address</th>
                  <th>Customer Remarks</th>
                  <th>Status</th>
                  <th>Rejection Reason</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {mtoShipments.map((items) => (
                  <Fragment key={items.id}>
                    <tr>
                      <td>{items.id}</td>
                      <td>{items.shipmentType}</td>
                      <td>{items.requestedPickupDate}</td>
                      <td>{items.scheduledPickupDate}</td>
                      <td>
                        {get(items.pickupAddress, 'street_address_1')} {get(items.pickupAddress, 'street_address_2')}{' '}
                        {get(items.pickupAddress, 'street_address_3')} {} {get(items.pickupAddress, 'state')}{' '}
                        {get(items.pickupAddress, 'postal_code')}
                      </td>
                      <td>
                        {get(items.secondaryPickupAddress, 'street_address_1')}{' '}
                        {get(items.secondaryPickupAddress, 'street_address_2')}{' '}
                        {get(items.secondaryPickupAddress, 'street_address_3')} {}{' '}
                        {get(items.secondaryPickupAddress, 'state')} {get(items.secondaryPickupAddress, 'postal_code')}
                      </td>
                      <td>
                        {get(items.destinationAddress, 'street_address_1')}{' '}
                        {get(items.destinationAddress, 'street_address_2')}{' '}
                        {get(items.destinationAddress, 'street_address_3')} {} {get(items.destinationAddress, 'state')}{' '}
                        {get(items.destinationAddress, 'postal_code')}
                      </td>
                      <td>
                        {get(items.secondaryDeliveryAddress, 'street_address_1')}{' '}
                        {get(items.secondaryDeliveryAddress, 'street_address_2')}{' '}
                        {get(items.secondaryDeliveryAddress, 'street_address_3')} {}{' '}
                        {get(items.secondaryDeliveryAddress, 'state')}{' '}
                        {get(items.secondaryDeliveryAddress, 'postal_code')}
                      </td>
                      <td>{items.customerRemarks}</td>
                      <td>{items.status}</td>
                      <td>{items.rejectionReason}</td>
                      <td>
                        <ApproveRejectModal
                          showModal={items.status === 'SUBMITTED'}
                          approveBtnOnClick={() =>
                            this.props
                              .patchMTOShipmentStatus(get(moveTaskOrder, 'id'), items.id, 'APPROVED', items.eTag)
                              .then(() => this.props.getMTOServiceItems(items.moveTaskOrderID))
                          }
                          rejectBtnOnClick={(rejectionReason) =>
                            this.props.patchMTOShipmentStatus(
                              get(moveTaskOrder, 'id'),
                              items.id,
                              'REJECTED',
                              items.eTag,
                              rejectionReason,
                            )
                          }
                        />
                      </td>
                    </tr>
                  </Fragment>
                ))}
              </tbody>
            </table>

            <h2>Shipment Agents</h2>
            <table>
              <thead>
                <tr>
                  <th>id</th>
                  <th>Agent Type</th>
                  <th>First Name</th>
                  <th>Last Name</th>
                  <th>Email</th>
                  <th>Phone</th>
                </tr>
              </thead>
              <tbody>
                {mtoAgents.map((shipmentAgent) => (
                  <Fragment key={get(shipmentAgent[0], 'id')}>
                    <tr>
                      <td>{get(shipmentAgent[0], 'id')}</td>
                      <td>{get(shipmentAgent[0], 'agentType')}</td>
                      <td>{get(shipmentAgent[0], 'firstName')}</td>
                      <td>{get(shipmentAgent[0], 'lastName')}</td>
                      <td>{get(shipmentAgent[0], 'email')}</td>
                      <td>{get(shipmentAgent[0], 'phone')}</td>
                    </tr>
                  </Fragment>
                ))}
              </tbody>
            </table>

            <h2>MTO Service Items</h2>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Move Task Order ID</th>
                  <th>Rate Engine Service ID</th>
                  <th>Rate Engine Service Code</th>
                  <th>Rate Engine Service Name</th>
                </tr>
              </thead>
              <tbody>
                {mtoServiceItems.map((items) => (
                  <Fragment key={items.id}>
                    <tr>
                      <td>{items.id}</td>
                      <td>{items.moveTaskOrderID}</td>
                      <td>{items.reServiceID}</td>
                      <td>{items.reServiceCode}</td>
                      <td>{items.reServiceName}</td>
                    </tr>
                  </Fragment>
                ))}
              </tbody>
            </table>

            <div>
              <button
                onClick={() => this.props.updateMoveTaskOrderStatus(moveTaskOrder.id, moveTaskOrder.eTag, ['MS', 'CS'])}
              >
                Send to Prime
              </button>
            </div>
          </>
        )}
      </>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  const moveOrderId = ownProps.match.params.moveOrderId;
  const moveOrder = selectMoveOrder(state, moveOrderId);
  const moveTaskOrders = selectMoveTaskOrders(state, moveOrderId);
  let moveTaskOrderId;
  if (moveTaskOrders[0]) {
    moveTaskOrderId = moveTaskOrders[0].id;
  }
  return {
    moveOrder,
    customer: selectCustomer(state, ownProps.match.params.customerId),
    mtoServiceItems: selectMTOServiceItems(state, moveOrderId),
    mtoShipments: selectMTOShipments(state, moveOrderId),
    // TODO: Change when we start making use of multiple move task orders
    moveTaskOrder: moveTaskOrders[0],
    mtoAgents: selectMTOAgents(state, moveTaskOrderId),
  };
};

const mapDispatchToProps = {
  getMoveOrder,
  getAllMoveTaskOrders,
  updateMoveTaskOrderStatus,
  getCustomer,
  getMTOServiceItems,
  getMTOShipments,
  patchMTOShipmentStatus,
  getMTOAgentList,
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomerDetails);
