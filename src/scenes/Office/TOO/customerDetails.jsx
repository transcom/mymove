import React from 'react';
import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';
import { getEntitlements, updateMoveTaskOrderStatus, getCustomerInfo } from 'shared/Entities/modules/moveTaskOrders';
import { selectMoveTaskOrder } from 'shared/Entities/modules/moveTaskOrders';

class CustomerDetails extends React.Component {
  componentDidMount() {
    this.props.getEntitlements('fake_move_task_order_id');
    this.props.getCustomerInfo('fake id');
  }

  render() {
    const { entitlements, moveTaskOrder, customer } = this.props;
    const fakeMoveTaskOrderID = '5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
    const NTS = entitlements && entitlements.nonTemporaryStorage ? 'Y' : 'N';
    const POV = entitlements && entitlements.privatelyOwnedVehicle ? 'Y' : 'N';
    const moveTaskOrderNonTemporaryStorage = get(moveTaskOrder, 'entitlements.nonTemporaryStorage') ? 'Y' : 'N';
    const moveTaskOrderPrivatelyOwnedVehicle = get(moveTaskOrder, 'entitlements.privatelyOwnedVehicle') ? 'Y' : 'N';
    return (
      <>
        <h1>Customer Deets Page</h1>
        {customer && (
          <>
            <h2>Customer Info</h2>
            <dl>
              <dt>Full Name</dt>
              <dd>
                {customer.first_name} {customer.middle_name} {customer.last_name}
              </dd>
              <dt>Service Branch / Agency</dt>
              <dd>{customer.agency}</dd>
              <dt>Rank / Grade</dt>
              <dd>{customer.grade}</dd>
              <dt>Email</dt>
              <dd>{customer.email}</dd>
              <dt>Phone</dt>
              <dd>{customer.telephone}</dd>
              <dt>Origin Duty Station</dt>
              <dd>{customer.origin_duty_station}</dd>
              <dt>Destination Duty Station</dt>
              <dd>{customer.destination_duty_station}</dd>
              <dt>Pickup Address</dt>
              <dd>{customer.pickup_address}</dd>
              {/* TODO does this belong on customer or something else?*/}
              <dt>Dependents Authorized</dt>
              <dd>{customer.dependents_authorized ? 'Y' : 'N'}</dd>
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
              <dt>NTS Entitlement</dt>
              <dd>{NTS}</dd>
              <dt>POV Entitlement</dt>
              <dd>{POV}</dd>
            </dl>
          </>
        )}
        {!isEmpty(moveTaskOrder) && (
          <>
            <h2>Move Task Order</h2>
            <dl>
              <dt>Origin Duty Station</dt>
              <dd>{get(moveTaskOrder, 'originDutyStation')}</dd>
              <dt>Destination Duty Station</dt>
              <dd>{get(moveTaskOrder, 'destinationDutyStation')}</dd>
              <dt>Pickup Address</dt>
              <dd>{JSON.stringify(get(moveTaskOrder, 'pickupAddress'))}</dd>
              <dt>Destination Address</dt>
              <dd>{JSON.stringify(get(moveTaskOrder, 'destinationAddress'))}</dd>
              <dt>Requested Pickup Date</dt>
              <dd>{get(moveTaskOrder, 'requestedPickupDate')}</dd>
              <dt>Customer Remarks</dt>
              <dd>{get(moveTaskOrder, 'remarks')}</dd>
              <dt>Service Items</dt>
              <dd>{JSON.stringify(get(moveTaskOrder, 'serviceItems'))}</dd>
              <dt>Status</dt>
              <dd>{get(moveTaskOrder, 'status')}</dd>
              <dt>Weight Entitlement</dt>
              {/*TODO do the entitlements belong on MTO? Was thinking yes based on the ticket, but since have a separate endpoint wasn't sure*/}
              <dd>{get(moveTaskOrder, 'entitlements.totalWeightSelf')}</dd>
              <dt>SIT Entitlement</dt>
              <dd>{get(moveTaskOrder, 'entitlements.storageInTransit')}</dd>
              <dt>NTS Entitlement</dt>
              <dd>{moveTaskOrderNonTemporaryStorage}</dd>
              <dt>POV Entitlement</dt>
              <dd>{moveTaskOrderPrivatelyOwnedVehicle}</dd>
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

const mapStateToProps = state => {
  //TODO hard coding mto for now
  const fakeMoveTaskOrderID = '5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
  const entitlements = get(state, 'entities.entitlements');
  const moveTaskOrder = selectMoveTaskOrder(state, fakeMoveTaskOrderID);
  const customer = get(state, 'entities.customer', {});
  return {
    entitlements: entitlements && Object.values(entitlements).length > 0 ? Object.values(entitlements)[0] : null,
    moveTaskOrder,
    customer: Object.values(customer)[0] || null,
  };
};

const mapDispatchToProps = {
  getEntitlements,
  updateMoveTaskOrderStatus,
  getCustomerInfo,
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(CustomerDetails);
